// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/tasks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/registries"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContainerRegistryTaskResource struct{}

var (
	_ sdk.ResourceWithUpdate        = ContainerRegistryTaskResource{}
	_ sdk.ResourceWithCustomizeDiff = ContainerRegistryTaskResource{}
)

type AgentConfig struct {
	CPU int `tfschema:"cpu"`
}

type Platform struct {
	OS           string `tfschema:"os"`
	Architecture string `tfschema:"architecture"`
	Variant      string `tfschema:"variant"`
}

type DockerStep struct {
	ContextPath        string            `tfschema:"context_path"`
	ContextAccessToken string            `tfschema:"context_access_token"`
	DockerfilePath     string            `tfschema:"dockerfile_path"`
	ImageNames         []string          `tfschema:"image_names"`
	IsPushEnabled      bool              `tfschema:"push_enabled"`
	IsCacheEnabled     bool              `tfschema:"cache_enabled"`
	Target             string            `tfschema:"target"`
	Arguments          map[string]string `tfschema:"arguments"`
	SecretArguments    map[string]string `tfschema:"secret_arguments"`
}

type FileTaskStep struct {
	ContextPath        string            `tfschema:"context_path"`
	ContextAccessToken string            `tfschema:"context_access_token"`
	TaskFilePath       string            `tfschema:"task_file_path"`
	ValueFilePath      string            `tfschema:"value_file_path"`
	Values             map[string]string `tfschema:"values"`
	SecretValues       map[string]string `tfschema:"secret_values"`
}

