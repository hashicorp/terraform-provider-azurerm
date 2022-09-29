package cdn

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
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

var cdnFrontDoorCustomDomainResourceName = "azurerm_cdn_frontdoor_custom_domain"
var cdnFrontDoorRouteResourceName = "azurerm_cdn_frontdoor_route"

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
		return err
	}

	// create the resource id
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("generating UUID: %+v", err)
	}
	id := parse.NewFrontDoorRouteDisableLinkToDefaultDomainID(routeId.SubscriptionId, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName, uuid)

	// Lock the resources to make sure this will be an atomic operation
	// as this is a one to many association potentially(e.g. route -> custom domain(s))...
	locks.ByName(routeId.RouteName, cdnFrontDoorRouteResourceName)
	defer locks.UnlockByName(routeId.RouteName, cdnFrontDoorRouteResourceName)

	// Lock all of the custom domains associated with this unlink default domain resource
	// prolly overkill but we don't want the custom domains moving while we are doing this
	// operation as well...
	for _, v := range customDomains {
		customDomainId, err := parse.FrontDoorCustomDomainID(v.(string))
		if err != nil {
			return err
		}

		locks.ByName(customDomainId.CustomDomainName, cdnFrontDoorCustomDomainResourceName)
		defer locks.UnlockByName(customDomainId.CustomDomainName, cdnFrontDoorCustomDomainResourceName)
	}

	routeResp, err := routeClient.Get(routeCtx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *routeId, err)
	}

	props := routeResp.RouteProperties
	if props == nil {
		return fmt.Errorf("%s properties are 'nil': %+v", *routeId, err)
	}

	resourceCustomDomains := d.Get("cdn_frontdoor_custom_domain_ids").([]interface{})
	routeCustomDomains := flattenCdnFrontdoorRouteActivatedResourceArray(props.CustomDomains)

	// make sure its valid to disable the LinkToDefaultDomain on this route...
	if len(routeCustomDomains) == 0 {
		return fmt.Errorf("it is invalid to disable the 'LinkToDefaultDomain' for the CDN Front Door Route(Name: %s) since the route does not have any CDN Front Door Custom Domains associated with it", routeId.RouteName)
	}

	// validate the custom domains...
	if err := validateCustomDomanLinkToDefaultDomainState(resourceCustomDomains, routeCustomDomains, routeId.RouteName, routeId.ProfileName); err != nil {
		return err
	}

	// If it is already disabled do not update the route...
	if props.LinkToDefaultDomain != cdn.LinkToDefaultDomainDisabled {
		updateProps := azuresdkhacks.RouteUpdatePropertiesParameters{
			CustomDomains: expandCdnFrontdoorRouteActivatedResourceArray(customDomains),
		}

		// Since this unlink default domain resource always set the value to false
		updateProps.LinkToDefaultDomain = cdn.LinkToDefaultDomainDisabled

		// NOTE: You must pull the Cache Configuration from the existing route else you will get a diff, because a nil value means disabled
		if props.CacheConfiguration != nil {
			updateProps.CacheConfiguration = props.CacheConfiguration
		}

		updatePrarams := azuresdkhacks.RouteUpdateParameters{
			RouteUpdatePropertiesParameters: &updateProps,
		}

		future, err := workaroundsClient.Update(routeCtx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName, updatePrarams)
		if err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		if err = future.WaitForCompletionRef(routeCtx, routeClient.Client); err != nil {
			return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceCdnFrontDoorRouteDisableLinkToDefaultDomainRead(d, meta)
}

func resourceCdnFrontDoorRouteDisableLinkToDefaultDomainRead(d *pluginsdk.ResourceData, meta interface{}) error {
	routeClient := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	routeCtx, routeCancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer routeCancel()

	customDomainClient := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	customDomainCtx, customDomaincancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer customDomaincancel()

	routeId, err := parse.FrontDoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
	if err != nil {
		return fmt.Errorf("unable to parse CDN Front Door Route ID: %+v", err)
	}

	// Make sure the route still exist...
	routeResp, err := routeClient.Get(routeCtx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *routeId, err)
	}

	props := routeResp.RouteProperties
	if props == nil {
		return fmt.Errorf("%s properties are 'nil': %+v", *routeId, err)
	}

	// Make sure all of the custom domains still exist...
	routeCustomDomains := flattenCdnFrontdoorRouteActivatedResourceArray(props.CustomDomains)
	resourceCustomDomains := d.Get("cdn_frontdoor_custom_domain_ids").([]interface{})
	for _, v := range resourceCustomDomains {
		customDomainId, err := parse.FrontDoorCustomDomainID(v.(string))
		if err != nil {
			return fmt.Errorf("unable to parse CDN Front Door Custom Domain ID: %+v", err)
		}

		_, err = customDomainClient.Get(customDomainCtx, customDomainId.ResourceGroup, customDomainId.ProfileName, customDomainId.CustomDomainName)
		if err != nil {
			return fmt.Errorf("retrieving existing %s: %+v", customDomainId, err)
		}
	}

	// Only do the validation if you are not deleting the resource...
	if d.HasChange("cdn_frontdoor_route_id") {
		if _, newRoute := d.GetChange("cdn_frontdoor_route_id"); newRoute != "" {
			if len(routeCustomDomains) == 0 {
				return fmt.Errorf("there are currently no CDN Front Door Custom Domains associated with the CDN Front Door Route(Name: %q). Please remove the resource from your configuration file", routeId.RouteName)
			}

			// validate the custom domains...
			if err := validateCustomDomanLinkToDefaultDomainState(resourceCustomDomains, routeCustomDomains, routeId.RouteName, routeId.ProfileName); err != nil {
				return err
			}

			// In case someone updates this value in portal...
			if props.LinkToDefaultDomain != cdn.LinkToDefaultDomainDisabled {
				return fmt.Errorf("the 'LinkToDefaultDomain' field has been 'enabled' on the CDN Front Door Route(Name: %q). Please revert this value to 'disabled' before proceeding", routeId.RouteName)
			}
		}
	}

	d.Set("cdn_frontdoor_route_id", routeId.ID())
	d.Set("cdn_frontdoor_custom_domain_ids", resourceCustomDomains)

	return nil
}

