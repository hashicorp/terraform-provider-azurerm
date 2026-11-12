// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package database_migration_project

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2021-06-30/projectresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatabaseMigrationProjectDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DatabaseMigration.ProjectsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
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
