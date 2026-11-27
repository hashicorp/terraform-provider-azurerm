// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_private_link_scoped_service

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2019-10-17-preview/privatelinkscopedresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	monitorClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorPrivateLinkScopedServiceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := monitorClient.From(meta).PrivateLinkScopedResourcesClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := privatelinkscopedresources.ParseScopedResourceID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
