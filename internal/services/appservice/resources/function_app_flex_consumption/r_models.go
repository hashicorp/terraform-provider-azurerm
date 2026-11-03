// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package function_app_flex_consumption

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

type FunctionAppFlexConsumptionModel struct {
	Name          string `tfschema:"name"`
	ResourceGroup string `tfschema:"resource_group_name"`
	Location      string `tfschema:"location"`
	ServicePlanId string `tfschema:"service_plan_id"`

	Enabled                          bool                       `tfschema:"enabled"`
	AppSettings                      map[string]string          `tfschema:"app_settings"`
	StickySettings                   []helpers.StickySettings   `tfschema:"sticky_settings"`
	AuthSettings                     []helpers.AuthSettings     `tfschema:"auth_settings"`
	AuthV2Settings                   []helpers.AuthV2Settings   `tfschema:"auth_settings_v2"`
	ClientCertEnabled                bool                       `tfschema:"client_certificate_enabled"`
	ClientCertMode                   string                     `tfschema:"client_certificate_mode"`
	ClientCertExclusionPaths         string                     `tfschema:"client_certificate_exclusion_paths"`
	ConnectionStrings                []helpers.ConnectionString `tfschema:"connection_string"`
	PublicNetworkAccess              bool                       `tfschema:"public_network_access_enabled"`
	HttpsOnly                        bool                       `tfschema:"https_only"`
	VirtualNetworkSubnetID           string                     `tfschema:"virtual_network_subnet_id"`
	ZipDeployFile                    string                     `tfschema:"zip_deploy_file"`
	PublishingDeployBasicAuthEnabled bool                       `tfschema:"webdeploy_publish_basic_authentication_enabled"`

	StorageContainerType          string                                         `tfschema:"storage_container_type"`
	StorageContainerEndpoint      string                                         `tfschema:"storage_container_endpoint"`
	StorageAuthType               string                                         `tfschema:"storage_authentication_type"`
	StorageAccessKey              string                                         `tfschema:"storage_access_key"`
	StorageUserAssignedIdentityID string                                         `tfschema:"storage_user_assigned_identity_id"`
	RuntimeName                   string                                         `tfschema:"runtime_name"`
	RuntimeVersion                string                                         `tfschema:"runtime_version"`
	MaximumInstanceCount          int64                                          `tfschema:"maximum_instance_count"`
	InstanceMemoryInMB            int64                                          `tfschema:"instance_memory_in_mb"`
	HttpConcurrency               int64                                          `tfschema:"http_concurrency"`
	AlwaysReady                   []FunctionAppAlwaysReady                       `tfschema:"always_ready"`
	SiteConfig                    []helpers.SiteConfigFunctionAppFlexConsumption `tfschema:"site_config"`
	Identity                      []identity.ModelSystemAssignedUserAssigned     `tfschema:"identity"`
	Tags                          map[string]string                              `tfschema:"tags"`

	CustomDomainVerificationId    string   `tfschema:"custom_domain_verification_id"`
	DefaultHostname               string   `tfschema:"default_hostname"`
	HostingEnvId                  string   `tfschema:"hosting_environment_id"`
	Kind                          string   `tfschema:"kind"`
	OutboundIPAddresses           string   `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList         []string `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses   string   `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList []string `tfschema:"possible_outbound_ip_address_list"`

	SiteCredentials []helpers.SiteCredential `tfschema:"site_credential"`
}

func (r Resource) ModelObject() any {
	return &FunctionAppFlexConsumptionModel{}
}
