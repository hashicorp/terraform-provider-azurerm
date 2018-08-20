package azurerm

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/trafficmanager/mgmt/2017-05-01/trafficmanager"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmTrafficManagerEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmTrafficManagerEndpointCreate,
		Read:   resourceArmTrafficManagerEndpointRead,
		Update: resourceArmTrafficManagerEndpointCreate,
		Delete: resourceArmTrafficManagerEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 30),
			Update: schema.DefaultTimeout(time.Minute * 30),
			Delete: schema.DefaultTimeout(time.Minute * 30),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"profile_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameDiffSuppressSchema(),

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
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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
				StateFunc:        azureRMNormalizeLocation,
				DiffSuppressFunc: azureRMSuppressLocationDiff,
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
		},
	}
}

func resourceArmTrafficManagerEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).trafficManagerEndpointsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for TrafficManager Endpoint creation.")

	name := d.Get("name").(string)
	profileName := d.Get("profile_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		// first check if there's one in this subscription requiring import
		resp, err := client.Get(ctx, resourceGroup, profileName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for the existence of Traffic Manager Endpoint %q (Profile %q / Resource Group %q): %+v", name, profileName, resourceGroup, err)
			}
		}

		if resp.ID != nil {
			return tf.ImportAsExistsError("azurerm_traffic_manager_endpoint", *resp.ID)
		}
	}

	endpointType := d.Get("type").(string)
	fullEndpointType := fmt.Sprintf("Microsoft.Network/TrafficManagerProfiles/%s", endpointType)

	params := trafficmanager.Endpoint{
		Name:               &name,
		Type:               &fullEndpointType,
		EndpointProperties: getArmTrafficManagerEndpointProperties(d),
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(tf.TimeoutForCreateUpdate(d)))
	defer cancel()
	_, err := client.CreateOrUpdate(waitCtx, resourceGroup, profileName, endpointType, name, params)
	if err != nil {
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
	client := meta.(*ArmClient).trafficManagerEndpointsClient

	id, err := parseAzureResourceID(d.Id())
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

	ctx := meta.(*ArmClient).StopContext
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
	}

	return nil
}

func resourceArmTrafficManagerEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).trafficManagerEndpointsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	// TODO: fix the ID
	resGroup := id.ResourceGroup
	endpointType := d.Get("type").(string)
	profileName := id.Path["trafficManagerProfiles"]

	// endpoint name is keyed by endpoint type in ARM ID
	name := id.Path[endpointType]
	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutDelete))
	defer cancel()
	resp, err := client.Delete(waitCtx, resGroup, profileName, endpointType, name)
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
		//TODO? Workaround for upstream behavior: if the target is blank instead of nil, the REST API will throw a 500 error. Remove if/when no longer necessary
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

	return &endpointProps
}
