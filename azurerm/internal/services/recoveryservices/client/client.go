package client

import (
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2016-06-01/recoveryservices"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2017-07-01/backup"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-01-10/siterecovery"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ProtectedItemsClient            *backup.ProtectedItemsClient
	ProtectionPoliciesClient        *backup.ProtectionPoliciesClient
	VaultsClient                    *recoveryservices.VaultsClient
	FabricClient                    func(resourceGroupName string, vaultName string) siterecovery.ReplicationFabricsClient
	ProtectionContainerClient       func(resourceGroupName string, vaultName string) siterecovery.ReplicationProtectionContainersClient
	ReplicationPoliciesClient       func(resourceGroupName string, vaultName string) siterecovery.ReplicationPoliciesClient
	ContainerMappingClient          func(resourceGroupName string, vaultName string) siterecovery.ReplicationProtectionContainerMappingsClient
	NetworkMappingClient            func(resourceGroupName string, vaultName string) siterecovery.ReplicationNetworkMappingsClient
	ReplicationMigrationItemsClient func(resourceGroupName string, vaultName string) siterecovery.ReplicationProtectedItemsClient
}

func NewClient(o *common.ClientOptions) *Client {
	vaultsClient := recoveryservices.NewVaultsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vaultsClient.Client, o.ResourceManagerAuthorizer)

	protectedItemsClient := backup.NewProtectedItemsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&protectedItemsClient.Client, o.ResourceManagerAuthorizer)

	protectionPoliciesClient := backup.NewProtectionPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&protectionPoliciesClient.Client, o.ResourceManagerAuthorizer)

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
		ProtectedItemsClient:            &protectedItemsClient,
		ProtectionPoliciesClient:        &protectionPoliciesClient,
		VaultsClient:                    &vaultsClient,
		FabricClient:                    fabricClient,
		ProtectionContainerClient:       protectionContainerClient,
		ReplicationPoliciesClient:       replicationPoliciesClient,
		ContainerMappingClient:          containerMappingClient,
		NetworkMappingClient:            networkMappingClient,
		ReplicationMigrationItemsClient: replicationMigrationItemsClient,
	}
}
