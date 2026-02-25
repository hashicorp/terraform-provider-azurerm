package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ComputeFleetClient *fleets.FleetsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	computeFleetsClient, err := fleets.NewFleetsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building compute fleet client: %+v", err)
	}
	o.Configure(computeFleetsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ComputeFleetClient: computeFleetsClient,
	}, nil
}
