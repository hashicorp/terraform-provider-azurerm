// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package databasemigration

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2021-06-30/projectresource"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databasemigration/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDatabaseMigrationProject() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDatabaseMigrationProjectCreateUpdate,
		Read:   resourceDatabaseMigrationProjectRead,
		Update: resourceDatabaseMigrationProjectCreateUpdate,
		Delete: resourceDatabaseMigrationProjectDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := projectresource.ParseProjectID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ProjectName,
			},

			"service_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServiceName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"source_platform": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(projectresource.ProjectSourcePlatformMongoDb),
					string(projectresource.ProjectSourcePlatformMySQL),
					string(projectresource.ProjectSourcePlatformPostgreSql),
					string(projectresource.ProjectSourcePlatformSQL),
					string(projectresource.ProjectSourcePlatformUnknown),
				}, false),
			},

			"target_platform": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(projectresource.ProjectTargetPlatformAzureDbForMySql),
					string(projectresource.ProjectTargetPlatformAzureDbForPostgreSql),
					string(projectresource.ProjectTargetPlatformMongoDb),
					string(projectresource.ProjectTargetPlatformSQLDB),
					string(projectresource.ProjectTargetPlatformSQLMI),
					string(projectresource.ProjectTargetPlatformUnknown),
				}, false),
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceDatabaseMigrationProjectCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DatabaseMigration.ProjectsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := projectresource.NewProjectID(subscriptionId, d.Get("resource_group_name").(string), d.Get("service_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
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

	location := azure.NormalizeLocation(d.Get("location").(string))
	sourcePlatform := d.Get("source_platform").(string)
	targetPlatform := d.Get("target_platform").(string)
	t := d.Get("tags").(map[string]interface{})

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
	return resourceDatabaseMigrationProjectRead(d, meta)
}

func resourceDatabaseMigrationProjectRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceDatabaseMigrationProjectDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DatabaseMigration.ProjectsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := projectresource.ParseProjectID(d.Id())
	if err != nil {
		return err
	}

	opts := projectresource.ProjectsDeleteOperationOptions{
		DeleteRunningTasks: utils.Bool(false),
	}
	if _, err := client.ProjectsDelete(ctx, *id, opts); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
