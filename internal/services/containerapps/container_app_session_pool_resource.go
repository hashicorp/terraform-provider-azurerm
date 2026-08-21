// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containerapps

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/containerappssessionpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity

type ContainerAppSessionPoolResource struct{}

var (
	_ sdk.ResourceWithCustomizeDiff = ContainerAppSessionPoolResource{}
	_ sdk.ResourceWithIdentity      = ContainerAppSessionPoolResource{}
	_ sdk.ResourceWithUpdate        = ContainerAppSessionPoolResource{}
)

type ContainerAppSessionPoolModel struct {
	Name                      string                                     `tfschema:"name"`
	ResourceGroupName         string                                     `tfschema:"resource_group_name"`
	Location                  string                                     `tfschema:"location"`
	ContainerType             string                                     `tfschema:"container_type"`
	MaxConcurrentSessions     int64                                      `tfschema:"max_concurrent_sessions"`
	ContainerAppEnvironmentId string                                     `tfschema:"container_app_environment_id"`
	CooldownPeriodInSeconds   int64                                      `tfschema:"cooldown_period_in_seconds"`
	CustomContainerTemplate   []CustomContainerTemplateModel             `tfschema:"custom_container_template"`
	Identity                  []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	LifecycleType             string                                     `tfschema:"lifecycle_type"`
	SessionManagedIdentities  []string                                   `tfschema:"session_managed_identities"`
	MaxAlivePeriodInSeconds   int64                                      `tfschema:"max_alive_period_in_seconds"`
	NetworkEgressEnabled      bool                                       `tfschema:"network_egress_enabled"`
	ReadySessionInstances     int64                                      `tfschema:"ready_session_instances"`
	Secrets                   []SessionPoolSecretModel                   `tfschema:"secret"`
	Tags                      map[string]string                          `tfschema:"tags"`

	NodeCount              int64  `tfschema:"node_count"`
	PoolManagementEndpoint string `tfschema:"pool_management_endpoint"`
}

type CustomContainerTemplateModel struct {
	Containers        []SessionContainerModel `tfschema:"container"`
	IngressTargetPort int64                   `tfschema:"ingress_target_port"`
	Registry          []helpers.Registry      `tfschema:"registry"`
}

type SessionContainerModel struct {
	Name    string                    `tfschema:"name"`
	Image   string                    `tfschema:"image"`
	Args    []string                  `tfschema:"args"`
	Command []string                  `tfschema:"command"`
	Cpu     float64                   `tfschema:"cpu"`
	Env     []helpers.ContainerEnvVar `tfschema:"env"`
	Memory  string                    `tfschema:"memory"`
}

type SessionPoolSecretModel struct {
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
}

func (r ContainerAppSessionPoolResource) ResourceType() string {
	return "azurerm_container_app_session_pool"
}

func (r ContainerAppSessionPoolResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return containerappssessionpools.ValidateSessionPoolID
}

func (r ContainerAppSessionPoolResource) Identity() resourceids.ResourceId {
	return &containerappssessionpools.SessionPoolId{}
}

func (r ContainerAppSessionPoolResource) ModelObject() interface{} {
	return &ContainerAppSessionPoolModel{}
}

func (r ContainerAppSessionPoolResource) CustomizeDiff() sdk.ResourceFunc {
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

			switch config.LifecycleType {
			case string(containerappssessionpools.LifecycleTypeTimed):
				if !isConfigured("cooldown_period_in_seconds") {
					return errors.New("`cooldown_period_in_seconds` must be set when `lifecycle_type` is `Timed`")
				}
				if isConfigured("max_alive_period_in_seconds") {
					return errors.New("`max_alive_period_in_seconds` cannot be set when `lifecycle_type` is `Timed`, use `cooldown_period_in_seconds` instead")
				}

			case string(containerappssessionpools.LifecycleTypeOnContainerExit):
				if config.ContainerType != string(containerappssessionpools.ContainerTypeCustomContainer) {
					return errors.New("`lifecycle_type` can only be set to `OnContainerExit` when `container_type` is `CustomContainer`")
				}
				if isConfigured("cooldown_period_in_seconds") {
					return errors.New("`cooldown_period_in_seconds` cannot be set when `lifecycle_type` is `OnContainerExit`, use `max_alive_period_in_seconds` instead")
				}
				if !isConfigured("max_alive_period_in_seconds") {
					return errors.New("`max_alive_period_in_seconds` must be set when `lifecycle_type` is `OnContainerExit`")
				}
			}

			return nil
		},
	}
}

