# genv

genv is a tool for generating dotenv files by retrieving values from third-party services like AWS Secrets Manager

# Install

```bash
$ go install github.com/mrtc0/genv/cmd/genv@latest
```

# Supported Secrets Providers

- [x] AWS Secrets Manager

# Usage

The following is given by running `genv -h`:

```bash
Usage:
  genv [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  gen         Generate .env file
  help        Help about any command
  outdated    Show outdated envs in the dotenv file.
  run         Run a command with environment variables from .env file

Flags:
  -h, --help   help for genv

Use "genv [command] --help" for more information about a command.
```


# Getting Started

## Generate .env file

Create a YAML file that defines the third-party secret provider and environment variables.

> **NOTE**: genv looks for .genv.yaml by default, but you can specify a different file with the --config option.

```yaml
# .genv.yaml
secretProvider:
  aws:
    - id: example-account
      service: SecretsManager
      region: us-east-1
      auth:
        # If you want to use a specific AWS profile, specify it here
        profile: default
        # If you want to use a specific Shared Credentials File
        # sharedCredentialsFiles: ["/path/to/credentials"]
        # If you want to use a specific Shared Configuration File
        # sharedConfigFiles: ["/path/to/config"]
    - id: another-account
      service: SecretsManager
      region: us-west-2

envs:
  APP_ENV:
    value: "development"
  API_KEY:
    secretRef:
      provider: example-account
      key: apikey
  DB_PASSWORD:
    secretRef:
      provider: another-account
      key: db-credentials
      property: ".password"
```

In this example, we define an environment variable `APP_ENV=development`, while the environment variables `API_KEY` and `DB_PASSWORD` are retrieved from AWS Secrets Manager.
When the value stored in Secrets Manager is in JSON format, you can specify a property using the `property` field.

Run `genv gen`, a `.env` file is generated based on the above configuration.

```shell
$ genv gen

$ cat .env
APP_ENV=development
API_KEY=this-is-a-secret
DB_PASSWORD=password
```

## Detect outdated environment variable definitions

The `genv outdated` command compares the environment variables defined in genv.yaml with the environment variables in the .env file.

```shell
$ genv outdated
~ DB_PASSWORD  =  "password" => "new-password"

Error: outdated envs found
exit status 1
```

If you want to ignore changes in environment variable values, use the `--ignore-value` option. With this option, values won't be retrieved from authentication providers.
This is useful when you want to avoid accessing credential providers.

```shell
$ genv outdated --ignore-value
+ DB_HOST      =  "(value not retrieved)"
- DB_PASSWORD  =  "(value not retrieved)"

Error: outdated envs found
exit status 1
```
