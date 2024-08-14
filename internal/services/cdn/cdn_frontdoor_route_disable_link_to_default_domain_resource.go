// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontDoorRouteDisableLinkToDefaultDomain() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorRouteDisableLinkToDefaultDomainCreate,
		Read:   resourceCdnFrontDoorRouteDisableLinkToDefaultDomainRead,
		Update: resourceCdnFrontDoorRouteDisableLinkToDefaultDomainUpdate,
		Delete: resourceCdnFrontDoorRouteDisableLinkToDefaultDomainDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			if _, err := parse.FrontDoorRouteDisableLinkToDefaultDomainID(id); err != nil {
				return err
			}
			return nil
		}),

		DeprecationMessage: "The `azurerm_cdn_frontdoor_route_disable_link_to_default_domain` resource has been deprecated and will be removed in v4.0 of the AzureRM provider in favour of the 'link_to_default_domain' property in the `azurerm_cdn_frontdoor_route` resource.",

		Schema: map[string]*pluginsdk.Schema{
			"cdn_frontdoor_route_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorRouteID,
			},

			"cdn_frontdoor_custom_domain_ids": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.FrontDoorCustomDomainID,
				},
			},
		},
	}
}

func resourceCdnFrontDoorRouteDisableLinkToDefaultDomainCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	routeClient := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	workaroundsClient := azuresdkhacks.NewCdnFrontDoorRoutesWorkaroundClient(routeClient)
	routeCtx, routeCancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer routeCancel()

	log.Printf("[INFO] preparing arguments for CDN FrontDoor Route Unlink Default Domain")

	customDomains := d.Get("cdn_frontdoor_custom_domain_ids").([]interface{})

	routeId, err := parse.FrontDoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
	if err != nil {
		return fmt.Errorf("creating Front Door Route Disable Link To Default Domain: %+v", err)
	}

	// create the resource id
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("generating UUID: %+v", err)
	}

	id := parse.NewFrontDoorRouteDisableLinkToDefaultDomainID(routeId.SubscriptionId, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName, uuid)

	locks.ByName(routeId.RouteName, cdnFrontDoorRouteResourceName)
	defer locks.UnlockByName(routeId.RouteName, cdnFrontDoorRouteResourceName)

	for _, v := range customDomains {
		customDomainId, err := parse.FrontDoorCustomDomainID(v.(string))
		if err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		locks.ByName(customDomainId.CustomDomainName, cdnFrontDoorCustomDomainResourceName)
		defer locks.UnlockByName(customDomainId.CustomDomainName, cdnFrontDoorCustomDomainResourceName)
	}

	existing, err := routeClient.Get(routeCtx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("creating %s: %s was not found", id, routeId)
		}

		return fmt.Errorf("retrieving existing %s: %+v", *routeId, err)
	}

	props := existing.RouteProperties
	if props == nil {
		return fmt.Errorf("creating %s: %s properties are 'nil'", id, *routeId)
	}

	resourceCustomDomains := d.Get("cdn_frontdoor_custom_domain_ids").([]interface{})
	routeCustomDomains, err := flattenCustomDomainActivatedResourceArray(props.CustomDomains)
	if err != nil {
		return err
	}

	// make sure its valid to disable the LinkToDefaultDomain on this route...
	if len(routeCustomDomains) == 0 {
		return fmt.Errorf("creating %s: it is invalid to disable the 'LinkToDefaultDomain' for the CDN Front Door Route(Name: %s) since the route does not have any CDN Front Door Custom Domains associated with it", id, routeId.RouteName)
	}

	// validate the custom domains...
	if err := validateCustomDomainLinkToDefaultDomainState(resourceCustomDomains, routeCustomDomains, routeId.RouteName, routeId.ProfileName); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// If it is already disabled do not update the route...
	if props.LinkToDefaultDomain != cdn.LinkToDefaultDomainDisabled {
		updateProps := azuresdkhacks.RouteUpdatePropertiesParameters{
			CustomDomains: expandCustomDomainActivatedResourceArray(customDomains),
		}

		// Since this unlink default domain resource always set the value to false
		updateProps.LinkToDefaultDomain = cdn.LinkToDefaultDomainDisabled

		// NOTE: You must pull the Cache Configuration from the existing route else you will get a diff, because a nil value means disabled
		if props.CacheConfiguration != nil {
			updateProps.CacheConfiguration = props.CacheConfiguration
		}

		updateParams := azuresdkhacks.RouteUpdateParameters{
			RouteUpdatePropertiesParameters: &updateProps,
		}

		future, err := workaroundsClient.Update(routeCtx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName, updateParams)
		if err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		if err = future.WaitForCompletionRef(routeCtx, routeClient.Client); err != nil {
			return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())
	d.Set("cdn_frontdoor_route_id", routeId.ID())
	d.Set("cdn_frontdoor_custom_domain_ids", customDomains)

	return resourceCdnFrontDoorRouteDisableLinkToDefaultDomainRead(d, meta)
}

