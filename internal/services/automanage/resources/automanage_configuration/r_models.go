// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automanage_configuration

type ConfigurationModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`

	Antimalware               []AntimalwareConfiguration           `tfschema:"antimalware"`
	AzureSecurityBaseline     []AzureSecurityBaselineConfiguration `tfschema:"azure_security_baseline"`
	Backup                    []BackupConfiguration                `tfschema:"backup"`
	LogAnalyticsEnabled       bool                                 `tfschema:"log_analytics_enabled"`
	AutomationAccountEnabled  bool                                 `tfschema:"automation_account_enabled"`
	BootDiagnosticsEnabled    bool                                 `tfschema:"boot_diagnostics_enabled"`
	DefenderForCloudEnabled   bool                                 `tfschema:"defender_for_cloud_enabled"`
	GuestConfigurationEnabled bool                                 `tfschema:"guest_configuration_enabled"`
	StatusChangeAlertEnabled  bool                                 `tfschema:"status_change_alert_enabled"`

	Location string            `tfschema:"location"`
	Tags     map[string]string `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &ConfigurationModel{}
}
