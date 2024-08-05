package queues

import (
	"context"
)

type StorageQueue interface {
	Delete(ctx context.Context, queueName string) (DeleteResponse, error)
	GetMetaData(ctx context.Context, queueName string) (GetMetaDataResponse, error)
	SetMetaData(ctx context.Context, queueName string, input SetMetaDataInput) (SetMetaDataResponse, error)
	Create(ctx context.Context, queueName string, input CreateInput) (CreateResponse, error)
	GetResourceManagerResourceID(subscriptionID, resourceGroup, accountName, queueName string) string
	SetServiceProperties(ctx context.Context, input SetStorageServicePropertiesInput) (SetStorageServicePropertiesResponse, error)
	GetServiceProperties(ctx context.Context) (GetStorageServicePropertiesResponse, error)
}
