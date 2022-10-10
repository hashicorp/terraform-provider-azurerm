package cdn

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
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
		Update: resourceCdnFrontDoorCustomDomainAssociationUpdate,
		Delete: resourceCdnFrontDoorCustomDomainAssociationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.FrontDoorCustomDomainAssociationID(id)
			return err
		}, importCdnFrontDoorCustomDomainAssociation()),

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

	cdId, err := customDomainNullable(d.Get("cdn_frontdoor_custom_domain_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFrontDoorCustomDomainAssociationID(cdId.SubscriptionId, cdId.ResourceGroup, cdId.ProfileName, cdId.CustomDomainName)

	d.SetId(id.ID())

	return resourceCdnFrontDoorCustomDomainAssociationRead(d, meta)
}

func resourceCdnFrontDoorCustomDomainAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := customDomainNullable(d.Get("cdn_frontdoor_custom_domain_id").(string))
	if err != nil {
		return err
	}

	// id will be nill if you are deleting the resource
	if id != nil {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
		if err != nil {
			if utils.ResponseWasNotFound(existing.Response) {
				d.SetId("")
				return fmt.Errorf("CDN FrontDoor Custom Domain(Resource Group: %q Name: %q) was not found", id.ResourceGroup, id.CustomDomainName)
			}

			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}

		// make sure the routes exist and are valid for this custom domain...
		routes, err := flattenRoutes(d, meta, id)
		if err != nil {
			return err
		}

		d.Set("cdn_frontdoor_custom_domain_id", id.ID())
		d.Set("cdn_frontdoor_route_ids", routes)
	}

	return nil
}

func resourceCdnFrontDoorCustomDomainAssociationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	// if this value has not changed there is nothing to do so just return
	if !d.HasChange("cdn_frontdoor_route_ids") {
		return nil
	}

	return resourceCdnFrontDoorCustomDomainAssociationRead(d, meta)
}

func resourceCdnFrontDoorCustomDomainAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	// since you are deleting the resource you cannot grab the value from the config
	// because it will be empty, you have to get it from the states old value...
	oCdId, _ := d.GetChange("cdn_frontdoor_custom_domain_id")
	id, _ := customDomainNullable(oCdId.(string))

	oRids, _ := d.GetChange("cdn_frontdoor_route_ids")
	oR := oRids.([]interface{})

	v, _, err := routesInsensitively(oR)
	if err != nil {
		return err
	}

	if len(*v) != 0 {
		if err := removeCustomDomainAssociationFromRoutes(d, meta, v, id); err != nil {
			return err
		}
	}

	d.SetId("")

	return nil
}

func customDomainNullable(input string) (*parse.FrontDoorCustomDomainId, error) {
	if len(input) == 0 || input == "" {
		return nil, nil
	}

	v, err := parse.FrontDoorCustomDomainIDInsensitively(input)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func flattenRoutes(d *pluginsdk.ResourceData, meta interface{}, id *parse.FrontDoorCustomDomainId) ([]interface{}, error) {
	out := make([]interface{}, 0)
	o, n := d.GetChange("cdn_frontdoor_route_ids")
	oRoutes := o.([]interface{})
	nRoutes := n.([]interface{})

	if len(nRoutes) == 0 || nRoutes == nil || id == nil {
		return out, nil
	}

	oIds, _, err := routesInsensitively(oRoutes)
	if err != nil {
		return out, err
	}

	nIds, result, err := routesInsensitively(nRoutes)
	if err != nil {
		return out, err
	}

	// validate the new routes...
	if len(*nIds) != 0 {
		for _, v := range *nIds {
			// Make sure the route exists and get the routes custom domain association list...
			associations, _, err := getRouteProperties(d, meta, &v, "cdn_frontdoor_custom_domain_association")
			if err != nil {
				return out, err
			}

			// Make sure the custom domain is in the routes association list
			if len(associations) == 0 || !sliceContainsString(associations, id.ID()) {
				return out, fmt.Errorf("the CDN FrontDoor Route(Name: %q) is currently not associated with the CDN FrontDoor Custom Domain(Name: %q). Please remove the CDN FrontDoor Route from your 'cdn_frontdoor_custom_domain_association' configuration block", v.RouteName, id.CustomDomainName)
			}
		}

		if err := validateCustomDomainRoutes(nIds, id); err != nil {
			return out, err
		}
	}

	if len(oRoutes) != 0 {
		// now get the delta between the old and the new list, if any custom domains were removed from
		// the list we need to remove the custom domain association from those routes...
		if delta, _ := routeDelta(oIds, nIds); len(*delta) != 0 {
			if err = removeCustomDomainAssociationFromRoutes(d, meta, delta, id); err != nil {
				return out, err
			}
		}
	}

	return result, nil
}
