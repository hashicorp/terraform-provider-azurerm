// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_agent

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func mssqlJobAgentSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Elastic Job Agent.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlJobAgentName,
		},

		"database_id": {
			Description:  "The ID of the database to store metadata for this Elastic Job Agent.",
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
			Description:  "The name of the SKU to use for this Elastic Job Agent.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      SqlJobAgentSkuJA100,
			ValidateFunc: validation.StringInSlice(PossibleValuesForJobAgentSku(), false),
		},

		"tags": commonschema.Tags(),
	}
}
