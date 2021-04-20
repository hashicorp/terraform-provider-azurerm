package client

import (
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	GroupsClient            *graphrbac.GroupsClient
	RoleAssignmentsClient   *authorization.RoleAssignmentsClient
	RoleDefinitionsClient   *authorization.RoleDefinitionsClient
	ServicePrincipalsClient *graphrbac.ServicePrincipalsClient
}

func NewClient(o *common.ClientOptions) *Client {
	groupsClient := graphrbac.NewGroupsClientWithBaseURI(o.GraphEndpoint, o.TenantID)
	o.ConfigureClient(&groupsClient.Client, o.GraphAuthorizer)

	roleAssignmentsClient := authorization.NewRoleAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&roleAssignmentsClient.Client, o.ResourceManagerAuthorizer)

	roleDefinitionsClient := authorization.NewRoleDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&roleDefinitionsClient.Client, o.ResourceManagerAuthorizer)

	servicePrincipalsClient := graphrbac.NewServicePrincipalsClientWithBaseURI(o.GraphEndpoint, o.TenantID)
	o.ConfigureClient(&servicePrincipalsClient.Client, o.GraphAuthorizer)

	return &Client{
		GroupsClient:            &groupsClient,
		RoleAssignmentsClient:   &roleAssignmentsClient,
		RoleDefinitionsClient:   &roleDefinitionsClient,
		ServicePrincipalsClient: &servicePrincipalsClient,
	}
}
