// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dedicated_host_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/dedicatedhostgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	computeClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDedicatedHostGroupUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := computeClient.From(meta).DedicatedHostGroupsClient
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := commonids.ParseDedicatedHostGroupID(d.Id())
	if err != nil {
		return err
	}

	payload := dedicatedhostgroups.DedicatedHostGroupUpdate{
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if _, err := client.Update(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceDedicatedHostGroupRead(d, meta)
}
