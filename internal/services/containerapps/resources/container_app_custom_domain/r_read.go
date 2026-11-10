// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_custom_domain

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/containerapps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/parse"
)

func (a Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ContainerAppClient
			state := ContainerAppCustomDomainResourceModel{}
			id, err := parse.ContainerAppCustomDomainID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			containerAppId := containerapps.NewContainerAppID(id.SubscriptionId, id.ResourceGroupName, id.ContainerAppName)

			containerApp, err := client.Get(ctx, containerAppId)
			if err != nil {
				if response.WasNotFound(containerApp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", containerAppId, err)
			}

			if containerApp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", containerAppId)
			}
			model := containerApp.Model

			if model.Properties == nil || model.Properties.Configuration == nil || model.Properties.Configuration.Ingress == nil {
				return fmt.Errorf("could not read Ingress configuration for %s", containerAppId)
			}

			ingress := *model.Properties.Configuration.Ingress
			found := false
			if customDomains := ingress.CustomDomains; customDomains != nil {
				for _, v := range *customDomains {
					if strings.EqualFold(v.Name, id.CustomDomainName) {
						found = true
						state.Name = id.CustomDomainName
						state.ContainerAppId = containerAppId.ID()
						if pointer.From(v.CertificateId) != "" {
							// The `v.CertificateId` returned from API has two possible values. when using an Azure created Managed Certificate,
							// its format is "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/managedEnvironments/%s/managedCertificates/%s",
							// another format is "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/managedEnvironments/%s/certificates/%s",
							// both cases are handled here to avoid parsing error.
							certId, err1 := managedenvironments.ParseCertificateIDInsensitively(pointer.From(v.CertificateId))
							if err1 != nil {
								managedCertId, err2 := managedenvironments.ParseManagedCertificateID(pointer.From(v.CertificateId))
								if err2 != nil {
									return err1
								}
								state.ManagedCertificateId = managedCertId.ID()
							} else {
								state.CertificateId = certId.ID()
							}
						}

						state.BindingType = pointer.FromEnum(v.BindingType)
					}
				}
			}

			if !found {
				return metadata.MarkAsGone(id)
			}

			return metadata.Encode(&state)
		},
	}
}
