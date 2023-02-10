package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization" // nolint: staticcheck // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/rolemanagementpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	RoleAssignmentsClient        *authorization.RoleAssignmentsClient
	RoleDefinitionsClient        *authorization.RoleDefinitionsClient
	RoleManagementPoliciesClient *rolemanagementpolicies.RoleManagementPoliciesClient
}

func NewClient(o *common.ClientOptions) *Client {
	roleAssignmentsClient := authorization.NewRoleAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&roleAssignmentsClient.Client, o.ResourceManagerAuthorizer)

	roleDefinitionsClient := authorization.NewRoleDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&roleDefinitionsClient.Client, o.ResourceManagerAuthorizer)

	roleManagementPoliciesClient := rolemanagementpolicies.NewRoleManagementPoliciesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&roleManagementPoliciesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		RoleAssignmentsClient:        &roleAssignmentsClient,
		RoleDefinitionsClient:        &roleDefinitionsClient,
		RoleManagementPoliciesClient: &roleManagementPoliciesClient,
	}
}
