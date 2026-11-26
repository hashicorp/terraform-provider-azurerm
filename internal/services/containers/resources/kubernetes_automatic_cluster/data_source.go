// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_automatic_cluster

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/kubernetes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = DataSource{}

type DataSource struct{}

type KubernetesAutomaticClusterDataSourceModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	Location          string `tfschema:"location"`

	APIServerAccess          []APIServerAccessDataSourceModel           `tfschema:"api_server_access"`
	CurrentKubernetesVersion string                                     `tfschema:"current_kubernetes_version"`
	DNSPrefix                string                                     `tfschema:"dns_prefix"`
	FQDN                     string                                     `tfschema:"fully_qualified_domain_name"`
	HostedSystemProfile      []HostedSystemProfileDataSourceModel       `tfschema:"hosted_system"`
	PortalFQDN               string                                     `tfschema:"portal_fully_qualified_domain_name"`
	PrivateCluster           []PrivateClusterDataSourceModel            `tfschema:"private_cluster"`
	PrivateFQDN              string                                     `tfschema:"private_fully_qualified_domain_name"`
	Identity                 []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	KubernetesVersion        string                                     `tfschema:"kubernetes_version"`
	KubeConfig               []KubeConfigModel                          `tfschema:"kube_config"`
	KubeConfigRaw            string                                     `tfschema:"kube_config_raw"`
	KubeletIdentity          []KubeletIdentityDataSourceModel           `tfschema:"kubelet_identity"`
	NodeResourceGroup        string                                     `tfschema:"node_resource_group"`
	NodeResourceGroupID      string                                     `tfschema:"node_resource_group_id"`
	ServiceMeshProfile       []ServiceMeshProfileDataSourceModel        `tfschema:"service_mesh"`
	WebAppRoutingIngress     []WebAppRoutingIngressDataSourceModel      `tfschema:"web_app_routing_ingress"`
	Tags                     map[string]any                             `tfschema:"tags"`
}

type APIServerAccessDataSourceModel struct {
	AuthorizedIPRanges []string `tfschema:"authorized_ip_ranges"`
	SubnetID           string   `tfschema:"subnet_id"`
}

type HostedSystemProfileDataSourceModel struct {
	NodeSubnetID       string `tfschema:"node_subnet_id"`
	SystemNodeSubnetID string `tfschema:"system_node_subnet_id"`
}

type WebAppRoutingIngressDataSourceModel struct {
	DNSZoneIDs             []string                               `tfschema:"dns_zone_ids"`
	DefaultNginxController string                                 `tfschema:"default_nginx_controller"`
	IstioEnabled           bool                                   `tfschema:"istio_enabled"`
	WebAppRoutingIdentity  []WebAppRoutingIdentityDataSourceModel `tfschema:"web_app_routing_identity"`
}

