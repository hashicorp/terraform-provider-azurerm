// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_custom_domain_association

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorCustomDomainAssociationUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.AFDCustomDomainsClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if d.HasChange("cdn_frontdoor_route_ids") {
		cdId, err := afddomains.ParseCustomDomainID(d.Get("cdn_frontdoor_custom_domain_id").(string))
		if err != nil {
			return err
		}

		id := parse.NewFrontDoorCustomDomainAssociationID(cdId.SubscriptionId, cdId.ResourceGroupName, cdId.ProfileName, cdId.CustomDomainName)

		customDomainId := afddomains.NewCustomDomainID(cdId.SubscriptionId, cdId.ResourceGroupName, cdId.ProfileName, cdId.CustomDomainName)
		existing, err := client.AFDCustomDomainsGet(ctx, customDomainId)
		if err != nil {
			if response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("updating %s: %s was not found", id, cdId)
			}

			return fmt.Errorf("updating %s: %+v", id, err)
		}

		// make sure the routes exist and are valid for this custom domain...
		routes, err := validateRoutes(ctx, d, meta, cdId)
		if err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}

		d.Set("cdn_frontdoor_route_ids", routes)
	}

	return resourceCdnFrontDoorCustomDomainAssociationRead(d, meta)
}
