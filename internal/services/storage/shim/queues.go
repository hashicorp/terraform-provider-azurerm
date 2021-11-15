package shim

import (
	"context"
)

type StorageQueuesWrapper interface {
	Create(ctx context.Context, resourceGroup, accountName, queueName string, metaData map[string]string) error
	Delete(ctx context.Context, resourceGroup, accountName, queueName string) error
	Exists(ctx context.Context, resourceGroup, accountName, queueName string) (*bool, error)
	Get(ctx context.Context, resourceGroup, accountName, queueName string) (*StorageQueueProperties, error)
	UpdateMetaData(ctx context.Context, resourceGroup, accountName, queueName string, metaData map[string]string) error
}

type StorageQueueProperties struct {
	MetaData map[string]string
}
