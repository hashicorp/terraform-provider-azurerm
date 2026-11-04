// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resourcegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceResourceGroupUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Resource.ResourceGroupsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseResourceGroupIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	patch := resourcegroups.ResourceGroupPatchable{}

	if d.HasChange("tags") {
		patch.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if _, err := client.Update(ctx, *id, patch); err != nil {
		return fmt.Errorf("updating %q: %+v", *id, err)
	}

	return resourceResourceGroupRead(d, meta)
}
