package shim

import (
	"context"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/queue/queues"
)

type DataPlaneStorageQueueWrapper struct {
	client *queues.Client
}

func NewDataPlaneStorageQueueWrapper(client *queues.Client) StorageQueuesWrapper {
	return DataPlaneStorageQueueWrapper{
		client: client,
	}
}

func (w DataPlaneStorageQueueWrapper) Create(ctx context.Context, _, accountName, queueName string, metaData map[string]string) error {
	_, err := w.client.Create(ctx, accountName, queueName, metaData)
	return err
}

func (w DataPlaneStorageQueueWrapper) Delete(ctx context.Context, _, accountName, queueName string) error {
	_, err := w.client.Delete(ctx, accountName, queueName)
	return err
}

func (w DataPlaneStorageQueueWrapper) Exists(ctx context.Context, _, accountName, queueName string) (*bool, error) {
	existing, err := w.client.GetMetaData(ctx, accountName, queueName)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return utils.Bool(false), nil
		}
		return nil, err
	}

	return utils.Bool(true), nil
}

func (w DataPlaneStorageQueueWrapper) Get(ctx context.Context, _, accountName, queueName string) (*StorageQueueProperties, error) {
	existing, err := w.client.GetMetaData(ctx, accountName, queueName)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return nil, nil
		}
		return nil, err
	}

	return &StorageQueueProperties{
		MetaData: existing.MetaData,
	}, nil
}

func (w DataPlaneStorageQueueWrapper) UpdateMetaData(ctx context.Context, _, accountName, queueName string, metaData map[string]string) error {
	_, err := w.client.SetMetaData(ctx, accountName, queueName, metaData)
	return err
}
