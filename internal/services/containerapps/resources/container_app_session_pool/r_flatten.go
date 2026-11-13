// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_session_pool

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/containerappssessionpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *containerappssessionpools.SessionPoolId, model *containerappssessionpools.SessionPool) error {
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
			state.ContainerType = pointer.FromEnum(props.ContainerType)
			state.CustomContainerTemplate = flattenContainerAppSessionPoolCustomContainerTemplate(props.CustomContainerTemplate)
			state.SessionManagedIdentities = flattenContainerAppSessionPoolManagedIdentitySettings(props.ManagedIdentitySettings)
			state.NodeCount = pointer.From(props.NodeCount)
			state.PoolManagementEndpoint = pointer.From(props.PoolManagementEndpoint)
			state.Secrets = flattenContainerAppSessionPoolSecrets(props.Secrets, existing.Secrets)

			if dynamicPool := props.DynamicPoolConfiguration; dynamicPool != nil && dynamicPool.LifecycleConfiguration != nil {
				state.LifecycleConfiguration = flattenContainerAppSessionPoolLifecycleConfiguration(dynamicPool.LifecycleConfiguration)
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

func flattenContainerAppSessionPoolLifecycleConfiguration(input *containerappssessionpools.LifecycleConfiguration) []LifecycleConfigurationModel {
	if input == nil {
		return []LifecycleConfigurationModel{}
	}

	return []LifecycleConfigurationModel{
		{
			CooldownPeriodInSeconds: pointer.From(input.CooldownPeriodInSeconds),
			LifecycleType:           pointer.FromEnum(input.LifecycleType),
			MaxAlivePeriodInSeconds: pointer.From(input.MaxAlivePeriodInSeconds),
		},
	}
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