func resourceCdnFrontDoorRouteDisableLinkToDefaultDomainUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	routeClient := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	routeCtx, routeCancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer routeCancel()

	routeId, err := parse.FrontDoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
	if err != nil {
		return fmt.Errorf("unable to parse CDN Front Door Route ID: %+v", err)
	}

	routeResp, err := routeClient.Get(routeCtx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *routeId, err)
	}

	props := routeResp.RouteProperties
	if props == nil {
		return fmt.Errorf("%s properties are 'nil': %+v", *routeId, err)
	}

	resourceCustomDomains := d.Get("cdn_frontdoor_custom_domain_ids").([]interface{})
	routeCustomDomains := flattenCdnFrontdoorRouteActivatedResourceArray(props.CustomDomains)

	// validate the custom domains...
	if err := validateCustomDomanLinkToDefaultDomainState(resourceCustomDomains, routeCustomDomains, routeId.RouteName, routeId.ProfileName); err != nil {
		return err
	}

	d.Set("cdn_frontdoor_custom_domain_ids", resourceCustomDomains)

	return nil
}

func resourceCdnFrontDoorRouteDisableLinkToDefaultDomainDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	routeClient := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	workaroundsClient := azuresdkhacks.NewCdnFrontDoorRoutesWorkaroundClient(routeClient)
	routeCtx, routeCancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer routeCancel()

	id, err := parse.FrontDoorRouteDisableLinkToDefaultDomainID(d.Id())
	if err != nil {
		return err
	}

	currentRoute := d.Get("cdn_frontdoor_route_id").(string)

	// If this delete was due to a change you need to revert
	// the old route not the new route...
	if d.HasChange("cdn_frontdoor_route_id") {
		if oldRoute, _ := d.GetChange("cdn_frontdoor_route_id"); oldRoute.(string) != "" {
			currentRoute = oldRoute.(string)
		}
	}

	routeId, err := parse.FrontDoorRouteID(currentRoute)
	if err != nil {
		return err
	}

	locks.ByName(routeId.RouteName, cdnFrontDoorRouteResourceName)
	defer locks.UnlockByName(routeId.RouteName, cdnFrontDoorRouteResourceName)

	routeResp, err := routeClient.Get(routeCtx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName)
	if err != nil {
		// No Op: Since the route and the custom domain resources will be deleted before this resource
		// in the destroy scenario a 404 can be ignored as it means that it was already destroyed
		// so you do not need to set the "link to default domain" value on the route anymore...
		if utils.ResponseWasNotFound(routeResp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving existing %s: %+v", *routeId, err)
	}

	props := routeResp.RouteProperties
	if props == nil {
		return fmt.Errorf("%s properties are 'nil': %+v", *routeId, err)
	}

	updateProps := azuresdkhacks.RouteUpdatePropertiesParameters{
		CustomDomains: props.CustomDomains,
	}

	// NOTE: You must pull the Cache Configuration from the existing route else you will
	// get a diff because the API sees nil as disabled
	if props.CacheConfiguration != nil {
		updateProps.CacheConfiguration = props.CacheConfiguration
	}

	// NOTE: Only update LinkToDefaultDomain to enabled if there are not any custom domains associated with the route
	routeCustomDomains := flattenCdnFrontdoorRouteActivatedResourceArray(props.CustomDomains)

	if len(routeCustomDomains) == 0 {
		// only update the route if it is currently in the disabled state...
		if updateProps.LinkToDefaultDomain == cdn.LinkToDefaultDomainDisabled {
			updateProps.LinkToDefaultDomain = cdn.LinkToDefaultDomainEnabled

			updatePrarams := azuresdkhacks.RouteUpdateParameters{
				RouteUpdatePropertiesParameters: &updateProps,
			}

			future, err := workaroundsClient.Update(routeCtx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName, updatePrarams)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			if err = future.WaitForCompletionRef(routeCtx, routeClient.Client); err != nil {
				return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
			}
		}
	}

	// Everything was successful
	d.SetId("")

	return nil
}
