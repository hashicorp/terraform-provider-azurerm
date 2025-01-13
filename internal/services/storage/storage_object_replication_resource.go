// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/objectreplicationpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// TODO: @tombuildsstuff: this wants a state migration to move the ID to `{id1}|{id2}` to match other resources

func resourceStorageObjectReplication() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageObjectReplicationCreate,
		Read:   resourceStorageObjectReplicationRead,
		Update: resourceStorageObjectReplicationUpdate,
		Delete: resourceStorageObjectReplicationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ObjectReplicationID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"source_storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"destination_storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"rules": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"source_container_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.StorageContainerName,
						},

						"destination_container_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.StorageContainerName,
						},

						// Possible values are "" means "OnlyNewObjects", "1601-01-01T00:00:00Z" means "Everything" and timeStamp "2020-10-21T16:00:00Z"
						"copy_blobs_created_after": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "OnlyNewObjects",
							ValidateFunc: validate.ObjectReplicationCopyBlobsCreatedAfter,
						},

						"filter_out_blobs_with_prefix": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"source_object_replication_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"destination_object_replication_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceStorageObjectReplicationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ResourceManager.ObjectReplicationPolicies
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

	srcId := objectreplicationpolicies.NewObjectReplicationPolicyID(srcAccount.SubscriptionId, srcAccount.ResourceGroupName, srcAccount.StorageAccountName, "default")
	dstId := objectreplicationpolicies.NewObjectReplicationPolicyID(dstAccount.SubscriptionId, dstAccount.ResourceGroupName, dstAccount.StorageAccountName, "default")

	resp, err := client.List(ctx, *dstAccount)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("checking for present of existing Storage Object Replication for destination %q): %+v", dstAccount, err)
		}
	}
	if resp.Model != nil && resp.Model.Value != nil {
		for _, existing := range *resp.Model.Value {
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

	props := objectreplicationpolicies.ObjectReplicationPolicy{
		Properties: &objectreplicationpolicies.ObjectReplicationPolicyProperties{
			SourceAccount:      srcAccount.ID(),
			DestinationAccount: dstAccount.ID(),
			Rules:              expandArmObjectReplicationRuleArray(d.Get("rules").(*pluginsdk.Set).List()),
		},
	}

	// create in dest storage account
	dstResp, err := client.CreateOrUpdate(ctx, dstId, props)
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
	if _, err := client.CreateOrUpdate(ctx, srcId, props); err != nil {
		return fmt.Errorf("creating Storage Object Replication %q for source storage account name %q: %+v", srcId.ObjectReplicationPolicyId, srcId.StorageAccountName, err)
	}

	d.SetId(parse.NewObjectReplicationID(srcId, dstId).ID())

	return resourceStorageObjectReplicationRead(d, meta)
}

func resourceStorageObjectReplicationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ResourceManager.ObjectReplicationPolicies
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ObjectReplicationID(d.Id())
	if err != nil {
		return err
	}

	srcAccount := commonids.NewStorageAccountID(id.Src.SubscriptionId, id.Src.ResourceGroupName, id.Src.StorageAccountName)
	dstAccount := commonids.NewStorageAccountID(id.Dst.SubscriptionId, id.Dst.ResourceGroupName, id.Dst.StorageAccountName)

	props := objectreplicationpolicies.ObjectReplicationPolicy{
		Properties: &objectreplicationpolicies.ObjectReplicationPolicyProperties{
			SourceAccount:      srcAccount.ID(),
			DestinationAccount: dstAccount.ID(),
			Rules:              expandArmObjectReplicationRuleArray(d.Get("rules").(*pluginsdk.Set).List()),
		},
	}

	// update in dest storage account
	resp, err := client.CreateOrUpdate(ctx, id.Dst, props)
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
	if _, err := client.CreateOrUpdate(ctx, id.Src, props); err != nil {
		return fmt.Errorf("updating %q for source storage account name %q: %+v", id, id.Src.StorageAccountName, err)
	}

	return resourceStorageObjectReplicationRead(d, meta)
}

func resourceStorageObjectReplicationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ResourceManager.ObjectReplicationPolicies
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ObjectReplicationID(d.Id())
	if err != nil {
		return err
	}

	dstResp, err := client.Get(ctx, id.Dst)
	if err != nil {
		if response.WasNotFound(dstResp.HttpResponse) {
			log.Printf("[INFO] storage object replication %q (dst) does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	srcResp, err := client.Get(ctx, id.Src)
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
	return nil
}

func resourceStorageObjectReplicationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ResourceManager.ObjectReplicationPolicies
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ObjectReplicationID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.Dst); err != nil {
		return fmt.Errorf("deleting %q: %+v", id.Dst, err)
	}

	if _, err := client.Delete(ctx, id.Src); err != nil {
		return fmt.Errorf("deleting %q : %+v", id.Dst, err)
	}
	return nil
}

func expandArmObjectReplicationRuleArray(input []interface{}) *[]objectreplicationpolicies.ObjectReplicationPolicyRule {
	results := make([]objectreplicationpolicies.ObjectReplicationPolicyRule, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		result := objectreplicationpolicies.ObjectReplicationPolicyRule{
			SourceContainer:      v["source_container_name"].(string),
			DestinationContainer: v["destination_container_name"].(string),
			Filters: &objectreplicationpolicies.ObjectReplicationPolicyFilter{
				MinCreationTime: utils.String(expandArmObjectReplicationMinCreationTime(v["copy_blobs_created_after"].(string))),
			},
		}

		if r, ok := v["name"].(string); ok && r != "" {
			result.RuleId = utils.String(r)
		}

		if f, ok := v["filter_out_blobs_with_prefix"]; ok {
			result.Filters.PrefixMatch = utils.ExpandStringSlice(f.(*pluginsdk.Set).List())
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

func flattenObjectReplicationRules(input *[]objectreplicationpolicies.ObjectReplicationPolicyRule) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		destinationContainer := item.DestinationContainer
		sourceContainer := item.SourceContainer

		var ruleId string
		if item.RuleId != nil {
			ruleId = *item.RuleId
		}

		var minCreationTime string
		if item.Filters != nil && item.Filters.MinCreationTime != nil {
			minCreationTime = *item.Filters.MinCreationTime
		}

		var prefix []interface{}
		if item.Filters != nil && item.Filters.PrefixMatch != nil {
			prefix = utils.FlattenStringSlice(item.Filters.PrefixMatch)
		}

		v := map[string]interface{}{
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
