package main

import (
	"context"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/sashabaranov/go-openai"
)

func main() {
	t := &Terminal{
		in:      os.Stdin,
		out:     os.Stdout,
		noColor: false,
	}
	err := do(context.Background(), t)
	if err != nil {
		t.WriteError(err.Error())
		os.Exit(1)
	}
}

type config struct {
	OpenAIKey   SecretString `env:"OPENAI_API_KEY,notEmpty"`
	OpenAIModel string       `env:"OPENAI_MODEL" envDefault:"gpt-4o"`
	Seed        *int         `env:"SEED"`
}

// SecretString is a string that should not be printed, avoid accidental logging
type SecretString string

func (s SecretString) String() string {
	return "[REDACTED]"
}

func do(ctx context.Context, terminal *Terminal) error {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	tools := &Tools{
		skipExecutionConfirmation: false,
		terminal:                  terminal,
	}

	client := openai.NewClient(string(cfg.OpenAIKey))

	ku := NewKupilot(tools, client, terminal, cfg.Seed, cfg.OpenAIModel, false)
	err := ku.Run(ctx)
	if err != nil {
		return fmt.Errorf("failed to run kupilot: %w", err)
	}
	return nil
}
