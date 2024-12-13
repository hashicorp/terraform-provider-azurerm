package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/fabric/2023-11-01/fabriccapacities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	FabricCapacitiesClient *fabriccapacities.FabricCapacitiesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	fabricCapacitiesClient, err := fabriccapacities.NewFabricCapacitiesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building capacities client: %+v", err)
	}
	o.Configure(fabricCapacitiesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		FabricCapacitiesClient: fabricCapacitiesClient,
	}, nil
}
