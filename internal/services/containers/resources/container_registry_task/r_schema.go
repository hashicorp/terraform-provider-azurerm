// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_task

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/tasks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Container Registry Task.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ContainerRegistryTaskName,
		},
		"container_registry_id": {
			Description:  "The ID of the Container Registry that this Container Registry Task resides in.",
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
						Description:  "The operating system type required for the task.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(tasks.PossibleValuesForOS(), false),
					},
					"architecture": {
						Description:  "The OS architecture.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(tasks.PossibleValuesForArchitecture(), false),
					},
					"variant": {
						Description:  "The variant of the CPU.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(tasks.PossibleValuesForVariant(), false),
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
						Description:  "The Dockerfile path relative to the source context.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"context_path": {
						Description:  "The URL (absolute or relative) of the source context for this step. If the context is a URL you can reference a specific branch or folder via `#branch:folder`.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"context_access_token": {
						Description:  "The token (Git PAT or SAS token of Storage Account blob) associated with the context for this step.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"image_names": {
						Description: "A list of fully qualified image names including the repository and tag.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"push_enabled": {
						Description: "Should the image built be pushed to the registry or not?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
					"cache_enabled": {
						Description: "Should the image cache be enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
					"target": {
						Description:  "The name of the target build stage for the docker build.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"arguments": {
						Description: "A map of arguments to be used when executing this step.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"secret_arguments": {
						Description: "A map of *secret* arguments to be used when executing this step.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Sensitive:   true,
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
						Description:  "The task template file path relative to the source context.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"value_file_path": {
						Description:  "The parameters file path relative to the source context.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"context_path": {
						Description:  "The URL (absolute or relative) of the source context for this step.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"context_access_token": {
						Description:  "The token (Git PAT or SAS token of Storage Account blob) associated with the context for this step.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"values": {
						Description: "A map of values that can be passed when running a task.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"secret_values": {
						Description: "A map of secret values that can be passed when running a task.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Sensitive:   true,
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
						Description:  "The (optionally base64 encoded) content of the build template.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						StateFunc:    userDataStateFunc,
					},
					"value_content": {
						Description:  "The (optionally base64 encoded) content of the build parameters.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						StateFunc:    userDataStateFunc,
					},
					"context_path": {
						Description:  "The URL (absolute or relative) of the source context for this step.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"context_access_token": {
						Description:  "The token (Git PAT or SAS token of Storage Account blob) associated with the context for this step.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"values": {
						Description: "A map of values that can be passed when running a task.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"secret_values": {
						Description: "A map of secret values that can be passed when running a task.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Sensitive:   true,
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
						Description:  "The name which should be used for this trigger.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"type": {
						Description:  "The type of the trigger.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(tasks.PossibleValuesForBaseImageTriggerType(), false),
					},
					"enabled": {
						Description: "Should the trigger be enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
					"update_trigger_endpoint": {
						Description:  "The endpoint URL for receiving the trigger.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"update_trigger_payload_type": {
						Description:  "Type of payload body for the trigger.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(tasks.PossibleValuesForUpdateTriggerPayloadType(), false),
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
						Description:  "The name which should be used for this trigger.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"events": {
						Description: "A list of source events corresponding to the trigger.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(tasks.PossibleValuesForSourceTriggerEvent(), false),
						},
					},
					"source_type": {
						Description:  "The type of the source control service.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(tasks.PossibleValuesForSourceControlType(), false),
					},
					"repository_url": {
						Description:  "The full URL to the source code repository.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"branch": {
						Description: "The branch name of the source code.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"authentication": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"token_type": {
									Description:  "The type of the token.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(tasks.PossibleValuesForTokenType(), false),
								},
								"token": {
									Description:  "The access token used to access the source control provider.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"refresh_token": {
									Description:  "The refresh token used to refresh the access token.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"scope": {
									Description:  "The scope of the access token.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"expire_in_seconds": {
									Description: "Time in seconds that the token remains valid.",
									Type:        pluginsdk.TypeInt,
									Optional:    true,
								},
							},
						},
					},
					"enabled": {
						Description: "Should the trigger be enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
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
						Description:  "The name which should be used for this trigger.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"schedule": {
						Description:  "The CRON expression for the task schedule.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"enabled": {
						Description: "Should the trigger be enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
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
									Description:  "The login mode for the source registry.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(tasks.PossibleValuesForSourceRegistryLoginMode(), false),
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
									Description:  "The login server of the custom Container Registry.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"username": {
									Description:  "The username for logging into the custom Container Registry.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"password": {
									Description:  "The password for logging into the custom Container Registry.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"identity": {
									Description: "The Managed Identity assigned to this custom credential. For User Assigned Identity, the value is the client ID of the identity. For System Assigned Identity, the value is `[system]`.",
									Type:        pluginsdk.TypeString,
									Optional:    true,
									ValidateFunc: validation.Any(
										validation.StringInSlice([]string{"[system]"}, false),
										commonids.ValidateUserAssignedIdentityID,
									),
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
						Description:  "The number of cores required for the Container Registry Task.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntInSlice([]int{2}),
					},
				},
			},
			ConflictsWith: []string{"agent_pool_name"},
		},
		"agent_pool_name": {
			Description:   "The name of the dedicated Container Registry Agent Pool for this Container Registry Task.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"agent_setting"},
		},
		"enabled": {
			Description: "Should this Container Registry Task be enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},
		"is_system_task": {
			Description: "Whether this Container Registry Task is a system task.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},
		"log_template": {
			Description:  "The template that describes the run log artifact.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"timeout_in_seconds": {
			Description:  "The timeout of this Container Registry Task in seconds. The valid range lies from 300 to 28800.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(300, 28800),
			Default:      3600,
		},
		"tags": commonschema.Tags(),
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff
			isSystemTask := rd.Get("is_system_task").(bool)

			if isSystemTask {
				invalidProps := []string{"platform", "docker_step", "file_step", "encoded_step", "base_image_trigger", "source_trigger", "timer_trigger"}
				for _, prop := range invalidProps {
					if v := rd.Get(prop).([]any); len(v) != 0 {
						return fmt.Errorf("system task can't specify `%s`", prop)
					}
				}
			} else {
				if v := rd.Get("platform").([]any); len(v) == 0 {
					return fmt.Errorf("non-system task have to specify `platform`")
				}

				dockerStep := rd.Get("docker_step").([]any)
				fileTaskStep := rd.Get("file_step").([]any)
				encodedTaskStep := rd.Get("encoded_step").([]any)
				if len(dockerStep)+len(fileTaskStep)+len(encodedTaskStep) == 0 {
					return fmt.Errorf("non-system task have to specify one of `docker_step`, `file_step` and `encoded_step`")
				}
			}

			return nil
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
