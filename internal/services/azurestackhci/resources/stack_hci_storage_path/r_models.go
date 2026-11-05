// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_storage_path

func (Resource) ModelObject() any {
	return &StackHCIStoragePathResourceModel{}
}

type StackHCIStoragePathResourceModel struct {
	Name              string         `tfschema:"name"`
	ResourceGroupName string         `tfschema:"resource_group_name"`
	Location          string         `tfschema:"location"`
	CustomLocationId  string         `tfschema:"custom_location_id"`
	Path              string         `tfschema:"path"`
	Tags              map[string]any `tfschema:"tags"`
}
