// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package static_web_app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/sdkhacks"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient

			config := StaticWebAppResourceModel{}

			if err := metadata.Decode(&config); err != nil {
				return err
			}

			id, err := staticsites.ParseStaticSiteID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.GetStaticSite(ctx, *id)
			if err != nil || existing.Model == nil {
				return fmt.Errorf("retrieving %s for update: %+v", *id, err)
			}

			model := *existing.Model

			if metadata.ResourceData.HasChange("identity") {
				ident, err := identity.ExpandSystemAndUserAssignedMapFromModel(config.Identity)
				if err != nil {
					return err
				}
				model.Identity = ident
				// If we're changing to `Free` we must remove the Identity first
				if strings.EqualFold(string(resourceproviders.SkuNameFree), config.SkuTier) {
					if err := client.CreateOrUpdateStaticSiteThenPoll(ctx, *id, model); err != nil {
						return fmt.Errorf("creating %s: %+v", id, err)
					}
					// Once removed, the identity payload needs to be nilled or the API validation for `Free` will reject the request
					model.Identity = nil
				}
			}

			if metadata.ResourceData.HasChanges("sku_tier", "sku_size") {
				model.Sku = &staticsites.SkuDescription{
					Name: pointer.To(config.SkuSize),
					Tier: pointer.To(config.SkuTier),
				}
			}

			if metadata.ResourceData.HasChange("configuration_file_changes_enabled") {
				model.Properties.AllowConfigFileUpdates = pointer.To(config.ConfigFileChanges)
			}

			if metadata.ResourceData.HasChange("preview_environments_enabled") {
				if !config.PreviewEnvironments {
					model.Properties.StagingEnvironmentPolicy = pointer.To(staticsites.StagingEnvironmentPolicyDisabled)
				} else {
					model.Properties.StagingEnvironmentPolicy = pointer.To(staticsites.StagingEnvironmentPolicyEnabled)
				}
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				if !config.PublicNetworkAccess {
					model.Properties.PublicNetworkAccess = pointer.To(helpers.PublicNetworkAccessDisabled)
				} else {
					model.Properties.PublicNetworkAccess = pointer.To(helpers.PublicNetworkAccessEnabled)
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = pointer.To(config.Tags)
			}

			if err := client.CreateOrUpdateStaticSiteThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("app_settings") {
				appSettings := staticsites.StringDictionary{
					Properties: pointer.To(config.AppSettings),
				}

				if _, err = client.CreateOrUpdateStaticSiteAppSettings(ctx, *id, appSettings); err != nil {
					return fmt.Errorf("updating app settings for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("basic_auth") {
				sdkHackClient := sdkhacks.NewStaticWebAppClient(client)
				authProps := staticsites.StaticSiteBasicAuthPropertiesARMResource{}
				if len(config.BasicAuth) > 0 {
					auth := config.BasicAuth[0]
					authProps.Properties = &staticsites.StaticSiteBasicAuthPropertiesARMResourceProperties{
						ApplicableEnvironmentsMode: auth.Environments,
						Password:                   pointer.To(auth.Password),
						SecretState:                pointer.To("Password"),
					}
				} else {
					authProps.Properties = &staticsites.StaticSiteBasicAuthPropertiesARMResourceProperties{
						ApplicableEnvironmentsMode: "SpecifiedEnvironments",
						SecretState:                pointer.To(""),
					}
				}

				if _, err := sdkHackClient.CreateOrUpdateBasicAuth(ctx, *id, authProps); err != nil {
					return fmt.Errorf("setting basic auth on %s: %+v", *id, err)
				}
			}

			if metadata.ResourceData.HasChanges("repository_url", "repository_branch", "repository_token") {
				model.Properties.RepositoryURL = pointer.To(config.RepositoryUrl)
				model.Properties.Branch = pointer.To(config.RepositoryBranch)
				model.Properties.RepositoryToken = pointer.To(config.RepositoryToken)
			}

			return nil
		},
	}
}
