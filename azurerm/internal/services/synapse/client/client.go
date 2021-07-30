package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/2019-06-01-preview/managedvirtualnetwork"
	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/2020-08-01-preview/accesscontrol"
	"github.com/Azure/azure-sdk-for-go/services/synapse/mgmt/2021-03-01/synapse"
	"github.com/Azure/go-autorest/autorest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	FirewallRulesClient                              *synapse.IPFirewallRulesClient
	PrivateLinkHubsClient                            *synapse.PrivateLinkHubsClient
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

	privateLinkHubsClient := synapse.NewPrivateLinkHubsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&privateLinkHubsClient.Client, o.ResourceManagerAuthorizer)

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
		PrivateLinkHubsClient:                            &privateLinkHubsClient,
		SparkPoolClient:                                  &sparkPoolClient,
		SqlPoolClient:                                    &sqlPoolClient,
		SqlPoolTransparentDataEncryptionClient:           &sqlPoolTransparentDataEncryptionClient,
		WorkspaceClient:                                  &workspaceClient,
		WorkspaceAadAdminsClient:                         &workspaceAadAdminsClient,
		WorkspaceManagedIdentitySQLControlSettingsClient: &workspaceManagedIdentitySQLControlSettingsClient,

		synapseAuthorizer: o.SynapseAuthorizer,
	}
}

func (client Client) RoleDefinitionsClient(workspaceName, synapseEndpointSuffix string) (*accesscontrol.RoleDefinitionsClient, error) {
	if client.synapseAuthorizer == nil {
		return nil, fmt.Errorf("Synapse is not supported in this Azure Environment")
	}
	endpoint := buildEndpoint(workspaceName, synapseEndpointSuffix)
	roleDefinitionsClient := accesscontrol.NewRoleDefinitionsClient(endpoint)
	roleDefinitionsClient.Client.Authorizer = client.synapseAuthorizer
	return &roleDefinitionsClient, nil
}

func (client Client) RoleAssignmentsClient(workspaceName, synapseEndpointSuffix string) (*accesscontrol.RoleAssignmentsClient, error) {
	if client.synapseAuthorizer == nil {
		return nil, fmt.Errorf("Synapse is not supported in this Azure Environment")
	}
	endpoint := buildEndpoint(workspaceName, synapseEndpointSuffix)
	roleAssignmentsClient := accesscontrol.NewRoleAssignmentsClient(endpoint)
	roleAssignmentsClient.Client.Authorizer = client.synapseAuthorizer
	return &roleAssignmentsClient, nil
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
