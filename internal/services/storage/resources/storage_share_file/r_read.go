// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_share_file

import (
	"fmt"
	"log"

	clientsAccount "github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	storageClient2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/files"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/shares"
)

func resourceStorageShareFileRead(d *pluginsdk.ResourceData, meta any) error {
	storageClient := storageClient2.From(meta)
	subscriptionId := clientsAccount.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(clientsAccount.StopContext(meta), d)
	defer cancel()

	id, err := files.ParseFileID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for File %q (Share %q): %s", id.AccountId.AccountName, id.FileName, id.ShareName, err)
	}
	if account == nil {
		log.Printf("[WARN] Unable to determine Storage Account for %s - assuming removed & removing from state", id)
		d.SetId("")
		return nil
	}

	client, err := storageClient.FileShareFilesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building File Share Client for Storage Account %s: %s", id.AccountId, err)
	}

	props, err := client.GetProperties(ctx, id.ShareName, id.DirectoryPath, id.FileName)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.Set("name", id.FileName)
	d.Set("path", id.DirectoryPath)
	d.Set("storage_share_url", shares.NewShareID(id.AccountId, id.ShareName).ID())

	if err = d.Set("metadata", helpers.FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("setting `metadata`: %s", err)
	}
	d.Set("content_type", props.ContentType)
	d.Set("content_encoding", props.ContentEncoding)

	// Set the ContentMD5 value to md5 hash in hex
	contentMD5 := ""
	if props.ContentMD5 != "" {
		contentMD5, err = helpers.ConvertBase64ToHexEncoding(props.ContentMD5)
		if err != nil {
			return fmt.Errorf("converting hex to base64 encoding for content_md5: %v", err)
		}
	}
	d.Set("content_md5", contentMD5)

	d.Set("content_disposition", props.ContentDisposition)

	if props.ContentLength == nil {
		return fmt.Errorf("file share file properties %q returned no information about the content-length", id.FileName)
	}

	d.Set("content_length", int(*props.ContentLength))

	return nil
}
