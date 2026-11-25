// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_source_control

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/sourcecontrol"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (m Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"automation_account_id": {
			Description:  "The ID of Automation Account to manage this Source Control.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: sourcecontrol.ValidateAutomationAccountID,
		},

		"name": {
			Description:  "The name which should be used for this Automation Source Control.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},

		"repository_url": {
			Description:  "The Repository URL of the source control.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(0, 2000),
		},

		"branch": {
			Description:  "The repo branch of the Source Control. Empty value is valid only for `VsoTfvc`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(0, 255),
		},

		"folder_path": {
			Description:  "The folder path of the source control. This Path must be relative.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(0, 255),
		},

		"automatic_sync": {
			Description: "Whether to automatically sync the Source Control.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"publish_runbook_enabled": {
			Description: "Whether auto publish the Source Control.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"source_control_type": {
			Description:  "The source type of Source Control.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(sourcecontrol.PossibleValuesForSourceType(), true),
		},

		"description": {
			Description: "A short description of the Source Control.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"security": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"token": {
						Description:  "The access token of specified repo.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(0, 1024),
					},

					"refresh_token": {
						Description:  "The refresh token of specified repository.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringLenBetween(0, 1024),
					},

					"token_type": {
						Description:  "The token type.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(sourcecontrol.PossibleValuesForTokenType(), false),
					},
				},
			},
		},
	}
}

func (m Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
