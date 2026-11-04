// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_standard

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

type LogicAppResourceModel struct {
	Name                        string                                     `tfschema:"name"`
	ResourceGroupName           string                                     `tfschema:"resource_group_name"`
	Location                    string                                     `tfschema:"location"`
	AppServicePlanId            string                                     `tfschema:"app_service_plan_id"`
	AppSettings                 map[string]string                          `tfschema:"app_settings"`
	UseExtensionBundle          bool                                       `tfschema:"use_extension_bundle"`
	BundleVersion               string                                     `tfschema:"bundle_version"`
	ClientAffinityEnabled       bool                                       `tfschema:"client_affinity_enabled"`
	ClientCertificateMode       string                                     `tfschema:"client_certificate_mode"`
	Enabled                     bool                                       `tfschema:"enabled"`
	FtpPublishBasicAuthEnabled  bool                                       `tfschema:"ftp_publish_basic_authentication_enabled"`
	HTTPSOnly                   bool                                       `tfschema:"https_only"`
	Identity                    []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	KeyvaultReferenceIdentityId string                                     `tfschema:"key_vault_reference_identity_id"`
	SCMPublishBasicAuthEnabled  bool                                       `tfschema:"scm_publish_basic_authentication_enabled"`
	SiteConfig                  []helpers.LogicAppSiteConfig               `tfschema:"site_config"`
	ConnectionStrings           []helpers.ConnectionString                 `tfschema:"connection_string"`
	StorageAccountName          string                                     `tfschema:"storage_account_name"`
	StorageAccountAccessKey     string                                     `tfschema:"storage_account_access_key"`
	StorageKeyVaultSecretID     string                                     `tfschema:"storage_key_vault_secret_id"`
	PublicNetworkAccess         string                                     `tfschema:"public_network_access"`
	StorageAccountShareName     string                                     `tfschema:"storage_account_share_name"`
	Version                     string                                     `tfschema:"version"`
	VNETContentShareEnabled     bool                                       `tfschema:"vnet_content_share_enabled"`
	VirtualNetworkSubnetId      string                                     `tfschema:"virtual_network_subnet_id"`
	Tags                        map[string]string                          `tfschema:"tags"`

	CustomDomainVerificationId  string                           `tfschema:"custom_domain_verification_id"`
	DefaultHostname             string                           `tfschema:"default_hostname"`
	Kind                        string                           `tfschema:"kind"`
	OutboundIpAddresses         string                           `tfschema:"outbound_ip_addresses"`
	PossibleOutboundIpAddresses string                           `tfschema:"possible_outbound_ip_addresses"`
	SiteCredential              []helpers.SiteCredentialLogicApp `tfschema:"site_credential"`
}

func (r Resource) ModelObject() any {
	return &LogicAppResourceModel{}
}
