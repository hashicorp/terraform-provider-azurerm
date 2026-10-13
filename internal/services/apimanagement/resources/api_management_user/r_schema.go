// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"email": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"last_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"confirmation": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(user.PossibleValuesForConfirmation(), false),
		},

		"note": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"password": {
			Type:      pluginsdk.TypeString,
			Optional:  true,
			Sensitive: true,
		},

		"state": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice([]string{
				string(user.UserStateActive),
				string(user.UserStateBlocked),
				string(user.UserStatePending),
			}, false),
		},
	}
}
