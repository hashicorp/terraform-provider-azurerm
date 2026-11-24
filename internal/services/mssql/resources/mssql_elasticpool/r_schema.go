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
			Description:  "The name of the elastic pool. This needs to be globally unique.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlElasticPoolName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"server_name": {
			Description:  "The name of the SQL Server on which to create the elastic pool.",
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
						Description: "The SKU Name for this Elasticpool. The name of the SKU, will be either `vCore` based or `DTU` based. Possible `DTU` based values are `BasicPool`, `StandardPool`, `PremiumPool` while possible `vCore` based values are `GP_Gen4`, `GP_Gen5`, `GP_Fsv2`, `GP_DC`, `BC_Gen4`, `BC_Gen5`, `BC_DC`, `HS_PRMS`, `HS_MOPRMS`, or `HS_Gen5`.",
						Type:        pluginsdk.TypeString,
						Required:    true,
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
						Description:  "The scale up/out capacity, representing server's compute units. For more information see the documentation for your Elasticpool configuration: [vCore-based](https://learn.microsoft.com/azure/sql-database/sql-database-vcore-resource-limits-elastic-pools) or [DTU-based](https://learn.microsoft.com/azure/sql-database/sql-database-dtu-resource-limits-elastic-pools).",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"tier": {
						Description: "The tier of the particular SKU. For more information see the documentation for your Elasticpool configuration: [vCore-based](https://learn.microsoft.com/azure/sql-database/sql-database-vcore-resource-limits-elastic-pools) or [DTU-based](https://learn.microsoft.com/azure/sql-database/sql-database-dtu-resource-limits-elastic-pools).",
						Type:        pluginsdk.TypeString,
						Required:    true,
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
						Description: "The `family` of hardware `Gen4`, `Gen5`, `Fsv2`, `MOPRMS`, `PRMS`, or `DC`.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
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
			Description:  "The name of the Public Maintenance Configuration window to apply to the elastic pool.",
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
						Description:  "The minimum capacity all databases are guaranteed.",
						Type:         pluginsdk.TypeFloat,
						Required:     true,
						ValidateFunc: validation.FloatAtLeast(0.0),
					},

					"max_capacity": {
						Description:  "The maximum capacity any one database can consume.",
						Type:         pluginsdk.TypeFloat,
						Required:     true,
						ValidateFunc: validation.FloatAtLeast(0.0),
					},
				},
			},
		},

		"max_size_bytes": {
			Description: "The max data size of the elastic pool in bytes.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			// Note: O+C because this value is computed based on `max_size_gb` if not set
			Computed:      true,
			ConflictsWith: []string{"max_size_gb"},
			ValidateFunc:  validation.IntAtLeast(0),
		},

		"max_size_gb": {
			Description: "The max data size of the elastic pool in gigabytes.",
			Type:        pluginsdk.TypeFloat,
			Optional:    true,
			// Note: O+C because this value is computed based on `max_size_bytes` if not set
			Computed:      true,
			ConflictsWith: []string{"max_size_bytes"},
			ValidateFunc:  validation.FloatAtLeast(0),
		},

		"zone_redundant": {
			Description: "Whether this elastic pool is zone redundant. `tier` needs to be `Premium` for `DTU` based or `BusinessCritical` for `vCore` based `sku`.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		// NOTE: The implementation of 'enclave_type' in the API differs slightly between database
		// and elasticpools. Database does not allow the 'Default' value to be passed for DW or
		// DC skus, where elasticpools allows 'Default' but will error if you try to set the
		// 'enclave_type' to 'VBS' for DC skus...
		"enclave_type": {
			Description:  "The type of enclave to be used by the elastic pool. When `enclave_type` is not specified (e.g., the default) enclaves are not enabled on the elastic pool. Once enabled (e.g., by specifying `Default` or `VBS`) removing the `enclave_type` field from the configuration file will force the creation of a new resource.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(databases.PossibleValuesForAlwaysEncryptedEnclaveType(), false),
		},

		"license_type": {
			Description:  "The license type applied to this database.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(elasticpools.PossibleValuesForElasticPoolLicenseType(), false),
		},

		"high_availability_replica_count": {
			Description: "The number of high availability replicas for the elastic pool.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			// NOTE: O+C can only be set for Hyperscale skus, which have a default value of 1
			Computed:     true,
			ValidateFunc: validation.IntBetween(0, 4),
		},

		"tags": commonschema.Tags(),
	}
}
