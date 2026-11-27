// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_lustre_file_system

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/amlfilesystems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageCache.AmlFilesystems

			id, err := amlfilesystems.ParseAmlFilesystemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ManagedLustreFileSystemModel{}
			if model := resp.Model; model != nil {
				state.Name = id.AmlFilesystemName
				state.ResourceGroupName = id.ResourceGroupName
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				identity, err := identity.FlattenUserAssignedMapToModel(model.Identity)
				if err != nil {
					return err
				}
				state.Identity = pointer.From(identity)

				if properties := model.Properties; properties != nil {
					state.SubnetId = properties.FilesystemSubnet
					state.StorageCapacityInTb = int64(properties.StorageCapacityTiB)
					state.MaintenanceWindow = flattenManagedLustreFileSystemMaintenanceWindow(properties.MaintenanceWindow)
					if properties.ClientInfo != nil && properties.ClientInfo.MgsAddress != nil {
						state.MgsAddress = *properties.ClientInfo.MgsAddress
					}
					state.HsmSetting = flattenManagedLustreFileSystemHsmSetting(properties.Hsm)
					state.Zones = pointer.From(model.Zones)
					state.EncryptionKey = flattenManagedLustreFileSystemEncryptionKey(properties.EncryptionSettings)
					state.RootSquashSettings = flattenRootSquashSettings(properties.RootSquashSettings)

					if v := model.Sku; v != nil {
						state.SkuName = pointer.From(v.Name)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

type HsmSetting struct {
	ContainerId        string `tfschema:"container_id"`
	LoggingContainerId string `tfschema:"logging_container_id"`
	ImportPrefix       string `tfschema:"import_prefix"`
}

type EncryptionKey struct {
	KeyURL        string `tfschema:"key_url"`
	SourceVaultId string `tfschema:"source_vault_id"`
}

type RootSquashSetting struct {
	Mode            string `tfschema:"mode"`
	NoSquashNidList string `tfschema:"no_squash_nids"`
	SquashGID       int64  `tfschema:"squash_gid"`
	SquashUID       int64  `tfschema:"squash_uid"`
}

type MaintenanceWindow struct {
	DayOfWeek      amlfilesystems.MaintenanceDayOfWeekType `tfschema:"day_of_week"`
	TimeOfDayInUTC string                                  `tfschema:"time_of_day_in_utc"`
}

type SkuProperties struct {
	Name           string
	MinimumStorage int
}

func GetSkuPropertiesByName(skuName string) *SkuProperties {
	for _, sku := range PossibleSkuProperties() {
		if skuName == sku.Name {
			return pointer.To(sku)
		}
	}

	return nil
}

func PossibleValuesForSkuName() []string {
	skus := make([]string, 0)

	for _, sku := range PossibleSkuProperties() {
		skus = append(skus, sku.Name)
	}

	return skus
}

func PossibleSkuProperties() []SkuProperties {
	return []SkuProperties{
		{
			Name:           "AMLFS-Durable-Premium-40",
			MinimumStorage: 48,
		},
		{
			Name:           "AMLFS-Durable-Premium-125",
			MinimumStorage: 16,
		},
		{
			Name:           "AMLFS-Durable-Premium-250",
			MinimumStorage: 8,
		},
		{
			Name:           "AMLFS-Durable-Premium-500",
			MinimumStorage: 4,
		},
	}
}
