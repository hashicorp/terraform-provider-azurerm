// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package database_migration_service

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2021-06-30/serviceresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatabaseMigrationServiceUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DatabaseMigration.ServicesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := serviceresource.ParseServiceID(d.Id())
	if err != nil {
		return err
	}

	parameters := serviceresource.DataMigrationService{
		// location isn't update-able but if we don't supply the current value the SDK sends an empty string instead which errors on the API side
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.ServicesUpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceDatabaseMigrationServiceRead(d, meta)
}
