// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package private_link_service

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privatelinkservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePrivateLinkServiceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.PrivateLinkServices
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatelinkservices.ParsePrivateLinkServiceID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
