package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type Tools struct {
	skipExecutionConfirmation bool
	terminal                  *Terminal
}

func NewTools(terminal *Terminal) *Tools {
	return &Tools{
		skipExecutionConfirmation: false,
		terminal:                  terminal,
	}
}

var ToolsSchema = schema()

func schema() []byte {
	return []byte(`{"properties":{"script":{"type":"string","description":"The bash script to run"}},"additionalProperties":false,"type":"object","required":["script"]}`)
}

func (t *Tools) Call(requests []openai.ToolCall) ([]openai.ChatCompletionMessage, error) {
	responses := make([]openai.ChatCompletionMessage, 0, len(requests))
	for _, request := range requests {
		response := t.Exec(request)
		responses = append(responses, response)
	}
	return responses, nil
}

type ScriptParams struct {
	Script string `json:"script" jsonschema:"description=The bash script to run"`
}

func (t *Tools) execScript(params ScriptParams) (string, error) {
	if err := t.confirmExecution(params); err != nil {
		return "", err
	}

	cmd := exec.Command("bash", "-c", params.Script)
	output, err := cmd.CombinedOutput()
	t.terminal.Write(fmt.Sprintf("Output:\n%s\n", output))
	if err != nil {
		return string(output), err
	}
	if len(output) > 10000 {
		return string(output[:10000]) + "\nOutput truncated, full output is printed for the user", nil
	}
	return string(output), nil
}

func (t *Tools) confirmExecution(params ScriptParams) error {
	if t.skipExecutionConfirmation {
		return nil
	}
	t.terminal.WriteWarning(fmt.Sprintf("About to execute:\n%s\nDo you want to proceed? (y/n) :", params.Script))
	text, err := t.terminal.Read()
	if err != nil {
		return err
	}
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Trim(text, " ")
	if strings.ToLower(text) != "y" && strings.ToLower(text) != "yes" {
		return fmt.Errorf("execution cancelled by user")
	}
	return nil
}

func (t *Tools) Exec(req openai.ToolCall) openai.ChatCompletionMessage {
	response := openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleTool,
		ToolCallID: req.ID,
		Name:       req.Function.Name,
	}

	if req.Type != openai.ToolTypeFunction {
		t.terminal.WriteError(fmt.Sprintf("Unknown tool type: %s\n", req.Type))
		return response
	}

	switch req.Function.Name {
	case "script":
		var params ScriptParams
		err := json.Unmarshal([]byte(req.Function.Arguments), &params)
		if err != nil {
			response.Content = err.Error()
			return response
		}
		content, err := t.execScript(params)
		if err != nil {
			response.Content = err.Error()
			t.terminal.WriteError(err.Error())
			return response
		}
		response.Content = content
		return response
	}

	errMsg := fmt.Sprintf("Invalid tool: %s\n", req.Function.Name)
	t.terminal.WriteError(errMsg)
	response.Content = errMsg
	response.Name = "error"
	return response
}
