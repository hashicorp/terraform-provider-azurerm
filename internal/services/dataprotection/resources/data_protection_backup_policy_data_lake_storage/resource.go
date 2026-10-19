// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_data_lake_storage

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name data_protection_backup_policy_data_lake_storage -service-package-name dataprotection -properties "name" -compare-values "subscription_id:data_protection_backup_vault_id,resource_group_name:data_protection_backup_vault_id,backup_vault_name:data_protection_backup_vault_id"

type Criteria struct {
	AbsoluteCriteria     string   `tfschema:"absolute_criteria"`
	DaysOfWeek           []string `tfschema:"days_of_week"`
	MonthsOfYear         []string `tfschema:"months_of_year"`
	ScheduledBackupTimes []string `tfschema:"scheduled_backup_times"`
	WeeksOfMonth         []string `tfschema:"weeks_of_month"`
}

type Resource struct{}

var (
	_ sdk.Resource             = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &basebackuppolicyresources.BackupPolicyId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_data_protection_backup_policy_data_lake_storage"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return basebackuppolicyresources.ValidateBackupPolicyID
}
