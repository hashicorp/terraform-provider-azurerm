// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_watcher_flow_log

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/flowlogs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkWatcherFlowLogDelete(d *pluginsdk.ResourceData, meta any) error {
	client := networkClient.From(meta).FlowLogs
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := flowlogs.ParseFlowLogID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.TargetResourceId == "" {
		return fmt.Errorf("retrieving %s: `properties` or `properties.TargetResourceID` was nil", id)
	}

	targetResourceId := resp.Model.Properties.TargetResourceId
	if nsgId, err := networksecuritygroups.ParseNetworkSecurityGroupIDInsensitively(resp.Model.Properties.TargetResourceId); err == nil {
		targetResourceId = nsgId.ID()
	} else if vnetId, err := commonids.ParseVirtualNetworkIDInsensitively(resp.Model.Properties.TargetResourceId); err == nil {
		targetResourceId = vnetId.ID()
	}

	locks.ByID(targetResourceId)
	defer locks.UnlockByID(targetResourceId)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %v", id, err)
	}

	return nil
}
