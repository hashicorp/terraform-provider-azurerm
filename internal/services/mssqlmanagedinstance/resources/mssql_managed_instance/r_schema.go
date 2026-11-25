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
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	arguments := map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the SQL Managed Instance. This needs to be globally unique within Azure.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlManagedInstanceServerName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku_name": {
			Description: "The SKU Name for the SQL Managed Instance.",
			Type:        schema.TypeString,
			Required:    true,
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
			Description: "What type of license the Managed Instance will use.",
			Type:        schema.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"LicenseIncluded",
				"BasePrice",
			}, true),
		},

		"storage_size_in_gb": {
			Description:  "Maximum storage space for the SQL Managed instance. This should be a multiple of 32 (GB).",
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(32, 32768),
		},

		"subnet_id": {
			Description:  "The subnet resource ID that the SQL Managed Instance will be associated with.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"vcores": {
			Description: "Number of cores that should be assigned to the SQL Managed Instance.",
			Type:        schema.TypeInt,
			Required:    true,
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
			Description: "The administrator login name for the new SQL Managed Instance.",
			Type:        schema.TypeString,
			Optional:    true,
			// Note: O+C because Azure returns a generated value if `azure_active_directory_administrator.azuread_authentication_only_enabled` is `true`.
			// which leads to unnecessary resource recreation on subsequent plans where Terraform tries to remove it.
			Computed:     true,
			ForceNew:     true,
			AtLeastOneOf: []string{"administrator_login", "azure_active_directory_administrator"},
			RequiredWith: []string{"administrator_login", "administrator_login_password"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"administrator_login_password": {
			Description:  "The password associated with the `administrator_login` user. Needs to comply with Azure's [Password Policy](https://msdn.microsoft.com/library/ms161959.aspx).",
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
						Description:  "The login username of the Azure AD Administrator of this SQL Managed Instance.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"object_id": {
						Description:  "The object ID of the Azure AD Administrator of this SQL Managed Instance.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},

					"principal_type": {
						Description:  "The principal type of the Azure AD Administrator of this SQL Managed Instance.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(managedinstances.PossibleValuesForPrincipalType(), false),
					},

					"azuread_authentication_only_enabled": {
						Description: "Whether only Azure AD authentication can be used to log in to this SQL Managed Instance. When `true`, the `administrator_login` and `administrator_login_password` properties can be omitted.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"tenant_id": {
						Description:  "The tenant ID of the Azure AD Administrator of this SQL Managed Instance.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},

		"collation": {
			Description:  "How the SQL Managed Instance will be collated.",
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "SQL_Latin1_General_CP1_CI_AS",
			ValidateFunc: validation.StringIsNotEmpty,
			ForceNew:     true,
		},

		"database_format": {
			Description:  "The internal format of the SQL Managed Instance databases specific to the SQL engine version.",
			Type:         schema.TypeString,
			Optional:     true,
			Default:      managedinstances.ManagedInstanceDatabaseFormatSQLServerTwoZeroTwoTwo,
			ValidateFunc: validation.StringInSlice(managedinstances.PossibleValuesForManagedInstanceDatabaseFormat(), false),
		},

		"dns_zone_partner_id": {
			Description:  "The ID of the SQL Managed Instance which will share the DNS zone. This is a prerequisite for creating an `azurerm_mssql_managed_instance_failover_group`. Setting this after creation forces a new resource to be created.",
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseSqlManagedInstanceIDInsensitively),
		},

		"hybrid_secondary_usage": {
			Description:  "The hybrid secondary usage for disaster recovery of the SQL Managed Instance.",
			Type:         schema.TypeString,
			Optional:     true,
			Default:      managedinstances.HybridSecondaryUsageActive,
			ValidateFunc: validation.StringInSlice(managedinstances.PossibleValuesForHybridSecondaryUsage(), false),
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"general_purpose_v2_enabled": {
			Description: "Whether the SQL Managed Instance should use the Next-gen General Purpose service tier.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"maintenance_configuration_name": {
			Description: "The name of the Public Maintenance Configuration window to apply to the SQL Managed Instance.",
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "SQL_Default",
			ValidateFunc: validation.Any(
				validation.StringInSlice([]string{
					"SQL_Default",
				}, false),
				validation.StringMatch(regexp.MustCompile(`^SQL_[A-Za-z0-9]+_MI_\d+$`), "expected a name in the format `SQL_{Location}_MI_{Number}` or `SQL_Default`"),
			),
		},

		"minimum_tls_version": {
			Description: "The Minimum TLS Version.",
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "1.2",
			ValidateFunc: validation.StringInSlice([]string{
				"1.2",
			}, false),
		},

		"pricing_model": {
			Description:  "The pricing model for the SQL Managed Instance.",
			Type:         schema.TypeString,
			Optional:     true,
			Default:      string(managedinstances.PricingModelRegular),
			ValidateFunc: validation.StringInSlice(managedinstances.PossibleValuesForPricingModel(), false),
		},

		"proxy_override": {
			Description: "How the SQL Managed Instance will be accessed.",
			Type:        schema.TypeString,
			Optional:    true,
			Default:     string(managedinstances.ManagedInstanceProxyOverrideRedirect),
			ValidateFunc: validation.StringInSlice([]string{
				string(managedinstances.ManagedInstanceProxyOverrideRedirect),
				string(managedinstances.ManagedInstanceProxyOverrideProxy),
			}, false),
		},

		"public_data_endpoint_enabled": {
			Description: "Whether the public data endpoint is enabled.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"service_principal_type": {
			Description: "The Service Principal type.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(managedinstances.ServicePrincipalTypeSystemAssigned),
			}, false),
		},

		"storage_account_type": {
			Description: "The Storage Account type used to store backups for this database.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     StorageAccountTypeGRS,
			ValidateFunc: validation.StringInSlice([]string{
				StorageAccountTypeGRS,
				StorageAccountTypeLRS,
				StorageAccountTypeZRS,
				StorageAccountTypeGZRS,
			}, false),
		},

		"storage_iops": {
			Description: "The storage IOPS for the SQL Managed Instance. This can only be specified when `general_purpose_v2_enabled` is `true`.",
			Type:        schema.TypeInt,
			Optional:    true,
			// NOTE: O+C - Azure returns a calculated IOPS value for GPv2 instances when `storage_iops` is omitted.
			Computed:     true,
			ValidateFunc: validation.IntBetween(300, 80000),
		},

		"timezone_id": {
			Description:  "The TimeZone ID that the SQL Managed Instance will be operating in.",
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "UTC",
			ValidateFunc: validation.StringIsNotEmpty,
			ForceNew:     true,
		},

		"zone_redundant_enabled": {
			Description: "Whether the SQL Managed Instance is zone redundant.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"tags": commonschema.Tags(),
	}

	if !features.SixPointOh() {
		arguments["pricing_model"] = &pluginsdk.Schema{
			Description: "The pricing model for the SQL Managed Instance.",
			Type:        schema.TypeString,
			Optional:    true,
			// NOTE: O+C to preserve existing free instances when omitted, avoiding an unintended upgrade to paid pricing in 5.x.
			Computed:     true,
			ValidateFunc: validation.StringInSlice(managedinstances.PossibleValuesForPricingModel(), false),
		}
	}

	return arguments
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dns_zone": {
			Description: "The DNS Zone where the SQL Managed Instance is located.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"fqdn": {
			Description: "The fully qualified domain name of the Azure Managed SQL Instance.",
			Type:        schema.TypeString,
			Computed:    true,
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
			// check for non-zero values as a null value is stored in state as `0`, attribute validation prevents a user from explicitly defniing this as `0`
			if storageIOps, ok := rawConfig["storage_iops"]; ok && storageIOps.IsKnown() && !storageIOps.IsNull() && rd.Get("storage_iops") != 0 {
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

			oldPricingModel, newPricingModel := rd.GetChange("pricing_model")
			if rd.NewValueKnown("pricing_model") && oldPricingModel.(string) == string(managedinstances.PricingModelRegular) && newPricingModel.(string) == string(managedinstances.PricingModelFreemium) {
				if err := rd.ForceNew("pricing_model"); err != nil {
					return err
				}
			}

			if newPricingModel.(string) == string(managedinstances.PricingModelFreemium) {
				if sku := rd.Get("sku_name").(string); rd.NewValueKnown("sku_name") && sku != "GP_Gen5" {
					return fmt.Errorf("`sku_name` must be `GP_Gen5` when `pricing_model` is `Freemium`, got `%s`", sku)
				}

				vcores := rd.Get("vcores").(int)
				if rd.NewValueKnown("vcores") && vcores != 4 && vcores != 8 {
					return fmt.Errorf("`vcores` must be `4` or `8` when `pricing_model` is `Freemium`, got `%d`", vcores)
				}

				storageSizeInGb := rd.Get("storage_size_in_gb").(int)
				if rd.NewValueKnown("storage_size_in_gb") && storageSizeInGb != 64 {
					return fmt.Errorf("`storage_size_in_gb` must be `64` when `pricing_model` is `Freemium`, got `%d`", storageSizeInGb)
				}

				storageAccountType := rd.Get("storage_account_type").(string)
				if rd.NewValueKnown("storage_account_type") && storageAccountType != StorageAccountTypeLRS {
					return fmt.Errorf("`storage_account_type` must be `LRS` when `pricing_model` is `Freemium`, got `%s`", storageAccountType)
				}

				if rd.NewValueKnown("zone_redundant_enabled") && rd.Get("zone_redundant_enabled").(bool) {
					return errors.New("`zone_redundant_enabled` must be `false` when `pricing_model` is `Freemium`")
				}
			}

			return nil
		},
	}
}
