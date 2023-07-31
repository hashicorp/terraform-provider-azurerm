package client

import (
	"fmt"

	paloalto_2022_08_29 "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	*paloalto_2022_08_29.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	client, err := paloalto_2022_08_29.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building clients for Network: %+v", err)
	}

	return &Client{
		Client: client,
	}, nil
}
