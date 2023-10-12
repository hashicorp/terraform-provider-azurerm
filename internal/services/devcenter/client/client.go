package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/attachednetworkconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/devcenters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DevCenterClient                  *devcenters.DevCentersClient
	AttachedNetworkConnectionsClient *attachednetworkconnections.AttachedNetworkConnectionsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	devCenterClient, err := devcenters.NewDevCentersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Devcenter client: %+v", err)
	}
	o.Configure(devCenterClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		DevCenterClient: devCenterClient,
	}, nil
}
