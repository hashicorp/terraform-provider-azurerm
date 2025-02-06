package tables

import (
	"context"
)

type StorageTable interface {
	Delete(ctx context.Context, tableName string) (resp DeleteTableResponse, err error)
	Exists(ctx context.Context, tableName string) (resp TableExistsResponse, err error)
	GetACL(ctx context.Context, tableName string) (resp GetACLResponse, err error)
	Create(ctx context.Context, tableName string) (resp CreateTableResponse, err error)
	GetResourceManagerResourceID(subscriptionID, resourceGroup, accountName, tableName string) string
	Query(ctx context.Context, input QueryInput) (resp GetResponse, err error)
	SetACL(ctx context.Context, tableName string, acls []SignedIdentifier) (resp SetACLResponse, err error)
}
