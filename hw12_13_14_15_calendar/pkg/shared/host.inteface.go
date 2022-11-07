package shared

import "context"

type IHost interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
