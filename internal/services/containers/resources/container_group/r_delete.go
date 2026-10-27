// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerinstance/2025-09-01/containerinstance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceContainerGroupDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Containers.ContainerInstanceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := containerinstance.ParseContainerGroupID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.ContainerGroupsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			// already deleted
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	if model := existing.Model; model != nil {
		props := model.Properties
		if subnetIDs := props.SubnetIds; subnetIDs != nil && len(*subnetIDs) != 0 {
			// Avoid parallel deletion if "subnet_ids" are given.
			for _, item := range *subnetIDs {
				subnet, err := commonids.ParseSubnetID(item.Id)
				if err != nil {
					return fmt.Errorf(`parsing subnet id %q: %v`, item.Id, err)
				}

				locks.ByID(subnet.ID())
				defer locks.UnlockByID(subnet.ID())
			}
		}
	}

	if err := client.ContainerGroupsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	return nil
}
