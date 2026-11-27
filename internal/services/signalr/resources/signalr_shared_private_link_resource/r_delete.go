// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_shared_private_link_resource

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSignalrSharedPrivateLinkDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SignalR.SignalRClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := signalr.ParseSharedPrivateLinkResourceIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	if err := client.SharedPrivateLinkResourcesDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}
