// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package relay_namespace_authorization_rule

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/relay/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Azure Relay Namespace Authorization Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"namespace_name": {
			Description:  "Name of the Azure Relay Namespace for which this Azure Relay Namespace Authorization Rule will be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(6, 50),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"listen": {
			Description: "Grants listen access to this Authorization Rule.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"manage": {
			Description: "Grants manage access to this Authorization Rule. When this property is `true` - both `listen` and `send` must be set to `true` too.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"send": {
			Description: "Grants send access to this Authorization Rule.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"primary_connection_string": {
			Description: "The Primary Connection String for the Azure Relay Namespace Authorization Rule.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"primary_key": {
			Description: "The Primary Key for the Azure Relay Namespace Authorization Rule.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_connection_string": {
			Description: "The Secondary Connection String for the Azure Relay Namespace Authorization Rule.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_key": {
			Description: "The Secondary Key for the Azure Relay Namespace Authorization Rule.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			return helpers.AuthorizationRuleCustomizeDiff(ctx, metadata.ResourceDiff, nil)
		},
	}
}
