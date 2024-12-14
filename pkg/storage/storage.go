package storage

import (
	"context"
	"fmt"

	"github.com/bonaysoft/gofbot/apis/message/v1alpha1"
)

type Manager interface {
	Start(ctx context.Context) error
	List(ctx context.Context) ([]v1alpha1.Message, error)
}

type construct func() (Manager, error)

var storages = map[string]construct{
	"file": NewStorage,
}

func New(name string) (Manager, error) {
	storageConstruct, ok := storages[name]
	if !ok {
		return nil, fmt.Errorf("%s not exist", name)
	}

	return storageConstruct()
}
