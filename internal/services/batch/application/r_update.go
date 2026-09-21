// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/application"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBatchApplicationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.ApplicationClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
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
			AllowUpdates:   pointer.To(allowUpdates),
			DefaultVersion: pointer.To(defaultVersion),
			DisplayName:    pointer.To(displayName),
		},
	}

	if _, err := client.Update(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceBatchApplicationRead(d, meta)
}
