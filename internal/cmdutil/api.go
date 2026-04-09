package cmdutil

import (
	"fmt"
	"net/http"

	"github.com/madaima/qmai-cli/internal/client"
	"github.com/madaima/qmai-cli/internal/core"
	"github.com/madaima/qmai-cli/internal/keychain"
)

// ApiClient creates an authenticated open platform API client.
func (f *Factory) ApiClient() (*client.Client, error) {
	cfg, err := f.Config()
	if err != nil {
		return nil, err
	}

	profile := cfg.Profiles[f.EffectiveProfile()]
	if profile == nil {
		return nil, fmt.Errorf("%w: run 'qmai auth login'", core.ErrNotAuthenticated)
	}

	if profile.OpenId == "" || profile.GrantCode == "" {
		return nil, fmt.Errorf("%w: missing openId/grantCode, run 'qmai auth login'", core.ErrNotAuthenticated)
	}

	// Get openKey from keychain
	kc := f.Keychain
	if kc == nil {
		kc = keychain.NewOSKeychain()
	}

	openKey, err := kc.Get(f.EffectiveProfile())
	if err != nil || openKey == "" {
		return nil, fmt.Errorf("%w: missing openKey, run 'qmai auth login'", core.ErrNotAuthenticated)
	}

	httpClient := &http.Client{}

	return client.NewClient(cfg.BaseURL(), profile.OpenId, profile.GrantCode, openKey, httpClient, f.Debug), nil
}