type EncodedTaskStep struct {
	ContextPath        string            `tfschema:"context_path"`
	ContextAccessToken string            `tfschema:"context_access_token"`
	TaskContent        string            `tfschema:"task_content"`
	ValueContent       string            `tfschema:"value_content"`
	Values             map[string]string `tfschema:"values"`
	SecretValues       map[string]string `tfschema:"secret_values"`
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

type SourceSetting struct{}

type SourceTrigger struct {
	Name          string   `tfschema:"name"`
	Enabled       bool     `tfschema:"enabled"`
	Events        []string `tfschema:"events"`
	SourceType    string   `tfschema:"source_type"`
	RepositoryURL string   `tfschema:"repository_url"`
	Branch        string   `tfschema:"branch"`
	Auth          []Auth   `tfschema:"authentication"`
}

type TimerTrigger struct {
	Name     string `tfschema:"name"`
	Enabled  bool   `tfschema:"enabled"`
	Schedule string `tfschema:"schedule"`
}

type RegistryCredential struct {
	Source []SourceRegistryCredential `tfschema:"source"`
	Custom []CustomRegistryCredential `tfschema:"custom"`
}

type SourceRegistryCredential struct {
	LoginMode string `tfschema:"login_mode"`
}

type CustomRegistryCredential struct {
	LoginServer string `tfschema:"login_server"`
	UserName    string `tfschema:"username"`
	Password    string `tfschema:"password"`
	Identity    string `tfschema:"identity"`
}

type ContainerRegistryTaskModel struct {
	Name                string               `tfschema:"name"`
	ContainerRegistryId string               `tfschema:"container_registry_id"`
	AgentConfig         []AgentConfig        `tfschema:"agent_setting"`
	AgentPoolName       string               `tfschema:"agent_pool_name"`
	IsSystemTask        bool                 `tfschema:"is_system_task"`
	LogTemplate         string               `tfschema:"log_template"`
	Platform            []Platform           `tfschema:"platform"`
	Enabled             bool                 `tfschema:"enabled"`
	TimeoutInSec        int                  `tfschema:"timeout_in_seconds"`
	DockerStep          []DockerStep         `tfschema:"docker_step"`
	FileTaskStep        []FileTaskStep       `tfschema:"file_step"`
	EncodedTaskStep     []EncodedTaskStep    `tfschema:"encoded_step"`
	BaseImageTrigger    []BaseImageTrigger   `tfschema:"base_image_trigger"`
	SourceTrigger       []SourceTrigger      `tfschema:"source_trigger"`
	TimerTrigger        []TimerTrigger       `tfschema:"timer_trigger"`
	RegistryCredential  []RegistryCredential `tfschema:"registry_credential"`
	Tags                map[string]string    `tfschema:"tags"`
}

func userDataStateFunc(v interface{}) string {
	switch s := v.(type) {
	case string:
		return utils.Base64EncodeIfNot(s)
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
			ValidateFunc: registries.ValidateRegistryID,
		},
		"platform": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"os": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(tasks.OSWindows),
							string(tasks.OSLinux),
						}, false),
					},
					"architecture": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(tasks.ArchitectureAmdSixFour),
							string(tasks.ArchitectureArm),
							string(tasks.ArchitectureArmSixFour),
							string(tasks.ArchitectureThreeEightSix),
							string(tasks.ArchitectureXEightSix),
						}, false),
					},
					"variant": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(tasks.VariantVSix),
							string(tasks.VariantVSeven),
							string(tasks.VariantVEight),
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
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"context_access_token": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"image_names": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"push_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"cache_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"target": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"arguments": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"secret_arguments": {
						Type:      pluginsdk.TypeMap,
						Optional:  true,
						Sensitive: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
			ConflictsWith: []string{"file_step", "encoded_step"},
		},
		"file_step": {
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
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"values": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"secret_values": {
						Type:      pluginsdk.TypeMap,
						Optional:  true,
						Sensitive: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
			ConflictsWith: []string{"docker_step", "encoded_step"},
		},
		"encoded_step": {
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
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"values": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"secret_values": {
						Type:      pluginsdk.TypeMap,
						Optional:  true,
						Sensitive: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
			ConflictsWith: []string{"docker_step", "file_step"},
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
							string(tasks.BaseImageTriggerTypeAll),
							string(tasks.BaseImageTriggerTypeRuntime),
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
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"update_trigger_payload_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(tasks.UpdateTriggerPayloadTypeDefault),
							string(tasks.UpdateTriggerPayloadTypeToken),
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
					"events": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								string(tasks.SourceTriggerEventCommit),
								string(tasks.SourceTriggerEventPullrequest),
							}, false),
						},
					},
					"source_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(tasks.SourceControlTypeGithub),
							string(tasks.SourceControlTypeVisualStudioTeamService),
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
					"authentication": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"token_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(tasks.TokenTypePAT),
										string(tasks.TokenTypeOAuth),
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
									Sensitive:    true,
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
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
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
		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"registry_credential": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"source": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"login_mode": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(tasks.SourceRegistryLoginModeNone),
										string(tasks.SourceRegistryLoginModeDefault),
									}, false),
								},
							},
						},
						AtLeastOneOf: []string{"registry_credential.0.source", "registry_credential.0.custom"},
					},
					"custom": {
						Type:      pluginsdk.TypeSet,
						Sensitive: true,
						Optional:  true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"login_server": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"username": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"password": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"identity": {
									// TODO - 4.0: this should be `user_assigned_identity_id`?
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
						AtLeastOneOf: []string{"registry_credential.0.source", "registry_credential.0.custom"},
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
			ForceNew: true,
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
		"tags": commonschema.Tags(),
	}
}

func (r ContainerRegistryTaskResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff
			isSystemTask := rd.Get("is_system_task").(bool)

			if isSystemTask {
				invalidProps := []string{"platform", "docker_step", "file_step", "encoded_step", "base_image_trigger", "source_trigger", "timer_trigger"}
				for _, prop := range invalidProps {
					if v := rd.Get(prop).([]interface{}); len(v) != 0 {
						return fmt.Errorf("system task can't specify `%s`", prop)
					}
				}
			} else {
				if v := rd.Get("platform").([]interface{}); len(v) == 0 {
					return fmt.Errorf("non-system task have to specify `platform`")
				}

				dockerStep := rd.Get("docker_step").([]interface{})
				fileTaskStep := rd.Get("file_step").([]interface{})
				encodedTaskStep := rd.Get("encoded_step").([]interface{})
				if len(dockerStep)+len(fileTaskStep)+len(encodedTaskStep) == 0 {
					return fmt.Errorf("non-system task have to specify one of `docker_step`, `file_step` and `encoded_step`")
				}
			}

			return nil
		},
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
	return tasks.ValidateTaskID
}

