// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_node_pool

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/agentpools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/kubernetes_cluster_node_pool/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const azureKubernetesClusterNodePoolResourceName = "azurerm_kubernetes_cluster_node_pool"

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "name" -compare-values "managed_cluster_name:kubernetes_cluster_id,resource_group_name:kubernetes_cluster_id,subscription_id:kubernetes_cluster_id" -test-name "manualScaleConfig"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKubernetesClusterNodePoolCreate,
		Read:   resourceKubernetesClusterNodePoolRead,
		Update: resourceKubernetesClusterNodePoolUpdate,
		Delete: resourceKubernetesClusterNodePoolDelete,

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&agentpools.AgentPoolId{}),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&agentpools.AgentPoolId{}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KubernetesClusterNodePoolV0ToV1{},
		}),

		Schema: resourceKubernetesClusterNodePoolSchema(),

		CustomizeDiff: pluginsdk.CustomDiffInSequence(
			pluginsdk.ForceNewIfChange("os_sku", func(ctx context.Context, old, new, meta any) bool {
				oldStr := old.(string)
				newStr := new.(string)

				// Ubuntu and AzureLinux are currently the only allowed Linux OSSKU Migration targets.
				if !strings.HasPrefix(oldStr, string(agentpools.OSSKUUbuntu)) && !strings.HasPrefix(oldStr, string(agentpools.OSSKUAzureLinux)) {
					return true
				}

				if !strings.HasPrefix(newStr, string(agentpools.OSSKUUbuntu)) && !strings.HasPrefix(newStr, string(agentpools.OSSKUAzureLinux)) {
					return true
				}

				return false
			}),
			// The behaviour of the API requires this, but this could be removed when https://github.com/Azure/azure-rest-api-specs/issues/27373 has been addressed
			pluginsdk.ForceNewIfChange("upgrade_settings.0.drain_timeout_in_minutes", func(ctx context.Context, old, new, meta any) bool {
				return old != 0 && new == 0
			}),
			pluginsdk.ForceNewIfChange("upgrade_settings.0.undrainable_node_behavior", func(ctx context.Context, old, new, meta any) bool {
				return old != "" && new == ""
			}),
			func(ctx context.Context, d *pluginsdk.ResourceDiff, meta any) error {
				priority := d.Get("priority").(string)
				isSpot := priority == string(agentpools.ScaleSetPrioritySpot)

				upgradeSettingsRaw := d.Get("upgrade_settings").([]any)
				if len(upgradeSettingsRaw) == 0 || upgradeSettingsRaw[0] == nil {
					return nil
				}

				upgradeSettings := upgradeSettingsRaw[0].(map[string]any)
				maxSurgeRaw := upgradeSettings["max_surge"].(string)
				maxUnavailableRaw := upgradeSettings["max_unavailable"].(string)

				if isSpot {
					if maxSurgeRaw != "" {
						return fmt.Errorf("`max_surge` cannot be set when `priority` is set to `Spot`. Spot pools do not support `max_surge`")
					}
					if maxUnavailableRaw != "" {
						return fmt.Errorf("`max_unavailable` cannot be set when `priority` is set to `Spot`. Spot pools do not support `max_unavailable`")
					}
				} else if maxSurgeRaw == "" && maxUnavailableRaw == "" {
					return fmt.Errorf("either `max_surge` or `max_unavailable` must be specified in `upgrade_settings` when `priority` is not `Spot`")
				}

				return nil
			},
		),
	}
}
