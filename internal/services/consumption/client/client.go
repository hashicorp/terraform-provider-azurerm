package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/consumption/2019-10-01/budgets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	BudgetsClient *budgets.BudgetsClient
}

func NewClient(o *common.ClientOptions) *Client {
	budgetsClient := budgets.NewBudgetsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&budgetsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		BudgetsClient: &budgetsClient,
	}
}
