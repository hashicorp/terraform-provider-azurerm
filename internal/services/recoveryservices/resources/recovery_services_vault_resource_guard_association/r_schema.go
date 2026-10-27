// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package recovery_services_vault_resource_guard_association

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (r Resource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&vaults.VaultId{}),

		"resource_guard_id": commonschema.ResourceIDReferenceRequiredForceNew(&resourceguardresources.ResourceGuardId{}),
	}
}

func (r Resource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}
