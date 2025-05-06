# genv

genv is a tool for generating dotenv files by retrieving values from third-party services like AWS Secrets Manager

# Install

```bash
$ go install github.com/mrtc0/genv/cmd/genv@latest
```

# Supported Secrets Providers

- [x] AWS Secrets Manager

# Guide

## 1. Create a genv setting file

Now, create a setting file that describes the third-party secrets provider and the definition of the environment variables to be generated.
genv reads a file named `genv.yaml` by default.

```yaml
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

## 2. Generate .env file

Run the following command to generate a `.env` file.

```bash
$ genv gen

$ cat .env
APP_ENV=development
API_KEY=dummy-api-key
DB_PASSWORD=dummy-db-password
```
