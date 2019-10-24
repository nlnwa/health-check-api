package controller

import "context"

type apiKeyCredentials struct {
	Key string
}

// implement credentials.PerRPCCredentials interface
func (akc apiKeyCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "apikey " + akc.Key,
	}, nil
}

// implement credentials.PerRPCCredentials interface
func (akc apiKeyCredentials) RequireTransportSecurity() bool {
	return false
}
