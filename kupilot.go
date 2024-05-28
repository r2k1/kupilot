package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/glamour"
	"github.com/sashabaranov/go-openai"
)

type Kupilot struct {
	tools  *Tools
	openai *openai.Client
	msgs   []openai.ChatCompletionMessage
}

var SysMessage = openai.ChatCompletionMessage{
	Role:    openai.ChatMessageRoleSystem,
	Content: "You are a kubernetes expert, your job is to help users with their kubernetes questions. You can ask for more information if needed. You can also ask for clarification if you are unsure about something. You have read access to the kubernetes cluster. Be concise. Output of every function call is printed to the user, don't repeat it. If output is truncated you can modify the script to limit the scope. Respond in plaint text, don't use markdown",
}

func NewKupilot(tools *Tools, aiclient *openai.Client) *Kupilot {
	return &Kupilot{
		tools:  tools,
		openai: aiclient,
		msgs:   []openai.ChatCompletionMessage{SysMessage},
	}
}

func (k *Kupilot) Run(ctx context.Context) error {
	fmt.Println("How can I help you?")
	for {
		userInput, err := read()
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
		Model:    openai.GPT4o,
		Messages: k.msgs,
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
		return fmt.Errorf("failed to create chat completion stream: %w", err)
	}

	agentMsg := resp.Choices[0].Message

	out, err := glamour.Render(agentMsg.Content, "dark")
	if err != nil {
		fmt.Println(agentMsg.Content)
	} else {
		fmt.Print(out)
	}

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

func read() (string, error) {
	fmt.Printf("\n> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read user input: %w", err)
	}
	userInput := scanner.Text()
	if userInput == "exit" || userInput == "quit" {
		os.Exit(0)
	}
	return userInput, nil
}
