package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-06-01/resourcevalidationclient"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	PreflightClient *resourcevalidationclient.ResourceValidationClientClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	preflightCLient, err := resourcevalidationclient.NewResourceValidationClientClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Preflight client: %+v", err)
	}
	o.Configure(preflightCLient.Client, o.Authorizers.ResourceManager)

	return &Client{
		PreflightClient: preflightCLient,
	}, nil
}
