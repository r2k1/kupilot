package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/glamour"
	"github.com/sashabaranov/go-openai"
)

type Kupilot struct {
	tools    *Tools
	openai   OpenAIClient
	msgs     []openai.ChatCompletionMessage
	terminal *Terminal
	seed     *int
	model    string
	noColor  bool
}

//go:generate moq -out open_ai_client_moq_test.go . OpenAIClient
type OpenAIClient interface {
	CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

var SysMessage = openai.ChatCompletionMessage{
	Role: openai.ChatMessageRoleSystem,
	Content: `You are a kubernetes expert, your job is to help users with their kubernetes questions.
You can ask for more information if needed. You can also ask for clarification if you are unsure about something.
You have read access to the kubernetes cluster. Be concise. Output of every function call is printed to the user, don't repeat it. 
If output is truncated you can modify the script to limit the scope of the output.`,
}

func NewKupilot(tools *Tools, aiclient OpenAIClient, terminal *Terminal, seed *int, model string, noColor bool) *Kupilot {
	return &Kupilot{
		tools:    tools,
		openai:   aiclient,
		msgs:     []openai.ChatCompletionMessage{SysMessage},
		terminal: terminal,
		seed:     seed,
		model:    model,
		noColor:  noColor,
	}
}

func (k *Kupilot) Run(ctx context.Context) error {
	if k.seed != nil {
		k.terminal.WriteInfo(fmt.Sprintf("Using seed: %d\n", *k.seed))
	}
	k.terminal.Write("Hello! Kupilot here, how can I help you?\n")
	for {
		userInput, err := k.terminal.Read()
		if err != nil {
			return err
		}

		k.msgs = append(k.msgs, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: userInput,
		})

		if err = k.askGPT(ctx); err != nil {
			return fmt.Errorf("failed to ask GPT: %w", err)
		}
	}
}

func (k *Kupilot) askGPT(ctx context.Context) error {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	_ = s.Color("cyan")
	s.Start()
	resp, err := k.openai.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    k.model,
		Messages: k.msgs,
		Seed:     k.seed,
		Tools: []openai.Tool{
			{
				Type: openai.ToolTypeFunction,
				Function: &openai.FunctionDefinition{
					Name:        "script",
					Description: "Run a bash script, full output is printed for the user but can be truncated for the agent",
					Parameters:  json.RawMessage(ToolsSchema),
				},
			},
		},
	})
	s.Stop()
	if err != nil {
		return fmt.Errorf("failed to submit openai request: %w", err)
	}

	agentMsg := resp.Choices[0].Message

	k.writeMessage(agentMsg)

	k.msgs = append(k.msgs, agentMsg)

	if agentMsg.ToolCalls == nil {
		return nil
	}

	toolMsgs, err := k.tools.Call(agentMsg.ToolCalls)
	if err != nil {
		return fmt.Errorf("failed to call tool: %w", err)
	}
	k.msgs = append(k.msgs, toolMsgs...)
	return k.askGPT(ctx)
}

func (k *Kupilot) writeMessage(agentMsg openai.ChatCompletionMessage) {
	if !k.noColor {
		out, err := glamour.Render(agentMsg.Content, "dark")
		if err == nil {
			k.terminal.Write(out)
			return
		}
	}
	k.terminal.Write(agentMsg.Content)
}
