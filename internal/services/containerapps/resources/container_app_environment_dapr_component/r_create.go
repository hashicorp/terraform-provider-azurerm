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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.DaprComponentsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var daprComponent ContainerAppEnvironmentDaprComponentModel

			if err := metadata.Decode(&daprComponent); err != nil {
				return err
			}

			managedEnvironmentId, err := daprcomponents.ParseManagedEnvironmentID(daprComponent.ManagedEnvironmentId)
			if err != nil {
				return err
			}

			id := daprcomponents.NewDaprComponentID(subscriptionId, managedEnvironmentId.ResourceGroupName, managedEnvironmentId.ManagedEnvironmentName, daprComponent.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			daprComponentRequest := daprcomponents.DaprComponent{
				Properties: &daprcomponents.DaprComponentProperties{
					ComponentType: pointer.To(daprComponent.ComponentType),
					IgnoreErrors:  pointer.To(daprComponent.IgnoreErrors),
					InitTimeout:   pointer.To(daprComponent.InitTimeout),
					Metadata:      expandDaprComponentPropertiesMetadata(daprComponent.Metadata),
					Secrets:       helpers.ExpandDaprSecrets(daprComponent.Secrets),
					Scopes:        pointer.To(daprComponent.Scopes),
					Version:       pointer.To(daprComponent.Version),
				},
			}

			if len(daprComponent.Scopes) > 0 {
				daprComponentRequest.Properties.Scopes = &daprComponent.Scopes
			}

			if _, err := client.CreateOrUpdate(ctx, id, daprComponentRequest); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandDaprComponentPropertiesMetadata(input []helpers.DaprMetadata) *[]daprcomponents.DaprMetadata {
	if len(input) == 0 {
		return nil
	}

	result := make([]daprcomponents.DaprMetadata, 0)

	for _, v := range input {
		d := daprcomponents.DaprMetadata{
			Name: pointer.To(v.Name),
		}
		if v.Value != "" {
			d.Value = pointer.To(v.Value)
		}
		if v.SecretName != "" {
			d.SecretRef = pointer.To(v.SecretName)
		}
		result = append(result, d)
	}

	return &result
}
