package shared

import "context"

type IDbClient interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	GetConnection() interface{}
}
