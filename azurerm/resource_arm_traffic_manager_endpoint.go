package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/trafficmanager/mgmt/2017-05-01/trafficmanager"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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

	log.Printf("[INFO] preparing arguments for TrafficManager Endpoint creation.")

	name := d.Get("name").(string)
	endpointType := d.Get("type").(string)
	fullEndpointType := fmt.Sprintf("Microsoft.Network/TrafficManagerProfiles/%s", endpointType)
	profileName := d.Get("profile_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	params := trafficmanager.Endpoint{
		Name:               &name,
		Type:               &fullEndpointType,
		EndpointProperties: getArmTrafficManagerEndpointProperties(d),
	}

	ctx := meta.(*ArmClient).StopContext
	_, err := client.CreateOrUpdate(ctx, resourceGroup, profileName, endpointType, name, params)
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

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	endpointType := d.Get("type").(string)
	profileName := id.Path["trafficManagerProfiles"]

	// endpoint name is keyed by endpoint type in ARM ID
	name := id.Path[endpointType]
	ctx := meta.(*ArmClient).StopContext
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