func (r ContainerRegistryTaskResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.ContainerRegistryClient_v2019_06_01_preview.Tasks
			registryClient := metadata.Client.Containers.ContainerRegistryClient_v2021_08_01_preview.Registries

			var model ContainerRegistryTaskModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			registryId, err := registries.ParseRegistryID(model.ContainerRegistryId)
			if err != nil {
				return err
			}

			registry, err := registryClient.Get(ctx, *registryId)
			if err != nil {
				return fmt.Errorf("getting registry %s: %+v", registryId, err)
			}

			id := tasks.NewTaskID(registryId.SubscriptionId, registryId.ResourceGroupName, registryId.RegistryName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			status := tasks.TaskStatusDisabled
			if model.Enabled {
				status = tasks.TaskStatusEnabled
			}

			expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			params := tasks.Task{
				Properties: &tasks.TaskProperties{
					Platform:           expandRegistryTaskPlatform(model.Platform),
					Step:               expandRegistryTaskStep(model),
					Trigger:            expandRegistryTaskTrigger(model),
					Status:             pointer.To(status),
					IsSystemTask:       &model.IsSystemTask,
					Timeout:            utils.Int64(int64(model.TimeoutInSec)),
					Credentials:        expandRegistryTaskCredentials(model.RegistryCredential),
					AgentConfiguration: expandRegistryTaskAgentProperties(model.AgentConfig),
				},

				// The location of the task must be the same as the registry, otherwise the API will raise error complaining can't find the registry.
				Location: location.Normalize(registry.Model.Location),
				Identity: expandedIdentity,
				Tags:     &model.Tags,
			}

			if model.AgentPoolName != "" {
				params.Properties.AgentPoolName = &model.AgentPoolName
			}
			if model.LogTemplate != "" {
				params.Properties.LogTemplate = &model.LogTemplate
			}

			if err := client.CreateThenPoll(ctx, id, params); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
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
			client := metadata.Client.Containers.ContainerRegistryClient_v2019_06_01_preview.Tasks
			id, err := tasks.ParseTaskID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			registryId := registries.NewRegistryID(id.SubscriptionId, id.ResourceGroupName, id.RegistryName)

			task, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(task.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			var diffOrStateModel ContainerRegistryTaskModel
			if err := metadata.Decode(&diffOrStateModel); err != nil {
				return fmt.Errorf("decoding model from diff/state %+v", err)
			}

			var (
				agentConfig        []AgentConfig
				agentPoolName      string
				isSystemTask       bool
				logTemplate        string
				platform           []Platform
				enabled            bool
				timeoutInSec       int
				dockerStep         []DockerStep
				fileTaskStep       []FileTaskStep
				encodedTaskStep    []EncodedTaskStep
				baseImageTrigger   []BaseImageTrigger
				sourceTrigger      []SourceTrigger
				timerTrigger       []TimerTrigger
				registryCredential []RegistryCredential
				tag                map[string]string
			)

			if model := task.Model; model != nil {
				if model.Tags != nil {
					tag = *model.Tags
				}

				flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				if err := metadata.ResourceData.Set("identity", flattenedIdentity); err != nil {
					return fmt.Errorf("setting `identity`: %+v", err)
				}

				if props := model.Properties; props != nil {
					agentConfig = flattenRegistryTaskAgentProperties(props.AgentConfiguration)
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
					enabled = *props.Status == tasks.TaskStatusEnabled
					if props.Timeout != nil {
						timeoutInSec = int(*props.Timeout)
					}
					dockerStep = flattenRegistryTaskDockerStep(props.Step, diffOrStateModel)
					fileTaskStep = flattenRegistryTaskFileTaskStep(props.Step, diffOrStateModel)
					encodedTaskStep = flattenRegistryTaskEncodedTaskStep(props.Step, diffOrStateModel)
					if trigger := props.Trigger; trigger != nil {
						baseImageTrigger = flattenRegistryTaskBaseImageTrigger(trigger.BaseImageTrigger, diffOrStateModel)
						sourceTrigger = flattenRegistryTaskSourceTriggers(trigger.SourceTriggers, diffOrStateModel)
						timerTrigger = flattenRegistryTaskTimerTriggers(trigger.TimerTriggers)
					}
					registryCredential = flattenRegistryTaskCredentials(props.Credentials, diffOrStateModel)
				}
			}

			model := ContainerRegistryTaskModel{
				Name:                id.TaskName,
				ContainerRegistryId: registryId.ID(),
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
				RegistryCredential:  registryCredential,
				Tags:                tag,
			}

			if err := metadata.Encode(&model); err != nil {
				return fmt.Errorf("encoding: %+v", err)
			}

			return nil
		},
	}
}

