// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package traffic_manager_external_endpoint

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/trafficmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	trafficmanagerClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExternalEndpointCreate(d *pluginsdk.ResourceData, meta any) error {
	client := trafficmanagerClient.From(meta).EndpointsClient
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	profileId, err := profiles.ParseTrafficManagerProfileID(d.Get("profile_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `profile_id`: %+v", err)
	}

	id := trafficmanagers.NewEndpointTypeID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.TrafficManagerProfileName, trafficmanagers.EndpointTypeExternalEndpoints, d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.EndpointsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_traffic_manager_external_endpoint", id.ID())
		}
	}

	status := trafficmanagers.EndpointStatusEnabled
	if !d.Get("enabled").(bool) {
		status = trafficmanagers.EndpointStatusDisabled
	}

	params := trafficmanagers.Endpoint{
		Name: pointer.To(id.EndpointName),
		Type: pointer.To(fmt.Sprintf("Microsoft.Network/trafficManagerProfiles/%s", trafficmanagers.EndpointTypeExternalEndpoints)),
		Properties: &trafficmanagers.EndpointProperties{
			AlwaysServe:    pointer.To(trafficmanagers.AlwaysServeDisabled),
			CustomHeaders:  helpers.ExpandEndpointCustomHeaderConfig(d.Get("custom_header").([]any)),
			EndpointStatus: &status,
			Target:         pointer.To(d.Get("target").(string)),
			Subnets:        helpers.ExpandEndpointSubnetConfig(d.Get("subnet").([]any)),
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

	if endpointLocation := d.Get("endpoint_location").(string); endpointLocation != "" {
		params.Properties.EndpointLocation = pointer.To(endpointLocation)
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
	return resourceExternalEndpointRead(d, meta)
}
