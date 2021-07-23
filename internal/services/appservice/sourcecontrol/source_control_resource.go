package sourcecontrol

import (
	"context"
	"fmt"
	"strings"
	"time"

<<<<<<< HEAD:internal/services/appservice/sourcecontrol/source_control_resource.go
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-12-01/web"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
=======
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-01-15/web"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/sdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appservice/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appservice/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
>>>>>>> 7c6535352 (rebase and bump to API version 2021-01-15 now it is available):azurerm/internal/services/appservice/sourcecontrol/source_control_resource.go
)

type AppServiceSourceControlResource struct{}

type AppServiceSourceControlModel struct {
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

var _ sdk.Resource = AppServiceSourceControlResource{}

func (r AppServiceSourceControlResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppID,
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
		},

		"use_manual_integration": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"github_action_configuration": githubActionConfigSchema(),

		"use_mercurial": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"rollback_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},
	}
}

func (r AppServiceSourceControlResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"scm_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"uses_github_action": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
	}
}

func (r AppServiceSourceControlResource) ModelObject() interface{} {
	return AppServiceSourceControlModel{}
}

func (r AppServiceSourceControlResource) ResourceType() string {
	return "azurerm_app_service_source_control" // TODO - Does this name fit the new convention?
}

func (r AppServiceSourceControlResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var appSourceControl AppServiceSourceControlModel

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

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AppServiceSourceControlResource) Read() sdk.ResourceFunc {
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

			props := *appSourceControl.SiteSourceControlProperties

			state := AppServiceSourceControlModel{
				AppID:                     id.ID(),
				SCMType:                   string(siteConfig.ScmType),
				RepoURL:                   utils.NormalizeNilableString(props.RepoURL),
				Branch:                    utils.NormalizeNilableString(props.Branch),
				ManualIntegration:         *props.IsManualIntegration,
				UseMercurial:              *props.IsMercurial,
				RollbackEnabled:           *props.DeploymentRollbackEnabled,
				UsesGithubAction:          *props.IsGitHubAction,
				GithubActionConfiguration: flattenGitHubActionConfiguration(props.GitHubActionConfiguration),
				LocalGitSCM:               siteConfig.ScmType == web.ScmTypeLocalGit,
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AppServiceSourceControlResource) Delete() sdk.ResourceFunc {
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

func (r AppServiceSourceControlResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	// This is a meta resource with a 1:1 relationship with the service it's pointed at so we use the same ID
	return validate.WebAppID
}
