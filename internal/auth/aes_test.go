package auth

import "testing"

func TestAESEncryptDecrypt_DocExample(t *testing.T) {
	// Test data from the open platform documentation.
	// aesKey = openKey (input to MD5 for key derivation)
	// ivSeed = openId (input to MD5 for IV derivation)
	openKey := "4g7Ku9WE5T1ba84f55dece41d93001dbb722f90efef6EOTVI3"
	openId := "bb9203417a6c7da10d65e8eef01e61ae"
	plaintext := "13815331845"
	expectedCiphertext := "Rv/q+dGCwn3t+Fi0y+SYuA=="

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
	openId := "d14c1559e87b747d577c834b275a4310"
	openKey := "LyvrkvkxRkG2R6aM55bXpPwjYAbkEXTbVnKwfDYvVHjNwNFAmx"
	plaintext := "hello world 测试"

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
