package main

import (
	"context"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/sashabaranov/go-openai"
)

func main() {
	t := NewTerminal(os.Stdin, os.Stdout)
	err := do(context.Background(), t)
	if err != nil {
		t.WriteError(err.Error())
		os.Exit(1)
	}
}

type config struct {
	OpenAIKey SecretString `env:"OPENAI_API_KEY,notEmpty"`
	Seed      *int         `env:"SEED"`
}

type SecretString string

func (s SecretString) String() string {
	return "[REDACTED]"
}

func do(ctx context.Context, terminal *Terminal) error {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	tools := NewTools(terminal)

	client := openai.NewClient(string(cfg.OpenAIKey))

	ku := NewKupilot(tools, client, terminal, cfg.Seed)
	err := ku.Run(ctx)
	if err != nil {
		return fmt.Errorf("failed to run kupilot: %w", err)
	}
	return nil
}
