// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_app

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func expandSpringCloudAppIdentity(input []any) (*appplatform.ManagedIdentityProperties, error) {
	config, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := appplatform.ManagedIdentityProperties{
		Type: appplatform.ManagedIdentityType(string(config.Type)),
	}
	if config.Type == identity.TypeUserAssigned || config.Type == identity.TypeSystemAssignedUserAssigned {
		out.UserAssignedIdentities = make(map[string]*appplatform.UserAssignedManagedIdentity)
		for k := range config.IdentityIds {
			out.UserAssignedIdentities[k] = &appplatform.UserAssignedManagedIdentity{
				// intentionally empty
			}
		}
	}

	return &out, nil
}

func expandSpringCloudAppPersistentDisk(input []any) *appplatform.PersistentDisk {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	raw := input[0].(map[string]any)
	return &appplatform.PersistentDisk{
		SizeInGB:  pointer.To(int32(raw["size_in_gb"].(int))),
		MountPath: pointer.To(raw["mount_path"].(string)),
	}
}

func expandAppCustomPersistentDiskResourceArray(input []any, id appplatform_rm.AppId) *[]appplatform.CustomPersistentDiskResource {
	results := make([]appplatform.CustomPersistentDiskResource, 0)
	for _, item := range input {
		v := item.(map[string]any)
		results = append(results, appplatform.CustomPersistentDiskResource{
			StorageID: pointer.To(appplatform_rm.NewStorageID(id.SubscriptionId, id.ResourceGroupName, id.SpringName, v["storage_name"].(string)).ID()),
			CustomPersistentDiskProperties: &appplatform.AzureFileVolume{
				ShareName:    pointer.To(v["share_name"].(string)),
				MountPath:    pointer.To(v["mount_path"].(string)),
				ReadOnly:     pointer.To(v["read_only_enabled"].(bool)),
				MountOptions: pluginsdk.ExpandStringSlice(v["mount_options"].(*pluginsdk.Set).List()),
			},
		})
	}
	return &results
}

func expandSpringCloudAppIngressSetting(input []any) *appplatform.IngressSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	raw := input[0].(map[string]any)

	return &appplatform.IngressSettings{
		ReadTimeoutInSeconds: pointer.To(int32(raw["read_timeout_in_seconds"].(int))),
		SendTimeoutInSeconds: pointer.To(int32(raw["send_timeout_in_seconds"].(int))),
		SessionAffinity:      appplatform.SessionAffinity(raw["session_affinity"].(string)),
		SessionCookieMaxAge:  pointer.To(int32(raw["session_cookie_max_age"].(int))),
		BackendProtocol:      appplatform.BackendProtocol(raw["backend_protocol"].(string)),
	}
}
