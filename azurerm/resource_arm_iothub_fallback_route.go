package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/iothub/mgmt/2018-12-01-preview/devices"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmIotHubFallbackRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIotHubFallbackRouteCreateUpdate,
		Read:   resourceArmIotHubFallbackRouteRead,
		Update: resourceArmIotHubFallbackRouteCreateUpdate,
		Delete: resourceArmIotHubFallbackRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"iothub_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"source": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "DeviceMessages",
				ValidateFunc: validation.StringInSlice([]string{
					"DeviceJobLifecycleEvents",
					"DeviceLifecycleEvents",
					"DeviceMessages",
					"Invalid",
					"TwinChangeEvents",
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
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringLenBetween(0, 64),
				},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceArmIotHubFallbackRouteCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).iothub.ResourceClient
	ctx := meta.(*ArmClient).StopContext

	iothubName := d.Get("iothub_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	azureRMLockByName(iothubName, iothubResourceName)
	defer azureRMUnlockByName(iothubName, iothubResourceName)

	iothub, err := client.Get(ctx, resourceGroup, iothubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", iothubName, resourceGroup)
		}

		return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	routing := iothub.Properties.Routing

	if routing == nil {
		routing = &devices.RoutingProperties{}
	}

	resourceId := fmt.Sprintf("%s/FallbackRoute/defined", *iothub.ID)

	if d.IsNewResource() && routing.FallbackRoute != nil && requireResourcesToBeImported {
		return tf.ImportAsExistsError("azurerm_iothub_fallback_route", resourceId)
	}

	source := d.Get("source").(string)
	condition := d.Get("condition").(string)
	endpointNamesRaw := d.Get("endpoint_names").([]interface{})
	isEnabled := d.Get("enabled").(bool)

	fallbackRoute := devices.FallbackRouteProperties{
		Source:        &source,
		Condition:     &condition,
		EndpointNames: utils.ExpandStringSlice(endpointNamesRaw),
		IsEnabled:     &isEnabled,
	}

	routing.FallbackRoute = &fallbackRoute

	future, err := client.CreateOrUpdate(ctx, resourceGroup, iothubName, iothub, "")
	if err != nil {
		return fmt.Errorf("Error creating/updating IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the completion of the creating/updating of IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	d.SetId(resourceId)

	return resourceArmIotHubFallbackRouteRead(d, meta)
}

func resourceArmIotHubFallbackRouteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).iothub.ResourceClient
	ctx := meta.(*ArmClient).StopContext

	parsedIothubRouteId, err := parseAzureResourceID(d.Id())

	if err != nil {
		return err
	}

	resourceGroup := parsedIothubRouteId.ResourceGroup
	iothubName := parsedIothubRouteId.Path["IotHubs"]

	iothub, err := client.Get(ctx, resourceGroup, iothubName)
	if err != nil {
		return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	d.Set("iothub_name", iothubName)
	d.Set("resource_group_name", resourceGroup)

	if iothub.Properties == nil || iothub.Properties.Routing == nil {
		return nil
	}

	route := iothub.Properties.Routing.FallbackRoute

	d.Set("source", route.Source)
	d.Set("condition", route.Condition)
	d.Set("enabled", route.IsEnabled)
	d.Set("endpoint_names", route.EndpointNames)

	return nil
}

func resourceArmIotHubFallbackRouteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).iothub.ResourceClient
	ctx := meta.(*ArmClient).StopContext

	parsedIothubRouteId, err := parseAzureResourceID(d.Id())

	if err != nil {
		return err
	}

	resourceGroup := parsedIothubRouteId.ResourceGroup
	iothubName := parsedIothubRouteId.Path["IotHubs"]

	azureRMLockByName(iothubName, iothubResourceName)
	defer azureRMUnlockByName(iothubName, iothubResourceName)

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

	iothub.Properties.Routing.FallbackRoute = nil

	future, err := client.CreateOrUpdate(ctx, resourceGroup, iothubName, iothub, "")
	if err != nil {
		return fmt.Errorf("Error updating IotHub %q (Resource Group %q) with Fallback Route: %+v", iothubName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for IotHub %q (Resource Group %q) to finish updating Fallback Route: %+v", iothubName, resourceGroup, err)
	}

	return nil
}
