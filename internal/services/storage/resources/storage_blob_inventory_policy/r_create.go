// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_blob_inventory_policy

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/blobinventorypolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	storageClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageBlobInventoryPolicyCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := account.From(meta).SubscriptionId
	client := storageClient.From(meta).ResourceManager.BlobInventoryPolicies
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	accountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	// This Resource is 1:1 with a Storage Account, therefore we use it's Resource ID for this Resource
	// however we want to ensure it's in the same subscription, so we'll build this up here
	id := commonids.NewStorageAccountID(subscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName)
	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %q: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_storage_blob_inventory_policy", id.ID())
			}
		}
	}

	rules, err := expandBlobInventoryPolicyRules(d.Get("rules").(*pluginsdk.Set).List())
	if err != nil {
		return err
	}

	payload := blobinventorypolicies.BlobInventoryPolicy{
		Properties: &blobinventorypolicies.BlobInventoryPolicyProperties{
			Policy: blobinventorypolicies.BlobInventoryPolicySchema{
				Enabled: true,
				Type:    blobinventorypolicies.InventoryRuleTypeInventory,
				Rules:   rules,
			},
		},
	}
	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id, pluginsdk.ResourceTypeForIdentityVirtual); err != nil {
		return err
	}

	return resourceStorageBlobInventoryPolicyRead(d, meta)
}

func expandBlobInventoryPolicyRules(input []any) ([]blobinventorypolicies.BlobInventoryPolicyRule, error) {
	results := make([]blobinventorypolicies.BlobInventoryPolicyRule, 0)
	for _, item := range input {
		v := item.(map[string]any)

		filters, err := expandBlobInventoryPolicyFilter(v["filter"].([]any), v["scope"].(string))
		if err != nil {
			return nil, fmt.Errorf("%s rule is invalid: %+v", v["name"].(string), err)
		}

		results = append(results, blobinventorypolicies.BlobInventoryPolicyRule{
			Enabled:     true,
			Name:        v["name"].(string),
			Destination: v["storage_container_name"].(string),
			Definition: blobinventorypolicies.BlobInventoryPolicyDefinition{
				Format:       blobinventorypolicies.Format(v["format"].(string)),
				Schedule:     blobinventorypolicies.Schedule(v["schedule"].(string)),
				ObjectType:   blobinventorypolicies.ObjectType(v["scope"].(string)),
				SchemaFields: *pluginsdk.ExpandStringSlice(v["schema_fields"].([]any)),
				Filters:      filters,
			},
		})
	}
	return results, nil
}

func expandBlobInventoryPolicyFilter(input []any, objectType string) (*blobinventorypolicies.BlobInventoryPolicyFilter, error) {
	if len(input) == 0 {
		return nil, nil
	}
	v := input[0].(map[string]any)
	policyFilter := &blobinventorypolicies.BlobInventoryPolicyFilter{
		PrefixMatch:         pluginsdk.ExpandStringSlice(v["prefix_match"].(*pluginsdk.Set).List()),
		ExcludePrefix:       pluginsdk.ExpandStringSlice(v["exclude_prefixes"].(*pluginsdk.Set).List()),
		BlobTypes:           pluginsdk.ExpandStringSlice(v["blob_types"].(*pluginsdk.Set).List()),
		IncludeBlobVersions: pointer.To(v["include_blob_versions"].(bool)),
		IncludeDeleted:      pointer.To(v["include_deleted"].(bool)),
		IncludeSnapshots:    pointer.To(v["include_snapshots"].(bool)),
	}

	// If the objectType is Container, the following values must be nil when passed to the API
	if objectType == string(blobinventorypolicies.ObjectTypeContainer) {
		if len(*policyFilter.BlobTypes) > 0 || *policyFilter.IncludeBlobVersions || *policyFilter.IncludeSnapshots {
			return nil, errors.New("`blobTypes`, `includeBlobVersions`, `includeSnapshots` cannot be used with objectType `Container`")
		}
		policyFilter.BlobTypes = nil
		policyFilter.IncludeBlobVersions = nil
		policyFilter.IncludeSnapshots = nil
	}

	if v["created_within_days"].(int) > 0 {
		policyFilter.CreationTime = &blobinventorypolicies.BlobInventoryCreationTime{
			LastNDays: pointer.To(int64(v["created_within_days"].(int))),
		}
	}

	return policyFilter, nil
}
