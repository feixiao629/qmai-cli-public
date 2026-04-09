package auth

// This file previously contained AuthTransport for cookie-based auth injection.
// With the open platform migration, authentication is handled directly in
// client.Client.Call() via request body signing.
//
// The transport layer is no longer needed for auth purposes.
// See internal/client/client.go for the new auth flow.
