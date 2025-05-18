package genv

import (
	"context"

	"github.com/mrtc0/genv/diff"
)

const (
	unfetchedValuePlaceholder = "(value not fetched)"
)

// Diff compares the environment variables defined in the config with the
// environment variables in the dotenv map.
func Diff(ctx context.Context, cfg *Config, envMap map[string]string, nameOnly bool) (*diff.Diff, error) {
	if nameOnly {
		return DiffEnvName(ctx, cfg, envMap)
	}

	return DiffEnv(ctx, cfg, envMap)
}

// DiffEnvName compares the environment variables defined in the config with
// the environment variables in the dotenv map, but only the names of the
// environment variables are used to take the difference.
func DiffEnvName(ctx context.Context, cfg *Config, envMap map[string]string) (*diff.Diff, error) {
	definedEnv := make(map[string]string)
	for key := range cfg.Envs {
		definedEnv[key] = unfetchedValuePlaceholder
	}

	scrubbedEnvMap := make(map[string]string)
	for key := range envMap {
		scrubbedEnvMap[key] = unfetchedValuePlaceholder
	}

	return diffEnvMap(scrubbedEnvMap, definedEnv), nil
}

// DiffEnv compares the environment variables defined in the config with
// the environment variables in the dotenv map, including the values of the
// environment variables.
func DiffEnv(ctx context.Context, cfg *Config, envMap map[string]string) (*diff.Diff, error) {
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

	return diffEnvMap(envMap, fetched), nil
}

func diffEnvMap(old, new map[string]string) *diff.Diff {
	d := diff.DiffEnvMap(old, new)
	return &d
}
