// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster

import (
	"context"
	"fmt"
	"maps"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/kubernetes_cluster/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceKubernetesClusterCreate,
		Read:   resourceKubernetesClusterRead,
		Update: resourceKubernetesClusterUpdate,
		Delete: resourceKubernetesClusterDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := commonids.ParseKubernetesClusterID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta any) ([]*pluginsdk.ResourceData, error) {
			resourceId, err := commonids.ParseKubernetesClusterID(d.Id())
			if err != nil {
				return []*pluginsdk.ResourceData{d}, err
			}

			client := meta.(*clients.Client).Containers.KubernetesClustersClient
			resp, err := client.Get(ctx, *resourceId)
			if err != nil || resp.Model == nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("retrieving %s: %+v", *resourceId, err)
			}

			if resp.Model.Sku == nil || resp.Model.Sku.Name == nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("importing %s: SKU information is missing", resourceId)
			}

			if pointer.From(resp.Model.Sku.Name) == managedclusters.ManagedClusterSKUNameAutomatic {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("importing %s: use the `azurerm_kubernetes_automatic_cluster` resource for clusters using the `Automatic` SKU", d.Id())
			}

			return []*pluginsdk.ResourceData{d}, nil
		}),

		CustomizeDiff: pluginsdk.CustomDiffInSequence(
			// The behaviour of the API requires this, but this could be removed when https://github.com/Azure/azure-rest-api-specs/issues/27373 has been addressed
			pluginsdk.ForceNewIfChange("default_node_pool.0.upgrade_settings.0.drain_timeout_in_minutes", func(ctx context.Context, old, new, meta any) bool {
				return old != 0 && new == 0
			}),
			pluginsdk.ForceNewIfChange("default_node_pool.0.upgrade_settings.0.undrainable_node_behavior", func(ctx context.Context, old, new, meta any) bool {
				return old != "" && new == ""
			}),
			pluginsdk.ForceNewIfChange("oidc_issuer_enabled", func(ctx context.Context, old, new, meta any) bool {
				// The API does not allow disabling OIDC after it has been enabled
				return old.(bool) && !new.(bool)
			}),
			// Migration of `identity` to `service_principal` is not allowed, the other way around is
			pluginsdk.ForceNewIfChange("service_principal.0.client_id", func(ctx context.Context, old, new, meta any) bool {
				return old == "msi" || old == ""
			}),
			pluginsdk.ForceNewIfChange("windows_profile.0.gmsa", func(ctx context.Context, old, new, meta any) bool {
				return len(old.([]any)) != 0 && len(new.([]any)) == 0
			}),
			pluginsdk.ForceNewIfChange("windows_profile.0.gmsa.0.dns_server", func(ctx context.Context, old, new, meta any) bool {
				return old != "" && new == ""
			}),
			pluginsdk.ForceNewIfChange("windows_profile.0.gmsa.0.root_domain", func(ctx context.Context, old, new, meta any) bool {
				return old != "" && new == ""
			}),
			pluginsdk.ForceNewIfChange("api_server_access_profile.0.subnet_id", func(ctx context.Context, old, new, meta any) bool {
				return old != "" && new == ""
			}),
			pluginsdk.ForceNewIf("default_node_pool.0.name", func(ctx context.Context, d *schema.ResourceDiff, meta any) bool {
				old, new := d.GetChange("default_node_pool.0.name")
				defaultName := d.Get("default_node_pool.0.name")
				tempName := d.Get("default_node_pool.0.temporary_name_for_rotation")

				// if the default node pool name has been set to temporary_name_for_rotation it means resizing failed
				// we should not try to recreate the cluster, another apply will attempt the resize again
				if old != "" && old == tempName {
					return new != defaultName
				}
				return true
			}),
			pluginsdk.ForceNewIfChange("network_profile.0.network_data_plane", func(ctx context.Context, old, new, meta any) bool {
				oldStr := old.(string)
				newStr := new.(string)

				if oldStr == "" {
					return false
				}

				if oldStr == "azure" && newStr == "cilium" {
					return false
				}

				return true
			}),
			func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
				if d.HasChange("network_profile.0.network_plugin") {
					old, new := d.GetChange("network_profile.0.network_plugin")
					oldStr := old.(string)
					newStr := new.(string)

					if oldStr == "kubenet" && newStr == "azure" {
						networkPluginMode := d.Get("network_profile.0.network_plugin_mode").(string)
						if strings.EqualFold(networkPluginMode, string(managedclusters.NetworkPluginModeOverlay)) {
							return nil
						}
					}

					return d.ForceNew("network_profile.0.network_plugin")
				}
				return nil
			},
			func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
				if d.HasChange("oidc_issuer_enabled") {
					d.SetNewComputed("oidc_issuer_url")
				}
				return nil
			},
			pluginsdk.ForceNewIfChange("network_profile.0.network_plugin_mode", func(ctx context.Context, _, new, meta any) bool {
				return !strings.EqualFold(new.(string), string(managedclusters.NetworkPluginModeOverlay))
			}),
			pluginsdk.ForceNewIfChange("network_profile.0.network_policy", func(ctx context.Context, old, new, meta any) bool {
				// Azure supports in-place upgrade from Azure Network Policy Manager to Cilium and Calico to Cilium
				// Other transitions (removing policy, etc.) require cluster recreation
				oldStr := old.(string)
				newStr := new.(string)

				if oldStr == "" {
					return false
				}

				if oldStr == "azure" && newStr == "cilium" {
					return false
				}
				if oldStr == "calico" && newStr == "cilium" {
					return false
				}

				return true
			}),
			pluginsdk.ForceNewIfChange("network_profile.0.pod_cidr", func(ctx context.Context, old, new, meta any) bool {
				return old.(string) != ""
			}),
			pluginsdk.ForceNewIfChange("network_profile.0.pod_cidrs", func(ctx context.Context, old, new, meta any) bool {
				return len(old.([]any)) > 0
			}),
			func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
				outboundType := d.Get("network_profile.0.outbound_type").(string)
				artifactSource := d.Get("bootstrap_profile.0.artifact_source").(string)

				if outboundType == string(managedclusters.OutboundTypeNone) && artifactSource != string(managedclusters.ArtifactSourceCache) {
					return fmt.Errorf("when `network_profile.0.outbound_type` is set to `none`, `bootstrap_profile.0.artifact_source` must be set to `Cache`")
				}

				return nil
			},
			func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
				if len(d.Get("network_profile.0.advanced_networking").([]any)) == 1 {
					securityEnabled := d.Get("network_profile.0.advanced_networking.0.security_enabled").(bool)

					if securityEnabled {
						if d.Get("network_profile.0.network_data_plane").(string) != string(managedclusters.NetworkDataplaneCilium) {
							return fmt.Errorf("when `network_profile.0.advanced_networking` has `security_enabled` set to `true`, `network_profile.0.network_data_plane` must be set to `%s`", managedclusters.NetworkDataplaneCilium)
						}
						if d.Get("network_profile.0.network_plugin").(string) != string(managedclusters.NetworkPluginAzure) {
							return fmt.Errorf("when `network_profile.0.advanced_networking` has `security_enabled` set to `true`, `network_profile.0.network_plugin` must be set to `%s`", managedclusters.NetworkPluginAzure)
						}
					}
				}
				return nil
			},
		),
		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KubernetesClusterV0ToV1{},
			1: migration.KubernetesClusterV1ToV2{},
		}),

		Schema: kubernetesClusterSchema(),
	}

	maps.Copy(resource.Schema, schemaKubernetesAddOns())

	return resource
}
