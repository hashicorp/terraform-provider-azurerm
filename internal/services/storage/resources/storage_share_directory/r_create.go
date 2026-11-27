// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_share_directory

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/directories"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/shares"
)

func resourceStorageShareDirectoryCreate(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	directoryName := d.Get("name").(string)
	metaDataRaw := d.Get("metadata").(map[string]any)
	metaData := helpers.ExpandMetaData(metaDataRaw)

	var (
		storageShareId *shares.ShareId
		err            error
	)
	storageShareId, err = shares.ParseShareID(d.Get("storage_share_url").(string), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	if storageShareId == nil {
		return errors.New("determining storage share ID")
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, storageShareId.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Directory %q (Share %q): %v", storageShareId.AccountId.AccountName, directoryName, storageShareId.ShareName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", storageShareId.AccountId.AccountName)
	}

	accountId, err := accounts.ParseAccountID(storageShareId.ID(), storageClient.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
	}

	id := directories.NewDirectoryID(*accountId, storageShareId.ShareName, directoryName)

	client, err := storageClient.FileShareDirectoriesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building File Share Directories Client: %v", err)
	}

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, storageShareId.ShareName, directoryName)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_storage_share_directory", id.ID())
		}
	}

	input := directories.CreateDirectoryInput{
		MetaData: metaData,
	}
	if _, err = client.Create(ctx, storageShareId.ShareName, directoryName, input); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}
	d.SetId(id.ID())

	// Storage Share Directories are eventually consistent
	poller := custompollers.NewEventualConsistencyPoller(5, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.Get(pollerCtx, id.ShareName, id.DirectoryPath)
		return resp.HttpResponse, err
	}, &custompollers.EventualConsistencyPollerOptions{
		Interval:              10 * time.Second,
		TargetStatusCode:      pointer.To(http.StatusOK),
		RetryErrorStatusCodes: []int{http.StatusNotFound},
	})
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become available: %v", id, err)
	}

	return resourceStorageShareDirectoryRead(d, meta)
}
