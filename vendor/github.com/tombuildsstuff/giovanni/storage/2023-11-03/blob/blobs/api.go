package blobs

import (
	"context"
	"os"
)

type StorageBlob interface {
	AppendBlock(ctx context.Context, containerName string, blobName string, input AppendBlockInput) (AppendBlockResponse, error)
	Copy(ctx context.Context, containerName string, blobName string, input CopyInput) (CopyResponse, error)
	AbortCopy(ctx context.Context, containerName string, blobName string, input AbortCopyInput) (CopyAbortResponse, error)
	CopyAndWait(ctx context.Context, containerName string, blobName string, input CopyInput) error
	Delete(ctx context.Context, containerName string, blobName string, input DeleteInput) (DeleteResponse, error)
	DeleteSnapshot(ctx context.Context, containerName string, blobName string, input DeleteSnapshotInput) (DeleteSnapshotResponse, error)
	DeleteSnapshots(ctx context.Context, containerName string, blobName string, input DeleteSnapshotsInput) (DeleteSnapshotsResponse, error)
	Get(ctx context.Context, containerName string, blobName string, input GetInput) (GetResponse, error)
	GetBlockList(ctx context.Context, containerName string, blobName string, input GetBlockListInput) (GetBlockListResponse, error)
	GetPageRanges(ctx context.Context, containerName, blobName string, input GetPageRangesInput) (GetPageRangesResponse, error)
	IncrementalCopyBlob(ctx context.Context, containerName string, blobName string, input IncrementalCopyBlobInput) (IncrementalCopyBlob, error)
	AcquireLease(ctx context.Context, containerName string, blobName string, input AcquireLeaseInput) (AcquireLeaseResponse, error)
	BreakLease(ctx context.Context, containerName string, blobName string, input BreakLeaseInput) (BreakLeaseResponse, error)
	ChangeLease(ctx context.Context, containerName string, blobName string, input ChangeLeaseInput) (ChangeLeaseResponse, error)
	ReleaseLease(ctx context.Context, containerName string, blobName string, input ReleaseLeaseInput) (ReleaseLeaseResponse, error)
	RenewLease(ctx context.Context, containerName string, blobName string, input RenewLeaseInput) (RenewLeaseResponse, error)
	SetMetaData(ctx context.Context, containerName string, blobName string, input SetMetaDataInput) (SetMetaDataResponse, error)
	GetProperties(ctx context.Context, containerName string, blobName string, input GetPropertiesInput) (GetPropertiesResponse, error)
	SetProperties(ctx context.Context, containerName string, blobName string, input SetPropertiesInput) (SetPropertiesResponse, error)
	PutAppendBlob(ctx context.Context, containerName string, blobName string, input PutAppendBlobInput) (PutAppendBlobResponse, error)
	PutBlock(ctx context.Context, containerName string, blobName string, input PutBlockInput) (PutBlockResponse, error)
	PutBlockBlob(ctx context.Context, containerName string, blobName string, input PutBlockBlobInput) (PutBlockBlobResponse, error)
	PutBlockBlobFromFile(ctx context.Context, containerName string, blobName string, file *os.File, input PutBlockBlobInput) error
	PutBlockList(ctx context.Context, containerName string, blobName string, input PutBlockListInput) (PutBlockListResponse, error)
	PutBlockFromURL(ctx context.Context, containerName string, blobName string, input PutBlockFromURLInput) (PutBlockFromURLResponse, error)
	PutPageBlob(ctx context.Context, containerName string, blobName string, input PutPageBlobInput) (PutPageBlobResponse, error)
	PutPageClear(ctx context.Context, containerName string, blobName string, input PutPageClearInput) (PutPageClearResponse, error)
	PutPageUpdate(ctx context.Context, containerName string, blobName string, input PutPageUpdateInput) (PutPageUpdateResponse, error)
	SetTier(ctx context.Context, containerName string, blobName string, input SetTierInput) (SetTierResponse, error)
	Snapshot(ctx context.Context, containerName string, blobName string, input SnapshotInput) (SnapshotResponse, error)
	GetSnapshotProperties(ctx context.Context, containerName string, blobName string, input GetSnapshotPropertiesInput) (GetPropertiesResponse, error)
	Undelete(ctx context.Context, containerName string, blobName string) (UndeleteResponse, error)
}
