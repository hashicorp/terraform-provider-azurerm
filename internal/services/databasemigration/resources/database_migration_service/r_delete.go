// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package database_migration_service

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2021-06-30/serviceresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	databasemigrationClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/databasemigration/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatabaseMigrationServiceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := databasemigrationClient.From(meta).ServicesClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
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
