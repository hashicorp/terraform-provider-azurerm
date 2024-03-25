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

func (c Client) AddToCache(accountName string, props storage.Account) error {
	accountsLock.Lock()
	defer accountsLock.Unlock()

	account, err := populateAccountDetails(accountName, props)
	if err != nil {
		return err
	}

	storageAccountsCache[accountName] = *account

	return nil
}

func (c Client) RemoveAccountFromCache(accountName string) {
	accountsLock.Lock()
	delete(storageAccountsCache, accountName)
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
		if v.Name == nil {
			continue
		}

		account, err := populateAccountDetails(*v.Name, v)
		if err != nil {
			return nil, err
		}

		storageAccountsCache[*v.Name] = *account
	}

	if existing, ok := storageAccountsCache[accountName]; ok {
		return &existing, nil
	}

	return nil, nil
}

func populateAccountDetails(accountName string, props storage.Account) (*accountDetails, error) {
	if props.ID == nil {
		return nil, fmt.Errorf("`id` was nil for Account %q", accountName)
	}

	accountId := *props.ID
	id, err := commonids.ParseStorageAccountID(accountId)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as a Resource ID: %+v", accountId, err)
	}

	account := accountDetails{
		StorageAccountId: *id,
		Kind:             props.Kind,
	}

	if props.AccountProperties != nil {
		account.IsHnsEnabled = pointer.From(props.AccountProperties.IsHnsEnabled)

		if props.AccountProperties.PrimaryEndpoints != nil {
			if props.AccountProperties.PrimaryEndpoints.Blob != nil {
				endpoint := strings.TrimSuffix(*props.AccountProperties.PrimaryEndpoints.Blob, "/")
				account.primaryBlobEndpoint = pointer.To(endpoint)
			}
			if props.AccountProperties.PrimaryEndpoints.Dfs != nil {
				endpoint := strings.TrimSuffix(*props.AccountProperties.PrimaryEndpoints.Dfs, "/")
				account.primaryDfsEndpoint = pointer.To(endpoint)
			}
			if props.AccountProperties.PrimaryEndpoints.File != nil {
				endpoint := strings.TrimSuffix(*props.AccountProperties.PrimaryEndpoints.File, "/")
				account.primaryFileEndpoint = pointer.To(endpoint)
			}
			if props.AccountProperties.PrimaryEndpoints.Queue != nil {
				endpoint := strings.TrimSuffix(*props.AccountProperties.PrimaryEndpoints.Queue, "/")
				account.primaryQueueEndpoint = pointer.To(endpoint)
			}
			if props.AccountProperties.PrimaryEndpoints.Table != nil {
				endpoint := strings.TrimSuffix(*props.AccountProperties.PrimaryEndpoints.Table, "/")
				account.primaryTableEndpoint = pointer.To(endpoint)
			}
		}
	}

	return &account, nil
}
