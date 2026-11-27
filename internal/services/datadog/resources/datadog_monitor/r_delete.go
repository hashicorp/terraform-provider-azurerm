// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/datadogmonitorresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	datadogClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatadogMonitorDelete(d *pluginsdk.ResourceData, meta any) error {
	client := datadogClient.From(meta).DatadogMonitorResources
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
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
