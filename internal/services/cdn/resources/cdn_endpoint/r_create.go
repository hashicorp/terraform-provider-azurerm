// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_endpoint

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnEndpointCreate(d *pluginsdk.ResourceData, meta any) error {
	endpointsClient := meta.(*clients.Client).Cdn.EndpointsClient
	profilesClient := meta.(*clients.Client).Cdn.ProfilesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := endpoints.NewEndpointID(subscriptionId, d.Get("resource_group_name").(string), d.Get("profile_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := endpointsClient.Get(ctx, id.ResourceGroupName, id.ProfileName, id.EndpointName)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.Response.Response) {
			return tf.ImportAsExistsError("azurerm_cdn_endpoint", id.ID())
		}
	}

	location := location.Normalize(d.Get("location").(string))
	httpAllowed := d.Get("is_http_allowed").(bool)
	httpsAllowed := d.Get("is_https_allowed").(bool)
	cachingBehaviour := d.Get("querystring_caching_behaviour").(string)
	originPath := d.Get("origin_path").(string)
	probePath := d.Get("probe_path").(string)
	optimizationType := d.Get("optimization_type").(string)
	t := d.Get("tags").(map[string]any)

	endpoint := cdn.Endpoint{
		Location: &location,
		EndpointProperties: &cdn.EndpointProperties{
			IsHTTPAllowed:              &httpAllowed,
			IsHTTPSAllowed:             &httpsAllowed,
			QueryStringCachingBehavior: cdn.QueryStringCachingBehavior(cachingBehaviour),
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("origin_host_header"); ok {
		endpoint.OriginHostHeader = pointer.To(v.(string))
	}

	if _, ok := d.GetOk("content_types_to_compress"); ok {
		contentTypes := expandArmCdnEndpointContentTypesToCompress(d)
		endpoint.ContentTypesToCompress = &contentTypes
	}

	if _, ok := d.GetOk("geo_filter"); ok {
		endpoint.GeoFilters = expandCdnEndpointGeoFilters(d)
	}

	if v, ok := d.GetOk("is_compression_enabled"); ok {
		endpoint.IsCompressionEnabled = pointer.To(v.(bool))
	}

	if optimizationType != "" {
		endpoint.OptimizationType = cdn.OptimizationType(optimizationType)
	}

	if originPath != "" {
		endpoint.OriginPath = pointer.To(originPath)
	}

	if probePath != "" {
		endpoint.ProbePath = pointer.To(probePath)
	}

	origins := expandAzureRmCdnEndpointOrigins(d)
	if len(origins) > 0 {
		endpoint.Origins = &origins
	}

	profile, err := profilesClient.Get(ctx, id.ResourceGroupName, id.ProfileName)
	if err != nil {
		return fmt.Errorf("retrieving parent CDN Profile for %s: %+v", id, err)
	}

	if profile.Sku != nil {
		globalDeliveryRulesRaw := d.Get("global_delivery_rule").([]any)
		deliveryRulesRaw := d.Get("delivery_rule").([]any)
		deliveryPolicy, err := expandArmCdnEndpointDeliveryPolicy(globalDeliveryRulesRaw, deliveryRulesRaw)
		if err != nil {
			return fmt.Errorf("expanding `global_delivery_rule` or `delivery_rule`: %s", err)
		}

		if profile.Sku.Name != cdn.SkuNameStandardMicrosoft && len(*deliveryPolicy.Rules) > 0 {
			return fmt.Errorf("`global_delivery_rule` and `delivery_rule` are only allowed when `Standard_Microsoft` sku is used. Profile sku:  %s", profile.Sku.Name)
		}

		if profile.Sku.Name == cdn.SkuNameStandardMicrosoft {
			endpoint.DeliveryPolicy = deliveryPolicy
		}
	}

	future, err := endpointsClient.Create(ctx, id.ResourceGroupName, id.ProfileName, id.EndpointName, endpoint)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, endpointsClient.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnEndpointRead(d, meta)
}
