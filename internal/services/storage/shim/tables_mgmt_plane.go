package shim

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/utils"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/table/tables"
)

type ResourceManagerStorageTableWrapper struct {
	client *storage.TableClient
}

func NewManagementPlaneStorageTableWrapper(client *storage.TableClient) StorageTableWrapper {
	return ResourceManagerStorageTableWrapper{
		client: client,
	}
}

func (w ResourceManagerStorageTableWrapper) Create(ctx context.Context, resourceGroup string, accountName string, tableName string) error {
	_, err := w.client.Create(ctx, resourceGroup, accountName, tableName, &storage.Table{})
	return err
}

func (w ResourceManagerStorageTableWrapper) Delete(ctx context.Context, resourceGroup string, accountName string, tableName string) error {
	_, err := w.client.Delete(ctx, resourceGroup, accountName, tableName)
	return err
}

func (w ResourceManagerStorageTableWrapper) Exists(ctx context.Context, resourceGroup string, accountName string, tableName string) (*bool, error) {
	table, err := w.client.Get(ctx, resourceGroup, accountName, tableName)
	if err != nil {
		if utils.ResponseWasNotFound(table.Response) {
			return utils.Bool(false), nil
		}

		return nil, err
	}

	return utils.Bool(table.TableProperties != nil), nil
}

func (w ResourceManagerStorageTableWrapper) GetACLs(ctx context.Context, resourceGroup string, accountName string, tableName string) (*[]tables.SignedIdentifier, error) {
	table, err := w.client.Get(ctx, resourceGroup, accountName, tableName)
	if err != nil {
		return nil, err
	}
	if prop := table.TableProperties; prop != nil {
		identifiers := w.mapToACLs(prop.SignedIdentifiers)
		return &identifiers, nil
	}
	return nil, nil
}

func (w ResourceManagerStorageTableWrapper) UpdateACLs(ctx context.Context, resourceGroup string, accountName string, tableName string, acls []tables.SignedIdentifier) error {
	identifiers, err := w.mapFromACLs(acls)
	if err != nil {
		return err
	}
	_, err = w.client.Update(ctx, resourceGroup, accountName, tableName, &storage.Table{
		TableProperties: &storage.TableProperties{
			SignedIdentifiers: identifiers,
		},
	})
	return err
}

func (w ResourceManagerStorageTableWrapper) mapToACLs(input *[]storage.TableSignedIdentifier) []tables.SignedIdentifier {
	if input == nil {
		return nil
	}
	var output []tables.SignedIdentifier
	for _, identifier := range *input {
		var id string
		if identifier.ID != nil {
			id = *identifier.ID
		}

		var policy tables.AccessPolicy
		if identifier.AccessPolicy != nil {
			var expiry string
			if identifier.AccessPolicy.ExpiryTime != nil {
				expiry = identifier.AccessPolicy.ExpiryTime.String()
			}
			var start string
			if identifier.AccessPolicy.StartTime != nil {
				start = identifier.AccessPolicy.StartTime.String()
			}
			var permission string
			if identifier.AccessPolicy.Permission != nil {
				permission = *identifier.AccessPolicy.Permission
			}
			policy = tables.AccessPolicy{
				Start:      start,
				Expiry:     expiry,
				Permission: permission,
			}
		}
		output = append(output, tables.SignedIdentifier{
			Id:           id,
			AccessPolicy: policy,
		})
	}
	return output
}

func (w ResourceManagerStorageTableWrapper) mapFromACLs(input []tables.SignedIdentifier) (*[]storage.TableSignedIdentifier, error) {
	if len(input) == 0 {
		return nil, nil
	}

	var output []storage.TableSignedIdentifier
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
		output = append(output, storage.TableSignedIdentifier{
			ID: &identifier.Id,
			AccessPolicy: &storage.TableAccessPolicy{
				StartTime:  &date.Time{Time: start},
				ExpiryTime: &date.Time{Time: expiry},
				Permission: &policy.Permission,
			},
		})
	}
	return &output, nil
}
