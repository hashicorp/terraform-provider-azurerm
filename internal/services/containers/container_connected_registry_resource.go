package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/containerregistry/mgmt/2021-08-01-preview/containerregistry" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
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
			ValidateFunc: validate.RegistryID,
		},

		"parent_registry_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.Any(validate.ContainerConnectedRegistryID, validate.RegistryID),
		},

		"sync_token_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ContainerRegistryTokenID,
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
			Default:  string(containerregistry.ConnectedRegistryModeReadWrite),
			ValidateFunc: validation.StringInSlice(
				[]string{
					string(containerregistry.ConnectedRegistryModeMirror),
					string(containerregistry.ConnectedRegistryModeReadOnly),
					string(containerregistry.ConnectedRegistryModeReadWrite),
					string(containerregistry.ConnectedRegistryModeRegistry),
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
				ValidateFunc: validate.ContainerRegistryTokenID,
			},
		},

		"log_level": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  containerregistry.LogLevelNone,
			ValidateFunc: validation.StringInSlice(
				[]string{
					string(containerregistry.LogLevelNone),
					string(containerregistry.LogLevelDebug),
					string(containerregistry.LogLevelInformation),
					string(containerregistry.LogLevelWarning),
					string(containerregistry.LogLevelError),
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
	return validate.ContainerConnectedRegistryID
}

func (r ContainerConnectedRegistryResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.ConnectedRegistriesClient

			var model ContainerConnectedRegistryModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			rid, err := parse.RegistryID(model.ContainerRegistryId)
			if err != nil {
				return fmt.Errorf("parsing parent container registry id: %v", err)
			}
			id := parse.NewContainerConnectedRegistryID(rid.SubscriptionId, rid.ResourceGroup, rid.Name, model.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.RegistryName, id.ConnectedRegistryName)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			auditLogStatus := containerregistry.AuditLogStatusDisabled
			if model.AuditLogEnabled {
				auditLogStatus = containerregistry.AuditLogStatusEnabled
			}

			notifications, err := r.expandRepoNotifications(model.RepoNotifications)
			if err != nil {
				return fmt.Errorf("expanding `notification`: %+v", err)
			}

			params := containerregistry.ConnectedRegistry{
				ConnectedRegistryProperties: &containerregistry.ConnectedRegistryProperties{
					Mode: containerregistry.ConnectedRegistryMode(model.Mode),
					Parent: &containerregistry.ParentProperties{
						SyncProperties: &containerregistry.SyncProperties{
							TokenID:    utils.String(model.SyncTokenId),
							Schedule:   utils.String(model.SyncSchedule),
							SyncWindow: utils.String(model.SyncWindow),
							MessageTTL: utils.String(model.SyncMessageTTL),
						},
					},
					ClientTokenIds: &model.ClientTokenIds,
					Logging: &containerregistry.LoggingProperties{
						LogLevel:       containerregistry.LogLevel(model.LogLevel),
						AuditLogStatus: auditLogStatus,
					},
					NotificationsList: notifications,
				},
			}

			if model.ParentRegistryId != "" {
				if pid, err := parse.RegistryID(model.ParentRegistryId); err == nil {
					params.ConnectedRegistryProperties.Parent.ID = utils.String(pid.ID())
				} else if pid, err := parse.ContainerConnectedRegistryID(model.ParentRegistryId); err == nil {
					params.ConnectedRegistryProperties.Parent.ID = utils.String(pid.ID())
				}
			}

			future, err := client.Create(ctx, id.ResourceGroup, id.RegistryName, id.ConnectedRegistryName, params)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
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
			client := metadata.Client.Containers.ConnectedRegistriesClient
			id, err := parse.ContainerConnectedRegistryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.RegistryName, id.ConnectedRegistryName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			rid := parse.NewRegistryID(id.SubscriptionId, id.ResourceGroup, id.RegistryName)

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

			if props := existing.ConnectedRegistryProperties; props != nil {
				mode = string(props.Mode)
				if props.NotificationsList != nil {
					notificationList = *props.NotificationsList
				}
				if props.ClientTokenIds != nil {
					clientTokenIds = *props.ClientTokenIds
				}
				if logging := props.Logging; logging != nil {
					logLevel = string(logging.LogLevel)
					auditLogEnabled = logging.AuditLogStatus == containerregistry.AuditLogStatusEnabled
				}
				if parent := props.Parent; parent != nil {
					if parent.ID != nil {
						if pid, err := parse.RegistryIDInsensitively(*parent.ID); err == nil {
							parentRegistryId = pid.ID()
						} else if pid, err := parse.ContainerConnectedRegistryID(*parent.ID); err == nil {
							parentRegistryId = pid.ID()
						}
					}
					if sync := parent.SyncProperties; sync != nil {
						if sync.TokenID != nil {
							syncTokenId = *sync.TokenID
						}
						if sync.Schedule != nil {
							syncSchedule = *sync.Schedule
						}
						if sync.MessageTTL != nil {
							syncMessageTTL = *sync.MessageTTL
						}
						if sync.SyncWindow != nil {
							syncWindow = *sync.SyncWindow
						}
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
			client := metadata.Client.Containers.ConnectedRegistriesClient

			id, err := parse.ContainerConnectedRegistryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.RegistryName, id.ConnectedRegistryName)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				if !response.WasNotFound(future.Response()) {
					return fmt.Errorf("waiting for removal of %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (r ContainerConnectedRegistryResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.ContainerConnectedRegistryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ContainerConnectedRegistryModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Containers.ConnectedRegistriesClient

			existing, err := client.Get(ctx, id.ResourceGroup, id.RegistryName, id.ConnectedRegistryName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if props := existing.ConnectedRegistryProperties; props != nil {
				if metadata.ResourceData.HasChange("mode") {
					props.Mode = containerregistry.ConnectedRegistryMode(state.Mode)
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
						logging.LogLevel = containerregistry.LogLevel(state.LogLevel)
					}
					if metadata.ResourceData.HasChange("audit_log_enabled") {
						logging.AuditLogStatus = containerregistry.AuditLogStatusDisabled
						if state.AuditLogEnabled {
							logging.AuditLogStatus = containerregistry.AuditLogStatusEnabled
						}
					}
				}
				if parent := props.Parent; parent != nil {
					if sync := parent.SyncProperties; sync != nil {
						if metadata.ResourceData.HasChange("sync_token_id") {
							sync.TokenID = &state.SyncTokenId
						}
						if metadata.ResourceData.HasChange("sync_schedule") {
							sync.Schedule = &state.SyncSchedule
						}
						if metadata.ResourceData.HasChange("sync_message_ttl") {
							sync.MessageTTL = &state.SyncMessageTTL
						}
						if metadata.ResourceData.HasChange("sync_window") {
							sync.SyncWindow = &state.SyncWindow
						}
					}
				}
			}

			future, err := client.Create(ctx, id.ResourceGroup, id.RegistryName, id.ConnectedRegistryName, existing)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", id, err)
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
