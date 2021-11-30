package containers

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"

	legacyacr "github.com/Azure/azure-sdk-for-go/services/preview/containerregistry/mgmt/2019-06-01-preview/containerregistry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	msivalidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/msi/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContainerRegistryTaskResource struct{}

//var _ sdk.ResourceWithUpdate = ContainerRegistryTaskResource{}

type AgentConfig struct {
	CPU int `tfschema:"cpu"`
}

type Platform struct {
	OS           string `tfschema:"os"`
	Architecture string `tfschema:"architecture"`
	Variant      string `tfschema:"variant"`
}

type Argument struct {
	Name     string `tfschema:"name"`
	Value    string `tfschema:"value"`
	IsSecret bool   `tfschema:"is_secret"`
}

type Value struct {
	Name     string `tfschema:"name"`
	Value    string `tfschema:"value"`
	IsSecret bool   `tfschema:"is_secret"`
}

type DockerStep struct {
	ContextPath        string     `tfschema:"context_path"`
	ContextAccessToken string     `tfschema:"context_access_token"`
	Arguments          []Argument `tfschema:"argument"`
	DockerfilePath     string     `tfschema:"dockerfile_path"`
	ImageNames         []string   `tfschema:"image_names"`
	IsPushEnabled      bool       `tfschema:"is_push_enabled"`
	IsCacheEnabled     bool       `tfschema:"is_cache_enabled"`
	Target             string     `tfschema:"target"`
}

type FileTaskStep struct {
	ContextPath        string  `tfschema:"context_path"`
	ContextAccessToken string  `tfschema:"context_access_token"`
	TaskFilePath       string  `tfschema:"task_file_path"`
	ValueFilePath      string  `tfschema:"value_file_path"`
	Values             []Value `tfschema:"value"`
}

type EncodedTaskStep struct {
	ContextPath        string  `tfschema:"context_path"`
	ContextAccessToken string  `tfschema:"context_access_token"`
	TaskContent        string  `tfschema:"task_content"`
	ValueContent       string  `tfschema:"value_content"`
	Values             []Value `tfschema:"value"`
}

type BaseImageTrigger struct {
	Name                     string `tfschema:"name"`
	Type                     string `tfschema:"type"`
	Enabled                  bool   `tfschema:"enabled"`
	UpdateTriggerEndpoint    string `tfschema:"update_trigger_endpoint"`
	UpdateTriggerPayloadType string `tfschema:"update_trigger_payload_type"`
}

type Auth struct {
	TokenType    string `tfschema:"token_type"`
	Token        string `tfschema:"token"`
	RefreshToken string `tfschema:"refresh_token"`
	Scope        string `tfschema:"scope"`
	ExpireInSec  int    `tfschema:"expire_in_seconds"`
}

type SourceSetting struct {
	SourceType    string `tfschema:"source_type"`
	RepositoryURL string `tfschema:"repository_url"`
	Branch        string `tfschema:"branch"`
	Auth          []Auth `tfschema:"auth"`
}

type SourceTrigger struct {
	Name          string          `tfschema:"name"`
	Enabled       bool            `tfschema:"enabled"`
	Events        []string        `tfschema:"events"`
	SourceSetting []SourceSetting `tfschema:"source_setting"`
}

type TimerTrigger struct {
	Name     string `tfschema:"name"`
	Enabled  bool   `tfschema:"enabled"`
	Schedule string `tfschema:"schedule"`
}

type Identity struct {
	Type        string   `tfschema:"type"`
	IdentityIds []string `tfschema:"identity_ids"`
	PrincipalId string   `tfschema:"principal_id"`
	TenantId    string   `tfschema:"tenant_id"`
}

