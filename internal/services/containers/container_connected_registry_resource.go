// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/connectedregistries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/tokens"
	tfvalidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContainerConnectedRegistryResource struct{}

var _ sdk.ResourceWithUpdate = ContainerConnectedRegistryResource{}

type ContainerConnectedRegistryModel struct {
	Name                string                   `tfschema:"name"`
	ContainerRegistryId string                   `tfschema:"container_registry_id"`
	ParentRegistryId    string                   `tfschema:"parent_registry_id"`
	SyncTokenId         string                   `tfschema:"sync_token_id"`
	SyncSchedule        string                   `tfschema:"sync_schedule"`
	SyncMessageTTL      string                   `tfschema:"sync_message_ttl"`
	SyncWindow          string                   `tfschema:"sync_window"`
	Mode                string                   `tfschema:"mode"`
	RepoNotifications   []RepositoryNotification `tfschema:"notification"`
	ClientTokenIds      []string                 `tfschema:"client_token_ids"`
	LogLevel            string                   `tfschema:"log_level"`
	AuditLogEnabled     bool                     `tfschema:"audit_log_enabled"`
}

type RepositoryNotification struct {
	Name   string `tfschema:"name"`
	Tag    string `tfschema:"tag"`
	Digest string `tfschema:"digest"`
	Action string `tfschema:"action"`
}

func (r ContainerConnectedRegistryResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ContainerRegistryName,
		},

		"container_registry_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: registries.ValidateRegistryID,
		},

		"parent_registry_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.Any(connectedregistries.ValidateConnectedRegistryID, registries.ValidateRegistryID),
		},

		"sync_token_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: tokens.ValidateTokenID,
		},

		"sync_schedule": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "* * * * *",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"sync_message_ttl": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "P1D",
			ValidateFunc: tfvalidate.ISO8601DurationBetween("P1D", "P90D"),
		},

		"sync_window": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: tfvalidate.ISO8601DurationBetween("PT3H", "P7D"),
		},

		"mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(connectedregistries.ConnectedRegistryModeReadWrite),
			ValidateFunc: validation.StringInSlice(
				[]string{
					string(connectedregistries.ConnectedRegistryModeMirror),
					string(connectedregistries.ConnectedRegistryModeReadOnly),
					string(connectedregistries.ConnectedRegistryModeReadWrite),
					string(connectedregistries.ConnectedRegistryModeRegistry),
				},
				false,
			),
		},

		"notification": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"tag": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"digest": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"action": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(parse.RepositoryNotificationActionPush),
							string(parse.RepositoryNotificationActionDelete),
							string(parse.RepositoryNotificationActionAny),
						}, false),
					},
				},
			},
		},

		"client_token_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: tokens.ValidateTokenID,
			},
		},

		"log_level": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  connectedregistries.LogLevelNone,
			ValidateFunc: validation.StringInSlice(
				[]string{
					string(connectedregistries.LogLevelNone),
					string(connectedregistries.LogLevelDebug),
					string(connectedregistries.LogLevelInformation),
					string(connectedregistries.LogLevelWarning),
					string(connectedregistries.LogLevelError),
				},
				false,
			),
		},

		"audit_log_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func (r ContainerConnectedRegistryResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ContainerConnectedRegistryResource) ResourceType() string {
	return "azurerm_container_connected_registry"
}

func (r ContainerConnectedRegistryResource) ModelObject() interface{} {
	return &ContainerConnectedRegistryModel{}
}

func (r ContainerConnectedRegistryResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return connectedregistries.ValidateConnectedRegistryID
}

func (r ContainerConnectedRegistryResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.ContainerRegistryClient_v2021_08_01_preview.ConnectedRegistries

			var model ContainerConnectedRegistryModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			rid, err := registries.ParseRegistryID(model.ContainerRegistryId)
			if err != nil {
				return fmt.Errorf("parsing parent container registry id: %v", err)
			}
			id := connectedregistries.NewConnectedRegistryID(rid.SubscriptionId, rid.ResourceGroupName, rid.RegistryName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
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
							Schedule:   utils.String(model.SyncSchedule),
							SyncWindow: utils.String(model.SyncWindow),
							MessageTtl: model.SyncMessageTTL,
						},
					},
					ClientTokenIds: &model.ClientTokenIds,
					Logging: &connectedregistries.LoggingProperties{
						LogLevel:       pointer.To(connectedregistries.LogLevel(model.LogLevel)),
						AuditLogStatus: pointer.To(auditLogStatus),
					},
					NotificationsList: notifications,
				},
			}

			if model.ParentRegistryId != "" {
				if pid, err := registries.ParseRegistryID(model.ParentRegistryId); err == nil {
					params.Properties.Parent.Id = utils.String(pid.ID())
				} else if pid, err := connectedregistries.ParseConnectedRegistryID(model.ParentRegistryId); err == nil {
					params.Properties.Parent.Id = utils.String(pid.ID())
				}
			}

			if err := client.CreateThenPoll(ctx, id, params); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ContainerConnectedRegistryResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.ContainerRegistryClient_v2021_08_01_preview.ConnectedRegistries
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

func (r ContainerConnectedRegistryResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.ContainerRegistryClient_v2021_08_01_preview.ConnectedRegistries

			id, err := connectedregistries.ParseConnectedRegistryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ContainerConnectedRegistryResource) Update() sdk.ResourceFunc {
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

			client := metadata.Client.Containers.ContainerRegistryClient_v2021_08_01_preview.ConnectedRegistries

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
						logging.LogLevel = pointer.To(connectedregistries.LogLevel(state.LogLevel))
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

func (r ContainerConnectedRegistryResource) expandRepoNotifications(input []RepositoryNotification) (*[]string, error) {
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

func (r ContainerConnectedRegistryResource) flattenRepoNotifications(input []string) ([]RepositoryNotification, error) {
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
			Name:   notification.Artifact.Name,
			Tag:    notification.Artifact.Tag,
			Digest: notification.Artifact.Digest,
			Action: string(notification.Action),
		})
	}

	return output, nil
}
