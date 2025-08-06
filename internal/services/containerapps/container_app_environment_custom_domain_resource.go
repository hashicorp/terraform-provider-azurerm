// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/managedenvironments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ContainerAppEnvironmentCustomDomainResource struct{}

type ContainerAppEnvironmentCustomDomainModel struct {
	ManagedEnvironmentId string `tfschema:"container_app_environment_id"`

	CertificatePassword string `tfschema:"certificate_password"`
	CertificateValue    string `tfschema:"certificate_blob_base64"`
	DnsSuffix           string `tfschema:"dns_suffix"`
}

var _ sdk.ResourceWithUpdate = ContainerAppEnvironmentCustomDomainResource{}

func (r ContainerAppEnvironmentCustomDomainResource) ModelObject() interface{} {
	return &ContainerAppEnvironmentCustomDomainModel{}
}

func (r ContainerAppEnvironmentCustomDomainResource) ResourceType() string {
	return "azurerm_container_app_environment_custom_domain"
}

func (r ContainerAppEnvironmentCustomDomainResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return managedenvironments.ValidateManagedEnvironmentID
}

func (r ContainerAppEnvironmentCustomDomainResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedenvironments.ValidateManagedEnvironmentID,
			Description:  "The Container App Managed Environment ID to configure this Custom Domain on.",
		},

		"certificate_blob_base64": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsBase64,
			Description:  "The Custom Domain Certificate Private Key as a base64 encoded PFX or PEM.",
		},

		"certificate_password": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "The Custom Domain Certificate password.",
		},

		"dns_suffix": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The Custom Domain DNS suffix for this Container App Environment.",
		},
	}
}

func (r ContainerAppEnvironmentCustomDomainResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ContainerAppEnvironmentCustomDomainResource) Create() sdk.ResourceFunc {
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
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			// Check if this resource needs import
			if customDomain := existing.Model.Properties.CustomDomainConfiguration; customDomain != nil && customDomain.DnsSuffix != nil {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
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

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ContainerAppEnvironmentCustomDomainResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ManagedEnvironmentClient
			id, err := managedenvironments.ParseManagedEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state := ContainerAppEnvironmentCustomDomainModel{}

			if model := existing.Model; model != nil {
				if props := model.Properties; props != nil {
					if customdomain := props.CustomDomainConfiguration; customdomain.DnsSuffix != nil {
						state.DnsSuffix = pointer.From(customdomain.DnsSuffix)
						if certValue, ok := metadata.ResourceData.GetOk("certificate_blob_base64"); ok {
							state.CertificateValue = certValue.(string)
						}
						if certPassword, ok := metadata.ResourceData.GetOk("certificate_password"); ok {
							state.CertificatePassword = certPassword.(string)
						}
						state.ManagedEnvironmentId = metadata.ResourceData.Id()
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ContainerAppEnvironmentCustomDomainResource) Delete() sdk.ResourceFunc {
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

			existing.Model.Properties.CustomDomainConfiguration = nil

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ContainerAppEnvironmentCustomDomainResource) Update() sdk.ResourceFunc {
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
			if metadata.ResourceData.HasChange("dns_suffix") ||
				metadata.ResourceData.HasChange("certificate_blob_base64") ||
				metadata.ResourceData.HasChange("certificate_password") {
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

func findLogAnalyticsWorkspaceSecret(ctx context.Context, client *workspaces.WorkspacesClient, subscriptionId, targetCustomerId string) (string, error) {
	parsedSubscriptionId := commonids.NewSubscriptionID(subscriptionId)

	resp, err := client.List(ctx, parsedSubscriptionId)
	if err != nil {
		return "", err
	}

	if resp.Model == nil {
		return "", fmt.Errorf("model was nil")
	}

	if resp.Model.Value == nil {
		return "", fmt.Errorf("value was nil")
	}

	for _, law := range *resp.Model.Value {
		if law.Properties != nil && law.Properties.CustomerId != nil && *law.Properties.CustomerId == targetCustomerId && law.Id != nil {
			id, err := workspaces.ParseWorkspaceIDInsensitively(*law.Id)
			if err != nil {
				return "", fmt.Errorf("parsing ID or %s: %+v", *law.Id, err)
			}
			keys, err := client.SharedKeysGetSharedKeys(ctx, *id)
			if err != nil {
				return "", fmt.Errorf("retrieving access keys for %s: %+v", *law.Id, err)
			}
			if keys.Model.PrimarySharedKey == nil {
				return "", fmt.Errorf("reading shared key for %s", *law.Id)
			}
			return *keys.Model.PrimarySharedKey, nil
		}
	}

	return "", fmt.Errorf("no matching workspace found")
}
