// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_schedule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/schedules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	devtestlabsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDevTestLabSchedulesDelete(d *pluginsdk.ResourceData, meta any) error {
	client := devtestlabsClient.From(meta).LabSchedulesClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := schedules.ParseLabScheduleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return err
	}
	return nil
}
