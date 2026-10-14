// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package linux_web_app

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

type LinuxWebAppModel struct {
	Name                               string                                     `tfschema:"name"`
	ResourceGroup                      string                                     `tfschema:"resource_group_name"`
	Location                           string                                     `tfschema:"location"`
	ServicePlanId                      string                                     `tfschema:"service_plan_id"`
	AppSettings                        map[string]string                          `tfschema:"app_settings"`
	StickySettings                     []helpers.StickySettings                   `tfschema:"sticky_settings"`
	AuthSettings                       []helpers.AuthSettings                     `tfschema:"auth_settings"`
	AuthV2Settings                     []helpers.AuthV2Settings                   `tfschema:"auth_settings_v2"`
	Backup                             []helpers.Backup                           `tfschema:"backup"`
	ClientAffinityEnabled              bool                                       `tfschema:"client_affinity_enabled"`
	ClientCertEnabled                  bool                                       `tfschema:"client_certificate_enabled"`
	ClientCertMode                     string                                     `tfschema:"client_certificate_mode"`
	ClientCertExclusionPaths           string                                     `tfschema:"client_certificate_exclusion_paths"`
	Enabled                            bool                                       `tfschema:"enabled"`
	HttpsOnly                          bool                                       `tfschema:"https_only"`
	VirtualNetworkBackupRestoreEnabled bool                                       `tfschema:"virtual_network_backup_restore_enabled"`
	VirtualNetworkSubnetID             string                                     `tfschema:"virtual_network_subnet_id"`
	KeyVaultReferenceIdentityID        string                                     `tfschema:"key_vault_reference_identity_id"`
	LogsConfig                         []helpers.LogsConfig                       `tfschema:"logs"`
	SiteConfig                         []helpers.SiteConfigLinux                  `tfschema:"site_config"`
	StorageAccounts                    []helpers.StorageAccount                   `tfschema:"storage_account"`
	ConnectionStrings                  []helpers.ConnectionString                 `tfschema:"connection_string"`
	ZipDeployFile                      string                                     `tfschema:"zip_deploy_file"`
	Tags                               map[string]string                          `tfschema:"tags"`
	CustomDomainVerificationId         string                                     `tfschema:"custom_domain_verification_id"`
	HostingEnvId                       string                                     `tfschema:"hosting_environment_id"`
	DefaultHostname                    string                                     `tfschema:"default_hostname"`
	Kind                               string                                     `tfschema:"kind"`
	Identity                           []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	OutboundIPAddresses                string                                     `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList              []string                                   `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses        string                                     `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList      []string                                   `tfschema:"possible_outbound_ip_address_list"`
	PublicNetworkAccess                bool                                       `tfschema:"public_network_access_enabled"`
	PublishingDeployBasicAuthEnabled   bool                                       `tfschema:"webdeploy_publish_basic_authentication_enabled"`
	PublishingFTPBasicAuthEnabled      bool                                       `tfschema:"ftp_publish_basic_authentication_enabled"`
	SiteCredentials                    []helpers.SiteCredential                   `tfschema:"site_credential"`
	EndToEndTLSEncryptionEnabled       bool                                       `tfschema:"end_to_end_tls_encryption_enabled"`
	VnetImagePullEnabled               bool                                       `tfschema:"vnet_image_pull_enabled"`
}

func (r Resource) ModelObject() any {
	return &LinuxWebAppModel{}
}
