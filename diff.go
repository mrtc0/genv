package genv

import (
	"context"

	"github.com/mrtc0/genv/diff"
)

func Diff(ctx context.Context, cfg *Config, dotenv map[string]string, nameOnly bool) (*diff.Diff, error) {
	if nameOnly {
		return DiffEnvName(ctx, cfg, dotenv)
	}

	return DiffEnv(ctx, cfg, dotenv)
}

func DiffEnvName(ctx context.Context, cfg *Config, dotenv map[string]string) (*diff.Diff, error) {
	definedEnv := make(map[string]string)
	for key := range cfg.Envs {
		definedEnv[key] = ""
	}

	scrubbed := make(map[string]string)
	for key := range dotenv {
		scrubbed[key] = ""
	}

	diff := diff.DiffEnvMap(definedEnv, scrubbed)
	return &diff, nil
}

func DiffEnv(ctx context.Context, cfg *Config, dotenv map[string]string) (*diff.Diff, error) {
	generator, err := NewDotenvGenerator(ctx, DotenvGeneratorConfig{
		Config: cfg,
	})
	if err != nil {
		return nil, err
	}

	fetched, err := generator.FetchSecrets(ctx)
	if err != nil {
		return nil, err
	}

	diff := diff.DiffEnvMap(dotenv, fetched)
	return &diff, nil
}

func extractKeys(envMap map[string]string) map[string]string {
	keys := make(map[string]string, len(envMap))
	for k := range envMap {
		keys[k] = ""
	}
	return keys
}

func diffEnvMap(a, b map[string]string) *diff.Diff {
	d := diff.DiffEnvMap(a, b)
	return &d
}
