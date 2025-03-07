package database

import "context"

type DatabaseInterface interface {
	Ping(ctx context.Context) error
	Close() error
}
