// Copyright IBM Corp. 2014, 2025
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
	model := existing.Model
	props := model.Properties

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
			// Ensure that requested Capacity is always between min and max to support updating to not overlapping autoscale ranges
			if *sku.Capacity < optimizedAutoScale.Minimum {
				sku.Capacity = pointer.To(optimizedAutoScale.Minimum)
			}
			if *sku.Capacity > optimizedAutoScale.Maximum {
				sku.Capacity = pointer.To(optimizedAutoScale.Maximum)
			}
			if optimizedAutoScale.Minimum > optimizedAutoScale.Maximum {
				return fmt.Errorf("`optimized_auto_scaling.maximum_instances` must be >= `optimized_auto_scaling.minimum_instances`")
			}
		}
		model.Sku = *sku
		// optimized_auto_scale can't be updated when the sku is also being updated, so we'll update sku first and then optimized_auto_scale
		if d.HasChange("optimized_auto_scale") && d.HasChange("sku") {
			model.Properties.OptimizedAutoscale = nil
			if err := client.CreateOrUpdateThenPoll(ctx, *id, *model, clusters.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
		}
		props.OptimizedAutoscale = optimizedAutoScale
	}

	if d.HasChange("location") {
		model.Location = location.Normalize(d.Get("location").(string))
	}

	if d.HasChange("identity") {
		expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		model.Identity = expandedIdentity
	}

	if d.HasChange("allowed_fqdns") {
		props.AllowedFqdnList = expandKustoListString(d.Get("allowed_fqdns").([]any))
	}

	if d.HasChange("allowed_ip_ranges") {
		props.AllowedIPRangeList = expandKustoListString(d.Get("allowed_ip_ranges").([]any))
	}

	if d.HasChange("auto_stop_enabled") {
		props.EnableAutoStop = pointer.To(d.Get("auto_stop_enabled").(bool))
	}

	if d.HasChange("disk_encryption_enabled") {
		props.EnableDiskEncryption = pointer.To(d.Get("disk_encryption_enabled").(bool))
	}

	if d.HasChange("double_encryption_enabled") {
		props.EnableDoubleEncryption = pointer.To(d.Get("double_encryption_enabled").(bool))
	}

	if d.HasChange("language_extension") {
		props.LanguageExtensions = expandKustoClusterLanguageExtensionList(d.Get("language_extension").([]any))
	}

	if d.HasChange("outbound_network_access_restricted") {
		restrictOutboundNetworkAccess := clusters.ClusterNetworkAccessFlagDisabled
		if d.Get("outbound_network_access_restricted").(bool) {
			restrictOutboundNetworkAccess = clusters.ClusterNetworkAccessFlagEnabled
		}
		props.RestrictOutboundNetworkAccess = pointer.To(restrictOutboundNetworkAccess)
	}

	if d.HasChange("purge_enabled") {
		props.EnablePurge = pointer.To(d.Get("purge_enabled").(bool))
	}

	if d.HasChange("public_ip_type") {
		publicIPType := clusters.PublicIPType(d.Get("public_ip_type").(string))
		props.PublicIPType = pointer.To(publicIPType)
	}

	if d.HasChange("public_network_access_enabled") {
		publicNetworkAccess := clusters.PublicNetworkAccessEnabled
		if !d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = clusters.PublicNetworkAccessDisabled
		}
		props.PublicNetworkAccess = pointer.To(publicNetworkAccess)
	}

	if d.HasChange("streaming_ingestion_enabled") {
		props.EnableStreamingIngest = pointer.To(d.Get("streaming_ingestion_enabled").(bool))
	}

	if d.HasChange("trusted_external_tenants") {
		props.TrustedExternalTenants = expandTrustedExternalTenants(d.Get("trusted_external_tenants").([]any))
	}

	if d.HasChange("zones") {
		zoneList := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
		model.Zones = pointer.To(zoneList)
		if len(zoneList) > 0 {
			model.Zones = &zoneList
		} else {
			model.Zones = nil
		}
	}

	if d.HasChange("tags") {
		model.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	model.Properties = props

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *model, clusters.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceKustoClusterRead(d, meta)
}
