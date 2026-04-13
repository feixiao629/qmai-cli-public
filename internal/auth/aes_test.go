package auth

import "testing"

func TestAESEncryptDecrypt_DocExample(t *testing.T) {
	// Fake credential material used to keep the example deterministic.
	openKey := "fake-open-key-for-test-only-0002"
	openId := "fake-open-id-for-test-only-0002"
	plaintext := "13800000000"
	expectedCiphertext := "IjvTj7Umt0Kq+s/b1jd4iQ=="

	encrypted, err := AESEncrypt(plaintext, openId, openKey)
	if err != nil {
		t.Fatalf("AESEncrypt: %v", err)
	}
	if encrypted != expectedCiphertext {
		t.Errorf("AESEncrypt() = %q, want %q", encrypted, expectedCiphertext)
	}

	// Verify decryption roundtrip
	decrypted, err := AESDecrypt(expectedCiphertext, openId, openKey)
	if err != nil {
		t.Fatalf("AESDecrypt: %v", err)
	}
	if decrypted != plaintext {
		t.Errorf("AESDecrypt() = %q, want %q", decrypted, plaintext)
	}
}

func TestAESEncryptDecrypt_Roundtrip(t *testing.T) {
	openId := "fake-open-id-for-test-only-0003"
	openKey := "fake-open-key-for-test-only-0003"
	plaintext := "hello world test"

	encrypted, err := AESEncrypt(plaintext, openId, openKey)
	if err != nil {
		t.Fatalf("AESEncrypt: %v", err)
	}

	decrypted, err := AESDecrypt(encrypted, openId, openKey)
	if err != nil {
		t.Fatalf("AESDecrypt: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("roundtrip: got %q, want %q", decrypted, plaintext)
	}
}

func TestDeriveAESKey(t *testing.T) {
	key := deriveAESKey("test")
	if len(key) != 32 {
		t.Errorf("key length = %d, want 32", len(key))
	}
}

func TestDeriveAESIV(t *testing.T) {
	iv := deriveAESIV("test")
	if len(iv) != 16 {
		t.Errorf("iv length = %d, want 16", len(iv))
	}
}
