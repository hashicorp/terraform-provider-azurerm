package client

import (
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2016-06-01/recoveryservices"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-01-10/siterecovery"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2019-05-13/backup"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ProtectableItemsClient           *backup.ProtectableItemsClient
	ProtectedItemsClient             *backup.ProtectedItemsClient
	ProtectedItemsGroupClient        *backup.ProtectedItemsGroupClient
	ProtectionPoliciesClient         *backup.ProtectionPoliciesClient
	BackupProtectionContainersClient *backup.ProtectionContainersClient
	BackupOperationStatusesClient    *backup.OperationStatusesClient
	VaultsClient                     *recoveryservices.VaultsClient
	VaultsConfigsClient              *backup.ResourceVaultConfigsClient // Not sure why this is in backup, but https://github.com/Azure/azure-sdk-for-go/issues/7279
	FabricClient                     func(resourceGroupName string, vaultName string) siterecovery.ReplicationFabricsClient
	ProtectionContainerClient        func(resourceGroupName string, vaultName string) siterecovery.ReplicationProtectionContainersClient
	ReplicationPoliciesClient        func(resourceGroupName string, vaultName string) siterecovery.ReplicationPoliciesClient
	ContainerMappingClient           func(resourceGroupName string, vaultName string) siterecovery.ReplicationProtectionContainerMappingsClient
	NetworkMappingClient             func(resourceGroupName string, vaultName string) siterecovery.ReplicationNetworkMappingsClient
	ReplicationMigrationItemsClient  func(resourceGroupName string, vaultName string) siterecovery.ReplicationProtectedItemsClient
}

func NewClient(o *common.ClientOptions) *Client {
	vaultConfigsClient := backup.NewResourceVaultConfigsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vaultConfigsClient.Client, o.ResourceManagerAuthorizer)

	vaultsClient := recoveryservices.NewVaultsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vaultsClient.Client, o.ResourceManagerAuthorizer)

	protectableItemsClient := backup.NewProtectableItemsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&protectableItemsClient.Client, o.ResourceManagerAuthorizer)

	protectedItemsClient := backup.NewProtectedItemsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&protectedItemsClient.Client, o.ResourceManagerAuthorizer)

	protectedItemsGroupClient := backup.NewProtectedItemsGroupClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&protectedItemsGroupClient.Client, o.ResourceManagerAuthorizer)

	protectionPoliciesClient := backup.NewProtectionPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&protectionPoliciesClient.Client, o.ResourceManagerAuthorizer)

	backupProtectionContainersClient := backup.NewProtectionContainersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&backupProtectionContainersClient.Client, o.ResourceManagerAuthorizer)

	backupOperationStatusesClient := backup.NewOperationStatusesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&backupOperationStatusesClient.Client, o.ResourceManagerAuthorizer)

	fabricClient := func(resourceGroupName string, vaultName string) siterecovery.ReplicationFabricsClient {
		client := siterecovery.NewReplicationFabricsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, resourceGroupName, vaultName)
		o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)
		return client
	}

	protectionContainerClient := func(resourceGroupName string, vaultName string) siterecovery.ReplicationProtectionContainersClient {
		client := siterecovery.NewReplicationProtectionContainersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, resourceGroupName, vaultName)
		o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)
		return client
	}

	replicationPoliciesClient := func(resourceGroupName string, vaultName string) siterecovery.ReplicationPoliciesClient {
		client := siterecovery.NewReplicationPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, resourceGroupName, vaultName)
		o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)
		return client
	}

	containerMappingClient := func(resourceGroupName string, vaultName string) siterecovery.ReplicationProtectionContainerMappingsClient {
		client := siterecovery.NewReplicationProtectionContainerMappingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, resourceGroupName, vaultName)
		o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)
		return client
	}

	networkMappingClient := func(resourceGroupName string, vaultName string) siterecovery.ReplicationNetworkMappingsClient {
		client := siterecovery.NewReplicationNetworkMappingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, resourceGroupName, vaultName)
		o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)
		return client
	}

	replicationMigrationItemsClient := func(resourceGroupName string, vaultName string) siterecovery.ReplicationProtectedItemsClient {
		client := siterecovery.NewReplicationProtectedItemsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, resourceGroupName, vaultName)
		o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)
		return client
	}

	return &Client{
		ProtectableItemsClient:           &protectableItemsClient,
		ProtectedItemsClient:             &protectedItemsClient,
		ProtectedItemsGroupClient:        &protectedItemsGroupClient,
		ProtectionPoliciesClient:         &protectionPoliciesClient,
		BackupProtectionContainersClient: &backupProtectionContainersClient,
		BackupOperationStatusesClient:    &backupOperationStatusesClient,
		VaultsClient:                     &vaultsClient,
		VaultsConfigsClient:              &vaultConfigsClient,
		FabricClient:                     fabricClient,
		ProtectionContainerClient:        protectionContainerClient,
		ReplicationPoliciesClient:        replicationPoliciesClient,
		ContainerMappingClient:           containerMappingClient,
		NetworkMappingClient:             networkMappingClient,
		ReplicationMigrationItemsClient:  replicationMigrationItemsClient,
	}
}
