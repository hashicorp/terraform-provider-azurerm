package storage

import (
	"slices"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func expandAzureRmStorageAccountIdentity(input []interface{}) (*storage.Identity, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := storage.Identity{
		Type: storage.IdentityType(string(expanded.Type)),
	}

	// work around the Swagger defining `SystemAssigned,UserAssigned` rather than `SystemAssigned, UserAssigned`
	if expanded.Type == identity.TypeSystemAssignedUserAssigned {
		out.Type = storage.IdentityTypeSystemAssignedUserAssigned
	}

	// 'Failed to perform resource identity operation. Status: 'BadRequest'. Response:
	// {"error":{"code":"BadRequest",
	//  "message":"The request format was unexpected, a non-UserAssigned identity type should not contain: userAssignedIdentities"
	// }}
	// Upstream issue: https://github.com/Azure/azure-rest-api-specs/issues/17650
	if len(expanded.IdentityIds) > 0 {
		userAssignedIdentities := make(map[string]*storage.UserAssignedIdentity)
		for id := range expanded.IdentityIds {
			userAssignedIdentities[id] = &storage.UserAssignedIdentity{}
		}
		out.UserAssignedIdentities = userAssignedIdentities
	}

	return &out, nil
}

func expandEdgeZone(input string) *storage.ExtendedLocation {
	normalized := edgezones.Normalize(input)
	if normalized == "" {
		return nil
	}

	return &storage.ExtendedLocation{
		Name: utils.String(normalized),
		Type: storage.ExtendedLocationTypesEdgeZone,
	}
}

func availableFunctionalityForAccountLegacy(kind storage.Kind, tier storage.SkuTier, replicationType string) storageAccountServiceSupportLevel {
	// FileStorage doesn't support blob
	supportBlob := kind != storage.KindFileStorage

	// Queue is only supported for Storage and StorageV2, in Standard sku tier.
	supportQueue := tier == storage.SkuTierStandard && (kind == storage.KindStorageV2 ||
		(kind == storage.KindStorage &&
			// Per local test, only LRS/GRS/RAGRS Storage V1 accounts support queue endpoint.
			// GZRS and RAGZRS is invalid, while ZRS is valid but has no queue endpoint.
			slices.Contains([]string{"LRS", "GRS", "RAGRS"}, replicationType)))

	// File share is only supported for StorageV2 and FileStorage.
	// See: https://docs.microsoft.com/en-us/azure/storage/files/storage-files-planning#management-concepts
	// Per test, the StorageV2 with Premium sku tier also doesn't support file share.
	supportShare := kind == storage.KindFileStorage || (tier != storage.SkuTierPremium && (kind == storage.KindStorageV2 ||
		(kind == storage.KindStorage &&
			// Per local test, only LRS/GRS/RAGRS Storage V1 accounts support file endpoint.
			// GZRS and RAGZRS is invalid, while ZRS is valid but has no file endpoint.
			slices.Contains([]string{"LRS", "GRS", "RAGRS"}, replicationType))))

	// Static Website is only supported for StorageV2 (not for Storage(v1)) and BlockBlobStorage
	supportStaticWebSite := kind == storage.KindStorageV2 || kind == storage.KindBlockBlobStorage

	return storageAccountServiceSupportLevel{
		supportBlob:          supportBlob,
		supportQueue:         supportQueue,
		supportShare:         supportShare,
		supportStaticWebsite: supportStaticWebSite,
	}
}
