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
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/parse"
)

func (a Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ContainerAppClient

			model := ContainerAppCustomDomainResourceModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			containerAppId, err := containerapps.ParseContainerAppID(model.ContainerAppId)
			if err != nil {
				return err
			}

			locks.ByID(containerAppId.ID())
			defer locks.UnlockByID(containerAppId.ID())

			id := parse.NewContainerAppCustomDomainId(containerAppId.SubscriptionId, containerAppId.ResourceGroupName, containerAppId.ContainerAppName, model.Name)

			var certificateId *managedenvironments.CertificateId
			if model.CertificateId != "" {
				certificateId, err = managedenvironments.ParseCertificateID(model.CertificateId)
				if err != nil {
					return err
				}
			}

			containerApp, err := client.Get(ctx, *containerAppId)
			if err != nil || containerApp.Model == nil {
				return fmt.Errorf("retrieving %s to create %s", containerAppId, id)
			}

			props := containerApp.Model.Properties
			if props == nil || props.Configuration == nil {
				return fmt.Errorf("could not retrieve properties of %s", containerAppId)
			}

			config := *props.Configuration

			if config.Ingress == nil {
				return fmt.Errorf("specified Container App (%s) has no Ingress configuration for Custom Domains", containerAppId)
			}

			// Delta-updates need the secrets back from the list API, or we'll end up removing them or erroring out.
			secretsResp, err := client.ListSecrets(ctx, *containerAppId)
			if err != nil || secretsResp.Model == nil {
				if !response.WasStatusCode(secretsResp.HttpResponse, http.StatusNoContent) {
					return fmt.Errorf("retrieving secrets for update for %s: %+v", *containerAppId, err)
				}
			}
			props.Configuration.Secrets = helpers.UnpackContainerSecretsCollection(secretsResp.Model)

			ingress := *config.Ingress

			customDomains := make([]containerapps.CustomDomain, 0)
			exists := false
			if existingCustomDomains := ingress.CustomDomains; existingCustomDomains != nil {
				for _, v := range *existingCustomDomains {
					if strings.EqualFold(v.Name, model.Name) {
						exists = true
						if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
							return metadata.ResourceRequiresImport(Resource{}.ResourceType(), id)
						}
					}
				}

				customDomains = *existingCustomDomains
			}

			customDomain := containerapps.CustomDomain{
				Name:        model.Name,
				BindingType: pointer.To(containerapps.BindingTypeDisabled),
			}

			if certificateId != nil {
				customDomain.CertificateId = pointer.To(certificateId.ID())
				customDomain.BindingType = pointer.ToEnum[containerapps.BindingType](model.BindingType)
			}

			if !exists {
				customDomains = append(customDomains, customDomain)
			}

			containerApp.Model.Properties.Configuration.Ingress.CustomDomains = pointer.To(customDomains)

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, *containerAppId, *containerApp.Model, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
