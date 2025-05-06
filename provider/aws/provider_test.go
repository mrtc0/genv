package aws_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/mrtc0/genv/provider/aws"
	"github.com/mrtc0/genv/provider/aws/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfig(t *testing.T) {
	mocks.InitSessionTestEnv(t)

	type arrange struct {
		SharedConfigFile     string
		SharedCredentialFile string
		SSOSessionName       string
	}

	type expect struct {
		region      string
		credentials awssdk.Credentials
	}

	testCases := map[string]struct {
		arrange           arrange
		awsProviderConfig *aws.AwsProviderConfig
		expect            expect
	}{
		"with SharedCredentialFile": {
			arrange: arrange{
				SharedConfigFile: `
[default]
region = us-east-1
`,
				SharedCredentialFile: `
[default]
aws_access_key_id=AKIAIOSFODNN7EXAMPLE
aws_secret_access_key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
`,
			},
			awsProviderConfig: &aws.AwsProviderConfig{
				ID:      "aws",
				Service: aws.AWSSecretsManager,
				Region:  "ap-northeast-1",
				Auth:    aws.AwsAuth{},
			},
			expect: expect{
				region: "ap-northeast-1",
				credentials: awssdk.Credentials{
					AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
					SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				},
			},
		},
		"with profile": {
			arrange: arrange{
				SharedConfigFile: `
[profile user1]
region = us-east-1
`,
				SharedCredentialFile: `
[user1]
aws_access_key_id=AKIAIOSFODNN7EXAMPLE
aws_secret_access_key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
`,
			},
			awsProviderConfig: &aws.AwsProviderConfig{
				ID:      "aws",
				Service: aws.AWSSecretsManager,
				Region:  "ap-northeast-1",
				Auth: aws.AwsAuth{
					Profile: "user1",
				},
			},
			expect: expect{
				region: "ap-northeast-1",
				credentials: awssdk.Credentials{
					AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
					SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				},
			},
		},
		"with profile and sso-session": {
			arrange: arrange{
				SharedConfigFile: fmt.Sprintf(`
[profile dev]
sso_session = my-sso
sso_account_id = %s
sso_role_name = %s

[sso-session my-sso]
sso_region = ap-northeast-1
sso_start_url = https://my-sso-portal.awsapps.com/start
`, mocks.StsGetRoleCredentialsAccountId, mocks.StsGetRoleCredentialsRoleName),
				SSOSessionName: "my-sso",
			},
			awsProviderConfig: &aws.AwsProviderConfig{
				ID:      "aws",
				Service: aws.AWSSecretsManager,
				Region:  "ap-northeast-1",
				Auth: aws.AwsAuth{
					Profile: "dev",
				},
			},
			expect: expect{
				region: "ap-northeast-1",
				credentials: awssdk.Credentials{
					AccessKeyID:     mocks.StsGetRoleCredentialsAccessKeyId,
					SecretAccessKey: mocks.StsGetRoleCredentialsSecretAccessKey,
					SessionToken:    mocks.StsGetRoleCredentialsSessionToken,
					AccountID:       mocks.StsGetRoleCredentialsAccountId,
					CanExpire:       true,
				},
			},
		},
	}

	for name, tt := range testCases {
		tt := tt

		t.Run(name, func(t *testing.T) {
			ts := mocks.MockAwsApiServer(t, []*mocks.MockEndpoint{
				mocks.MockStsGetCallerIdentityValidEndpoint,
				mocks.MockStsGetRoleCredentialsValidEndpoint,
			})
			defer ts.Close()

			tt.awsProviderConfig.Endpoint = ts.URL

			if tt.arrange.SSOSessionName != "" {
				require.NoError(t, mocks.SsoTestSetup(t, tt.arrange.SSOSessionName))
			}

			if tt.arrange.SharedConfigFile != "" {
				file, err := os.CreateTemp("", "aws-sdk-go-base-shared-configuration-file")
				require.NoError(t, err, "failed to create temporary shared configuration file")

				defer os.Remove(file.Name())

				err = os.WriteFile(file.Name(), []byte(tt.arrange.SharedConfigFile), 0600)
				require.NoError(t, err, "failed to write temporary shared configuration file")

				tt.awsProviderConfig.Auth.SharedConfigFiles = []string{file.Name()}
			}

			if tt.arrange.SharedCredentialFile != "" {
				file, err := os.CreateTemp("", "aws-sdk-go-base-shared-credentials-file")
				require.NoError(t, err, "failed to create temporary shared credentials file")

				defer os.Remove(file.Name())

				err = os.WriteFile(file.Name(), []byte(tt.arrange.SharedCredentialFile), 0600)
				require.NoError(t, err, "failed to write temporary shared credentials file")

				tt.awsProviderConfig.Auth.SharedCredentialsFiles = []string{file.Name()}
			}

			ctx := context.Background()

			cfg, err := aws.GetAWSConfig(ctx, tt.awsProviderConfig)
			assert.NoError(t, err)
			assert.Equal(t, tt.expect.region, cfg.Region)

			creds, err := cfg.Credentials.Retrieve(ctx)
			assert.NoError(t, err)

			if diff := cmp.Diff(creds, tt.expect.credentials, cmpopts.IgnoreFields(awssdk.Credentials{}, "Expires", "Source")); diff != "" {
				assert.Fail(t, "unexpected credentials", "(- got, + expected)\n%s", diff)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	providerConfig := &aws.AwsProviderConfig{
		ID:      "aws",
		Service: aws.AWSSecretsManager,
		Region:  "ap-northeast-1",
		Auth:    aws.AwsAuth{},
	}

	provider := aws.NewProvider(providerConfig)
	client, err := provider.NewClient(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
