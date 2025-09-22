package onepassword

import (
	"context"

	"github.com/mrtc0/genv/provider"
	"github.com/mrtc0/genv/provider/onepassword/op"
)

var _ provider.Provider = &Provider{}

type Provider struct {
	account string
}

func NewProvider(account string) provider.Provider {
	return &Provider{account: account}
}

func (p *Provider) NewClient(ctx context.Context) (provider.SecretClient, error) {
	return op.NewOPClient(p.account), nil
}