func (r ContainerAppSessionPoolResource) Arguments() map[string]*pluginsdk.Schema {
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

		"cooldown_period_in_seconds": {
			Type:          pluginsdk.TypeInt,
			Optional:      true,
			ValidateFunc:  validation.IntBetween(300, 3600),
			ConflictsWith: []string{"max_alive_period_in_seconds"},
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

		"lifecycle_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(containerappssessionpools.LifecycleTypeTimed),
			ValidateFunc: validation.StringInSlice(containerappssessionpools.PossibleValuesForLifecycleType(), false),
		},

		"max_alive_period_in_seconds": {
			Type:          pluginsdk.TypeInt,
			Optional:      true,
			ValidateFunc:  validation.IntAtLeast(1),
			ConflictsWith: []string{"cooldown_period_in_seconds"},
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

func (r ContainerAppSessionPoolResource) Attributes() map[string]*pluginsdk.Schema {
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

func (r ContainerAppSessionPoolResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.SessionPoolClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ContainerAppSessionPoolModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := containerappssessionpools.NewSessionPoolID(subscriptionId, config.ResourceGroupName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(config.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			customContainerTemplate, err := expandContainerAppSessionPoolCustomContainerTemplate(config.CustomContainerTemplate)
			if err != nil {
				return fmt.Errorf("expanding `custom_container_template`: %+v", err)
			}

			networkStatus := containerappssessionpools.SessionNetworkStatusEgressDisabled
			if config.NetworkEgressEnabled {
				networkStatus = containerappssessionpools.SessionNetworkStatusEgressEnabled
			}

			payload := containerappssessionpools.SessionPool{
				Identity: expandedIdentity,
				Location: location.Normalize(config.Location),
				Properties: &containerappssessionpools.SessionPoolProperties{
					ContainerType:            pointer.To(containerappssessionpools.ContainerType(config.ContainerType)),
					CustomContainerTemplate:  customContainerTemplate,
					DynamicPoolConfiguration: expandContainerAppSessionPoolDynamicPoolConfiguration(config),
					ManagedIdentitySettings:  expandContainerAppSessionPoolManagedIdentitySettings(config.SessionManagedIdentities),
					PoolManagementType:       pointer.To(containerappssessionpools.PoolManagementTypeDynamic),
					ScaleConfiguration:       expandContainerAppSessionPoolScaleConfiguration(config),
					Secrets:                  expandContainerAppSessionPoolSecrets(config.Secrets),
					SessionNetworkConfiguration: &containerappssessionpools.SessionNetworkConfiguration{
						Status: pointer.To(networkStatus),
					},
				},
				Tags: pointer.To(config.Tags),
			}

			if config.ContainerAppEnvironmentId != "" {
				payload.Properties.EnvironmentId = pointer.To(config.ContainerAppEnvironmentId)
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func (r ContainerAppSessionPoolResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.SessionPoolClient

			id, err := containerappssessionpools.ParseSessionPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config ContainerAppSessionPoolModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			payload := *existing.Model
			props := payload.Properties

			if metadata.ResourceData.HasChange("custom_container_template") {
				customContainerTemplate, err := expandContainerAppSessionPoolCustomContainerTemplate(config.CustomContainerTemplate)
				if err != nil {
					return fmt.Errorf("expanding `custom_container_template`: %+v", err)
				}
				props.CustomContainerTemplate = customContainerTemplate
			}

			if metadata.ResourceData.HasChanges("cooldown_period_in_seconds", "lifecycle_type", "max_alive_period_in_seconds") {
				props.DynamicPoolConfiguration = expandContainerAppSessionPoolDynamicPoolConfiguration(config)
			}

			if metadata.ResourceData.HasChanges("max_concurrent_sessions", "ready_session_instances") {
				props.ScaleConfiguration = expandContainerAppSessionPoolScaleConfiguration(config)
			}

			// clearing this sends a null, which the API resets to the `None` lifecycle stage
			if metadata.ResourceData.HasChange("session_managed_identities") {
				props.ManagedIdentitySettings = expandContainerAppSessionPoolManagedIdentitySettings(config.SessionManagedIdentities)
			}

			if metadata.ResourceData.HasChange("network_egress_enabled") {
				networkStatus := containerappssessionpools.SessionNetworkStatusEgressDisabled
				if config.NetworkEgressEnabled {
					networkStatus = containerappssessionpools.SessionNetworkStatusEgressEnabled
				}

				props.SessionNetworkConfiguration = &containerappssessionpools.SessionNetworkConfiguration{
					Status: pointer.To(networkStatus),
				}
			}

			// The API doesn't return secret values, so these are always sent from the config to avoid clearing them.
			props.Secrets = expandContainerAppSessionPoolSecrets(config.Secrets)

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(config.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(config.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ContainerAppSessionPoolResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.SessionPoolClient

			id, err := containerappssessionpools.ParseSessionPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r ContainerAppSessionPoolResource) flatten(metadata sdk.ResourceMetaData, id *containerappssessionpools.SessionPoolId, model *containerappssessionpools.SessionPool) error {
	var existing ContainerAppSessionPoolModel
	if err := metadata.Decode(&existing); err != nil {
		return fmt.Errorf("decoding: %+v", err)
	}

	state := ContainerAppSessionPoolModel{
		Name:              id.SessionPoolName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.Tags = pointer.From(model.Tags)

		flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		state.Identity = flattenedIdentity

		if props := model.Properties; props != nil {
			if props.EnvironmentId != nil {
				environmentId, err := managedenvironments.ParseManagedEnvironmentIDInsensitively(pointer.From(props.EnvironmentId))
				if err != nil {
					return fmt.Errorf("parsing `container_app_environment_id`: %+v", err)
				}
				state.ContainerAppEnvironmentId = environmentId.ID()
			}
			state.ContainerType = string(pointer.From(props.ContainerType))
			state.CustomContainerTemplate = flattenContainerAppSessionPoolCustomContainerTemplate(props.CustomContainerTemplate)
			state.SessionManagedIdentities = flattenContainerAppSessionPoolManagedIdentitySettings(props.ManagedIdentitySettings)
			state.NodeCount = pointer.From(props.NodeCount)
			state.PoolManagementEndpoint = pointer.From(props.PoolManagementEndpoint)
			state.Secrets = flattenContainerAppSessionPoolSecrets(props.Secrets, existing.Secrets)

			if dynamicPool := props.DynamicPoolConfiguration; dynamicPool != nil && dynamicPool.LifecycleConfiguration != nil {
				lifecycle := dynamicPool.LifecycleConfiguration
				state.CooldownPeriodInSeconds = pointer.From(lifecycle.CooldownPeriodInSeconds)
				state.LifecycleType = string(pointer.From(lifecycle.LifecycleType))
				state.MaxAlivePeriodInSeconds = pointer.From(lifecycle.MaxAlivePeriodInSeconds)
			}

			if scale := props.ScaleConfiguration; scale != nil {
				state.MaxConcurrentSessions = pointer.From(scale.MaxConcurrentSessions)
				state.ReadySessionInstances = pointer.From(scale.ReadySessionInstances)
			}

			if network := props.SessionNetworkConfiguration; network != nil {
				state.NetworkEgressEnabled = pointer.From(network.Status) == containerappssessionpools.SessionNetworkStatusEgressEnabled
			}
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r ContainerAppSessionPoolResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.SessionPoolClient

			id, err := containerappssessionpools.ParseSessionPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandContainerAppSessionPoolCustomContainerTemplate(input []CustomContainerTemplateModel) (*containerappssessionpools.CustomContainerTemplate, error) {
	if len(input) == 0 {
		return nil, nil
	}

	v := input[0]
	result := &containerappssessionpools.CustomContainerTemplate{
		Containers: expandContainerAppSessionPoolContainers(v.Containers),
	}

	if v.IngressTargetPort != 0 {
		result.Ingress = &containerappssessionpools.SessionIngress{
			TargetPort: pointer.To(v.IngressTargetPort),
		}
	}

	if len(v.Registry) > 0 {
		registry := v.Registry[0]
		if err := helpers.ValidateContainerAppRegistry(registry); err != nil {
			return nil, err
		}

		result.RegistryCredentials = &containerappssessionpools.SessionRegistryCredentials{
			Server: pointer.To(registry.Server),
		}

		if registry.Identity != "" {
			result.RegistryCredentials.Identity = pointer.To(registry.Identity)
		}
		if registry.PasswordSecretRef != "" {
			result.RegistryCredentials.PasswordSecretRef = pointer.To(registry.PasswordSecretRef)
		}
		if registry.UserName != "" {
			result.RegistryCredentials.Username = pointer.To(registry.UserName)
		}
	}

	return result, nil
}

func expandContainerAppSessionPoolContainers(input []SessionContainerModel) *[]containerappssessionpools.SessionContainer {
	if len(input) == 0 {
		return nil
	}

	result := make([]containerappssessionpools.SessionContainer, 0, len(input))
	for _, v := range input {
		container := containerappssessionpools.SessionContainer{
			Env:   expandContainerAppSessionPoolEnvironmentVars(v.Env),
			Image: pointer.To(v.Image),
			Name:  pointer.To(v.Name),
		}

		if len(v.Args) > 0 {
			container.Args = pointer.To(v.Args)
		}
		if len(v.Command) > 0 {
			container.Command = pointer.To(v.Command)
		}
		if v.Cpu != 0 || v.Memory != "" {
			container.Resources = &containerappssessionpools.SessionContainerResources{}
			if v.Cpu != 0 {
				container.Resources.Cpu = pointer.To(v.Cpu)
			}
			if v.Memory != "" {
				container.Resources.Memory = pointer.To(v.Memory)
			}
		}

		result = append(result, container)
	}

	return &result
}

func expandContainerAppSessionPoolEnvironmentVars(input []helpers.ContainerEnvVar) *[]containerappssessionpools.EnvironmentVar {
	if len(input) == 0 {
		return nil
	}

	result := make([]containerappssessionpools.EnvironmentVar, 0, len(input))
	for _, v := range input {
		env := containerappssessionpools.EnvironmentVar{
			Name: pointer.To(v.Name),
		}

		if v.SecretReference != "" {
			env.SecretRef = pointer.To(v.SecretReference)
		} else {
			env.Value = pointer.To(v.Value)
		}

		result = append(result, env)
	}

	return &result
}

func expandContainerAppSessionPoolDynamicPoolConfiguration(input ContainerAppSessionPoolModel) *containerappssessionpools.DynamicPoolConfiguration {
	// The API rejects a null `dynamicPoolConfiguration` when `poolManagementType` is `Dynamic`.
	lifecycle := &containerappssessionpools.LifecycleConfiguration{
		LifecycleType: pointer.To(containerappssessionpools.LifecycleType(input.LifecycleType)),
	}

	switch containerappssessionpools.LifecycleType(input.LifecycleType) {
	case containerappssessionpools.LifecycleTypeTimed:
		if input.CooldownPeriodInSeconds != 0 {
			lifecycle.CooldownPeriodInSeconds = pointer.To(input.CooldownPeriodInSeconds)
		}

	case containerappssessionpools.LifecycleTypeOnContainerExit:
		if input.MaxAlivePeriodInSeconds != 0 {
			lifecycle.MaxAlivePeriodInSeconds = pointer.To(input.MaxAlivePeriodInSeconds)
		}
	}

	return &containerappssessionpools.DynamicPoolConfiguration{
		LifecycleConfiguration: lifecycle,
	}
}

func expandContainerAppSessionPoolScaleConfiguration(input ContainerAppSessionPoolModel) *containerappssessionpools.ScaleConfiguration {
	result := &containerappssessionpools.ScaleConfiguration{
		MaxConcurrentSessions: pointer.To(input.MaxConcurrentSessions),
	}

	if input.ReadySessionInstances != 0 {
		result.ReadySessionInstances = pointer.To(input.ReadySessionInstances)
	}

	return result
}

func expandContainerAppSessionPoolManagedIdentitySettings(input []string) *[]containerappssessionpools.ManagedIdentitySetting {
	if len(input) == 0 {
		return nil
	}

	result := make([]containerappssessionpools.ManagedIdentitySetting, 0, len(input))
	for _, v := range input {
		result = append(result, containerappssessionpools.ManagedIdentitySetting{
			Identity:  v,
			Lifecycle: pointer.To(containerappssessionpools.IdentitySettingsLifeCycleMain),
		})
	}

	return &result
}

func expandContainerAppSessionPoolSecrets(input []SessionPoolSecretModel) *[]containerappssessionpools.SessionPoolSecret {
	if len(input) == 0 {
		return nil
	}

	result := make([]containerappssessionpools.SessionPoolSecret, 0, len(input))
	for _, v := range input {
		result = append(result, containerappssessionpools.SessionPoolSecret{
			Name:  pointer.To(v.Name),
			Value: pointer.To(v.Value),
		})
	}

	return &result
}

func flattenContainerAppSessionPoolCustomContainerTemplate(input *containerappssessionpools.CustomContainerTemplate) []CustomContainerTemplateModel {
	if input == nil {
		return []CustomContainerTemplateModel{}
	}

	result := CustomContainerTemplateModel{
		Containers: flattenContainerAppSessionPoolContainers(input.Containers),
	}

	if ingress := input.Ingress; ingress != nil {
		result.IngressTargetPort = pointer.From(ingress.TargetPort)
	}

	if credentials := input.RegistryCredentials; credentials != nil {
		result.Registry = []helpers.Registry{
			{
				Identity:          normaliseContainerAppSessionPoolIdentity(pointer.From(credentials.Identity)),
				PasswordSecretRef: pointer.From(credentials.PasswordSecretRef),
				Server:            pointer.From(credentials.Server),
				UserName:          pointer.From(credentials.Username),
			},
		}
	}

	return []CustomContainerTemplateModel{result}
}

func flattenContainerAppSessionPoolContainers(input *[]containerappssessionpools.SessionContainer) []SessionContainerModel {
	if input == nil {
		return []SessionContainerModel{}
	}

	result := make([]SessionContainerModel, 0, len(*input))
	for _, v := range *input {
		container := SessionContainerModel{
			Args:    pointer.From(v.Args),
			Command: pointer.From(v.Command),
			Env:     flattenContainerAppSessionPoolEnvironmentVars(v.Env),
			Image:   pointer.From(v.Image),
			Name:    pointer.From(v.Name),
		}

		if resources := v.Resources; resources != nil {
			container.Cpu = pointer.From(resources.Cpu)
			container.Memory = pointer.From(resources.Memory)
		}

		result = append(result, container)
	}

	return result
}

func flattenContainerAppSessionPoolEnvironmentVars(input *[]containerappssessionpools.EnvironmentVar) []helpers.ContainerEnvVar {
	if input == nil {
		return []helpers.ContainerEnvVar{}
	}

	result := make([]helpers.ContainerEnvVar, 0, len(*input))
	for _, v := range *input {
		result = append(result, helpers.ContainerEnvVar{
			Name:            pointer.From(v.Name),
			SecretReference: pointer.From(v.SecretRef),
			Value:           pointer.From(v.Value),
		})
	}

	return result
}

func flattenContainerAppSessionPoolManagedIdentitySettings(input *[]containerappssessionpools.ManagedIdentitySetting) []string {
	if input == nil {
		return []string{}
	}

	result := make([]string, 0, len(*input))
	for _, v := range *input {
		// the API adds a `None` entry for every assigned identity that isn't exposed to the sessions
		if pointer.From(v.Lifecycle) == containerappssessionpools.IdentitySettingsLifeCycleNone {
			continue
		}

		result = append(result, normaliseContainerAppSessionPoolIdentity(v.Identity))
	}

	return result
}

// The API returns `System` as `system`, and can return User Assigned Identity IDs with inconsistent segment casing.
func normaliseContainerAppSessionPoolIdentity(input string) string {
	if strings.EqualFold(input, "System") {
		return "System"
	}

	if id, err := commonids.ParseUserAssignedIdentityIDInsensitively(input); err == nil {
		return id.ID()
	}

	return input
}

func flattenContainerAppSessionPoolSecrets(input *[]containerappssessionpools.SessionPoolSecret, existing []SessionPoolSecretModel) []SessionPoolSecretModel {
	if input == nil {
		return []SessionPoolSecretModel{}
	}

	// `value` is flagged `x-ms-secret` in the API definition so it's never returned and has to be carried over from state.
	valuesFromState := make(map[string]string, len(existing))
	for _, v := range existing {
		valuesFromState[v.Name] = v.Value
	}

	result := make([]SessionPoolSecretModel, 0, len(*input))
	for _, v := range *input {
		name := pointer.From(v.Name)
		result = append(result, SessionPoolSecretModel{
			Name:  name,
			Value: valuesFromState[name],
		})
	}

	return result
}
