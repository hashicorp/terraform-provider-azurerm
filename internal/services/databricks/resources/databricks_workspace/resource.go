// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks_workspace

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -test-params "premium"

func RegisterResource() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceDatabricksWorkspaceCreate,
		Read:   resourceDatabricksWorkspaceRead,
		Update: resourceDatabricksWorkspaceUpdate,
		Delete: resourceDatabricksWorkspaceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&workspaces.WorkspaceId{}),
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&workspaces.WorkspaceId{}),
		},

		Schema: databricksWorkspaceSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
				_, customerEncryptionEnabled := d.GetChange("customer_managed_key_enabled")
				_, defaultStorageFirewallEnabled := d.GetChange("default_storage_firewall_enabled")
				_, infrastructureEncryptionEnabled := d.GetChange("infrastructure_encryption_enabled")
				_, publicNetworkAccess := d.GetChange("public_network_access_enabled")
				_, requireNsgRules := d.GetChange("network_security_group_rules_required")
				_, backendPool := d.GetChange("load_balancer_backend_address_pool_id")
				_, managedServicesCMK := d.GetChange("managed_services_cmk_key_vault_key_id")
				_, managedDiskCMK := d.GetChange("managed_disk_cmk_key_vault_key_id")
				_, enhancedSecurityCompliance := d.GetChange("enhanced_security_compliance")

				oldSku, newSku := d.GetChange("sku")

				// Disabling Public Network Access means that this is a Private Endpoint Workspace
				// Having a Load Balancer Backend Address Pool means the this is a Secure Cluster Connectivity Workspace
				// You cannot have a Private Endpoint Workspace and a Secure Cluster Connectivity Workspace definitions in
				// the same workspace configuration...
				if !publicNetworkAccess.(bool) {
					if requireNsgRules.(string) == string(workspaces.RequiredNsgRulesAllRules) {
						return fmt.Errorf("having `network_security_group_rules_required` set to %q and `public_network_access_enabled` set to `false` is an invalid configuration", string(workspaces.RequiredNsgRulesAllRules))
					}
					if backendPool.(string) != "" {
						return fmt.Errorf("having `load_balancer_backend_address_pool_id` defined and having `public_network_access_enabled` set to `false` is an invalid configuration")
					}
				}

				if d.HasChange("sku") {
					if newSku == "trial" {
						log.Printf("[DEBUG] recreate databricks workspace, cannot be migrated to %s", newSku)
						d.ForceNew("sku")
					} else {
						log.Printf("[DEBUG] databricks workspace can be upgraded from %s to %s", oldSku, newSku)
					}
				}

				if (customerEncryptionEnabled.(bool) || defaultStorageFirewallEnabled.(bool) || len(enhancedSecurityCompliance.([]any)) > 0 || infrastructureEncryptionEnabled.(bool) || managedServicesCMK.(string) != "" || managedDiskCMK.(string) != "") && !strings.EqualFold("premium", newSku.(string)) {
					return fmt.Errorf("`customer_managed_key_enabled`, `default_storage_firewall_enabled`, `enhanced_security_compliance`, `infrastructure_encryption_enabled`, `managed_disk_cmk_key_vault_key_id` and `managed_services_cmk_key_vault_key_id` are only available with a `premium` workspace `sku`, got %q", newSku)
				}

				return nil
			}),

			// Once compliance security profile has been enabled, disabling it will force a workspace replacement
			pluginsdk.ForceNewIfChange("enhanced_security_compliance.0.compliance_security_profile_enabled", func(ctx context.Context, old, new, meta any) bool {
				return old.(bool) && !new.(bool)
			}),

			// Once a compliance standard is enabled, disabling it will force a workspace replacement
			pluginsdk.ForceNewIfChange("enhanced_security_compliance.0.compliance_security_profile_standards", func(ctx context.Context, old, new, meta any) bool {
				removedStandards := old.(*pluginsdk.Set).Difference(new.(*pluginsdk.Set))
				return removedStandards.Len() > 0
			}),

			// Compliance security profile requires automatic cluster update and enhanced security monitoring to be enabled
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
				_, complianceSecurityProfileEnabled := d.GetChange("enhanced_security_compliance.0.compliance_security_profile_enabled")
				_, automaticClusterUpdateEnabled := d.GetChange("enhanced_security_compliance.0.automatic_cluster_update_enabled")
				_, enhancedSecurityMonitoringEnabled := d.GetChange("enhanced_security_compliance.0.enhanced_security_monitoring_enabled")

				if complianceSecurityProfileEnabled.(bool) && (!automaticClusterUpdateEnabled.(bool) || !enhancedSecurityMonitoringEnabled.(bool)) {
					return fmt.Errorf("`automatic_cluster_update_enabled` and `enhanced_security_monitoring_enabled` must be set to true when `compliance_security_profile_enabled` is set to true")
				}

				return nil
			}),

			// compliance standards cannot be specified without enabling compliance profile
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
				_, complianceSecurityProfileEnabled := d.GetChange("enhanced_security_compliance.0.compliance_security_profile_enabled")
				_, complianceStandards := d.GetChange("enhanced_security_compliance.0.compliance_security_profile_standards")

				if !complianceSecurityProfileEnabled.(bool) && complianceStandards.(*pluginsdk.Set).Len() > 0 {
					return fmt.Errorf("`compliance_security_profile_standards` cannot be set when `compliance_security_profile_enabled` is false")
				}

				return nil
			}),

			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
				// Neither of these arguments can be removed once set
				for _, k := range []string{"managed_disk_cmk_key_vault_key_id", "managed_services_cmk_key_vault_key_id"} {
					o, n := d.GetChange(k)

					if o.(string) != "" && n.(string) == "" {
						// Check RawConfig to prevent replacements on `(known after apply)` values
						rawConfig := d.GetRawConfig()
						if rawConfig.IsNull() || !rawConfig.IsKnown() {
							return nil
						}
						rawValues := rawConfig.AsValueMap()

						if !rawValues[k].IsNull() {
							return nil
						}

						if err := d.ForceNew(k); err != nil {
							return err
						}
					}
				}

				return nil
			}),
		),
	}

	if !features.SixPointOh() {
		resource.Schema["managed_services_cmk_key_vault_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
			DiffSuppressFunc: func(_, o, n string, _ *pluginsdk.ResourceData) bool {
				// Suppress removal diff for 5.x since that does not require an update
				return o != "" && n == ""
			},
			Deprecated: "`managed_services_cmk_key_vault_id` has been deprecated and will be removed in v6.0 of the AzureRM provider. This property is no longer required for cross-subscription scenarios.",
		}

		resource.Schema["managed_disk_cmk_key_vault_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
			DiffSuppressFunc: func(_, o, n string, _ *pluginsdk.ResourceData) bool {
				// Suppress removal diff for 5.x since that does not require an update
				return o != "" && n == ""
			},
			Deprecated: "`managed_disk_cmk_key_vault_id` has been deprecated and will be removed in v6.0 of the AzureRM provider. This property is no longer required for cross-subscription scenarios.",
		}
	}

	return resource
}
