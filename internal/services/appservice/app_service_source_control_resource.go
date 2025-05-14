// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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
			ValidateFunc: commonids.ValidateWebAppID,
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

			id, err := commonids.ParseWebAppID(appSourceControl.AppID)
			if err != nil {
				return err
			}

			existing, err := client.GetConfiguration(ctx, *id)
			if err != nil || existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("checking for existing Source Control configuration on %s: %+v", id, err)
			}
			if pointer.From(existing.Model.Properties.ScmType) != webapps.ScmTypeNone {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
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

				if ghaConfig := expandGithubActionConfig(appSourceControl.GithubActionConfiguration, usesLinux); ghaConfig != nil {
					sourceControl.Properties.GitHubActionConfiguration = ghaConfig
				}

				_, err = client.UpdateSourceControl(ctx, *id, sourceControl)
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
					SCMType:                   string(pointer.From(siteConfig.Model.Properties.ScmType)),
					RepoURL:                   pointer.From(props.RepoURL),
					Branch:                    pointer.From(props.Branch),
					ManualIntegration:         pointer.From(props.IsManualIntegration),
					UseMercurial:              pointer.From(props.IsMercurial),
					RollbackEnabled:           pointer.From(props.DeploymentRollbackEnabled),
					UsesGithubAction:          pointer.From(props.IsGitHubAction),
					GithubActionConfiguration: flattenGitHubActionConfiguration(props.GitHubActionConfiguration),
					LocalGitSCM:               pointer.From(siteConfig.Model.Properties.ScmType) == webapps.ScmTypeLocalGit,
				}
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
			id, err := commonids.ParseWebAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			sitePatch := webapps.SitePatchResource{
				Properties: &webapps.SitePatchResourceProperties{
					SiteConfig: &webapps.SiteConfig{
						ScmType: pointer.To(webapps.ScmTypeNone),
					},
				},
			}
			if _, err := client.Update(ctx, *id, sitePatch); err != nil {
				return fmt.Errorf("setting App Source Control Type for %s: %v", id, err)
			}

			if _, err := client.DeleteSourceControl(ctx, *id, webapps.DefaultDeleteSourceControlOperationOptions()); err != nil {
				return fmt.Errorf("deleting Source Control for %s: %v", id, err)
			}

			return nil
		},
	}
}

func (r SourceControlResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	// This is a meta resource with a 1:1 relationship with the service it's pointed at so we use the same ID
	return commonids.ValidateAppServiceID
}
