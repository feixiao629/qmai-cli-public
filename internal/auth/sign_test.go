package auth

import "testing"

func TestComputeToken(t *testing.T) {
	// Test data from the open platform documentation
	openKey := "LyvrkvkxRkG2R6aM55bXpPwjYAbkEXTbVnKwfDYvVHjNwNFAmx"
	openId := "d14c1559e87b747d577c834b275a4310"
	grantCode := "ba67d4fa46"
	var timestamp int64 = 1465185768
	var nonce int64 = 11886

	// The expected token from docs (URL-encoded base64)
	expected := "cFw0t9IuvL9jVo9qAzk0qMcw5BM%3D"

	got := ComputeToken(openId, grantCode, openKey, nonce, timestamp)
	if got != expected {
		t.Errorf("ComputeToken() = %q, want %q", got, expected)
	}
}
