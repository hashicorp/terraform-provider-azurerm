// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_service_source_control

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := commonids.ParseWebAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient

			appSourceControl, err := client.GetSourceControl(ctx, *id)
			if err != nil || appSourceControl.Model == nil || appSourceControl.Model.Properties == nil {
				if response.WasNotFound(appSourceControl.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Source Control for %s: %v", id, err)
			}

			siteConfig, err := client.GetConfiguration(ctx, *id)
			if err != nil || siteConfig.Model == nil || siteConfig.Model.Properties == nil {
				return fmt.Errorf("reading App for Source Control %s: %v", id, err)
			}

			if pointer.From(siteConfig.Model.Properties.ScmType) == webapps.ScmTypeNone {
				metadata.Logger.Infof("App %s SCMType is `None` removing Source Control resource from state", id.SiteName)
				metadata.ResourceData.SetId("")
			}

			state := SourceControlModel{}
			if model := appSourceControl.Model; model != nil {
				props := model.Properties
				state = SourceControlModel{
					AppID:                     id.ID(),
					SCMType:                   pointer.FromEnum(siteConfig.Model.Properties.ScmType),
					RepoURL:                   pointer.From(props.RepoURL),
					Branch:                    pointer.From(props.Branch),
					ManualIntegration:         pointer.From(props.IsManualIntegration),
					UseMercurial:              pointer.From(props.IsMercurial),
					RollbackEnabled:           pointer.From(props.DeploymentRollbackEnabled),
					UsesGithubAction:          pointer.From(props.IsGitHubAction),
					GithubActionConfiguration: helpers.FlattenGitHubActionConfiguration(props.GitHubActionConfiguration),
					LocalGitSCM:               pointer.From(siteConfig.Model.Properties.ScmType) == webapps.ScmTypeLocalGit,
				}
			}

			return metadata.Encode(&state)
		},
	}
}
