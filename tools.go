package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type Tools struct {
	skipExecutionConfirmation bool
}

func NewTools() *Tools {
	return &Tools{
		skipExecutionConfirmation: false,
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
	fmt.Printf("Output:\n%s\n", output)
	if err != nil {
		return string(output), err
	}
	fmt.Printf("Output:\n%s\n", output)
	if len(output) > 10000 {
		return string(output[:10000]) + "\nOutput truncated, full output is printed for the user", nil
	}
	return string(output), nil
}

func (t *Tools) confirmExecution(params ScriptParams) error {
	if t.skipExecutionConfirmation {
		return nil
	}
	reader := bufio.NewReader(os.Stdin)
	printLnCyan(fmt.Sprintf("About to execute:\n%s\nDo you want to proceed? (y/n):", params.Script))
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	if strings.ToLower(text) != "y" {
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
		fmt.Printf("Unknown tool type: %s\n", req.Type)
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
			printlnRed(err)
			return response
		}
		response.Content = content
		return response
	}

	fmt.Printf("Unknown tool: %s\n", req.Function.Name)
	response.Content = "Unknown tool"
	return response
}

func printlnRed(input ...any) {
	redColor := "\033[31m"
	resetColor := "\033[0m"
	out := []any{redColor}
	out = append(out, input...)
	out = append(out, resetColor)
	fmt.Println(out...)
}

func printlnGrey(input ...any) {
	redColor := "\033[90m"
	resetColor := "\033[0m"
	out := []any{redColor}
	out = append(out, input...)
	out = append(out, resetColor)
	fmt.Println(out...)
}

func printLnCyan(input ...any) {
	cyanColor := "\033[36m"
	resetColor := "\033[0m"
	out := []any{cyanColor}
	out = append(out, input...)
	out = append(out, resetColor)
	fmt.Println(out...)

}
