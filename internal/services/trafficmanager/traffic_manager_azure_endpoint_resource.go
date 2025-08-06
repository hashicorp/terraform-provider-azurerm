// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package trafficmanager

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/endpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	azSchema "github.com/hashicorp/terraform-provider-azurerm/internal/tf/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAzureEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAzureEndpointCreate,
		Read:   resourceAzureEndpointRead,
		Update: resourceAzureEndpointUpdate,
		Delete: resourceAzureEndpointDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			endpointType, err := endpoints.ParseEndpointTypeID(id)
			if err != nil {
				return err
			}

			if endpointType.EndpointType != endpoints.EndpointTypeAzureEndpoints {
				return fmt.Errorf("this resource only supports `AzureEndpoints` but got %s", string(endpointType.EndpointType))
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

			"target_resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
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

func resourceAzureEndpointCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.EndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := profiles.ParseTrafficManagerProfileID(d.Get("profile_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `profile_id`: %+v", err)
	}

	id := endpoints.NewEndpointTypeID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.TrafficManagerProfileName, endpoints.EndpointTypeAzureEndpoints, d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_traffic_manager_azure_endpoint", id.ID())
	}

	status := endpoints.EndpointStatusEnabled
	if !d.Get("enabled").(bool) {
		status = endpoints.EndpointStatusDisabled
	}

	params := endpoints.Endpoint{
		Name: utils.String(id.EndpointName),
		Type: utils.String(fmt.Sprintf("Microsoft.Network/trafficManagerProfiles/%s", endpoints.EndpointTypeAzureEndpoints)),
		Properties: &endpoints.EndpointProperties{
			CustomHeaders:    expandEndpointCustomHeaderConfig(d.Get("custom_header").([]interface{})),
			AlwaysServe:      pointer.To(endpoints.AlwaysServeDisabled),
			EndpointStatus:   &status,
			TargetResourceId: utils.String(d.Get("target_resource_id").(string)),
			Subnets:          expandEndpointSubnetConfig(d.Get("subnet").([]interface{})),
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

	inputMappings := d.Get("geo_mappings").([]interface{})
	geoMappings := make([]string, 0)
	for _, v := range inputMappings {
		geoMappings = append(geoMappings, v.(string))
	}
	if len(geoMappings) > 0 {
		params.Properties.GeoMapping = &geoMappings
	}

	if _, err := client.CreateOrUpdate(ctx, id, params); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceAzureEndpointRead(d, meta)
}

func resourceAzureEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			d.Set("target_resource_id", props.TargetResourceId)
			d.Set("weight", props.Weight)
			d.Set("priority", props.Priority)
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

func resourceAzureEndpointUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.EndpointsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := profiles.ParseTrafficManagerProfileID(d.Get("profile_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `profile_id`: %+v", err)
	}

	id := endpoints.NewEndpointTypeID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.TrafficManagerProfileName, endpoints.EndpointTypeAzureEndpoints, d.Get("name").(string))

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

	if d.HasChange("target_resource_id") {
		params.Properties.TargetResourceId = utils.String(d.Get("target_resource_id").(string))
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
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceAzureEndpointRead(d, meta)
}

func resourceAzureEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
