package client

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	preflightClient "github.com/hashicorp/terraform-provider-azurerm/internal/preflight/sdk"
)

type Client struct {
	PreflightClient *preflightClient.PreflightClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	preflightCLient, err := preflightClient.NewResourceValidationClientClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Preflight client: %+v", err)
	}
	o.Configure(preflightCLient.Client, o.Authorizers.ResourceManager)

	return &Client{
		PreflightClient: preflightCLient,
	}, nil
}