func resourceCdnFrontDoorRouteDisableLinkToDefaultDomainRead(d *pluginsdk.ResourceData, meta interface{}) error {
	routeClient := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	routeCtx, routeCancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer routeCancel()

	customDomainClient := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	customDomainCtx, customDomainCancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer customDomainCancel()

	id, err := parse.FrontDoorRouteDisableLinkToDefaultDomainID(d.Id())
	if err != nil {
		return err
	}

	routeId, err := parse.FrontDoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
	if err != nil {
		return fmt.Errorf("front door route disable link to default domain: %+v", err)
	}

	// Make sure the route still exist...
	existing, err := routeClient.Get(routeCtx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving existing %s: %+v", *routeId, err)
	}

	customDomains := d.Get("cdn_frontdoor_custom_domain_ids").([]interface{})
	for _, v := range customDomains {
		cdId, err := parse.FrontDoorCustomDomainID(v.(string))
		if err != nil {
			return fmt.Errorf("%s: unable to parse CDN Front Door Custom Domain ID: %+v", id, err)
		}

		_, err = customDomainClient.Get(customDomainCtx, cdId.ResourceGroup, cdId.ProfileName, cdId.CustomDomainName)
		if err != nil {
			return fmt.Errorf("retrieving existing %s: %+v", cdId, err)
		}
	}

	return nil
}

func resourceCdnFrontDoorRouteDisableLinkToDefaultDomainUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	routeClient := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	workaroundsClient := azuresdkhacks.NewCdnFrontDoorRoutesWorkaroundClient(routeClient)
	routeCtx, routeCancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer routeCancel()

	if d.HasChange("cdn_frontdoor_custom_domain_ids") {
		customDomains := d.Get("cdn_frontdoor_custom_domain_ids").([]interface{})

		routeId, err := parse.FrontDoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
		if err != nil {
			return fmt.Errorf("updating Front Door Route Disable Link To Default Domain: %+v", err)
		}

		id, err := parse.FrontDoorRouteDisableLinkToDefaultDomainID(d.Id())
		if err != nil {
			return err
		}

		locks.ByName(routeId.RouteName, cdnFrontDoorRouteResourceName)
		defer locks.UnlockByName(routeId.RouteName, cdnFrontDoorRouteResourceName)

		for _, v := range customDomains {
			customDomainId, err := parse.FrontDoorCustomDomainID(v.(string))
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			locks.ByName(customDomainId.CustomDomainName, cdnFrontDoorCustomDomainResourceName)
			defer locks.UnlockByName(customDomainId.CustomDomainName, cdnFrontDoorCustomDomainResourceName)
		}

		existing, err := routeClient.Get(routeCtx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName)
		if err != nil {
			if utils.ResponseWasNotFound(existing.Response) {
				d.SetId("")
				return nil
			}

			return fmt.Errorf("%s: retrieving existing %s: %+v", *id, *routeId, err)
		}

		props := existing.RouteProperties
		if props == nil {
			return fmt.Errorf("updating %s: %s properties are 'nil'", id, *routeId)
		}

		resourceCustomDomains := d.Get("cdn_frontdoor_custom_domain_ids").([]interface{})
		routeCustomDomains, err := flattenCustomDomainActivatedResourceArray(props.CustomDomains)
		if err != nil {
			return err
		}

		// make sure its valid to disable the LinkToDefaultDomain on this route...
		if len(routeCustomDomains) == 0 {
			return fmt.Errorf("updating %s: it is invalid to disable the 'LinkToDefaultDomain' for the CDN Front Door Route(Name: %s) since the route does not have any CDN Front Door Custom Domains associated with it", id, routeId.RouteName)
		}

		// validate the custom domains...
		if err := validateCustomDomainLinkToDefaultDomainState(resourceCustomDomains, routeCustomDomains, routeId.RouteName, routeId.ProfileName); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}

		// If it is already disabled do not update the route...
		if props.LinkToDefaultDomain != cdn.LinkToDefaultDomainDisabled {
			updateProps := azuresdkhacks.RouteUpdatePropertiesParameters{
				CustomDomains: expandCustomDomainActivatedResourceArray(customDomains),
			}

			// Since this unlink default domain resource always set the value to false
			updateProps.LinkToDefaultDomain = cdn.LinkToDefaultDomainDisabled

			// NOTE: You must pull the Cache Configuration from the existing route else you will get a diff, because a nil value means disabled
			if props.CacheConfiguration != nil {
				updateProps.CacheConfiguration = props.CacheConfiguration
			}

			updateParams := azuresdkhacks.RouteUpdateParameters{
				RouteUpdatePropertiesParameters: &updateProps,
			}

			future, err := workaroundsClient.Update(routeCtx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName, updateParams)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			if err = future.WaitForCompletionRef(routeCtx, routeClient.Client); err != nil {
				return fmt.Errorf("waiting for the update of %s: %+v", id, err)
			}
		}

		d.Set("cdn_frontdoor_route_id", routeId.ID())
		d.Set("cdn_frontdoor_custom_domain_ids", customDomains)
	}

	return resourceCdnFrontDoorRouteDisableLinkToDefaultDomainRead(d, meta)
}

