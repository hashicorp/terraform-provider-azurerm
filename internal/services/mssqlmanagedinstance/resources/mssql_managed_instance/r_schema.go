// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlManagedInstanceServerName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku_name": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"BC_Gen4",
				"BC_Gen5",
				"BC_Gen8IH",
				"BC_Gen8IM",
				"GP_Gen4",
				"GP_Gen5",
				"GP_Gen8IH",
				"GP_Gen8IM",
			}, false),
		},

		"license_type": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"LicenseIncluded",
				"BasePrice",
			}, true),
		},

		"storage_size_in_gb": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(32, 32768),
		},

		"subnet_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"vcores": {
			Type:     schema.TypeInt,
			Required: true,
			ValidateFunc: validation.IntInSlice([]int{
				4,
				6,
				8,
				10,
				12,
				16,
				20,
				24,
				32,
				40,
				48,
				56,
				64,
				80,
				96,
				128,
			}),
		},

		"administrator_login": {
			Type:     schema.TypeString,
			Optional: true,
			// Note: O+C because Azure returns a generated value if `azure_active_directory_administrator.azuread_authentication_only_enabled` is `true`.
			// which leads to unnecessary resource recreation on subsequent plans where Terraform tries to remove it.
			Computed:     true,
			ForceNew:     true,
			AtLeastOneOf: []string{"administrator_login", "azure_active_directory_administrator"},
			RequiredWith: []string{"administrator_login", "administrator_login_password"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"administrator_login_password": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			AtLeastOneOf: []string{"administrator_login_password", "azure_active_directory_administrator"},
			RequiredWith: []string{"administrator_login", "administrator_login_password"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"azure_active_directory_administrator": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"login_username": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"object_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},

					"principal_type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(managedinstances.PossibleValuesForPrincipalType(), false),
					},

					"azuread_authentication_only_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"tenant_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},

		"collation": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "SQL_Latin1_General_CP1_CI_AS",
			ValidateFunc: validation.StringIsNotEmpty,
			ForceNew:     true,
		},

		"database_format": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      managedinstances.ManagedInstanceDatabaseFormatSQLServerTwoZeroTwoTwo,
			ValidateFunc: validation.StringInSlice(managedinstances.PossibleValuesForManagedInstanceDatabaseFormat(), false),
		},

		"dns_zone_partner_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseSqlManagedInstanceIDInsensitively),
		},

		"hybrid_secondary_usage": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      managedinstances.HybridSecondaryUsageActive,
			ValidateFunc: validation.StringInSlice(managedinstances.PossibleValuesForHybridSecondaryUsage(), false),
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"general_purpose_v2_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"maintenance_configuration_name": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "SQL_Default",
			ValidateFunc: validation.Any(
				validation.StringInSlice([]string{
					"SQL_Default",
				}, false),
				validation.StringMatch(regexp.MustCompile(`^SQL_[A-Za-z0-9]+_MI_\d+$`), "expected a name in the format `SQL_{Location}_MI_{Number}` or `SQL_Default`"),
			),
		},

		"minimum_tls_version": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "1.2",
			ValidateFunc: validation.StringInSlice([]string{
				"1.2",
			}, false),
		},

		"proxy_override": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  string(managedinstances.ManagedInstanceProxyOverrideRedirect),
			ValidateFunc: validation.StringInSlice([]string{
				string(managedinstances.ManagedInstanceProxyOverrideRedirect),
				string(managedinstances.ManagedInstanceProxyOverrideProxy),
			}, false),
		},

		"public_data_endpoint_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"service_principal_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(managedinstances.ServicePrincipalTypeSystemAssigned),
			}, false),
		},

		"storage_account_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  StorageAccountTypeGRS,
			ValidateFunc: validation.StringInSlice([]string{
				StorageAccountTypeGRS,
				StorageAccountTypeLRS,
				StorageAccountTypeZRS,
				StorageAccountTypeGZRS,
			}, false),
		},

		"storage_iops": {
			Type:     schema.TypeInt,
			Optional: true,
			// NOTE: O+C - Azure returns a calculated IOPS value for GPv2 instances when `storage_iops` is omitted.
			Computed:     true,
			ValidateFunc: validation.IntBetween(300, 80000),
		},

		"timezone_id": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "UTC",
			ValidateFunc: validation.StringIsNotEmpty,
			ForceNew:     true,
		},

		"zone_redundant_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dns_zone": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"fqdn": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			// dns_zone_partner_id can only be set on init
			if oldVal, newVal := rd.GetChange("dns_zone_partner_id"); oldVal.(string) == "" && newVal.(string) != "" {
				if err := rd.ForceNew("dns_zone_partner_id"); err != nil {
					return err
				}
			}

			// system-assigned identity can't be removed due to https://github.com/Azure/azure-rest-api-specs/issues/16838
			if oldVal, newVal := rd.GetChange("identity.#"); oldVal.(int) == 1 && newVal.(int) == 0 {
				if err := rd.ForceNew("identity"); err != nil {
					return err
				}
			}

			// Once the `AlwaysUpToDate` is enabled, you can't go back to `SQLServer2022` update policy.
			if oldVal, newVal := rd.GetChange("database_format"); oldVal.(string) == string(managedinstances.ManagedInstanceDatabaseFormatAlwaysUpToDate) && newVal.(string) == string(managedinstances.ManagedInstanceDatabaseFormatSQLServerTwoZeroTwoTwo) {
				if err := rd.ForceNew("database_format"); err != nil {
					return err
				}
			}

			_, aadAdminOk := rd.GetOk("azure_active_directory_administrator")
			authOnlyEnabled := rd.Get("azure_active_directory_administrator.0.azuread_authentication_only_enabled").(bool)
			rawConfig := rd.GetRawConfig().AsValueMap()
			adminLogin := rawConfig["administrator_login"]
			adminPassword := rawConfig["administrator_login_password"]

			if aadAdminOk && !authOnlyEnabled && (adminLogin.IsNull() || adminPassword.IsNull()) {
				return fmt.Errorf("`administrator_login` and `administrator_login_password` are required when `azuread_authentication_only_enabled` is false")
			}

			// Unknown values can come from expressions and may later resolve to null, so only
			// validate `storage_iops` rules when the argument is known to be explicitly configured.
			if storageIOps, ok := rawConfig["storage_iops"]; ok && storageIOps.IsKnown() && !storageIOps.IsNull() {
				if !rd.Get("general_purpose_v2_enabled").(bool) {
					return errors.New("`storage_iops` can only be set when `general_purpose_v2_enabled` is `true`")
				}

				if sku := rd.Get("sku_name").(string); strings.HasPrefix(sku, "BC_") {
					return fmt.Errorf("`storage_iops` is not supported on Business Critical SKUs, got SKU `%s`", sku)
				}
			}

			if sku := rd.Get("sku_name").(string); strings.HasPrefix(sku, "BC_") && rd.Get("general_purpose_v2_enabled").(bool) {
				return fmt.Errorf("`general_purpose_v2_enabled` cannot be set to `true` on Business Critical SKUs, got SKU `%s`", sku)
			}

			// Zone redundancy is not available for Next-gen General Purpose instances.
			// https://learn.microsoft.com/azure/azure-sql/managed-instance/high-availability-sla-local-zone-redundancy#next-gen-general-purpose-service-tier
			if rd.Get("zone_redundant_enabled").(bool) && rd.Get("general_purpose_v2_enabled").(bool) {
				return errors.New("`zone_redundant_enabled` cannot be set to `true` when `general_purpose_v2_enabled` is `true`")
			}

			return nil
		},
	}
}
