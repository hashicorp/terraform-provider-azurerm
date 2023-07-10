// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/web/2022-09-01/web"
)

type SourceControlResource struct{}

type SourceControlModel struct {
	AppID                     string                      `tfschema:"app_id"`
	SCMType                   string                      `tfschema:"scm_type"`
	RepoURL                   string                      `tfschema:"repo_url"`
	Branch                    string                      `tfschema:"branch"`
	LocalGitSCM               bool                        `tfschema:"use_local_git"`
	ManualIntegration         bool                        `tfschema:"use_manual_integration"`
	UseMercurial              bool                        `tfschema:"use_mercurial"`
	RollbackEnabled           bool                        `tfschema:"rollback_enabled"`
	UsesGithubAction          bool                        `tfschema:"uses_github_action"`
	GithubActionConfiguration []GithubActionConfiguration `tfschema:"github_action_configuration"`
}

var _ sdk.Resource = SourceControlResource{}

func (r SourceControlResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppID,
			Description:  "The ID of the Windows or Linux Web App.",
		},

		"repo_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{
				"branch",
			},
			Description: "The URL for the repository.",
		},

		"branch": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{
				"repo_url",
			},
			Description: "The branch name to use for deployments.",
		},

		"use_local_git": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
			ConflictsWith: []string{
				"repo_url",
				"branch",
				"use_manual_integration",
				"uses_github_action",
				"github_action_configuration",
				"use_mercurial",
				"rollback_enabled",
			},
			Description: "Should the App use local Git configuration.",
		},

		"use_manual_integration": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
			Description: "Should code be deployed manually. Set to `false` to enable continuous integration, such as webhooks into online repos such as GitHub. Defaults to `false`.",
		},

		"github_action_configuration": githubActionConfigSchema(),

		"use_mercurial": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
			Description: "The repository specified is Mercurial. Defaults to `false`.",
		},

		"rollback_enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
			Description: "Should the Deployment Rollback be enabled? Defaults to `false`.",
		},
	}
}

func (r SourceControlResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"scm_type": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The SCM Type in use. This value is decoded by the service from the repository information supplied.",
		},

		"uses_github_action": {
			Type:        pluginsdk.TypeBool,
			Computed:    true,
			Description: "Indicates if the Slot uses a GitHub action for deployment. This value is decoded by the service from the repository information supplied.",
		},
	}
}

func (r SourceControlResource) ModelObject() interface{} {
	return &SourceControlModel{}
}

func (r SourceControlResource) ResourceType() string {
	return "azurerm_app_service_source_control" // TODO - Does this name fit the new convention?
}

func (r SourceControlResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var appSourceControl SourceControlModel

			if err := metadata.Decode(&appSourceControl); err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient

			id, err := parse.WebAppID(appSourceControl.AppID)
			if err != nil {
				return err
			}

			existing, err := client.GetConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil || existing.SiteConfig == nil {
				return fmt.Errorf("checking for existing Source Control configuration on %s: %+v", id, err)
			}
			if existing.SiteConfig.ScmType != web.ScmTypeNone {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			if appSourceControl.LocalGitSCM {
				sitePatch := web.SitePatchResource{
					SitePatchResourceProperties: &web.SitePatchResourceProperties{
						SiteConfig: &web.SiteConfig{
							ScmType: web.ScmTypeLocalGit,
						},
					},
				}

				if _, err := client.Update(ctx, id.ResourceGroup, id.SiteName, sitePatch); err != nil {
					return fmt.Errorf("setting App Source Control Type for %s: %v", id, err)
				}
			} else {
				app, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
				if err != nil || app.Kind == nil {
					return fmt.Errorf("reading site to determine O/S type for %s: %+v", id, err)
				}

				usesLinux := false
				if strings.Contains(strings.ToLower(*app.Kind), "linux") {
					usesLinux = true
				}

				sourceControl := web.SiteSourceControl{
					SiteSourceControlProperties: &web.SiteSourceControlProperties{
						IsManualIntegration:       utils.Bool(appSourceControl.ManualIntegration),
						DeploymentRollbackEnabled: utils.Bool(appSourceControl.RollbackEnabled),
						IsMercurial:               utils.Bool(appSourceControl.UseMercurial),
					},
				}

				if appSourceControl.RepoURL != "" {
					sourceControl.SiteSourceControlProperties.RepoURL = utils.String(appSourceControl.RepoURL)
				}

				if appSourceControl.Branch != "" {
					sourceControl.SiteSourceControlProperties.Branch = utils.String(appSourceControl.Branch)
				}

				if ghaConfig := expandGithubActionConfig(appSourceControl.GithubActionConfiguration, usesLinux); ghaConfig != nil {
					sourceControl.SiteSourceControlProperties.GitHubActionConfiguration = ghaConfig
				}

				_, err = client.UpdateSourceControl(ctx, id.ResourceGroup, id.SiteName, sourceControl)
				if err != nil {
					return fmt.Errorf("creating Source Control configuration for %s: %v", id, err)
				}
			}

			// TODO - Need to introduce polling for deployment statuses to avoid 409's elsewhere

			metadata.SetID(id)
			return nil
		},
	}
}

func (r SourceControlResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.WebAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient

			appSourceControl, err := client.GetSourceControl(ctx, id.ResourceGroup, id.SiteName)
			if err != nil || appSourceControl.SiteSourceControlProperties == nil {
				if utils.ResponseWasNotFound(appSourceControl.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Source Control for %s: %v", id, err)
			}

			siteConfig, err := client.GetConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading App for Source Control %s: %v", id, err)
			}

			if siteConfig.ScmType == web.ScmTypeNone {
				metadata.Logger.Infof("App %s SCMType is `None` removing Source Control resource from state", id.SiteName)
				metadata.ResourceData.SetId("")
			}

			props := *appSourceControl.SiteSourceControlProperties

			state := SourceControlModel{
				AppID:                     id.ID(),
				SCMType:                   string(siteConfig.ScmType),
				RepoURL:                   utils.NormalizeNilableString(props.RepoURL),
				Branch:                    utils.NormalizeNilableString(props.Branch),
				ManualIntegration:         utils.NormaliseNilableBool(props.IsManualIntegration),
				UseMercurial:              utils.NormaliseNilableBool(props.IsMercurial),
				RollbackEnabled:           utils.NormaliseNilableBool(props.DeploymentRollbackEnabled),
				UsesGithubAction:          utils.NormaliseNilableBool(props.IsGitHubAction),
				GithubActionConfiguration: flattenGitHubActionConfiguration(props.GitHubActionConfiguration),
				LocalGitSCM:               siteConfig.ScmType == web.ScmTypeLocalGit,
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SourceControlResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := parse.WebAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			sitePatch := web.SitePatchResource{
				SitePatchResourceProperties: &web.SitePatchResourceProperties{
					SiteConfig: &web.SiteConfig{
						ScmType: web.ScmTypeNone,
					},
				},
			}
			if _, err := client.Update(ctx, id.ResourceGroup, id.SiteName, sitePatch); err != nil {
				return fmt.Errorf("setting App Source Control Type for %s: %v", id, err)
			}

			if _, err := client.DeleteSourceControl(ctx, id.ResourceGroup, id.SiteName, ""); err != nil {
				return fmt.Errorf("deleting Source Control for %s: %v", id, err)
			}

			return nil
		},
	}
}

func (r SourceControlResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	// This is a meta resource with a 1:1 relationship with the service it's pointed at so we use the same ID
	return validate.WebAppID
}
