// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_user

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/user"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementUserSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"user_id": schemaz.SchemaApiManagementUserName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"first_name": {
			Description:  "The first name for this user.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"email": {
			Description:  "The email address associated with this user.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"last_name": {
			Description:  "The last name for this user.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"confirmation": {
			Description:  "The kind of confirmation email which will be sent to this user.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(user.PossibleValuesForConfirmation(), false),
		},

		"note": {
			Description: "A note about this user.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"password": {
			Description: "The password associated with this user.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Sensitive:   true,
		},

		"state": {
			Description: "The state of this user.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice([]string{
				string(user.UserStateActive),
				string(user.UserStateBlocked),
				string(user.UserStatePending),
			}, false),
		},
	}
}
