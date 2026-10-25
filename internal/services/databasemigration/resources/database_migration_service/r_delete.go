// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package database_migration_service

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2021-06-30/serviceresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatabaseMigrationServiceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DatabaseMigration.ServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := serviceresource.ParseServiceID(d.Id())
	if err != nil {
		return err
	}

	opts := serviceresource.ServicesDeleteOperationOptions{
		DeleteRunningTasks: pointer.To(false),
	}
	if err := client.ServicesDeleteThenPoll(ctx, *id, opts); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
