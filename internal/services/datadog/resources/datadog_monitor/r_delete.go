// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/datadogmonitorresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatadogMonitorDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Datadog.DatadogMonitorResources
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := datadogmonitorresources.ParseMonitorID(d.Id())
	if err != nil {
		return err
	}

	if err := client.MonitorsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting of %s: %+v", *id, err)
	}

	return nil
}
