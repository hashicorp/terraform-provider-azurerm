// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/clusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoClusterUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseKustoClusterID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.KustoClusterName, "azurerm_kusto_cluster")
	defer locks.UnlockByName(id.KustoClusterName, "azurerm_kusto_cluster")

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	patch := clusters.ClusterUpdate{}
	props := clusters.ClusterProperties{}
	hasPropertyChanges := false

	if d.HasChanges("sku", "optimized_auto_scale") {
		sku, err := expandKustoClusterSku(d.Get("sku").([]any))
		if err != nil {
			return err
		}

		optimizedAutoScale := expandOptimizedAutoScale(d.Get("optimized_auto_scale").([]any))
		if optimizedAutoScale != nil && optimizedAutoScale.IsEnabled {
			if sku.Capacity == nil {
				return fmt.Errorf("sku.capacity cannot be empty")
			}
			if optimizedAutoScale.Minimum > optimizedAutoScale.Maximum {
				return fmt.Errorf("`optimized_auto_scaling.maximum_instances` must be >= `optimized_auto_scaling.minimum_instances`")
			}
			originalCapacity := *sku.Capacity
			if *sku.Capacity < optimizedAutoScale.Minimum {
				sku.Capacity = pointer.To(optimizedAutoScale.Minimum)
			}
			if *sku.Capacity > optimizedAutoScale.Maximum {
				sku.Capacity = pointer.To(optimizedAutoScale.Maximum)
			}
			if *sku.Capacity != originalCapacity {
				patch.Sku = sku
			}
		}

		if d.HasChange("optimized_auto_scale") {
			if optimizedAutoScale == nil {
				optimizedAutoScale = &clusters.OptimizedAutoscale{Version: 1}
			}
			props.OptimizedAutoscale = optimizedAutoScale
			hasPropertyChanges = true
		}
		if d.HasChange("sku") {
			patch.Sku = sku
		}
	}

	if d.HasChange("location") {
		patch.Location = pointer.To(location.Normalize(d.Get("location").(string)))
	}

	if d.HasChange("identity") {
		expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		patch.Identity = expandedIdentity
	}

	if d.HasChange("allowed_fqdns") {
		props.AllowedFqdnList = expandKustoListString(d.Get("allowed_fqdns").([]any))
		hasPropertyChanges = true
	}

	if d.HasChange("allowed_ip_ranges") {
		props.AllowedIPRangeList = expandKustoListString(d.Get("allowed_ip_ranges").([]any))
		hasPropertyChanges = true
	}

	if d.HasChange("auto_stop_enabled") {
		props.EnableAutoStop = pointer.To(d.Get("auto_stop_enabled").(bool))
		hasPropertyChanges = true
	}

	if d.HasChange("customer_managed_key") {
		customerManagedKey, err := expandKustoClusterCustomerManagedKey(d.Get("customer_managed_key").([]any))
		if err != nil {
			return fmt.Errorf("expanding `customer_managed_key`: %+v", err)
		}
		if customerManagedKey == nil {
			customerManagedKey = &clusters.KeyVaultProperties{}
		}
		props.KeyVaultProperties = customerManagedKey
		hasPropertyChanges = true
	}

	if d.HasChange("disk_encryption_enabled") {
		props.EnableDiskEncryption = pointer.To(d.Get("disk_encryption_enabled").(bool))
		hasPropertyChanges = true
	}

	if d.HasChange("double_encryption_enabled") {
		props.EnableDoubleEncryption = pointer.To(d.Get("double_encryption_enabled").(bool))
		hasPropertyChanges = true
	}

	if d.HasChange("language_extension") {
		props.LanguageExtensions = expandKustoClusterLanguageExtensionList(d.Get("language_extension").([]any))
		hasPropertyChanges = true
	}

	if d.HasChange("outbound_network_access_restricted") {
		restrictOutboundNetworkAccess := clusters.ClusterNetworkAccessFlagDisabled
		if d.Get("outbound_network_access_restricted").(bool) {
			restrictOutboundNetworkAccess = clusters.ClusterNetworkAccessFlagEnabled
		}
		props.RestrictOutboundNetworkAccess = pointer.To(restrictOutboundNetworkAccess)
		hasPropertyChanges = true
	}

	if d.HasChange("purge_enabled") {
		props.EnablePurge = pointer.To(d.Get("purge_enabled").(bool))
		hasPropertyChanges = true
	}

	if d.HasChange("public_ip_type") {
		props.PublicIPType = pointer.To(clusters.PublicIPType(d.Get("public_ip_type").(string)))
		hasPropertyChanges = true
	}

	if d.HasChange("public_network_access_enabled") {
		publicNetworkAccess := clusters.PublicNetworkAccessEnabled
		if !d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = clusters.PublicNetworkAccessDisabled
		}
		props.PublicNetworkAccess = pointer.To(publicNetworkAccess)
		hasPropertyChanges = true
	}

	if d.HasChange("streaming_ingestion_enabled") {
		props.EnableStreamingIngest = pointer.To(d.Get("streaming_ingestion_enabled").(bool))
		hasPropertyChanges = true
	}

	if d.HasChange("trusted_external_tenants") {
		props.TrustedExternalTenants = expandTrustedExternalTenants(d.Get("trusted_external_tenants").([]any))
		hasPropertyChanges = true
	}

	if d.HasChange("zones") {
		zoneList := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
		if len(zoneList) > 0 {
			patch.Zones = pointer.To(zoneList)
		} else {
			patch.Zones = nil
		}
	}

	if d.HasChange("tags") {
		patch.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if hasPropertyChanges {
		patch.Properties = &props
	}

	// SKU and autoscale must be changed in separate operations by the Kusto API.
	autoScaleIsEnabled := existing.Model.Properties.OptimizedAutoscale != nil && existing.Model.Properties.OptimizedAutoscale.IsEnabled
	if patch.Sku != nil && (d.HasChange("optimized_auto_scale") || autoScaleIsEnabled) {
		if autoScaleIsEnabled {
			disableAutoScalePatch := clusters.ClusterUpdate{
				Properties: &clusters.ClusterProperties{
					OptimizedAutoscale: &clusters.OptimizedAutoscale{Version: 1},
				},
			}
			if err := client.UpdateThenPoll(ctx, *id, disableAutoScalePatch, clusters.UpdateOperationOptions{}); err != nil {
				return fmt.Errorf("disabling optimized autoscale on %s before updating SKU: %+v", id, err)
			}
		}

		skuPatch := clusters.ClusterUpdate{Sku: patch.Sku}
		if err := client.UpdateThenPoll(ctx, *id, skuPatch, clusters.UpdateOperationOptions{}); err != nil {
			return fmt.Errorf("updating %s SKU: %+v", id, err)
		}
		patch.Sku = nil

		if autoScaleIsEnabled {
			if patch.Properties == nil {
				patch.Properties = &clusters.ClusterProperties{}
			}
			if patch.Properties.OptimizedAutoscale == nil {
				patch.Properties.OptimizedAutoscale = existing.Model.Properties.OptimizedAutoscale
			}
		}
	}

	if err := client.UpdateThenPoll(ctx, *id, patch, clusters.UpdateOperationOptions{}); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceKustoClusterRead(d, meta)
}
