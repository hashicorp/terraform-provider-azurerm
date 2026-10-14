// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_endpoint

import (
	"fmt"

	//nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnEndpointRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.EndpointsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := endpoints.ParseEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.ProfileName, id.EndpointName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.EndpointName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("profile_name", id.ProfileName)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.EndpointProperties; props != nil {
		d.Set("fqdn", props.HostName)
		d.Set("is_http_allowed", props.IsHTTPAllowed)
		d.Set("is_https_allowed", props.IsHTTPSAllowed)
		d.Set("querystring_caching_behaviour", props.QueryStringCachingBehavior)
		d.Set("origin_host_header", props.OriginHostHeader)
		d.Set("origin_path", props.OriginPath)
		d.Set("probe_path", props.ProbePath)
		d.Set("optimization_type", string(props.OptimizationType))

		d.Set("is_compression_enabled", pointer.From(props.IsCompressionEnabled))

		if err := d.Set("content_types_to_compress", flattenAzureRMCdnEndpointContentTypes(props.ContentTypesToCompress)); err != nil {
			return fmt.Errorf("setting `content_types_to_compress`: %+v", err)
		}

		if err := d.Set("geo_filter", flattenCdnEndpointGeoFilters(props.GeoFilters)); err != nil {
			return fmt.Errorf("setting `geo_filter`: %+v", err)
		}

		if err := d.Set("origin", flattenAzureRMCdnEndpointOrigin(props.Origins)); err != nil {
			return fmt.Errorf("setting `origin`: %+v", err)
		}

		flattenedDeliveryPolicies, err := flattenEndpointDeliveryPolicy(props.DeliveryPolicy)
		if err != nil {
			return err
		}
		if err := d.Set("global_delivery_rule", flattenedDeliveryPolicies.globalDeliveryRules); err != nil {
			return fmt.Errorf("setting `global_delivery_rule`: %+v", err)
		}
		if err := d.Set("delivery_rule", flattenedDeliveryPolicies.deliveryRules); err != nil {
			return fmt.Errorf("setting `delivery_rule`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
