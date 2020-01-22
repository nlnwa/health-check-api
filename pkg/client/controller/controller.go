// Copyright 2018 National Library of Norway
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package aggregator contains an aggregator service client
package controller

import (
	"context"
	"fmt"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Options struct {
	Host   string
	Port   int
	ApiKey string
}

// Client represents the client to the aggregator service.
type Client struct {
	address string // address in the form "host:port"
	cred    credentials.PerRPCCredentials
}

// New creates a new client with the specified address and apiKey.
func New(options Options) Client {
	return Client{
		address: options.Host + ":" + strconv.FormatInt(int64(options.Port), 10),
		cred:    apiKeyCredentials{Key: options.ApiKey},
	}
}

// Dial makes a connection to the gRPC service.
func (ac Client) dial(ctx context.Context) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx, ac.address, grpc.WithInsecure(), grpc.WithPerRPCCredentials(ac.cred))
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", ac.address, err)
	}
	return conn, nil
}

// Hangup closes the connection to the gRPC service.
func (ac Client) hangup(conn *grpc.ClientConn) {
	if conn != nil {
		_ = conn.Close()
	}
}
