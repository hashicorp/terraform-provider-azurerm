// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package batch_application

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/application"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	batchClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBatchApplicationUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := batchClient.From(meta).ApplicationClient
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := application.ParseApplicationID(d.Id())
	if err != nil {
		return err
	}

	allowUpdates := d.Get("allow_updates").(bool)
	defaultVersion := d.Get("default_version").(string)
	displayName := d.Get("display_name").(string)

	parameters := application.Application{
		Properties: &application.ApplicationProperties{
			AllowUpdates:   new(allowUpdates),
			DefaultVersion: new(defaultVersion),
			DisplayName:    new(displayName),
		},
	}

	if _, err := client.Update(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceBatchApplicationRead(d, meta)
}
