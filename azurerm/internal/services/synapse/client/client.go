package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/2019-06-01-preview/managedvirtualnetwork"
	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/2020-02-01-preview/accesscontrol"
	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/2019-06-01-preview/synapse"
	"github.com/Azure/go-autorest/autorest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	FirewallRulesClient                              *synapse.IPFirewallRulesClient
	SparkPoolClient                                  *synapse.BigDataPoolsClient
	SqlPoolClient                                    *synapse.SQLPoolsClient
	SqlPoolTransparentDataEncryptionClient           *synapse.SQLPoolTransparentDataEncryptionsClient
	WorkspaceClient                                  *synapse.WorkspacesClient
	WorkspaceAadAdminsClient                         *synapse.WorkspaceAadAdminsClient
	WorkspaceManagedIdentitySQLControlSettingsClient *synapse.WorkspaceManagedIdentitySQLControlSettingsClient

	synapseAuthorizer autorest.Authorizer
}

func NewClient(o *common.ClientOptions) *Client {
	firewallRuleClient := synapse.NewIPFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&firewallRuleClient.Client, o.ResourceManagerAuthorizer)

	// the service team hopes to rename it to sparkPool, so rename the sdk here
	sparkPoolClient := synapse.NewBigDataPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sparkPoolClient.Client, o.ResourceManagerAuthorizer)

	sqlPoolClient := synapse.NewSQLPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlPoolClient.Client, o.ResourceManagerAuthorizer)

	sqlPoolTransparentDataEncryptionClient := synapse.NewSQLPoolTransparentDataEncryptionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlPoolTransparentDataEncryptionClient.Client, o.ResourceManagerAuthorizer)

	workspaceClient := synapse.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&workspaceClient.Client, o.ResourceManagerAuthorizer)

	workspaceAadAdminsClient := synapse.NewWorkspaceAadAdminsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&workspaceAadAdminsClient.Client, o.ResourceManagerAuthorizer)

	workspaceManagedIdentitySQLControlSettingsClient := synapse.NewWorkspaceManagedIdentitySQLControlSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&workspaceManagedIdentitySQLControlSettingsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		FirewallRulesClient:                              &firewallRuleClient,
		SparkPoolClient:                                  &sparkPoolClient,
		SqlPoolClient:                                    &sqlPoolClient,
		SqlPoolTransparentDataEncryptionClient:           &sqlPoolTransparentDataEncryptionClient,
		WorkspaceClient:                                  &workspaceClient,
		WorkspaceAadAdminsClient:                         &workspaceAadAdminsClient,
		WorkspaceManagedIdentitySQLControlSettingsClient: &workspaceManagedIdentitySQLControlSettingsClient,

		synapseAuthorizer: o.SynapseAuthorizer,
	}
}

func (client Client) AccessControlClient(workspaceName, synapseEndpointSuffix string) (*accesscontrol.BaseClient, error) {
	if client.synapseAuthorizer == nil {
		return nil, fmt.Errorf("Synapse is not supported in this Azure Environment")
	}
	endpoint := buildEndpoint(workspaceName, synapseEndpointSuffix)
	accessControlClient := accesscontrol.New(endpoint)
	accessControlClient.Client.Authorizer = client.synapseAuthorizer
	return &accessControlClient, nil
}

func (client Client) ManagedPrivateEndpointsClient(workspaceName, synapseEndpointSuffix string) (*managedvirtualnetwork.ManagedPrivateEndpointsClient, error) {
	if client.synapseAuthorizer == nil {
		return nil, fmt.Errorf("Synapse is not supported in this Azure Environment")
	}
	endpoint := buildEndpoint(workspaceName, synapseEndpointSuffix)
	managedPrivateEndpointsClient := managedvirtualnetwork.NewManagedPrivateEndpointsClient(endpoint)
	managedPrivateEndpointsClient.Client.Authorizer = client.synapseAuthorizer
	return &managedPrivateEndpointsClient, nil
}

func buildEndpoint(workspaceName string, synapseEndpointSuffix string) string {
	return fmt.Sprintf("https://%s.%s", workspaceName, synapseEndpointSuffix)
}
