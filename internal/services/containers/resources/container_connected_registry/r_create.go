// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_connected_registry

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/connectedregistries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.ContainerRegistryClient.ConnectedRegistries

			var model ContainerConnectedRegistryModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			rid, err := registries.ParseRegistryID(model.ContainerRegistryId)
			if err != nil {
				return fmt.Errorf("parsing parent container registry id: %v", err)
			}
			id := connectedregistries.NewConnectedRegistryID(rid.SubscriptionId, rid.ResourceGroupName, rid.RegistryName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			auditLogStatus := connectedregistries.AuditLogStatusDisabled
			if model.AuditLogEnabled {
				auditLogStatus = connectedregistries.AuditLogStatusEnabled
			}

			notifications, err := r.expandRepoNotifications(model.RepoNotifications)
			if err != nil {
				return fmt.Errorf("expanding `notification`: %+v", err)
			}

			params := connectedregistries.ConnectedRegistry{
				Properties: &connectedregistries.ConnectedRegistryProperties{
					Mode: connectedregistries.ConnectedRegistryMode(model.Mode),
					Parent: connectedregistries.ParentProperties{
						SyncProperties: connectedregistries.SyncProperties{
							TokenId:    model.SyncTokenId,
							Schedule:   pointer.To(model.SyncSchedule),
							SyncWindow: pointer.To(model.SyncWindow),
							MessageTtl: model.SyncMessageTTL,
						},
					},
					ClientTokenIds: &model.ClientTokenIds,
					Logging: &connectedregistries.LoggingProperties{
						LogLevel:       pointer.ToEnum[connectedregistries.LogLevel](model.LogLevel),
						AuditLogStatus: pointer.To(auditLogStatus),
					},
					NotificationsList: notifications,
				},
			}

			if model.ParentRegistryId != "" {
				if pid, err := registries.ParseRegistryID(model.ParentRegistryId); err == nil {
					params.Properties.Parent.Id = pointer.To(pid.ID())
				} else if pid, err := connectedregistries.ParseConnectedRegistryID(model.ParentRegistryId); err == nil {
					params.Properties.Parent.Id = pointer.To(pid.ID())
				}
			}

			if err := client.CreateCallbackThenPoll(ctx, id, params, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r Resource) expandRepoNotifications(input []RepositoryNotification) (*[]string, error) {
	if len(input) == 0 {
		return nil, nil
	}

	result := make([]string, 0)

	for _, e := range input {
		if e.Digest != "" && e.Tag != "" {
			return nil, fmt.Errorf("notification %q has both `digest` and `tag` specified, only one of them is allowed", e.Name)
		}
		notification := parse.RepositoryNotification{
			Artifact: parse.Artifact{
				Name:   e.Name,
				Tag:    e.Tag,
				Digest: e.Digest,
			},
			Action: parse.RepositoryNotificationAction(e.Action),
		}
		result = append(result, notification.String())
	}

	return &result, nil
}
