// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2026-03-15/openapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2026-03-15/sqldedicatedgateway"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/firewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/roles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ClustersClient            *clusters.ClustersClient
	ConfigurationsClient      *configurations.ConfigurationsClient
	FirewallRulesClient       *firewallrules.FirewallRulesClient
	OpenapisClient            *openapis.OpenapisClient
	RolesClient               *roles.RolesClient
	SqlDedicatedGatewayClient *sqldedicatedgateway.SqlDedicatedGatewayClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	clustersClient, err := clusters.NewClustersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Clusters client: %+v", err)
	}
	o.Configure(clustersClient.Client, o.Authorizers.ResourceManager)

	configurationsClient, err := configurations.NewConfigurationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Configurations client: %+v", err)
	}
	o.Configure(configurationsClient.Client, o.Authorizers.ResourceManager)

	firewallRulesClient, err := firewallrules.NewFirewallRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FirewallRules client: %+v", err)
	}
	o.Configure(firewallRulesClient.Client, o.Authorizers.ResourceManager)

	openapisClient, err := openapis.NewOpenapisClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Openapis client: %+v", err)
	}
	o.Configure(openapisClient.Client, o.Authorizers.ResourceManager)

	rolesClient, err := roles.NewRolesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Roles client: %+v", err)
	}
	o.Configure(rolesClient.Client, o.Authorizers.ResourceManager)

	sqlDedicatedGatewayClient, err := sqldedicatedgateway.NewSqlDedicatedGatewayClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Sql Dedicated Gateway client: %+v", err)
	}
	o.Configure(sqlDedicatedGatewayClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ClustersClient:            clustersClient,
		ConfigurationsClient:      configurationsClient,
		FirewallRulesClient:       firewallRulesClient,
		OpenapisClient:            openapisClient,
		RolesClient:               rolesClient,
		SqlDedicatedGatewayClient: sqlDedicatedGatewayClient,
	}, nil
}
