// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/servers"
	flexibleServers "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

const (
	springCloudAppMysqlAssociationKeyDatabase = "databaseName"
	springCloudAppMysqlAssociationKeyUsername = "username"
)

func resourceSpringCloudAppMysqlAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		DeprecationMessage: features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_app_mysql_association` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information."),

		Create: resourceSpringCloudAppMysqlAssociationCreateUpdate,
		Read:   resourceSpringCloudAppMysqlAssociationRead,
		Update: resourceSpringCloudAppMysqlAssociationCreateUpdate,
		Delete: resourceSpringCloudAppMysqlAssociationDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SpringCloudAppMySqlAssociationV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.SpringCloudAppAssociationID(id)
			return err
		}, importSpringCloudAppAssociation(springCloudAppAssociationTypeMysql)),

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
				ValidateFunc: validate.SpringCloudAppAssociationName,
			},

			"spring_cloud_app_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudAppID,
			},

			"mysql_server_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					servers.ValidateServerID,
					flexibleServers.ValidateFlexibleServerID,
				),
			},

			"database_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"username": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"password": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceSpringCloudAppMysqlAssociationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.BindingsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	appId, err := parse.SpringCloudAppID(d.Get("spring_cloud_app_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewSpringCloudAppAssociationID(appId.SubscriptionId, appId.ResourceGroup, appId.SpringName, appId.AppName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.BindingName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_spring_cloud_app_mysql_association", id.ID())
		}
	}

	bindingResource := appplatform.BindingResource{
		Properties: &appplatform.BindingResourceProperties{
			BindingParameters: map[string]*string{
				springCloudAppMysqlAssociationKeyDatabase: utils.String(d.Get("database_name").(string)),
				springCloudAppMysqlAssociationKeyUsername: utils.String(d.Get("username").(string)),
			},
			Key:        utils.String(d.Get("password").(string)),
			ResourceID: utils.String(d.Get("mysql_server_id").(string)),
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.BindingName, bindingResource)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSpringCloudAppMysqlAssociationRead(d, meta)
}

func resourceSpringCloudAppMysqlAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.BindingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudAppAssociationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.BindingName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud App Association %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.Set("name", id.BindingName)
	d.Set("spring_cloud_app_id", parse.NewSpringCloudAppID(id.SubscriptionId, id.ResourceGroup, id.SpringName, id.AppName).ID())
	if props := resp.Properties; props != nil {
		d.Set("mysql_server_id", props.ResourceID)

		databaseName := ""
		if v, ok := props.BindingParameters[springCloudAppMysqlAssociationKeyDatabase]; ok && v != nil {
			databaseName = *v
		}
		d.Set("database_name", databaseName)

		username := ""
		if v, ok := props.BindingParameters[springCloudAppMysqlAssociationKeyUsername]; ok && v != nil {
			username = *v
		}
		d.Set("username", username)
	}
	return nil
}

func resourceSpringCloudAppMysqlAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.BindingsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudAppAssociationID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.BindingName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %q: %+v", id, err)
	}

	return nil
}
