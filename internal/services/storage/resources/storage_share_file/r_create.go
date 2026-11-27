// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_share_file

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/files"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/shares"
)

func resourceStorageShareFileCreate(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	var (
		storageShareId *shares.ShareId
		err            error
	)
	storageShareId, err = shares.ParseShareID(d.Get("storage_share_url").(string), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	fileName := d.Get("name").(string)
	path := d.Get("path").(string)

	account, err := storageClient.FindAccount(ctx, subscriptionId, storageShareId.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for File %q (Share %q): %v", storageShareId.AccountId.AccountName, fileName, storageShareId.ShareName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", storageShareId.AccountId.AccountName)
	}

	accountId, err := accounts.ParseAccountID(storageShareId.ID(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	id := files.NewFileID(*accountId, storageShareId.ShareName, path, fileName)

	client, err := storageClient.FileShareFilesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building File Share Directories Client: %s", err)
	}

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.GetProperties(ctx, storageShareId.ShareName, path, fileName)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_storage_share_file", id.ID())
		}
	}

	input := files.CreateInput{
		MetaData:           helpers.ExpandMetaData(d.Get("metadata").(map[string]any)),
		ContentType:        pointer.To(d.Get("content_type").(string)),
		ContentEncoding:    pointer.To(d.Get("content_encoding").(string)),
		ContentDisposition: pointer.To(d.Get("content_disposition").(string)),
	}

	if v, ok := d.GetOk("content_md5"); ok {
		// Azure uses a Base64 encoded representation of the standard MD5 sum of the file
		contentMD5, err := helpers.ConvertHexToBase64Encoding(v.(string))
		if err != nil {
			return fmt.Errorf("failed to hex decode then base64 encode `content_md5` value: %s", err)
		}
		input.ContentMD5 = &contentMD5
	}

	var file *os.File
	if v, ok := d.GetOk("source"); ok {
		file, err = os.Open(v.(string))
		if err != nil {
			return fmt.Errorf("opening file: %s", err)
		}
		defer file.Close()
	} else if v, ok := d.GetOk("source_content"); ok {
		file, err = os.CreateTemp(os.TempDir(), "share-file-")
		if err != nil {
			return fmt.Errorf("creating temporary file: %s", err)
		}
		defer os.Remove(file.Name())
		defer file.Close()

		if _, err = file.WriteString(v.(string)); err != nil {
			return fmt.Errorf("writing `source_content` to temporary file: %s", err)
		}
	}

	if file != nil {
		info, err := file.Stat()
		if err != nil {
			return fmt.Errorf("'stat'-ing File %q (File Share %q / Account %q): %v", fileName, storageShareId.ShareName, storageShareId.AccountId.AccountName, err)
		}

		if info.Size() == 0 {
			return fmt.Errorf("file %q (File Share %q / Account %q) is empty", fileName, storageShareId.ShareName, storageShareId.AccountId.AccountName)
		}

		input.ContentLength = info.Size()
	}

	if _, err = client.Create(ctx, storageShareId.ShareName, path, fileName, input); err != nil {
		return fmt.Errorf("creating File %q (File Share %q / Account %q): %v", fileName, storageShareId.ShareName, storageShareId.AccountId.AccountName, err)
	}

	if file != nil {
		if err = client.PutFile(ctx, storageShareId.ShareName, path, fileName, file, 4); err != nil {
			return fmt.Errorf("uploading File: %q (File Share %q / Account %q): %v", fileName, storageShareId.ShareName, storageShareId.AccountId.AccountName, err)
		}
	}

	d.SetId(id.ID())

	return resourceStorageShareFileRead(d, meta)
}
