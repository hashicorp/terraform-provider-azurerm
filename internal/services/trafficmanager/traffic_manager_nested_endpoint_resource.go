// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package trafficmanager

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
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

func resourceNestedEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNestedEndpointCreateUpdate,
		Read:   resourceNestedEndpointRead,
		Update: resourceNestedEndpointCreateUpdate,
		Delete: resourceNestedEndpointDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			endpointType, err := endpoints.ParseEndpointTypeID(id)
			if err != nil {
				return err
			}

			if endpointType.EndpointType != endpoints.EndpointTypeNestedEndpoints {
				return fmt.Errorf("this resource only supports `NestedEndpoints` but got %s", string(endpointType.EndpointType))
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

			"minimum_child_endpoints": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"minimum_required_child_endpoints_ipv4": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
			},

			"minimum_required_child_endpoints_ipv6": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
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

func resourceNestedEndpointCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.EndpointsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := profiles.ParseTrafficManagerProfileID(d.Get("profile_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `profile_id`: %+v", err)
	}

	id := endpoints.NewEndpointTypeID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.TrafficManagerProfileName, endpoints.EndpointTypeNestedEndpoints, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_traffic_manager_nested_endpoint", id.ID())
		}
	}

	status := endpoints.EndpointStatusEnabled
	if !d.Get("enabled").(bool) {
		status = endpoints.EndpointStatusDisabled
	}

	params := endpoints.Endpoint{
		Name: utils.String(id.EndpointName),
		Type: utils.String(fmt.Sprintf("Microsoft.Network/trafficManagerProfiles/%s", endpoints.EndpointTypeNestedEndpoints)),
		Properties: &endpoints.EndpointProperties{
			CustomHeaders:     expandEndpointCustomHeaderConfig(d.Get("custom_header").([]interface{})),
			EndpointStatus:    &status,
			MinChildEndpoints: utils.Int64(int64(d.Get("minimum_child_endpoints").(int))),
			TargetResourceId:  utils.String(d.Get("target_resource_id").(string)),
			Subnets:           expandEndpointSubnetConfig(d.Get("subnet").([]interface{})),
		},
	}

	if weight := d.Get("weight").(int); weight != 0 {
		params.Properties.Weight = utils.Int64(int64(weight))
	}

	minChildEndpointsIPv4 := d.Get("minimum_required_child_endpoints_ipv4").(int)
	if minChildEndpointsIPv4 > 0 {
		params.Properties.MinChildEndpointsIPv4 = utils.Int64(int64(minChildEndpointsIPv4))
	}

	minChildEndpointsIPv6 := d.Get("minimum_required_child_endpoints_ipv6").(int)
	if minChildEndpointsIPv6 > 0 {
		params.Properties.MinChildEndpointsIPv6 = utils.Int64(int64(minChildEndpointsIPv6))
	}

	if priority := d.Get("priority").(int); priority != 0 {
		params.Properties.Priority = utils.Int64(int64(priority))
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
	return resourceNestedEndpointRead(d, meta)
}

func resourceNestedEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			d.Set("minimum_child_endpoints", props.MinChildEndpoints)
			d.Set("minimum_required_child_endpoints_ipv4", props.MinChildEndpointsIPv4)
			d.Set("minimum_required_child_endpoints_ipv6", props.MinChildEndpointsIPv6)
			d.Set("priority", props.Priority)
			d.Set("endpoint_location", props.EndpointLocation)
			d.Set("geo_mappings", props.GeoMapping)

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

func resourceNestedEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
