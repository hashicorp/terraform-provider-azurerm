// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_dapr_component

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/daprcomponents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.DaprComponentsClient
			id, err := daprcomponents.ParseDaprComponentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			daprComponentResp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(daprComponentResp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var state ContainerAppEnvironmentDaprComponentModel

			state.Name = id.DaprComponentName
			state.ManagedEnvironmentId = daprcomponents.NewManagedEnvironmentID(id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName).ID()

			if model := daprComponentResp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Version = pointer.From(props.Version)
					state.ComponentType = pointer.From(props.ComponentType)
					state.Scopes = pointer.From(props.Scopes)
					state.InitTimeout = pointer.From(props.InitTimeout)
					state.IgnoreErrors = pointer.From(props.IgnoreErrors)
					state.Metadata = flattenDaprComponentPropertiesMetadata(props.Metadata)
				}
			}

			secretsResp, err := client.ListSecrets(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving secrets for %s: %+v", *id, err)
			}

			state.Secrets = helpers.FlattenContainerAppDaprSecrets(secretsResp.Model)

			return metadata.Encode(&state)
		},
	}
}

func flattenDaprComponentPropertiesMetadata(input *[]daprcomponents.DaprMetadata) []helpers.DaprMetadata {
	if input == nil {
		return []helpers.DaprMetadata{}
	}

	result := make([]helpers.DaprMetadata, 0)
	for _, v := range *input {
		result = append(result, helpers.DaprMetadata{
			Name:       pointer.From(v.Name),
			SecretName: pointer.From(v.SecretRef),
			Value:      pointer.From(v.Value),
		})
	}

	return result
}
