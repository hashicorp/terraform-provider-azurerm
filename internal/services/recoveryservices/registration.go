// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/backup_container_storage_account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/backup_policy_file_share"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/backup_policy_vm"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/backup_policy_vm_workload"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/backup_protected_file_share"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/backup_protected_vm"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/recovery_services_vault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/recovery_services_vault_resource_guard_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/site_recovery_fabric"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/site_recovery_hyperv_network_mapping"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/site_recovery_hyperv_replication_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/site_recovery_hyperv_replication_policy_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/site_recovery_network_mapping"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/site_recovery_protection_container"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/site_recovery_protection_container_mapping"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/site_recovery_replicated_vm"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/site_recovery_replication_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/site_recovery_replication_recovery_plan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/site_recovery_services_vault_hyperv_site"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/site_recovery_vmware_replicated_vm"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/site_recovery_vmware_replication_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/resources/site_recovery_vmware_replication_policy_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/recovery-services"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		recovery_services_vault.DataSource{},
		site_recovery_replication_recovery_plan.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		backup_policy_vm_workload.Resource{},
		site_recovery_hyperv_network_mapping.Resource{},
		site_recovery_hyperv_replication_policy_association.Resource{},
		site_recovery_services_vault_hyperv_site.Resource{},
		site_recovery_hyperv_replication_policy.Resource{},
		site_recovery_replication_recovery_plan.Resource{},
		recovery_services_vault_resource_guard_association.Resource{},
		site_recovery_vmware_replicated_vm.Resource{},
		site_recovery_vmware_replication_policy_association.Resource{},
		site_recovery_vmware_replication_policy.Resource{},
	}
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Recovery Services"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Recovery Services",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_backup_policy_file_share":           backup_policy_file_share.RegisterDataSource(),
		"azurerm_backup_policy_vm":                   backup_policy_vm.RegisterDataSource(),
		"azurerm_site_recovery_fabric":               site_recovery_fabric.RegisterDataSource(),
		"azurerm_site_recovery_protection_container": site_recovery_protection_container.RegisterDataSource(),
		"azurerm_site_recovery_replication_policy":   site_recovery_replication_policy.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	// todo - this package should probably be split into backup, recovery, and site recovery?
	return map[string]*pluginsdk.Resource{
		"azurerm_backup_container_storage_account":           backup_container_storage_account.RegisterResource(),
		"azurerm_backup_policy_file_share":                   backup_policy_file_share.RegisterResource(),
		"azurerm_backup_policy_vm":                           backup_policy_vm.RegisterResource(),
		"azurerm_backup_protected_file_share":                backup_protected_file_share.RegisterResource(),
		"azurerm_backup_protected_vm":                        backup_protected_vm.RegisterResource(),
		"azurerm_recovery_services_vault":                    recovery_services_vault.RegisterResource(),
		"azurerm_site_recovery_fabric":                       site_recovery_fabric.RegisterResource(),
		"azurerm_site_recovery_network_mapping":              site_recovery_network_mapping.RegisterResource(),
		"azurerm_site_recovery_protection_container":         site_recovery_protection_container.RegisterResource(),
		"azurerm_site_recovery_protection_container_mapping": site_recovery_protection_container_mapping.RegisterResource(),
		"azurerm_site_recovery_replicated_vm":                site_recovery_replicated_vm.RegisterResource(),
		"azurerm_site_recovery_replication_policy":           site_recovery_replication_policy.RegisterResource(),
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
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

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{
		backup_policy_vm.ListResource{},
	}
}
