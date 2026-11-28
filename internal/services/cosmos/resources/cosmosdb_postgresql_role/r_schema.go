// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_postgresql_role

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/roles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Azure Cosmos DB for PostgreSQL Role.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RoleName,
		},

		"cluster_id": {
			Description:  "The resource ID of the Azure Cosmos DB for PostgreSQL Cluster.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: roles.ValidateServerGroupsv2ID,
		},

		"password": {
			Description:  "The password of the Azure Cosmos DB for PostgreSQL Role.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringLenBetween(8, 256),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
