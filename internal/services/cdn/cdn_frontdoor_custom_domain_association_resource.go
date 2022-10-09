package cdn

import (
	"fmt"
	"log"
	"time"

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
			if _, err := parse.FrontDoorCustomDomainAssociationID(id); err != nil {
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
	log.Printf("[INFO] preparing arguments for CDN FrontDoor Route <-> CDN FrontDoor Custom Domain Association creation")

	if customDomain, err := expandCustomDomain(d.Get("cdn_frontdoor_custom_domain_id").(string)); err != nil {
		return err
	} else if customDomain != nil {
		id := parse.NewFrontDoorCustomDomainAssociationID(customDomain.SubscriptionId, customDomain.ResourceGroup, customDomain.ProfileName, customDomain.CustomDomainName)

		// TODO: Get import error to work
		// if !utils.ResponseWasNotFound(existing.Response) {
		// 	return tf.ImportAsExistsError("azurerm_cdn_frontdoor_custom_domain_association", id.ID())
		// }

		d.SetId(id.ID())
	}

	return resourceCdnFrontDoorCustomDomainAssociationRead(d, meta)
}

func resourceCdnFrontDoorCustomDomainAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	if customDomain, err := expandCustomDomain(d.Get("cdn_frontdoor_custom_domain_id").(string)); err != nil {
		return err
	} else if customDomain != nil {
		d.Set("cdn_frontdoor_custom_domain_id", customDomain.ID())
		d.Set("cdn_frontdoor_route_ids", d.Get("cdn_frontdoor_route_ids").([]interface{}))
	}

	return nil
}

func resourceCdnFrontDoorCustomDomainAssociationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	if customDomain, err := expandCustomDomain(d.Get("cdn_frontdoor_custom_domain_id").(string)); err != nil {
		return err
	} else if customDomain != nil {
		_, err := expandRouteIds(d, meta, customDomain)
		if err != nil {
			return err
		}
	}

	return resourceCdnFrontDoorCustomDomainAssociationRead(d, meta)
}

func resourceCdnFrontDoorCustomDomainAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	// since you are deleting the resource you cannot grab the value from the config
	// because it will be empty, you have to get it from the states old value...
	oldCustomDomain, _ := d.GetChange("cdn_frontdoor_custom_domain_id")
	customDomain, _ := expandCustomDomain(oldCustomDomain.(string))

	old, _ := d.GetChange("cdn_frontdoor_route_ids")
	oldRoutes := old.([]interface{})

	routes, _, err := normalizeRouteIds(oldRoutes)
	if err != nil {
		return err
	}

	if len(*routes) != 0 {
		if err := removeCustomDomainAssociationFromRoutes(d, meta, routes, customDomain); err != nil {
			return err
		}
	}

	d.SetId("")

	return nil
}

func expandCustomDomain(input string) (*parse.FrontDoorCustomDomainId, error) {
	if len(input) == 0 || input == "" {
		return nil, nil
	}

	customDomain, err := parse.FrontDoorCustomDomainIDInsensitively(input)
	if err != nil {
		return nil, err
	}

	return customDomain, nil
}

func expandRouteIds(d *pluginsdk.ResourceData, meta interface{}, customDomain *parse.FrontDoorCustomDomainId) ([]interface{}, error) {
	out := make([]interface{}, 0)

	old, new := d.GetChange("cdn_frontdoor_route_ids")
	oldRoutes := old.([]interface{})
	newRoutes := new.([]interface{})

	if len(newRoutes) == 0 || newRoutes == nil || customDomain == nil {
		return out, nil
	}

	oldRouteIds, _, err := normalizeRouteIds(oldRoutes)
	if err != nil {
		return nil, err
	}

	newRouteIds, out, err := normalizeRouteIds(newRoutes)
	if err != nil {
		return nil, err
	}

	// validate the new routes...
	if len(*newRouteIds) != 0 {
		for _, newRoute := range *newRouteIds {
			// Make sure the route exists and get the routes custom domain association list...
			routeAssociations, _, err := getRouteProperties(d, meta, &newRoute, "cdn_frontdoor_custom_domain_association")
			if err != nil {
				return nil, err
			}

			// Make sure the custom domain is in the routes association list
			if len(routeAssociations) == 0 || !sliceContainsString(routeAssociations, customDomain.ID()) {
				return nil, fmt.Errorf(notAssociatedErr, newRoute.RouteName, customDomain.CustomDomainName)
			}
		}

		if err := validateCustomDomainRoutes(newRouteIds, customDomain); err != nil {
			return nil, err
		}
	}

	if len(oldRoutes) != 0 {
		// now get the delta between the old and the new list, if any custom domains were removed from
		// the list we need to remove the custom domain association from those routes...
		if delta, _ := routeDelta(oldRouteIds, newRouteIds); len(*delta) != 0 {
			if err = removeCustomDomainAssociationFromRoutes(d, meta, delta, customDomain); err != nil {
				return nil, err
			}
		}
	}

	return out, nil
}