func (r ContainerRegistryTaskResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.ContainerRegistryClient_v2019_06_01_preview.Tasks

			id, err := tasks.ParseTaskID(metadata.ResourceData.Id())
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

func (r ContainerRegistryTaskResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := tasks.ParseTaskID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Containers.ContainerRegistryClient_v2019_06_01_preview.Tasks
			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: +%v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("model is nil for %s", *id)
			}

			var model ContainerRegistryTaskModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("platform") {
				existing.Model.Properties.Platform = expandRegistryTaskPlatform(model.Platform)
			}
			if metadata.ResourceData.HasChange("docker_step") || metadata.ResourceData.HasChange("file_step") || metadata.ResourceData.HasChange("encoded_step") {
				existing.Model.Properties.Step = expandRegistryTaskStep(model)
			}

			if metadata.ResourceData.HasChange("base_image_trigger") || metadata.ResourceData.HasChange("source_trigger") || metadata.ResourceData.HasChange("timer_trigger") {
				existing.Model.Properties.Trigger = expandRegistryTaskTrigger(model)
			}

			if existing.Model.Properties.Trigger != nil {
				if !metadata.ResourceData.HasChange("source_triggers") && existing.Model.Properties.Trigger.SourceTriggers != nil {
					// For update that is not affecting source_triggers, we need to patch the source_triggers to include the properties missing in the response of GET.
					existing.Model.Properties.Trigger.SourceTriggers = patchRegistryTaskTriggerSourceTrigger(*existing.Model.Properties.Trigger.SourceTriggers, model)
				}
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				existing.Model.Identity = expandedIdentity
			}

			// Deliberately always set "registry_credential" as the custom registry's credentials are not returned by API, but are required for a PUT request.
			existing.Model.Properties.Credentials = expandRegistryTaskCredentials(model.RegistryCredential)

			if metadata.ResourceData.HasChange("agent_setting") {
				existing.Model.Properties.AgentConfiguration = expandRegistryTaskAgentProperties(model.AgentConfig)
			}
			if metadata.ResourceData.HasChange("agent_pool_name") && model.AgentPoolName != "" {
				existing.Model.Properties.AgentPoolName = &model.AgentPoolName
			}
			if metadata.ResourceData.HasChange("enabled") {
				status := tasks.TaskStatusDisabled
				if model.Enabled {
					status = tasks.TaskStatusEnabled
				}
				existing.Model.Properties.Status = pointer.To(status)
			}
			if metadata.ResourceData.HasChange("log_template") && model.LogTemplate != "" {
				existing.Model.Properties.LogTemplate = &model.LogTemplate
			}
			if metadata.ResourceData.HasChange("timeout_in_seconds") {
				existing.Model.Properties.Timeout = utils.Int64(int64(model.TimeoutInSec))
			}
			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = &model.Tags
			}

			// Due to the fact that the service doesn't honor explicitly set to null fields in the PATCH request,
			// we can not use PATCH (i.e. the Update) here.
			if err := client.CreateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandRegistryTaskTrigger(model ContainerRegistryTaskModel) *tasks.TriggerProperties {
	baseImageTrigger := expandRegistryTaskBaseImageTrigger(model.BaseImageTrigger)
	sourceTriggers := expandRegistryTaskSourceTriggers(model.SourceTrigger)
	timerTriggers := expandRegistryTaskTimerTriggers(model.TimerTrigger)
	if baseImageTrigger == nil && sourceTriggers == nil && timerTriggers == nil {
		return nil
	}
	return &tasks.TriggerProperties{
		BaseImageTrigger: baseImageTrigger,
		SourceTriggers:   sourceTriggers,
		TimerTriggers:    timerTriggers,
	}
}

