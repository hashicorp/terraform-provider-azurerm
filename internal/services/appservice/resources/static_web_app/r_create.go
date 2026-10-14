// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package static_web_app

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/sdkhacks"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			model := StaticWebAppResourceModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := staticsites.NewStaticSiteID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GetStaticSite(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			envelope := staticsites.StaticSiteARMResource{
				Location: location.Normalize(model.Location),
				Sku: &staticsites.SkuDescription{
					Name: pointer.To(model.SkuSize),
					Tier: pointer.To(model.SkuTier),
				},
				Tags: pointer.To(model.Tags),
			}

			ident, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding identity for %s: %+v", id, err)
			}
			if ident.Type != identity.TypeNone {
				envelope.Identity = ident
			}

			props := &staticsites.StaticSite{
				AllowConfigFileUpdates:   pointer.To(model.ConfigFileChanges),
				StagingEnvironmentPolicy: pointer.To(staticsites.StagingEnvironmentPolicyEnabled),
				PublicNetworkAccess:      pointer.To(helpers.PublicNetworkAccessEnabled),
			}

			if !model.PreviewEnvironments {
				props.StagingEnvironmentPolicy = pointer.To(staticsites.StagingEnvironmentPolicyDisabled)
			}

			// Check if repository URL, branch, or token are set
			if model.RepositoryUrl != "" || model.RepositoryBranch != "" || model.RepositoryToken != "" {
				props.Branch = pointer.To(model.RepositoryBranch)
				props.RepositoryURL = pointer.To(model.RepositoryUrl)
				props.RepositoryToken = pointer.To(model.RepositoryToken)
			}

			if !model.PublicNetworkAccess {
				props.PublicNetworkAccess = pointer.To(helpers.PublicNetworkAccessDisabled)
			}

			envelope.Properties = props

			if err := client.CreateOrUpdateStaticSiteCallbackThenPoll(ctx, id, envelope, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			if len(model.AppSettings) > 0 {
				appSettings := staticsites.StringDictionary{
					Properties: pointer.To(model.AppSettings),
				}

				if _, err = client.CreateOrUpdateStaticSiteAppSettings(ctx, id, appSettings); err != nil {
					return fmt.Errorf("updating app settings for %s: %+v", id, err)
				}
			}

			if len(model.BasicAuth) > 0 {
				sdkHackClient := sdkhacks.NewStaticWebAppClient(client)

				auth := model.BasicAuth[0]

				authProps := staticsites.StaticSiteBasicAuthPropertiesARMResource{
					Properties: &staticsites.StaticSiteBasicAuthPropertiesARMResourceProperties{
						ApplicableEnvironmentsMode: auth.Environments,
						Password:                   pointer.To(auth.Password),
						SecretState:                pointer.To("Password"),
					},
				}

				if _, err := sdkHackClient.CreateOrUpdateBasicAuth(ctx, id, authProps); err != nil {
					return fmt.Errorf("setting basic auth on %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}
