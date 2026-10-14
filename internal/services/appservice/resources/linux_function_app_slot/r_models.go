// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package linux_function_app_slot

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

type LinuxFunctionAppSlotModel struct {
	Name                               string                                     `tfschema:"name"`
	FunctionAppID                      string                                     `tfschema:"function_app_id"`
	ServicePlanID                      string                                     `tfschema:"service_plan_id"`
	StorageAccountName                 string                                     `tfschema:"storage_account_name"`
	StorageAccountKey                  string                                     `tfschema:"storage_account_access_key"`
	StorageUsesMSI                     bool                                       `tfschema:"storage_uses_managed_identity"` // Storage uses MSI not account key
	StorageKeyVaultSecretID            string                                     `tfschema:"storage_key_vault_secret_id"`
	AppSettings                        map[string]string                          `tfschema:"app_settings"`
	AuthSettings                       []helpers.AuthSettings                     `tfschema:"auth_settings"`
	AuthV2Settings                     []helpers.AuthV2Settings                   `tfschema:"auth_settings_v2"`
	Backup                             []helpers.Backup                           `tfschema:"backup"` // Not supported on Dynamic or Basic plans
	BuiltinLogging                     bool                                       `tfschema:"builtin_logging_enabled"`
	ClientCertEnabled                  bool                                       `tfschema:"client_certificate_enabled"`
	ClientCertMode                     string                                     `tfschema:"client_certificate_mode"`
	ClientCertExclusionPaths           string                                     `tfschema:"client_certificate_exclusion_paths"`
	ConnectionStrings                  []helpers.ConnectionString                 `tfschema:"connection_string"`
	DailyMemoryTimeQuota               int64                                      `tfschema:"daily_memory_time_quota"` // TODO - Value ignored in for linux apps, even in Consumption plans?
	Enabled                            bool                                       `tfschema:"enabled"`
	FunctionExtensionsVersion          string                                     `tfschema:"functions_extension_version"`
	ForceDisableContentShare           bool                                       `tfschema:"content_share_force_disabled"`
	HttpsOnly                          bool                                       `tfschema:"https_only"`
	KeyVaultReferenceIdentityID        string                                     `tfschema:"key_vault_reference_identity_id"`
	SiteConfig                         []helpers.SiteConfigLinuxFunctionAppSlot   `tfschema:"site_config"`
	Tags                               map[string]string                          `tfschema:"tags"`
	VirtualNetworkBackupRestoreEnabled bool                                       `tfschema:"virtual_network_backup_restore_enabled"`
	VirtualNetworkSubnetID             string                                     `tfschema:"virtual_network_subnet_id"`
	CustomDomainVerificationId         string                                     `tfschema:"custom_domain_verification_id"`
	HostingEnvId                       string                                     `tfschema:"hosting_environment_id"`
	DefaultHostname                    string                                     `tfschema:"default_hostname"`
	Kind                               string                                     `tfschema:"kind"`
	OutboundIPAddresses                string                                     `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList              []string                                   `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses        string                                     `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList      []string                                   `tfschema:"possible_outbound_ip_address_list"`
	PublicNetworkAccess                bool                                       `tfschema:"public_network_access_enabled"`
	PublishingDeployBasicAuthEnabled   bool                                       `tfschema:"webdeploy_publish_basic_authentication_enabled"`
	PublishingFTPBasicAuthEnabled      bool                                       `tfschema:"ftp_publish_basic_authentication_enabled"`
	SiteCredentials                    []helpers.SiteCredential                   `tfschema:"site_credential"`
	StorageAccounts                    []helpers.StorageAccount                   `tfschema:"storage_account"`
	Identity                           []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	VnetImagePullEnabled               bool                                       `tfschema:"vnet_image_pull_enabled"`
	EndToEndTLSEncryptionEnabled       bool                                       `tfschema:"end_to_end_tls_encryption_enabled"`
}

func (r Resource) ModelObject() any {
	return &LinuxFunctionAppSlotModel{}
}
