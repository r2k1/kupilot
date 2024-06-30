package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKupilot(t *testing.T) {
	ctx := context.Background()
	out := &bytes.Buffer{}
	terminal := &Terminal{
		in: &Input{
			msgs: []string{"test\n", "y\n"},
		},
		out:     out,
		noColor: true,
	}
	tools := &Tools{
		skipExecutionConfirmation: false,
		terminal:                  terminal,
	}
	calls := -1

	// delete test/test.txt
	err := os.WriteFile("test/test.txt", []byte("test content\n"), 0644)
	require.NoError(t, err)
	ai := &OpenAIClientMock{
		CreateChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
			assert.Equal(t, "GPT4o", req.Model)
			calls++
			switch calls {
			case 0:
				return openai.ChatCompletionResponse{
					Choices: []openai.ChatCompletionChoice{
						{
							Message: openai.ChatCompletionMessage{
								ToolCalls: []openai.ToolCall{
									{
										ID:   "1",
										Type: openai.ToolTypeFunction,
										Function: openai.FunctionCall{
											Name:      "script",
											Arguments: `{"script": "echo 'test content' > test/test.txt"}`,
										},
									},
								},
							},
						},
					},
				}, nil
			case 1, 2:
				return openai.ChatCompletionResponse{
					Choices: []openai.ChatCompletionChoice{
						{
							Message: openai.ChatCompletionMessage{
								Content: "Hello from GPT",
							},
						},
					},
				}, nil
			default:
				return openai.ChatCompletionResponse{}, errors.New("unexpected call")
			}
		},
	}
	ku := NewKupilot(tools, ai, terminal, nil, "GPT4o", true)
	err = ku.Run(ctx)
	assert.ErrorIs(t, err, doneErr)
	assert.Equal(t, "Hello! Kupilot here, how can I help you?\n\n> About to execute:\necho 'test content' > test/test.txt\nDo you want to proceed? (y/n) :\n> Output:\n\nHello from GPT\n> ", out.String())
	content, err := os.ReadFile("test/test.txt")
	require.NoError(t, err)
	assert.Equal(t, "test content\n", string(content))
}

type Input struct {
	msgs  []string
	index int
}

var doneErr = errors.New("DONE")

func (i *Input) Read(p []byte) (int, error) {
	if i.index >= len(i.msgs) {
		return 0, doneErr
	}
	n := copy(p, i.msgs[i.index])
	i.index++
	return n, nil
}
