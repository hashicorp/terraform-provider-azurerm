// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package federated_identity_credential

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Required: true,
			ForceNew: true,
			Type:     pluginsdk.TypeString,
		},

		"user_assigned_identity_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"audience": {
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Required: true,
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
		},

		"issuer": {
			Required: true,
			Type:     pluginsdk.TypeString,
		},

		"subject": {
			Required: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
