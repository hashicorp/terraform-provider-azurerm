// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourcePostgreSQLConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePostgreSQLConfigurationCreateUpdate,
		Read:   resourcePostgreSQLConfigurationRead,
		Delete: resourcePostgreSQLConfigurationDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := configurations.ParseConfigurationID(id)
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
				ForceNew: true,
			},
		},
	}
}

func resourcePostgreSQLConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ConfigurationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := configurations.NewConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))
	// TODO: support RequiresImport - this is possible to tell if it's the non-default value from the API (see Delete)

	locks.ByName(id.ServerName, postgreSQLServerResourceName)
	defer locks.UnlockByName(id.ServerName, postgreSQLServerResourceName)

	properties := configurations.Configuration{
		Properties: &configurations.ConfigurationProperties{
			Value: utils.String(d.Get("value").(string)),
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourcePostgreSQLConfigurationRead(d, meta)
}

func resourcePostgreSQLConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurations.ParseConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ConfigurationName)
	d.Set("server_name", id.ServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("value", props.Value)
		}
	}

	return nil
}

func resourcePostgreSQLConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ConfigurationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurations.ParseConfigurationID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ServerName, postgreSQLServerResourceName)
	defer locks.UnlockByName(id.ServerName, postgreSQLServerResourceName)

	// "delete" = resetting this to the default value
	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	defaultValue := ""
	if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.DefaultValue != nil {
		defaultValue = *resp.Model.Properties.DefaultValue
	}

	properties := configurations.Configuration{
		Properties: &configurations.ConfigurationProperties{
			// we can alternatively set `source: "system-default"`
			Value: &defaultValue,
		},
	}

	if err = client.CreateOrUpdateThenPoll(ctx, *id, properties); err != nil {
		return fmt.Errorf("resetting %s to it's default value: %+v", *id, err)
	}

	return nil
}
