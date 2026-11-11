// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_custom_domain

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			lawClient := metadata.Client.LogAnalytics.SharedKeyWorkspacesClient
			client := metadata.Client.ContainerApps.ManagedEnvironmentClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id, err := managedenvironments.ParseManagedEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			model := ContainerAppEnvironmentCustomDomainModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
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

			// If custom domain dns suffix or its certificate changed, update all the required attributes
			if metadata.ResourceData.HasChanges("dns_suffix", "certificate_blob_base64", "certificate_password") {
				existing.Model.Properties.CustomDomainConfiguration = &managedenvironments.CustomDomainConfiguration{
					DnsSuffix:           pointer.To(model.DnsSuffix),
					CertificateValue:    pointer.To(model.CertificateValue),
					CertificatePassword: pointer.To(model.CertificatePassword),
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
