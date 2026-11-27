// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_object_replication

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/objectreplicationpolicyoperationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	storageClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageObjectReplicationUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := storageClient.From(meta).ResourceManager.ObjectReplicationPolicyOperationGroup
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := parse.ObjectReplicationID(d.Id())
	if err != nil {
		return err
	}

	srcAccount := commonids.NewStorageAccountID(id.Src.SubscriptionId, id.Src.ResourceGroupName, id.Src.StorageAccountName)
	dstAccount := commonids.NewStorageAccountID(id.Dst.SubscriptionId, id.Dst.ResourceGroupName, id.Dst.StorageAccountName)

	props := objectreplicationpolicyoperationgroup.ObjectReplicationPolicy{
		Properties: &objectreplicationpolicyoperationgroup.ObjectReplicationPolicyProperties{
			SourceAccount:      srcAccount.ID(),
			DestinationAccount: dstAccount.ID(),
			Rules:              expandArmObjectReplicationRuleArray(d.Get("rules").(*pluginsdk.Set).List()),
			Metrics: &objectreplicationpolicyoperationgroup.ObjectReplicationPolicyPropertiesMetrics{
				Enabled: pointer.To(d.Get("metrics_enabled").(bool)),
			},
		},
	}

	// update in dest storage account
	resp, err := client.ObjectReplicationPoliciesCreateOrUpdate(ctx, id.Dst, props)
	if err != nil {
		return fmt.Errorf("updating %q for destination storage account name %q: %+v", id, id.Dst.StorageAccountName, err)
	}
	if resp.Model == nil {
		return fmt.Errorf("nil model returned for Storage Object Replication for destination storage account name %q ID", id.Dst.StorageAccountName)
	}
	if resp.Model.Properties == nil {
		return fmt.Errorf("nil properties returned for Storage Object Replication for destination storage account name %q ID", id.Dst.StorageAccountName)
	}

	// update in source storage account, update policy Id and ruleId
	props.Properties.Rules = resp.Model.Properties.Rules
	if _, err := client.ObjectReplicationPoliciesCreateOrUpdate(ctx, id.Src, props); err != nil {
		return fmt.Errorf("updating %q for source storage account name %q: %+v", id, id.Src.StorageAccountName, err)
	}

	return resourceStorageObjectReplicationRead(d, meta)
}
