// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package trafficmanager

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/endpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	azSchema "github.com/hashicorp/terraform-provider-azurerm/internal/tf/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceExternalEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceExternalEndpointCreate,
		Read:   resourceExternalEndpointRead,
		Update: resourceExternalEndpointUpdate,
		Delete: resourceExternalEndpointDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			endpointType, err := endpoints.ParseEndpointTypeID(id)
			if err != nil {
				return err
			}

			if endpointType.EndpointType != endpoints.EndpointTypeExternalEndpoints {
				return fmt.Errorf("this resource only supports `ExternalEndpoints` but got %s", string(endpointType.EndpointType))
			}

			return nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: profiles.ValidateTrafficManagerProfileID,
			},

			"target": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"weight": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 1000),
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"always_serve_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"custom_header": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"value": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
					},
				},
			},

			"priority": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				// NOTE: O+C the API dynamically increments the default value for priority depending on the number of endpoints
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 1000),
			},

			"endpoint_location": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Computed:         true,
				StateFunc:        location.StateFunc,
				DiffSuppressFunc: location.DiffSuppressFunc,
			},

			"geo_mappings": {
				Type:     pluginsdk.TypeList,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Optional: true,
			},

			"subnet": {
				Type:     pluginsdk.TypeList,
				ForceNew: true,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"first": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azValidate.IPv4Address,
						},
						"last": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: azValidate.IPv4Address,
						},
						"scope": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 32),
						},
					},
				},
			},
		},
	}
}

func resourceExternalEndpointCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.EndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := profiles.ParseTrafficManagerProfileID(d.Get("profile_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `profile_id`: %+v", err)
	}

	id := endpoints.NewEndpointTypeID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.TrafficManagerProfileName, endpoints.EndpointTypeExternalEndpoints, d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_traffic_manager_external_endpoint", id.ID())
	}

	status := endpoints.EndpointStatusEnabled
	if !d.Get("enabled").(bool) {
		status = endpoints.EndpointStatusDisabled
	}

	params := endpoints.Endpoint{
		Name: utils.String(id.EndpointName),
		Type: utils.String(fmt.Sprintf("Microsoft.Network/trafficManagerProfiles/%s", endpoints.EndpointTypeExternalEndpoints)),
		Properties: &endpoints.EndpointProperties{
			AlwaysServe:    pointer.To(endpoints.AlwaysServeDisabled),
			CustomHeaders:  expandEndpointCustomHeaderConfig(d.Get("custom_header").([]interface{})),
			EndpointStatus: &status,
			Target:         utils.String(d.Get("target").(string)),
			Subnets:        expandEndpointSubnetConfig(d.Get("subnet").([]interface{})),
		},
	}

	if alwaysServe := d.Get("always_serve_enabled").(bool); alwaysServe {
		params.Properties.AlwaysServe = pointer.To(endpoints.AlwaysServeEnabled)
	}

	if priority := d.Get("priority").(int); priority != 0 {
		params.Properties.Priority = utils.Int64(int64(priority))
	}

	if weight := d.Get("weight").(int); weight != 0 {
		params.Properties.Weight = utils.Int64(int64(weight))
	}

	if endpointLocation := d.Get("endpoint_location").(string); endpointLocation != "" {
		params.Properties.EndpointLocation = utils.String(endpointLocation)
	}

	inputMappings := d.Get("geo_mappings").([]interface{})
	geoMappings := make([]string, 0)
	for _, v := range inputMappings {
		geoMappings = append(geoMappings, v.(string))
	}
	if len(geoMappings) > 0 {
		params.Properties.GeoMapping = &geoMappings
	}

	if _, err := client.CreateOrUpdate(ctx, id, params); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceExternalEndpointRead(d, meta)
}

func resourceExternalEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.EndpointsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := endpoints.ParseEndpointTypeID(d.Id())
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

	d.Set("name", id.EndpointName)
	d.Set("profile_id", profiles.NewTrafficManagerProfileID(id.SubscriptionId, id.ResourceGroupName, id.TrafficManagerProfileName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			enabled := true
			if props.EndpointStatus != nil && *props.EndpointStatus == endpoints.EndpointStatusDisabled {
				enabled = false
			}
			d.Set("enabled", enabled)
			d.Set("target", props.Target)
			d.Set("weight", props.Weight)
			d.Set("priority", props.Priority)
			d.Set("endpoint_location", props.EndpointLocation)
			d.Set("geo_mappings", props.GeoMapping)

			if props.AlwaysServe != nil && *props.AlwaysServe == endpoints.AlwaysServeEnabled {
				d.Set("always_serve_enabled", true)
			} else {
				d.Set("always_serve_enabled", false)
			}

			if err := d.Set("custom_header", flattenEndpointCustomHeaderConfig(props.CustomHeaders)); err != nil {
				return fmt.Errorf("setting `custom_header`: %s", err)
			}
			if err := d.Set("subnet", flattenEndpointSubnetConfig(props.Subnets)); err != nil {
				return fmt.Errorf("setting `subnet`: %s", err)
			}
		}
	}

	return nil
}

func resourceExternalEndpointUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.EndpointsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := profiles.ParseTrafficManagerProfileID(d.Get("profile_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `profile_id`: %+v", err)
	}

	id := endpoints.NewEndpointTypeID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.TrafficManagerProfileName, endpoints.EndpointTypeExternalEndpoints, d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("checking for presence of existing %s: %v", id, err)
	}

	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("model/properties was nil for %s", id)
	}

	params := *existing.Model

	if d.HasChange("enabled") {
		status := endpoints.EndpointStatusEnabled
		if !d.Get("enabled").(bool) {
			status = endpoints.EndpointStatusDisabled
		}
		params.Properties.EndpointStatus = pointer.To(status)
	}

	if d.HasChange("always_serve_enabled") {
		alwaysServe := endpoints.AlwaysServeDisabled
		if d.Get("always_serve_enabled").(bool) {
			alwaysServe = endpoints.AlwaysServeEnabled
		}
		params.Properties.AlwaysServe = pointer.To(alwaysServe)
	}

	if d.HasChange("custom_header") {
		params.Properties.CustomHeaders = expandEndpointCustomHeaderConfig(d.Get("custom_header").([]interface{}))
	}

	if d.HasChange("target") {
		params.Properties.Target = utils.String(d.Get("target").(string))
	}

	if d.HasChange("subnet") {
		params.Properties.Subnets = expandEndpointSubnetConfig(d.Get("subnet").([]interface{}))
	}

	if d.HasChange("priority") {
		if priority := d.Get("priority").(int); priority != 0 {
			params.Properties.Priority = utils.Int64(int64(priority))
		}
	}

	if d.HasChange("weight") {
		if weight := d.Get("weight").(int); weight != 0 {
			params.Properties.Weight = utils.Int64(int64(weight))
		}
	}

	if d.HasChange("endpoint_location") {
		if endpointLocation := d.Get("endpoint_location").(string); endpointLocation != "" {
			params.Properties.EndpointLocation = utils.String(endpointLocation)
		} else {
			params.Properties.EndpointLocation = nil
		}
	}

	if d.HasChange("geo_mappings") {
		inputMappings := d.Get("geo_mappings").([]interface{})
		geoMappings := make([]string, 0)
		for _, v := range inputMappings {
			geoMappings = append(geoMappings, v.(string))
		}
		if len(geoMappings) > 0 {
			params.Properties.GeoMapping = &geoMappings
		} else {
			params.Properties.GeoMapping = nil
		}
	}
	if _, err := client.CreateOrUpdate(ctx, id, params); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	return resourceExternalEndpointRead(d, meta)
}

func resourceExternalEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.EndpointsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := endpoints.ParseEndpointTypeID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
