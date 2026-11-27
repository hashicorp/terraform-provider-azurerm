// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2026-10-01/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2026-10-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	postgresClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFlexibleServerConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := account.From(meta).SubscriptionId
	client := postgresClient.From(meta).FlexibleServersConfigurationsClient
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
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
				if features.From(meta).PostgresqlFlexibleServer.RestartServerOnConfigurationValueChange {
					if err := postgresClient.From(meta).RestartServer(ctx, *serverId); err != nil {
						return fmt.Errorf("restarting server %s: %+v", serverId, err)
					}
				}
			}
		}
	}

	return resourceFlexibleServerConfigurationRead(d, meta)
}
