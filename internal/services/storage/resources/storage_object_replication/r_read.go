// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_object_replication

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/objectreplicationpolicyoperationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	storageClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageObjectReplicationRead(d *pluginsdk.ResourceData, meta any) error {
	client := storageClient.From(meta).ResourceManager.ObjectReplicationPolicyOperationGroup
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := parse.ObjectReplicationID(d.Id())
	if err != nil {
		return err
	}

	dstResp, err := client.ObjectReplicationPoliciesGet(ctx, id.Dst)
	if err != nil {
		if response.WasNotFound(dstResp.HttpResponse) {
			log.Printf("[INFO] storage object replication %q (dst) does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	srcResp, err := client.ObjectReplicationPoliciesGet(ctx, id.Src)
	if err != nil {
		if response.WasNotFound(srcResp.HttpResponse) {
			log.Printf("[INFO] storage object replication %q (src) does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	if model := dstResp.Model; model != nil {
		if props := dstResp.Model.Properties; props != nil {
			d.Set("source_storage_account_id", commonids.NewStorageAccountID(id.Src.SubscriptionId, id.Src.ResourceGroupName, id.Src.StorageAccountName).ID())
			d.Set("destination_storage_account_id", commonids.NewStorageAccountID(id.Dst.SubscriptionId, id.Dst.ResourceGroupName, id.Dst.StorageAccountName).ID())
			if err := d.Set("rules", flattenObjectReplicationRules(props.Rules)); err != nil {
				return fmt.Errorf("setting `rules`: %+v", err)
			}
			d.Set("source_object_replication_id", id.Src.ID())
			d.Set("destination_object_replication_id", id.Dst.ID())
		}
	}

	if model := srcResp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("metrics_enabled", props.Metrics != nil && pointer.From(props.Metrics.Enabled))
		}
	}

	return nil
}

func flattenObjectReplicationRules(input *[]objectreplicationpolicyoperationgroup.ObjectReplicationPolicyRule) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		destinationContainer := item.DestinationContainer
		sourceContainer := item.SourceContainer

		ruleId := pointer.From(item.RuleId)

		var minCreationTime string
		if item.Filters != nil && item.Filters.MinCreationTime != nil {
			minCreationTime = *item.Filters.MinCreationTime
		}

		var prefix []any
		if item.Filters != nil && item.Filters.PrefixMatch != nil {
			prefix = pluginsdk.FlattenSlice(item.Filters.PrefixMatch)
		}

		v := map[string]any{
			"destination_container_name":   destinationContainer,
			"source_container_name":        sourceContainer,
			"copy_blobs_created_after":     flattenArmObjectReplicationMinCreationTime(minCreationTime),
			"filter_out_blobs_with_prefix": pluginsdk.NewSet(pluginsdk.HashString, prefix),
			"name":                         ruleId,
		}
		results = append(results, v)
	}
	return results
}

func flattenArmObjectReplicationMinCreationTime(input string) string {
	switch input {
	case "1601-01-01T00:00:00Z":
		return "Everything"
	case "":
		return "OnlyNewObjects"
	default:
		return input
	}
}
