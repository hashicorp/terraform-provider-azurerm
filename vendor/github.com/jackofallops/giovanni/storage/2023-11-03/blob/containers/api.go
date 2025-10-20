package containers

import (
	"context"
)

type StorageContainer interface {
	Create(ctx context.Context, containerName string, input CreateInput) (CreateResponse, error)
	Delete(ctx context.Context, containerName string) (DeleteResponse, error)
	GetProperties(ctx context.Context, containerName string, input GetPropertiesInput) (GetPropertiesResponse, error)
	AcquireLease(ctx context.Context, containerName string, input AcquireLeaseInput) (AcquireLeaseResponse, error)
	BreakLease(ctx context.Context, containerName string, input BreakLeaseInput) (BreakLeaseResponse, error)
	ChangeLease(ctx context.Context, containerName string, input ChangeLeaseInput) (ChangeLeaseResponse, error)
	ReleaseLease(ctx context.Context, containerName string, input ReleaseLeaseInput) (ReleaseLeaseResponse, error)
	RenewLease(ctx context.Context, containerName string, input RenewLeaseInput) (RenewLeaseResponse, error)
	ListBlobs(ctx context.Context, containerName string, input ListBlobsInput) (ListBlobsResponse, error)
	GetResourceManagerResourceID(subscriptionID, resourceGroup, accountName, containerName string) string
	SetAccessControl(ctx context.Context, containerName string, input SetAccessControlInput) (SetAccessControlResponse, error)
	SetMetaData(ctx context.Context, containerName string, metaData SetMetaDataInput) (SetMetaDataResponse, error)
}
