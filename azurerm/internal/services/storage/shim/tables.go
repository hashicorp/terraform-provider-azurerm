package shim

import (
	"context"

	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/table/tables"
)

type StorageTableWrapper interface {
	Create(ctx context.Context, resourceGroup string, accountName string, tableName string) error
	Delete(ctx context.Context, resourceGroup string, accountName string, tableName string) error
	Exists(ctx context.Context, resourceGroup string, accountName string, tableName string) (*bool, error)
	GetACLs(ctx context.Context, resourceGroup string, accountName string, tableName string) (*[]tables.SignedIdentifier, error)
	UpdateACLs(ctx context.Context, resourceGroup string, accountName string, tableName string, acls []tables.SignedIdentifier) error
}
