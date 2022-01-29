package trafficmanager

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"

	"github.com/hashicorp/terraform-provider-azurerm/utils"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/sdk/2018-08-01/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	azSchema "github.com/hashicorp/terraform-provider-azurerm/internal/tf/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmTrafficManagerExternalEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmTrafficManagerExternalEndpointCreateUpdate,
		Read:   resourceArmTrafficManagerExternalEndpointRead,
		Update: resourceArmTrafficManagerExternalEndpointCreateUpdate,
		Delete: resourceArmTrafficManagerExternalEndpointDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := endpoints.ParseEndpointTypeID(id)
			return err
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

			"profile_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"target": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"weight": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 1000),
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
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
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 1000),
			},

			"endpoint_location": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
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
							ValidateFunc: validate.IPv4Address,
						},
						"last": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validate.IPv4Address,
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

func resourceArmTrafficManagerExternalEndpointCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.EndpointsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for TrafficManager Endpoint creation.")

	name := d.Get("name").(string)
	fullEndpointType := fmt.Sprintf("Microsoft.Network/TrafficManagerProfiles/%s", endpoints.EndpointTypeExternalEndpoints)
	profileName := d.Get("profile_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := endpoints.NewEndpointTypeID(subscriptionId, resourceGroup, profileName, endpoints.EndpointTypeExternalEndpoints, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_traffic_manager_external_endpoint", id.ID())
		}
	}

	status := endpoints.EndpointStatusEnabled
	if !d.Get("enabled").(bool) {
		status = endpoints.EndpointStatusDisabled
	}

	params := endpoints.Endpoint{
		Name: &name,
		Type: &fullEndpointType,
		Properties: &endpoints.EndpointProperties{
			Target:         utils.String(d.Get("target").(string)),
			EndpointStatus: &status,
			Weight:         utils.Int64(int64(d.Get("weight").(int))),
		},
	}

	headerSlice := make([]endpoints.EndpointPropertiesCustomHeadersInlined, 0)
	for _, header := range d.Get("custom_header").([]interface{}) {
		headerBlock := header.(map[string]interface{})
		headerSlice = append(headerSlice, endpoints.EndpointPropertiesCustomHeadersInlined{
			Name:  utils.String(headerBlock["name"].(string)),
			Value: utils.String(headerBlock["value"].(string)),
		})
	}

	if len(headerSlice) > 0 {
		params.Properties.CustomHeaders = &headerSlice
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

	subnetSlice := make([]endpoints.EndpointPropertiesSubnetsInlined, 0)
	for _, subnet := range d.Get("subnet").([]interface{}) {
		subnetBlock := subnet.(map[string]interface{})
		if subnetBlock["scope"].(int) == 0 && subnetBlock["first"].(string) != "0.0.0.0" {
			subnetSlice = append(subnetSlice, endpoints.EndpointPropertiesSubnetsInlined{
				First: utils.String(subnetBlock["first"].(string)),
				Last:  utils.String(subnetBlock["last"].(string)),
			})
		} else {
			subnetSlice = append(subnetSlice, endpoints.EndpointPropertiesSubnetsInlined{
				First: utils.String(subnetBlock["first"].(string)),
				Scope: utils.Int64(int64(subnetBlock["scope"].(int))),
			})
		}
	}
	if len(subnetSlice) > 0 {
		params.Properties.Subnets = &subnetSlice
	}

	if _, err := client.CreateOrUpdate(ctx, id, params); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceArmTrafficManagerExternalEndpointRead(d, meta)
}

func resourceArmTrafficManagerExternalEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("profile_name", id.ProfileName)
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

			if err := d.Set("custom_header", flattenAzureRMTrafficManagerEndpointCustomHeaderConfig(props.CustomHeaders)); err != nil {
				return fmt.Errorf("setting `custom_header`: %s", err)
			}
			if err := d.Set("subnet", flattenAzureRMTrafficManagerEndpointSubnetConfig(props.Subnets)); err != nil {
				return fmt.Errorf("setting `subnet`: %s", err)
			}
		}
	}

	return nil
}

func resourceArmTrafficManagerExternalEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.EndpointsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := endpoints.ParseEndpointTypeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
