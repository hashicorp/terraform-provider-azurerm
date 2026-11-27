// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_source_control_slot

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"slot_id": {
			Description:  "The ID of the Linux or Windows Web App Slot.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: webapps.ValidateSlotID,
		},

		"repo_url": {
			Description:  "The branch name to use for deployments.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{
				"branch",
			},
		},

		"branch": {
			Description:  "The URL for the repository.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{
				"repo_url",
			},
		},

		"use_local_git": {
			Description: "Whether the Slot should use local Git configuration.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
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
			Description: "Whether code should be deployed manually. Set to `true` to disable continuous integration, such as webhooks into online repos such as GitHub.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"github_action_configuration": helpers.GithubActionConfigSchema(),

		"use_mercurial": {
			Description: "The repository specified is Mercurial.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"rollback_enabled": {
			Description: "Whether the Deployment Rollback should be enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"scm_type": {
			Description: "The SCM Type in use. This value is decoded by the service from the repository information supplied.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"uses_github_action": {
			Description: "Whether the Slot uses a GitHub action for deployment. This value is decoded by the service from the repository information supplied.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},
	}
}