type ContainerRegistryTaskModel struct {
	Name                string                 `tfschema:"name"`
	ContainerRegistryId string                 `tfschema:"container_registry_id"`
	Identity            []Identity             `tfschema:"identity"`
	AgentConfig         []AgentConfig          `tfschema:"agent_setting"`
	AgentPoolName       string                 `tfschema:"agent_pool_name"`
	IsSystemTask        bool                   `tfschema:"is_system_task"`
	LogTemplate         string                 `tfschema:"log_template"`
	Platform            []Platform             `tfschema:"platform_setting"`
	Enabled             bool                   `tfschema:"enabled"`
	TimeoutInSec        int                    `tfschema:"timeout_in_seconds"`
	DockerStep          []DockerStep           `tfschema:"docker_step"`
	FileTaskStep        []FileTaskStep         `tfschema:"file_task_step"`
	EncodedTaskStep     []EncodedTaskStep      `tfschema:"encoded_task_step"`
	BaseImageTrigger    []BaseImageTrigger     `tfschema:"base_image_trigger"`
	SourceTrigger       []SourceTrigger        `tfschema:"source_trigger"`
	TimerTrigger        []TimerTrigger         `tfschema:"timer_trigger"`
	Tags                map[string]interface{} `tfschema:"tags"`
}

func userDataStateFunc(v interface{}) string {
	switch s := v.(type) {
	case string:
		s = utils.Base64EncodeIfNot(s)
		hash := sha1.Sum([]byte(s))
		return hex.EncodeToString(hash[:])
	default:
		return ""
	}
}

func (r ContainerRegistryTaskResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ContainerRegistryTaskName,
		},
		"container_registry_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RegistryID,
		},
		"platform_setting": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"os": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(legacyacr.Windows),
							string(legacyacr.Linux),
						}, false),
					},
					"architecture": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(legacyacr.Amd64),
							string(legacyacr.Arm),
							string(legacyacr.Arm64),
							string(legacyacr.ThreeEightSix),
							string(legacyacr.X86),
						}, false),
					},
					"variant": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(legacyacr.V6),
							string(legacyacr.V7),
							string(legacyacr.V8),
						}, false),
					},
				},
			},
		},
		"docker_step": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"dockerfile_path": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"context_path": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"context_access_token": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"argument": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"value": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"is_secret": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
							},
						},
					},
					"image_names": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"is_push_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"is_cache_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"target": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
			ConflictsWith: []string{"file_task_step", "encoded_task_step"},
		},
		"file_task_step": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"task_file_path": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"value_file_path": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"context_path": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"context_access_token": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"value": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"value": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"is_secret": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
							},
						},
					},
				},
			},
			ConflictsWith: []string{"docker_step", "encoded_task_step"},
		},
		"encoded_task_step": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"task_content": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						StateFunc:    userDataStateFunc,
					},
					"value_content": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						StateFunc:    userDataStateFunc,
					},
					"context_path": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"context_access_token": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"value": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"value": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"is_secret": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
							},
						},
					},
				},
			},
			ConflictsWith: []string{"docker_step", "file_task_step"},
		},
		"base_image_trigger": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(legacyacr.All),
							string(legacyacr.Runtime),
						}, false),
					},
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"update_trigger_endpoint": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"update_trigger_payload_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(legacyacr.UpdateTriggerPayloadTypeDefault),
							string(legacyacr.UpdateTriggerPayloadTypeDefault),
						}, false),
					},
				},
			},
		},
		"source_trigger": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"source_setting": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"source_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(legacyacr.Github),
										string(legacyacr.VisualStudioTeamService),
									}, false),
								},
								"repository_url": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"branch": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"auth": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*schema.Schema{
											"token_type": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringInSlice([]string{
													string(legacyacr.PAT),
													string(legacyacr.OAuth),
												}, false),
											},
											"token": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												Sensitive:    true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
											"refresh_token": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
											"scope": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
											"expire_in_seconds": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"events": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
		"timer_trigger": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"schedule": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},
		"identity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(legacyacr.UserAssigned),
							string(legacyacr.SystemAssigned),
							string(legacyacr.SystemAssignedUserAssigned),
						}, false),
					},
					"identity_ids": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: msivalidate.UserAssignedIdentityID,
						},
					},
					"principal_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
		"agent_setting": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"cpu": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
				},
			},
		},
		"agent_pool_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},
		"is_system_task": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
		"log_template": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"timeout_in_seconds": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(300, 28800),
			Default:      3600,
		},
		"tags": tags.ForceNewSchema(),
	}
}

