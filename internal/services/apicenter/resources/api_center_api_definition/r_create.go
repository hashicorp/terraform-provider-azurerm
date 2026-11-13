// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_api_definition

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apidefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apiversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiCenter.ApiDefinitionsClient

			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ApiCenterApiDefinitionModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			if config.Title == "" {
				config.Title = config.Name
			}

			api, err := apiversions.ParseVersionID(config.ApiVersionId)
			if err != nil {
				return err
			}

			id := apidefinitions.NewDefinitionID(subscriptionId, api.ResourceGroupName, api.ServiceName, "default", api.ApiName, api.VersionName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := apidefinitions.ApiDefinition{
				Name: pointer.To(config.Name),
				Properties: &apidefinitions.ApiDefinitionProperties{
					Description: pointer.To(config.Description),
					Title:       config.Title,
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			importSpec := expandApiDefinitionPropertiesSpecification(config.Specification)
			if importSpec != nil {
				if err := client.ImportSpecificationThenPoll(ctx, id, *importSpec); err != nil {
					return fmt.Errorf("importing specification for %s: %+v", id, err)
				}
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}
