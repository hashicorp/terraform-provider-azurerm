package migration

import (
	"context"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

// Default:  string(sql.CreateModeDefault),

var _ pluginsdk.StateUpgrade = DatabaseV0ToV1{}

type DatabaseV0ToV1 struct {
}

func (d DatabaseV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"server_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"auto_pause_delay_in_minutes": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"create_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Computed: true,
		},

		"collation": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"elastic_pool_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"extended_auditing_policy": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			Computed:   true,
			Deprecated: "the `extended_auditing_policy` block has been moved to `azurerm_mssql_server_extended_auditing_policy` and `azurerm_mssql_database_extended_auditing_policy`. This block will be removed in version 3.0 of the provider.",
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			MaxItems:   1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"storage_account_access_key": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						Sensitive: true,
					},

					"storage_endpoint": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"storage_account_access_key_is_secondary": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"retention_in_days": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"log_monitoring_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},

		"license_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		//lintignore:XS003
		"long_term_retention_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					// WeeklyRetention - The weekly retention policy for an LTR backup in an ISO 8601 format.
					"weekly_retention": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					// MonthlyRetention - The monthly retention policy for an LTR backup in an ISO 8601 format.
					"monthly_retention": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					// YearlyRetention - The yearly retention policy for an LTR backup in an ISO 8601 format.
					"yearly_retention": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					// WeekOfYear - The week of year to take the yearly backup in an ISO 8601 format.
					"week_of_year": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
					},
				},
			},
		},

		"short_term_retention_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"retention_days": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
				},
			},
		},

		"max_size_gb": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"min_capacity": {
			Type:     pluginsdk.TypeFloat,
			Optional: true,
			Computed: true,
		},

		"restore_point_in_time": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
		},

		"recover_database_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"restore_dropped_database_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"read_replica_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"read_scale": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"sample_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"sku_name": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"creation_source_database_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Computed: true,
		},

		"storage_account_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "GRS",
		},

		"zone_redundant": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"threat_detection_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disabled_alerts": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Set:      pluginsdk.HashString,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"Sql_Injection",
								"Sql_Injection_Vulnerability",
								"Access_Anomaly",
							}, true),
						},
					},

					"email_account_admins": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						Default:          "Disabled",
					},

					"email_addresses": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Set: pluginsdk.HashString,
					},

					"retention_days": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"state": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						Default:          "Disabled",
					},

					"storage_account_access_key": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						Sensitive: true,
					},

					"storage_endpoint": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"use_server_default": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						Default:          "Disabled",
					},
				},
			},
		},

		"tags": tags.Schema(),
	}
}

func (d DatabaseV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		log.Printf("[DEBUG] Upgrading from Database V0 to V1..")
		existing := rawState["create_mode"]
		if existing == nil {
			log.Printf("[DEBUG] Setting `create_mode` to `Default`")
			rawState["create_mode"] = "Default"
		}

		log.Printf("[DEBUG] Upgraded from Database V0 to V1..")
		return rawState, nil
	}
}
