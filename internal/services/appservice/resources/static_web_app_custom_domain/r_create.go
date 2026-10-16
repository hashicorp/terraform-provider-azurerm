// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package static_web_app_custom_domain

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient

			model := StaticWebAppCustomDomainResourceModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			siteId, err := staticsites.ParseStaticSiteID(model.StaticSiteId)
			if err != nil {
				return err
			}

			id := staticsites.NewCustomDomainID(siteId.SubscriptionId, siteId.ResourceGroupName, siteId.StaticSiteName, model.DomainName)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GetStaticSiteCustomDomain(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			customDomain := staticsites.StaticSiteCustomDomainRequestPropertiesARMResource{
				Properties: &staticsites.StaticSiteCustomDomainRequestPropertiesARMResourceProperties{
					ValidationMethod: pointer.To(model.ValidationType),
				},
			}

			if strings.EqualFold(model.ValidationType, helpers.ValidationTypeCName) {
				if err = client.CreateOrUpdateStaticSiteCustomDomainCallbackThenPoll(ctx, id, customDomain, metadata.SetIDCallback(&id)); err != nil {
					return fmt.Errorf("creating %s: %+v", id, err)
				}
				metadata.SetID(id)
			} else {
				if _, err := client.CreateOrUpdateStaticSiteCustomDomain(ctx, id, customDomain); err != nil {
					return fmt.Errorf("creating %s: %+v", id, err)
				}

				metadata.SetID(id)

				deadline, ok := ctx.Deadline()
				if !ok {
					return fmt.Errorf("internal-error: context was missing a deadline")
				}

				stateConf := &pluginsdk.StateChangeConf{
					Pending: []string{
						string(staticsites.CustomDomainStatusRetrievingValidationToken),
					},
					Target: []string{
						string(staticsites.CustomDomainStatusValidating),
					},
					MinTimeout: 20 * time.Second,
					Timeout:    time.Until(deadline),
					Refresh: func() (any, string, error) {
						domain, err := client.GetStaticSiteCustomDomain(ctx, id)
						if err != nil {
							return domain, "Error", fmt.Errorf("retrieving %s: %+v", id, err)
						}

						if domain.Model == nil || domain.Model.Properties == nil {
							return nil, "Failed", fmt.Errorf("`properties` was missing from the response")
						}
						return domain, pointer.FromEnum(domain.Model.Properties.Status), nil
					},
				}

				if _, err := stateConf.WaitForStateContext(ctx); err != nil {
					return fmt.Errorf("waiting for DNS Validation after Creation of %s %+v", id, err)
				}
			}

			// Once validated the token value is zeroed,
			domain, err := client.GetStaticSiteCustomDomain(ctx, id)
			if err != nil {
				return fmt.Errorf("reading validation token for %s: %+v", id, err)
			}
			if m := domain.Model; m != nil {
				if m.Properties != nil {
					if err = metadata.ResourceData.Set("validation_token", m.Properties.ValidationToken); err != nil {
						return fmt.Errorf("setting validation_token value for %s: %+v", id, err)
					}
				}
			}

			return nil
		},
	}
}
