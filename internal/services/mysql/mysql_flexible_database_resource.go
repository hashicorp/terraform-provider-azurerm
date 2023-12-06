// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mysql

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/databases"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMySqlFlexibleDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMySqlFlexibleDatabaseCreate,
		Read:   resourceMySqlFlexibleDatabaseRead,
		Delete: resourceMySqlFlexibleDatabaseDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := databases.ParseDatabaseID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"server_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerName,
			},

			"charset": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"collation": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceMySqlFlexibleDatabaseCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FlexibleServers.Databases
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := databases.NewDatabaseID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_mysql_flexible_database", id.ID())
		}
	}

	payload := databases.Database{
		Properties: &databases.DatabaseProperties{
			Charset:   pointer.To(d.Get("charset").(string)),
			Collation: pointer.To(d.Get("collation").(string)),
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	d.SetId(id.ID())
	return resourceMySqlFlexibleDatabaseRead(d, meta)
}

func resourceMySqlFlexibleDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FlexibleServers.Databases
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := databases.ParseDatabaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.DatabaseName)
	d.Set("server_name", id.FlexibleServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("charset", props.Charset)
			d.Set("collation", props.Collation)
		}
	}

	return nil
}

func resourceMySqlFlexibleDatabaseDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FlexibleServers.Databases
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := databases.ParseDatabaseID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
