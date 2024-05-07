// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package databasemigration

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2021-06-30/serviceresource"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databasemigration/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDatabaseMigrationService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDatabaseMigrationServiceCreate,
		Read:   resourceDatabaseMigrationServiceRead,
		Update: resourceDatabaseMigrationServiceUpdate,
		Delete: resourceDatabaseMigrationServiceDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := serviceresource.ParseServiceID(id)
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
				ValidateFunc: validate.ServiceName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateSubnetID,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					// No const defined in go sdk, the literal listed below is derived from the response of listskus endpoint.
					// See: https://docs.microsoft.com/en-us/rest/api/datamigration/resourceskus/listskus
					"Premium_4vCores",
					"Standard_1vCores",
					"Standard_2vCores",
					"Standard_4vCores",
				}, false),
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceDatabaseMigrationServiceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DatabaseMigration.ServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := serviceresource.NewServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.ServicesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_database_migration_service", id.ID())
		}
	}

	parameters := serviceresource.DataMigrationService{
		Location: azure.NormalizeLocation(d.Get("location").(string)),
		Properties: &serviceresource.DataMigrationServiceProperties{
			VirtualSubnetId: d.Get("subnet_id").(string),
		},
		Sku: &serviceresource.ServiceSku{
			Name: utils.String(d.Get("sku_name").(string)),
		},
		Kind: utils.String("Cloud"), // currently only "Cloud" is supported, hence hardcode here
	}
	if t, ok := d.GetOk("tags"); ok {
		parameters.Tags = tags.Expand(t.(map[string]interface{}))
	}

	if err := client.ServicesCreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDatabaseMigrationServiceRead(d, meta)
}

func resourceDatabaseMigrationServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DatabaseMigration.ServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := serviceresource.ParseServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ServicesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		if props := model.Properties; props != nil {
			d.Set("subnet_id", props.VirtualSubnetId)
		}
		d.Set("sku_name", model.Sku.Name)

		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceDatabaseMigrationServiceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DatabaseMigration.ServicesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := serviceresource.ParseServiceID(d.Id())
	if err != nil {
		return err
	}

	parameters := serviceresource.DataMigrationService{
		// location isn't update-able but if we don't supply the current value the SDK sends an empty string instead which errors on the API side
		Location: azure.NormalizeLocation(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.ServicesUpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceDatabaseMigrationServiceRead(d, meta)
}

func resourceDatabaseMigrationServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DatabaseMigration.ServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := serviceresource.ParseServiceID(d.Id())
	if err != nil {
		return err
	}

	opts := serviceresource.ServicesDeleteOperationOptions{
		DeleteRunningTasks: utils.Bool(false),
	}
	if err := client.ServicesDeleteThenPoll(ctx, *id, opts); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
