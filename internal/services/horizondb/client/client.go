// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/defaultparametergroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/firewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/passwordauthentications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ClustersClient                *clusters.ClustersClient
	DefaultParameterGroupsClient  *defaultparametergroups.DefaultParameterGroupsClient
	FirewallRulesClient           *firewallrules.FirewallRulesClient
	PasswordAuthenticationsClient *passwordauthentications.PasswordAuthenticationsClient
	PoolsClient                   *pools.PoolsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	clustersClient, err := clusters.NewClustersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Clusters client: %+v", err)
	}
	o.Configure(clustersClient.Client, o.Authorizers.ResourceManager)

	defaultParameterGroupsClient, err := defaultparametergroups.NewDefaultParameterGroupsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DefaultParameterGroups client: %+v", err)
	}
	o.Configure(defaultParameterGroupsClient.Client, o.Authorizers.ResourceManager)

	firewallRulesClient, err := firewallrules.NewFirewallRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FirewallRules client: %+v", err)
	}
	o.Configure(firewallRulesClient.Client, o.Authorizers.ResourceManager)

	passwordAuthenticationsClient, err := passwordauthentications.NewPasswordAuthenticationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building PasswordAuthentications client: %+v", err)
	}
	o.Configure(passwordAuthenticationsClient.Client, o.Authorizers.ResourceManager)

	poolsClient, err := pools.NewPoolsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Pools client: %+v", err)
	}
	o.Configure(poolsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ClustersClient:                clustersClient,
		DefaultParameterGroupsClient:  defaultParameterGroupsClient,
		FirewallRulesClient:           firewallRulesClient,
		PasswordAuthenticationsClient: passwordAuthenticationsClient,
		PoolsClient:                   poolsClient,
	}, nil
}
