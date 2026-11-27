// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package database_migration_project

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2021-06-30/projectresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	databasemigrationClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/databasemigration/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatabaseMigrationProjectDelete(d *pluginsdk.ResourceData, meta any) error {
	client := databasemigrationClient.From(meta).ProjectsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := projectresource.ParseProjectID(d.Id())
	if err != nil {
		return err
	}

	opts := projectresource.ProjectsDeleteOperationOptions{
		DeleteRunningTasks: pointer.To(false),
	}
	if _, err := client.ProjectsDelete(ctx, *id, opts); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
