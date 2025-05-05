// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2024-08-01/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2024-08-01/serverrestart"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePostgresqlFlexibleServerConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFlexibleServerConfigurationUpdate,
		Read:   resourceFlexibleServerConfigurationRead,
		Update: resourceFlexibleServerConfigurationUpdate,
		Delete: resourceFlexibleServerConfigurationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := configurations.ParseConfigurationID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"server_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: configurations.ValidateFlexibleServerID,
			},

			"value": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceFlexibleServerConfigurationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Postgres.FlexibleServersConfigurationsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Postgresql Flexible Server configuration creation.")

	serverId, err := configurations.ParseFlexibleServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}
	id := configurations.NewConfigurationID(subscriptionId, serverId.ResourceGroupName, serverId.FlexibleServerName, d.Get("name").(string))

	locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
	defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

	props := configurations.ConfigurationForUpdate{
		Properties: &configurations.ConfigurationProperties{
			Value:  pointer.To(d.Get("value").(string)),
			Source: pointer.To("user-override"),
		},
	}

	if err := client.UpdateThenPoll(ctx, id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	resp, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if model := resp.Model; model != nil && model.Properties != nil {
		props := model.Properties

		if isDynamicConfig := props.IsDynamicConfig; isDynamicConfig != nil && !*isDynamicConfig {
			if isReadOnly := props.IsReadOnly; isReadOnly != nil && !*isReadOnly {
				if meta.(*clients.Client).Features.PostgresqlFlexibleServer.RestartServerOnConfigurationValueChange {
					restartClient := meta.(*clients.Client).Postgres.ServerRestartClient
					restartServerId := serverrestart.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName)

					if err = restartClient.ServersRestartThenPoll(ctx, restartServerId, serverrestart.RestartParameter{}); err != nil {
						return fmt.Errorf("restarting server %s: %+v", id, err)
					}
				}
			}
		}
	}

	d.SetId(id.ID())

	return resourceFlexibleServerConfigurationRead(d, meta)
}

func resourceFlexibleServerConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Postgres.FlexibleServersConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurations.ParseConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found, removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for %s: %+v", id, err)
	}

	d.Set("name", id.ConfigurationName)
	d.Set("server_id", configurations.NewFlexibleServerID(subscriptionId, id.ResourceGroupName, id.FlexibleServerName).ID())

	if resp.Model != nil && resp.Model.Properties != nil {
		d.Set("value", resp.Model.Properties.Value)
	}

	return nil
}

func resourceFlexibleServerConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.FlexibleServersConfigurationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurations.ParseConfigurationID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
	defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	defaultValue := ""
	if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.DefaultValue != nil {
		defaultValue = *resp.Model.Properties.DefaultValue
	}

	props := configurations.ConfigurationForUpdate{
		Properties: &configurations.ConfigurationProperties{
			Value:  &defaultValue,
			Source: pointer.To("user-override"),
		},
	}

	if err = client.UpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		props := model.Properties

		if isDynamicConfig := props.IsDynamicConfig; isDynamicConfig != nil && !*isDynamicConfig {
			if isReadOnly := props.IsReadOnly; isReadOnly != nil && !*isReadOnly {
				if meta.(*clients.Client).Features.PostgresqlFlexibleServer.RestartServerOnConfigurationValueChange {
					restartClient := meta.(*clients.Client).Postgres.ServerRestartClient
					restartServerId := serverrestart.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName)

					if err = restartClient.ServersRestartThenPoll(ctx, restartServerId, serverrestart.RestartParameter{}); err != nil {
						return fmt.Errorf("restarting server %s: %+v", id, err)
					}
				}
			}
		}
	}

	return nil
}
