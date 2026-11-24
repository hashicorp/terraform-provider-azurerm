// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_group

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementChildName("The name of the API Management Workspace Group."),

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&group.WorkspaceId{}),

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 300),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"external_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(group.GroupTypeCustom),
			// `group.GroupTypeSystem` is excluded here because system groups are predefined by API Management
			// and cannot be created through the API.
			ValidateFunc: validation.StringInSlice([]string{
				string(group.GroupTypeCustom),
				string(group.GroupTypeExternal),
			}, false),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff
			groupType := rd.Get("type").(string)
			isExternalIdUnset := rd.GetRawConfig().AsValueMap()["external_id"].IsNull()

			if groupType == string(group.GroupTypeExternal) && isExternalIdUnset {
				return errors.New("`external_id` must be specified when `type` is set to `external`")
			}
			if groupType == string(group.GroupTypeCustom) && !isExternalIdUnset {
				return errors.New("`type` must be set to `external` when `external_id` is specified`")
			}
			return nil
		},
		Timeout: 5 * time.Minute,
	}
}
