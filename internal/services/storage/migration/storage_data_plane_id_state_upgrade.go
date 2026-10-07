// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

// resolveStorageAccountIDForStateUpgrade determines the parent storage account ID during a state migration
// from a data plane ID to a resource manager ID. It first attempts to use the `resource_manager_id` attribute
// (only if it is a valid id) and falls back to looking the account up by name for states created before then.
func resolveStorageAccountIDForStateUpgrade(ctx context.Context, meta any, rawState map[string]any, parseResourceManagerID func(string) (*commonids.StorageAccountId, error)) (*commonids.StorageAccountId, error) {
	if resourceManagerID, ok := rawState["resource_manager_id"].(string); ok {
		// The `resource_manager_id` can be malformed (see: #32950), in which case fallbacks to find account.
		if storageAccountID, err := parseResourceManagerID(resourceManagerID); err == nil {
			return storageAccountID, nil
		}
	}

	// `resource_manager_id` was introduced in v3.39.0
	// The find account logic is only here for edge cases where users upgrade from < v3.39.0 directly to >= 5.0.0
	client := meta.(*clients.Client).Storage
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId

	storageAccountNameRaw, ok := rawState["storage_account_name"]
	if !ok {
		return nil, errors.New("expected a `storage_account_name` attribute to be present in state")
	}

	storageAccountName, ok := storageAccountNameRaw.(string)
	if !ok {
		return nil, fmt.Errorf("expected `storage_account_name` to be of type string, got %T", storageAccountNameRaw)
	}

	// This may seem like an excessive timeout, however, populating the accounts via the list API could take a significant amount of time
	// when the subscription contains a large number of accounts and the account cache is not already populated.
	findCtx, cancel := context.WithTimeout(ctx, 1*time.Hour)
	defer cancel()

	log.Printf("[DEBUG] searching for a storage account by name (`%s`) in subscription (`%s`)", storageAccountName, subscriptionID)
	account, err := client.FindAccount(findCtx, subscriptionID, storageAccountName)
	if err != nil || account == nil {
		return nil, fmt.Errorf("locating a storage account by name (`%s`) in subscription (`%s`): %w", storageAccountName, subscriptionID, err)
	}

	return &account.StorageAccountId, nil
}