func expandRegistryTaskBaseImageTrigger(triggers []BaseImageTrigger) *tasks.BaseImageTrigger {
	if len(triggers) == 0 {
		return nil
	}

	trigger := triggers[0]
	status := tasks.TriggerStatusDisabled
	if trigger.Enabled {
		status = tasks.TriggerStatusEnabled
	}
	out := &tasks.BaseImageTrigger{
		Name:                 trigger.Name,
		BaseImageTriggerType: tasks.BaseImageTriggerType(trigger.Type),
		Status:               &status,
	}
	if trigger.UpdateTriggerEndpoint != "" {
		out.UpdateTriggerEndpoint = &trigger.UpdateTriggerEndpoint
	}
	if trigger.UpdateTriggerPayloadType != "" {
		out.UpdateTriggerPayloadType = pointer.To(tasks.UpdateTriggerPayloadType(trigger.UpdateTriggerPayloadType))
	}
	return out
}

func flattenRegistryTaskBaseImageTrigger(trigger *tasks.BaseImageTrigger, model ContainerRegistryTaskModel) []BaseImageTrigger {
	if trigger == nil {
		return nil
	}

	payloadType := ""
	if v := trigger.UpdateTriggerPayloadType; v != nil {
		payloadType = string(*v)
	}
	obj := BaseImageTrigger{
		Type:                     string(trigger.BaseImageTriggerType),
		Enabled:                  *trigger.Status == tasks.TriggerStatusEnabled,
		UpdateTriggerPayloadType: payloadType,
		Name:                     trigger.Name,
	}

	// UpdateTriggerEndpoint is not returned from API, setting it from config.
	if len(model.BaseImageTrigger) == 1 {
		obj.UpdateTriggerEndpoint = model.BaseImageTrigger[0].UpdateTriggerEndpoint
	}

	return []BaseImageTrigger{obj}
}

func expandRegistryTaskSourceTriggers(triggers []SourceTrigger) *[]tasks.SourceTrigger {
	if len(triggers) == 0 {
		return nil
	}
	out := make([]tasks.SourceTrigger, 0, len(triggers))
	for _, trigger := range triggers {
		status := tasks.TriggerStatusDisabled
		if trigger.Enabled {
			status = tasks.TriggerStatusEnabled
		}
		sourceTrigger := tasks.SourceTrigger{
			Name:   trigger.Name,
			Status: &status,
			SourceRepository: tasks.SourceProperties{
				SourceControlType: tasks.SourceControlType(trigger.SourceType),
				RepositoryUrl:     trigger.RepositoryURL,
			},
		}
		if len(trigger.Events) != 0 {
			events := make([]tasks.SourceTriggerEvent, 0, len(trigger.Events))
			for _, event := range trigger.Events {
				events = append(events, tasks.SourceTriggerEvent(event))
			}
			sourceTrigger.SourceTriggerEvents = events
		}

		if trigger.Branch != "" {
			sourceTrigger.SourceRepository.Branch = &trigger.Branch
		}
		if len(trigger.Auth) != 0 {
			sourceTrigger.SourceRepository.SourceControlAuthProperties = expandRegistryTaskAuthInfo(trigger.Auth[0])
		}
		out = append(out, sourceTrigger)
	}
	return &out
}

func flattenRegistryTaskSourceTriggers(triggers *[]tasks.SourceTrigger, model ContainerRegistryTaskModel) []SourceTrigger {
	if triggers == nil {
		return nil
	}
	out := make([]SourceTrigger, 0, len(*triggers))
	for i, trigger := range *triggers {
		obj := SourceTrigger{
			Enabled: *trigger.Status == tasks.TriggerStatusEnabled,
		}
		obj.Name = trigger.Name

		if trigger.SourceTriggerEvents != nil {
			events := make([]string, 0, len(trigger.SourceTriggerEvents))
			for _, event := range trigger.SourceTriggerEvents {
				events = append(events, string(event))
			}
			obj.Events = events
		}

		obj.SourceType = string(trigger.SourceRepository.SourceControlType)
		obj.RepositoryURL = trigger.SourceRepository.RepositoryUrl
		if trigger.SourceRepository.Branch != nil {
			obj.Branch = *trigger.SourceRepository.Branch
		}

		// Auth is not returned from API, setting it from config.
		if len(model.SourceTrigger) > i {
			obj.Auth = model.SourceTrigger[i].Auth
		}

		out = append(out, obj)
	}
	return out
}

