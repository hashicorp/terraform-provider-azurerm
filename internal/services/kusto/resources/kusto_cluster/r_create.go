// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/clusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoClusterCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewKustoClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil && !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_kusto_cluster", id.ID())
		}
	}

	locks.ByName(id.KustoClusterName, "azurerm_kusto_cluster")
	defer locks.UnlockByName(id.KustoClusterName, "azurerm_kusto_cluster")

	sku, err := expandKustoClusterSku(d.Get("sku").([]any))
	if err != nil {
		return fmt.Errorf("expanding `sku`: %+v", err)
	}

	optimizedAutoScale := expandOptimizedAutoScale(d.Get("optimized_auto_scale").([]any))

	if optimizedAutoScale != nil && optimizedAutoScale.IsEnabled {
		if sku.Capacity == nil {
			return fmt.Errorf("sku.capacity could not be empty")
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

	publicNetworkAccess := clusters.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = clusters.PublicNetworkAccessDisabled
	}

	publicIPType := clusters.PublicIPType(d.Get("public_ip_type").(string))

	clusterProperties := clusters.ClusterProperties{
		OptimizedAutoscale:     optimizedAutoScale,
		EnableAutoStop:         pointer.To(d.Get("auto_stop_enabled").(bool)),
		EnableDiskEncryption:   pointer.To(d.Get("disk_encryption_enabled").(bool)),
		EnableDoubleEncryption: pointer.To(d.Get("double_encryption_enabled").(bool)),
		EnableStreamingIngest:  pointer.To(d.Get("streaming_ingestion_enabled").(bool)),
		EnablePurge:            pointer.To(d.Get("purge_enabled").(bool)),
		PublicNetworkAccess:    &publicNetworkAccess,
		PublicIPType:           &publicIPType,
		TrustedExternalTenants: expandTrustedExternalTenants(d.Get("trusted_external_tenants").([]any)),
	}

	if v, ok := d.GetOk("allowed_fqdns"); ok {
		clusterProperties.AllowedFqdnList = expandKustoListString(v.([]any))
	}

	if v, ok := d.GetOk("allowed_ip_ranges"); ok {
		clusterProperties.AllowedIPRangeList = expandKustoListString(v.([]any))
	}

	restrictOutboundNetworkAccess := clusters.ClusterNetworkAccessFlagDisabled
	if v, ok := d.GetOk("outbound_network_access_restricted"); ok {
		if v.(bool) {
			restrictOutboundNetworkAccess = clusters.ClusterNetworkAccessFlagEnabled
		}
	}
	clusterProperties.RestrictOutboundNetworkAccess = &restrictOutboundNetworkAccess

	if v, ok := d.GetOk("language_extension"); ok {
		clusterProperties.LanguageExtensions = expandKustoClusterLanguageExtensionList(v.([]any))
	}

	expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	kustoCluster := clusters.Cluster{
		Location:   location.Normalize(d.Get("location").(string)),
		Identity:   expandedIdentity,
		Sku:        *sku,
		Properties: &clusterProperties,
		Tags:       tags.Expand(d.Get("tags").(map[string]any)),
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		kustoCluster.Zones = &zones
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, kustoCluster, clusters.CreateOrUpdateOperationOptions{}, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceKustoClusterRead(d, meta)
}
