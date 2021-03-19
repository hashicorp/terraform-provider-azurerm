package migration

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
)

// Default:  string(sql.CreateModeDefault),

func DatabaseV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		Version: 0,
		Type:    databaseV0V1Schema().CoreConfigSchema().ImpliedType(),
		Upgrade: databaseUpgradeV0ToV1,
	}
}

func databaseV0V1Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"auto_pause_delay_in_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"create_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"collation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"elastic_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"extended_auditing_policy": {
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				Deprecated: "the `extended_auditing_policy` block has been moved to `azurerm_mssql_server_extended_auditing_policy` and `azurerm_mssql_database_extended_auditing_policy`. This block will be removed in version 3.0 of the provider.",
				ConfigMode: schema.SchemaConfigModeAttr,
				MaxItems:   1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_account_access_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},

						"storage_endpoint": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"storage_account_access_key_is_secondary": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"retention_in_days": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"log_monitoring_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},

			"license_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"long_term_retention_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// WeeklyRetention - The weekly retention policy for an LTR backup in an ISO 8601 format.
						"weekly_retention": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						// MonthlyRetention - The monthly retention policy for an LTR backup in an ISO 8601 format.
						"monthly_retention": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						// YearlyRetention - The yearly retention policy for an LTR backup in an ISO 8601 format.
						"yearly_retention": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						// WeekOfYear - The week of year to take the yearly backup in an ISO 8601 format.
						"week_of_year": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"short_term_retention_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"retention_days": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"max_size_gb": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"min_capacity": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},

			"restore_point_in_time": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
			},

			"recover_database_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"restore_dropped_database_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"read_replica_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"read_scale": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"sample_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"sku_name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"creation_source_database_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"storage_account_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "GRS",
			},

			"zone_redundant": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"threat_detection_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disabled_alerts": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      schema.HashString,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"Sql_Injection",
									"Sql_Injection_Vulnerability",
									"Access_Anomaly",
								}, true),
							},
						},

						"email_account_admins": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							Default:          "Disabled",
						},

						"email_addresses": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: schema.HashString,
						},

						"retention_days": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"state": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							Default:          "Disabled",
						},

						"storage_account_access_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},

						"storage_endpoint": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"use_server_default": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							Default:          "Disabled",
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func databaseUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Upgrading from Database V0 to V1..")
	existing := rawState["create_mode"]
	if existing == nil {
		log.Printf("[DEBUG] Setting `create_mode` to `Default`")
		rawState["create_mode"] = "Default"
	}

	log.Printf("[DEBUG] Upgraded from Database V0 to V1..")
	return rawState, nil
}