func expandRegistryTaskAuthInfo(auth Auth) *tasks.AuthInfo {
	out := tasks.AuthInfo{
		TokenType: tasks.TokenType(auth.TokenType),
		Token:     auth.Token,
	}
	if auth.RefreshToken != "" {
		out.RefreshToken = &auth.RefreshToken
	}
	if auth.Scope != "" {
		out.Scope = &auth.Scope
	}
	if auth.ExpireInSec != 0 {
		out.ExpiresIn = utils.Int64(int64(auth.ExpireInSec))
	}
	return &out
}

func expandRegistryTaskTimerTriggers(triggers []TimerTrigger) *[]tasks.TimerTrigger {
	if len(triggers) == 0 {
		return nil
	}
	out := make([]tasks.TimerTrigger, 0, len(triggers))
	for _, trigger := range triggers {
		status := tasks.TriggerStatusDisabled
		if trigger.Enabled {
			status = tasks.TriggerStatusEnabled
		}
		timerTrigger := tasks.TimerTrigger{
			Name:     trigger.Name,
			Schedule: trigger.Schedule,
			Status:   &status,
		}
		out = append(out, timerTrigger)
	}
	return &out
}

func flattenRegistryTaskTimerTriggers(triggers *[]tasks.TimerTrigger) []TimerTrigger {
	if triggers == nil {
		return nil
	}
	out := make([]TimerTrigger, 0, len(*triggers))
	for _, trigger := range *triggers {
		obj := TimerTrigger{
			Enabled: *trigger.Status == tasks.TriggerStatusEnabled,
		}
		obj.Name = trigger.Name
		obj.Schedule = trigger.Schedule
		out = append(out, obj)
	}
	return out
}

