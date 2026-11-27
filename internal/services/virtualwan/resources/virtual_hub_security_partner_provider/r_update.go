// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_security_partner_provider

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/securitypartnerproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualHubSecurityPartnerProviderUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).SecurityPartnerProviders
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := securitypartnerproviders.ParseSecurityPartnerProviderID(d.Id())
	if err != nil {
		return err
	}

	parameters := securitypartnerproviders.TagsObject{}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if _, err := client.UpdateTags(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceVirtualHubSecurityPartnerProviderRead(d, meta)
}
