// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/shim"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/blobs"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/containers"
	"github.com/jackofallops/giovanni/storage/2023-11-03/datalakestore/filesystems"
	"github.com/jackofallops/giovanni/storage/2023-11-03/datalakestore/paths"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/directories"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/files"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/shares"
	"github.com/jackofallops/giovanni/storage/2023-11-03/queue/queues"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/entities"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/tables"
)

type DataPlaneOperation struct {
	SupportsAadAuthentication       bool
	SupportsSharedKeyAuthentication bool

	sharedKeyAuthenticationType auth.SharedKeyType
}

func (Client) DataPlaneOperationSupportingAnyAuthMethod() DataPlaneOperation {
	return DataPlaneOperation{
		SupportsAadAuthentication:       true,
		SupportsSharedKeyAuthentication: true,
	}
}

func (Client) DataPlaneOperationSupportingOnlySharedKeyAuth() DataPlaneOperation {
	return DataPlaneOperation{
		SupportsAadAuthentication:       false,
		SupportsSharedKeyAuthentication: true,
	}
}

func (c Client) configureDataPlane(ctx context.Context, clientName, resourceIdentifier string, baseClient client.BaseClient, account AccountDetails, operation DataPlaneOperation) error {
	if operation.SupportsAadAuthentication && c.authConfigForAzureAD != nil {
		api := c.authConfigForAzureAD.Environment.Storage.WithResourceIdentifier(resourceIdentifier)
		storageAuth, err := auth.NewAuthorizerFromCredentials(ctx, *c.authConfigForAzureAD, api)
		if err != nil {
			return fmt.Errorf("unable to build authorizer for Storage API: %+v", err)
		}

		baseClient.SetAuthorizer(storageAuth)
		return nil
	}

	if operation.SupportsSharedKeyAuthentication {
		accountKey, err := account.AccountKey(ctx, c)
		if err != nil {
			return fmt.Errorf("retrieving Storage Account Key: %s", err)
		}

		storageAuth, err := auth.NewSharedKeyAuthorizer(account.StorageAccountId.StorageAccountName, *accountKey, operation.sharedKeyAuthenticationType)
		if err != nil {
			return fmt.Errorf("building Shared Key Authorizer for %s client: %+v", clientName, err)
		}

		baseClient.SetAuthorizer(storageAuth)
		return nil
	}

	return fmt.Errorf("building %s client: no configured authentication types are supported", clientName)
}

func (c Client) AccountsDataPlaneClient(ctx context.Context, account AccountDetails, operation DataPlaneOperation) (*accounts.Client, error) {
	const clientName = "Blob Storage Accounts"
	operation.sharedKeyAuthenticationType = auth.SharedKey

	baseUri, err := account.DataPlaneEndpoint(EndpointTypeBlob)
	if err != nil {
		return nil, err
	}

	apiClient, err := accounts.NewWithBaseUri(*baseUri)
	if err != nil {
		return nil, fmt.Errorf("building %s client: %+v", clientName, err)
	}

	err = c.configureDataPlane(ctx, clientName, *baseUri, apiClient.Client, account, operation)
	if err != nil {
		return nil, err
	}

	return apiClient, nil
}

func (c Client) BlobsDataPlaneClient(ctx context.Context, account AccountDetails, operation DataPlaneOperation) (*blobs.Client, error) {
	const clientName = "Blob Storage Blobs"
	operation.sharedKeyAuthenticationType = auth.SharedKey

	baseUri, err := account.DataPlaneEndpoint(EndpointTypeBlob)
	if err != nil {
		return nil, err
	}

	apiClient, err := blobs.NewWithBaseUri(*baseUri)
	if err != nil {
		return nil, fmt.Errorf("building %s client: %+v", clientName, err)
	}

	err = c.configureDataPlane(ctx, clientName, *baseUri, apiClient.Client, account, operation)
	if err != nil {
		return nil, err
	}

	return apiClient, nil
}

func (c Client) ContainersDataPlaneClient(ctx context.Context, account AccountDetails, operation DataPlaneOperation) (shim.StorageContainerWrapper, error) {
	const clientName = "Blob Storage Containers"
	operation.sharedKeyAuthenticationType = auth.SharedKey

	baseUri, err := account.DataPlaneEndpoint(EndpointTypeBlob)
	if err != nil {
		return nil, err
	}

	apiClient, err := containers.NewWithBaseUri(*baseUri)
	if err != nil {
		return nil, fmt.Errorf("building %s client: %+v", clientName, err)
	}

	err = c.configureDataPlane(ctx, clientName, *baseUri, apiClient.Client, account, operation)
	if err != nil {
		return nil, err
	}

	return shim.NewDataPlaneStorageContainerWrapper(apiClient), nil
}

func (c Client) DataLakeFilesystemsDataPlaneClient(ctx context.Context, account AccountDetails, operation DataPlaneOperation) (*filesystems.Client, error) {
	const clientName = "Data Lake Gen2 Filesystems"
	operation.sharedKeyAuthenticationType = auth.SharedKey

	baseUri, err := account.DataPlaneEndpoint(EndpointTypeDfs)
	if err != nil {
		return nil, err
	}

	apiClient, err := filesystems.NewWithBaseUri(*baseUri)
	if err != nil {
		return nil, fmt.Errorf("building %s client: %+v", clientName, err)
	}

	err = c.configureDataPlane(ctx, clientName, *baseUri, apiClient.Client, account, operation)
	if err != nil {
		return nil, err
	}

	return apiClient, nil
}