func resourceCdnFrontDoorRouteDisableLinkToDefaultDomainDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	workaroundsClient := azuresdkhacks.NewCdnFrontDoorRoutesWorkaroundClient(client)
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorRouteDisableLinkToDefaultDomainID(d.Id())
	if err != nil {
		return err
	}

	oldRoute, _ := d.GetChange("cdn_frontdoor_route_id")

	route, err := parse.FrontDoorRouteID(oldRoute.(string))
	if err != nil {
		return err
	}

	locks.ByName(route.RouteName, cdnFrontDoorRouteResourceName)
	defer locks.UnlockByName(route.RouteName, cdnFrontDoorRouteResourceName)

	resp, err := client.Get(ctx, route.ResourceGroup, route.ProfileName, route.AfdEndpointName, route.RouteName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving existing %s: %+v", *route, err)
	}

	props := resp.RouteProperties
	if props == nil {
		return fmt.Errorf("deleting %s: %s properties are 'nil'", *id, *route)
	}

	updateProps := azuresdkhacks.RouteUpdatePropertiesParameters{
		CustomDomains: props.CustomDomains,
	}

	// NOTE: You must pull the Cache Configuration from the existing route else you will
	// get a diff because the API sees nil as disabled
	if props.CacheConfiguration != nil {
		updateProps.CacheConfiguration = props.CacheConfiguration
	}

	customDomains, err := flattenCustomDomainActivatedResourceArray(props.CustomDomains)
	if err != nil {
		return err
	}

	// NOTE: Only update LinkToDefaultDomain to enabled if there are not any custom domains associated with the route
	if len(customDomains) == 0 {
		// only update the route if it is currently in the disabled state...
		if updateProps.LinkToDefaultDomain == cdn.LinkToDefaultDomainDisabled {
			updateProps.LinkToDefaultDomain = cdn.LinkToDefaultDomainEnabled

			updateParams := azuresdkhacks.RouteUpdateParameters{
				RouteUpdatePropertiesParameters: &updateProps,
			}

			future, err := workaroundsClient.Update(ctx, route.ResourceGroup, route.ProfileName, route.AfdEndpointName, route.RouteName, updateParams)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
			}
		}
	}

	// Everything was successful
	d.SetId("")

	return nil
}
