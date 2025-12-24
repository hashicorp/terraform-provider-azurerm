// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{
		newDataProtectionBackupInstanceProtectAction,
	}
}

func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}

func (r Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

func (r Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

var (
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.FrameworkServiceRegistration               = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/data-protection"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "DataProtection"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"DataProtection",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_data_protection_backup_vault": dataSourceDataProtectionBackupVault(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	resources := map[string]*pluginsdk.Resource{
		"azurerm_data_protection_backup_instance_blob_storage": resourceDataProtectionBackupInstanceBlobStorage(),
		"azurerm_data_protection_backup_instance_disk":         resourceDataProtectionBackupInstanceDisk(),
		"azurerm_data_protection_backup_policy_blob_storage":   resourceDataProtectionBackupPolicyBlobStorage(),
		"azurerm_data_protection_backup_policy_disk":           resourceDataProtectionBackupPolicyDisk(),
		"azurerm_data_protection_backup_vault":                 resourceDataProtectionBackupVault(),
		"azurerm_data_protection_resource_guard":               resourceDataProtectionResourceGuard(),
	}

	if !features.FivePointOh() {
		resources["azurerm_data_protection_backup_instance_postgresql"] = resourceDataProtectionBackupInstancePostgreSQL()
		resources["azurerm_data_protection_backup_policy_postgresql"] = resourceDataProtectionBackupPolicyPostgreSQL()
	}

	return resources
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		DataProtectionBackupPolicyDataLakeStorageResource{},
		DataProtectionBackupPolicyKubernatesClusterResource{},
		DataProtectionBackupPolicyMySQLFlexibleServerResource{},
		DataProtectionBackupPolicyPostgreSQLFlexibleServerResource{},
		DataProtectionBackupInstanceKubernatesClusterResource{},
		DataProtectionBackupInstanceMySQLFlexibleServerResource{},
		DataProtectionBackupInstancePostgreSQLFlexibleServerResource{},
		DataProtectionBackupVaultCustomerManagedKeyResource{},
	}
}
