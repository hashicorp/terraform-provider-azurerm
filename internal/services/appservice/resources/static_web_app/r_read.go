// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package static_web_app

import (
	"context"
	"fmt"
	"strings"
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

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient

			id, err := staticsites.ParseStaticSiteID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			staticSite, err := client.GetStaticSite(ctx, *id)
			if err != nil {
				if response.WasNotFound(staticSite.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := StaticWebAppResourceModel{
				Name:              id.StaticSiteName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := staticSite.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				ident, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening identity for %s: %+v", *id, err)
				}
				state.Identity = pointer.From(ident)

				state.Tags = pointer.From(model.Tags)
				if props := model.Properties; props != nil {
					state.ConfigFileChanges = pointer.From(props.AllowConfigFileUpdates)
					state.DefaultHostName = pointer.From(props.DefaultHostname)
					state.PreviewEnvironments = pointer.From(props.StagingEnvironmentPolicy) == staticsites.StagingEnvironmentPolicyEnabled

					state.RepositoryUrl = pointer.From(props.RepositoryURL)
					state.RepositoryBranch = pointer.From(props.Branch)

					// Token isn't returned in the response, so we need to grab it from the config
					if repositoryToken, ok := metadata.ResourceData.GetOk("repository_token"); ok {
						state.RepositoryToken = repositoryToken.(string)
					}

					state.PublicNetworkAccess = !strings.EqualFold(pointer.From(props.PublicNetworkAccess), helpers.PublicNetworkAccessDisabled)
				}

				if sku := model.Sku; sku != nil {
					state.SkuSize = pointer.From(sku.Name)
					state.SkuTier = pointer.From(sku.Tier)
				}

				sec, err := client.ListStaticSiteSecrets(ctx, *id)
				if err != nil || sec.Model == nil {
					return fmt.Errorf("retrieving secrets for %s: %+v", *id, err)
				}

				if secProps := sec.Model.Properties; secProps != nil {
					state.ApiKey = pointer.From(secProps)["apiKey"]
				}
			}

			appSettings, err := client.ListStaticSiteAppSettings(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving app_settings for %s: %+v", *id, err)
			}
			if appSettingsModel := appSettings.Model; appSettingsModel != nil {
				state.AppSettings = pointer.From(appSettingsModel.Properties)
			}

			sdkHackClient := sdkhacks.NewStaticWebAppClient(client)
			auth, err := sdkHackClient.GetBasicAuth(ctx, *id)
			if err != nil && !response.WasNotFound(auth.HttpResponse) { // If basic auth is not configured then this 404's
				return fmt.Errorf("retrieving auth config for %s: %+v", *id, err)
			}
			if !response.WasNotFound(auth.HttpResponse) {
				if authModel := auth.Model; authModel != nil && authModel.Properties != nil && !strings.EqualFold(authModel.Properties.ApplicableEnvironmentsMode, EnvironmentsTypeSpecifiedEnvironments) {
					state.BasicAuth = []BasicAuth{
						{
							Password:     metadata.ResourceData.Get("basic_auth.0.password").(string), // Grab this from the config, if we can.
							Environments: authModel.Properties.ApplicableEnvironmentsMode,
						},
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}
