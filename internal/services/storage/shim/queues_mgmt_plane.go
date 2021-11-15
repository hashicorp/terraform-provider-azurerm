package shim

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/utils"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-04-01/storage"
)

type ResourceManagerStorageQueueWrapper struct {
	client *storage.QueueClient
}

func NewMgmtPlaneStorageQueueWrapper(client *storage.QueueClient) StorageQueuesWrapper {
	return ResourceManagerStorageQueueWrapper{
		client: client,
	}
}

func (w ResourceManagerStorageQueueWrapper) Create(ctx context.Context, resourceGroup, accountName, queueName string, metaData map[string]string) error {
	rmInput := storage.Queue{
		QueueProperties: &storage.QueueProperties{
			Metadata: mapStringToMapStringPtr(metaData),
		},
	}
	_, err := w.client.Create(ctx, resourceGroup, accountName, queueName, rmInput)
	return err
}

func (w ResourceManagerStorageQueueWrapper) Delete(ctx context.Context, resourceGroup, accountName, queueName string) error {
	resp, err := w.client.Delete(ctx, resourceGroup, accountName, queueName)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return err
	}
	return nil
}

func (w ResourceManagerStorageQueueWrapper) Exists(ctx context.Context, resourceGroup, accountName, queueName string) (*bool, error) {
	queue, err := w.client.Get(ctx, resourceGroup, accountName, queueName)
	if err != nil {
		if utils.ResponseWasNotFound(queue.Response) {
			return utils.Bool(false), nil
		}

		return nil, err
	}

	return utils.Bool(queue.QueueProperties != nil), nil
}

func (w ResourceManagerStorageQueueWrapper) Get(ctx context.Context, resourceGroup, accountName, queueName string) (*StorageQueueProperties, error) {
	queue, err := w.client.Get(ctx, resourceGroup, accountName, queueName)
	if err != nil {
		if utils.ResponseWasNotFound(queue.Response) {
			return nil, nil
		}

		return nil, err
	}

	if queue.QueueProperties == nil {
		return nil, fmt.Errorf("`properties` is null in the API response")
	}

	output := StorageQueueProperties{
		MetaData: mapStringPtrToMapString(queue.QueueProperties.Metadata),
	}
	return &output, nil
}

func (w ResourceManagerStorageQueueWrapper) UpdateMetaData(ctx context.Context, resourceGroup, accountName, queueName string, metaData map[string]string) error {
	rmInput := storage.Queue{
		QueueProperties: &storage.QueueProperties{
			Metadata: mapStringToMapStringPtr(metaData),
		},
	}
	_, err := w.client.Update(ctx, resourceGroup, accountName, queueName, rmInput)
	return err
}
