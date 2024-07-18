// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
)

var (
	storageAccountsCache = map[string]AccountDetails{}

	cacheAccountsLock    = sync.RWMutex{}
	cacheCredentialsLock = sync.RWMutex{}
)

type EndpointType string

const (
	EndpointTypeBlob  = "blob"
	EndpointTypeDfs   = "dfs"
	EndpointTypeFile  = "file"
	EndpointTypeQueue = "queue"
	EndpointTypeTable = "table"
)

type AccountDetails struct {
	Kind             storageaccounts.Kind
	IsHnsEnabled     bool
	StorageAccountId commonids.StorageAccountId

	accountKey *string

	// primaryBlobEndpoint is the Primary Blob Endpoint for the Data Plane API for this Storage Account
	// e.g. `https://{account}.blob.core.windows.net`
	primaryBlobEndpoint *string

	// primaryDfsEndpoint is the Primary Dfs Endpoint for the Data Plane API for this Storage Account
	// e.g. `https://sale.dfs.core.windows.net`
	primaryDfsEndpoint *string

	// primaryFileEndpoint is the Primary File Endpoint for the Data Plane API for this Storage Account
	// e.g. `https://{account}.file.core.windows.net`
	primaryFileEndpoint *string

	// primaryQueueEndpoint is the Primary Queue Endpoint for the Data Plane API for this Storage Account
	// e.g. `https://{account}.queue.core.windows.net`
	primaryQueueEndpoint *string

	// primaryTableEndpoint is the Primary Table Endpoint for the Data Plane API for this Storage Account
	// e.g. `https://{account}.table.core.windows.net`
	primaryTableEndpoint *string
}

func (ad *AccountDetails) AccountKey(ctx context.Context, client Client) (*string, error) {
	cacheCredentialsLock.Lock()
	defer cacheCredentialsLock.Unlock()

	if ad.accountKey != nil {
		return ad.accountKey, nil
	}

	log.Printf("[DEBUG] Cache Miss - looking up the account key for %s..", ad.StorageAccountId)
	opts := storageaccounts.DefaultListKeysOperationOptions()
	opts.Expand = pointer.To(storageaccounts.ListKeyExpandKerb)
	listKeysResp, err := client.ResourceManager.StorageAccounts.ListKeys(ctx, ad.StorageAccountId, opts)
	if err != nil {
		return nil, fmt.Errorf("listing Keys for %s: %+v", ad.StorageAccountId, err)
	}

	if model := listKeysResp.Model; model != nil && model.Keys != nil {
		for _, key := range *model.Keys {
			if key.Permissions == nil || key.Value == nil {
				continue
			}

			if *key.Permissions == storageaccounts.KeyPermissionFull {
				ad.accountKey = key.Value
				break
			}
		}
	}

	if ad.accountKey == nil {
		return nil, fmt.Errorf("unable to determine the Write Key for %s", ad.StorageAccountId)
	}

	// force-cache this
	storageAccountsCache[ad.StorageAccountId.StorageAccountName] = *ad

	return ad.accountKey, nil
}

func (ad *AccountDetails) DataPlaneEndpoint(endpointType EndpointType) (*string, error) {
	var baseUri *string
	switch endpointType {
	case EndpointTypeBlob:
		baseUri = ad.primaryBlobEndpoint

	case EndpointTypeDfs:
		baseUri = ad.primaryDfsEndpoint

	case EndpointTypeFile:
		baseUri = ad.primaryFileEndpoint

	case EndpointTypeQueue:
		baseUri = ad.primaryQueueEndpoint

	case EndpointTypeTable:
		baseUri = ad.primaryTableEndpoint

	default:
		return nil, fmt.Errorf("internal-error: unrecognised endpoint type %q when building storage client", endpointType)
	}

	if baseUri == nil {
		return nil, fmt.Errorf("determining %s endpoint for %s: missing primary endpoint", endpointType, ad.StorageAccountId)
	}
	return baseUri, nil
}

func (c Client) AddToCache(accountId commonids.StorageAccountId, account storageaccounts.StorageAccount) error {
	cacheAccountsLock.Lock()
	defer cacheAccountsLock.Unlock()

	accountDetails, err := populateAccountDetails(accountId, account)
	if err != nil {
		return err
	}
	storageAccountsCache[accountId.StorageAccountName] = *accountDetails

	return nil
}

func (c Client) RemoveAccountFromCache(accountId commonids.StorageAccountId) {
	cacheAccountsLock.Lock()
	delete(storageAccountsCache, accountId.StorageAccountName)
	cacheAccountsLock.Unlock()
}

func (c Client) FindAccount(ctx context.Context, subscriptionIdRaw, accountName string) (*AccountDetails, error) {
	cacheAccountsLock.Lock()
	defer cacheAccountsLock.Unlock()

	if existing, ok := storageAccountsCache[accountName]; ok {
		return &existing, nil
	}

	subscriptionId := commonids.NewSubscriptionID(subscriptionIdRaw)
	listResult, err := c.ResourceManager.StorageAccounts.ListComplete(ctx, subscriptionId)
	if err != nil {
		return nil, fmt.Errorf("listing Storage Accounts within %s: %+v", subscriptionId, err)
	}
	for _, item := range listResult.Items {
		if item.Id == nil || item.Name == nil {
			continue
		}

		storageAccountId, err := commonids.ParseStorageAccountIDInsensitively(*item.Id)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", *item.Id, err)
		}

		account, err := populateAccountDetails(*storageAccountId, item)
		if err != nil {
			return nil, fmt.Errorf("populating details for %s: %+v", *storageAccountId, err)
		}

		storageAccountsCache[storageAccountId.StorageAccountName] = *account
	}

	if existing, ok := storageAccountsCache[accountName]; ok {
		return &existing, nil
	}

	return nil, nil
}

func populateAccountDetails(accountId commonids.StorageAccountId, account storageaccounts.StorageAccount) (*AccountDetails, error) {
	out := AccountDetails{
		Kind:             pointer.From(account.Kind),
		StorageAccountId: accountId,
	}

	if account.Properties == nil {
		return nil, fmt.Errorf("populating details for %s: `model.Properties` was nil", accountId)
	}
	if account.Properties.PrimaryEndpoints == nil {
		return nil, fmt.Errorf("populating details for %s: `model.Properties.PrimaryEndpoints` was nil", accountId)
	}

	props := *account.Properties
	out.IsHnsEnabled = pointer.From(props.IsHnsEnabled)

	endpoints := *props.PrimaryEndpoints
	if endpoints.Blob != nil {
		endpoint := strings.TrimSuffix(*endpoints.Blob, "/")
		out.primaryBlobEndpoint = pointer.To(endpoint)
	}
	if endpoints.Dfs != nil {
		endpoint := strings.TrimSuffix(*endpoints.Dfs, "/")
		out.primaryDfsEndpoint = pointer.To(endpoint)
	}
	if endpoints.File != nil {
		endpoint := strings.TrimSuffix(*endpoints.File, "/")
		out.primaryFileEndpoint = pointer.To(endpoint)
	}
	if endpoints.Queue != nil {
		endpoint := strings.TrimSuffix(*endpoints.Queue, "/")
		out.primaryQueueEndpoint = pointer.To(endpoint)
	}
	if endpoints.Table != nil {
		endpoint := strings.TrimSuffix(*endpoints.Table, "/")
		out.primaryTableEndpoint = pointer.To(endpoint)
	}

	return &out, nil
}
