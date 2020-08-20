package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/2019-06-01-preview/synapse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	BigDataPoolClient        *synapse.BigDataPoolsClient
	FirewallRulesClient      *synapse.IPFirewallRulesClient
	WorkspaceClient          *synapse.WorkspacesClient
	WorkspaceAadAdminsClient *synapse.WorkspaceAadAdminsClient
}

func NewClient(o *common.ClientOptions) *Client {
	bigDataPoolClient := synapse.NewBigDataPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&bigDataPoolClient.Client, o.ResourceManagerAuthorizer)

	firewallRuleClient := synapse.NewIPFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&firewallRuleClient.Client, o.ResourceManagerAuthorizer)

	workspaceClient := synapse.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&workspaceClient.Client, o.ResourceManagerAuthorizer)

	workspaceAadAdminsClient := synapse.NewWorkspaceAadAdminsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&workspaceAadAdminsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		BigDataPoolClient:        &bigDataPoolClient,
		FirewallRulesClient:      &firewallRuleClient,
		WorkspaceClient:          &workspaceClient,
		WorkspaceAadAdminsClient: &workspaceAadAdminsClient,
	}
}
