// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_named_value

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the API Management Workspace Named Value.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z](?:[a-zA-Z0-9-]{0,254}[a-zA-Z0-9])?$`),
				"`name` must be 1-256 characters, starting with a letter, and using only letters, numbers, or hyphens",
			),
		},

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&workspace.WorkspaceId{}, "API Management Workspace"),

		"display_name": {
			Description: "The display name of this API Management Workspace Named Value.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[A-Za-z0-9-._]{1,256}$`),
				"`display_name` must be 1-256 characters and can only contain letters, numbers, hyphens, periods, and underscores",
			),
		},

		"secret": {
			Description: "Whether the API Management Workspace Named Value is secret.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		// NOTE: this is not the common `tags` attribute used by most Azure resources, this field accepts a list of strings rather than a map of key-value pairs
		"tags": {
			Description: "A list of tags to be applied to the API Management Workspace Named Value.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"value": {
			Description:  "The value of this API Management Workspace Named Value.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"value", "value_from_key_vault"},
		},

		"value_from_key_vault": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"value", "value_from_key_vault"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"secret_id": {
						Description:  "The resource ID of the Key Vault Secret.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret),
					},
					"identity_client_id": {
						Description:  "The client ID of the User Assigned Identity, for the API Management Service, which will be used to access the Key Vault Secret. The System Assigned Identity will be used if not specified.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ApiManagementWorkspaceNamedValueModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if len(model.ValueFromKeyVault) > 0 && !model.Secret {
				return errors.New("`secret` must be set to `true` when `value_from_key_vault` is specified")
			}

			return nil
		},
	}
}
