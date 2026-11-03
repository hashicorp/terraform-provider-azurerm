// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster_customer_managed_key

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func kustoClusterCustomerManagedKeySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKustoClusterID,
		},

		"key_vault_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
			ExactlyOneOf: []string{"managed_hsm_key_id", "key_vault_id"},
		},

		"managed_hsm_key_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.Any(validate.ManagedHSMDataPlaneVersionedKeyID, validate.ManagedHSMDataPlaneVersionlessKeyID),
			ExactlyOneOf: []string{"managed_hsm_key_id", "key_vault_id"},
		},

		"key_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"key_vault_id"},
		},

		"key_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"user_identity": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},
	}
}
