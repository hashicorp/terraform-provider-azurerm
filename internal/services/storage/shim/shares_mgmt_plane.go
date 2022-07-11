package shim

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-04-01/storage"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/file/shares"

	"github.com/Azure/go-autorest/autorest/date"
)

type ResourceManagerStorageShareWrapper struct {
	client *storage.FileSharesClient
}

func NewManagementPlaneStorageShareWrapper(client *storage.FileSharesClient) StorageShareWrapper {
	return ResourceManagerStorageShareWrapper{
		client: client,
	}
}

func (w ResourceManagerStorageShareWrapper) Create(ctx context.Context, resourceGroup, accountName, shareName string, input shares.CreateInput) error {
	rmInput := storage.FileShare{
		FileShareProperties: &storage.FileShareProperties{
			ShareQuota:       utils.Int32(int32(input.QuotaInGB)),
			EnabledProtocols: storage.EnabledProtocols(input.EnabledProtocol),
			Metadata:         mapStringToMapStringPtr(input.MetaData),
		},
	}
	_, err := w.client.Create(ctx, resourceGroup, accountName, shareName, rmInput, "")
	return err
}

func (w ResourceManagerStorageShareWrapper) Delete(ctx context.Context, resourceGroup, accountName, shareName string) error {
	_, err := w.client.Delete(ctx, resourceGroup, accountName, shareName, "", "")
	return err
}

func (w ResourceManagerStorageShareWrapper) Exists(ctx context.Context, resourceGroup, accountName, shareName string) (*bool, error) {
	share, err := w.client.Get(ctx, resourceGroup, accountName, shareName, "", "")
	if err != nil {
		if utils.ResponseWasNotFound(share.Response) {
			return utils.Bool(false), nil
		}

		return nil, err
	}

	return utils.Bool(share.FileShareProperties != nil), nil
}

func (w ResourceManagerStorageShareWrapper) Get(ctx context.Context, resourceGroup, accountName, shareName string) (*StorageShareProperties, error) {
	share, err := w.client.Get(ctx, resourceGroup, accountName, shareName, "", "")
	if err != nil {
		if utils.ResponseWasNotFound(share.Response) {
			return nil, nil
		}

		return nil, err
	}

	if share.FileShareProperties == nil {
		return nil, fmt.Errorf("`properties` is null in the API response")
	}

	quotaGB := 0
	if quotaPtr := share.FileShareProperties.ShareQuota; quotaPtr != nil {
		quotaGB = int(*quotaPtr)
	}

	// Currently, the service won't return the `enabledProtocols` in the GET response when set to `SMB`. See: https://github.com/Azure/azure-rest-api-specs/issues/16782.
	// Once above issue being addressed, we can remove below code.
	protocol := shares.SMB
	if share.FileShareProperties.EnabledProtocols != "" {
		protocol = shares.ShareProtocol(share.FileShareProperties.EnabledProtocols)
	}

	output := StorageShareProperties{
		ACLs:            w.mapToACLs(share.FileShareProperties.SignedIdentifiers),
		MetaData:        mapStringPtrToMapString(share.FileShareProperties.Metadata),
		QuotaGB:         quotaGB,
		EnabledProtocol: protocol,
	}
	return &output, nil
}

func (w ResourceManagerStorageShareWrapper) UpdateACLs(ctx context.Context, resourceGroup, accountName, shareName string, acls []shares.SignedIdentifier) error {
	identifiers, err := w.mapFromACLs(acls)
	if err != nil {
		return err
	}
	rmInput := storage.FileShare{
		FileShareProperties: &storage.FileShareProperties{
			SignedIdentifiers: identifiers,
		},
	}
	_, err = w.client.Update(ctx, resourceGroup, accountName, shareName, rmInput)
	return err
}

func (w ResourceManagerStorageShareWrapper) UpdateMetaData(ctx context.Context, resourceGroup, accountName, shareName string, metaData map[string]string) error {
	rmInput := storage.FileShare{
		FileShareProperties: &storage.FileShareProperties{
			Metadata: mapStringToMapStringPtr(metaData),
		},
	}
	_, err := w.client.Update(ctx, resourceGroup, accountName, shareName, rmInput)
	return err
}

func (w ResourceManagerStorageShareWrapper) UpdateQuota(ctx context.Context, resourceGroup, accountName, shareName string, quotaGB int) error {
	rmInput := storage.FileShare{
		FileShareProperties: &storage.FileShareProperties{
			ShareQuota: utils.Int32(int32(quotaGB)),
		},
	}
	_, err := w.client.Update(ctx, resourceGroup, accountName, shareName, rmInput)
	return err
}

func (w ResourceManagerStorageShareWrapper) UpdateTier(ctx context.Context, resourceGroup string, accountName string, shareName string, tier shares.AccessTier) error {
	rmInput := storage.FileShare{
		FileShareProperties: &storage.FileShareProperties{
			AccessTier: storage.ShareAccessTier(tier),
		},
	}
	_, err := w.client.Update(ctx, resourceGroup, accountName, shareName, rmInput)
	return err
}

func (w ResourceManagerStorageShareWrapper) mapToACLs(input *[]storage.SignedIdentifier) []shares.SignedIdentifier {
	if input == nil {
		return nil
	}
	var output []shares.SignedIdentifier
	for _, identifier := range *input {
		var id string
		if identifier.ID != nil {
			id = *identifier.ID
		}

		var policy shares.AccessPolicy
		if identifier.AccessPolicy != nil {
			var expiry string
			if identifier.AccessPolicy.Expiry != nil {
				expiry = identifier.AccessPolicy.Expiry.String()
			}
			var start string
			if identifier.AccessPolicy.Start != nil {
				start = identifier.AccessPolicy.Start.String()
			}
			var permission string
			if identifier.AccessPolicy.Permission != nil {
				permission = *identifier.AccessPolicy.Permission
			}
			policy = shares.AccessPolicy{
				Start:      start,
				Expiry:     expiry,
				Permission: permission,
			}
		}
		output = append(output, shares.SignedIdentifier{
			Id:           id,
			AccessPolicy: policy,
		})
	}
	return output
}

func (w ResourceManagerStorageShareWrapper) mapFromACLs(input []shares.SignedIdentifier) (*[]storage.SignedIdentifier, error) {
	if len(input) == 0 {
		return nil, nil
	}

	var output []storage.SignedIdentifier
	for _, identifier := range input {
		policy := identifier.AccessPolicy
		start, err := date.ParseTime(time.RFC3339, policy.Start)
		if err != nil {
			return nil, fmt.Errorf("parsing start time of the ACL %q: %w", policy.Start, err)
		}
		expiry, err := date.ParseTime(time.RFC3339, policy.Expiry)
		if err != nil {
			return nil, fmt.Errorf("parsing expiry time of the ACL %q: %w", policy.Expiry, err)
		}
		output = append(output, storage.SignedIdentifier{
			ID: &identifier.Id,
			AccessPolicy: &storage.AccessPolicy{
				Start:      &date.Time{Time: start},
				Expiry:     &date.Time{Time: expiry},
				Permission: &policy.Permission,
			},
		})
	}
	return &output, nil
}
