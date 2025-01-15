package directories

import (
	"context"
)

type StorageDirectory interface {
	Delete(ctx context.Context, shareName, path string) (resp DeleteResponse, err error)
	GetMetaData(ctx context.Context, shareName, path string) (resp GetMetaDataResponse, err error)
	SetMetaData(ctx context.Context, shareName, path string, input SetMetaDataInput) (resp SetMetaDataResponse, err error)
	Create(ctx context.Context, shareName, path string, input CreateDirectoryInput) (resp CreateDirectoryResponse, err error)
	Get(ctx context.Context, shareName, path string) (resp GetResponse, err error)
}
