// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_custom_domain_association

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorCustomDomainAssociationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.AFDCustomDomainsClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorCustomDomainAssociationID(d.Id())
	if err != nil {
		return err
	}

	customDomainId := afddomains.NewCustomDomainID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.AssociationName)
	resp, err := client.AFDCustomDomainsGet(ctx, customDomainId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return nil
}

func validateRoutes(ctx context.Context, d *pluginsdk.ResourceData, meta any, id *afddomains.CustomDomainId) ([]any, error) {
	out := make([]any, 0)
	o, n := d.GetChange("cdn_frontdoor_route_ids")
	oRoutes := o.([]any)
	nRoutes := n.([]any)

	if len(nRoutes) == 0 || id == nil {
		return out, nil
	}

	oIds, _, err := helpers.ExpandRoutes(oRoutes)
	if err != nil {
		return out, err
	}

	nIds, result, err := helpers.ExpandRoutes(nRoutes)
	if err != nil {
		return out, err
	}

	// validate the new routes...
	if len(*nIds) != 0 {
		for _, v := range *nIds {
			// Make sure the route exists and get the routes custom domain association list...
			associations, _, err := helpers.GetRouteProperties(ctx, meta, &v, "cdn_frontdoor_custom_domain_association")
			if err != nil {
				return out, err
			}

			// Make sure the custom domain is in the routes association list
			if len(associations) == 0 || !helpers.SliceContainsString(associations, id.ID()) {
				return out, fmt.Errorf("the CDN FrontDoor Route(Name: %q) is currently not associated with the CDN FrontDoor Custom Domain(Name: %q). Please remove the CDN FrontDoor Route from your 'cdn_frontdoor_custom_domain_association' configuration block", v.RouteName, id.CustomDomainName)
			}
		}

		if err := helpers.ValidateCustomDomainRoutes(nIds, id); err != nil {
			return out, err
		}
	}

	if len(oRoutes) != 0 {
		// now get the delta between the old and the new list, if any custom domains were removed from
		// the list we need to remove the custom domain association from those routes...
		if delta, _ := helpers.RouteDelta(oIds, nIds); len(*delta) != 0 {
			if err = helpers.RemoveCustomDomainAssociationFromRoutes(ctx, meta, delta, id); err != nil {
				return out, err
			}
		}
	}

	return result, nil
}
