// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package database_migration_project

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2021-06-30/projectresource"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatabaseMigrationProjectCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DatabaseMigration.ProjectsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := projectresource.NewProjectID(subscriptionId, d.Get("resource_group_name").(string), d.Get("service_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.ProjectsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_database_migration_project", id.ID())
			}
		}
	}

	location := location.Normalize(d.Get("location").(string))
	sourcePlatform := d.Get("source_platform").(string)
	targetPlatform := d.Get("target_platform").(string)
	t := d.Get("tags").(map[string]any)

	parameters := projectresource.Project{
		Location: location,
		Properties: &projectresource.ProjectProperties{
			SourcePlatform: projectresource.ProjectSourcePlatform(sourcePlatform),
			TargetPlatform: projectresource.ProjectTargetPlatform(targetPlatform),
		},
		Tags: tags.Expand(t),
	}

	if _, err := client.ProjectsCreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}
	return resourceDatabaseMigrationProjectRead(d, meta)
}
