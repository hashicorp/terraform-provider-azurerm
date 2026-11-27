// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package database_migration_project

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2021-06-30/projectresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	databasemigrationClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/databasemigration/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databasemigration/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Database Migration Project.",
		Read:        dataSourceDatabaseMigrationProjectRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "Name of the database migration project.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ProjectName,
			},

			"service_name": {
				Description:  "Name of the database migration service where resource belongs to.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ServiceName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"source_platform": {
				Description: "The platform type of the migration source.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"target_platform": {
				Description: "The platform type of the migration target.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceDatabaseMigrationProjectRead(d *pluginsdk.ResourceData, meta any) error {
	client := databasemigrationClient.From(meta).ProjectsClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id := projectresource.NewProjectID(subscriptionId, d.Get("resource_group_name").(string), d.Get("service_name").(string), d.Get("name").(string))
	resp, err := client.ProjectsGet(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

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
