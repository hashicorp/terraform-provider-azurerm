// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_route

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func resourceIotHubRouteCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewRouteID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_name").(string), d.Get("name").(string))

	locks.ByName(id.IotHubName, IothubResourceName)
	defer locks.UnlockByName(id.IotHubName, IothubResourceName)

	iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		if response.WasNotFound(iothub.Response.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", id.IotHubName, id.ResourceGroup)
		}

		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
	}

	source := devices.RoutingSource(d.Get("source").(string))
	condition := d.Get("condition").(string)
	endpointNamesRaw := d.Get("endpoint_names").([]any)
	isEnabled := d.Get("enabled").(bool)

	route := devices.RouteProperties{
		Name:          &id.Name,
		Source:        source,
		Condition:     &condition,
		EndpointNames: pluginsdk.ExpandStringSlice(endpointNamesRaw),
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
					if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
						return tf.ImportAsExistsError("azurerm_iothub_route", id.ID())
					}
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

	d.SetId(id.ID())

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the completion of the creating/updating of %s: %+v", id, err)
	}

	return resourceIotHubRouteRead(d, meta)
}
