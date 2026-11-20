// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_blob_inventory_policy

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/blobinventorypolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageBlobInventoryPolicyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Storage.ResourceManager.BlobInventoryPolicies
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Blob Inventory Policy for %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		d.Set("storage_account_id", id.ID())
		if props := model.Properties; props != nil {
			if !props.Policy.Enabled {
				log.Printf("[INFO] Blob Inventory Policy is not enabled for %s - removing from state", *id)
				d.SetId("")
				return nil
			}

			d.Set("rules", flattenBlobInventoryPolicyRules(props.Policy.Rules))
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id, pluginsdk.ResourceTypeForIdentityVirtual)
}

func flattenBlobInventoryPolicyRules(input []blobinventorypolicies.BlobInventoryPolicyRule) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range input {
		if !item.Enabled {
			continue
		}

		results = append(results, map[string]any{
			"name":                   item.Name,
			"storage_container_name": item.Destination,
			"format":                 string(item.Definition.Format),
			"schedule":               string(item.Definition.Schedule),
			"scope":                  string(item.Definition.ObjectType),
			"schema_fields":          item.Definition.SchemaFields,
			"filter":                 flattenBlobInventoryPolicyFilter(item.Definition.Filters),
		})
	}
	return results
}

func flattenBlobInventoryPolicyFilter(input *blobinventorypolicies.BlobInventoryPolicyFilter) []any {
	if input == nil {
		return make([]any, 0)
	}

	var createdWithinDays int
	if input.CreationTime != nil {
		createdWithinDays = int(pointer.From(input.CreationTime.LastNDays))
	}

	return []any{
		map[string]any{
			"blob_types":            pluginsdk.FlattenSlice(input.BlobTypes),
			"created_within_days":   createdWithinDays,
			"include_blob_versions": pointer.From(input.IncludeBlobVersions),
			"include_deleted":       pointer.From(input.IncludeDeleted),
			"include_snapshots":     pointer.From(input.IncludeSnapshots),
			"prefix_match":          pluginsdk.FlattenSlice(input.PrefixMatch),
			"exclude_prefixes":      pluginsdk.FlattenSlice(input.ExcludePrefix),
		},
	}
}
