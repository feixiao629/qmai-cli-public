package auth

import "testing"

func TestComputeToken(t *testing.T) {
	// Fake credential material used to keep the test deterministic.
	openKey := "fake-open-key-for-test-only-0001"
	openId := "fake-open-id-for-test-only-0001"
	grantCode := "fake-grant-code-0001"
	var timestamp int64 = 1712563200
	var nonce int64 = 24680

	expected := "ySypL4SRFMmfSrCqOiFE01LK3ks%3D"

	got := ComputeToken(openId, grantCode, openKey, nonce, timestamp)
	if got != expected {
		t.Errorf("ComputeToken() = %q, want %q", got, expected)
	}
}
