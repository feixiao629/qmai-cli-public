package keychain

import "testing"

func TestMemoryKeychain_SetGetDelete(t *testing.T) {
	kc := NewMemoryKeychain()

	// Set a token
	if err := kc.Set("prod", "token-abc"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get the token back
	tok, err := kc.Get("prod")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if tok != "token-abc" {
		t.Errorf("Get = %q, want %q", tok, "token-abc")
	}

	// Overwrite the token
	if err := kc.Set("prod", "token-xyz"); err != nil {
		t.Fatalf("Set (overwrite) failed: %v", err)
	}
	tok, err = kc.Get("prod")
	if err != nil {
		t.Fatalf("Get after overwrite failed: %v", err)
	}
	if tok != "token-xyz" {
		t.Errorf("Get after overwrite = %q, want %q", tok, "token-xyz")
	}

	// Delete the token
	if err := kc.Delete("prod"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	tok, err = kc.Get("prod")
	if err != nil {
		t.Fatalf("Get after delete failed: %v", err)
	}
	if tok != "" {
		t.Errorf("Get after delete = %q, want empty string", tok)
	}
}

func TestMemoryKeychain_GetNonExistent(t *testing.T) {
	kc := NewMemoryKeychain()

	tok, err := kc.Get("nonexistent")
	if err != nil {
		t.Fatalf("Get non-existent should not error, got: %v", err)
	}
	if tok != "" {
		t.Errorf("Get non-existent = %q, want empty string", tok)
	}
}

func TestMemoryKeychain_MultipleProfiles(t *testing.T) {
	kc := NewMemoryKeychain()

	_ = kc.Set("dev", "dev-token")
	_ = kc.Set("prod", "prod-token")

	devTok, _ := kc.Get("dev")
	prodTok, _ := kc.Get("prod")

	if devTok != "dev-token" {
		t.Errorf("dev token = %q, want %q", devTok, "dev-token")
	}
	if prodTok != "prod-token" {
		t.Errorf("prod token = %q, want %q", prodTok, "prod-token")
	}

	// Delete one, the other should remain
	_ = kc.Delete("dev")
	devTok, _ = kc.Get("dev")
	prodTok, _ = kc.Get("prod")
	if devTok != "" {
		t.Errorf("dev token after delete = %q, want empty", devTok)
	}
	if prodTok != "prod-token" {
		t.Errorf("prod token after deleting dev = %q, want %q", prodTok, "prod-token")
	}
}

func TestMemoryKeychain_DeleteNonExistent(t *testing.T) {
	kc := NewMemoryKeychain()

	// Delete non-existent key should not error
	if err := kc.Delete("nonexistent"); err != nil {
		t.Fatalf("Delete non-existent should not error, got: %v", err)
	}
}

func TestMemoryKeychain_ImplementsInterface(t *testing.T) {
	var _ Keychain = NewMemoryKeychain()
}
