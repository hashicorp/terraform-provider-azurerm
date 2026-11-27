// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_postgresql_firewall_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for the Azure Cosmos DB for PostgreSQL Firewall Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FirewallRuleName,
		},

		"cluster_id": {
			Description:  "The resource ID of the Azure Cosmos DB for PostgreSQL Cluster.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: firewallrules.ValidateServerGroupsv2ID,
		},

		"end_ip_address": {
			Description:  "The end IP address of the Azure Cosmos DB for PostgreSQL Firewall Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsIPv4Address,
		},

		"start_ip_address": {
			Description:  "The start IP address of the Azure Cosmos DB for PostgreSQL Firewall Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsIPv4Address,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