func (r ContainerRegistryTaskResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ContainerRegistryTaskResource) ResourceType() string {
	return "azurerm_container_registry_task"
}

func (r ContainerRegistryTaskResource) ModelObject() interface{} {
	return &ContainerRegistryTaskModel{}
}

func (r ContainerRegistryTaskResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ContainerRegistryTaskID
}

func (r ContainerRegistryTaskResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.TasksClient
			registryClient := metadata.Client.Containers.RegistriesClient

			var model ContainerRegistryTaskModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			registryId, err := parse.RegistryID(model.ContainerRegistryId)
			if err != nil {
				return err
			}

			registry, err := registryClient.Get(ctx, registryId.ResourceGroup, registryId.Name)
			if err != nil {
				return fmt.Errorf("getting registry %s: %+v", registryId, err)
			}

			id := parse.NewContainerRegistryTaskID(registryId.SubscriptionId, registryId.ResourceGroup, registryId.Name, model.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.RegistryName, id.TaskName)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			status := legacyacr.TaskStatusDisabled
			if model.Enabled {
				status = legacyacr.TaskStatusEnabled
			}

			params := legacyacr.Task{
				TaskProperties: &legacyacr.TaskProperties{
					Platform:     expandRegistryTaskPlatform(model.Platform[0]),
					Step:         expandRegistryTaskStep(model),
					Trigger:      expandRegistryTaskTrigger(model),
					Status:       status,
					IsSystemTask: &model.IsSystemTask,
					Timeout:      utils.Int32(int32(model.TimeoutInSec)),
					// TODO
					//Credentials:  nil,
				},

				// The location of the task must be the same as the registry, otherwise the API will raise error complaining can't find the registry.
				Location: registry.Location,
				Identity: expandRegistryTaskIdentity(model.Identity),
				Tags:     tags.Expand(model.Tags),
			}

			if len(model.AgentConfig) != 0 {
				agentConfig := model.AgentConfig[0]
				params.TaskProperties.AgentConfiguration = &legacyacr.AgentProperties{CPU: utils.Int32(int32(agentConfig.CPU))}

			}

			if model.AgentPoolName != "" {
				params.TaskProperties.AgentPoolName = &model.AgentPoolName
			}

			if model.LogTemplate != "" {
				params.TaskProperties.LogTemplate = &model.LogTemplate
			}

			future, err := client.Create(ctx, id.ResourceGroup, id.RegistryName, id.TaskName, params)
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

func (r ContainerRegistryTaskResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.TasksClient
			id, err := parse.ContainerRegistryTaskID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			registryId := parse.NewRegistryID(id.SubscriptionId, id.ResourceGroup, id.RegistryName)

			task, err := client.Get(ctx, id.ResourceGroup, id.RegistryName, id.TaskName)
			if err != nil {
				if utils.ResponseWasNotFound(task.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			var (
				agentConfig      []AgentConfig
				agentPoolName    string
				isSystemTask     bool
				logTemplate      string
				platform         []Platform
				enabled          bool
				timeoutInSec     int
				dockerStep       []DockerStep
				fileTaskStep     []FileTaskStep
				encodedTaskStep  []EncodedTaskStep
				baseImageTrigger []BaseImageTrigger
				sourceTrigger    []SourceTrigger
				timerTrigger     []TimerTrigger
			)
			if props := task.TaskProperties; props != nil {
				if cfg := props.AgentConfiguration; cfg != nil {
					cpu := 0
					if cfg.CPU != nil {
						cpu = int(*cfg.CPU)
					}
					agentConfig = []AgentConfig{{CPU: cpu}}
				}
				if props.AgentPoolName != nil {
					agentPoolName = *props.AgentPoolName
				}
				if props.IsSystemTask != nil {
					isSystemTask = *props.IsSystemTask
				}
				if props.LogTemplate != nil {
					logTemplate = *props.LogTemplate
				}
				platform = flattenRegistryTaskPlatform(props.Platform)
				enabled = props.Status == legacyacr.TaskStatusEnabled
				if props.Timeout != nil {
					timeoutInSec = int(*props.Timeout)
				}
				dockerStep = flattenRegistryTaskDockerStep(props.Step)
				fileTaskStep = flattenRegistryTaskFileTaskStep(props.Step)
				encodedTaskStep = flattenRegistryTaskEncodedTaskStep(props.Step)
				if trigger := props.Trigger; trigger != nil {
					baseImageTrigger = flattenRegistryTaskBaseImageTrigger(trigger.BaseImageTrigger)
					sourceTrigger = flattenRegistryTaskSourceTriggers(trigger.SourceTriggers)
					timerTrigger = flattenRegistryTaskTimerTriggers(trigger.TimerTriggers)
				}
			}
			model := ContainerRegistryTaskModel{
				Name:                id.TaskName,
				ContainerRegistryId: registryId.ID(),
				Identity:            flattenRegistryTaskIdentity(task.Identity),
				AgentConfig:         agentConfig,
				AgentPoolName:       agentPoolName,
				IsSystemTask:        isSystemTask,
				LogTemplate:         logTemplate,
				Platform:            platform,
				Enabled:             enabled,
				TimeoutInSec:        timeoutInSec,
				DockerStep:          dockerStep,
				FileTaskStep:        fileTaskStep,
				EncodedTaskStep:     encodedTaskStep,
				BaseImageTrigger:    baseImageTrigger,
				SourceTrigger:       sourceTrigger,
				TimerTrigger:        timerTrigger,
				Tags:                tags.Flatten(task.Tags),
			}

			return metadata.Encode(&model)
		},
	}
}

func (r ContainerRegistryTaskResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.TasksClient

			id, err := parse.ContainerRegistryTaskID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.RegistryName, id.TaskName)
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

//func (r ContainerRegistryTaskResource) Update() sdk.ResourceFunc {
//	return sdk.ResourceFunc{
//		Timeout: 30 * time.Minute,
//		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
//			id, err := parse.ContainerRegistryTaskID(metadata.ResourceData.Id())
//			if err != nil {
//				return err
//			}
//
//			var state ContainerRegistryTaskModel
//			if err := metadata.Decode(&state); err != nil {
//				return err
//			}
//
//			client := metadata.Client.Containers.TasksClient
//
//			// TODO: define the SDK patch object
//			// patch := {}
//			// if metadata.ResourceData.HasChange(?) {
//			// 	patch.? = ?
//			// }
//
//			future, err = client.Update(ctx, id.ResourceGroup, id.Name, patch)
//			if err != nil {
//				return fmt.Errorf("updating %s: %+v", id, err)
//			}
//			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
//				return fmt.Errorf("waiting for update of %s: %+v", id, err)
//			}
//
//			return nil
//		},
//	}
//}

func expandRegistryTaskTrigger(model ContainerRegistryTaskModel) *legacyacr.TriggerProperties {
	return &legacyacr.TriggerProperties{
		BaseImageTrigger: expandRegistryTaskBaseImageTrigger(model.BaseImageTrigger),
		SourceTriggers:   expandRegistryTaskSourceTriggers(model.SourceTrigger),
		TimerTriggers:    expandRegistryTaskTimerTriggers(model.TimerTrigger),
	}
}

func expandRegistryTaskBaseImageTrigger(triggers []BaseImageTrigger) *legacyacr.BaseImageTrigger {
	if len(triggers) == 0 {
		return nil
	}

	trigger := triggers[0]
	status := legacyacr.TriggerStatusDisabled
	if trigger.Enabled {
		status = legacyacr.TriggerStatusEnabled
	}
	out := &legacyacr.BaseImageTrigger{
		Name:                 &trigger.Name,
		BaseImageTriggerType: legacyacr.BaseImageTriggerType(trigger.Type),
		Status:               status,
	}
	if trigger.UpdateTriggerEndpoint != "" {
		out.UpdateTriggerEndpoint = &trigger.UpdateTriggerEndpoint
	}
	if trigger.UpdateTriggerPayloadType != "" {
		out.UpdateTriggerPayloadType = legacyacr.UpdateTriggerPayloadType(trigger.UpdateTriggerPayloadType)
	}
	return out
}

func flattenRegistryTaskBaseImageTrigger(trigger *legacyacr.BaseImageTrigger) []BaseImageTrigger {
	if trigger == nil {
		return nil
	}

	obj := BaseImageTrigger{
		Type:                     string(trigger.BaseImageTriggerType),
		Enabled:                  trigger.Status == legacyacr.TriggerStatusEnabled,
		UpdateTriggerPayloadType: string(trigger.UpdateTriggerPayloadType),
	}

	if trigger.Name != nil {
		obj.Name = *trigger.Name
	}
	if trigger.UpdateTriggerEndpoint != nil {
		obj.UpdateTriggerEndpoint = *trigger.UpdateTriggerEndpoint
	}
	return []BaseImageTrigger{obj}
}

func expandRegistryTaskSourceTriggers(triggers []SourceTrigger) *[]legacyacr.SourceTrigger {
	if len(triggers) == 0 {
		return nil
	}
	out := make([]legacyacr.SourceTrigger, 0, len(triggers))
	for _, trigger := range triggers {
		status := legacyacr.TriggerStatusDisabled
		if trigger.Enabled {
			status = legacyacr.TriggerStatusEnabled
		}
		sourceTrigger := legacyacr.SourceTrigger{
			Name:             &trigger.Name,
			Status:           status,
			SourceRepository: expandRegistryTaskSourceRepository(trigger.SourceSetting[0]),
		}
		if len(trigger.Events) != 0 {
			events := make([]legacyacr.SourceTriggerEvent, 0, len(trigger.Events))
			for _, event := range trigger.Events {
				events = append(events, legacyacr.SourceTriggerEvent(event))
			}
			sourceTrigger.SourceTriggerEvents = &events
		}
		out = append(out, sourceTrigger)
	}
	return &out
}

func flattenRegistryTaskSourceTriggers(triggers *[]legacyacr.SourceTrigger) []SourceTrigger {
	if triggers == nil {
		return nil
	}
	out := make([]SourceTrigger, 0, len(*triggers))
	for _, trigger := range *triggers {
		obj := SourceTrigger{
			Enabled:       trigger.Status == legacyacr.TriggerStatusEnabled,
			SourceSetting: flattenRegistryTaskSourceRepository(trigger.SourceRepository),
		}
		if trigger.Name != nil {
			obj.Name = *trigger.Name
		}
		if trigger.SourceTriggerEvents != nil {
			events := make([]string, 0, len(*trigger.SourceTriggerEvents))
			for _, event := range *trigger.SourceTriggerEvents {
				events = append(events, string(event))
			}
			obj.Events = events
		}
		out = append(out, obj)
	}
	return out
}

func expandRegistryTaskSourceRepository(setting SourceSetting) *legacyacr.SourceProperties {
	out := legacyacr.SourceProperties{
		SourceControlType: legacyacr.SourceControlType(setting.SourceType),
		RepositoryURL:     &setting.RepositoryURL,
	}
	if setting.Branch != "" {
		out.Branch = &setting.Branch
	}
	if len(setting.Auth) != 0 {
		out.SourceControlAuthProperties = expandRegistryTaskAuthInfo(setting.Auth[0])
	}
	return &out
}

func flattenRegistryTaskSourceRepository(repository *legacyacr.SourceProperties) []SourceSetting {
	if repository == nil {
		return nil
	}
	obj := SourceSetting{
		SourceType: string(repository.SourceControlType),
		Auth:       flattenRegistryTaskAuthInfo(repository.SourceControlAuthProperties),
	}
	if repository.RepositoryURL != nil {
		obj.RepositoryURL = *repository.RepositoryURL
	}
	if repository.Branch != nil {
		obj.Branch = *repository.Branch
	}
	return []SourceSetting{obj}
}

func expandRegistryTaskAuthInfo(auth Auth) *legacyacr.AuthInfo {
	out := legacyacr.AuthInfo{
		TokenType: legacyacr.TokenType(auth.TokenType),
		Token:     &auth.Token,
	}
	if auth.RefreshToken != "" {
		out.RefreshToken = &auth.RefreshToken
	}
	if auth.Scope != "" {
		out.Scope = &auth.Scope
	}
	if auth.ExpireInSec != 0 {
		out.ExpiresIn = utils.Int32(int32(auth.ExpireInSec))
	}
	return &out
}

func flattenRegistryTaskAuthInfo(auth *legacyacr.AuthInfo) []Auth {
	if auth == nil {
		return nil
	}
	obj := Auth{
		TokenType: string(auth.TokenType),
	}
	if auth.Token != nil {
		obj.Token = *auth.Token
	}
	if auth.RefreshToken != nil {
		obj.RefreshToken = *auth.RefreshToken
	}
	if auth.Scope != nil {
		obj.Scope = *auth.Scope
	}
	if auth.ExpiresIn != nil {
		obj.ExpireInSec = int(*auth.ExpiresIn)
	}
	return []Auth{obj}
}

func expandRegistryTaskTimerTriggers(triggers []TimerTrigger) *[]legacyacr.TimerTrigger {
	if len(triggers) == 0 {
		return nil
	}
	out := make([]legacyacr.TimerTrigger, 0, len(triggers))
	for _, trigger := range triggers {
		status := legacyacr.TriggerStatusDisabled
		if trigger.Enabled {
			status = legacyacr.TriggerStatusEnabled
		}
		timerTrigger := legacyacr.TimerTrigger{
			Name:     &trigger.Name,
			Schedule: &trigger.Schedule,
			Status:   status,
		}
		out = append(out, timerTrigger)
	}
	return &out
}

func flattenRegistryTaskTimerTriggers(triggers *[]legacyacr.TimerTrigger) []TimerTrigger {
	if triggers == nil {
		return nil
	}
	out := make([]TimerTrigger, 0, len(*triggers))
	for _, trigger := range *triggers {
		obj := TimerTrigger{
			Enabled: trigger.Status == legacyacr.TriggerStatusEnabled,
		}
		if trigger.Name != nil {
			obj.Name = *trigger.Name
		}
		if trigger.Schedule != nil {
			obj.Schedule = *trigger.Schedule
		}
		out = append(out, obj)
	}
	return out
}

func expandRegistryTaskStep(model ContainerRegistryTaskModel) legacyacr.BasicTaskStepProperties {
	switch {
	case len(model.DockerStep) != 0:
		return expandRegistryTaskDockerStep(model.DockerStep[0])
	case len(model.FileTaskStep) != 0:
		return expandRegistryTaskFileTaskStep(model.FileTaskStep[0])
	case len(model.EncodedTaskStep) != 0:
		return expandRegistryTaskEncodedTaskStep(model.EncodedTaskStep[0])
	}
	return nil
}

func expandRegistryTaskDockerStep(step DockerStep) legacyacr.DockerBuildStep {
	out := legacyacr.DockerBuildStep{
		Type:           legacyacr.TypeDocker,
		DockerFilePath: &step.DockerfilePath,
		IsPushEnabled:  &step.IsPushEnabled,
		NoCache:        utils.Bool(!step.IsCacheEnabled),
	}
	if step.ContextPath != "" {
		out.ContextPath = &step.ContextPath
	}
	if step.ContextAccessToken != "" {
		out.ContextAccessToken = &step.ContextAccessToken
	}
	if len(step.ImageNames) != 0 {
		out.ImageNames = &step.ImageNames
	}
	if step.Target != "" {
		out.Target = &step.Target
	}
	if len(step.Arguments) != 0 {
		out.Arguments = expandRegistryTaskArguments(step.Arguments)
	}
	return out
}

func flattenRegistryTaskDockerStep(step legacyacr.BasicTaskStepProperties) []DockerStep {
	dockerStep, ok := step.AsDockerBuildStep()
	if !ok {
		return nil
	}

	obj := DockerStep{
		Arguments: flattenRegistryTaskArguments(dockerStep.Arguments),
	}

	if dockerStep.ContextPath != nil {
		obj.ContextPath = *dockerStep.ContextPath
	}
	if dockerStep.ContextAccessToken != nil {
		obj.ContextAccessToken = *dockerStep.ContextAccessToken
	}
	if dockerStep.DockerFilePath != nil {
		obj.DockerfilePath = *dockerStep.DockerFilePath
	}
	if dockerStep.ImageNames != nil {
		obj.ImageNames = *dockerStep.ImageNames
	}
	if dockerStep.IsPushEnabled != nil {
		obj.IsPushEnabled = *dockerStep.IsPushEnabled
	}
	if dockerStep.NoCache != nil {
		obj.IsCacheEnabled = !(*dockerStep.NoCache)
	}
	if dockerStep.Target != nil {
		obj.Target = *dockerStep.Target
	}
	return []DockerStep{obj}
}

func expandRegistryTaskFileTaskStep(step FileTaskStep) legacyacr.FileTaskStep {
	out := legacyacr.FileTaskStep{
		Type:         legacyacr.TypeFileTask,
		TaskFilePath: &step.TaskFilePath,
	}
	if step.ValueFilePath != "" {
		out.ValuesFilePath = &step.ValueFilePath
	}
	if step.ContextPath != "" {
		out.ContextPath = &step.ContextPath
	}
	if step.ContextAccessToken != "" {
		out.ContextAccessToken = &step.ContextAccessToken
	}
	if len(step.Values) != 0 {
		out.Values = expandRegistryTaskValues(step.Values)
	}
	return out
}

func flattenRegistryTaskFileTaskStep(step legacyacr.BasicTaskStepProperties) []FileTaskStep {
	fileTaskStep, ok := step.AsFileTaskStep()
	if !ok {
		return nil
	}
	obj := FileTaskStep{
		Values: flattenRegistryTaskValues(fileTaskStep.Values),
	}
	if fileTaskStep.ContextPath != nil {
		obj.ContextPath = *fileTaskStep.ContextPath
	}
	if fileTaskStep.ContextAccessToken != nil {
		obj.ContextAccessToken = *fileTaskStep.ContextAccessToken
	}
	if fileTaskStep.TaskFilePath != nil {
		obj.TaskFilePath = *fileTaskStep.TaskFilePath
	}
	if fileTaskStep.ValuesFilePath != nil {
		obj.ValueFilePath = *fileTaskStep.ValuesFilePath
	}
	return []FileTaskStep{obj}
}

func expandRegistryTaskEncodedTaskStep(step EncodedTaskStep) legacyacr.EncodedTaskStep {
	out := legacyacr.EncodedTaskStep{
		Type:               legacyacr.TypeEncodedTask,
		EncodedTaskContent: utils.String(utils.Base64EncodeIfNot(step.TaskContent)),
	}
	if step.ContextPath != "" {
		out.ContextPath = &step.ContextPath
	}
	if step.ContextAccessToken != "" {
		out.ContextAccessToken = &step.ContextAccessToken
	}
	if step.ValueContent != "" {
		out.EncodedValuesContent = utils.String(utils.Base64EncodeIfNot(step.ValueContent))
	}
	if len(step.Values) != 0 {
		out.Values = expandRegistryTaskValues(step.Values)
	}
	return out
}

func flattenRegistryTaskEncodedTaskStep(step legacyacr.BasicTaskStepProperties) []EncodedTaskStep {
	encodedTaskStep, ok := step.AsEncodedTaskStep()
	if !ok {
		return nil
	}

	obj := EncodedTaskStep{
		Values: flattenRegistryTaskValues(encodedTaskStep.Values),
	}

	if encodedTaskStep.ContextPath != nil {
		obj.ContextPath = *encodedTaskStep.ContextPath
	}
	if encodedTaskStep.ContextAccessToken != nil {
		obj.ContextAccessToken = *encodedTaskStep.ContextAccessToken
	}
	if encodedTaskStep.EncodedTaskContent != nil {
		obj.TaskContent = *encodedTaskStep.EncodedTaskContent
	}
	if encodedTaskStep.EncodedValuesContent != nil {
		obj.ValueContent = *encodedTaskStep.EncodedValuesContent
	}

	return []EncodedTaskStep{obj}
}

func expandRegistryTaskArguments(arguments []Argument) *[]legacyacr.Argument {
	if len(arguments) == 0 {
		return nil
	}
	out := make([]legacyacr.Argument, 0, len(arguments))
	for _, argument := range arguments {
		out = append(out, legacyacr.Argument{
			Name:     &argument.Name,
			Value:    &argument.Value,
			IsSecret: &argument.IsSecret,
		})
	}
	return &out
}

func flattenRegistryTaskArguments(arguments *[]legacyacr.Argument) []Argument {
	if arguments == nil {
		return nil
	}
	out := make([]Argument, 0, len(*arguments))
	for _, argument := range *arguments {
		obj := Argument{}
		if argument.Name != nil {
			obj.Name = *argument.Name
		}
		if argument.Value != nil {
			obj.Value = *argument.Value
		}
		if argument.IsSecret != nil {
			obj.IsSecret = *argument.IsSecret
		}
		out = append(out, obj)
	}
	return out
}

func expandRegistryTaskValues(values []Value) *[]legacyacr.SetValue {
	if len(values) == 0 {
		return nil
	}
	out := make([]legacyacr.SetValue, 0, len(values))
	for _, value := range values {
		out = append(out, legacyacr.SetValue{
			Name:     &value.Name,
			Value:    &value.Value,
			IsSecret: &value.IsSecret,
		})
	}
	return &out
}

func flattenRegistryTaskValues(values *[]legacyacr.SetValue) []Value {
	if values == nil {
		return nil
	}
	out := make([]Value, 0, len(*values))
	for _, value := range *values {
		obj := Value{}
		if value.Name != nil {
			obj.Name = *value.Name
		}
		if value.Value != nil {
			obj.Value = *value.Value
		}
		if value.IsSecret != nil {
			obj.IsSecret = *value.IsSecret
		}
		out = append(out, obj)
	}
	return out
}

func expandRegistryTaskIdentity(identities []Identity) *legacyacr.IdentityProperties {
	if len(identities) == 0 {
		return nil
	}

	identity := identities[0]
	out := &legacyacr.IdentityProperties{
		Type: legacyacr.ResourceIdentityType(identity.Type),
	}
	if len(identity.IdentityIds) > 0 {
		out.UserAssignedIdentities = map[string]*legacyacr.UserIdentityProperties{}
		for _, identityId := range identity.IdentityIds {
			out.UserAssignedIdentities[identityId] = &legacyacr.UserIdentityProperties{}
		}
	}
	return out
}

func flattenRegistryTaskIdentity(identity *legacyacr.IdentityProperties) []Identity {
	if identity == nil {
		return nil
	}

	obj := Identity{
		Type: string(identity.Type),
	}

	if identity.PrincipalID != nil {
		obj.PrincipalId = *identity.PrincipalID
	}
	if identity.TenantID != nil {
		obj.TenantId = *identity.TenantID
	}

	var identityIds []string
	for id := range identity.UserAssignedIdentities {
		identityIds = append(identityIds, id)
	}
	obj.IdentityIds = identityIds

	return []Identity{obj}
}

func expandRegistryTaskPlatform(platform Platform) *legacyacr.PlatformProperties {
	out := &legacyacr.PlatformProperties{
		Os: legacyacr.OS(platform.OS),
	}
	if arch := platform.Architecture; arch != "" {
		out.Architecture = legacyacr.Architecture(arch)
	}
	if variant := platform.Variant; variant != "" {
		out.Variant = legacyacr.Variant(variant)
	}
	return out
}

func flattenRegistryTaskPlatform(platform *legacyacr.PlatformProperties) []Platform {
	if platform == nil {
		return nil
	}
	return []Platform{{
		OS:           string(platform.Os),
		Architecture: string(platform.Architecture),
		Variant:      string(platform.Variant),
	}}
}
