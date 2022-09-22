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
)

func resourceCdnFrontDoorCustomDomainAssociationDefault() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorCustomDomainAssociationDefaultCreate,
		Read:   resourceCdnFrontDoorCustomDomainAssociationDefaultRead,
		Delete: resourceCdnFrontDoorCustomDomainAssociationDefaultDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			if _, err := parse.FrontDoorRouteID(id); err != nil {
				return err
			}
			return nil
		}),

		Schema: map[string]*pluginsdk.Schema{
			"cdn_frontdoor_custom_domain_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorCustomDomainID,
			},

			"cdn_frontdoor_route_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorRouteID,
			},
		},
	}
}

func resourceCdnFrontDoorCustomDomainAssociationDefaultCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	customDomainClient := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	customDomainCtx, customDomainCancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer customDomainCancel()

	routeClient := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	workaroundsClient := azuresdkhacks.NewCdnFrontDoorRoutesWorkaroundClient(routeClient)
	routeCtx, routeCancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer routeCancel()

	log.Printf("[INFO] preparing arguments for CDN FrontDoor Custom Domain <-> CDN FrontDoor Route Association Default creation")

	customDomainId, err := parse.FrontDoorCustomDomainID(d.Get("cdn_frontdoor_custom_domain_id").(string))
	if err != nil {
		return err
	}

	routeId, err := parse.FrontDoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
	if err != nil {
		return err
	}

	// create the association default id
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("generating UUID for the 'azurerm_cdn_frontdoor_custom_domain_association_default': %+v", err)
	}
	assocId := parse.NewFrontDoorCustomDomainAssociationDefaultID(routeId.SubscriptionId, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, uuid, routeId.RouteName, customDomainId.CustomDomainName)

	// Lock both of the resources to make sure this will be an atomic operation
	// as this is a one to many association potentially(e.g. route -> custom domain(s))...
	locks.ByName(routeId.RouteName, cdnFrontDoorRouteResourceName)
	defer locks.UnlockByName(routeId.RouteName, cdnFrontDoorRouteResourceName)

	locks.ByName(customDomainId.CustomDomainName, cdnFrontDoorCustomDomainResourceName)
	defer locks.UnlockByName(customDomainId.CustomDomainName, cdnFrontDoorCustomDomainResourceName)

	// NOTE: Make sure the route and the custom domain both exists since we are attempting
	// to add an association with the custom domain to the route resource
	_, err = customDomainClient.Get(customDomainCtx, customDomainId.ResourceGroup, customDomainId.ProfileName, customDomainId.CustomDomainName)
	if err != nil {
		return fmt.Errorf("azurerm_cdn_frontdoor_custom_domain_association_default: retrieving existing %s: %+v", *customDomainId, err)
	}

	routeResp, err := routeClient.Get(routeCtx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName)
	if err != nil {
		return fmt.Errorf("azurerm_cdn_frontdoor_custom_domain_association_default: retrieving existing %s: %+v", *routeId, err)
	}

	props := routeResp.RouteProperties
	if props == nil {
		return fmt.Errorf("azurerm_cdn_frontdoor_custom_domain_association_default: %s properties are 'nil': %+v", *routeId, err)
	}

	// NOTE: This is needed due to the association resource
	// upon initial creation of the route, the route must
	// have at least one "domain" associated with it by having
	// either a custom domain ID associated with it or by having
	// the routes "link to default domain" field set to true. the
	// workaround for that issue is to expose two different versions
	// of the association resource, one that always has the "link to
	// default domain" set to true and another association resource
	// that has "link to default domain" always set to false.

	customDomains := flattenCdnFrontdoorRouteActivatedResourceArray(props.CustomDomains)
	if len(customDomains) != 0 {
		// Check to make sure the custom domain is not already associated with the route
		if RouteContainsCustomDomain(customDomains, customDomainId.ID()) {
			return fmt.Errorf("azurerm_cdn_frontdoor_custom_domain_association_default: the CDN FrontDoor Custom Domain %q is already associated with %s", customDomainId.ID(), *routeId)
		}
	}

	// it is not a duplicate custom domain for the route, add the custom domain
	// to the route to make the association...
	customDomains = append(customDomains, customDomainId.ID())

	updateProps := azuresdkhacks.RouteUpdatePropertiesParameters{
		CustomDomains: expandCdnFrontdoorRouteActivatedResourceArray(customDomains),
	}

	// Since this is the link to default domain association resource always set the
	// value to true if it is set to false
	if props.LinkToDefaultDomain == cdn.LinkToDefaultDomainDisabled {
		updateProps.LinkToDefaultDomain = cdn.LinkToDefaultDomainEnabled
	}

	// NOTE: You must pull the Cache Configuration from the existing route else you will get a diff
	if props.CacheConfiguration != nil {
		updateProps.CacheConfiguration = props.CacheConfiguration
	}

	updatePrarams := azuresdkhacks.RouteUpdateParameters{
		RouteUpdatePropertiesParameters: &updateProps,
	}

	future, err := workaroundsClient.Update(routeCtx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName, updatePrarams)
	if err != nil {
		return fmt.Errorf("azurerm_cdn_frontdoor_custom_domain_association_default: creating %s: %+v", assocId, err)
	}
	if err = future.WaitForCompletionRef(routeCtx, routeClient.Client); err != nil {
		return fmt.Errorf("azurerm_cdn_frontdoor_custom_domain_association_default: waiting for the creation of %s: %+v", assocId, err)
	}

	// Everything was successful
	d.SetId(assocId.ID())

	return resourceCdnFrontDoorCustomDomainAssociationDefaultRead(d, meta)
}

