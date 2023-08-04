// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/consumption/2019-10-01/budgets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	BudgetsClient *budgets.BudgetsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	budgetsClient, err := budgets.NewBudgetsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Budgets client: %+v", err)
	}
	o.Configure(budgetsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		BudgetsClient: budgetsClient,
	}, nil
}
