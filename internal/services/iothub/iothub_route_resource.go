// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	devices "github.com/tombuildsstuff/kermit/sdk/iothub/2022-04-30-preview/iothub"
)

func resourceIotHubRoute() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubRouteCreateUpdate,
		Read:   resourceIotHubRouteRead,
		Update: resourceIotHubRouteCreateUpdate,
		Delete: resourceIotHubRouteDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.IoTHubRouteV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.RouteID(id)
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-_.a-zA-Z0-9]{1,64}$"),
					"Route Name name can only include alphanumeric characters, periods, underscores, hyphens, has a maximum length of 64 characters, and must be unique.",
				),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"iothub_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"source": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(devices.RoutingSourceDeviceConnectionStateEvents),
					string(devices.RoutingSourceDeviceJobLifecycleEvents),
					string(devices.RoutingSourceDeviceLifecycleEvents),
					string(devices.RoutingSourceDeviceMessages),
					string(devices.RoutingSourceDigitalTwinChangeEvents),
					string(devices.RoutingSourceInvalid),
					string(devices.RoutingSourceTwinChangeEvents),
				}, false),
			},
			"condition": {
				// The condition is a string value representing device-to-cloud message routes query expression
				// https://docs.microsoft.com/en-us/azure/iot-hub/iot-hub-devguide-query-language#device-to-cloud-message-routes-query-expressions
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "true",
			},
			"endpoint_names": {
				Type: pluginsdk.TypeList,
				// Currently only one endpoint is allowed. With that comment from Microsoft, we'll leave this open to enhancement when they add multiple endpoint support.
				MaxItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				Required: true,
			},
			"enabled": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceIotHubRouteCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewRouteID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_name").(string), d.Get("name").(string))

	locks.ByName(id.IotHubName, IothubResourceName)
	defer locks.UnlockByName(id.IotHubName, IothubResourceName)

	iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", id.IotHubName, id.ResourceGroup)
		}

		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
	}

	source := devices.RoutingSource(d.Get("source").(string))
	condition := d.Get("condition").(string)
	endpointNamesRaw := d.Get("endpoint_names").([]interface{})
	isEnabled := d.Get("enabled").(bool)

	route := devices.RouteProperties{
		Name:          &id.Name,
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
			if strings.EqualFold(*existingRoute.Name, id.Name) {
				if d.IsNewResource() {
					return tf.ImportAsExistsError("azurerm_iothub_route", id.ID())
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
		return fmt.Errorf("unable to find %s", id)
	}

	routing.Routes = &routes

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, iothub, "")
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the completion of the creating/updating of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIotHubRouteRead(d, meta)
}

func resourceIotHubRouteRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.RouteID(d.Id())
	if err != nil {
		return err
	}

	iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("iothub_name", id.IotHubName)
	d.Set("resource_group_name", id.ResourceGroup)

	if iothub.Properties == nil || iothub.Properties.Routing == nil {
		d.SetId("")
		return nil
	}

	exist := false

	if routes := iothub.Properties.Routing.Routes; routes != nil {
		for _, route := range *routes {
			if route.Name != nil {
				if strings.EqualFold(*route.Name, id.Name) {
					exist = true
					d.Set("source", route.Source)
					d.Set("condition", route.Condition)
					d.Set("enabled", route.IsEnabled)
					d.Set("endpoint_names", route.EndpointNames)
				}
			}
		}
	}

	if !exist {
		d.SetId("")
	}

	return nil
}

func resourceIotHubRouteDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.RouteID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.IotHubName, IothubResourceName)
	defer locks.UnlockByName(id.IotHubName, IothubResourceName)

	iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", id.IotHubName, id.ResourceGroup)
		}

		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
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
			if !strings.EqualFold(*route.Name, id.Name) {
				updatedRoutes = append(updatedRoutes, route)
			}
		}
	}

	iothub.Properties.Routing.Routes = &updatedRoutes

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, iothub, "")
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to finish updating: %+v", id, err)
	}

	return nil
}
