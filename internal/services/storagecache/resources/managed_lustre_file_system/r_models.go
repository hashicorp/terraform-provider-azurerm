// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_lustre_file_system

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ManagedLustreFileSystemModel struct {
	Name                string                       `tfschema:"name"`
	ResourceGroupName   string                       `tfschema:"resource_group_name"`
	Location            string                       `tfschema:"location"`
	HsmSetting          []HsmSetting                 `tfschema:"hsm_setting"`
	Identity            []identity.ModelUserAssigned `tfschema:"identity"`
	EncryptionKey       []EncryptionKey              `tfschema:"encryption_key"`
	MaintenanceWindow   []MaintenanceWindow          `tfschema:"maintenance_window"`
	MgsAddress          string                       `tfschema:"mgs_address"`
	SkuName             string                       `tfschema:"sku_name"`
	RootSquashSettings  []RootSquashSetting          `tfschema:"root_squash"`
	StorageCapacityInTb int64                        `tfschema:"storage_capacity_in_tb"`
	SubnetId            string                       `tfschema:"subnet_id"`
	Zones               []string                     `tfschema:"zones"`
	Tags                map[string]string            `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &ManagedLustreFileSystemModel{}
}
