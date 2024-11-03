package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/astronomer/2023-08-01/organizations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	OrganizationsClient *organizations.OrganizationsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	organizationsClient, err := organizations.NewOrganizationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building appliances Client: %+v", err)
	}
	o.Configure(organizationsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		OrganizationsClient: organizationsClient,
	}, nil
}
