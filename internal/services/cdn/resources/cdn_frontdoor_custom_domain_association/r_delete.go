// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_custom_domain_association

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorCustomDomainAssociationDelete(d *pluginsdk.ResourceData, meta any) error {
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// since you are deleting the resource you cannot grab the value from the config
	// because it will be empty, you have to get it from the states old value...
	oCdId, _ := d.GetChange("cdn_frontdoor_custom_domain_id")

	cdId, err := afddomains.ParseCustomDomainID(oCdId.(string))
	if err != nil {
		return err
	}

	id, err := parse.FrontDoorCustomDomainAssociationID(d.Id())
	if err != nil {
		return err
	}

	oRids, _ := d.GetChange("cdn_frontdoor_route_ids")
	oR := oRids.([]any)

	v, _, err := helpers.ExpandRoutes(oR)
	if err != nil {
		return err
	}

	if len(*v) != 0 {
		if err := helpers.RemoveCustomDomainAssociationFromRoutes(ctx, meta, v, cdId); err != nil {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	d.SetId("")

	return nil
}
