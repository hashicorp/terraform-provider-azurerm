package files

import (
	"context"
	"os"
)

type StorageFile interface {
	PutByteRange(ctx context.Context, shareName string, path string, fileName string, input PutByteRangeInput) (PutRangeResponse, error)
	GetByteRange(ctx context.Context, shareName string, path string, fileName string, input GetByteRangeInput) (GetByteRangeResponse, error)
	ClearByteRange(ctx context.Context, shareName string, path string, fileName string, input ClearByteRangeInput) (ClearByteRangeResponse, error)
	SetProperties(ctx context.Context, shareName string, path string, fileName string, input SetPropertiesInput) (SetPropertiesResponse, error)
	PutFile(ctx context.Context, shareName string, path string, fileName string, file *os.File, parallelism int) error
	Copy(ctx context.Context, shareName, path, fileName string, input CopyInput) (CopyResponse, error)
	SetMetaData(ctx context.Context, shareName string, path string, fileName string, input SetMetaDataInput) (SetMetaDataResponse, error)
	GetMetaData(ctx context.Context, shareName string, path string, fileName string) (GetMetaDataResponse, error)
	AbortCopy(ctx context.Context, shareName string, path string, fileName string, input CopyAbortInput) (CopyAbortResponse, error)
	GetFile(ctx context.Context, shareName string, path string, fileName string, input GetFileInput) (GetFileResponse, error)
	ListRanges(ctx context.Context, shareName, path, fileName string) (ListRangesResponse, error)
	GetProperties(ctx context.Context, shareName string, path string, fileName string) (GetResponse, error)
	Delete(ctx context.Context, shareName string, path string, fileName string) (DeleteResponse, error)
	Create(ctx context.Context, shareName string, path string, fileName string, input CreateInput) (CreateResponse, error)
	CopyAndWait(ctx context.Context, shareName, path, fileName string, input CopyInput) (CopyResponse, error)
}
