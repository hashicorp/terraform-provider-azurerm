// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/advisor/2023-01-01/getrecommendations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/advisor/2023-01-01/suppressions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	RecommendationsClient *getrecommendations.GetRecommendationsClient
	SuppressionsClient    *suppressions.SuppressionsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	recommendationsClient, err := getrecommendations.NewGetRecommendationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building recommendations client: %+v", err)
	}
	o.Configure(recommendationsClient.Client, o.Authorizers.ResourceManager)

	suppressionsClient, err := suppressions.NewSuppressionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building suppressions client: %+v", err)
	}
	o.Configure(suppressionsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		RecommendationsClient: recommendationsClient,
		SuppressionsClient:    suppressionsClient,
	}, nil
}
