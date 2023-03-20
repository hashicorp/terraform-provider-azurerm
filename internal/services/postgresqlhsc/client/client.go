package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	FirewallRulesClient *firewallrules.FirewallRulesClient
}

func NewClient(o *common.ClientOptions) *Client {
	firewallRulesClient := firewallrules.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&firewallRulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		FirewallRulesClient: &firewallRulesClient,
	}
}
