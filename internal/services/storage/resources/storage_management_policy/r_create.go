// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_management_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/managementpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageManagementPolicyCreateOrUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Storage.ResourceManager.ManagementPolicies
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	rid, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	// The name of the Storage Account Management Policy. It should always be 'default' (from https://learn.microsoft.com/rest/api/storagerp/managementpolicies/createorupdate)
	mgmtPolicyId := parse.NewStorageAccountManagementPolicyID(rid.SubscriptionId, rid.ResourceGroupName, rid.StorageAccountName, "default")

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, *rid)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", mgmtPolicyId, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_storage_management_policy", mgmtPolicyId.ID())
			}
		}
	}

	parameters := managementpolicies.ManagementPolicy{
		Name: &mgmtPolicyId.ManagementPolicyName,
	}

	armRules, err := expandStorageManagementPolicyRules(d)
	if err != nil {
		return fmt.Errorf("expanding %s: %+v", mgmtPolicyId, err)
	}

	parameters.Properties = &managementpolicies.ManagementPolicyProperties{
		Policy: managementpolicies.ManagementPolicySchema{
			Rules: armRules,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, *rid, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", mgmtPolicyId, err)
	}

	d.SetId(mgmtPolicyId.ID())

	return resourceStorageManagementPolicyRead(d, meta)
}