func (c Client) DataLakePathsDataPlaneClient(ctx context.Context, account AccountDetails, operation DataPlaneOperation) (*paths.Client, error) {
	const clientName = "Data Lake Gen2 Paths"
	operation.sharedKeyAuthenticationType = auth.SharedKey

	baseUri, err := account.DataPlaneEndpoint(EndpointTypeDfs)
	if err != nil {
		return nil, err
	}

	apiClient, err := paths.NewWithBaseUri(*baseUri)
	if err != nil {
		return nil, fmt.Errorf("building %s client: %+v", clientName, err)
	}

	err = c.configureDataPlane(ctx, clientName, *baseUri, apiClient.Client, account, operation)
	if err != nil {
		return nil, err
	}

	return apiClient, nil
}

func (c Client) FileShareDirectoriesDataPlaneClient(ctx context.Context, account AccountDetails, operation DataPlaneOperation) (*directories.Client, error) {
	const clientName = "File Storage Share Directories"
	operation.sharedKeyAuthenticationType = auth.SharedKey

	baseUri, err := account.DataPlaneEndpoint(EndpointTypeFile)
	if err != nil {
		return nil, err
	}

	apiClient, err := directories.NewWithBaseUri(*baseUri)
	if err != nil {
		return nil, fmt.Errorf("building %s client: %+v", clientName, err)
	}

	err = c.configureDataPlane(ctx, clientName, *baseUri, apiClient.Client, account, operation)
	if err != nil {
		return nil, err
	}

	return apiClient, nil
}

func (c Client) FileShareFilesDataPlaneClient(ctx context.Context, account AccountDetails, operation DataPlaneOperation) (*files.Client, error) {
	const clientName = "File Storage Share Files"
	operation.sharedKeyAuthenticationType = auth.SharedKey

	baseUri, err := account.DataPlaneEndpoint(EndpointTypeFile)
	if err != nil {
		return nil, err
	}

	apiClient, err := files.NewWithBaseUri(*baseUri)
	if err != nil {
		return nil, fmt.Errorf("building %s client: %+v", clientName, err)
	}

	err = c.configureDataPlane(ctx, clientName, *baseUri, apiClient.Client, account, operation)
	if err != nil {
		return nil, err
	}

	return apiClient, nil
}

func (c Client) FileSharesDataPlaneClient(ctx context.Context, account AccountDetails, operation DataPlaneOperation) (shim.StorageShareWrapper, error) {
	const clientName = "File Storage Shares"
	operation.sharedKeyAuthenticationType = auth.SharedKey

	baseUri, err := account.DataPlaneEndpoint(EndpointTypeFile)
	if err != nil {
		return nil, err
	}

	apiClient, err := shares.NewWithBaseUri(*baseUri)
	if err != nil {
		return nil, fmt.Errorf("building %s client: %+v", clientName, err)
	}

	err = c.configureDataPlane(ctx, clientName, *baseUri, apiClient.Client, account, operation)
	if err != nil {
		return nil, err
	}

	return shim.NewDataPlaneStorageShareWrapper(apiClient), nil
}

func (c Client) QueuesDataPlaneClient(ctx context.Context, account AccountDetails, operation DataPlaneOperation) (shim.StorageQueuesWrapper, error) {
	const clientName = "File Storage Queue Queues"
	operation.sharedKeyAuthenticationType = auth.SharedKey

	baseUri, err := account.DataPlaneEndpoint(EndpointTypeQueue)
	if err != nil {
		return nil, err
	}

	apiClient, err := queues.NewWithBaseUri(*baseUri)
	if err != nil {
		return nil, fmt.Errorf("building %s client: %+v", clientName, err)
	}

	err = c.configureDataPlane(ctx, clientName, *baseUri, apiClient.Client, account, operation)
	if err != nil {
		return nil, err
	}

	return shim.NewDataPlaneStorageQueueWrapper(apiClient), nil
}

func (c Client) TableEntityDataPlaneClient(ctx context.Context, account AccountDetails, operation DataPlaneOperation) (*entities.Client, error) {
	const clientName = "Table Storage Share Entities"
	operation.sharedKeyAuthenticationType = auth.SharedKeyTable

	baseUri, err := account.DataPlaneEndpoint(EndpointTypeTable)
	if err != nil {
		return nil, err
	}

	apiClient, err := entities.NewWithBaseUri(*baseUri)
	if err != nil {
		return nil, fmt.Errorf("building %s client: %+v", clientName, err)
	}

	err = c.configureDataPlane(ctx, clientName, *baseUri, apiClient.Client, account, operation)
	if err != nil {
		return nil, err
	}

	return apiClient, nil
}

func (c Client) TablesDataPlaneClient(ctx context.Context, account AccountDetails, operation DataPlaneOperation) (shim.StorageTableWrapper, error) {
	const clientName = "Table Storage Share Tables"
	operation.sharedKeyAuthenticationType = auth.SharedKeyTable

	baseUri, err := account.DataPlaneEndpoint(EndpointTypeTable)
	if err != nil {
		return nil, err
	}

	apiClient, err := tables.NewWithBaseUri(*baseUri)
	if err != nil {
		return nil, fmt.Errorf("building %s client: %+v", clientName, err)
	}

	err = c.configureDataPlane(ctx, clientName, *baseUri, apiClient.Client, account, operation)
	if err != nil {
		return nil, err
	}

	return shim.NewDataPlaneStorageTableWrapper(apiClient), nil
}
