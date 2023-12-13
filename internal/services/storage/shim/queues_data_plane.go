// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shim

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/queue/queues"
)

type DataPlaneStorageQueueWrapper struct {
	client *queues.Client
}

func NewDataPlaneStorageQueueWrapper(client *queues.Client) StorageQueuesWrapper {
	return DataPlaneStorageQueueWrapper{
		client: client,
	}
}

func (w DataPlaneStorageQueueWrapper) Create(ctx context.Context, _, queueName string, metaData map[string]string) error {
	input := queues.CreateInput{
		MetaData: metaData,
	}
	_, err := w.client.Create(ctx, queueName, input)
	return err
}

func (w DataPlaneStorageQueueWrapper) Delete(ctx context.Context, _, queueName string) error {
	_, err := w.client.Delete(ctx, queueName)
	return err
}

func (w DataPlaneStorageQueueWrapper) Exists(ctx context.Context, _, queueName string) (*bool, error) {
	existing, err := w.client.GetMetaData(ctx, queueName)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse.Response) {
			return utils.Bool(false), nil
		}
		return nil, err
	}

	return utils.Bool(true), nil
}

func (w DataPlaneStorageQueueWrapper) Get(ctx context.Context, _, queueName string) (*StorageQueueProperties, error) {
	props, err := w.client.GetMetaData(ctx, queueName)
	if err != nil {
		if response.WasNotFound(props.HttpResponse.Response) {
			return nil, nil
		}
		return nil, err
	}

	return &StorageQueueProperties{
		MetaData: props.MetaData,
	}, nil
}

func (w DataPlaneStorageQueueWrapper) GetServiceProperties(ctx context.Context) (*queues.StorageServiceProperties, error) {
	serviceProps, err := w.client.GetServiceProperties(ctx)
	if err != nil {
		if response.WasNotFound(serviceProps.HttpResponse.Response) {
			return nil, nil
		}
		return nil, err
	}

	return &serviceProps.StorageServiceProperties, nil
}

func (w DataPlaneStorageQueueWrapper) UpdateMetaData(ctx context.Context, _, queueName string, metaData map[string]string) error {
	input := queues.SetMetaDataInput{
		MetaData: metaData,
	}
	_, err := w.client.SetMetaData(ctx, queueName, input)
	return err
}

func (w DataPlaneStorageQueueWrapper) UpdateServiceProperties(ctx context.Context, _ string, properties queues.StorageServiceProperties) error {
	input := queues.SetStorageServicePropertiesInput{
		Properties: properties,
	}
	_, err := w.client.SetServiceProperties(ctx, input)
	return err
}
