// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_customer_managed_key

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"data_factory_id": {
			Description:  "The ID of the Data Factory Resource the Customer Managed Key will be associated with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},
		"customer_managed_key_id": {
			Description:  "The ID the of the Customer Managed Key to associate with the Data Factory.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
		},
		"user_assigned_identity_id": {
			Description:  "The User Assigned Identity ID that will be used to access Key Vaults that contain the encryption keys.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
