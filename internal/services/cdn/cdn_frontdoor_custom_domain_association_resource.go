package cdn

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var cdnFrontDoorCustomDomainResourceName = "azurerm_cdn_frontdoor_custom_domain"
var cdnFrontDoorRouteResourceName = "azurerm_cdn_frontdoor_route"
var notAssociatedErr = "the CDN FrontDoor Route(Name: %q) is currently not associated with the CDN FrontDoor Custom Domain(Name: %q). Please remove the CDN FrontDoor Route from your 'cdn_frontdoor_custom_domain_association' configuration block"

func resourceCdnFrontDoorCustomDomainAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorCustomDomainAssociationCreate,
		Read:   resourceCdnFrontDoorCustomDomainAssociationRead,
		Update: resourceCdnFrontDoorCustomDomainAssociationUpdate,
		Delete: resourceCdnFrontDoorCustomDomainAssociationDelete,

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
			"cdn_frontdoor_custom_domain_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorCustomDomainID,
			},

			"cdn_frontdoor_route_ids": {
				Type:     pluginsdk.TypeList,
				Required: true,

				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.FrontDoorRouteID,
				},
			},
		},
	}
}

func resourceCdnFrontDoorCustomDomainAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	log.Printf("[INFO] preparing arguments for CDN FrontDoor Custom Domain Association")
	customDomain := d.Get("cdn_frontdoor_custom_domain_id").(string)
	configRouteIds := d.Get("cdn_frontdoor_route_ids").([]interface{})

	customDomainId, err := parse.FrontDoorCustomDomainIDInsensitively(customDomain)
	if err != nil {
		return err
	}

	// create the resource id
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("generating UUID: %+v", err)
	}
	id := parse.NewFrontDoorCustomDomainAssociationID(customDomainId.SubscriptionId, customDomainId.ResourceGroup, customDomainId.ProfileName, customDomainId.CustomDomainName, uuid)

	routeIds, normalizedRoutes, err := normalizeRouteIds(configRouteIds)
	if err != nil {
		return err
	}

	for _, routeId := range *routeIds {
		// Make sure the route exists and get the routes custom domain association list...
		routeAssociations, _, err := getRouteProperties(d, meta, &routeId, "cdn_frontdoor_custom_domain_association")
		if err != nil {
			return err
		}

		// Make sure the custom domain is in the routes association list
		if len(routeAssociations) == 0 || !sliceContainsString(routeAssociations, customDomainId.ID()) {
			return fmt.Errorf(notAssociatedErr, routeId.RouteName, customDomainId.CustomDomainName)
		}
	}

	// validate the routes...
	if len(*routeIds) != 0 {
		if err := validateCustomDomainRoutes(routeIds, customDomainId); err != nil {
			return err
		}
	}

	d.SetId(id.ID())
	d.Set("cdn_frontdoor_custom_domain_id", customDomainId.ID())
	d.Set("cdn_frontdoor_route_ids", normalizedRoutes)

	return resourceCdnFrontDoorCustomDomainAssociationRead(d, meta)
}

func resourceCdnFrontDoorCustomDomainAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	customDomain := d.Get("cdn_frontdoor_custom_domain_id").(string)
	configRouteIds := d.Get("cdn_frontdoor_route_ids").([]interface{})

	customDomainId, err := parse.FrontDoorCustomDomainIDInsensitively(customDomain)
	if err != nil {
		return err
	}

	routeIds, normalizedRoutes, err := normalizeRouteIds(configRouteIds)
	if err != nil {
		return err
	}

	for _, route := range *routeIds {
		// Make sure the route exists and get the routes custom domain association list...
		routeAssociations, _, err := getRouteProperties(d, meta, &route, "cdn_frontdoor_custom_domain_association")
		if err != nil {
			return err
		}

		// Make sure the custom domain is in the routes association list
		if len(routeAssociations) == 0 || !sliceContainsString(routeAssociations, customDomainId.ID()) {
			return fmt.Errorf(notAssociatedErr, route.RouteName, customDomainId.CustomDomainName)
		}
	}

	// validate the routes...
	if len(*routeIds) != 0 {
		if err := validateCustomDomainRoutes(routeIds, customDomainId); err != nil {
			return err
		}
	}

	d.Set("cdn_frontdoor_custom_domain_id", customDomainId.ID())
	d.Set("cdn_frontdoor_route_ids", normalizedRoutes)

	return nil
}

func resourceCdnFrontDoorCustomDomainAssociationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	customDomain := d.Get("cdn_frontdoor_custom_domain_id").(string)

	customDomainId, err := parse.FrontDoorCustomDomainIDInsensitively(customDomain)
	if err != nil {
		return err
	}

	if d.HasChange("cdn_frontdoor_route_ids") {
		old, new := d.GetChange("cdn_frontdoor_route_ids")
		oldRoutes := old.([]interface{})
		newRoutes := new.([]interface{})

		oldRouteIds, _, err := normalizeRouteIds(oldRoutes)
		if err != nil {
			return err
		}

		newRouteIds, newNormalizedRoutes, err := normalizeRouteIds(newRoutes)
		if err != nil {
			return err
		}

		// validate the new routes...
		if len(*newRouteIds) != 0 {
			for _, newRoute := range *newRouteIds {
				// Make sure the route exists and get the routes custom domain association list...
				routeAssociations, _, err := getRouteProperties(d, meta, &newRoute, "cdn_frontdoor_custom_domain_association")
				if err != nil {
					return err
				}

				// Make sure the custom domain is in the routes association list
				if len(routeAssociations) == 0 || !sliceContainsString(routeAssociations, customDomainId.ID()) {
					return fmt.Errorf(notAssociatedErr, newRoute.RouteName, customDomainId.CustomDomainName)
				}
			}

			if err := validateCustomDomainRoutes(newRouteIds, customDomainId); err != nil {
				return err
			}
		}

		// now get the delta between the old and the new list, if any custom domains were removed from
		// the list we need to remove the custom domain association from those routes...
		if delta, _ := routeDelta(oldRouteIds, newRouteIds); len(*delta) != 0 {
			if err = removeCustomDomainAssociationFromRoutes(d, meta, delta, customDomainId); err != nil {
				return err
			}

			d.Set("cdn_frontdoor_route_ids", newNormalizedRoutes)
		}
	}

	return nil
}

func resourceCdnFrontDoorCustomDomainAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	customDomain := d.Get("cdn_frontdoor_custom_domain_id").(string)
	configRouteIds := d.Get("cdn_frontdoor_route_ids").([]interface{})

	customDomainId, err := parse.FrontDoorCustomDomainIDInsensitively(customDomain)
	if err != nil {
		return err
	}

	routeIds, _, err := normalizeRouteIds(configRouteIds)
	if err != nil {
		return err
	}

	if len(*routeIds) != 0 {
		if err := removeCustomDomainAssociationFromRoutes(d, meta, routeIds, customDomainId); err != nil {
			return err
		}
	}

	d.SetId("")

	return nil
}
