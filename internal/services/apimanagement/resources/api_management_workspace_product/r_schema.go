// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_product

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,78}[a-zA-Z0-9])?$`),
				"`name` must be 1–80 characters, using only letters, numbers, or hyphens, and not starting or ending with a hyphen."),
		},

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&workspace.WorkspaceId{}),

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"published_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"require_approval_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"require_subscription_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"subscriptions_limit": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"terms": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if !metadata.ResourceDiff.Get("require_subscription_enabled").(bool) {
				if !metadata.ResourceDiff.GetRawConfig().AsValueMap()["subscriptions_limit"].IsNull() {
					return fmt.Errorf("`require_subscription_enabled` must be set to `true` when `subscriptions_limit` is specified")
				}

				if !metadata.ResourceDiff.GetRawConfig().AsValueMap()["require_approval_enabled"].IsNull() {
					return fmt.Errorf("`require_subscription_enabled` must be set to `true` when `require_approval_enabled` is specified")
				}
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
