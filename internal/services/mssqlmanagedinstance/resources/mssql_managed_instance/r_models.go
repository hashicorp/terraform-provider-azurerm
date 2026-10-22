// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type MsSqlManagedInstanceModel struct {
	AdministratorLogin                string                              `tfschema:"administrator_login"`
	AdministratorLoginPassword        string                              `tfschema:"administrator_login_password"`
	Collation                         string                              `tfschema:"collation"`
	DnsZonePartnerId                  string                              `tfschema:"dns_zone_partner_id"`
	DnsZone                           string                              `tfschema:"dns_zone"`
	Fqdn                              string                              `tfschema:"fqdn"`
	Identity                          []identity.SystemOrUserAssignedList `tfschema:"identity"`
	GeneralPurposeV2Enabled           bool                                `tfschema:"general_purpose_v2_enabled"`
	LicenseType                       string                              `tfschema:"license_type"`
	Location                          string                              `tfschema:"location"`
	MaintenanceConfigurationName      string                              `tfschema:"maintenance_configuration_name"`
	MinimumTlsVersion                 string                              `tfschema:"minimum_tls_version"`
	Name                              string                              `tfschema:"name"`
	ProxyOverride                     string                              `tfschema:"proxy_override"`
	PublicDataEndpointEnabled         bool                                `tfschema:"public_data_endpoint_enabled"`
	ResourceGroupName                 string                              `tfschema:"resource_group_name"`
	ServicePrincipalType              string                              `tfschema:"service_principal_type"`
	SkuName                           string                              `tfschema:"sku_name"`
	StorageAccountType                string                              `tfschema:"storage_account_type"`
	StorageIOps                       int64                               `tfschema:"storage_iops"`
	StorageSizeInGb                   int64                               `tfschema:"storage_size_in_gb"`
	SubnetId                          string                              `tfschema:"subnet_id"`
	TimezoneId                        string                              `tfschema:"timezone_id"`
	VCores                            int64                               `tfschema:"vcores"`
	AzureActiveDirectoryAdministrator []AzureActiveDirectoryAdministrator `tfschema:"azure_active_directory_administrator"`
	ZoneRedundantEnabled              bool                                `tfschema:"zone_redundant_enabled"`
	Tags                              map[string]string                   `tfschema:"tags"`
	DatabaseFormat                    string                              `tfschema:"database_format"`
	HybridSecondaryUsage              string                              `tfschema:"hybrid_secondary_usage"`
}

func (r Resource) ModelObject() any {
	return &MsSqlManagedInstanceModel{}
}
