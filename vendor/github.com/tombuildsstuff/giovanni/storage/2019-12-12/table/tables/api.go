package tables

import (
	"context"

	"github.com/Azure/go-autorest/autorest"
)

type StorageTable interface {
	Delete(ctx context.Context, accountName, tableName string) (result autorest.Response, err error)
	Exists(ctx context.Context, accountName, tableName string) (result autorest.Response, err error)
	GetACL(ctx context.Context, accountName, tableName string) (result GetACLResult, err error)
	Create(ctx context.Context, accountName, tableName string) (result autorest.Response, err error)
	GetResourceID(accountName, tableName string) string
	Query(ctx context.Context, accountName string, metaDataLevel MetaDataLevel) (result GetResult, err error)
	SetACL(ctx context.Context, accountName, tableName string, acls []SignedIdentifier) (result autorest.Response, err error)
}
