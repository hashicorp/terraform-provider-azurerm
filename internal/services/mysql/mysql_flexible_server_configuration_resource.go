// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mysql

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/configurations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name mysql_flexible_server_configuration -service-package-name mysql -test-name characterSetServer -properties "name,resource_group_name,flexible_server_name:server_name" -known-values "subscription_id:data.Subscriptions.Primary"

var mysqlFlexibleServerConfigurationResourceName = "azurerm_mysql_flexible_server_configuration"

func resourceMySQLFlexibleServerConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMySQLFlexibleServerConfigurationCreate,
		Read:   resourceMySQLFlexibleServerConfigurationRead,
		Update: resourceMySQLFlexibleServerConfigurationUpdate,
		Delete: resourceMySQLFlexibleServerConfigurationDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&configurations.ConfigurationId{}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&configurations.ConfigurationId{}),
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

			"value": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
		},
	}
}

func resourceMySQLFlexibleServerConfigurationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FlexibleServers.Configurations
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM MySQL Configuration creation.")

	payload := configurations.Configuration{
		Properties: &configurations.ConfigurationProperties{
			Value: pointer.To(d.Get("value").(string)),
		},
	}

	// NOTE: this resource intentionally doesn't support Requires Import
	//       since a fallback route is created by default

	id := configurations.NewConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))

	locks.ByName(id.FlexibleServerName, mysqlFlexibleServerResourceName)
	defer locks.UnlockByName(id.FlexibleServerName, mysqlFlexibleServerResourceName)

	if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceMySQLFlexibleServerConfigurationRead(d, meta)
}

func resourceMySQLFlexibleServerConfigurationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FlexibleServers.Configurations
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM MySQL Configuration update.")

	id, err := configurations.ParseConfigurationID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.FlexibleServerName, mysqlFlexibleServerResourceName)
	defer locks.UnlockByName(id.FlexibleServerName, mysqlFlexibleServerResourceName)

	payload := configurations.Configuration{
		Properties: &configurations.ConfigurationProperties{
			Value: pointer.To(d.Get("value").(string)),
		},
	}

	if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %v", id, err)
	}

	return resourceMySQLFlexibleServerConfigurationRead(d, meta)
}

func resourceMySQLFlexibleServerConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FlexibleServers.Configurations
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurations.ParseConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return resourceMySQLFlexibleServerConfigurationFlatten(d, id, resp.Model)
}

func resourceMySQLFlexibleServerConfigurationFlatten(d *pluginsdk.ResourceData, id *configurations.ConfigurationId, dbConfig *configurations.Configuration) error {
	d.Set("name", id.ConfigurationName)
	d.Set("server_name", id.FlexibleServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	value := ""
	if dbConfig != nil {
		if props := dbConfig.Properties; props != nil {
			value = *props.Value
		}
	}
	d.Set("value", value)

	return pluginsdk.SetResourceIdentityData(d, id)
}

func resourceMySQLFlexibleServerConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FlexibleServers.Configurations
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurations.ParseConfigurationID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.FlexibleServerName, mysqlFlexibleServerResourceName)
	defer locks.UnlockByName(id.FlexibleServerName, mysqlFlexibleServerResourceName)

	// "delete" = resetting this to the default value
	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	payload := configurations.Configuration{
		Properties: &configurations.ConfigurationProperties{
			// we can alternatively set `source: "system-default"`
			Value: resp.Model.Properties.DefaultValue,
		},
	}

	if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("resetting %s to it's default value: %+v", *id, err)
	}

	return nil
}
