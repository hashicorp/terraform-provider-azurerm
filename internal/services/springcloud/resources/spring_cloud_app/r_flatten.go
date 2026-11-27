// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_app

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func flattenSpringCloudAppIdentity(input *appplatform.ManagedIdentityProperties) (*[]any, error) {
	var transform *identity.SystemAndUserAssignedMap
	if input != nil {
		transform = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
		for k, v := range input.UserAssignedIdentities {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientID,
				PrincipalId: v.PrincipalID,
			}
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(transform)
}

func flattenSpringCloudAppPersistentDisk(input *appplatform.PersistentDisk) []any {
	if input == nil {
		return make([]any, 0)
	}

	sizeInGB := 0
	if input.SizeInGB != nil {
		sizeInGB = int(*input.SizeInGB)
	}

	return []any{
		map[string]any{
			"size_in_gb": sizeInGB,
			"mount_path": pointer.From(input.MountPath),
		},
	}
}

func flattenAppCustomPersistentDiskResourceArray(input *[]appplatform.CustomPersistentDiskResource) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var storageName string
		if item.StorageID != nil {
			// The returned value has inconsistent casing
			// TODO: Remove the normalization codes once the following issue is fixed.
			// Issue: https://github.com/Azure/azure-rest-api-specs/issues/22205
			if id, err := appplatform_rm.ParseStorageIDInsensitively(*item.StorageID); err == nil {
				storageName = id.StorageName
			}
		}
		var mountPath string
		var shareName string
		var readOnly bool
		var mountOptions *[]string
		if item.CustomPersistentDiskProperties != nil {
			if prop, ok := item.CustomPersistentDiskProperties.AsAzureFileVolume(); ok && prop != nil {
				if prop.MountPath != nil {
					mountPath = *prop.MountPath
				}
				if prop.ShareName != nil {
					shareName = *prop.ShareName
				}
				if prop.ReadOnly != nil {
					readOnly = *prop.ReadOnly
				}
				mountOptions = prop.MountOptions
			}
		}

		results = append(results, map[string]any{
			"storage_name":      storageName,
			"mount_path":        mountPath,
			"share_name":        shareName,
			"mount_options":     set.FromStringSliceNilable(mountOptions),
			"read_only_enabled": readOnly,
		})
	}
	return results
}

func flattenSpringCloudAppIngressSettings(input *appplatform.IngressSettings) any {
	if input == nil {
		return make([]any, 0)
	}
	var readTimeout, sendTimeout, maxAge int32
	backendProtocol := string(input.BackendProtocol)
	sessionAffinity := string(input.SessionAffinity)
	if input.ReadTimeoutInSeconds != nil {
		readTimeout = *input.ReadTimeoutInSeconds
	}
	if input.SendTimeoutInSeconds != nil {
		sendTimeout = *input.SendTimeoutInSeconds
	}
	if input.SessionCookieMaxAge != nil {
		maxAge = *input.SessionCookieMaxAge
	}
	return []any{map[string]any{
		"backend_protocol":        backendProtocol,
		"read_timeout_in_seconds": readTimeout,
		"send_timeout_in_seconds": sendTimeout,
		"session_affinity":        sessionAffinity,
		"session_cookie_max_age":  maxAge,
	}}
}
