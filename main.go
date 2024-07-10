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
	OpenAIKey         SecretString `env:"OPENAI_API_KEY,notEmpty"`
	OpenAIModel       string       `env:"OPENAI_MODEL" envDefault:"gpt-4o"`
	AzureOpenAIAPIURL string       `env:"AZURE_OPENAI_API_URL"`
	Seed              *int         `env:"SEED"`
	NoColor           bool         `env:"NO_COLOR"`
}

func (c config) NewAIConfig() openai.ClientConfig {
	if c.AzureOpenAIAPIURL != "" {
		return openai.DefaultAzureConfig(string(c.OpenAIKey), c.AzureOpenAIAPIURL)
	}
	return openai.DefaultConfig(string(c.OpenAIKey))
}

// SecretString is a string that should not be printed, avoid accidental logging
type SecretString string

func (s SecretString) String() string {
	return "[REDACTED]"
}

func do(ctx context.Context) error {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	terminal := &Terminal{
		in:      os.Stdin,
		out:     os.Stdout,
		noColor: cfg.NoColor,
	}
	tools := &Tools{
		skipExecutionConfirmation: false,
		terminal:                  terminal,
	}

	aiClient := openai.NewClientWithConfig(cfg.NewAIConfig())

	ku := &Kupilot{
		tools:    tools,
		openai:   aiClient,
		terminal: terminal,
		seed:     cfg.Seed,
		model:    cfg.OpenAIModel,
		noColor:  cfg.NoColor,
	}

	err := ku.Run(ctx)
	if err != nil {
		return fmt.Errorf("failed to run kupilot: %w", err)
	}
	return nil
}
