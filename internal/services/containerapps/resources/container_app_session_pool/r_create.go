// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_session_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/containerappssessionpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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
					ContainerType:            pointer.ToEnum[containerappssessionpools.ContainerType](config.ContainerType),
					CustomContainerTemplate:  customContainerTemplate,
					DynamicPoolConfiguration: expandContainerAppSessionPoolDynamicPoolConfiguration(config.LifecycleConfiguration),
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

func expandContainerAppSessionPoolDynamicPoolConfiguration(input []LifecycleConfigurationModel) *containerappssessionpools.DynamicPoolConfiguration {
	// The API rejects a null `dynamicPoolConfiguration` when `poolManagementType` is `Dynamic`.
	config := input[0]
	lifecycle := &containerappssessionpools.LifecycleConfiguration{
		LifecycleType: pointer.ToEnum[containerappssessionpools.LifecycleType](config.LifecycleType),
	}

	switch containerappssessionpools.LifecycleType(config.LifecycleType) {
	case containerappssessionpools.LifecycleTypeTimed:
		if config.CooldownPeriodInSeconds != 0 {
			lifecycle.CooldownPeriodInSeconds = pointer.To(config.CooldownPeriodInSeconds)
		}

	case containerappssessionpools.LifecycleTypeOnContainerExit:
		if config.MaxAlivePeriodInSeconds != 0 {
			lifecycle.MaxAlivePeriodInSeconds = pointer.To(config.MaxAlivePeriodInSeconds)
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
