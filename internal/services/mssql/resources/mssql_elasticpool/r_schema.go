// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_elasticpool

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/elasticpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func mssqlElasticpoolSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlElasticPoolName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"server_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlServerName,
		},

		"sku": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"BasicPool",
							"StandardPool",
							"PremiumPool",
							"GP_Gen4",
							"GP_Gen5",
							"GP_Fsv2",
							"GP_DC",
							"BC_Gen4",
							"BC_Gen5",
							"BC_DC",
							"HS_Gen5",
							"HS_PRMS",
							"HS_MOPRMS",
						}, false),
					},

					"capacity": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"tier": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Basic",
							"Standard",
							"Premium",
							"GeneralPurpose",
							"BusinessCritical",
							"Hyperscale",
						}, false),
					},

					"family": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Gen4",
							"Gen5",
							"Fsv2",
							"DC",
							"MOPRMS",
							"PRMS",
						}, false),
					},
				},
			},
		},

		"maintenance_configuration_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "SQL_Default",
			ValidateFunc: validation.StringInSlice(helpers.ResourceMsSqlDatabaseMaintenanceNames(), false),
		},

		"per_database_settings": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"min_capacity": {
						Type:         pluginsdk.TypeFloat,
						Required:     true,
						ValidateFunc: validation.FloatAtLeast(0.0),
					},

					"max_capacity": {
						Type:         pluginsdk.TypeFloat,
						Required:     true,
						ValidateFunc: validation.FloatAtLeast(0.0),
					},
				},
			},
		},

		"max_size_bytes": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// Note: O+C because this value is computed based on `max_size_gb` if not set
			Computed:      true,
			ConflictsWith: []string{"max_size_gb"},
			ValidateFunc:  validation.IntAtLeast(0),
		},

		"max_size_gb": {
			Type:     pluginsdk.TypeFloat,
			Optional: true,
			// Note: O+C because this value is computed based on `max_size_bytes` if not set
			Computed:      true,
			ConflictsWith: []string{"max_size_bytes"},
			ValidateFunc:  validation.FloatAtLeast(0),
		},

		"zone_redundant": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		// NOTE: The implementation of 'enclave_type' in the API differs slightly between database
		// and elasticpools. Database does not allow the 'Default' value to be passed for DW or
		// DC skus, where elasticpools allows 'Default' but will error if you try to set the
		// 'enclave_type' to 'VBS' for DC skus...
		"enclave_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(databases.PossibleValuesForAlwaysEncryptedEnclaveType(), false),
		},

		"license_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(elasticpools.PossibleValuesForElasticPoolLicenseType(), false),
		},

		"high_availability_replica_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C can only be set for Hyperscale skus, which have a default value of 1
			Computed:     true,
			ValidateFunc: validation.IntBetween(0, 4),
		},

		"tags": commonschema.Tags(),
	}
}
