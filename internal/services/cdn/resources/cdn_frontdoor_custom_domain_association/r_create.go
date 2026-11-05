// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_custom_domain_association

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorCustomDomainAssociationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.AFDCustomDomainsClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	cdId, err := afddomains.ParseCustomDomainID(d.Get("cdn_frontdoor_custom_domain_id").(string))
	if err != nil {
		return err
	}

	// TODO: Composite ID?
	id := parse.NewFrontDoorCustomDomainAssociationID(cdId.SubscriptionId, cdId.ResourceGroupName, cdId.ProfileName, cdId.CustomDomainName)

	if _, err := client.AFDCustomDomainsGet(ctx, *cdId); err != nil {
		return fmt.Errorf("retrieving %s: %+v", cdId, err)
	}

	// make sure the routes exist and are valid for this custom domain...
	r, err := validateRoutes(ctx, d, meta, cdId)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("cdn_frontdoor_custom_domain_id", cdId.ID())
	d.Set("cdn_frontdoor_route_ids", r)

	return resourceCdnFrontDoorCustomDomainAssociationRead(d, meta)
}
