# genv

genv is a tool for generating dotenv files by retrieving values from third-party services like AWS Secrets Manager

# Install

```bash
$ go install github.com/mrtc0/genv/cmd/genv@latest
```

# Supported Secrets Providers

- [x] AWS Secrets Manager
- [x] 1Password (via CLI or Service Account)

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

# Supported Providers

## AWS Secrets Manager

If you want to use AWS Secrets Manager as a secret provider, you can configure it as follows:

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

## 1Password

If you want to use 1Password as a secret provider, you can configure it as follows:

```yaml
secretProvider:
  1password:
    - id: my.1password.com
      auth:
        # Possible values for method are "cli" and "service-account"
        # If omitted, defaults to "cli"
        # When using "cli" method, genv will execute the 1Password CLI (`op`) command.
        # ref. https://developer.1password.com/docs/cli
        method: cli
        # account is optional when using "cli" method.
        # If omitted, the default account configured in the `op` CLI will be used.
        account: <your-account-id>
    - id: example.1password.com
      auth:
        # If you want to use Service Account authentication,
        # you must set the OP_SERVICE_ACCOUNT_TOKEN environment variable.
        # ref. https://developer.1password.com/docs/service-accounts
        method: service-account

envs:
  PASSWORD:
    secretRef:
      provider: my.1password.com
      # For 1Password provider configurations, the key must be in the format of a secret reference URI.
      #   op://<vault-name>/<item-name>/[section-name/]<field-name>
      #   e.g., op://my-vault/my-item/password
      # See details: https://developer.1password.com/docs/cli/secret-references
      key: "op://some-vault/some-item/field"
  API_KEY:
    secretRef:
      provider: example.1password.com
      key: "op://some-vault/some-item/field"
```

Currently, genv supports 1Password CLI (`op` command) authentication and Service Account authentication. Authentication using 1Password Connect is not supported.
