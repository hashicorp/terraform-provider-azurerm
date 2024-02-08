// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mysql

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/serveradministrators"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMySQLAdministrator() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMySQLAdministratorCreateUpdate,
		Read:   resourceMySQLAdministratorRead,
		Update: resourceMySQLAdministratorCreateUpdate,
		Delete: resourceMySQLAdministratorDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AzureActiveDirectoryAdministratorID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"server_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"login": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"object_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourceMySQLAdministratorCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.MySqlClient.ServerAdministrators
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverName := d.Get("server_name").(string)
	resGroup := d.Get("resource_group_name").(string)
	login := d.Get("login").(string)

	id := parse.NewAzureActiveDirectoryAdministratorID(subscriptionId, resGroup, serverName, "activeDirectory")
	serverId := serveradministrators.NewServerID(subscriptionId, resGroup, serverName)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, serverId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_mysql_active_directory_administrator", id.ID())
		}
	}

	parameters := serveradministrators.ServerAdministratorResource{
		Properties: &serveradministrators.ServerAdministratorProperties{
			AdministratorType: serveradministrators.AdministratorTypeActiveDirectory,
			Login:             login,
			Sid:               d.Get("object_id").(string),
			TenantId:          d.Get("tenant_id").(string),
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, serverId, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return nil
}

func resourceMySQLAdministratorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.MySqlClient.ServerAdministrators
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AzureActiveDirectoryAdministratorID(d.Id())
	if err != nil {
		return err
	}

	serverId := serveradministrators.NewServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	resp, err := client.Get(ctx, serverId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("server_name", id.ServerName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("login", props.Login)
			d.Set("object_id", props.Sid)
			d.Set("tenant_id", props.TenantId)
		}
	}

	return nil
}

func resourceMySQLAdministratorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.MySqlClient.ServerAdministrators
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AzureActiveDirectoryAdministratorID(d.Id())
	if err != nil {
		return err
	}

	serverId := serveradministrators.NewServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	if err = client.DeleteThenPoll(ctx, serverId); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
