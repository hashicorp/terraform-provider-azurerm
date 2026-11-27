// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_share_file

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	clientsAccount "github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	storageClient2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/files"
)

func resourceStorageShareFileUpdate(d *pluginsdk.ResourceData, meta any) error {
	storageClient := storageClient2.From(meta)
	subscriptionId := clientsAccount.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForUpdate(clientsAccount.StopContext(meta), d)
	defer cancel()

	id, err := files.ParseFileID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for %s: %v", id.AccountId.AccountName, id, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", id.AccountId.AccountName)
	}

	client, err := storageClient.FileShareFilesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building File Share Files Client: %v", err)
	}

	existing, err := client.GetProperties(ctx, id.ShareName, id.DirectoryPath, id.FileName)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %v", id, err)
		}
	}

	if d.HasChanges("content_type", "content_encoding", "content_disposition") {
		input := files.SetPropertiesInput{
			ContentType:        pointer.To(d.Get("content_type").(string)),
			ContentEncoding:    pointer.To(d.Get("content_encoding").(string)),
			ContentDisposition: pointer.To(d.Get("content_disposition").(string)),
			ContentLength:      int64(d.Get("content_length").(int)),
			MetaData:           helpers.ExpandMetaData(d.Get("metadata").(map[string]any)),
		}

		if v, ok := d.GetOk("content_md5"); ok {
			// Azure uses a Base64 encoded representation of the standard MD5 sum of the file
			contentMD5, err := helpers.ConvertHexToBase64Encoding(v.(string))
			if err != nil {
				return fmt.Errorf("failed to hex decode then base64 encode `content_md5` value: %s", err)
			}
			input.ContentMD5 = &contentMD5
		}

		if _, err = client.SetProperties(ctx, id.ShareName, id.DirectoryPath, id.FileName, input); err != nil {
			return fmt.Errorf("creating %s: %v", id, err)
		}
	}

	return resourceStorageShareFileRead(d, meta)
}
