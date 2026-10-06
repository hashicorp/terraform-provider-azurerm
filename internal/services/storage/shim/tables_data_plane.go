// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package shim

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/tables"
)

type DataPlaneStorageTableWrapper struct {
	client *tables.Client
}

func NewDataPlaneStorageTableWrapper(client *tables.Client) StorageTableWrapper {
	return DataPlaneStorageTableWrapper{
		client: client,
	}
}

func (w DataPlaneStorageTableWrapper) Create(ctx context.Context, tableName string) error {
	_, err := w.client.Create(ctx, tableName)
	return err
}

func (w DataPlaneStorageTableWrapper) Delete(ctx context.Context, tableName string) error {
	_, err := w.client.Delete(ctx, tableName)
	return err
}

func (w DataPlaneStorageTableWrapper) Exists(ctx context.Context, tableName string) (*bool, error) {
	existing, err := w.client.Exists(ctx, tableName)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, err
	}
	return pointer.To(true), nil
}

func (w DataPlaneStorageTableWrapper) GetACLs(ctx context.Context, tableName string) (*[]tables.SignedIdentifier, error) {
	acls, err := w.client.GetACL(ctx, tableName)
	if err != nil {
		return nil, err
	}

	return &acls.SignedIdentifiers, nil
}

func (w DataPlaneStorageTableWrapper) UpdateACLs(ctx context.Context, tableName string, acls []tables.SignedIdentifier) error {
	_, err := w.client.SetACL(ctx, tableName, acls)
	return err
}

func (w DataPlaneStorageTableWrapper) GetServiceProperties(ctx context.Context) (*tables.StorageServiceProperties, error) {
	serviceProps, err := w.client.GetServiceProperties(ctx)
	if err != nil {
		if serviceProps.HttpResponse == nil {
			return nil, pollers.PollingDroppedConnectionError{
				Message: err.Error(),
			}
		}
		if response.WasNotFound(serviceProps.HttpResponse) {
			return nil, nil
		}
		return nil, err
	}

	return &serviceProps.StorageServiceProperties, nil
}

func (w DataPlaneStorageTableWrapper) UpdateServiceProperties(ctx context.Context, properties tables.StorageServiceProperties) error {
	input := tables.SetStorageServicePropertiesInput{
		Properties: properties,
	}
	_, err := w.client.SetServiceProperties(ctx, input)
	return err
}
