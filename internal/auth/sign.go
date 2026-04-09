package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// ComputeToken calculates the HmacSHA1 token for open platform authentication.
// It sorts the fields alphabetically, joins them as key=value pairs, signs with
// HmacSHA1 using openKey, then Base64 + URL encodes the result.
func ComputeToken(openId, grantCode, openKey string, nonce int64, timestamp int64) string {
	params := map[string]string{
		"openId":    openId,
		"grantCode": grantCode,
		"nonce":     fmt.Sprintf("%d", nonce),
		"timestamp": fmt.Sprintf("%d", timestamp),
	}

	// Sort keys alphabetically
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build the string to sign: key=value&key=value
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = k + "=" + params[k]
	}
	signStr := strings.Join(parts, "&")

	// HmacSHA1
	mac := hmac.New(sha1.New, []byte(openKey))
	mac.Write([]byte(signStr))
	signature := mac.Sum(nil)

	// Base64 encode then URL encode
	b64 := base64.StdEncoding.EncodeToString(signature)
	return url.QueryEscape(b64)
}
