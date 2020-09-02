package azuread

import (
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	GroupsClient *graphrbac.GroupsClient
}

func NewClient(o *common.ClientOptions) *Client {
	aadGroupsClient := graphrbac.NewGroupsClientWithBaseURI(o.GraphEndpoint, o.TenantID)
	o.ConfigureClient(&aadGroupsClient.Client, o.GraphAuthorizer)
	return &Client{
		GroupsClient: &aadGroupsClient,
	}
}
