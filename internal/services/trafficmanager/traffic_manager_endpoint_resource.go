package trafficmanager

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/sdk/2018-08-01/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// TODO: split and deprecate this resource prior to 3.0

func resourceArmTrafficManagerEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmTrafficManagerEndpointCreateUpdate,
		Read:   resourceArmTrafficManagerEndpointRead,
		Update: resourceArmTrafficManagerEndpointCreateUpdate,
		Delete: resourceArmTrafficManagerEndpointDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		DeprecationMessage: "The resource 'azurerm_traffic_manager_endpoint' has been deprecated in favour of 'azurerm_traffic_manager_azure_endpoint', 'azurerm_traffic_manager_external_endpoint', and 'azurerm_traffic_manager_nested_endpoint' and will be removed in version 3.0 of the Azure Provider.",

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

			"type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(endpoints.EndpointTypeAzureEndpoints),
					string(endpoints.EndpointTypeNestedEndpoints),
					string(endpoints.EndpointTypeExternalEndpoints),
				}, !features.ThreePointOh()),
				DiffSuppressFunc: suppress.CaseDifferenceV2Only,
			},

			"target": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				// when targeting an Azure resource the FQDN of that resource will be set as the target
				Computed: true,
			},

			"target_resource_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"endpoint_status": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(endpoints.EndpointStatusEnabled),
				ValidateFunc: validation.StringInSlice([]string{
					string(endpoints.EndpointStatusDisabled),
					string(endpoints.EndpointStatusEnabled),
				}, !features.ThreePointOh()),
				DiffSuppressFunc: suppress.CaseDifferenceV2Only,
			},

			"weight": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 1000),
			},

			"priority": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 1000),
			},

			// when targeting an Azure resource the location of that resource will be set on the endpoint
			"endpoint_location": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Computed:         true,
				StateFunc:        location.StateFunc,
				DiffSuppressFunc: location.DiffSuppressFunc,
			},

			// TODO 3.0: rename this to `minimum_child_endpoints` to align with the other twos below.
			"min_child_endpoints": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
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

			"geo_mappings": {
				Type:     pluginsdk.TypeList,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Optional: true,
			},

			"endpoint_monitor_status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"custom_header": {
				Type:     pluginsdk.TypeList,
				ForceNew: true,
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

func resourceArmTrafficManagerEndpointCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.EndpointsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for TrafficManager Endpoint creation.")

	name := d.Get("name").(string)
	endpointType := d.Get("type").(string)
	fullEndpointType := fmt.Sprintf("Microsoft.Network/TrafficManagerProfiles/%s", endpointType)
	profileName := d.Get("profile_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := endpoints.NewEndpointTypeID(subscriptionId, resourceGroup, profileName, endpoints.EndpointType(endpointType), name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_traffic_manager_endpoint", id.ID())
		}
	}

	params := endpoints.Endpoint{
		Name:       &name,
		Type:       &fullEndpointType,
		Properties: getArmTrafficManagerEndpointProperties(d),
	}

	if _, err := client.CreateOrUpdate(ctx, id, params); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceArmTrafficManagerEndpointRead(d, meta)
}

func resourceArmTrafficManagerEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.EndpointsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := endpoints.ParseEndpointTypeIDInsensitively(d.Id())
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
	d.Set("type", id.EndpointType)
	d.Set("profile_name", id.ProfileName)
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			endpointStatus := ""
			if props.EndpointStatus != nil {
				endpointStatus = string(*props.EndpointStatus)
			}
			d.Set("endpoint_status", endpointStatus)
			d.Set("target_resource_id", props.TargetResourceId)
			d.Set("target", props.Target)
			d.Set("weight", props.Weight)
			d.Set("priority", props.Priority)
			d.Set("endpoint_location", props.EndpointLocation)

			endPointMonitorStatus := ""
			if props.EndpointMonitorStatus != nil {
				endPointMonitorStatus = string(*props.EndpointMonitorStatus)
			}
			d.Set("endpoint_monitor_status", endPointMonitorStatus)
			d.Set("min_child_endpoints", props.MinChildEndpoints)
			d.Set("minimum_required_child_endpoints_ipv4", props.MinChildEndpointsIPv4)
			d.Set("minimum_required_child_endpoints_ipv6", props.MinChildEndpointsIPv6)
			d.Set("geo_mappings", props.GeoMapping)
			if err := d.Set("subnet", flattenEndpointSubnetConfig(props.Subnets)); err != nil {
				return fmt.Errorf("setting `subnet`: %s", err)
			}
			if err := d.Set("custom_header", flattenEndpointCustomHeaderConfig(props.CustomHeaders)); err != nil {
				return fmt.Errorf("setting `custom_header`: %s", err)
			}
		}
	}

	return nil
}

func resourceArmTrafficManagerEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.EndpointsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := endpoints.ParseEndpointTypeIDInsensitively(d.Id())
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

func getArmTrafficManagerEndpointProperties(d *pluginsdk.ResourceData) *endpoints.EndpointProperties {
	target := d.Get("target").(string)
	status := endpoints.EndpointStatus(d.Get("endpoint_status").(string))

	endpointProps := endpoints.EndpointProperties{
		Target:         &target,
		EndpointStatus: &status,
	}

	if resourceId := d.Get("target_resource_id").(string); resourceId != "" {
		endpointProps.TargetResourceId = utils.String(resourceId)
		// NOTE: Workaround for upstream behaviour: if the target is blank instead of nil, the REST API will throw a 500 error
		if target == "" {
			endpointProps.Target = nil
		}
	}

	if location := d.Get("endpoint_location").(string); location != "" {
		endpointProps.EndpointLocation = utils.String(location)
	}

	inputMappings := d.Get("geo_mappings").([]interface{})
	geoMappings := make([]string, 0)
	for _, v := range inputMappings {
		geoMappings = append(geoMappings, v.(string))
	}
	if len(geoMappings) > 0 {
		endpointProps.GeoMapping = &geoMappings
	}

	if weight := d.Get("weight").(int); weight != 0 {
		endpointProps.Weight = utils.Int64(int64(weight))
	}

	if priority := d.Get("priority").(int); priority != 0 {
		endpointProps.Priority = utils.Int64(int64(priority))
	}

	minChildEndpoints := d.Get("min_child_endpoints").(int)
	if minChildEndpoints > 0 {
		endpointProps.MinChildEndpoints = utils.Int64(int64(minChildEndpoints))
	}

	minChildEndpointsIPv4 := d.Get("minimum_required_child_endpoints_ipv4").(int)
	if minChildEndpointsIPv4 > 0 {
		endpointProps.MinChildEndpointsIPv4 = utils.Int64(int64(minChildEndpointsIPv4))
	}

	minChildEndpointsIPv6 := d.Get("minimum_required_child_endpoints_ipv6").(int)
	if minChildEndpointsIPv6 > 0 {
		endpointProps.MinChildEndpointsIPv6 = utils.Int64(int64(minChildEndpointsIPv6))
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
		endpointProps.Subnets = &subnetSlice
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
		endpointProps.CustomHeaders = &headerSlice
	}

	return &endpointProps
}

func flattenEndpointSubnetConfig(input *[]endpoints.EndpointPropertiesSubnetsInlined) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}
	for _, subnet := range *input {
		flatSubnet := make(map[string]interface{}, 3)
		if subnet.First != nil {
			flatSubnet["first"] = *subnet.First
		}
		if subnet.Last != nil {
			flatSubnet["last"] = *subnet.Last
		}
		if subnet.Scope != nil {
			flatSubnet["scope"] = int(*subnet.Scope)
		}
		result = append(result, flatSubnet)
	}
	return result
}

func flattenEndpointCustomHeaderConfig(input *[]endpoints.EndpointPropertiesCustomHeadersInlined) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}
	for _, header := range *input {
		flatHeader := make(map[string]interface{}, 2)
		if header.Name != nil {
			flatHeader["name"] = *header.Name
		}
		if header.Value != nil {
			flatHeader["value"] = *header.Value
		}
		result = append(result, flatHeader)
	}
	return result
}
