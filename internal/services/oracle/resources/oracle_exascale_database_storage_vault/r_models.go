// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_exascale_database_storage_vault

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

type ExascaleDatabaseStorageVaultResourceModel struct {
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`
	Zones             zones.Schema      `tfschema:"zones"`

	AdditionalFlashCachePercentage int64                                 `tfschema:"additional_flash_cache_percentage"`
	Description                    string                                `tfschema:"description"`
	DisplayName                    string                                `tfschema:"display_name"`
	HighCapacityDatabaseStorage    []ExascaleDatabaseStorageDetailsModel `tfschema:"high_capacity_database_storage"`

	TimeZone string `tfschema:"time_zone"`
}

func (Resource) ModelObject() any {
	return &Resource{}
}
