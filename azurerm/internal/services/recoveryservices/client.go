package recoveryservices

import (
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2016-06-01/recoveryservices"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2017-07-01/backup"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-01-10/siterecovery"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ProtectedItemsClient            *backup.ProtectedItemsGroupClient
	ProtectionPoliciesClient        *backup.ProtectionPoliciesClient
	VaultsClient                    *recoveryservices.VaultsClient
	FabricClient                    func(resourceGroupName string, vaultName string) siterecovery.ReplicationFabricsClient
	ProtectionContainerClient       func(resourceGroupName string, vaultName string) siterecovery.ReplicationProtectionContainersClient
	ReplicationPoliciesClient       func(resourceGroupName string, vaultName string) siterecovery.ReplicationPoliciesClient
	ContainerMappingClient          func(resourceGroupName string, vaultName string) siterecovery.ReplicationProtectionContainerMappingsClient
	NetworkMappingClient            func(resourceGroupName string, vaultName string) siterecovery.ReplicationNetworkMappingsClient
	ReplicationMigrationItemsClient func(resourceGroupName string, vaultName string) siterecovery.ReplicationProtectedItemsClient
}

func BuildClient(o *common.ClientOptions) *Client {

	VaultsClient := recoveryservices.NewVaultsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VaultsClient.Client, o.ResourceManagerAuthorizer)

	ProtectedItemsClient := backup.NewProtectedItemsGroupClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ProtectedItemsClient.Client, o.ResourceManagerAuthorizer)

	ProtectionPoliciesClient := backup.NewProtectionPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ProtectionPoliciesClient.Client, o.ResourceManagerAuthorizer)

	FabricClient := func(resourceGroupName string, vaultName string) siterecovery.ReplicationFabricsClient {
		client := siterecovery.NewReplicationFabricsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, resourceGroupName, vaultName)
		o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)
		return client
	}

	ProtectionContainerClient := func(resourceGroupName string, vaultName string) siterecovery.ReplicationProtectionContainersClient {
		client := siterecovery.NewReplicationProtectionContainersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, resourceGroupName, vaultName)
		o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)
		return client
	}

	ReplicationPoliciesClient := func(resourceGroupName string, vaultName string) siterecovery.ReplicationPoliciesClient {
		client := siterecovery.NewReplicationPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, resourceGroupName, vaultName)
		o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)
		return client
	}

	ContainerMappingClient := func(resourceGroupName string, vaultName string) siterecovery.ReplicationProtectionContainerMappingsClient {
		client := siterecovery.NewReplicationProtectionContainerMappingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, resourceGroupName, vaultName)
		o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)
		return client
	}

	NetworkMappingClient := func(resourceGroupName string, vaultName string) siterecovery.ReplicationNetworkMappingsClient {
		client := siterecovery.NewReplicationNetworkMappingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, resourceGroupName, vaultName)
		o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)
		return client
	}

	ReplicationMigrationItemsClient := func(resourceGroupName string, vaultName string) siterecovery.ReplicationProtectedItemsClient {
		client := siterecovery.NewReplicationProtectedItemsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, resourceGroupName, vaultName)
		o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)
		return client
	}

	return &Client{
		ProtectedItemsClient:            &ProtectedItemsClient,
		ProtectionPoliciesClient:        &ProtectionPoliciesClient,
		VaultsClient:                    &VaultsClient,
		FabricClient:                    FabricClient,
		ProtectionContainerClient:       ProtectionContainerClient,
		ReplicationPoliciesClient:       ReplicationPoliciesClient,
		ContainerMappingClient:          ContainerMappingClient,
		NetworkMappingClient:            NetworkMappingClient,
		ReplicationMigrationItemsClient: ReplicationMigrationItemsClient,
	}
}
