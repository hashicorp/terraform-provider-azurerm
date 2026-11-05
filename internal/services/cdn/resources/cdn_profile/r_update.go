// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_profile

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" //nolint:staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnProfileUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.ProfilesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if !d.HasChange("tags") {
		return nil
	}

	id, err := profiles.ParseProfileID(d.Id())
	if err != nil {
		return err
	}

	newTags := d.Get("tags").(map[string]any)

	props := cdn.ProfileUpdateParameters{
		Tags: tags.Expand(newTags),
	}

	future, err := client.Update(ctx, id.ResourceGroupName, id.ProfileName, props)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	return resourceCdnProfileRead(d, meta)
}
