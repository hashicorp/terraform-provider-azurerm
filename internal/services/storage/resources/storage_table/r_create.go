// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_table

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/tables"
)

func resourceStorageTableCreate(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	tableName := d.Get("name").(string)
	aclsRaw := d.Get("acl").(*pluginsdk.Set).List()
	acls := expandStorageTableACLs(aclsRaw)

	accountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}
	accountName := accountId.StorageAccountName

	account, err := storageClient.FindAccount(ctx, subscriptionId, accountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Table %q: %s", accountName, tableName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", accountName)
	}

	tablesDataPlaneClient, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Tables Client: %s", err)
	}

	id := parse.NewStorageTableResourceManagerID(subscriptionId, accountId.ResourceGroupName, accountName, "default", tableName)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		exists, err := tablesDataPlaneClient.Exists(ctx, tableName)
		if err != nil {
			return fmt.Errorf("checking for existing %s: %v", id, err)
		}
		if exists != nil && *exists {
			return tf.ImportAsExistsError("azurerm_storage_table", id.ID())
		}
	}

	if err = tablesDataPlaneClient.Create(ctx, tableName); err != nil {
		return fmt.Errorf("creating %s: %v", id.ID(), err)
	}

	d.SetId(id.ID())

	aclClient, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Tables Client: %v", err)
	}

	if err = aclClient.UpdateACLs(ctx, tableName, acls); err != nil {
		return fmt.Errorf("setting ACLs for %s: %v", id.ID(), err)
	}

	return resourceStorageTableRead(d, meta)
}

func expandStorageTableACLs(input []any) []tables.SignedIdentifier {
	results := make([]tables.SignedIdentifier, 0)

	for _, v := range input {
		vals := v.(map[string]any)

		policies := vals["access_policy"].([]any)
		policy := policies[0].(map[string]any)

		identifier := tables.SignedIdentifier{
			Id: vals["id"].(string),
			AccessPolicy: tables.AccessPolicy{
				Start:      policy["start"].(string),
				Expiry:     policy["expiry"].(string),
				Permission: policy["permissions"].(string),
			},
		}
		results = append(results, identifier)
	}

	return results
}
