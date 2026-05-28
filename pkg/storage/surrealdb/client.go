package surrealdb

import (
	"context"
	"fmt"
)

// Client defines the storage boundary used by Headlamp.
//
// The first implementation targets SurrealDB as the
// local graph/document operational store.
type Client struct {
	Endpoint  string
	Namespace string
	Database  string
	Username  string
	Password  string
}

// Connect validates configuration and establishes a connection.
//
// Actual SurrealDB driver integration will be added next.
func (c *Client) Connect(ctx context.Context) error {
	if c.Endpoint == "" {
		return fmt.Errorf("missing surrealdb endpoint")
	}

	if c.Namespace == "" {
		return fmt.Errorf("missing surrealdb namespace")
	}

	if c.Database == "" {
		return fmt.Errorf("missing surrealdb database")
	}

	return nil
}
