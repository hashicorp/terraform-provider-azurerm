// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_management_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageManagementPolicyDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Storage.ResourceManager.ManagementPolicies
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	rid, err := parse.StorageAccountManagementPolicyID(d.Id())
	if err != nil {
		return err
	}

	accountId := commonids.NewStorageAccountID(rid.SubscriptionId, rid.ResourceGroup, rid.StorageAccountName)

	if _, err := client.Delete(ctx, accountId); err != nil {
		return fmt.Errorf("deleting %s: %+v", rid, err)
	}
	return nil
}
