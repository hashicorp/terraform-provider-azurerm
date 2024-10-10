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
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SourceControlSlotResource struct{}

type SourceControlSlotModel struct {
	SlotID                    string                      `tfschema:"slot_id"`
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

var _ sdk.Resource = SourceControlSlotResource{}

func (r SourceControlSlotResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"slot_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: webapps.ValidateSlotID,
			Description:  "The ID of the Linux or Windows Web App Slot.",
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
			Description: "The branch name to use for deployments.",
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
			Description: "The URL for the repository",
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
			Description: "Should the Slot use local Git configuration.",
		},

		"use_manual_integration": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
			Description: "Should code be deployed manually. Set to `true` to disable continuous integration, such as webhooks into online repos such as GitHub. Defaults to `false`",
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
			Description: "Should the Deployment Rollback be enabled? Defaults to `false`",
		},
	}
}

func (r SourceControlSlotResource) Attributes() map[string]*pluginsdk.Schema {
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

func (r SourceControlSlotResource) ModelObject() interface{} {
	return &SourceControlSlotModel{}
}

func (r SourceControlSlotResource) ResourceType() string {
	return "azurerm_app_service_source_control_slot" // TODO - Does this name fit the new convention?
}

func (r SourceControlSlotResource) Create() sdk.ResourceFunc {
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

			existing, err := client.GetConfigurationSlot(ctx, *id)
			if err != nil || existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("checking for existing Source Control configuration on %s: %+v", id, err)
			}
			if pointer.From(existing.Model.Properties.ScmType) != webapps.ScmTypeNone {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
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
					sourceControl.Properties.RepoURL = utils.String(appSourceControlSlot.RepoURL)
				}

				if appSourceControlSlot.Branch != "" {
					sourceControl.Properties.Branch = utils.String(appSourceControlSlot.Branch)
				}

				if ghaConfig := expandGithubActionConfig(appSourceControlSlot.GithubActionConfiguration, usesLinux); ghaConfig != nil {
					sourceControl.Properties.GitHubActionConfiguration = ghaConfig
				}

				_, err = client.UpdateSourceControlSlot(ctx, *id, sourceControl)
				if err != nil {
					return fmt.Errorf("creating Source Control configuration for %s: %v", id, err)
				}
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r SourceControlSlotResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := webapps.ParseSlotID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient

			appSourceControl, err := client.GetSourceControlSlot(ctx, *id)
			if err != nil || appSourceControl.Model == nil || appSourceControl.Model.Properties == nil {
				if response.WasNotFound(appSourceControl.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Source Control for %s: %v", id, err)
			}

			siteConfig, err := client.GetConfigurationSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading App for Source Control %s: %v", id, err)
			}

			if pointer.From(siteConfig.Model.Properties.ScmType) == webapps.ScmTypeNone {
				metadata.Logger.Infof("App %s SCMType is `None` removing Source Control resource from state", id.SiteName)
				metadata.ResourceData.SetId("")
			}

			props := *appSourceControl.Model.Properties

			state := SourceControlSlotModel{
				SlotID:                    id.ID(),
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

			return metadata.Encode(&state)
		},
	}
}

func (r SourceControlSlotResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := webapps.ParseSlotID(metadata.ResourceData.Id())
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
			if _, err := client.UpdateSlot(ctx, *id, sitePatch); err != nil {
				return fmt.Errorf("setting App Source Control Type for %s: %v", id, err)
			}

			if _, err := client.DeleteSourceControlSlot(ctx, *id, webapps.DefaultDeleteSourceControlSlotOperationOptions()); err != nil {
				return fmt.Errorf("deleting Source Control for %s: %v", id, err)
			}

			return nil
		},
	}
}

func (r SourceControlSlotResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	// This is a meta resource with a 1:1 relationship with the slot it's pointed at, so we use the same ID
	return webapps.ValidateSlotID
}
