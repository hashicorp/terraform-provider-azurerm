// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shim

import (
	"context"

	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/table/tables"
)

type StorageTableWrapper interface {
	Create(ctx context.Context, resourceGroup string, accountName string, tableName string) error
	Delete(ctx context.Context, resourceGroup string, accountName string, tableName string) error
	Exists(ctx context.Context, resourceGroup string, accountName string, tableName string) (*bool, error)
	GetACLs(ctx context.Context, resourceGroup string, accountName string, tableName string) (*[]tables.SignedIdentifier, error)
	UpdateACLs(ctx context.Context, resourceGroup string, accountName string, tableName string, acls []tables.SignedIdentifier) error
}
