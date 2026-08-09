// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func LegacyStorageAccountResourceCustomizeDiff(_ context.Context, diff *pluginsdk.ResourceDiff, _ any) error {
	if strings.HasPrefix(diff.Id(), "/subscriptions/") && diff.HasChange("storage_account_id") {
		return diff.ForceNew("storage_account_id")
	}

	if diff.Id() != "" && !strings.HasPrefix(diff.Id(), "/subscriptions/") && diff.HasChange("storage_account_name") {
		oldAccountId, newAccountId := diff.GetChange("storage_account_id")
		oldName, newName := diff.GetChange("storage_account_name")

		if oldAccountId.(string) != "" && newName.(string) != "" {
			return diff.ForceNew("storage_account_name")
		}

		if oldName.(string) != "" && newName.(string) != "" {
			return diff.ForceNew("storage_account_name")
		}

		if oldName.(string) != "" && newName.(string) == "" && newAccountId.(string) != "" {
			parsedId, err := commonids.ParseStorageAccountID(newAccountId.(string))
			if err != nil {
				return err
			}
			if !strings.EqualFold(parsedId.StorageAccountName, oldName.(string)) {
				return diff.ForceNew("storage_account_id")
			}
		}
	}

	return nil
}
