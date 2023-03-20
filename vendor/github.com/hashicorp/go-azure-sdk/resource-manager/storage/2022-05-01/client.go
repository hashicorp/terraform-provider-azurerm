package v2022_05_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/blobcontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/blobinventorypolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/blobservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/deletedaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/encryptionscopes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/fileservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/fileshares"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/localusers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/managementpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/objectreplicationpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/privateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/privatelinkresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/queueservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/queueserviceproperties"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/skus"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/tableservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/tableserviceproperties"
)

type Client struct {
	BlobContainers             *blobcontainers.BlobContainersClient
	BlobInventoryPolicies      *blobinventorypolicies.BlobInventoryPoliciesClient
	BlobService                *blobservice.BlobServiceClient
	DeletedAccounts            *deletedaccounts.DeletedAccountsClient
	EncryptionScopes           *encryptionscopes.EncryptionScopesClient
	FileService                *fileservice.FileServiceClient
	FileShares                 *fileshares.FileSharesClient
	LocalUsers                 *localusers.LocalUsersClient
	ManagementPolicies         *managementpolicies.ManagementPoliciesClient
	ObjectReplicationPolicies  *objectreplicationpolicies.ObjectReplicationPoliciesClient
	PrivateEndpointConnections *privateendpointconnections.PrivateEndpointConnectionsClient
	PrivateLinkResources       *privatelinkresources.PrivateLinkResourcesClient
	QueueService               *queueservice.QueueServiceClient
	QueueServiceProperties     *queueserviceproperties.QueueServicePropertiesClient
	Skus                       *skus.SkusClient
	StorageAccounts            *storageaccounts.StorageAccountsClient
	TableService               *tableservice.TableServiceClient
	TableServiceProperties     *tableserviceproperties.TableServicePropertiesClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	blobContainersClient := blobcontainers.NewBlobContainersClientWithBaseURI(endpoint)
	configureAuthFunc(&blobContainersClient.Client)

	blobInventoryPoliciesClient := blobinventorypolicies.NewBlobInventoryPoliciesClientWithBaseURI(endpoint)
	configureAuthFunc(&blobInventoryPoliciesClient.Client)

	blobServiceClient := blobservice.NewBlobServiceClientWithBaseURI(endpoint)
	configureAuthFunc(&blobServiceClient.Client)

	deletedAccountsClient := deletedaccounts.NewDeletedAccountsClientWithBaseURI(endpoint)
	configureAuthFunc(&deletedAccountsClient.Client)

	encryptionScopesClient := encryptionscopes.NewEncryptionScopesClientWithBaseURI(endpoint)
	configureAuthFunc(&encryptionScopesClient.Client)

	fileServiceClient := fileservice.NewFileServiceClientWithBaseURI(endpoint)
	configureAuthFunc(&fileServiceClient.Client)

	fileSharesClient := fileshares.NewFileSharesClientWithBaseURI(endpoint)
	configureAuthFunc(&fileSharesClient.Client)

	localUsersClient := localusers.NewLocalUsersClientWithBaseURI(endpoint)
	configureAuthFunc(&localUsersClient.Client)

	managementPoliciesClient := managementpolicies.NewManagementPoliciesClientWithBaseURI(endpoint)
	configureAuthFunc(&managementPoliciesClient.Client)

	objectReplicationPoliciesClient := objectreplicationpolicies.NewObjectReplicationPoliciesClientWithBaseURI(endpoint)
	configureAuthFunc(&objectReplicationPoliciesClient.Client)

	privateEndpointConnectionsClient := privateendpointconnections.NewPrivateEndpointConnectionsClientWithBaseURI(endpoint)
	configureAuthFunc(&privateEndpointConnectionsClient.Client)

	privateLinkResourcesClient := privatelinkresources.NewPrivateLinkResourcesClientWithBaseURI(endpoint)
	configureAuthFunc(&privateLinkResourcesClient.Client)

	queueServiceClient := queueservice.NewQueueServiceClientWithBaseURI(endpoint)
	configureAuthFunc(&queueServiceClient.Client)

	queueServicePropertiesClient := queueserviceproperties.NewQueueServicePropertiesClientWithBaseURI(endpoint)
	configureAuthFunc(&queueServicePropertiesClient.Client)

	skusClient := skus.NewSkusClientWithBaseURI(endpoint)
	configureAuthFunc(&skusClient.Client)

	storageAccountsClient := storageaccounts.NewStorageAccountsClientWithBaseURI(endpoint)
	configureAuthFunc(&storageAccountsClient.Client)

	tableServiceClient := tableservice.NewTableServiceClientWithBaseURI(endpoint)
	configureAuthFunc(&tableServiceClient.Client)

	tableServicePropertiesClient := tableserviceproperties.NewTableServicePropertiesClientWithBaseURI(endpoint)
	configureAuthFunc(&tableServicePropertiesClient.Client)

	return Client{
		BlobContainers:             &blobContainersClient,
		BlobInventoryPolicies:      &blobInventoryPoliciesClient,
		BlobService:                &blobServiceClient,
		DeletedAccounts:            &deletedAccountsClient,
		EncryptionScopes:           &encryptionScopesClient,
		FileService:                &fileServiceClient,
		FileShares:                 &fileSharesClient,
		LocalUsers:                 &localUsersClient,
		ManagementPolicies:         &managementPoliciesClient,
		ObjectReplicationPolicies:  &objectReplicationPoliciesClient,
		PrivateEndpointConnections: &privateEndpointConnectionsClient,
		PrivateLinkResources:       &privateLinkResourcesClient,
		QueueService:               &queueServiceClient,
		QueueServiceProperties:     &queueServicePropertiesClient,
		Skus:                       &skusClient,
		StorageAccounts:            &storageAccountsClient,
		TableService:               &tableServiceClient,
		TableServiceProperties:     &tableServicePropertiesClient,
	}
}
