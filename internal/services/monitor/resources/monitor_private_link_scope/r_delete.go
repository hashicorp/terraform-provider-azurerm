// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_private_link_scope

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-07-01-preview/privatelinkscopesapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorPrivateLinkScopeDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.PrivateLinkScopesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatelinkscopesapis.ParsePrivateLinkScopeID(d.Id())
	if err != nil {
		return err
	}

	if err = client.PrivateLinkScopesDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
