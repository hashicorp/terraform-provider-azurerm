package shim

import (
	"context"

	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/containers"
)

type StorageContainerWrapper interface {
	Create(ctx context.Context, resourceGroup, accountName, containerName string, input containers.CreateInput) error
	Delete(ctx context.Context, resourceGroup, accountName, containerName string) error
	Exists(ctx context.Context, resourceGroup, accountName, containerName string) (*bool, error)
	Get(ctx context.Context, resourceGroup, accountName, containerName string) (*StorageContainerProperties, error)
	UpdateAccessLevel(ctx context.Context, resourceGroup, accountName, containerName string, level containers.AccessLevel) error
	UpdateMetaData(ctx context.Context, resourceGroup, accountName, containerName string, metadata map[string]string) error
}

type StorageContainerProperties struct {
	AccessLevel           containers.AccessLevel
	MetaData              map[string]string
	HasImmutabilityPolicy bool
	HasLegalHold          bool
}
