// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_connected_registry

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/connectedregistries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.ContainerRegistryClient.ConnectedRegistries
			id, err := connectedregistries.ParseConnectedRegistryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			rid := registries.NewRegistryID(id.SubscriptionId, id.ResourceGroupName, id.RegistryName)

			var (
				mode             string
				parentRegistryId string
				syncTokenId      string
				syncSchedule     string
				syncMessageTTL   string
				syncWindow       string
				notificationList []string
				clientTokenIds   []string
				logLevel         string
				auditLogEnabled  bool
			)

			if model := existing.Model; model != nil {
				if props := model.Properties; props != nil {
					mode = string(props.Mode)

					if props.NotificationsList != nil {
						notificationList = *props.NotificationsList
					}

					if props.ClientTokenIds != nil {
						clientTokenIds = *props.ClientTokenIds
					}

					if logging := props.Logging; logging != nil {
						logLevel = string(*logging.LogLevel)
						auditLogEnabled = *logging.AuditLogStatus == connectedregistries.AuditLogStatusEnabled
					}

					parent := props.Parent
					if parent.Id != nil {
						if pid, err := registries.ParseRegistryIDInsensitively(*parent.Id); err == nil {
							parentRegistryId = pid.ID()
						} else if pid, err := connectedregistries.ParseConnectedRegistryID(*parent.Id); err == nil {
							parentRegistryId = pid.ID()
						}
					}

					sync := parent.SyncProperties

					syncTokenId = sync.TokenId
					syncMessageTTL = sync.MessageTtl

					if sync.Schedule != nil {
						syncSchedule = *sync.Schedule
					}

					if sync.SyncWindow != nil {
						syncWindow = *sync.SyncWindow
					}
				}
			}

			notifications, err := r.flattenRepoNotifications(notificationList)
			if err != nil {
				return fmt.Errorf("flattening `notification`: %+v", err)
			}

			model := ContainerConnectedRegistryModel{
				Name:                id.ConnectedRegistryName,
				ContainerRegistryId: rid.ID(),
				ParentRegistryId:    parentRegistryId,
				SyncTokenId:         syncTokenId,
				SyncSchedule:        syncSchedule,
				SyncMessageTTL:      syncMessageTTL,
				SyncWindow:          syncWindow,
				Mode:                mode,
				RepoNotifications:   notifications,
				ClientTokenIds:      clientTokenIds,
				LogLevel:            logLevel,
				AuditLogEnabled:     auditLogEnabled,
			}

			return metadata.Encode(&model)
		},
	}
}

func (r Resource) flattenRepoNotifications(input []string) ([]RepositoryNotification, error) {
	if input == nil {
		return []RepositoryNotification{}, nil
	}

	output := make([]RepositoryNotification, 0)

	for _, e := range input {
		notification, err := parse.ParseRepositoryNotification(e)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", e, err)
		}
		output = append(output, RepositoryNotification{
			Name:   notification.Name,
			Tag:    notification.Tag,
			Digest: notification.Digest,
			Action: string(notification.Action),
		})
	}

	return output, nil
}

type RepositoryNotification struct {
	Name   string `tfschema:"name"`
	Tag    string `tfschema:"tag"`
	Digest string `tfschema:"digest"`
	Action string `tfschema:"action"`
}