func expandRegistryTaskStep(model ContainerRegistryTaskModel) tasks.TaskStepProperties {
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

func expandRegistryTaskDockerStep(step DockerStep) tasks.DockerBuildStep {
	out := tasks.DockerBuildStep{
		DockerFilePath: step.DockerfilePath,
		IsPushEnabled:  &step.IsPushEnabled,
		NoCache:        utils.Bool(!step.IsCacheEnabled),
		Arguments:      expandRegistryTaskArguments(step.Arguments, step.SecretArguments),
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

	return out
}

func flattenRegistryTaskDockerStep(step tasks.TaskStepProperties, model ContainerRegistryTaskModel) []DockerStep {
	if step == nil {
		return nil
	}

	dockerStep, ok := step.(tasks.DockerBuildStep)
	if !ok {
		return nil
	}

	obj := DockerStep{
		Arguments: flattenRegistryTaskArguments(dockerStep.Arguments),
	}

	if dockerStep.ContextPath != nil {
		obj.ContextPath = *dockerStep.ContextPath
	}
	obj.DockerfilePath = dockerStep.DockerFilePath

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

	if len(model.DockerStep) == 1 {
		// The ContextAccessToken is sensitive and won't return from API, setting it from the config.
		obj.ContextAccessToken = model.DockerStep[0].ContextAccessToken

		// The SecretArguments is sensitive and won't return from API, setting it from the config.
		obj.SecretArguments = model.DockerStep[0].SecretArguments
	}

	return []DockerStep{obj}
}

func expandRegistryTaskFileTaskStep(step FileTaskStep) tasks.FileTaskStep {
	out := tasks.FileTaskStep{
		TaskFilePath: step.TaskFilePath,
		Values:       expandRegistryTaskValues(step.Values, step.SecretValues),
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
	return out
}

func flattenRegistryTaskFileTaskStep(step tasks.TaskStepProperties, model ContainerRegistryTaskModel) []FileTaskStep {
	if step == nil {
		return nil
	}

	fileTaskStep, ok := step.(tasks.FileTaskStep)
	if !ok {
		return nil
	}

	obj := FileTaskStep{
		Values: flattenRegistryTaskValues(fileTaskStep.Values),
	}

	if fileTaskStep.ContextPath != nil {
		obj.ContextPath = *fileTaskStep.ContextPath
	}
	obj.TaskFilePath = fileTaskStep.TaskFilePath

	if fileTaskStep.ValuesFilePath != nil {
		obj.ValueFilePath = *fileTaskStep.ValuesFilePath
	}

	if len(model.FileTaskStep) == 1 {
		// The ContextAccessToken is sensitive and won't return from API, setting it from the config.
		obj.ContextAccessToken = model.FileTaskStep[0].ContextAccessToken

		// The SecretValues is sensitive and won't return from API, setting it from the config.
		obj.SecretValues = model.FileTaskStep[0].SecretValues
	}

	return []FileTaskStep{obj}
}

func expandRegistryTaskEncodedTaskStep(step EncodedTaskStep) tasks.EncodedTaskStep {
	out := tasks.EncodedTaskStep{
		EncodedTaskContent: utils.Base64EncodeIfNot(step.TaskContent),
		Values:             expandRegistryTaskValues(step.Values, step.SecretValues),
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
	return out
}

func flattenRegistryTaskEncodedTaskStep(step tasks.TaskStepProperties, model ContainerRegistryTaskModel) []EncodedTaskStep {
	if step == nil {
		return nil
	}

	encodedTaskStep, ok := step.(tasks.EncodedTaskStep)
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
	obj.TaskContent = encodedTaskStep.EncodedTaskContent

	if encodedTaskStep.EncodedValuesContent != nil {
		obj.ValueContent = *encodedTaskStep.EncodedValuesContent
	}

	if len(model.EncodedTaskStep) == 1 {
		// The ContextAccessToken is sensitive and won't return from API, setting it from the config.
		obj.ContextAccessToken = model.EncodedTaskStep[0].ContextAccessToken

		// The SecretValues is sensitive and won't return from API, setting it from the config.
		obj.SecretValues = model.EncodedTaskStep[0].SecretValues
	}

	return []EncodedTaskStep{obj}
}

func expandRegistryTaskArguments(arguments map[string]string, secretArguments map[string]string) *[]tasks.Argument {
	if len(arguments) == 0 && len(secretArguments) == 0 {
		return nil
	}
	out := make([]tasks.Argument, 0, len(arguments)+len(secretArguments))
	for k, v := range arguments {
		out = append(out, tasks.Argument{
			Name:     k,
			Value:    v,
			IsSecret: utils.Bool(false),
		})
	}
	for k, v := range secretArguments {
		out = append(out, tasks.Argument{
			Name:     k,
			Value:    v,
			IsSecret: utils.Bool(true),
		})
	}
	return &out
}

func flattenRegistryTaskArguments(arguments *[]tasks.Argument) map[string]string {
	if arguments == nil {
		return nil
	}

	args := map[string]string{}

	for _, argument := range *arguments {
		var (
			k        string
			v        string
			isSecret bool
		)
		k = argument.Name
		v = argument.Value

		if argument.IsSecret != nil {
			isSecret = *argument.IsSecret
		}

		// GET response won't return the value of secret arguments
		if isSecret {
			continue
		}

		args[k] = v
	}
	return args
}

func expandRegistryTaskValues(values map[string]string, secretValues map[string]string) *[]tasks.SetValue {
	if len(values) == 0 && len(secretValues) == 0 {
		return nil
	}
	out := make([]tasks.SetValue, 0, len(values)+len(secretValues))
	for k, v := range values {
		out = append(out, tasks.SetValue{
			Name:     k,
			Value:    v,
			IsSecret: utils.Bool(false),
		})
	}
	for k, v := range secretValues {
		out = append(out, tasks.SetValue{
			Name:     k,
			Value:    v,
			IsSecret: utils.Bool(true),
		})
	}
	return &out
}

func flattenRegistryTaskValues(values *[]tasks.SetValue) map[string]string {
	if values == nil {
		return nil
	}

	vals := map[string]string{}

	for _, value := range *values {
		var (
			k        string
			v        string
			isSecret bool
		)
		k = value.Name
		v = value.Value

		if value.IsSecret != nil {
			isSecret = *value.IsSecret
		}

		if isSecret {
			// GET response won't return the value of secret values
			continue
		}
		vals[k] = v
	}
	return vals
}

func expandRegistryTaskPlatform(input []Platform) *tasks.PlatformProperties {
	if len(input) == 0 {
		return nil
	}
	platform := input[0]
	out := &tasks.PlatformProperties{
		Os: tasks.OS(platform.OS),
	}
	if arch := platform.Architecture; arch != "" {
		out.Architecture = pointer.To(tasks.Architecture(arch))
	}
	if variant := platform.Variant; variant != "" {
		out.Variant = pointer.To(tasks.Variant(variant))
	}
	return out
}

func flattenRegistryTaskPlatform(platform *tasks.PlatformProperties) []Platform {
	if platform == nil {
		return nil
	}

	architecture := ""
	if v := platform.Architecture; v != nil {
		architecture = string(*v)
	}

	variant := ""
	if v := platform.Variant; v != nil {
		variant = string(*v)
	}

	return []Platform{{
		OS:           string(platform.Os),
		Architecture: architecture,
		Variant:      variant,
	}}
}

func expandRegistryTaskCredentials(input []RegistryCredential) *tasks.Credentials {
	if len(input) == 0 {
		return nil
	}

	return &tasks.Credentials{
		SourceRegistry:   expandSourceRegistryCredential(input[0].Source),
		CustomRegistries: pointer.To(expandCustomRegistryCredential(input[0].Custom)),
	}
}

func flattenRegistryTaskCredentials(input *tasks.Credentials, model ContainerRegistryTaskModel) []RegistryCredential {
	if input == nil {
		return nil
	}

	// The customRegistryCredentials is sensitive and won't return from API, setting it from the config.
	var custom []CustomRegistryCredential
	if len(model.RegistryCredential) == 1 {
		custom = model.RegistryCredential[0].Custom
	}

	return []RegistryCredential{
		{
			Source: flattenSourceRegistryCredential(input.SourceRegistry),
			Custom: custom,
		},
	}
}

func expandSourceRegistryCredential(input []SourceRegistryCredential) *tasks.SourceRegistryCredentials {
	if len(input) == 0 {
		return nil
	}

	return &tasks.SourceRegistryCredentials{LoginMode: pointer.To(tasks.SourceRegistryLoginMode(input[0].LoginMode))}
}

func flattenSourceRegistryCredential(input *tasks.SourceRegistryCredentials) []SourceRegistryCredential {
	if input == nil {
		return nil
	}

	return []SourceRegistryCredential{{LoginMode: string(*input.LoginMode)}}
}

func expandCustomRegistryCredential(input []CustomRegistryCredential) map[string]tasks.CustomRegistryCredentials {
	if len(input) == 0 {
		return nil
	}

	out := map[string]tasks.CustomRegistryCredentials{}
	for _, credential := range input {
		cred := tasks.CustomRegistryCredentials{}

		if credential.UserName != "" {
			usernameType := tasks.SecretObjectTypeOpaque
			if _, err := keyVaultParse.ParseNestedItemID(credential.UserName); err == nil {
				usernameType = tasks.SecretObjectTypeVaultsecret
			}
			cred.UserName = &tasks.SecretObject{
				Value: utils.String(credential.UserName),
				Type:  &usernameType,
			}
		}
		if credential.Password != "" {
			passwordType := tasks.SecretObjectTypeOpaque
			if _, err := keyVaultParse.ParseNestedItemID(credential.Password); err == nil {
				passwordType = tasks.SecretObjectTypeVaultsecret
			}
			cred.Password = &tasks.SecretObject{
				Value: utils.String(credential.Password),
				Type:  &passwordType,
			}
		}
		if credential.Identity != "" {
			cred.Identity = utils.String(credential.Identity)
		}
		out[credential.LoginServer] = cred
	}
	return out
}

func expandRegistryTaskAgentProperties(input []AgentConfig) *tasks.AgentProperties {
	if len(input) == 0 {
		return nil
	}

	agentConfig := input[0]
	return &tasks.AgentProperties{Cpu: utils.Int64(int64(agentConfig.CPU))}
}

func flattenRegistryTaskAgentProperties(input *tasks.AgentProperties) []AgentConfig {
	if input == nil {
		return nil
	}

	cpu := 0
	if input.Cpu != nil {
		cpu = int(*input.Cpu)
	}
	return []AgentConfig{{CPU: cpu}}
}

func patchRegistryTaskTriggerSourceTrigger(triggers []tasks.SourceTrigger, model ContainerRegistryTaskModel) *[]tasks.SourceTrigger {
	if len(triggers) != len(model.SourceTrigger) {
		return &triggers
	}

	result := make([]tasks.SourceTrigger, len(triggers))
	for i, trigger := range model.SourceTrigger {
		t := (triggers)[i]
		if len(trigger.Auth) == 0 {
			result[i] = t
			continue
		}

		t.SourceRepository.SourceControlAuthProperties = expandRegistryTaskAuthInfo(trigger.Auth[0])
		result[i] = t
	}
	return &result
}