func resourceCdnFrontDoorCustomDomainAssociationDefaultRead(d *pluginsdk.ResourceData, meta interface{}) error {
	routeClient := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	routeCtx, routeCancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer routeCancel()

	customDomainId, err := parse.FrontdoorCustomDomainID(d.Get("cdn_frontdoor_custom_domain_id").(string))
	if err != nil {
		return err
	}

	routeId, err := parse.FrontDoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
	if err != nil {
		return err
	}

	routeResp, err := routeClient.Get(routeCtx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName)
	if err != nil {
		return fmt.Errorf("azurerm_cdn_frontdoor_custom_domain_association_default: retrieving existing %s: %+v", *routeId, err)
	}

	props := routeResp.RouteProperties
	if props == nil {
		return fmt.Errorf("azurerm_cdn_frontdoor_custom_domain_association_default: %s properties are 'nil': %+v", *routeId, err)
	}

	// Check to make sure the custom domain is associated with the route
	customDomains := flattenCdnFrontdoorRouteActivatedResourceArray(props.CustomDomains)
	if !RouteContainsCustomDomain(customDomains, customDomainId.ID()) {
		return fmt.Errorf("azurerm_cdn_frontdoor_custom_domain_association_default: %s is not associated with expected route %s", *customDomainId, *routeId)
	}

	d.Set("cdn_frontdoor_custom_domain_id", customDomainId.ID())
	d.Set("cdn_frontdoor_route_id", routeId.ID())

	return nil
}

func resourceCdnFrontDoorCustomDomainAssociationDefaultDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	routeClient := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	workaroundsClient := azuresdkhacks.NewCdnFrontDoorRoutesWorkaroundClient(routeClient)
	routeCtx, routeCancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer routeCancel()

	assocId, err := parse.FrontDoorCustomDomainAssociationDefaultID(d.Id())
	if err != nil {
		return err
	}

	customDomainId, err := parse.FrontDoorCustomDomainID(d.Get("cdn_frontdoor_custom_domain_id").(string))
	if err != nil {
		return err
	}

	routeId, err := parse.FrontDoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(routeId.RouteName, cdnFrontDoorRouteResourceName)
	defer locks.UnlockByName(routeId.RouteName, cdnFrontDoorRouteResourceName)

	routeResp, err := routeClient.Get(routeCtx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName)
	if err != nil {
		return fmt.Errorf("azurerm_cdn_frontdoor_custom_domain_association_default: retrieving existing %s: %+v", *routeId, err)
	}

	props := routeResp.RouteProperties
	if props == nil {
		return fmt.Errorf("azurerm_cdn_frontdoor_custom_domain_association_default: %s properties are 'nil': %+v", *routeId, err)
	}

	updateProps := azuresdkhacks.RouteUpdatePropertiesParameters{
		CustomDomains: props.CustomDomains,
	}

	// NOTE: You must pull the Cache Configuration from the existing route else you will get a diff
	if props.CacheConfiguration != nil {
		updateProps.CacheConfiguration = props.CacheConfiguration
	}

	var updateDomains bool
	customDomains := flattenCdnFrontdoorRouteActivatedResourceArray(props.CustomDomains)

	if len(customDomains) != 0 {
		// Check to see if the custom domain is associated with the route
		if RouteContainsCustomDomain(customDomains, customDomainId.ID()) {
			updateDomains = true
		}
	}

	if updateDomains {
		// if we found the custom domain in the list of custom domains from the route
		// remove it from the list and update the route with the new list of custom domains
		removeCustomDomain := d.Get("cdn_frontdoor_custom_domain_id").(string)
		newCustomDomains := RemoveCustomDomain(customDomains, removeCustomDomain)

		updateProps = azuresdkhacks.RouteUpdatePropertiesParameters{
			CustomDomains: expandCdnFrontdoorRouteActivatedResourceArray(newCustomDomains),
		}
	}

	// since this is the custom domain association default always enable the
	// link to default domain if it is set to false
	if props.LinkToDefaultDomain == cdn.LinkToDefaultDomainDisabled {
		updateProps.LinkToDefaultDomain = cdn.LinkToDefaultDomainEnabled
	}

	updatePrarams := azuresdkhacks.RouteUpdateParameters{
		RouteUpdatePropertiesParameters: &updateProps,
	}

	future, err := workaroundsClient.Update(routeCtx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName, updatePrarams)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *assocId, err)
	}
	if err = future.WaitForCompletionRef(routeCtx, routeClient.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *assocId, err)
	}

	// Everything was successful
	d.SetId("")

	return nil
}
