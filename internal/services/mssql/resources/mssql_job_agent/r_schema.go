// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_agent

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func mssqlJobAgentSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlJobAgentName,
		},

		"database_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSqlDatabaseID,
		},

		"location": commonschema.Location(),

		"identity": commonschema.UserAssignedIdentityOptional(),

		// This is a top level argument rather than a block because while Azure accepts input for both sku name and capacity fields,
		// the capacity must always be equal to the number included in the sku name.
		"sku": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      helper.SqlJobAgentSkuJA100,
			ValidateFunc: validation.StringInSlice(helper.PossibleValuesForJobAgentSku(), false),
		},

		"tags": commonschema.Tags(),
	}
}
