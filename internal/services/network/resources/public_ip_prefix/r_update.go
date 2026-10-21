// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package public_ip_prefix

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/publicipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePublicIpPrefixUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.PublicIPPrefixes
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := publicipprefixes.ParsePublicIPPrefixID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("sku") {
		existing, err := client.Get(ctx, *id, publicipprefixes.DefaultGetOperationOptions())
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}
		if existing.Model == nil {
			return fmt.Errorf("retrieving %s: `model` was nil", id)
		}
		if existing.Model.Sku == nil {
			return fmt.Errorf("retrieving %s: `sku` was nil", id)
		}

		existing.Model.Sku.Name = pointer.ToEnum[publicipprefixes.PublicIPPrefixSkuName](d.Get("sku").(string))
		if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	if d.HasChange("tags") {
		payload := publicipprefixes.TagsObject{
			Tags: tags.Expand(d.Get("tags").(map[string]any)),
		}

		if _, err = client.UpdateTags(ctx, *id, payload); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourcePublicIpPrefixRead(d, meta)
}
