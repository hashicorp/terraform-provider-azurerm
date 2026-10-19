// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_virtual_hard_disk

func (Resource) ModelObject() any {
	return &StackHCIVirtualHardDiskResourceModel{}
}

type StackHCIVirtualHardDiskResourceModel struct {
	Name                  string         `tfschema:"name"`
	ResourceGroupName     string         `tfschema:"resource_group_name"`
	Location              string         `tfschema:"location"`
	CustomLocationId      string         `tfschema:"custom_location_id"`
	BlockSizeInBytes      int64          `tfschema:"block_size_in_bytes"`
	DiskFileFormat        string         `tfschema:"disk_file_format"`
	DiskSizeInGB          int64          `tfschema:"disk_size_in_gb"`
	DynamicEnabled        bool           `tfschema:"dynamic_enabled"`
	HypervGeneration      string         `tfschema:"hyperv_generation"`
	LogicalSectorInBytes  int64          `tfschema:"logical_sector_in_bytes"`
	PhysicalSectorInBytes int64          `tfschema:"physical_sector_in_bytes"`
	StoragePathId         string         `tfschema:"storage_path_id"`
	Tags                  map[string]any `tfschema:"tags"`
}
