// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_custom_domain

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			lawClient := metadata.Client.LogAnalytics.SharedKeyWorkspacesClient
			client := metadata.Client.ContainerApps.ManagedEnvironmentClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			// Get ID from input
			model := ContainerAppEnvironmentCustomDomainModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := managedenvironments.ParseManagedEnvironmentID(model.ManagedEnvironmentId)
			if err != nil {
				return err
			}

			// Prevent parallel create of the same resource
			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				if customDomain := existing.Model.Properties.CustomDomainConfiguration; customDomain != nil && customDomain.DnsSuffix != nil {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			// We need to check if a log analytic is attached and must get the shared key if it does
			if appsLogs := existing.Model.Properties.AppLogsConfiguration; appsLogs != nil && appsLogs.LogAnalyticsConfiguration != nil {
				sharedKey, err := findLogAnalyticsWorkspaceSecret(ctx, lawClient, subscriptionId, pointer.From(appsLogs.LogAnalyticsConfiguration.CustomerId))
				if err != nil {
					return fmt.Errorf("retrieving Log Analytics Workspace: %+v", err)
				}

				existing.Model.Properties.AppLogsConfiguration = &managedenvironments.AppLogsConfiguration{
					Destination: pointer.To("log-analytics"),
					LogAnalyticsConfiguration: &managedenvironments.LogAnalyticsConfiguration{
						CustomerId: appsLogs.LogAnalyticsConfiguration.CustomerId,
						SharedKey:  pointer.To(sharedKey),
					},
				}
			} else {
				existing.Model.Properties.AppLogsConfiguration = nil
			}

			existing.Model.Properties.CustomDomainConfiguration = &managedenvironments.CustomDomainConfiguration{
				DnsSuffix:           pointer.To(model.DnsSuffix),
				CertificateValue:    pointer.To(model.CertificateValue),
				CertificatePassword: pointer.To(model.CertificatePassword),
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, *id, *existing.Model, metadata.SetIDCallback(id)); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
