// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_service_source_control

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
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
			Computed:     true, // azignore:AZS007 - pre-existing violation
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
			Computed:     true, // azignore:AZS007 - pre-existing violation
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

		"github_action_configuration": helpers.GithubActionConfigSchema(),

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

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
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
