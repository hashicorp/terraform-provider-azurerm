// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_table

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/tables"
)

func resourceStorageTableRead(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	rmId, err := parse.StorageTableResourceManagerID(d.Id())
	if err != nil {
		return err
	}

	tableName := rmId.TableName
	accountName := rmId.StorageAccountName

	account, err := storageClient.GetAccount(ctx, commonids.NewStorageAccountID(rmId.SubscriptionId, rmId.ResourceGroup, rmId.StorageAccountName))
	if err != nil {
		return fmt.Errorf("retrieving Storage Account %q for Table %q: %v", accountName, tableName, err)
	}
	if account == nil {
		log.Printf("Unable to determine Resource Group for Storage Table %q (Account %s) - assuming removed & removing from state", tableName, accountName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Tables Client: %v", err)
	}

	exists, err := client.Exists(ctx, tableName)
	if err != nil {
		return fmt.Errorf("retrieving table %q: %v", tableName, err)
	}
	if exists == nil || !*exists {
		log.Printf("[DEBUG] table %q not found, removing from state", tableName)
		d.SetId("")
		return nil
	}

	aclClient, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Tables Client: %v", err)
	}

	acls, err := aclClient.GetACLs(ctx, tableName)
	if err != nil {
		return fmt.Errorf("retrieving ACLs for table %q: %v", tableName, err)
	}

	d.Set("name", tableName)
	d.Set("storage_account_id", commonids.NewStorageAccountID(rmId.SubscriptionId, rmId.ResourceGroup, rmId.StorageAccountName).ID())

	if err = d.Set("acl", flattenStorageTableACLs(acls)); err != nil {
		return fmt.Errorf("setting `acl`: %v", err)
	}

	return nil
}

func flattenStorageTableACLs(input *[]tables.SignedIdentifier) []any {
	result := make([]any, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		output := map[string]any{
			"id": v.Id,
			"access_policy": []any{
				map[string]any{
					"start":       v.AccessPolicy.Start,
					"expiry":      v.AccessPolicy.Expiry,
					"permissions": v.AccessPolicy.Permission,
				},
			},
		}

		result = append(result, output)
	}

	return result
}
