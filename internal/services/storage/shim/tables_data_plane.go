// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shim

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
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
