// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package traffic_manager_azure_endpoint

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/trafficmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAzureEndpointCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).TrafficManager.EndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := profiles.ParseTrafficManagerProfileID(d.Get("profile_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `profile_id`: %+v", err)
	}

	id := trafficmanagers.NewEndpointTypeID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.TrafficManagerProfileName, trafficmanagers.EndpointTypeAzureEndpoints, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.EndpointsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_traffic_manager_azure_endpoint", id.ID())
		}
	}

	status := trafficmanagers.EndpointStatusEnabled
	if !d.Get("enabled").(bool) {
		status = trafficmanagers.EndpointStatusDisabled
	}

	params := trafficmanagers.Endpoint{
		Name: pointer.To(id.EndpointName),
		Type: pointer.To(fmt.Sprintf("Microsoft.Network/trafficManagerProfiles/%s", trafficmanagers.EndpointTypeAzureEndpoints)),
		Properties: &trafficmanagers.EndpointProperties{
			CustomHeaders:    helpers.ExpandEndpointCustomHeaderConfig(d.Get("custom_header").([]any)),
			AlwaysServe:      pointer.To(trafficmanagers.AlwaysServeDisabled),
			EndpointStatus:   &status,
			TargetResourceId: pointer.To(d.Get("target_resource_id").(string)),
			Subnets:          helpers.ExpandEndpointSubnetConfig(d.Get("subnet").([]any)),
		},
	}

	if alwaysServe := d.Get("always_serve_enabled").(bool); alwaysServe {
		params.Properties.AlwaysServe = pointer.To(trafficmanagers.AlwaysServeEnabled)
	}

	if priority := d.Get("priority").(int); priority != 0 {
		params.Properties.Priority = pointer.To(int64(priority))
	}

	if weight := d.Get("weight").(int); weight != 0 {
		params.Properties.Weight = pointer.To(int64(weight))
	}

	inputMappings := d.Get("geo_mappings").([]any)
	geoMappings := make([]string, 0)
	for _, v := range inputMappings {
		geoMappings = append(geoMappings, v.(string))
	}
	if len(geoMappings) > 0 {
		params.Properties.GeoMapping = &geoMappings
	}

	if _, err := client.EndpointsCreateOrUpdate(ctx, id, params); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceAzureEndpointRead(d, meta)
}
