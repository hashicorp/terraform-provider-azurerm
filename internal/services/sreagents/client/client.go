// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-07-01/agents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AgentsClient *agents.AgentsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	agentsClient, err := agents.NewAgentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SRE Agents client: %+v", err)
	}
	o.Configure(agentsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AgentsClient: agentsClient,
	}, nil
}
