package iothub

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/iothub/mgmt/2020-03-01/devices"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"iothub_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
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
				Required: true,
				// Currently only one endpoint is allowed. With that comment from Microsoft, we'll leave this open to enhancement when they add multiple endpoint support.
				MaxItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.IoTHubEndpointName,
				},
			},

			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceArmIotHubFallbackRouteCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	iothubName := d.Get("iothub_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(iothubName, IothubResourceName)
	defer locks.UnlockByName(iothubName, IothubResourceName)

	iothub, err := client.Get(ctx, resourceGroup, iothubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", iothubName, resourceGroup)
		}

		return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	// NOTE: this resource intentionally doesn't support Requires Import
	//       since a fallback route is created by default

	routing := iothub.Properties.Routing

	if routing == nil {
		routing = &devices.RoutingProperties{}
	}

	routing.FallbackRoute = &devices.FallbackRouteProperties{
		Source:        utils.String(string(devices.RoutingSourceDeviceMessages)),
		Condition:     utils.String(d.Get("condition").(string)),
		EndpointNames: utils.ExpandStringSlice(d.Get("endpoint_names").([]interface{})),
		IsEnabled:     utils.Bool(d.Get("enabled").(bool)),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, iothubName, iothub, "")
	if err != nil {
		return fmt.Errorf("Error creating/updating IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the completion of the creating/updating of IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	resourceId := fmt.Sprintf("%s/FallbackRoute/defined", *iothub.ID)
	d.SetId(resourceId)

	return resourceArmIotHubFallbackRouteRead(d, meta)
}

func resourceArmIotHubFallbackRouteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parsedIothubRouteId, err := azure.ParseAzureResourceID(d.Id())
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

	if props := iothub.Properties; props != nil {
		if routing := props.Routing; routing != nil {
			if fallbackRoute := routing.FallbackRoute; fallbackRoute != nil {
				d.Set("condition", fallbackRoute.Condition)
				d.Set("enabled", fallbackRoute.IsEnabled)
				d.Set("endpoint_names", fallbackRoute.EndpointNames)
			}
		}
	}

	return nil
}

func resourceArmIotHubFallbackRouteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parsedIothubRouteId, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := parsedIothubRouteId.ResourceGroup
	iothubName := parsedIothubRouteId.Path["IotHubs"]

	locks.ByName(iothubName, IothubResourceName)
	defer locks.UnlockByName(iothubName, IothubResourceName)

	iothub, err := client.Get(ctx, resourceGroup, iothubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", iothubName, resourceGroup)
		}

		return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	if iothub.Properties == nil || iothub.Properties.Routing == nil || iothub.Properties.Routing.FallbackRoute == nil {
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
