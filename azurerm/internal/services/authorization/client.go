package authorization

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-09-01-preview/authorization"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	RoleAssignmentsClient *authorization.RoleAssignmentsClient
	RoleDefinitionsClient *authorization.RoleDefinitionsClient
}

func BuildClient(o *common.ClientOptions) *Client {

	RoleAssignmentsClient := authorization.NewRoleAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&RoleAssignmentsClient.Client, o.ResourceManagerAuthorizer)

	RoleDefinitionsClient := authorization.NewRoleDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&RoleDefinitionsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		RoleAssignmentsClient: &RoleAssignmentsClient,
		RoleDefinitionsClient: &RoleDefinitionsClient,
	}
}
