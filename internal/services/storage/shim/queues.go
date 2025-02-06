// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shim

import (
	"context"

	"github.com/jackofallops/giovanni/storage/2023-11-03/queue/queues"
)

type StorageQueuesWrapper interface {
	Create(ctx context.Context, queueName string, metaData map[string]string) error
	Delete(ctx context.Context, queueName string) error
	Exists(ctx context.Context, queueName string) (*bool, error)
	Get(ctx context.Context, queueName string) (*StorageQueueProperties, error)
	GetServiceProperties(ctx context.Context) (*queues.StorageServiceProperties, error)
	UpdateMetaData(ctx context.Context, queueName string, metaData map[string]string) error
	UpdateServiceProperties(ctx context.Context, properties queues.StorageServiceProperties) error
}

type StorageQueueProperties struct {
	MetaData map[string]string
}
