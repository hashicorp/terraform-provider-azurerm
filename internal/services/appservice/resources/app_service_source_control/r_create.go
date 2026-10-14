// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_service_source_control

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var appSourceControl SourceControlModel

			if err := metadata.Decode(&appSourceControl); err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient

			id, err := commonids.ParseWebAppID(appSourceControl.AppID)
			if err != nil {
				return err
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GetConfiguration(ctx, *id)
				if err != nil || existing.Model == nil || existing.Model.Properties == nil {
					return fmt.Errorf("checking for existing Source Control configuration on %s: %+v", id, err)
				}
				if pointer.From(existing.Model.Properties.ScmType) != webapps.ScmTypeNone {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			if appSourceControl.LocalGitSCM {
				sitePatch := webapps.SitePatchResource{
					Properties: &webapps.SitePatchResourceProperties{
						SiteConfig: &webapps.SiteConfig{
							ScmType: pointer.To(webapps.ScmTypeLocalGit),
						},
					},
				}

				if _, err := client.Update(ctx, *id, sitePatch); err != nil {
					return fmt.Errorf("setting App Source Control Type for %s: %v", id, err)
				}
			} else {
				app, err := client.Get(ctx, *id)
				if err != nil || app.Model == nil || app.Model.Kind == nil {
					return fmt.Errorf("reading site to determine O/S type for %s: %+v", id, err)
				}

				usesLinux := false
				if strings.Contains(strings.ToLower(*app.Model.Kind), "linux") {
					usesLinux = true
				}

				sourceControl := webapps.SiteSourceControl{
					Properties: &webapps.SiteSourceControlProperties{
						IsManualIntegration:       pointer.To(appSourceControl.ManualIntegration),
						DeploymentRollbackEnabled: pointer.To(appSourceControl.RollbackEnabled),
						IsMercurial:               pointer.To(appSourceControl.UseMercurial),
					},
				}

				if appSourceControl.RepoURL != "" {
					sourceControl.Properties.RepoURL = pointer.To(appSourceControl.RepoURL)
				}

				if appSourceControl.Branch != "" {
					sourceControl.Properties.Branch = pointer.To(appSourceControl.Branch)
				}

				if ghaConfig := helpers.ExpandGithubActionConfig(appSourceControl.GithubActionConfiguration, usesLinux); ghaConfig != nil {
					sourceControl.Properties.GitHubActionConfiguration = ghaConfig
				}

				if _, err = client.UpdateSourceControl(ctx, *id, sourceControl); err != nil {
					return fmt.Errorf("creating Source Control configuration for %s: %v", id, err)
				}
			}

			// TODO - Need to introduce polling for deployment statuses to avoid 409's elsewhere

			metadata.SetID(id)
			return nil
		},
	}
}
