// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_connected_registry

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/connectedregistries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := connectedregistries.ParseConnectedRegistryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ContainerConnectedRegistryModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Containers.ContainerRegistryClient.ConnectedRegistries

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model != nil && existing.Model.Properties != nil {
				props := existing.Model.Properties
				if metadata.ResourceData.HasChange("mode") {
					props.Mode = connectedregistries.ConnectedRegistryMode(state.Mode)
				}
				if metadata.ResourceData.HasChange("notification") {
					notifications, err := r.expandRepoNotifications(state.RepoNotifications)
					if err != nil {
						return fmt.Errorf("expanding `notification`: %+v", err)
					}
					props.NotificationsList = notifications
				}
				if metadata.ResourceData.HasChange("client_token_ids") {
					props.ClientTokenIds = &state.ClientTokenIds
				}
				if logging := props.Logging; logging != nil {
					if metadata.ResourceData.HasChange("log_level") {
						logging.LogLevel = pointer.ToEnum[connectedregistries.LogLevel](state.LogLevel)
					}
					if metadata.ResourceData.HasChange("audit_log_enabled") {
						logging.AuditLogStatus = pointer.To(connectedregistries.AuditLogStatusDisabled)
						if state.AuditLogEnabled {
							logging.AuditLogStatus = pointer.To(connectedregistries.AuditLogStatusEnabled)
						}
					}
				}

				sync := props.Parent.SyncProperties
				if metadata.ResourceData.HasChange("sync_token_id") {
					sync.TokenId = state.SyncTokenId
				}
				if metadata.ResourceData.HasChange("sync_schedule") {
					sync.Schedule = &state.SyncSchedule
				}
				if metadata.ResourceData.HasChange("sync_message_ttl") {
					sync.MessageTtl = state.SyncMessageTTL
				}
				if metadata.ResourceData.HasChange("sync_window") {
					sync.SyncWindow = &state.SyncWindow
				}
			}

			if err := client.CreateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
