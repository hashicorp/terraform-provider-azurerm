package azurerm

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/iothub/mgmt/2018-12-01-preview/devices"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmIotHubRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIotHubRouteCreateUpdate,
		Read:   resourceArmIotHubRouteRead,
		Update: resourceArmIotHubRouteCreateUpdate,
		Delete: resourceArmIotHubRouteDelete,
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
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-_.a-zA-Z0-9]{1,64}$"),
					"Route Name name can only include alphanumeric characters, periods, underscores, hyphens, has a maximum length of 64 characters, and must be unique.",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"iothub_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"source": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(devices.RoutingSourceDeviceJobLifecycleEvents),
					string(devices.RoutingSourceDeviceLifecycleEvents),
					string(devices.RoutingSourceDeviceMessages),
					string(devices.RoutingSourceInvalid),
					string(devices.RoutingSourceTwinChangeEvents),
				}, false),
			},
			"condition": {
				// The condition is a string value representing device-to-cloud message routes query expression
				// https://docs.microsoft.com/en-us/azure/iot-hub/iot-hub-devguide-query-language#device-to-cloud-message-routes-query-expressions
				Type:     schema.TypeString,
				Optional: true,
				Default:  "true",
			},
			"endpoint_names": {
				Type: schema.TypeList,
				// Currently only one endpoint is allowed. With that comment from Microsoft, we'll leave this open to enhancement when they add multiple endpoint support.
				MaxItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.NoEmptyStrings,
				},
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceArmIotHubRouteCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	iothubName := d.Get("iothub_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(iothubName, iothubResourceName)
	defer locks.UnlockByName(iothubName, iothubResourceName)

	iothub, err := client.Get(ctx, resourceGroup, iothubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", iothubName, resourceGroup)
		}

		return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	routeName := d.Get("name").(string)

	resourceId := fmt.Sprintf("%s/Routes/%s", *iothub.ID, routeName)

	source := devices.RoutingSource(d.Get("source").(string))
	condition := d.Get("condition").(string)
	endpointNamesRaw := d.Get("endpoint_names").([]interface{})
	isEnabled := d.Get("enabled").(bool)

	route := devices.RouteProperties{
		Name:          &routeName,
		Source:        source,
		Condition:     &condition,
		EndpointNames: utils.ExpandStringSlice(endpointNamesRaw),
		IsEnabled:     &isEnabled,
	}

	routing := iothub.Properties.Routing

	if routing == nil {
		routing = &devices.RoutingProperties{}
	}

	if routing.Routes == nil {
		routes := make([]devices.RouteProperties, 0)
		routing.Routes = &routes
	}

	routes := make([]devices.RouteProperties, 0)

	alreadyExists := false
	for _, existingRoute := range *routing.Routes {
		if existingRoute.Name != nil {
			if strings.EqualFold(*existingRoute.Name, routeName) {
				if d.IsNewResource() && features.ShouldResourcesBeImported() {
					return tf.ImportAsExistsError("azurerm_iothub_route", resourceId)
				}
				routes = append(routes, route)
				alreadyExists = true
			} else {
				routes = append(routes, existingRoute)
			}
		}
	}

	if d.IsNewResource() {
		routes = append(routes, route)
	} else if !alreadyExists {
		return fmt.Errorf("Unable to find Route %q defined for IotHub %q (Resource Group %q)", routeName, iothubName, resourceGroup)
	}

	routing.Routes = &routes

	future, err := client.CreateOrUpdate(ctx, resourceGroup, iothubName, iothub, "")
	if err != nil {
		return fmt.Errorf("Error creating/updating IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the completion of the creating/updating of IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	d.SetId(resourceId)

	return resourceArmIotHubRouteRead(d, meta)
}

func resourceArmIotHubRouteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	parsedIothubRouteId, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := parsedIothubRouteId.ResourceGroup
	iothubName := parsedIothubRouteId.Path["IotHubs"]
	routeName := parsedIothubRouteId.Path["Routes"]

	iothub, err := client.Get(ctx, resourceGroup, iothubName)
	if err != nil {
		return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	d.Set("name", routeName)
	d.Set("iothub_name", iothubName)
	d.Set("resource_group_name", resourceGroup)

	if iothub.Properties == nil || iothub.Properties.Routing == nil {
		return nil
	}

	if routes := iothub.Properties.Routing.Routes; routes != nil {
		for _, route := range *routes {
			if route.Name != nil {
				if strings.EqualFold(*route.Name, routeName) {
					d.Set("source", route.Source)
					d.Set("condition", route.Condition)
					d.Set("enabled", route.IsEnabled)
					d.Set("endpoint_names", route.EndpointNames)
				}
			}
		}
	}

	return nil
}

func resourceArmIotHubRouteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	parsedIothubRouteId, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := parsedIothubRouteId.ResourceGroup
	iothubName := parsedIothubRouteId.Path["IotHubs"]
	routeName := parsedIothubRouteId.Path["Routes"]

	locks.ByName(iothubName, iothubResourceName)
	defer locks.UnlockByName(iothubName, iothubResourceName)

	iothub, err := client.Get(ctx, resourceGroup, iothubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", iothubName, resourceGroup)
		}

		return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	if iothub.Properties == nil || iothub.Properties.Routing == nil {
		return nil
	}
	routes := iothub.Properties.Routing.Routes

	if routes == nil {
		return nil
	}

	updatedRoutes := make([]devices.RouteProperties, 0)
	for _, route := range *routes {
		if route.Name != nil {
			if !strings.EqualFold(*route.Name, routeName) {
				updatedRoutes = append(updatedRoutes, route)
			}
		}
	}

	iothub.Properties.Routing.Routes = &updatedRoutes

	future, err := client.CreateOrUpdate(ctx, resourceGroup, iothubName, iothub, "")
	if err != nil {
		return fmt.Errorf("Error updating IotHub %q (Resource Group %q) with Route %q: %+v", iothubName, resourceGroup, routeName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for IotHub %q (Resource Group %q) to finish updating Route %q: %+v", iothubName, resourceGroup, routeName, err)
	}

	return nil
}
