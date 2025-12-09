// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2025-09-01/firewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2025-09-01/mongoclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2025-09-01/users"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	FirewallRulesClient *firewallrules.FirewallRulesClient
	MongoClustersClient *mongoclusters.MongoClustersClient
	UsersClient         *users.UsersClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	firewallRulesClient, err := firewallrules.NewFirewallRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FirewallRules client: %+v", err)
	}
	o.Configure(firewallRulesClient.Client, o.Authorizers.ResourceManager)

	mongoClustersClient, err := mongoclusters.NewMongoClustersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building MongoClusters client: %+v", err)
	}
	o.Configure(mongoClustersClient.Client, o.Authorizers.ResourceManager)

	usersClient, err := users.NewUsersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Users client: %+v", err)
	}
	o.Configure(usersClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		FirewallRulesClient: firewallRulesClient,
		MongoClustersClient: mongoClustersClient,
		UsersClient:         usersClient,
	}, nil
}
