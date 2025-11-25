// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesprojects"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/deployments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/projectconnectionresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/raiblocklists"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/raipolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AccountsClient                  *cognitiveservicesaccounts.CognitiveServicesAccountsClient
	DeploymentsClient               *deployments.DeploymentsClient
	ProjectConnectionResourceClient *projectconnectionresource.ProjectConnectionResourceClient
	ProjectsClient                  *cognitiveservicesprojects.CognitiveServicesProjectsClient
	RaiBlocklistsClient             *raiblocklists.RaiBlocklistsClient
	RaiPoliciesClient               *raipolicies.RaiPoliciesClient
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

	projectConnectionResourceClient, err := projectconnectionresource.NewProjectConnectionResourceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Project Connection client: %+v", err)
	}
	o.Configure(projectConnectionResourceClient.Client, o.Authorizers.ResourceManager)

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
		AccountsClient:                  accountsClient,
		DeploymentsClient:               deploymentsClient,
		ProjectConnectionResourceClient: projectConnectionResourceClient,
		ProjectsClient:                  projectsClient,
		RaiBlocklistsClient:             raiBlobklistsClient,
		RaiPoliciesClient:               raiPoliciesClient,
	}, nil
}
