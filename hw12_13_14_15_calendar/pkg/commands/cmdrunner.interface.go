package commands

import "context"

type ICmdRunner interface {
	Name() string
	Description() string
	Init(args []string) error
	Run(ctx context.Context) error
}
