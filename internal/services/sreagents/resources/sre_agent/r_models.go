// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sre_agent

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type SREAgentModel struct {
	Name                string                                     `tfschema:"name"`
	ResourceGroupName   string                                     `tfschema:"resource_group_name"`
	Location            string                                     `tfschema:"location"`
	Identity            []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Action              []ActionModel                              `tfschema:"action"`
	ApplicationInsights []ApplicationInsightsModel                 `tfschema:"application_insights"`
	DefaultModel        []DefaultModel                             `tfschema:"default_model"`
	IncidentManagement  []IncidentManagementModel                  `tfschema:"incident_management"`
	KnowledgeGraph      []KnowledgeGraphModel                      `tfschema:"knowledge_graph"`
	Sandbox             []SandboxModel                             `tfschema:"sandbox"`
	SubnetId            string                                     `tfschema:"subnet_id"`
	UpgradeChannel      string                                     `tfschema:"upgrade_channel"`
	Tags                map[string]string                          `tfschema:"tags"`
	Endpoint            string                                     `tfschema:"endpoint"`
	OutboundIPAddresses []string                                   `tfschema:"outbound_ip_addresses"`
}

type ActionModel struct {
	UserAssignedIdentityId string `tfschema:"user_assigned_identity_id"`
	AccessLevel            string `tfschema:"access_level"`
	Mode                   string `tfschema:"mode"`
}

type ApplicationInsightsModel struct {
	AppId            string `tfschema:"application_id"`
	ConnectionString string `tfschema:"connection_string"`
}

type DefaultModel struct {
	Name     string `tfschema:"name"`
	Provider string `tfschema:"provider"`
}

type IncidentManagementModel struct {
	Type           string `tfschema:"type"`
	ConnectionKey  string `tfschema:"connection_key"`
	ConnectionName string `tfschema:"connection_name"`
	ConnectionURL  string `tfschema:"connection_url"`
	OnBehalfOfUser string `tfschema:"on_behalf_of_user"`
}

type KnowledgeGraphModel struct {
	UserAssignedIdentityId string   `tfschema:"user_assigned_identity_id"`
	ManagedResourceIds     []string `tfschema:"managed_resource_ids"`
}

type SandboxModel struct {
	Egress  []SandboxEgressModel  `tfschema:"egress"`
	Package []SandboxPackageModel `tfschema:"package"`
}

type SandboxEgressModel struct {
	AllowedCodeRepositories           []string                `tfschema:"allowed_code_repositories"`
	AllowedHosts                      []string                `tfschema:"allowed_hosts"`
	AllowedRegistries                 []string                `tfschema:"allowed_registries"`
	HTTPMCPServerNetworkAccessEnabled bool                    `tfschema:"http_mcp_server_network_access_enabled"`
	HTTPProxy                         []SandboxHTTPProxyModel `tfschema:"http_proxy"`
	Mode                              string                  `tfschema:"mode"`
}

type SandboxHTTPProxyModel struct {
	HTTPSURL string   `tfschema:"https_url"`
	HTTPURL  string   `tfschema:"http_url"`
	NoProxy  []string `tfschema:"no_proxy"`
}

type SandboxPackageModel struct {
	Name           string `tfschema:"name"`
	PackageManager string `tfschema:"package_manager"`
	Version        string `tfschema:"version"`
}

func (Resource) ModelObject() any {
	return &SREAgentModel{}
}
