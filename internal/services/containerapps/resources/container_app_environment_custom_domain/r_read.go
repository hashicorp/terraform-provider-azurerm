// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_custom_domain

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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
