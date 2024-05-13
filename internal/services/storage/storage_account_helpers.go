package storage

import (
	"fmt"
	"net/url"
	"slices"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resolveStorageAccountServiceSupportLevel(kind storage.Kind, tier storage.SkuTier, replicationType string) storageAccountServiceSupportLevel {
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

func flattenAzureRmStorageAccountIdentity(input *storage.Identity) (*[]interface{}, error) {
	var config *identity.SystemAndUserAssignedMap

	if input != nil {
		config = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: nil,
		}

		// work around the Swagger defining `SystemAssigned,UserAssigned` rather than `SystemAssigned, UserAssigned`
		if input.Type == storage.IdentityTypeSystemAssignedUserAssigned {
			config.Type = identity.TypeSystemAssignedUserAssigned
		}

		if input.PrincipalID != nil {
			config.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			config.TenantId = *input.TenantID
		}
		identityIds := make(map[string]identity.UserAssignedIdentityDetails)
		for k, v := range input.UserAssignedIdentities {
			if v == nil {
				continue
			}

			details := identity.UserAssignedIdentityDetails{}

			if v.ClientID != nil {
				details.ClientId = utils.String(*v.ClientID)
			}
			if v.PrincipalID != nil {
				details.PrincipalId = utils.String(*v.PrincipalID)
			}

			identityIds[k] = details
		}

		config.IdentityIds = identityIds
	}

	return identity.FlattenSystemAndUserAssignedMap(config)
}

func flattenAndSetAzureRmStorageAccountPrimaryEndpoints(d *pluginsdk.ResourceData, primary *storage.Endpoints, routingInputs *storage.RoutingPreference) error {
	if primary == nil {
		return fmt.Errorf("primary endpoints should not be empty")
	}

	if err := setEndpointAndHost(d, "primary", primary.Blob, "blob"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "primary", primary.Dfs, "dfs"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "primary", primary.File, "file"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "primary", primary.Queue, "queue"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "primary", primary.Table, "table"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "primary", primary.Web, "web"); err != nil {
		return err
	}

	// below null check is to avoid nullpointer scenarios when either of publish_internet_endpoints
	// or publish_microsoft_endpoints or both aren't set
	if routingInputs != nil && routingInputs.PublishInternetEndpoints != nil && *routingInputs.PublishInternetEndpoints {
		if err := setEndpointAndHost(d, "primary", primary.InternetEndpoints.Blob, "blob_internet"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "primary", primary.InternetEndpoints.Dfs, "dfs_internet"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "primary", primary.InternetEndpoints.File, "file_internet"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "primary", primary.InternetEndpoints.Web, "web_internet"); err != nil {
			return err
		}
	} else {
		if err := setEndpointAndHost(d, "primary", nil, "blob_internet"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "primary", nil, "dfs_internet"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "primary", nil, "file_internet"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "primary", nil, "web_internet"); err != nil {
			return err
		}
	}

	if routingInputs != nil && routingInputs.PublishMicrosoftEndpoints != nil && *routingInputs.PublishMicrosoftEndpoints {
		if err := setEndpointAndHost(d, "primary", primary.MicrosoftEndpoints.Blob, "blob_microsoft"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "primary", primary.MicrosoftEndpoints.Dfs, "dfs_microsoft"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "primary", primary.MicrosoftEndpoints.File, "file_microsoft"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "primary", primary.MicrosoftEndpoints.Web, "web_microsoft"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "primary", primary.MicrosoftEndpoints.Table, "table_microsoft"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "primary", primary.MicrosoftEndpoints.Queue, "queue_microsoft"); err != nil {
			return err
		}
	} else {
		if err := setEndpointAndHost(d, "primary", nil, "blob_microsoft"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "primary", nil, "dfs_microsoft"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "primary", nil, "file_microsoft"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "primary", nil, "web_microsoft"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "primary", nil, "table_microsoft"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "primary", nil, "queue_microsoft"); err != nil {
			return err
		}
	}

	return nil
}

func flattenAndSetAzureRmStorageAccountSecondaryEndpoints(d *pluginsdk.ResourceData, secondary *storage.Endpoints, routingInputs *storage.RoutingPreference) error {
	if secondary == nil {
		return nil
	}

	if err := setEndpointAndHost(d, "secondary", secondary.Blob, "blob"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "secondary", secondary.Dfs, "dfs"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "secondary", secondary.File, "file"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "secondary", secondary.Queue, "queue"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "secondary", secondary.Table, "table"); err != nil {
		return err
	}
	if err := setEndpointAndHost(d, "secondary", secondary.Web, "web"); err != nil {
		return err
	}

	if v := routingInputs; v != nil && v.PublishInternetEndpoints != nil && *v.PublishInternetEndpoints {
		if err := setEndpointAndHost(d, "secondary", secondary.InternetEndpoints.Blob, "blob_internet"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "secondary", secondary.InternetEndpoints.Dfs, "dfs_internet"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "secondary", secondary.InternetEndpoints.File, "file_internet"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "secondary", secondary.InternetEndpoints.Web, "web_internet"); err != nil {
			return err
		}
	}

	if routingInputs != nil && *routingInputs.PublishMicrosoftEndpoints {
		if err := setEndpointAndHost(d, "secondary", secondary.MicrosoftEndpoints.Blob, "blob_microsoft"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "secondary", secondary.MicrosoftEndpoints.Dfs, "dfs_microsoft"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "secondary", secondary.MicrosoftEndpoints.File, "file_microsoft"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "secondary", secondary.MicrosoftEndpoints.Web, "web_microsoft"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "secondary", secondary.MicrosoftEndpoints.Table, "table_microsoft"); err != nil {
			return err
		}

		if err := setEndpointAndHost(d, "secondary", secondary.MicrosoftEndpoints.Queue, "queue_microsoft"); err != nil {
			return err
		}
	}
	return nil
}

func setEndpointAndHost(d *pluginsdk.ResourceData, ordinalString string, endpointType *string, typeString string) error {
	var endpoint, host string
	if v := endpointType; v != nil {
		endpoint = *v

		u, err := url.Parse(*v)
		if err != nil {
			return fmt.Errorf("invalid %s endpoint for parsing: %q", typeString, *v)
		}
		host = u.Host
	}

	// lintignore: R001
	d.Set(fmt.Sprintf("%s_%s_endpoint", ordinalString, typeString), endpoint)
	// lintignore: R001
	d.Set(fmt.Sprintf("%s_%s_host", ordinalString, typeString), host)
	return nil
}

func getBlobConnectionString(blobEndpoint *string, acctName *string, acctKey *string) string {
	var endpoint string
	if blobEndpoint != nil {
		endpoint = *blobEndpoint
	}

	var name string
	if acctName != nil {
		name = *acctName
	}

	var key string
	if acctKey != nil {
		key = *acctKey
	}

	return fmt.Sprintf("DefaultEndpointsProtocol=https;BlobEndpoint=%s;AccountName=%s;AccountKey=%s", endpoint, name, key)
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

func flattenEdgeZone(input *storage.ExtendedLocation) string {
	if input == nil || input.Type != storage.ExtendedLocationTypesEdgeZone || input.Name == nil {
		return ""
	}
	return edgezones.NormalizeNilable(input.Name)
}
