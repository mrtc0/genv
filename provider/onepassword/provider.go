package onepassword

import (
	"context"

	"github.com/mrtc0/genv/provider"
	"github.com/mrtc0/genv/provider/onepassword/op"
	"github.com/mrtc0/genv/provider/onepassword/sdk"
)

type OnePasswordAuthMethod string

const (
	OnePasswordAuthMethodCLI            OnePasswordAuthMethod = "cli"
	OnePasswordAuthMethodServiceAccount OnePasswordAuthMethod = "service-account"
)

var _ provider.Provider = &Provider{}

type Provider struct {
	account string
	method  OnePasswordAuthMethod
}

type ProviderOption func(*Provider)

func NewProvider(opts ...ProviderOption) provider.Provider {
	p := &Provider{}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *Provider) NewClient(ctx context.Context) (provider.SecretClient, error) {
	switch p.method {
	case OnePasswordAuthMethodServiceAccount:
		return sdk.NewOnePasswordClient()
	case OnePasswordAuthMethodCLI:
		return op.NewOPClient(p.account), nil
	default:
		return op.NewOPClient(p.account), nil
	}
}

func WithAuthMethod(method OnePasswordAuthMethod) ProviderOption {
	return func(p *Provider) {
		p.method = method
	}
}

func WithAccount(account string) ProviderOption {
	return func(p *Provider) {
		p.account = account
	}
}
