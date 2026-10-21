// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry_task

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/tasks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/base64"
)

func (r Resource) Read() sdk.ResourceFunc {
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
				timeoutInSec       int64
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
						timeoutInSec = *props.Timeout
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

type AgentConfig struct {
	CPU int64 `tfschema:"cpu"`
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
	ExpireInSec  int64  `tfschema:"expire_in_seconds"`
}

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

func userDataStateFunc(v any) string {
	switch s := v.(type) {
	case string:
		return base64.EncodeIfNot(s)
	default:
		return ""
	}
}
