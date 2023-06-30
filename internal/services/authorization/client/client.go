package client

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization" // nolint: staticcheck // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roleassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roledefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	RoleAssignmentsClient    *authorization.RoleAssignmentsClient
	RoleDefinitionsClient    *authorization.RoleDefinitionsClient
	NewRoleAssignmentsClient *roleassignments.RoleAssignmentsClient
	NewRoleDefinitionsClient *roledefinitions.RoleDefinitionsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	roleAssignmentsClient := authorization.NewRoleAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&roleAssignmentsClient.Client, o.ResourceManagerAuthorizer)

	roleDefinitionsClient := authorization.NewRoleDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&roleDefinitionsClient.Client, o.ResourceManagerAuthorizer)

	newRoleAssignmentsClient, err := roleassignments.NewRoleAssignmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Role Assignment Client:  %+v", err)
	}
	o.Configure(newRoleAssignmentsClient.Client, o.Authorizers.ResourceManager)

	newRoleDefinitionsClient, err := roledefinitions.NewRoleDefinitionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Role Definition Client:  %+v", err)
	}
	o.Configure(newRoleDefinitionsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		RoleAssignmentsClient:    &roleAssignmentsClient,
		RoleDefinitionsClient:    &roleDefinitionsClient,
		NewRoleAssignmentsClient: newRoleAssignmentsClient,
		NewRoleDefinitionsClient: newRoleDefinitionsClient,
	}, nil
}
