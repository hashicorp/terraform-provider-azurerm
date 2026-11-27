// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_discovery_workspace

type StorageDiscoveryWorkspaceModel struct {
	Name              string                       `tfschema:"name"`
	ResourceGroupName string                       `tfschema:"resource_group_name"`
	Location          string                       `tfschema:"location"`
	Scope             []StorageDiscoveryScopeModel `tfschema:"scope"`
	WorkspaceRoots    []string                     `tfschema:"workspace_roots"`
	Description       string                       `tfschema:"description"`
	Sku               string                       `tfschema:"sku"`
	Tags              map[string]string            `tfschema:"tags"`
}

type StorageDiscoveryScopeModel struct {
	DisplayName   string            `tfschema:"display_name"`
	ResourceTypes []string          `tfschema:"resource_types"`
	TagKeysOnly   []string          `tfschema:"tag_keys_only"`
	Tags          map[string]string `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &StorageDiscoveryWorkspaceModel{}
}
