// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_custom_domain

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/containerapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/parse"
)

func (a Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ContainerAppClient

			id, err := parse.ContainerAppCustomDomainID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			containerAppId := containerapps.NewContainerAppID(id.SubscriptionId, id.ResourceGroupName, id.ContainerAppName)

			containerApp, err := client.Get(ctx, containerAppId)
			if err != nil || containerApp.Model == nil {
				return fmt.Errorf("retrieving %s to read %s", containerAppId, id)
			}

			model := containerApp.Model

			if model.Properties == nil || model.Properties.Configuration == nil || model.Properties.Configuration.Ingress == nil {
				return fmt.Errorf("could not read Ingress configuration for %s", containerAppId)
			}

			ingress := *model.Properties.Configuration.Ingress
			updatedCustomDomains := make([]containerapps.CustomDomain, 0)
			if customDomains := ingress.CustomDomains; customDomains != nil {
				for _, v := range *customDomains {
					if !strings.EqualFold(v.Name, id.CustomDomainName) {
						updatedCustomDomains = append(updatedCustomDomains, v)
					} else {
						// attempt to lock the cert if we have the ID
						certificateId := pointer.From(v.CertificateId)
						if certificateId != "" {
							locks.ByID(certificateId)
							defer locks.UnlockByID(certificateId)
						}
					}
				}
			}

			model.Properties.Configuration.Ingress.CustomDomains = pointer.To(updatedCustomDomains)

			// Delta-updates need the secrets back from the list API, or we'll end up removing them or erroring out.
			secretsResp, err := client.ListSecrets(ctx, containerAppId)
			if err != nil || secretsResp.Model == nil {
				if !response.WasStatusCode(secretsResp.HttpResponse, http.StatusNoContent) {
					return fmt.Errorf("retrieving secrets for update for %s: %+v", containerAppId, err)
				}
			}
			model.Properties.Configuration.Secrets = helpers.UnpackContainerSecretsCollection(secretsResp.Model)

			if err := client.CreateOrUpdateThenPoll(ctx, containerAppId, *model); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
