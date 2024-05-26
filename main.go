package main

import (
	"context"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/sashabaranov/go-openai"
)

func main() {
	err := do(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type config struct {
	OpenAIKey SecretString `env:"OPENAI_API_KEY,notEmpty"`
}

type SecretString string

func (s SecretString) String() string {
	return "[REDACTED]"
}

func do(ctx context.Context) error {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	tools := NewTools()

	client := openai.NewClient(string(cfg.OpenAIKey))

	ku := NewKupilot(tools, client)
	err := ku.Run(ctx)
	if err != nil {
		return fmt.Errorf("failed to run kupilot: %w", err)
	}
	return nil
}
