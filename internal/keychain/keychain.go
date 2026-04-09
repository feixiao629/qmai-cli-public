package keychain

// Keychain provides secure token storage
type Keychain interface {
	// Set stores a token for the given profile
	Set(profile, token string) error
	// Get retrieves the token for the given profile
	Get(profile string) (string, error)
	// Delete removes the token for the given profile
	Delete(profile string) error
}

const ServiceName = "qmai-cli"
