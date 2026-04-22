package v2025_06_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/blobcontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/blobinventorypolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/blobservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/deletedaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/encryptionscopes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/fileservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/fileserviceusageoperationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/fileshares"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/immutabilitypolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/localuseroperationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/managementpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/networksecurityperimeterconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/objectreplicationpolicyoperationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/openapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/privateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/queueservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/storageaccountmigrations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/storagequeues"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/storagetaskassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/tables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/tableservices"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	BlobContainers                         *blobcontainers.BlobContainersClient
	BlobInventoryPolicies                  *blobinventorypolicies.BlobInventoryPoliciesClient
	BlobServices                           *blobservices.BlobServicesClient
	DeletedAccounts                        *deletedaccounts.DeletedAccountsClient
	EncryptionScopes                       *encryptionscopes.EncryptionScopesClient
	FileServiceUsageOperationGroup         *fileserviceusageoperationgroup.FileServiceUsageOperationGroupClient
	FileServices                           *fileservices.FileServicesClient
	FileShares                             *fileshares.FileSharesClient
	ImmutabilityPolicies                   *immutabilitypolicies.ImmutabilityPoliciesClient
	LocalUserOperationGroup                *localuseroperationgroup.LocalUserOperationGroupClient
	ManagementPolicies                     *managementpolicies.ManagementPoliciesClient
	NetworkSecurityPerimeterConfigurations *networksecurityperimeterconfigurations.NetworkSecurityPerimeterConfigurationsClient
	ObjectReplicationPolicyOperationGroup  *objectreplicationpolicyoperationgroup.ObjectReplicationPolicyOperationGroupClient
	Openapis                               *openapis.OpenapisClient
	PrivateEndpointConnections             *privateendpointconnections.PrivateEndpointConnectionsClient
	QueueServices                          *queueservices.QueueServicesClient
	StorageAccountMigrations               *storageaccountmigrations.StorageAccountMigrationsClient
	StorageAccounts                        *storageaccounts.StorageAccountsClient
	StorageQueues                          *storagequeues.StorageQueuesClient
	StorageTaskAssignments                 *storagetaskassignments.StorageTaskAssignmentsClient
	TableServices                          *tableservices.TableServicesClient
	Tables                                 *tables.TablesClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	blobContainersClient, err := blobcontainers.NewBlobContainersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building BlobContainers client: %+v", err)
	}
	configureFunc(blobContainersClient.Client)

	blobInventoryPoliciesClient, err := blobinventorypolicies.NewBlobInventoryPoliciesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building BlobInventoryPolicies client: %+v", err)
	}
	configureFunc(blobInventoryPoliciesClient.Client)

	blobServicesClient, err := blobservices.NewBlobServicesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building BlobServices client: %+v", err)
	}
	configureFunc(blobServicesClient.Client)

	deletedAccountsClient, err := deletedaccounts.NewDeletedAccountsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DeletedAccounts client: %+v", err)
	}
	configureFunc(deletedAccountsClient.Client)

	encryptionScopesClient, err := encryptionscopes.NewEncryptionScopesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building EncryptionScopes client: %+v", err)
	}
	configureFunc(encryptionScopesClient.Client)

	fileServiceUsageOperationGroupClient, err := fileserviceusageoperationgroup.NewFileServiceUsageOperationGroupClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building FileServiceUsageOperationGroup client: %+v", err)
	}
	configureFunc(fileServiceUsageOperationGroupClient.Client)

	fileServicesClient, err := fileservices.NewFileServicesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building FileServices client: %+v", err)
	}
	configureFunc(fileServicesClient.Client)

	fileSharesClient, err := fileshares.NewFileSharesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building FileShares client: %+v", err)
	}
	configureFunc(fileSharesClient.Client)

	immutabilityPoliciesClient, err := immutabilitypolicies.NewImmutabilityPoliciesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ImmutabilityPolicies client: %+v", err)
	}
	configureFunc(immutabilityPoliciesClient.Client)

	localUserOperationGroupClient, err := localuseroperationgroup.NewLocalUserOperationGroupClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LocalUserOperationGroup client: %+v", err)
	}
	configureFunc(localUserOperationGroupClient.Client)

	managementPoliciesClient, err := managementpolicies.NewManagementPoliciesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ManagementPolicies client: %+v", err)
	}
	configureFunc(managementPoliciesClient.Client)

	networkSecurityPerimeterConfigurationsClient, err := networksecurityperimeterconfigurations.NewNetworkSecurityPerimeterConfigurationsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NetworkSecurityPerimeterConfigurations client: %+v", err)
	}
	configureFunc(networkSecurityPerimeterConfigurationsClient.Client)

	objectReplicationPolicyOperationGroupClient, err := objectreplicationpolicyoperationgroup.NewObjectReplicationPolicyOperationGroupClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ObjectReplicationPolicyOperationGroup client: %+v", err)
	}
	configureFunc(objectReplicationPolicyOperationGroupClient.Client)

	openapisClient, err := openapis.NewOpenapisClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Openapis client: %+v", err)
	}
	configureFunc(openapisClient.Client)

	privateEndpointConnectionsClient, err := privateendpointconnections.NewPrivateEndpointConnectionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PrivateEndpointConnections client: %+v", err)
	}
	configureFunc(privateEndpointConnectionsClient.Client)

	queueServicesClient, err := queueservices.NewQueueServicesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building QueueServices client: %+v", err)
	}
	configureFunc(queueServicesClient.Client)

	storageAccountMigrationsClient, err := storageaccountmigrations.NewStorageAccountMigrationsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building StorageAccountMigrations client: %+v", err)
	}
	configureFunc(storageAccountMigrationsClient.Client)

	storageAccountsClient, err := storageaccounts.NewStorageAccountsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building StorageAccounts client: %+v", err)
	}
	configureFunc(storageAccountsClient.Client)

	storageQueuesClient, err := storagequeues.NewStorageQueuesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building StorageQueues client: %+v", err)
	}
	configureFunc(storageQueuesClient.Client)

	storageTaskAssignmentsClient, err := storagetaskassignments.NewStorageTaskAssignmentsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building StorageTaskAssignments client: %+v", err)
	}
	configureFunc(storageTaskAssignmentsClient.Client)

	tableServicesClient, err := tableservices.NewTableServicesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TableServices client: %+v", err)
	}
	configureFunc(tableServicesClient.Client)

	tablesClient, err := tables.NewTablesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Tables client: %+v", err)
	}
	configureFunc(tablesClient.Client)

	return &Client{
		BlobContainers:                         blobContainersClient,
		BlobInventoryPolicies:                  blobInventoryPoliciesClient,
		BlobServices:                           blobServicesClient,
		DeletedAccounts:                        deletedAccountsClient,
		EncryptionScopes:                       encryptionScopesClient,
		FileServiceUsageOperationGroup:         fileServiceUsageOperationGroupClient,
		FileServices:                           fileServicesClient,
		FileShares:                             fileSharesClient,
		ImmutabilityPolicies:                   immutabilityPoliciesClient,
		LocalUserOperationGroup:                localUserOperationGroupClient,
		ManagementPolicies:                     managementPoliciesClient,
		NetworkSecurityPerimeterConfigurations: networkSecurityPerimeterConfigurationsClient,
		ObjectReplicationPolicyOperationGroup:  objectReplicationPolicyOperationGroupClient,
		Openapis:                               openapisClient,
		PrivateEndpointConnections:             privateEndpointConnectionsClient,
		QueueServices:                          queueServicesClient,
		StorageAccountMigrations:               storageAccountMigrationsClient,
		StorageAccounts:                        storageAccountsClient,
		StorageQueues:                          storageQueuesClient,
		StorageTaskAssignments:                 storageTaskAssignmentsClient,
		TableServices:                          tableServicesClient,
		Tables:                                 tablesClient,
	}, nil
}
