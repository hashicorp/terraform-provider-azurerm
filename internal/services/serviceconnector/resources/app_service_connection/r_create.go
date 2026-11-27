// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2024-04-01/servicelinker"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	serviceconnectorClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/serviceconnector/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/serviceconnector/helpers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model AppServiceConnectorResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := serviceconnectorClient.From(metadata.Client).ServiceLinkerClient

			id := servicelinker.NewScopedLinkerID(model.AppServiceId, model.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.LinkerGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			authInfo, err := helpers.ExpandServiceConnectorAuthInfoForCreate(model.AuthInfo)
			if err != nil {
				return fmt.Errorf("expanding `authentication`: %+v", err)
			}

			serviceConnectorProperties := servicelinker.LinkerProperties{
				AuthInfo: authInfo,
			}

			if storageAccountId, err := commonids.ParseStorageAccountID(model.TargetResourceId); err == nil {
				targetResourceId := fmt.Sprintf("%s/blobServices/default", storageAccountId.ID())
				serviceConnectorProperties.TargetService = servicelinker.AzureResource{
					Id: &targetResourceId,
				}
			} else {
				serviceConnectorProperties.TargetService = servicelinker.AzureResource{
					Id: &model.TargetResourceId,
				}
			}

			if model.SecretStore != nil {
				serviceConnectorProperties.SecretStore = helpers.ExpandSecretStore(model.SecretStore)
			}

			if model.ClientType != "" {
				serviceConnectorProperties.ClientType = pointer.ToEnum[servicelinker.ClientType](model.ClientType)
			}

			if model.VnetSolution != "" {
				vNetSolution := servicelinker.VNetSolution{
					Type: pointer.ToEnum[servicelinker.VNetSolutionType](model.VnetSolution),
				}
				serviceConnectorProperties.VNetSolution = &vNetSolution
			}

			props := servicelinker.LinkerResource{
				Id:         pointer.To(id.ID()),
				Name:       pointer.To(model.Name),
				Properties: serviceConnectorProperties,
			}

			if err := client.LinkerCreateOrUpdateCallbackThenPoll(ctx, id, props, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
