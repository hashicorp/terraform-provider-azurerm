// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_aad_diagnostic_setting

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/azureactivedirectory/2017-04-01/diagnosticsettings"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorAADDiagnosticSettingDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.AADDiagnosticSettingsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := diagnosticsettings.ParseDiagnosticSettingID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	waitForAADDiagnosticSettingToBeGone := waitForAADDiagnosticSettingToBeGonePoller{
		client: client,
		id:     *id,
	}
	initialDelayDuration := 15 * time.Second
	poller := pollers.NewPoller(&waitForAADDiagnosticSettingToBeGone, initialDelayDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be fully deleted: %+v", *id, err)
	}

	return nil
}
