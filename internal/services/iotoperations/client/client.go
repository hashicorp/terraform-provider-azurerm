package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/instance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	InstanceClient *instance.InstanceClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	instanceClient := instance.NewInstanceClientWithBaseURI(o.ResourceManagerEndpoint)
	instanceClient.Client.Authorizer = o.ResourceManagerAuthorizer

	return &Client{
		InstanceClient: instanceClient,
	}, nil
}