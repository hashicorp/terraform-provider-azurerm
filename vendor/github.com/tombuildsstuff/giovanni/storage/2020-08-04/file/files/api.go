package files

import (
	"context"
	"os"
	"time"

	"github.com/Azure/go-autorest/autorest"
)

type StorageFile interface {
	PutByteRange(ctx context.Context, accountName, shareName, path, fileName string, input PutByteRangeInput) (result autorest.Response, err error)
	GetByteRange(ctx context.Context, accountName, shareName, path, fileName string, input GetByteRangeInput) (result GetByteRangeResult, err error)
	ClearByteRange(ctx context.Context, accountName, shareName, path, fileName string, input ClearByteRangeInput) (result autorest.Response, err error)
	SetProperties(ctx context.Context, accountName, shareName, path, fileName string, input SetPropertiesInput) (result autorest.Response, err error)
	PutFile(ctx context.Context, accountName, shareName, path, fileName string, file *os.File, parallelism int) error
	Copy(ctx context.Context, accountName, shareName, path, fileName string, input CopyInput) (result CopyResult, err error)
	SetMetaData(ctx context.Context, accountName, shareName, path, fileName string, metaData map[string]string) (result autorest.Response, err error)
	GetMetaData(ctx context.Context, accountName, shareName, path, fileName string) (result GetMetaDataResult, err error)
	AbortCopy(ctx context.Context, accountName, shareName, path, fileName, copyID string) (result autorest.Response, err error)
	GetFile(ctx context.Context, accountName, shareName, path, fileName string, parallelism int) (result autorest.Response, outputBytes []byte, err error)
	GetResourceID(accountName, shareName, directoryName, filePath string) string
	ListRanges(ctx context.Context, accountName, shareName, path, fileName string) (result ListRangesResult, err error)
	GetProperties(ctx context.Context, accountName, shareName, path, fileName string) (result GetResult, err error)
	Delete(ctx context.Context, accountName, shareName, path, fileName string) (result autorest.Response, err error)
	Create(ctx context.Context, accountName, shareName, path, fileName string, input CreateInput) (result autorest.Response, err error)
	CopyAndWait(ctx context.Context, accountName, shareName, path, fileName string, input CopyInput, pollDuration time.Duration) (result CopyResult, err error)
}
