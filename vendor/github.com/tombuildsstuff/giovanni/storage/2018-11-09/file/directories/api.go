package directories

import (
	"context"

	"github.com/Azure/go-autorest/autorest"
)

type StorageDirectory interface {
	Delete(ctx context.Context, accountName, shareName, path string) (result autorest.Response, err error)
	GetMetaData(ctx context.Context, accountName, shareName, path string) (result GetMetaDataResult, err error)
	SetMetaData(ctx context.Context, accountName, shareName, path string, metaData map[string]string) (result autorest.Response, err error)
	Create(ctx context.Context, accountName, shareName, path string, metaData map[string]string) (result autorest.Response, err error)
	GetResourceID(accountName, shareName, directoryName string) string
	Get(ctx context.Context, accountName, shareName, path string) (result GetResult, err error)
}
