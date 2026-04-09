package core

import (
	"errors"
	"fmt"
)

var (
	ErrNotAuthenticated = errors.New("not authenticated, run 'qmai auth login' first")
	ErrNoActiveProfile  = errors.New("no active profile, run 'qmai config init' first")
	ErrProfileNotFound  = errors.New("profile not found")
	ErrTokenExpired     = errors.New("token expired, run 'qmai auth login' to refresh")
)

// APIError represents an error from the qmai API
type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("API error %d (%s): %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// CancelError indicates user cancelled an operation
type CancelError struct {
	Message string
}

func (e *CancelError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "operation cancelled"
}
