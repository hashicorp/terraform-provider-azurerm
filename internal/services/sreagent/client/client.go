// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-01-01/agents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	Agents *agents.AgentsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	client, err := agents.NewAgentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Agents client: %+v", err)
	}
	o.Configure(client.Client, o.Authorizers.ResourceManager)
	return &Client{Agents: client}, nil
}
