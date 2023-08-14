// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/advisor/2020-01-01/getrecommendations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	RecommendationsClient *getrecommendations.GetRecommendationsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	recommendationsClient, err := getrecommendations.NewGetRecommendationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Recommendations client: %+v", err)
	}
	o.Configure(recommendationsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		RecommendationsClient: recommendationsClient,
	}, nil
}
