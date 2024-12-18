// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2024-10-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2024-10-01/deployments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2024-10-01/raiblocklists"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AccountsClient      *cognitiveservicesaccounts.CognitiveServicesAccountsClient
	DeploymentsClient   *deployments.DeploymentsClient
	RaiBlocklistsClient *raiblocklists.RaiBlocklistsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	accountsClient, err := cognitiveservicesaccounts.NewCognitiveServicesAccountsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Accounts client: %+v", err)
	}
	o.Configure(accountsClient.Client, o.Authorizers.ResourceManager)

	deploymentsClient, err := deployments.NewDeploymentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Deployments client: %+v", err)
	}
	o.Configure(deploymentsClient.Client, o.Authorizers.ResourceManager)

	raiBlobklistsClient, err := raiblocklists.NewRaiBlocklistsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Rai Blocklists client: %+v", err)
	}
	o.Configure(raiBlobklistsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AccountsClient:      accountsClient,
		DeploymentsClient:   deploymentsClient,
		RaiBlocklistsClient: raiBlobklistsClient,
	}, nil
}
