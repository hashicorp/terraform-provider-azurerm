// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shim

import (
	"context"

	"github.com/jackofallops/giovanni/storage/2023-11-03/table/tables"
)

type StorageTableWrapper interface {
	Create(ctx context.Context, tableName string) error
	Delete(ctx context.Context, tableName string) error
	Exists(ctx context.Context, tableName string) (*bool, error)
	GetACLs(ctx context.Context, tableName string) (*[]tables.SignedIdentifier, error)
	UpdateACLs(ctx context.Context, tableName string, acls []tables.SignedIdentifier) error
}
