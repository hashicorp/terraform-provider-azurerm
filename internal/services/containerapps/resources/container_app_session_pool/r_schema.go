// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_session_pool

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/containerappssessionpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			var config ContainerAppSessionPoolModel
			if err := metadata.DecodeDiff(&config); err != nil {
				return err
			}

			rawConfig := metadata.ResourceDiff.GetRawConfig()
			isConfigured := func(name string) bool {
				// `AsValueMap` panics on a null raw config, which Terraform can send during a destroy plan
				if rawConfig.IsNull() {
					return false
				}

				value, ok := rawConfig.AsValueMap()[name]
				return ok && !value.IsNull()
			}

			if config.ContainerType != string(containerappssessionpools.ContainerTypeCustomContainer) && len(config.Identity) > 0 {
				return errors.New("`identity` can only be set when `container_type` is `CustomContainer`")
			}

			if config.ReadySessionInstances != 0 && config.ReadySessionInstances >= config.MaxConcurrentSessions {
				return errors.New("`ready_session_instances` must be less than `max_concurrent_sessions`")
			}

			switch config.ContainerType {
			case string(containerappssessionpools.ContainerTypeCustomContainer):
				if len(config.CustomContainerTemplate) == 0 {
					return errors.New("`custom_container_template` must be set when `container_type` is `CustomContainer`")
				}

				if !isConfigured("container_app_environment_id") {
					return errors.New("`container_app_environment_id` must be set when `container_type` is `CustomContainer`")
				}

				if !isConfigured("ready_session_instances") {
					return errors.New("`ready_session_instances` must be set when `container_type` is `CustomContainer`")
				}

			default:
				if len(config.CustomContainerTemplate) > 0 {
					return errors.New("`custom_container_template` can only be set when `container_type` is `CustomContainer`")
				}
			}

			if len(config.LifecycleConfiguration) == 0 {
				return nil
			}

			rawLifecycleConfiguration := rawConfig.AsValueMap()["lifecycle_configuration"]
			if rawLifecycleConfiguration.IsNull() || !rawLifecycleConfiguration.IsKnown() || len(rawLifecycleConfiguration.AsValueSlice()) == 0 {
				return nil
			}

			rawLifecycle := rawLifecycleConfiguration.AsValueSlice()[0]
			if rawLifecycle.IsNull() || !rawLifecycle.IsKnown() {
				return nil
			}

			rawLifecycleProperties := rawLifecycle.AsValueMap()
			cooldownPeriodInSecondsConfigured := !rawLifecycleProperties["cooldown_period_in_seconds"].IsNull()
			maxAlivePeriodInSecondsConfigured := !rawLifecycleProperties["max_alive_period_in_seconds"].IsNull()

			lifecycle := config.LifecycleConfiguration[0]
			switch lifecycle.LifecycleType {
			case string(containerappssessionpools.LifecycleTypeTimed):
				if !cooldownPeriodInSecondsConfigured {
					return errors.New("`lifecycle_configuration.0.cooldown_period_in_seconds` must be set when `lifecycle_configuration.0.lifecycle_type` is `Timed`")
				}
				if maxAlivePeriodInSecondsConfigured {
					return errors.New("`lifecycle_configuration.0.max_alive_period_in_seconds` cannot be set when `lifecycle_configuration.0.lifecycle_type` is `Timed`, use `lifecycle_configuration.0.cooldown_period_in_seconds` instead")
				}

			case string(containerappssessionpools.LifecycleTypeOnContainerExit):
				if config.ContainerType != string(containerappssessionpools.ContainerTypeCustomContainer) {
					return errors.New("`lifecycle_configuration.0.lifecycle_type` can only be set to `OnContainerExit` when `container_type` is `CustomContainer`")
				}
				if cooldownPeriodInSecondsConfigured {
					return errors.New("`lifecycle_configuration.0.cooldown_period_in_seconds` cannot be set when `lifecycle_configuration.0.lifecycle_type` is `OnContainerExit`, use `lifecycle_configuration.0.max_alive_period_in_seconds` instead")
				}
				if !maxAlivePeriodInSecondsConfigured {
					return errors.New("`lifecycle_configuration.0.max_alive_period_in_seconds` must be set when `lifecycle_configuration.0.lifecycle_type` is `OnContainerExit`")
				}
			}

			return nil
		},
	}
}

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-z][a-z0-9]{2,62}$`),
				"The `name` must be between 3 and 63 characters long, begin with a lowercase letter and contain only lowercase letters and numbers.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"container_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(containerappssessionpools.ContainerTypePythonLTS),
			ValidateFunc: validation.StringInSlice(containerappssessionpools.PossibleValuesForContainerType(), false),
		},

		"max_concurrent_sessions": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      5,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: managedenvironments.ValidateManagedEnvironmentID,
		},

		"custom_container_template": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"container": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validate.ContainerAppContainerName,
								},

								"image": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"args": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"command": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"cpu": {
									Type:         pluginsdk.TypeFloat,
									Optional:     true,
									ValidateFunc: validation.FloatAtLeast(0.1),
								},

								"env": helpers.ContainerEnvVarSchema(),

								"memory": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"ingress_target_port": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 65535),
					},

					"registry": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"server": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"identity": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.Any(
										commonids.ValidateUserAssignedIdentityID,
										validation.StringInSlice([]string{"System"}, false),
									),
								},

								"password_secret_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validate.SecretName,
								},

								"username": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"lifecycle_configuration": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"lifecycle_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(containerappssessionpools.LifecycleTypeTimed),
						ValidateFunc: validation.StringInSlice(containerappssessionpools.PossibleValuesForLifecycleType(), false),
					},

					"cooldown_period_in_seconds": {
						Type:          pluginsdk.TypeInt,
						Optional:      true,
						ValidateFunc:  validation.IntBetween(300, 3600),
						ConflictsWith: []string{"lifecycle_configuration.0.max_alive_period_in_seconds"},
					},

					"max_alive_period_in_seconds": {
						Type:          pluginsdk.TypeInt,
						Optional:      true,
						ValidateFunc:  validation.IntAtLeast(1),
						ConflictsWith: []string{"lifecycle_configuration.0.cooldown_period_in_seconds"},
					},
				},
			},
		},

		"network_egress_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"ready_session_instances": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"secret": {
			Type:      pluginsdk.TypeSet,
			Optional:  true,
			Sensitive: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.SecretName,
					},

					"value": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		// these must already be assigned via `identity`; listing one here exposes it to the code
		// running inside the sessions, which the API otherwise denies by default
		"session_managed_identities": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.Any(
					commonids.ValidateUserAssignedIdentityID,
					validation.StringInSlice([]string{"System"}, false),
				),
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"node_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"pool_management_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
