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
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoClusterRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseKustoClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.KustoClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(&model.Location))
		d.Set("zones", zones.FlattenUntyped(model.Zones))

		identity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %s", err)
		}

		if err := d.Set("sku", flattenKustoClusterSku(&model.Sku)); err != nil {
			return fmt.Errorf("setting `sku`: %+v", err)
		}

		if props := model.Properties; props != nil {
			if props.PublicNetworkAccess != nil {
				d.Set("public_network_access_enabled", *props.PublicNetworkAccess == clusters.PublicNetworkAccessEnabled)
			}

			if props.RestrictOutboundNetworkAccess != nil {
				d.Set("outbound_network_access_restricted", *props.RestrictOutboundNetworkAccess == clusters.ClusterNetworkAccessFlagEnabled)
			}

			if err := d.Set("optimized_auto_scale", flattenOptimizedAutoScale(props.OptimizedAutoscale)); err != nil {
				return fmt.Errorf("setting `optimized_auto_scale`: %+v", err)
			}
			d.Set("allowed_fqdns", props.AllowedFqdnList)
			d.Set("allowed_ip_ranges", props.AllowedIPRangeList)
			d.Set("double_encryption_enabled", props.EnableDoubleEncryption)
			d.Set("trusted_external_tenants", flattenTrustedExternalTenants(props.TrustedExternalTenants))
			d.Set("auto_stop_enabled", props.EnableAutoStop)
			d.Set("disk_encryption_enabled", props.EnableDiskEncryption)
			d.Set("streaming_ingestion_enabled", props.EnableStreamingIngest)
			d.Set("purge_enabled", props.EnablePurge)
			d.Set("uri", props.Uri)
			d.Set("data_ingestion_uri", props.DataIngestionUri)
			d.Set("public_ip_type", pointer.FromEnum(props.PublicIPType))

			d.Set("language_extension", flattenKustoClusterLanguageExtensionList(props.LanguageExtensions))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
