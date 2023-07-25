// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for CDN FrontDoor Route <-> CDN FrontDoor Custom Domain Association creation")

	cdId, err := parse.FrontDoorCustomDomainID(d.Get("cdn_frontdoor_custom_domain_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFrontDoorCustomDomainAssociationID(cdId.SubscriptionId, cdId.ResourceGroup, cdId.ProfileName, cdId.CustomDomainName)

	existing, err := client.Get(ctx, cdId.ResourceGroup, cdId.ProfileName, cdId.CustomDomainName)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("creating %s: %s was not found", id, cdId)
		}

		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// make sure the routes exist and are valid for this custom domain...
	routes, err := validateRoutes(d, meta, cdId)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("cdn_frontdoor_custom_domain_id", cdId.ID())
	d.Set("cdn_frontdoor_route_ids", routes)

	return resourceCdnFrontDoorCustomDomainAssociationRead(d, meta)
}

func resourceCdnFrontDoorCustomDomainAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorCustomDomainAssociationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AssociationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return nil
}

func resourceCdnFrontDoorCustomDomainAssociationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if d.HasChange("cdn_frontdoor_route_ids") {
		cdId, err := parse.FrontDoorCustomDomainID(d.Get("cdn_frontdoor_custom_domain_id").(string))
		if err != nil {
			return err
		}

		id := parse.NewFrontDoorCustomDomainAssociationID(cdId.SubscriptionId, cdId.ResourceGroup, cdId.ProfileName, cdId.CustomDomainName)

		existing, err := client.Get(ctx, cdId.ResourceGroup, cdId.ProfileName, cdId.CustomDomainName)
		if err != nil {
			if utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("updating %s: %s was not found", id, cdId)
			}

			return fmt.Errorf("updating %s: %+v", id, err)
		}

		// make sure the routes exist and are valid for this custom domain...
		routes, err := validateRoutes(d, meta, cdId)
		if err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}

		d.Set("cdn_frontdoor_route_ids", routes)
	}

	return resourceCdnFrontDoorCustomDomainAssociationRead(d, meta)
}

func resourceCdnFrontDoorCustomDomainAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	// since you are deleting the resource you cannot grab the value from the config
	// because it will be empty, you have to get it from the states old value...
	oCdId, _ := d.GetChange("cdn_frontdoor_custom_domain_id")

	cdId, err := parse.FrontDoorCustomDomainID(oCdId.(string))
	if err != nil {
		return err
	}

	id, err := parse.FrontDoorCustomDomainAssociationID(d.Id())
	if err != nil {
		return err
	}

	oRids, _ := d.GetChange("cdn_frontdoor_route_ids")
	oR := oRids.([]interface{})

	v, _, err := expandRoutes(oR)
	if err != nil {
		return err
	}

	if len(*v) != 0 {
		if err := removeCustomDomainAssociationFromRoutes(d, meta, v, cdId); err != nil {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	d.SetId("")

	return nil
}

func validateRoutes(d *pluginsdk.ResourceData, meta interface{}, id *parse.FrontDoorCustomDomainId) ([]interface{}, error) {
	out := make([]interface{}, 0)
	o, n := d.GetChange("cdn_frontdoor_route_ids")
	oRoutes := o.([]interface{})
	nRoutes := n.([]interface{})

	if len(nRoutes) == 0 || nRoutes == nil || id == nil {
		return out, nil
	}

	oIds, _, err := expandRoutes(oRoutes)
	if err != nil {
		return out, err
	}

	nIds, result, err := expandRoutes(nRoutes)
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
