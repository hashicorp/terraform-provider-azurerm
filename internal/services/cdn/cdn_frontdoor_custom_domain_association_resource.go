package cdn

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
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

func resourceCdnFrontDoorCustomDomainAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorCustomDomainAssociationCreate,
		Read:   resourceCdnFrontDoorCustomDomainAssociationRead,
		Delete: resourceCdnFrontDoorCustomDomainAssociationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
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

func resourceCdnFrontDoorCustomDomainAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	customDomainClient := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	customDomainCtx, customDomainCancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer customDomainCancel()

	routeClient := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	routeCtx, routeCancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer routeCancel()

	log.Printf("[INFO] preparing arguments for CDN FrontDoor Custom Domain <-> CDN FrontDoor Route Association creation")

	customDomainId, err := parse.FrontDoorCustomDomainID(d.Get("cdn_frontdoor_custom_domain_id").(string))
	if err != nil {
		return err
	}

	routeId, err := parse.FrontDoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
	if err != nil {
		return err
	}

	// Lock both of the resources to make sure this will be an atomic operation
	// as this is a one to many association potentially(e.g. route -> custom domain(s))...
	locks.ByName(routeId.RouteName, cdnFrontDoorRouteResourceName)
	defer locks.UnlockByName(routeId.RouteName, cdnFrontDoorRouteResourceName)

	locks.ByName(customDomainId.CustomDomainName, cdnFrontDoorCustomDomainResourceName)
	defer locks.UnlockByName(customDomainId.CustomDomainName, cdnFrontDoorCustomDomainResourceName)

	// NOTE: Make sure the route and the custom domain both exists since we are attempting
	// to add/remove an association with the custom domain to the route resource
	_, err = customDomainClient.Get(customDomainCtx, customDomainId.ResourceGroup, customDomainId.ProfileName, customDomainId.CustomDomainName)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *customDomainId, err)
	}

	routeResp, err := routeClient.Get(routeCtx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *routeId, err)
	}

	if props := routeResp.RouteProperties; props == nil {
		return fmt.Errorf("cdn frontdoor route %s properties are 'nil': %+v", *routeId, err)
	} else {
		customDomains := flattenCdnFrontdoorRouteActivatedResourceArray(props.CustomDomains)
		linkToDefaultDomain := props.LinkToDefaultDomain

		// It's empty just add it, but if it is empty we need to make sure that linkToDefaultDomain is set to true
		// since the route requires at least one custom domain to be associated with it...
		if len(customDomains) == 0 {
			// This seems weird, need to work this out... since the route and the custom domain will be created first, may need to move link to default domain into the association resources as well... :/
			if linkToDefaultDomain == cdn.LinkToDefaultDomainDisabled {
				return fmt.Errorf("at least one domain is required for the 'azurerm_cdn_frontdoor_route'. Please create a 'cdn_frontdoor_custom_domain_association' resource or set the 'link_to_default_domain_enabled' to 'true'")
			}
		} else {
			// Check to make sure the custom domain is not already associated with the route
			err := isCustomDomainAlreadyAssociatedWithRoute(customDomains, customDomainId.ID())
			if err != nil {
				return fmt.Errorf("%s %q", err, *routeId)
			}

			// it is not a duplicate custom domain for the route, add the association...
		}

		enabledState := props.EnabledState
		forwardingProtocol := props.ForwardingProtocol
		httpsRedirect := props.HTTPSRedirect
		originPath := props.OriginPath
		patternsToMatch := props.PatternsToMatch

		// these might be nil
		cacheConfiguration := props.CacheConfiguration
		originGroup := props.OriginGroup
		ruleSets := props.RuleSets
		supportedProtocols := props.SupportedProtocols

		routeProps := azuresdkhacks.RouteUpdatePropertiesParameters{
			CustomDomains:       expandCdnFrontdoorRouteActivatedResourceArray(customDomains),
			CacheConfiguration:  cacheConfiguration,
			EnabledState:        enabledState,
			ForwardingProtocol:  forwardingProtocol,
			HTTPSRedirect:       httpsRedirect,
			LinkToDefaultDomain: linkToDefaultDomain,
			OriginGroup:         originGroup,
			OriginPath:          originPath,
			PatternsToMatch:     patternsToMatch,
			RuleSets:            ruleSets,
			SupportedProtocols:  supportedProtocols,
		}

		err := associationCdnFrontDoorRouteUpdate(d, meta, routeProps)
		if err != nil {
			return err
		}
	}

	// Everything was successful, create the association resource ID...
	assocId := parse.NewFrontDoorCustomDomainAssociationID(routeId.SubscriptionId, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName, customDomainId.CustomDomainName)

	d.SetId(assocId.ID())

	return resourceCdnFrontDoorCustomDomainAssociationRead(d, meta)
}

func resourceCdnFrontDoorCustomDomainAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	return nil
}

func resourceCdnFrontDoorCustomDomainAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// here I will need to remove the custom domain id from the list, and update the route

	return nil
}

func associationCdnFrontDoorRouteUpdate(d *pluginsdk.ResourceData, meta interface{}, props azuresdkhacks.RouteUpdatePropertiesParameters) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	workaroundsClient := azuresdkhacks.NewCdnFrontDoorRoutesWorkaroundClient(client)
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorRouteID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("cdn_frontdoor_custom_domain_ids") {
		props.CustomDomains = expandCdnFrontdoorRouteActivatedResourceArray(customDomains)
	}

	if d.HasChange("cache") {
		props.CacheConfiguration = expandCdnFrontdoorRouteCacheConfiguration(d.Get("cache").([]interface{}))
	}

	if d.HasChange("enabled") {
		props.EnabledState = expandEnabledBool(d.Get("enabled").(bool))
	}

	if d.HasChange("forwarding_protocol") {
		props.ForwardingProtocol = cdn.ForwardingProtocol(d.Get("forwarding_protocol").(string))
	}

	if d.HasChange("https_redirect_enabled") {
		props.HTTPSRedirect = expandEnabledBoolToRouteHttpsRedirect(httpsRedirect)
	}

	if d.HasChange("link_to_default_domain_enabled") {
		props.LinkToDefaultDomain = expandEnabledBoolToLinkToDefaultDomain(linkToDefaultDomain)
	}

	if d.HasChange("cdn_frontdoor_origin_group_id") {
		props.OriginGroup = expandResourceReference(d.Get("cdn_frontdoor_origin_group_id").(string))
	}

	if d.HasChange("cdn_frontdoor_origin_path") {
		props.OriginPath = utils.String(d.Get("cdn_frontdoor_origin_path").(string))
	}

	if d.HasChange("patterns_to_match") {
		props.PatternsToMatch = utils.ExpandStringSlice(d.Get("patterns_to_match").([]interface{}))
	}

	if d.HasChange("cdn_frontdoor_rule_set_ids") {
		props.RuleSets = expandCdnFrontdoorRouteResourceReferenceArray(d.Get("cdn_frontdoor_rule_set_ids").([]interface{}))
	}

	if d.HasChange("supported_protocols") {
		props.SupportedProtocols = expandCdnFrontdoorRouteEndpointProtocolsArray(protocolsRaw)
	}

	payload := azuresdkhacks.RouteUpdateParameters{
		RouteUpdatePropertiesParameters: &props,
	}

	future, err := workaroundsClient.Update(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName, payload)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	return resourceCdnFrontDoorRouteRead(d, meta)
}

func isCustomDomainAlreadyAssociatedWithRoute(input []interface{}, customDomain string) error {
	for _, key := range input {
		v := key.(string)
		if strings.EqualFold(v, customDomain) {
			return fmt.Errorf("the %q CDN FrontDoor Custom Domain is already associated with the CDN FrontDoor Route", v)
		}
	}

	return nil
}
