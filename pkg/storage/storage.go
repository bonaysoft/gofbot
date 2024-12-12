package storage

import (
	"context"

	"github.com/bonaysoft/gofbot/apis/message/v1alpha1"
)

type Manager interface {
	Start(ctx context.Context) error
	List(ctx context.Context) ([]v1alpha1.Message, error)
}