type WebAppRoutingIdentityDataSourceModel struct {
	ClientID               string `tfschema:"client_id"`
	ObjectID               string `tfschema:"object_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

type KubeletIdentityDataSourceModel struct {
	ClientID               string `tfschema:"client_id"`
	ObjectID               string `tfschema:"object_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

type PrivateClusterDataSourceModel struct {
	PrivateClusterPublicFQDNEnabled bool   `tfschema:"public_fully_qualified_domain_name_enabled"`
	PrivateDNSZoneID                string `tfschema:"private_dns_zone_id"`
}

type ServiceMeshProfileDataSourceModel struct {
	Revisions                     []string                              `tfschema:"revisions"`
	InternalIngressGatewayEnabled bool                                  `tfschema:"internal_ingress_gateway_enabled"`
	ExternalIngressGatewayEnabled bool                                  `tfschema:"external_ingress_gateway_enabled"`
	ProxyRedirectMechanism        string                                `tfschema:"proxy_redirect_mechanism"`
	CertificateAuthority          []CertificateAuthorityDataSourceModel `tfschema:"certificate_authority"`
}

type CertificateAuthorityDataSourceModel struct {
	KeyVaultID          string `tfschema:"key_vault_id"`
	RootCertObjectName  string `tfschema:"root_certificate_object_name"`
	CertChainObjectName string `tfschema:"certificate_chain_object_name"`
	CertObjectName      string `tfschema:"certificate_object_name"`
	KeyObjectName       string `tfschema:"key_object_name"`
}

func (DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Managed Kubernetes Automatic Cluster.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"api_server_access": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"authorized_ip_ranges": {
						Description: "A list of IP ranges authorised to access the API server.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"subnet_id": {
						Description: "The ID of the subnet that the API server is accessible from.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"current_kubernetes_version": {
			Description: "Contains the current version of Kubernetes running on the Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"dns_prefix": {
			Description: "The DNS Prefix of the Managed Kubernetes Automatic Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"fully_qualified_domain_name": {
			Description: "The FQDN of the Managed Kubernetes Automatic Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"hosted_system": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"node_subnet_id": {
						Description: "The ID of the subnet used for the hosted system nodes.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"system_node_subnet_id": {
						Description: "The ID of the subnet used for the system nodes.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

		"kube_config": {
			Type:      pluginsdk.TypeList,
			Computed:  true,
			Sensitive: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"host": {
						Description: "The Kubernetes cluster server host.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"username": {
						Description: "A username used to authenticate to the Kubernetes cluster.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"password": {
						Description: "A password or token used to authenticate to the Kubernetes cluster.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"client_certificate": {
						Description: "Base64 encoded public certificate used by clients to authenticate to the Kubernetes cluster.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"client_key": {
						Description: "Base64 encoded private key used by clients to authenticate to the Kubernetes cluster.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"cluster_ca_certificate": {
						Description: "Base64 encoded public CA certificate used as the root of trust for the Kubernetes cluster.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
		},

		"kube_config_raw": {
			Description: "Base64 encoded Kubernetes configuration.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"kubelet_identity": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Description: "The Client ID of the user-defined Managed Identity assigned to the Kubelets.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"object_id": {
						Description: "The Object ID of the user-defined Managed Identity assigned to the Kubelets.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"user_assigned_identity_id": {
						Description: "The ID of the User Assigned Identity assigned to the Kubelets.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"kubernetes_version": {
			Description: "The version of Kubernetes used on the Managed Kubernetes Automatic Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"node_resource_group": {
			Description: "Auto-generated Resource Group containing AKS Cluster resources.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"node_resource_group_id": {
			Description: "The ID of the Resource Group containing the resources for this Managed Kubernetes Automatic Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"portal_fully_qualified_domain_name": {
			Description: "The FQDN used by the Azure Portal for this Managed Kubernetes Automatic Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"private_fully_qualified_domain_name": {
			Description: "The FQDN of this Managed Kubernetes Automatic Cluster when Private Link has been enabled. This name is only resolvable inside the Virtual Network where the Azure Kubernetes Service is located.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"private_cluster": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"public_fully_qualified_domain_name_enabled": {
						Description: "If the public FQDN for this Managed Kubernetes Automatic Cluster is enabled.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},
					"private_dns_zone_id": {
						Description: "The ID of the Private DNS Zone used by this Managed Kubernetes Automatic Cluster.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"service_mesh": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"internal_ingress_gateway_enabled": {
						Description: "If the Istio Internal Ingress Gateway is enabled.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},
					"external_ingress_gateway_enabled": {
						Description: "If the Istio External Ingress Gateway is enabled.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},
					"certificate_authority": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key_vault_id": {
									Description: "The resource ID of the Key Vault.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"root_certificate_object_name": {
									Description: "The root certificate object name in Azure Key Vault.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"certificate_chain_object_name": {
									Description: "The certificate chain object name in Azure Key Vault.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"certificate_object_name": {
									Description: "The intermediate certificate object name in Azure Key Vault.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"key_object_name": {
									Description: "The intermediate certificate private key object name in Azure Key Vault.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},
					"proxy_redirect_mechanism": {
						Description: "The proxy redirect mechanism configured for the Istio service mesh.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"revisions": {
						Description: "List of revisions of the Istio control plane. When an upgrade is not in progress, this holds one value. When a canary upgrade is in progress, this can hold two consecutive values. [Learn More](https://learn.microsoft.com/azure/aks/istio-upgrade).",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"tags": commonschema.TagsDataSource(),

		"web_app_routing_ingress": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dns_zone_ids": {
						Description: "A list of DNS Zone IDs associated with the Web App routing ingress.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"default_nginx_controller": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"istio_enabled": {
						Description: "If Istio is enabled for the Web App routing ingress.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},
					"web_app_routing_identity": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Description: "The Client ID of the user-defined Managed Identity used by the Web App Routing ingress controller.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"object_id": {
									Description: "The Object ID of the user-defined Managed Identity used by the Web App Routing ingress controller.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"user_assigned_identity_id": {
									Description: "The ID of the User Assigned Identity used by the Web App Routing ingress controller.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (DataSource) ModelObject() any {
	return &KubernetesAutomaticClusterDataSourceModel{}
}

func (DataSource) ResourceType() string {
	return "azurerm_kubernetes_automatic_cluster"
}

func (DataSource) Description() string {
	return "Use this data source to access information about an existing Managed Kubernetes Automatic Cluster (AKS)."
}

func (DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.KubernetesClustersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state KubernetesAutomaticClusterDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %w", err)
			}

			id := commonids.NewKubernetesClusterID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %w", id, err)
			}

			userCredentialsResp, err := client.ListClusterUserCredentials(ctx, id, managedclusters.ListClusterUserCredentialsOperationOptions{})
			if err != nil && !response.WasStatusCode(userCredentialsResp.HttpResponse, http.StatusForbidden) {
				return fmt.Errorf("retrieving User Credentials for %s: %w", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				if resp.Model.Sku == nil || resp.Model.Sku.Name == nil {
					return fmt.Errorf("importing %s: SKU information is missing", id)
				}

				if pointer.From(resp.Model.Sku.Name) != managedclusters.ManagedClusterSKUNameAutomatic {
					return fmt.Errorf("importing %s: specified Kubernetes Cluster is not using the SKU `Automatic`, got `%s`", id, pointer.From(resp.Model.Sku.Name))
				}

				state.Location = location.Normalize(model.Location)
				state.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					state.DNSPrefix = pointer.From(props.DnsPrefix)
					state.FQDN = pointer.From(props.Fqdn)
					state.PortalFQDN = pointer.From(props.AzurePortalFQDN)
					state.PrivateFQDN = pointer.From(props.PrivateFQDN)
					state.KubernetesVersion = pointer.From(props.KubernetesVersion)
					state.CurrentKubernetesVersion = pointer.From(props.CurrentKubernetesVersion)
					state.NodeResourceGroup = pointer.From(props.NodeResourceGroup)
					if nodeResourceGroup := pointer.From(props.NodeResourceGroup); nodeResourceGroup != "" {
						state.NodeResourceGroupID = commonids.NewResourceGroupID(id.SubscriptionId, nodeResourceGroup).ID()
					}
					state.PrivateCluster, err = flattenKubernetesAutomaticClusterDataSourcePrivateCluster(props.ApiServerAccessProfile)
					if err != nil {
						return fmt.Errorf("flattening `private_cluster`: %w", err)
					}

					// Flatten API server access profile
					state.APIServerAccess = flattenKubernetesAutomaticClusterDataSourceAPIServerAccessProfile(props.ApiServerAccessProfile)

					// Flatten hosted system profile
					state.HostedSystemProfile = flattenKubernetesAutomaticClusterDataSourceHostedSystemProfile(props.HostedSystemProfile)

					state.ServiceMeshProfile = flattenKubernetesAutomaticClusterDataSourceServiceMeshProfile(props.ServiceMeshProfile)

					// Flatten web app routing ingress
					webAppRoutingIngress, err := flattenKubernetesAutomaticClusterDataSourceWebAppRoutingIngress(props.IngressProfile)
					if err != nil {
						return fmt.Errorf("flattening `web_app_routing_ingress`: %w", err)
					}
					state.WebAppRoutingIngress = webAppRoutingIngress
				}

				// Flatten identity
				flattenedIdentity, err := identity.FlattenSystemOrUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %w", err)
				}
				state.Identity = pointer.From(flattenedIdentity)

				// Flatten kubelet identity
				if model.Properties.IdentityProfile != nil {
					kubeletIdentity, err := flattenKubernetesAutomaticClusterDataSourceIdentityProfile(model.Properties.IdentityProfile)
					if err != nil {
						return fmt.Errorf("flattening `kubelet_identity`: %w", err)
					}
					state.KubeletIdentity = kubeletIdentity
				}

				// Flatten kube configs
				if userCredentialsResp.Model != nil {
					kubeConfigRaw, kubeConfig := flattenKubernetesAutomaticClusterCredentials(userCredentialsResp.Model, "clusterUser")
					state.KubeConfigRaw = pointer.From(kubeConfigRaw)
					state.KubeConfig = kubeConfig
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenKubernetesAutomaticClusterCredentials(model *managedclusters.CredentialResults, configName string) (*string, []KubeConfigModel) {
	if model == nil || model.Kubeconfigs == nil || len(*model.Kubeconfigs) < 1 {
		return nil, []KubeConfigModel{}
	}

	for _, c := range *model.Kubeconfigs {
		if c.Name == nil || *c.Name != configName {
			continue
		}
		if kubeConfigRaw := c.Value; kubeConfigRaw != nil {
			rawConfig := *kubeConfigRaw
			if helpers.Base64IsEncoded(*kubeConfigRaw) {
				rawConfig = helpers.Base64Decode(*kubeConfigRaw)
			}

			var flattenedKubeConfig []KubeConfigModel

			if strings.Contains(rawConfig, "apiserver-id:") || strings.Contains(rawConfig, "exec") {
				kubeConfigAAD, err := kubernetes.ParseKubeConfigAAD(rawConfig)
				if err != nil {
					return pointer.To(rawConfig), []KubeConfigModel{}
				}

				flattenedKubeConfig = flattenKubernetesAutomaticClusterDataSourceKubeConfigAAD(*kubeConfigAAD)
			} else {
				kubeConfig, err := kubernetes.ParseKubeConfig(rawConfig)
				if err != nil {
					return pointer.To(rawConfig), []KubeConfigModel{}
				}

				flattenedKubeConfig = flattenKubernetesAutomaticClusterDataSourceKubeConfig(*kubeConfig)
			}

			return pointer.To(rawConfig), flattenedKubeConfig
		}
	}

	return nil, []KubeConfigModel{}
}

func flattenKubernetesAutomaticClusterDataSourceIdentityProfile(profile *map[string]managedclusters.UserAssignedIdentity) ([]KubeletIdentityDataSourceModel, error) {
	if profile == nil || *profile == nil {
		return []KubeletIdentityDataSourceModel{}, nil
	}

	kubeletidentity := (*profile)["kubeletidentity"]

	clientId := pointer.From(kubeletidentity.ClientId)
	objectId := pointer.From(kubeletidentity.ObjectId)

	userAssignedIdentityId := ""
	if resourceid := kubeletidentity.ResourceId; resourceid != nil {
		parsedId, err := commonids.ParseUserAssignedIdentityIDInsensitively(*resourceid)
		if err != nil {
			return nil, err
		}
		userAssignedIdentityId = parsedId.ID()
	}

	return []KubeletIdentityDataSourceModel{{
		ClientID:               clientId,
		ObjectID:               objectId,
		UserAssignedIdentityID: userAssignedIdentityId,
	}}, nil
}

func flattenKubernetesAutomaticClusterDataSourceKubeConfig(config kubernetes.KubeConfig) []KubeConfigModel {
	cluster := config.Clusters[0].Cluster
	user := config.Users[0].User
	name := config.Users[0].Name

	return []KubeConfigModel{{
		Host:                 cluster.Server,
		Username:             name,
		Password:             user.Token,
		ClientCertificate:    user.ClientCertificateData,
		ClientKey:            user.ClientKeyData,
		ClusterCACertificate: cluster.ClusterAuthorityData,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceKubeConfigAAD(config kubernetes.KubeConfigAAD) []KubeConfigModel {
	cluster := config.Clusters[0].Cluster
	name := config.Users[0].Name

	return []KubeConfigModel{{
		Host:                 cluster.Server,
		Username:             name,
		ClusterCACertificate: cluster.ClusterAuthorityData,
	}}
}

func flattenKubernetesAutomaticClusterDataSourcePrivateCluster(profile *managedclusters.ManagedClusterAPIServerAccessProfile) ([]PrivateClusterDataSourceModel, error) {
	if profile == nil {
		return []PrivateClusterDataSourceModel{}, nil
	}

	enablePrivateCluster := pointer.From(profile.EnablePrivateCluster)
	privateDNSZoneID := pointer.From(profile.PrivateDNSZone)

	if !enablePrivateCluster {
		return []PrivateClusterDataSourceModel{}, nil
	}

	switch {
	case strings.EqualFold(privateDNSZoneID, "System"):
		privateDNSZoneID = "System"
	case strings.EqualFold(privateDNSZoneID, "None"):
		privateDNSZoneID = "None"
	case privateDNSZoneID != "":
		parsedPrivateDNSZoneID, err := privatezones.ParsePrivateDnsZoneIDInsensitively(privateDNSZoneID)
		if err != nil {
			return nil, fmt.Errorf("parsing `private_cluster.0.private_dns_zone_id`: %+v", err)
		}
		privateDNSZoneID = parsedPrivateDNSZoneID.ID()
	}

	return []PrivateClusterDataSourceModel{
		{
			PrivateClusterPublicFQDNEnabled: pointer.From(profile.EnablePrivateClusterPublicFQDN),
			PrivateDNSZoneID:                privateDNSZoneID,
		},
	}, nil
}

func flattenKubernetesAutomaticClusterDataSourceServiceMeshProfile(profile *managedclusters.ServiceMeshProfile) []ServiceMeshProfileDataSourceModel {
	if profile == nil || profile.Mode != managedclusters.ServiceMeshModeIstio || profile.Istio == nil {
		return []ServiceMeshProfileDataSourceModel{}
	}

	revisions := pointer.From(profile.Istio.Revisions)

	internalIngressGatewayEnabled := false
	externalIngressGatewayEnabled := false
	proxyRedirectMechanism := ""

	if profile.Istio.Components != nil && profile.Istio.Components.IngressGateways != nil {
		for _, gateway := range *profile.Istio.Components.IngressGateways {
			if gateway.Mode == managedclusters.IstioIngressGatewayModeInternal {
				internalIngressGatewayEnabled = gateway.Enabled
			}
			if gateway.Mode == managedclusters.IstioIngressGatewayModeExternal {
				externalIngressGatewayEnabled = gateway.Enabled
			}
		}
		proxyRedirectMechanism = pointer.FromEnum(profile.Istio.Components.ProxyRedirectionMechanism)
	}

	certificateAuthority := flattenKubernetesAutomaticClusterDataSourceServiceMeshProfileCertificateAuthority(profile.Istio.CertificateAuthority)

	return []ServiceMeshProfileDataSourceModel{{
		InternalIngressGatewayEnabled: internalIngressGatewayEnabled,
		ExternalIngressGatewayEnabled: externalIngressGatewayEnabled,
		CertificateAuthority:          certificateAuthority,
		ProxyRedirectMechanism:        proxyRedirectMechanism,
		Revisions:                     revisions,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceServiceMeshProfileCertificateAuthority(certificateAuthority *managedclusters.IstioCertificateAuthority) []CertificateAuthorityDataSourceModel {
	if certificateAuthority == nil || certificateAuthority.Plugin == nil {
		return []CertificateAuthorityDataSourceModel{}
	}

	plugin := certificateAuthority.Plugin

	return []CertificateAuthorityDataSourceModel{{
		KeyVaultID:          pointer.From(plugin.KeyVaultId),
		RootCertObjectName:  pointer.From(plugin.RootCertObjectName),
		CertChainObjectName: pointer.From(plugin.CertChainObjectName),
		CertObjectName:      pointer.From(plugin.CertObjectName),
		KeyObjectName:       pointer.From(plugin.KeyObjectName),
	}}
}

func flattenKubernetesAutomaticClusterDataSourceAPIServerAccessProfile(profile *managedclusters.ManagedClusterAPIServerAccessProfile) []APIServerAccessDataSourceModel {
	if profile == nil {
		return []APIServerAccessDataSourceModel{}
	}

	hasAuthorizedIPRanges := profile.AuthorizedIPRanges != nil && len(*profile.AuthorizedIPRanges) > 0
	hasSubnetId := profile.SubnetId != nil && *profile.SubnetId != ""

	if !hasAuthorizedIPRanges && !hasSubnetId {
		return []APIServerAccessDataSourceModel{}
	}

	return []APIServerAccessDataSourceModel{{
		AuthorizedIPRanges: pointer.From(profile.AuthorizedIPRanges),
		SubnetID:           pointer.From(profile.SubnetId),
	}}
}

func flattenKubernetesAutomaticClusterDataSourceHostedSystemProfile(profile *managedclusters.ManagedClusterHostedSystemProfile) []HostedSystemProfileDataSourceModel {
	if profile == nil {
		return []HostedSystemProfileDataSourceModel{}
	}

	return []HostedSystemProfileDataSourceModel{{
		NodeSubnetID:       pointer.From(profile.NodeSubnetID),
		SystemNodeSubnetID: pointer.From(profile.SystemNodeSubnetID),
	}}
}

func flattenKubernetesAutomaticClusterDataSourceWebAppRoutingIngress(input *managedclusters.ManagedClusterIngressProfile) ([]WebAppRoutingIngressDataSourceModel, error) {
	if input == nil || input.WebAppRouting == nil || input.WebAppRouting.Enabled == nil || !*input.WebAppRouting.Enabled {
		return []WebAppRoutingIngressDataSourceModel{}, nil
	}

	dnsZoneIDs := make([]string, 0)
	if input.WebAppRouting.DnsZoneResourceIds != nil {
		dnsZoneIDs = pointer.From(input.WebAppRouting.DnsZoneResourceIds)
	}

	defaultNginxController := ""
	if input.WebAppRouting.Nginx != nil {
		ingressControllerType := pointer.From(input.WebAppRouting.Nginx.DefaultIngressControllerType)
		if ingressControllerType != managedclusters.NginxIngressControllerTypeNone {
			defaultNginxController = string(ingressControllerType)
		}
	}

	istioEnabled := input.WebAppRouting.GatewayAPIImplementations != nil &&
		input.WebAppRouting.GatewayAPIImplementations.AppRoutingIstio != nil &&
		pointer.From(input.WebAppRouting.GatewayAPIImplementations.AppRoutingIstio.Mode) == managedclusters.GatewayAPIIstioEnabledEnabled

	webAppRoutingIdentity := make([]WebAppRoutingIdentityDataSourceModel, 0)
	if input.WebAppRouting.Identity != nil {
		parsedResourceId, err := commonids.ParseUserAssignedIdentityIDInsensitively(pointer.From(input.WebAppRouting.Identity.ResourceId))
		if err != nil {
			return nil, fmt.Errorf("parsing `web_app_routing_ingress.0.web_app_routing_identity.0.user_assigned_identity_id`: %+v", err)
		}

		webAppRoutingIdentity = append(webAppRoutingIdentity, WebAppRoutingIdentityDataSourceModel{
			ClientID:               pointer.From(input.WebAppRouting.Identity.ClientId),
			ObjectID:               pointer.From(input.WebAppRouting.Identity.ObjectId),
			UserAssignedIdentityID: parsedResourceId.ID(),
		})
	}

	return []WebAppRoutingIngressDataSourceModel{{
		DNSZoneIDs:             dnsZoneIDs,
		IstioEnabled:           istioEnabled,
		DefaultNginxController: defaultNginxController,
		WebAppRoutingIdentity:  webAppRoutingIdentity,
	}}, nil
}
