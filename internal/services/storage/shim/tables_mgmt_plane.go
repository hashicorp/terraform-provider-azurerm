package shim

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-provider-azurerm/utils"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-04-01/storage"
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
	_, err := w.client.Create(ctx, resourceGroup, accountName, tableName)
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
	// TODO: Implement this once the following issue is resolved: https://github.com/Azure/azure-rest-api-specs/issues/17007
	return nil, errors.New("Storage Table management plane API doesn't support ACLs now")
}

func (w ResourceManagerStorageTableWrapper) UpdateACLs(ctx context.Context, resourceGroup string, accountName string, tableName string, acls []tables.SignedIdentifier) error {
	// TODO: Implement this once the following issue is resolved: https://github.com/Azure/azure-rest-api-specs/issues/17007
	return errors.New("Storage Table management plane API doesn't support ACLs now")
}
