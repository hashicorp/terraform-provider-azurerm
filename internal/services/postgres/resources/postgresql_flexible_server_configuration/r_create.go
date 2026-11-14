// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFlexibleServerConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Postgres.FlexibleServersConfigurationsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverId, err := servers.ParseFlexibleServerID(d.Get("server_id").(string))
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

	if err := client.UpdateCallbackThenPoll(ctx, id, props, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}
	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
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
					if err := meta.(*clients.Client).Postgres.RestartServer(ctx, *serverId); err != nil {
						return fmt.Errorf("restarting server %s: %+v", serverId, err)
					}
				}
			}
		}
	}

	return resourceFlexibleServerConfigurationRead(d, meta)
}
