package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/labservices/sdk/2021-10-01-preview/labplan"
)

type Client struct {
	LabPlanClient *labplan.LabPlanClient
}

func NewClient(o *common.ClientOptions) *Client {
	labPlanClient := labplan.NewLabPlanClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&labPlanClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		LabPlanClient: &labPlanClient,
	}
}
