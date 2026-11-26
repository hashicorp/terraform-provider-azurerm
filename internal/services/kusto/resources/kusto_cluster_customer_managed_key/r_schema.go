// Copyright IBM Corp. 2014, 2026
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
			Description:  "The ID of the Kusto Cluster.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKustoClusterID,
		},

		"key_vault_id": {
			Description:  "The ID of the Key Vault for CMK encryption.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
			ExactlyOneOf: []string{"managed_hsm_key_id", "key_vault_id"},
		},

		"managed_hsm_key_id": {
			Description:  "The Managed HSM Key ID for CMK encryption.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.Any(validate.ManagedHSMDataPlaneVersionedKeyID, validate.ManagedHSMDataPlaneVersionlessKeyID),
			ExactlyOneOf: []string{"managed_hsm_key_id", "key_vault_id"},
		},

		"key_name": {
			Description:  "The name of Key Vault Key.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"key_vault_id"},
		},

		"key_version": {
			Description:  "The version of Key Vault Key.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"user_identity": {
			Description:  "The User Assigned Identity that has access to the Key Vault Key. If not specified, System Assigned Identity will be used.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},
	}
}
