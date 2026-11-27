// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_private_link_hub

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapsePrivateLinkHubUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Synapse.PrivateLinkHubsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PrivateLinkHubID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("tags") {
		privateLinkHubPatchInfo := synapse.PrivateLinkHubPatchInfo{
			Tags: tags.Expand(d.Get("tags").(map[string]any)),
		}

		if _, err := client.Update(ctx, privateLinkHubPatchInfo, id.ResourceGroup, id.Name); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceSynapsePrivateLinkHubRead(d, meta)
}
