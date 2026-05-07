package genv

import (
	"fmt"
	"os"

	"github.com/mrtc0/genv/provider/onepassword"
	"gopkg.in/yaml.v3"
)

type Config struct {
	SecretProvider SecretProvider      `yaml:"secretProvider,omitempty"`
	Envs           map[string]EnvValue `yaml:"envs,omitempty"`
}

type SecretProvider struct {
	Aws         []AwsProvider         `yaml:"aws,omitempty"`
	GoogleCloud []GoogleCloudProvider `yaml:"googleCloud,omitempty"`
	OnePassword []OnePasswordProvider `yaml:"1password,omitempty"`
	Exec        []ExecProvider        `yaml:"exec,omitempty"`
}

type AwsProvider struct {
	ID      string  `yaml:"id"`
	Service string  `yaml:"service"`
	Region  string  `yaml:"region,omitempty"`
	Auth    AwsAuth `yaml:"auth,omitempty"`
}

type AwsAuth struct {
	Profile                string   `yaml:"profile,omitempty"`
	SharedCredentialsFiles []string `yaml:"sharedCredentialsFiles,omitempty"`
	SharedConfigFiles      []string `yaml:"sharedConfigFiles,omitempty"`
}

type GoogleCloudProvider struct {
	ID        string `yaml:"id"`
	Service   string `yaml:"service"`
	ProjectID string `yaml:"projectID"`
	Location  string `yaml:"location,omitempty"`
}

type OnePasswordProvider struct {
	ID   string          `yaml:"id"`
	Auth OnePasswordAuth `yaml:"auth,omitempty"`
}

// OnePasswordAuth represents the authentication configuration for 1Password
type OnePasswordAuth struct {
	// The authentication method to use for 1Password
	// Possible values are "cli" and "service-account"
	// If omitted, defaults to "cli"
	Method onepassword.OnePasswordAuthMethod `yaml:"method"`
	// The account to use for 1Password (only applicable when Method is CLI)
	Account string `yaml:"account,omitempty"`
}

// ExecCommand supports two YAML forms for specifying a command:
//
//	String form:   command: "vault kv get -format=json secret/myapp | jq .data"
//	Sequence form: command: ["vault", "kv", "get", "-format=json", "secret/myapp"]
//
// The string form is passed to "sh -c" so that shell features such as pipes
// and redirections work. The sequence form is executed directly via execve,
// which avoids shell interpretation and is safer when the arguments are known
// at configuration time.
type ExecCommand struct {
	Args []string
}

// UnmarshalYAML allows ExecCommand to accept either a plain string or a
// sequence of strings in YAML.
//
//   - Scalar (string): wrapped as ["sh", "-c", <value>] so the shell handles
//     pipes, redirections, etc.
//   - Sequence: decoded as-is and executed directly without a shell.
func (c *ExecCommand) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		c.Args = []string{"sh", "-c", value.Value}
	case yaml.SequenceNode:
		return value.Decode(&c.Args)
	default:
		return fmt.Errorf("command must be a string or sequence")
	}
	return nil
}

type ExecProvider struct {
	ID      string      `yaml:"id"`
	Command ExecCommand `yaml:"command"`
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
