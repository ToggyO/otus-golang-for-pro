package appcore

import "context"

type IStorage interface { // TODO
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	GetConnection() any
}
