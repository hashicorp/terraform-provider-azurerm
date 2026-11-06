// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_environment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/environments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiCenter.EnvironmentsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ApiCenterEnvironmentResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			if model.Title == "" {
				model.Title = model.Name
			}

			service, err := services.ParseServiceID(model.ServiceId)
			if err != nil {
				return err
			}

			// @favoretti: can't find any workspace creation buttons in the portal, assume everything defaults to "default" for now?
			id := environments.NewEnvironmentID(subscriptionId, service.ResourceGroupName, service.ServiceName, "default", model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing ApiCenter Environment %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			apiCenterEnvironmentProps := environments.EnvironmentProperties{
				Kind:  environments.EnvironmentKind(model.Type),
				Title: model.Title,
			}

			if model.Description != "" {
				apiCenterEnvironmentProps.Description = pointer.To(model.Description)
			}

			apiCenterEnvironmentProps.Onboarding = &environments.Onboarding{}
			if model.DevPortalUri != "" {
				apiCenterEnvironmentProps.Onboarding.DeveloperPortalUri = pointer.To([]string{model.DevPortalUri})
			}

			if model.Instructions != "" {
				apiCenterEnvironmentProps.Onboarding.Instructions = pointer.To(model.Instructions)
			}

			apiCenterEnvironmentProps.Server = &environments.EnvironmentServer{}

			if model.ServerType != "" {
				apiCenterEnvironmentProps.Server.Type = pointer.ToEnum[environments.EnvironmentServerType](model.ServerType)
			}

			if model.MgmtPortalUri != "" {
				apiCenterEnvironmentProps.Server.ManagementPortalUri = pointer.To([]string{model.MgmtPortalUri})
			}

			apiCenterEnvironment := environments.Environment{
				Name:       pointer.To(model.Name),
				Properties: pointer.To(apiCenterEnvironmentProps),
				Type:       pointer.To(model.Type),
			}

			if _, err = client.CreateOrUpdate(ctx, id, apiCenterEnvironment); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}
