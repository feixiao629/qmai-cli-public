package client

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/madaima/qmai-cli/internal/auth"
)

// APIResponse is the standard response from the open platform.
type APIResponse struct {
	Status  bool            `json:"status"`
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	TraceId string          `json:"traceId,omitempty"`
}

// Client is the open platform API client.
type Client struct {
	BaseURL   string
	OpenId    string
	GrantCode string
	OpenKey   string
	HTTP      *http.Client
	Debug     bool
}

// NewClient creates a new open platform client.
func NewClient(baseURL, openId, grantCode, openKey string, httpClient *http.Client, debug bool) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return &Client{
		BaseURL:   baseURL,
		OpenId:    openId,
		GrantCode: grantCode,
		OpenKey:   openKey,
		HTTP:      httpClient,
		Debug:     debug,
	}
}

// openPlatformRequest is the standard request body envelope.
type openPlatformRequest struct {
	OpenId    string      `json:"openId"`
	GrantCode string      `json:"grantCode"`
	Nonce     int64       `json:"nonce"`
	Timestamp int64       `json:"timestamp"`
	Token     string      `json:"token"`
	Params    interface{} `json:"params"`
}

// Call makes an open platform API call. All requests are POST with the standard envelope.
func (c *Client) Call(ctx context.Context, path string, params interface{}) (*APIResponse, error) {
	nonce := generateNonce()
	timestamp := time.Now().Unix()
	token := auth.ComputeToken(c.OpenId, c.GrantCode, c.OpenKey, nonce, timestamp)

	reqBody := openPlatformRequest{
		OpenId:    c.OpenId,
		GrantCode: c.GrantCode,
		Nonce:     nonce,
		Timestamp: timestamp,
		Token:     token,
		Params:    params,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := c.BaseURL + strings.TrimPrefix(path, "/")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if c.Debug {
		log.Printf("[DEBUG] POST %s | req=%s | resp=%s | %dms", url, string(data), string(respBody), time.Since(start).Milliseconds())
	}

	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
		}
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if !apiResp.Status || apiResp.Code != 0 {
		return nil, fmt.Errorf("API error (code=%d, traceId=%s): %s", apiResp.Code, apiResp.TraceId, apiResp.Message)
	}

	return &apiResp, nil
}

// generateNonce returns a random positive integer for request signing.
func generateNonce() int64 {
	n, err := rand.Int(rand.Reader, big.NewInt(99999))
	if err != nil {
		return time.Now().UnixNano() % 99999
	}
	return n.Int64() + 1
}
