package shares

import (
	"context"

	"github.com/Azure/go-autorest/autorest"
)

type StorageShare interface {
	SetACL(ctx context.Context, accountName, shareName string, acls []SignedIdentifier) (result autorest.Response, err error)
	GetSnapshot(ctx context.Context, accountName, shareName, snapshotShare string) (result GetSnapshotPropertiesResult, err error)
	GetStats(ctx context.Context, accountName, shareName string) (result GetStatsResult, err error)
	GetACL(ctx context.Context, accountName, shareName string) (result GetACLResult, err error)
	SetMetaData(ctx context.Context, accountName, shareName string, metaData map[string]string) (result autorest.Response, err error)
	GetMetaData(ctx context.Context, accountName, shareName string) (result GetMetaDataResult, err error)
	SetProperties(ctx context.Context, accountName, shareName string, newQuotaGB int) (result autorest.Response, err error)
	DeleteSnapshot(ctx context.Context, accountName, shareName string, shareSnapshot string) (result autorest.Response, err error)
	CreateSnapshot(ctx context.Context, accountName, shareName string, input CreateSnapshotInput) (result CreateSnapshotResult, err error)
	GetResourceID(accountName, shareName string) string
	GetResourceManagerResourceID(subscriptionID, resourceGroup, accountName, shareName string) string
	GetProperties(ctx context.Context, accountName, shareName string) (result GetPropertiesResult, err error)
	Delete(ctx context.Context, accountName, shareName string, deleteSnapshots bool) (result autorest.Response, err error)
	Create(ctx context.Context, accountName, shareName string, input CreateInput) (result autorest.Response, err error)
}
