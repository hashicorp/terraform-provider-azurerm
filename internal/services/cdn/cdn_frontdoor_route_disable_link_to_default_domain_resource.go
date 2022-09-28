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
				ForceNew: true,

				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.FrontDoorCustomDomainID,
				},
			},
		},
	}
}

func resourceCdnFrontDoorRouteDisableLinkToDefaultDomainCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	customDomainClient := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	customDomainCtx, customDomainCancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer customDomainCancel()

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

	// make sure the custom domains and the route all belong to the same profile
	for _, v := range customDomains {
		customDomainId, err := parse.FrontDoorCustomDomainID(v.(string))
		if err != nil {
			return err
		}

		if customDomainId.ProfileName != routeId.ProfileName {
			return fmt.Errorf("azurerm_cdn_frontdoor_route_disable_link_to_default_domain: the configuration is invalid, the Front Door Custom Domain(Name: %q, Profile: %q) and the Front Door Route(Name: %q, Profile: %q) must belong to the same Front Door Profile", customDomainId.CustomDomainName, customDomainId.ProfileName, routeId.RouteName, routeId.ProfileName)
		}
	}

	// create the resource id
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("generating UUID for the 'azurerm_cdn_frontdoor_route_disable_link_to_default_domain': %+v", err)
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

	// NOTE: Make sure the route and the custom domains exist since we are attempting
	// to add an association with the custom domain to the route resource
	for _, v := range customDomains {
		customDomainId, err := parse.FrontDoorCustomDomainID(v.(string))
		if err != nil {
			return err
		}

		_, err = customDomainClient.Get(customDomainCtx, customDomainId.ResourceGroup, customDomainId.ProfileName, customDomainId.CustomDomainName)
		if err != nil {
			return fmt.Errorf("azurerm_cdn_frontdoor_route_disable_link_to_default_domain: retrieving existing %s: %+v", *customDomainId, err)
		}
	}

	routeResp, err := routeClient.Get(routeCtx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName)
	if err != nil {
		return fmt.Errorf("azurerm_cdn_frontdoor_route_disable_link_to_default_domain: retrieving existing %s: %+v", *routeId, err)
	}

	props := routeResp.RouteProperties
	if props == nil {
		return fmt.Errorf("azurerm_cdn_frontdoor_route_disable_link_to_default_domain: %s properties are 'nil': %+v", *routeId, err)
	}

	// If it is already disabled skip updating the the resource...
	if props.LinkToDefaultDomain != cdn.LinkToDefaultDomainDisabled {
		customDomainsProps := flattenCdnFrontdoorRouteActivatedResourceArray(props.CustomDomains)
		if len(customDomainsProps) == 0 {
			return fmt.Errorf("azurerm_cdn_frontdoor_route_disable_link_to_default_domain: it is invalid to disable the 'link to default domain' field of the CDN Front Door Route it the route does not have at least one custom domain associated with it, got 0 associated CDN Front Door Custom Domains")
		}

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

	// Everything was successful
	d.SetId(id.ID())

	return resourceCdnFrontDoorRouteDisableLinkToDefaultDomainRead(d, meta)
}

func resourceCdnFrontDoorRouteDisableLinkToDefaultDomainRead(d *pluginsdk.ResourceData, meta interface{}) error {
	routeClient := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	routeCtx, routeCancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer routeCancel()

	routeId, err := parse.FrontDoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
	if err != nil {
		return fmt.Errorf("azurerm_cdn_frontdoor_route_disable_link_to_default_domain: %+v", err)
	}

	routeResp, err := routeClient.Get(routeCtx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName)
	if err != nil {
		return fmt.Errorf("azurerm_cdn_frontdoor_route_disable_link_to_default_domain: retrieving existing %s: %+v", *routeId, err)
	}

	props := routeResp.RouteProperties
	if props == nil {
		return fmt.Errorf("azurerm_cdn_frontdoor_route_disable_link_to_default_domain: %s properties are 'nil': %+v", *routeId, err)
	}

	// NOTE: to keep from throwing the below errors when you attempt to
	// remove the resource from your configuration file I had to implement
	// it with a d.HasChange check...
	if d.HasChange("cdn_frontdoor_route_id") {
		if _, new := d.GetChange("cdn_frontdoor_route_id"); new != "" {
			// Check to make sure the custom domains are associated with the route
			resourceCustomDomains := d.Get("cdn_frontdoor_custom_domain_ids").([]interface{})
			customDomains := flattenCdnFrontdoorRouteActivatedResourceArray(props.CustomDomains)

			// if it is not associated with the route raise an error...
			for _, v := range resourceCustomDomains {
				customDomain, err := parse.FrontDoorCustomDomainID(v.(string))
				if err != nil {
					return fmt.Errorf("azurerm_cdn_frontdoor_route_disable_link_to_default_domain: unable to parse %q: %+v", v.(string), err)
				}

				if !sliceContainsString(customDomains, customDomain.ID()) {
					return fmt.Errorf("azurerm_cdn_frontdoor_route_disable_link_to_default_domain: the custom domain %q is currently unassociated with the defined route %q. Please remove the 'azurerm_cdn_frontdoor_route_disable_link_to_default_domain' resource from your configuration file", customDomain.CustomDomainName, routeId.RouteName)
				}
			}

			d.Set("cdn_frontdoor_route_id", routeId.ID())
			d.Set("cdn_frontdoor_custom_domain_ids", customDomains)
		}
	}

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

		return fmt.Errorf("azurerm_cdn_frontdoor_route_disable_link_to_default_domain: retrieving existing %s: %+v", *routeId, err)
	}

	props := routeResp.RouteProperties
	if props == nil {
		return fmt.Errorf("azurerm_cdn_frontdoor_route_disable_link_to_default_domain: %s properties are 'nil': %+v", *routeId, err)
	}

	updateProps := azuresdkhacks.RouteUpdatePropertiesParameters{
		CustomDomains: props.CustomDomains,
	}

	// NOTE: You must pull the Cache Configuration from the existing route else you will
	// get a diff because the API sees nil as disabled
	if props.CacheConfiguration != nil {
		updateProps.CacheConfiguration = props.CacheConfiguration
	}

	// NOTE: If you are deleting this resource you are reseting the
	// value to the default of true...
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

	// Everything was successful
	d.SetId("")

	return nil
}
