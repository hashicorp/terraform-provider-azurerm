package shares

import (
	"context"
)

type StorageShare interface {
	SetACL(ctx context.Context, shareName string, input SetAclInput) (SetAclResponse, error)
	GetSnapshot(ctx context.Context, shareName string, input GetSnapshotPropertiesInput) (GetSnapshotPropertiesResponse, error)
	GetStats(ctx context.Context, shareName string) (GetStatsResponse, error)
	GetACL(ctx context.Context, shareName string) (GetACLResult, error)
	SetMetaData(ctx context.Context, shareName string, input SetMetaDataInput) (SetMetaDataResponse, error)
	GetMetaData(ctx context.Context, shareName string) (GetMetaDataResponse, error)
	SetProperties(ctx context.Context, shareName string, properties ShareProperties) (SetPropertiesResponse, error)
	DeleteSnapshot(ctx context.Context, accountName string, shareName string, shareSnapshot string) (DeleteSnapshotResponse, error)
	CreateSnapshot(ctx context.Context, shareName string, input CreateSnapshotInput) (CreateSnapshotResponse, error)
	GetResourceManagerResourceID(subscriptionID, resourceGroup, accountName, shareName string) string
	GetProperties(ctx context.Context, shareName string) (GetPropertiesResult, error)
	Delete(ctx context.Context, shareName string, input DeleteInput) (DeleteResponse, error)
	Create(ctx context.Context, shareName string, input CreateInput) (CreateResponse, error)
}
