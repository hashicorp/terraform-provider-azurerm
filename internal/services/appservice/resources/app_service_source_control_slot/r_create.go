// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_service_source_control_slot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var appSourceControlSlot SourceControlSlotModel

			if err := metadata.Decode(&appSourceControlSlot); err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient

			id, err := webapps.ParseSlotID(appSourceControlSlot.SlotID)
			if err != nil {
				return err
			}

			appId := commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroupName, id.SiteName).ID()
			locks.ByID(appId)
			defer locks.UnlockByID(appId)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GetConfigurationSlot(ctx, *id)
				if err != nil || existing.Model == nil || existing.Model.Properties == nil {
					return fmt.Errorf("checking for existing Source Control configuration on %s: %+v", id, err)
				}
				if pointer.From(existing.Model.Properties.ScmType) != webapps.ScmTypeNone {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			if appSourceControlSlot.LocalGitSCM {
				sitePatch := webapps.SitePatchResource{
					Properties: &webapps.SitePatchResourceProperties{
						SiteConfig: &webapps.SiteConfig{
							ScmType: pointer.To(webapps.ScmTypeLocalGit),
						},
					},
				}

				if _, err := client.UpdateSlot(ctx, *id, sitePatch); err != nil {
					return fmt.Errorf("setting App Source Control Type for %s: %v", id, err)
				}
			} else {
				app, err := client.GetSlot(ctx, *id)
				if err != nil || app.Model == nil || app.Model.Kind == nil {
					return fmt.Errorf("reading slot to determine O/S type for %s: %+v", id, err)
				}

				usesLinux := false
				if strings.Contains(strings.ToLower(*app.Model.Kind), "linux") {
					usesLinux = true
				}

				sourceControl := webapps.SiteSourceControl{
					Properties: &webapps.SiteSourceControlProperties{
						IsManualIntegration:       pointer.To(appSourceControlSlot.ManualIntegration),
						DeploymentRollbackEnabled: pointer.To(appSourceControlSlot.RollbackEnabled),
						IsMercurial:               pointer.To(appSourceControlSlot.UseMercurial),
					},
				}

				if appSourceControlSlot.RepoURL != "" {
					sourceControl.Properties.RepoURL = pointer.To(appSourceControlSlot.RepoURL)
				}

				if appSourceControlSlot.Branch != "" {
					sourceControl.Properties.Branch = pointer.To(appSourceControlSlot.Branch)
				}

				if ghaConfig := helpers.ExpandGithubActionConfig(appSourceControlSlot.GithubActionConfiguration, usesLinux); ghaConfig != nil {
					sourceControl.Properties.GitHubActionConfiguration = ghaConfig
				}

				if _, err = client.UpdateSourceControlSlot(ctx, *id, sourceControl); err != nil {
					return fmt.Errorf("creating Source Control configuration for %s: %v", id, err)
				}
			}

			metadata.SetID(id)
			return nil
		},
	}
}
