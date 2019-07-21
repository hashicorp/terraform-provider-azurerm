package recoveryservices

import (
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2016-06-01/recoveryservices"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2017-07-01/backup"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-01-10/siterecovery"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ProtectedItemsClient      backup.ProtectedItemsGroupClient
	ProtectionPoliciesClient  backup.ProtectionPoliciesClient
	VaultsClient              recoveryservices.VaultsClient
	FabricClient              func(resourceGroupName string, vaultName string) siterecovery.ReplicationFabricsClient
	ProtectionContainerClient func(resourceGroupName string, vaultName string) siterecovery.ReplicationProtectionContainersClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.VaultsClient = recoveryservices.NewVaultsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VaultsClient.Client, o.ResourceManagerAuthorizer)

	c.ProtectedItemsClient = backup.NewProtectedItemsGroupClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ProtectedItemsClient.Client, o.ResourceManagerAuthorizer)

	c.ProtectionPoliciesClient = backup.NewProtectionPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ProtectionPoliciesClient.Client, o.ResourceManagerAuthorizer)

	c.FabricClient = func(resourceGroupName string, vaultName string) siterecovery.ReplicationFabricsClient {
		client := siterecovery.NewReplicationFabricsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, resourceGroupName, vaultName)
		o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)
		return client
	}

	c.ProtectionContainerClient = func(resourceGroupName string, vaultName string) siterecovery.ReplicationProtectionContainersClient {
		client := siterecovery.NewReplicationProtectionContainersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, resourceGroupName, vaultName)
		o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)
		return client
	}

	return &c
}
