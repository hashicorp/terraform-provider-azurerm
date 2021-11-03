package queues

import (
	"context"

	"github.com/Azure/go-autorest/autorest"
)

type StorageQueue interface {
	Delete(ctx context.Context, accountName, queueName string) (result autorest.Response, err error)
	GetMetaData(ctx context.Context, accountName, queueName string) (result GetMetaDataResult, err error)
	SetMetaData(ctx context.Context, accountName, queueName string, metaData map[string]string) (result autorest.Response, err error)
	Create(ctx context.Context, accountName, queueName string, metaData map[string]string) (result autorest.Response, err error)
	GetResourceID(accountName, queueName string) string
	SetServiceProperties(ctx context.Context, accountName string, properties StorageServiceProperties) (result autorest.Response, err error)
	GetServiceProperties(ctx context.Context, accountName string) (result StorageServicePropertiesResponse, err error)
}
