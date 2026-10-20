// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_automatic_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type KubernetesAutomaticClusterModel struct {
	Name                   string                                     `tfschema:"name"`
	Location               string                                     `tfschema:"location"`
	ResourceGroupName      string                                     `tfschema:"resource_group_name"`
	APIServerAccessProfile []APIServerAccessProfileModel              `tfschema:"api_server_access"`
	HostedSystemProfile    []HostedSystemProfile                      `tfschema:"hosted_system"`
	Identity               []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	PrivateCluster         []PrivateClusterModel                      `tfschema:"private_cluster"`
	ServiceMeshProfile     []ServiceMeshProfileModel                  `tfschema:"service_mesh"`
	WebAppRoutingIngress   []WebAppRoutingIngressModel                `tfschema:"web_app_routing_ingress"`
	Tags                   map[string]any                             `tfschema:"tags"`
	// Computed fields
	CurrentKubernetesVersion string            `tfschema:"current_kubernetes_version"`
	FQDN                     string            `tfschema:"fully_qualified_domain_name"`
	PortalFQDN               string            `tfschema:"portal_fully_qualified_domain_name"`
	KubeConfig               []KubeConfigModel `tfschema:"kube_config"`
	KubeConfigRaw            string            `tfschema:"kube_config_raw"`
	NodeResourceGroupID      string            `tfschema:"node_resource_group_id"`
	OIDCIssuerURL            string            `tfschema:"oidc_issuer_url"`
	PrivateFQDN              string            `tfschema:"private_fully_qualified_domain_name"`
}

type APIServerAccessProfileModel struct {
	AuthorizedIPRanges []string `tfschema:"authorized_ip_ranges"`
	SubnetID           string   `tfschema:"subnet_id"`
}

type KubeConfigModel struct {
	Host                 string `tfschema:"host"`
	Username             string `tfschema:"username"`
	Password             string `tfschema:"password"`
	ClientCertificate    string `tfschema:"client_certificate"`
	ClientKey            string `tfschema:"client_key"`
	ClusterCACertificate string `tfschema:"cluster_ca_certificate"`
}

type PrivateClusterModel struct {
	PrivateClusterPublicFQDNEnabled bool   `tfschema:"public_fully_qualified_domain_name_enabled"`
	PrivateDNSZoneID                string `tfschema:"private_dns_zone_id"`
}

type ServiceMeshProfileModel struct {
	Revisions                     []string                    `tfschema:"revisions"`
	InternalIngressGatewayEnabled bool                        `tfschema:"internal_ingress_gateway_enabled"`
	ExternalIngressGatewayEnabled bool                        `tfschema:"external_ingress_gateway_enabled"`
	ProxyRedirectMechanism        string                      `tfschema:"proxy_redirect_mechanism"`
	CertificateAuthority          []CertificateAuthorityModel `tfschema:"certificate_authority"`
}

type CertificateAuthorityModel struct {
	KeyVaultID          string `tfschema:"key_vault_id"`
	RootCertObjectName  string `tfschema:"root_certificate_object_name"`
	CertObjectName      string `tfschema:"certificate_object_name"`
	CertChainObjectName string `tfschema:"certificate_chain_object_name"`
	KeyObjectName       string `tfschema:"key_object_name"`
}

type WebAppRoutingIngressModel struct {
	DNSZoneIDs             []string                     `tfschema:"dns_zone_ids"`
	DefaultNginxController string                       `tfschema:"default_nginx_controller"`
	IstioEnabled           bool                         `tfschema:"istio_enabled"`
	WebAppRoutingIdentity  []WebAppRoutingIdentityModel `tfschema:"web_app_routing_identity"`
}

type WebAppRoutingIdentityModel struct {
	ClientID               string `tfschema:"client_id"`
	ObjectID               string `tfschema:"object_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

func (r Resource) ModelObject() any {
	return &KubernetesAutomaticClusterModel{}
}
