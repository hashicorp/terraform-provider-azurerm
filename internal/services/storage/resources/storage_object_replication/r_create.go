// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_object_replication

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/objectreplicationpolicyoperationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageObjectReplicationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Storage.ResourceManager.ObjectReplicationPolicyOperationGroup
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	srcAccount, err := commonids.ParseStorageAccountID(d.Get("source_storage_account_id").(string))
	if err != nil {
		return err
	}
	dstAccount, err := commonids.ParseStorageAccountID(d.Get("destination_storage_account_id").(string))
	if err != nil {
		return err
	}

	srcId := objectreplicationpolicyoperationgroup.NewObjectReplicationPolicyID(srcAccount.SubscriptionId, srcAccount.ResourceGroupName, srcAccount.StorageAccountName, "default")
	dstId := objectreplicationpolicyoperationgroup.NewObjectReplicationPolicyID(dstAccount.SubscriptionId, dstAccount.ResourceGroupName, dstAccount.StorageAccountName, "default")

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		resp, err := client.ObjectReplicationPoliciesList(ctx, *dstAccount)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for present of existing Storage Object Replication for destination %q): %+v", dstAccount, err)
			}
		}
		if resp.Model != nil {
			for _, existing := range *resp.Model {
				if existing.Name != nil && *existing.Name != "" {
					if prop := existing.Properties; prop != nil && (
					// Storage allows either a storage account name (only when allowCrossTenantReplication of the SA is false) or a full resource id (both cases).
					// We should check for both cases.
					(prop.SourceAccount == srcAccount.StorageAccountName && prop.DestinationAccount == dstAccount.StorageAccountName) ||
						(strings.EqualFold(prop.SourceAccount, srcAccount.ID()) && strings.EqualFold(prop.DestinationAccount, dstAccount.ID()))) {
						srcId.ObjectReplicationPolicyId = *existing.Name
						dstId.ObjectReplicationPolicyId = *existing.Name
						return tf.ImportAsExistsError("azurerm_storage_object_replication", parse.NewObjectReplicationID(srcId, dstId).ID())
					}
				}
			}
		}
	}

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

	// create in dest storage account
	dstResp, err := client.ObjectReplicationPoliciesCreateOrUpdate(ctx, dstId, props)
	if err != nil {
		return fmt.Errorf("creating Storage Object Replication for destination storage account name %q: %+v", dstId.StorageAccountName, err)
	}

	if dstResp.Model == nil {
		return fmt.Errorf("nil model returned for Storage Object Replication for destination storage account name %q ID", dstId.StorageAccountName)
	}
	if dstResp.Model.Id == nil || *dstResp.Model.Id == "" {
		return fmt.Errorf("empty or nil ID returned for Storage Object Replication for destination storage account name %q ID", dstId.StorageAccountName)
	}
	if dstResp.Model.Name == nil || *dstResp.Model.Name == "" {
		return fmt.Errorf("empty or nil Name returned for Storage Object Replication for destination storage account name %q ID", dstAccount.StorageAccountName)
	}
	if dstResp.Model.Properties == nil {
		return fmt.Errorf("nil properties returned for Storage Object Replication for destination storage account name %q ID", dstAccount.StorageAccountName)
	}

	// Update the srcId and dstId using the returned computed object replication policy ID.
	srcId.ObjectReplicationPolicyId = *dstResp.Model.Name
	dstId.ObjectReplicationPolicyId = *dstResp.Model.Name

	// create in source storage account, update policy Id and ruleId which are computed from destination ORP
	props.Properties.Rules = dstResp.Model.Properties.Rules
	if _, err := client.ObjectReplicationPoliciesCreateOrUpdate(ctx, srcId, props); err != nil {
		return fmt.Errorf("creating Storage Object Replication %q for source storage account name %q: %+v", srcId.ObjectReplicationPolicyId, srcId.StorageAccountName, err)
	}

	d.SetId(parse.NewObjectReplicationID(srcId, dstId).ID())

	return resourceStorageObjectReplicationRead(d, meta)
}

func expandArmObjectReplicationRuleArray(input []any) *[]objectreplicationpolicyoperationgroup.ObjectReplicationPolicyRule {
	results := make([]objectreplicationpolicyoperationgroup.ObjectReplicationPolicyRule, 0)
	for _, item := range input {
		v := item.(map[string]any)
		result := objectreplicationpolicyoperationgroup.ObjectReplicationPolicyRule{
			SourceContainer:      v["source_container_name"].(string),
			DestinationContainer: v["destination_container_name"].(string),
			Filters: &objectreplicationpolicyoperationgroup.ObjectReplicationPolicyFilter{
				MinCreationTime: pointer.To(expandArmObjectReplicationMinCreationTime(v["copy_blobs_created_after"].(string))),
			},
		}

		if r, ok := v["name"].(string); ok && r != "" {
			result.RuleId = pointer.To(r)
		}

		if f, ok := v["filter_out_blobs_with_prefix"]; ok {
			result.Filters.PrefixMatch = pluginsdk.ExpandStringSlice(f.(*pluginsdk.Set).List())
		}

		results = append(results, result)
	}
	return &results
}

func expandArmObjectReplicationMinCreationTime(input string) string {
	switch input {
	case "Everything":
		return "1601-01-01T00:00:00Z"
	case "OnlyNewObjects":
		return ""
	default:
		return input
	}
}
