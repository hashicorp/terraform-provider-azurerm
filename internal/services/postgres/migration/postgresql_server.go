// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ pluginsdk.StateUpgrade = PostgresqlServerV0ToV1{}

type PostgresqlServerV0ToV1 struct{}

func (PostgresqlServerV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ServerName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"version": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(servers.PossibleValuesForServerVersion(), false),
		},

		"administrator_login": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.All(validation.StringIsNotWhiteSpace, validate.AdminUsernames),
		},

		"administrator_login_password": {
			Type:      pluginsdk.TypeString,
			Optional:  true,
			Sensitive: true,
		},

		"auto_grow_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"backup_retention_days": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(7, 35),
		},

		"geo_redundant_backup_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"create_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(servers.CreateModeDefault),
			ValidateFunc: validation.StringInSlice(servers.PossibleValuesForCreateMode(), false),
		},

		"creation_source_server_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: servers.ValidateServerID,
		},

		"identity": commonschema.SystemAssignedIdentityOptional(),

		"infrastructure_encryption_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"restore_point_in_time": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"storage_mb": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.All(
				validation.IntBetween(5120, 16777216),
				validation.IntDivisibleBy(1024),
			),
		},

		"ssl_minimal_tls_version_enforced": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(servers.MinimalTlsVersionEnumTLSOneTwo),
			ValidateFunc: validation.StringInSlice(servers.PossibleValuesForMinimalTlsVersionEnum(), false),
		},

		"ssl_enforcement_enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},

		"threat_detection_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						AtLeastOneOf: []string{
							"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
							"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
							"threat_detection_policy.0.storage_endpoint",
						},
					},

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
								"Data_Exfiltration",
								"Unsafe_Action",
							}, false),
						},
						AtLeastOneOf: []string{
							"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
							"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
							"threat_detection_policy.0.storage_endpoint",
						},
					},

					"email_account_admins": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						AtLeastOneOf: []string{
							"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
							"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
							"threat_detection_policy.0.storage_endpoint",
						},
					},

					"email_addresses": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							// todo email validation in code
						},
						Set: pluginsdk.HashString,
						AtLeastOneOf: []string{
							"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
							"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
							"threat_detection_policy.0.storage_endpoint",
						},
					},

					"retention_days": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
						AtLeastOneOf: []string{
							"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
							"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
							"threat_detection_policy.0.storage_endpoint",
						},
					},

					"storage_account_access_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{
							"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
							"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
							"threat_detection_policy.0.storage_endpoint",
						},
					},

					"storage_endpoint": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{
							"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
							"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
							"threat_detection_policy.0.storage_endpoint",
						},
					},
				},
			},
		},

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}

func (PostgresqlServerV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old
		//  /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBForPostgreSQL/servers/{serverName}
		// new:
		//  /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/servers/{serverName}
		// summary:
		// Check for `For` and swap to `for`
		oldId := rawState["id"].(string)
		if strings.Contains(oldId, "Microsoft.DBForPostgreSQL") {
			modifiedId := strings.ReplaceAll(oldId, "Microsoft.DBForPostgreSQL", "Microsoft.DBforPostgreSQL")

			newId, err := servers.ParseServerID(modifiedId)
			if err != nil {
				return rawState, err
			}
			log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
			rawState["id"] = newId.ID()
		}

		return rawState, nil
	}
}
