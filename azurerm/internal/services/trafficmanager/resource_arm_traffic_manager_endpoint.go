package trafficmanager

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/trafficmanager/mgmt/2018-04-01/trafficmanager"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmTrafficManagerEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmTrafficManagerEndpointCreateUpdate,
		Read:   resourceArmTrafficManagerEndpointRead,
		Update: resourceArmTrafficManagerEndpointCreateUpdate,
		Delete: resourceArmTrafficManagerEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"profile_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"azureEndpoints",
					"nestedEndpoints",
					"externalEndpoints",
				}, false),
			},

			"target": {
				Type:     schema.TypeString,
				Optional: true,
				// when targeting an Azure resource the FQDN of that resource will be set as the target
				Computed: true,
			},

			"target_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"endpoint_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(trafficmanager.EndpointStatusDisabled),
					string(trafficmanager.EndpointStatusEnabled),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 1000),
			},

			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 1000),
			},

			// when targeting an Azure resource the location of that resource will be set on the endpoint
			"endpoint_location": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				StateFunc:        location.StateFunc,
				DiffSuppressFunc: location.DiffSuppressFunc,
			},

			"min_child_endpoints": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"geo_mappings": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},

			"endpoint_monitor_status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"custom_header": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
					},
				},
			},

			"subnet": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"first": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.IPv4Address,
						},
						"last": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.IPv4Address,
						},
						"scope": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 32),
						},
					},
				},
			},
		},
	}
}

func resourceArmTrafficManagerEndpointCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.EndpointsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for TrafficManager Endpoint creation.")

	name := d.Get("name").(string)
	endpointType := d.Get("type").(string)
	fullEndpointType := fmt.Sprintf("Microsoft.Network/TrafficManagerProfiles/%s", endpointType)
	profileName := d.Get("profile_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, profileName, endpointType, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Traffic Manager Endpoint %q (Resource Group %q): %v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_traffic_manager_endpoint", *existing.ID)
		}
	}

	params := trafficmanager.Endpoint{
		Name:               &name,
		Type:               &fullEndpointType,
		EndpointProperties: getArmTrafficManagerEndpointProperties(d),
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, profileName, endpointType, name, params); err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, profileName, endpointType, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Traffic Manager Endpoint %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmTrafficManagerEndpointRead(d, meta)
}

func resourceArmTrafficManagerEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.EndpointsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup

	// lookup endpointType in Azure ID path
	var endpointType string
	typeRegex := regexp.MustCompile("azureEndpoints|externalEndpoints|nestedEndpoints")
	for k := range id.Path {
		if typeRegex.MatchString(k) {
			endpointType = k
		}
	}
	profileName := id.Path["trafficManagerProfiles"]
	name := id.Path[endpointType]

	resp, err := client.Get(ctx, resGroup, profileName, endpointType, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on TrafficManager Endpoint %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("resource_group_name", resGroup)
	d.Set("name", resp.Name)
	d.Set("type", endpointType)
	d.Set("profile_name", profileName)

	if props := resp.EndpointProperties; props != nil {
		d.Set("endpoint_status", string(props.EndpointStatus))
		d.Set("target_resource_id", props.TargetResourceID)
		d.Set("target", props.Target)
		d.Set("weight", props.Weight)
		d.Set("priority", props.Priority)
		d.Set("endpoint_location", props.EndpointLocation)
		d.Set("endpoint_monitor_status", props.EndpointMonitorStatus)
		d.Set("min_child_endpoints", props.MinChildEndpoints)
		d.Set("geo_mappings", props.GeoMapping)
		if err := d.Set("subnet", flattenAzureRMTrafficManagerEndpointSubnetConfig(props.Subnets)); err != nil {
			return fmt.Errorf("Error setting `subnet`: %s", err)
		}
		if err := d.Set("custom_header", flattenAzureRMTrafficManagerEndpointCustomHeaderConfig(props.CustomHeaders)); err != nil {
			return fmt.Errorf("Error setting `custom_header`: %s", err)
		}
	}
	return nil
}

func flattenAzureRMTrafficManagerEndpointSubnetConfig(input *[]trafficmanager.EndpointPropertiesSubnetsItem) []interface{} {
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

func flattenAzureRMTrafficManagerEndpointCustomHeaderConfig(input *[]trafficmanager.EndpointPropertiesCustomHeadersItem) []interface{} {
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

func resourceArmTrafficManagerEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).TrafficManager.EndpointsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	endpointType := d.Get("type").(string)
	profileName := id.Path["trafficManagerProfiles"]

	// endpoint name is keyed by endpoint type in ARM ID
	name := id.Path[endpointType]
	resp, err := client.Delete(ctx, resGroup, profileName, endpointType, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}

		return err
	}

	return nil
}

func getArmTrafficManagerEndpointProperties(d *schema.ResourceData) *trafficmanager.EndpointProperties {
	target := d.Get("target").(string)
	status := d.Get("endpoint_status").(string)

	endpointProps := trafficmanager.EndpointProperties{
		Target:         &target,
		EndpointStatus: trafficmanager.EndpointStatus(status),
	}

	if resourceId := d.Get("target_resource_id").(string); resourceId != "" {
		endpointProps.TargetResourceID = utils.String(resourceId)
		// NOTE: Workaround for upstream behavior: if the target is blank instead of nil, the REST API will throw a 500 error
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

	if minChildEndpoints := d.Get("min_child_endpoints").(int); minChildEndpoints != 0 {
		mci64 := int64(minChildEndpoints)
		endpointProps.MinChildEndpoints = &mci64
	}

	subnetSlice := make([]trafficmanager.EndpointPropertiesSubnetsItem, 0)
	for _, subnet := range d.Get("subnet").([]interface{}) {
		subnetBlock := subnet.(map[string]interface{})
		if subnetBlock["scope"].(int) == 0 && subnetBlock["first"].(string) != "0.0.0.0" {
			subnetSlice = append(subnetSlice, trafficmanager.EndpointPropertiesSubnetsItem{
				First: utils.String(subnetBlock["first"].(string)),
				Last:  utils.String(subnetBlock["last"].(string)),
			})
		} else {
			subnetSlice = append(subnetSlice, trafficmanager.EndpointPropertiesSubnetsItem{
				First: utils.String(subnetBlock["first"].(string)),
				Scope: utils.Int32(int32(subnetBlock["scope"].(int))),
			})
		}
	}
	if len(subnetSlice) > 0 {
		endpointProps.Subnets = &subnetSlice
	}

	headerSlice := make([]trafficmanager.EndpointPropertiesCustomHeadersItem, 0)
	for _, header := range d.Get("custom_header").([]interface{}) {
		headerBlock := header.(map[string]interface{})
		headerSlice = append(headerSlice, trafficmanager.EndpointPropertiesCustomHeadersItem{
			Name:  utils.String(headerBlock["name"].(string)),
			Value: utils.String(headerBlock["value"].(string)),
		})
	}
	if len(headerSlice) > 0 {
		endpointProps.CustomHeaders = &headerSlice
	}

	return &endpointProps
}
