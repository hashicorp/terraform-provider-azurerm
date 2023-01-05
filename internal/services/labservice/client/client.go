package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/labplan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	LabPlanClient *labplan.LabPlanClient
}

func NewClient(o *common.ClientOptions) *Client {
	LabPlanClient := labplan.NewLabPlanClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&LabPlanClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		LabPlanClient: &LabPlanClient,
	}
}
