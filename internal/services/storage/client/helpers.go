// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

var (
	storageAccountsCache = map[string]accountDetails{}

	accountsLock    = sync.RWMutex{}
	credentialsLock = sync.RWMutex{}
)

type EndpointType string

const (
	EndpointTypeBlob  = "blob"
	EndpointTypeDfs   = "dfs"
	EndpointTypeFile  = "file"
	EndpointTypeQueue = "queue"
	EndpointTypeTable = "table"
)

type accountDetails struct {
	Kind             storage.Kind
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

func (ad *accountDetails) AccountKey(ctx context.Context, client Client) (*string, error) {
	credentialsLock.Lock()
	defer credentialsLock.Unlock()

	if ad.accountKey != nil {
		return ad.accountKey, nil
	}

	log.Printf("[DEBUG] Cache Miss - looking up the account key for %s..", ad.StorageAccountId)
	props, err := client.AccountsClient.ListKeys(ctx, ad.StorageAccountId.ResourceGroupName, ad.StorageAccountId.StorageAccountName, storage.ListKeyExpandKerb)
	if err != nil {
		return nil, fmt.Errorf("listing Keys for %s: %+v", ad.StorageAccountId, err)
	}

	if props.Keys == nil || len(*props.Keys) == 0 || (*props.Keys)[0].Value == nil {
		return nil, fmt.Errorf("keys were nil for %s: %+v", ad.StorageAccountId, err)
	}

	keys := *props.Keys
	ad.accountKey = keys[0].Value

	// force-cache this
	storageAccountsCache[ad.StorageAccountId.StorageAccountName] = *ad

	return ad.accountKey, nil
}

func (ad *accountDetails) DataPlaneEndpoint(endpointType EndpointType) (*string, error) {
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

func (c Client) AddToCache(accountId commonids.StorageAccountId, props storage.Account) error {
	accountsLock.Lock()
	defer accountsLock.Unlock()

	account, err := populateAccountDetails(accountId, props)
	if err != nil {
		return err
	}

	storageAccountsCache[accountId.StorageAccountName] = *account

	return nil
}

func (c Client) RemoveAccountFromCache(accountId commonids.StorageAccountId) {
	accountsLock.Lock()
	delete(storageAccountsCache, accountId.StorageAccountName)
	accountsLock.Unlock()
}

func (c Client) FindAccount(ctx context.Context, accountName string) (*accountDetails, error) {
	accountsLock.Lock()
	defer accountsLock.Unlock()

	if existing, ok := storageAccountsCache[accountName]; ok {
		return &existing, nil
	}

	accountsPage, err := c.AccountsClient.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("retrieving storage accounts: %+v", err)
	}

	var accounts []storage.Account
	for accountsPage.NotDone() {
		accounts = append(accounts, accountsPage.Values()...)
		err = accountsPage.NextWithContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("retrieving next page of storage accounts: %+v", err)
		}
	}

	for _, v := range accounts {
		if v.ID == nil || v.Name == nil {
			continue
		}

		storageAccountId, err := commonids.ParseStorageAccountIDInsensitively(*v.ID)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", *v.ID, err)
		}

		account, err := populateAccountDetails(*storageAccountId, v)
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

func populateAccountDetails(accountId commonids.StorageAccountId, account storage.Account) (*accountDetails, error) {
	out := accountDetails{
		Kind:             account.Kind,
		StorageAccountId: accountId,
	}

	if props := account.AccountProperties; props != nil {
		out.IsHnsEnabled = pointer.From(props.IsHnsEnabled)

		if endpoints := props.PrimaryEndpoints; endpoints != nil {
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
		}
	}

	return &out, nil
}
