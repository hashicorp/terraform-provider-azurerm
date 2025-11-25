// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/accountconnectionresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesprojects"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/deployments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/raiblocklists"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/raipolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AccountConnectionResourceClient *accountconnectionresource.AccountConnectionResourceClient
	AccountsClient                  *cognitiveservicesaccounts.CognitiveServicesAccountsClient
	DeploymentsClient               *deployments.DeploymentsClient
	ProjectsClient                  *cognitiveservicesprojects.CognitiveServicesProjectsClient
	RaiBlocklistsClient             *raiblocklists.RaiBlocklistsClient
	RaiPoliciesClient               *raipolicies.RaiPoliciesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	accountConnectionResourceClient, err := accountconnectionresource.NewAccountConnectionResourceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Account Connection client: %+v", err)
	}
	o.Configure(accountConnectionResourceClient.Client, o.Authorizers.ResourceManager)

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

	projectsClient, err := cognitiveservicesprojects.NewCognitiveServicesProjectsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Projects client: %+v", err)
	}
	o.Configure(projectsClient.Client, o.Authorizers.ResourceManager)

	raiPoliciesClient, err := raipolicies.NewRaiPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Rai Policies client: %+v", err)
	}
	o.Configure(raiPoliciesClient.Client, o.Authorizers.ResourceManager)

	raiBlobklistsClient, err := raiblocklists.NewRaiBlocklistsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Rai Blocklists client: %+v", err)
	}
	o.Configure(raiBlobklistsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AccountConnectionResourceClient: accountConnectionResourceClient,
		AccountsClient:                  accountsClient,
		DeploymentsClient:               deploymentsClient,
		ProjectsClient:                  projectsClient,
		RaiBlocklistsClient:             raiBlobklistsClient,
		RaiPoliciesClient:               raiPoliciesClient,
	}, nil
}
