// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_data_lake_storage

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name data_protection_backup_instance_data_lake_storage -service-package-name dataprotection -properties "name" -compare-values "subscription_id:data_protection_backup_vault_id,resource_group_name:data_protection_backup_vault_id,backup_vault_name:data_protection_backup_vault_id"

type BackupDatasourceParameters struct {
	IncludedNamespaces          []string `tfschema:"included_namespaces"`
	IncludedResourceTypes       []string `tfschema:"included_resource_types"`
	ExcludedNamespaces          []string `tfschema:"excluded_namespaces"`
	ExcludedResourceTypes       []string `tfschema:"excluded_resource_types"`
	LabelSelectors              []string `tfschema:"label_selectors"`
	VolumeSnapshotEnabled       bool     `tfschema:"volume_snapshot_enabled"`
	ClusterScopeResourceEnabled bool     `tfschema:"cluster_scoped_resources_enabled"`
}

type Resource struct{}

var (
	_ sdk.Resource             = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &backupinstanceresources.BackupInstanceId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_data_protection_backup_instance_data_lake_storage"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backupinstanceresources.ValidateBackupInstanceID
}
