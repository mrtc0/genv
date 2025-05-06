package genv

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SecretProvider SecretProvider      `yaml:"secretProvider,omitempty"`
	Envs           map[string]EnvValue `yaml:"envs,omitempty"`
}

type SecretProvider struct {
	Aws []AwsProvider `yaml:"aws,omitempty"`
}

type AwsProvider struct {
	ID      string  `yaml:"id"`
	Service string  `yaml:"service"`
	Region  string  `yaml:"region,omitempty"`
	Auth    AwsAuth `yaml:"auth,omitempty"`
}

type AwsAuth struct {
	Profile                string   `yaml:"profile,omitempty"`
	Region                 string   `yaml:"region,omitempty"`
	SharedCredentialsFiles []string `yaml:"sharedCredentialsFiles,omitempty"`
	SharedConfigFiles      []string `yaml:"sharedConfigFiles,omitempty"`
}

type EnvValue struct {
	Value     string     `yaml:"value,omitempty"`
	SecretRef *SecretRef `yaml:"secretRef,omitempty"`
}

type SecretRef struct {
	Provider string `yaml:"provider,omitempty"`
	Key      string `yaml:"key,omitempty"`
	Property string `yaml:"property,omitempty"`
}

func LoadConfig(filePath string) (*Config, error) {
	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(f, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
