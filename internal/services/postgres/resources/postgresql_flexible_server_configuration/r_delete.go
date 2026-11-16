// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_configuration

import (
	"fmt"
	"slices"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFlexibleServerConfigurationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Postgres.FlexibleServersConfigurationsClient
	serverClient := meta.(*clients.Client).Postgres.FlexibleServersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurations.ParseConfigurationID(d.Id())
	if err != nil {
		return err
	}

	serverID := servers.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName)

	server, err := serverClient.Get(ctx, serverID)
	if err != nil {
		if response.WasNotFound(server.HttpResponse) {
			// if the server has already been deleted, we don't need to do anything
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", serverID, err)
	}

	isReplica, primaryId := false, ""
	if server.Model != nil && server.Model.Properties != nil {
		isReplica = slices.Contains([]servers.ReplicationRole{servers.ReplicationRoleAsyncReplica, servers.ReplicationRoleGeoAsyncReplica}, pointer.From(server.Model.Properties.ReplicationRole))
		primaryId = pointer.From(server.Model.Properties.SourceServerResourceId)
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

	// If the targeted server is a replica and the configuration item is one of the below, we'll reset the value to be equal to
	// that of the primary, otherwise Azure may return a validation error as these items cannot be less than the primary config (and the default could be).
	// reference: https://learn.microsoft.com/azure/postgresql/read-replica/concepts-read-replicas#scale
	if isReplica && primaryId != "" && slices.Contains([]string{"max_connections", "max_prepared_transactions", "max_locks_per_transaction", "max_wal_senders", "max_worker_processes"}, id.ConfigurationName) {
		pID, err := servers.ParseFlexibleServerIDInsensitively(primaryId)
		if err != nil {
			return err
		}

		primaryConfigurationID := configurations.NewConfigurationID(pID.SubscriptionId, pID.ResourceGroupName, pID.FlexibleServerName, id.ConfigurationName)
		primaryConfig, err := client.Get(ctx, primaryConfigurationID)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", primaryConfigurationID, err)
		}

		if primaryConfig.Model != nil && primaryConfig.Model.Properties != nil && primaryConfig.Model.Properties.Value != nil {
			defaultValue = *primaryConfig.Model.Properties.Value
		}
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
					if err := meta.(*clients.Client).Postgres.RestartServer(ctx, serverID); err != nil {
						return fmt.Errorf("restarting server %s: %+v", id, err)
					}
				}
			}
		}
	}

	return nil
}
