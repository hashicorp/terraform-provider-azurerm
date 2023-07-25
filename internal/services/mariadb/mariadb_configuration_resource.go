// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mariadb

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mariadb/2018-06-01/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mariadb/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMariaDbConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMariaDbConfigurationCreateUpdate,
		Read:   resourceMariaDbConfigurationRead,
		Delete: resourceMariaDbConfigurationDelete,
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
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

func resourceMariaDbConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MariaDB.ConfigurationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := configurations.NewConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))
	value := d.Get("value").(string)

	properties := configurations.Configuration{
		Properties: &configurations.ConfigurationProperties{
			Value: utils.String(value),
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMariaDbConfigurationRead(d, meta)
}

func resourceMariaDbConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MariaDB.ConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurations.ParseConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found", *id)
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
			value := ""
			if v := props.Value; v != nil {
				value = *v
			}
			d.Set("value", value)
		}
	}

	return nil
}

func resourceMariaDbConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MariaDB.ConfigurationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurations.ParseConfigurationID(d.Id())
	if err != nil {
		return err
	}

	// "delete" = resetting this to the default value
	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	defaultValue := ""

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if v := props.DefaultValue; v != nil {
				defaultValue = *v
			}
		}
	}

	properties := configurations.Configuration{
		Properties: &configurations.ConfigurationProperties{
			// we can alternatively set `source: "system-default"`
			Value: &defaultValue,
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, properties); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
