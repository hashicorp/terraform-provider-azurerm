// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/custompollers"
)

type storageAccountServiceSupportLevel struct {
	supportBlob          bool
	supportQueue         bool
	supportShare         bool
	supportStaticWebsite bool
}

func availableFunctionalityForAccount(kind storageaccounts.Kind, tier storageaccounts.SkuTier, replicationType string) storageAccountServiceSupportLevel {
	// FileStorage doesn't support blob
	supportBlob := kind != storageaccounts.KindFileStorage

	// Queue is only supported for Storage and StorageV2, in Standard sku tier.
	supportQueue := tier == storageaccounts.SkuTierStandard && (kind == storageaccounts.KindStorageVTwo ||
		(kind == storageaccounts.KindStorage &&
			// Per local test, only LRS/GRS/RAGRS Storage V1 accounts support queue endpoint.
			// GZRS and RAGZRS is invalid, while ZRS is valid but has no queue endpoint.
			slices.Contains([]string{"LRS", "GRS", "RAGRS"}, replicationType)))

	// File share is only supported for StorageV2 and FileStorage.
	// See: https://docs.microsoft.com/en-us/azure/storage/files/storage-files-planning#management-concepts
	// Per test, the StorageV2 with Premium sku tier also doesn't support file share.
	supportShare := kind == storageaccounts.KindFileStorage || (tier != storageaccounts.SkuTierPremium && (kind == storageaccounts.KindStorageVTwo ||
		(kind == storageaccounts.KindStorage &&
			// Per local test, only LRS/GRS/RAGRS Storage V1 accounts support file endpoint.
			// GZRS and RAGZRS is invalid, while ZRS is valid but has no file endpoint.
			slices.Contains([]string{"LRS", "GRS", "RAGRS"}, replicationType))))

	// Static Website is only supported for StorageV2 (not for Storage(v1)) and BlockBlobStorage
	supportStaticWebSite := kind == storageaccounts.KindStorageVTwo || kind == storageaccounts.KindBlockBlobStorage

	return storageAccountServiceSupportLevel{
		supportBlob:          supportBlob,
		supportQueue:         supportQueue,
		supportShare:         supportShare,
		supportStaticWebsite: supportStaticWebSite,
	}
}

func waitForDataPlaneToBecomeAvailableForAccount(ctx context.Context, client *client.Client, account *client.AccountDetails, supportLevel storageAccountServiceSupportLevel) error {
	initialDelayDuration := 10 * time.Second

	if supportLevel.supportBlob {
		log.Printf("[DEBUG] waiting for the Blob Service to become available")
		pollerType, err := custompollers.NewDataPlaneBlobContainersAvailabilityPoller(ctx, client, account)
		if err != nil {
			return fmt.Errorf("building Blob Service Poller: %+v", err)
		}
		poller := pollers.NewPoller(pollerType, initialDelayDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("waiting for the Blob Service to become available: %+v", err)
		}
	}

	if supportLevel.supportQueue {
		log.Printf("[DEBUG] waiting for the Queues Service to become available")
		pollerType, err := custompollers.NewDataPlaneQueuesAvailabilityPoller(ctx, client, account)
		if err != nil {
			return fmt.Errorf("building Queues Poller: %+v", err)
		}
		poller := pollers.NewPoller(pollerType, initialDelayDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("waiting for the Queues Service to become available: %+v", err)
		}
	}

	if supportLevel.supportShare {
		log.Printf("[DEBUG] waiting for the File Service to become available")
		pollerType, err := custompollers.NewDataPlaneFileShareAvailabilityPoller(client, account)
		if err != nil {
			return fmt.Errorf("building File Share Poller: %+v", err)
		}
		poller := pollers.NewPoller(pollerType, initialDelayDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("waiting for the File Service to become available: %+v", err)
		}
	}

	if supportLevel.supportStaticWebsite {
		log.Printf("[DEBUG] waiting for the Static Website to become available")
		pollerType, err := custompollers.NewDataPlaneStaticWebsiteAvailabilityPoller(ctx, client, account)
		if err != nil {
			return fmt.Errorf("building Static Website Poller: %+v", err)
		}
		poller := pollers.NewPoller(pollerType, initialDelayDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("waiting for the Static Website to become available: %+v", err)
		}
	}

	return nil
}
