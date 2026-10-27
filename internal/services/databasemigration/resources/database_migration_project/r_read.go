// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package database_migration_project

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2021-06-30/projectresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatabaseMigrationProjectRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DatabaseMigration.ProjectsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := projectresource.ParseProjectID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ProjectsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ProjectName)
	d.Set("service_name", id.ServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		if props := model.Properties; props != nil {
			d.Set("source_platform", string(props.SourcePlatform))
			d.Set("target_platform", string(props.TargetPlatform))
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}
