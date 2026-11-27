// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_mysql_flexible_server

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Backup Instance for the MySQL Flexible Server.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&basebackuppolicyresources.BackupVaultId{}, "Backup Vault"),

		"backup_policy_id": commonschema.ResourceIDReferenceRequired(&basebackuppolicyresources.BackupPolicyId{}, "Backup Policy"),

		"server_id": commonschema.ResourceIDReferenceRequiredForceNew(&servers.FlexibleServerId{}, "MySQL Flexible Server"),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"protection_state": {
			Description: "The protection state of the Backup Instance MySQL Flexible Server.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
