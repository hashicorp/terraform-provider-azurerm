// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_vault_resource_guard_association

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardproxybaseresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const (
	dataProtectionBackupVaultResourceGuardAssociationDeleteRequestName = "default"
	dataProtectionBackupVaultResourceGuardAssociationProxyName         = "DppResourceGuardProxy"
	dataProtectionBackupVaultResourceGuardAssociationResourceType      = "azurerm_data_protection_backup_vault_resource_guard_association"
)

//go:generate env GOFILE= go run ../../../../tools/generator-tests resourceidentity -resource-name data_protection_backup_vault_resource_guard_association -service-package-name dataprotection -compare-values "subscription_id:data_protection_backup_vault_id,resource_group_name:data_protection_backup_vault_id,backup_vault_name:data_protection_backup_vault_id" -known-values "name:DppResourceGuardProxy"

type Resource struct{}

var (
	_ sdk.Resource             = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &resourceguardproxybaseresources.BackupResourceGuardProxyId{}
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return resourceguardproxybaseresources.ValidateBackupResourceGuardProxyID
}

func (r Resource) ResourceType() string {
	return dataProtectionBackupVaultResourceGuardAssociationResourceType
}
