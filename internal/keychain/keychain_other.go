//go:build !darwin

package keychain

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

type osKeychain struct{}

func NewOSKeychain() Keychain {
	return &osKeychain{}
}

func (k *osKeychain) Set(profile, token string) error {
	if err := keyring.Set(ServiceName, profile, token); err != nil {
		return fmt.Errorf("keychain set failed: %w", err)
	}
	return nil
}

func (k *osKeychain) Get(profile string) (string, error) {
	token, err := keyring.Get(ServiceName, profile)
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", nil
		}
		return "", fmt.Errorf("keychain get failed: %w", err)
	}
	return token, nil
}

func (k *osKeychain) Delete(profile string) error {
	if err := keyring.Delete(ServiceName, profile); err != nil {
		if err == keyring.ErrNotFound {
			return nil
		}
		return fmt.Errorf("keychain delete failed: %w", err)
	}
	return nil
}
