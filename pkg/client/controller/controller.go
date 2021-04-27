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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"time"
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
	conn    *grpc.ClientConn
}

// New creates a new client with the specified address and apiKey.
func New(options Options) Client {
	c := Client{
		address: fmt.Sprintf("%s:%d", options.Host, options.Port),
		cred:    apiKeyCredentials{Key: options.ApiKey},
	}
	return c
}

// Connect establishes a connection to the gRPC service (lazily).
func (ac Client) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, ac.address, grpc.WithInsecure(), grpc.WithPerRPCCredentials(ac.cred))
	if err != nil {
		return fmt.Errorf("failed to Connect %s: %w", ac.address, err)
	}
	ac.conn = conn
	return nil
}

func (ac Client) Close() {
	if ac.conn != nil {
		_ = ac.conn.Close()
	}
}
