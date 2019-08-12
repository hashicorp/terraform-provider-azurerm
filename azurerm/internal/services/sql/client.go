package sql

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2015-05-01-preview/sql"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DatabasesClient                       sql.DatabasesClient
	DatabaseThreatDetectionPoliciesClient sql.DatabaseThreatDetectionPoliciesClient
	ElasticPoolsClient                    sql.ElasticPoolsClient
	FirewallRulesClient                   sql.FirewallRulesClient
	ServersClient                         sql.ServersClient
	ServerAzureADAdministratorsClient     sql.ServerAzureADAdministratorsClient
	VirtualNetworkRulesClient             sql.VirtualNetworkRulesClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	// SQL Azure
	c.DatabasesClient = sql.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DatabasesClient.Client, o.ResourceManagerAuthorizer)

	c.DatabaseThreatDetectionPoliciesClient = sql.NewDatabaseThreatDetectionPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DatabaseThreatDetectionPoliciesClient.Client, o.ResourceManagerAuthorizer)

	c.ElasticPoolsClient = sql.NewElasticPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ElasticPoolsClient.Client, o.ResourceManagerAuthorizer)

	c.FirewallRulesClient = sql.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.FirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	c.ServersClient = sql.NewServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ServersClient.Client, o.ResourceManagerAuthorizer)

	c.ServerAzureADAdministratorsClient = sql.NewServerAzureADAdministratorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ServerAzureADAdministratorsClient.Client, o.ResourceManagerAuthorizer)

	c.VirtualNetworkRulesClient = sql.NewVirtualNetworkRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VirtualNetworkRulesClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
