// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_automatic_cluster

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			if rd.Id() != "" && rd.HasChange("api_server_access.0.subnet_id") {
				oldValue, newValue := rd.GetChange("api_server_access.0.subnet_id")
				if oldValue.(string) != "" && newValue.(string) == "" {
					if err := rd.ForceNew("api_server_access.0.subnet_id"); err != nil {
						return err
					}
				}
			}

			if rd.HasChange("identity.0.type") {
				if err := rd.ForceNew("identity.0.type"); err != nil {
					return err
				}
			}

			identityType := rd.Get("identity.0.type").(string)
			apiServerSubnetIDIsSet := false

			if rawAPIServerAccess := rd.GetRawConfig().AsValueMap()["api_server_access"]; !rawAPIServerAccess.IsNull() {
				if rawAPIServerAccess.IsKnown() && len(rawAPIServerAccess.AsValueSlice()) > 0 {
					rawAPIServerAccessConfig := rawAPIServerAccess.AsValueSlice()[0]
					if !rawAPIServerAccessConfig.IsNull() {
						rawSubnetID := rawAPIServerAccessConfig.AsValueMap()["subnet_id"]
						apiServerSubnetIDIsSet = !rawSubnetID.IsNull()
					}
				}
			}

			if rd.Id() == "" {
				hostedSystem := make([]any, 0)
				if v, ok := rd.Get("hosted_system").([]any); ok {
					hostedSystem = v
				}
				if len(hostedSystem) == 0 {
					if !strings.EqualFold(identityType, string(identity.TypeSystemAssigned)) {
						return errors.New("when `hosted_system` is not configured, `identity.type` must be `SystemAssigned`")
					}

					if apiServerSubnetIDIsSet {
						return errors.New("when `hosted_system` is not configured, `api_server_access.subnet_id` can not be set")
					}
				}
				if len(hostedSystem) > 0 && hostedSystem[0] != nil {
					if !strings.EqualFold(identityType, string(identity.TypeUserAssigned)) {
						return errors.New("`hosted_system` requires `identity.type` to be `UserAssigned`")
					}

					if !apiServerSubnetIDIsSet {
						return errors.New("`hosted_system` requires `api_server_access.subnet_id` to be set")
					}
				}
			}

			privateCluster := rd.Get("private_cluster").([]any)
			if len(privateCluster) > 0 && privateCluster[0] != nil {
				privateClusterConfig := privateCluster[0].(map[string]any)
				privateDNSZoneID := privateClusterConfig["private_dns_zone_id"].(string)

				if privateDNSZoneID != "" {
					if privateDNSZoneID != "System" && privateDNSZoneID != "None" && !strings.EqualFold(identityType, string(identity.TypeUserAssigned)) {
						return errors.New("a user assigned identity must be used when using a custom private dns zone")
					}
				}
			}

			return nil
		},
	}
}

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Managed Kubernetes Cluster to create.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.KubernetesClusterName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"identity": commonschema.SystemOrUserAssignedIdentityRequired(),

		"api_server_access": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"authorized_ip_ranges": {
						Description:   "Set of authorized IP ranges to allow access to API server, e.g. `[\"198.51.100.0/24\"]`.",
						Type:          pluginsdk.TypeSet,
						Optional:      true,
						ConflictsWith: []string{"private_cluster"},
						AtLeastOneOf: []string{
							"api_server_access.0.authorized_ip_ranges",
							"api_server_access.0.subnet_id",
						},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsCIDRIPv4,
						},
					},

					"subnet_id": {
						Description:  "The ID of the Subnet where the API server endpoint is delegated to. Is required for bring your own networking.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateSubnetID,
						AtLeastOneOf: []string{
							"api_server_access.0.authorized_ip_ranges",
							"api_server_access.0.subnet_id",
						},
					},
				},
			},
		},

		"hosted_system": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// Note: O+C if no subnet ids are supplied, it will return the new, managed subnet ids
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"node_subnet_id": {
						Description:  "The ID of the Subnet where the user nodes are hosted. Is required for bring your own networking.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},

					"system_node_subnet_id": {
						Description:  "The ID of the Subnet where the system nodes are hosted. Is required for bring your own networking.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
				},
			},
		},

		"private_cluster": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"public_fully_qualified_domain_name_enabled": {
						Description: "Provisions a Public FQDN for the private cluster.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"private_dns_zone_id": {
						Description: "The ID of the Private DNS Zone which should be used for this Kubernetes Cluster.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
						Default:     "System",
						ValidateFunc: validation.Any(
							privatezones.ValidatePrivateDnsZoneID,
							validation.StringInSlice([]string{
								"System",
								"None",
							}, false),
						),
					},
				},
			},
		},

		"service_mesh": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"internal_ingress_gateway_enabled": {
						Description: "Enables Istio Internal Ingress Gateway.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"external_ingress_gateway_enabled": {
						Description: "Enables Istio External Ingress Gateway.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"proxy_redirect_mechanism": {
						Description:  "The mechanism used to redirect application traffic to the Istio sidecar proxy.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      managedclusters.ProxyRedirectionMechanismInitContainers,
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForProxyRedirectionMechanism(), false),
					},

					"certificate_authority": {
						Description: "This configuration allows you to bring your own root certificate and keys for Istio CA in the Istio-based service mesh add-on for Azure Kubernetes Service.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						MaxItems:    1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}, "Key Vault"),

								"root_certificate_object_name": {
									Description:  "The root certificate object name in Azure Key Vault.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"certificate_chain_object_name": {
									Description:  "The certificate chain object name in Azure Key Vault.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"certificate_object_name": {
									Description:  "The intermediate certificate object name in Azure Key Vault.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"key_object_name": {
									Description:  "The intermediate certificate private key object name in Azure Key Vault.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
					"revisions": {
						Description: "One to two Istio control plane revisions for managing minor upgrades using the canary upgrade process. For example, create the resource with `revisions` set to `[\"asm-1-27\"]`. To start the canary upgrade, change `revisions` to `[\"asm-1-27\", \"asm-1-28\"]`. To roll back the canary upgrade, revert to `[\"asm-1-27\"]`. To confirm the upgrade, change to `[\"asm-1-28\"]`.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						MinItems:    1,
						MaxItems:    2,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringStartsWithOneOf("asm-"),
						},
					},
				},
			},
		},

		"web_app_routing_ingress": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// NOTE: O+C - Azure provides default ingress configuration if not specified
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dns_zone_ids": {
						Description: "Resource IDs of the DNS zones to be associated with the Application Routing add-on. Public and private DNS zones can be in different Resource Groups, but all public DNS zones must be in the same Resource Group and all private DNS zones must be in the same Resource Group.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.Any(
								zones.ValidateDnsZoneID,
								privatezones.ValidatePrivateDnsZoneID,
							),
						},
					},

					"default_nginx_controller": {
						Description: "The ingress type for the default `NginxIngressController` custom resource.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ValidateFunc: validation.StringInSlice([]string{
							string(managedclusters.NginxIngressControllerTypeAnnotationControlled),
							string(managedclusters.NginxIngressControllerTypeInternal),
							string(managedclusters.NginxIngressControllerTypeExternal),
						}, false),
						AtLeastOneOf: []string{"web_app_routing_ingress.0.default_nginx_controller", "web_app_routing_ingress.0.istio_enabled"},
					},

					"istio_enabled": {
						Description:  "Enables Istio as a Gateway API implementation.",
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						Default:      false,
						AtLeastOneOf: []string{"web_app_routing_ingress.0.default_nginx_controller", "web_app_routing_ingress.0.istio_enabled"},
					},

					"web_app_routing_identity": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Description: "The Client ID of the user-defined Managed Identity used for Web App Routing.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},

								"object_id": {
									Description: "The Object ID of the user-defined Managed Identity used for Web App Routing.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},

								"user_assigned_identity_id": {
									Description: "The ID of the User Assigned Identity used for Web App Routing.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"current_kubernetes_version": {
			Description: "The current version running on the Azure Kubernetes Managed Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"fully_qualified_domain_name": {
			Description: "The FQDN of the Azure Kubernetes Managed Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

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
			Description: "Raw Kubernetes config to be used by [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) and other compatible tools.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"node_resource_group_id": {
			Description: "The ID of the Resource Group containing the resources for this Managed Kubernetes Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"oidc_issuer_url": {
			Description: "The OIDC issuer URL that is associated with the cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"portal_fully_qualified_domain_name": {
			Description: "The FQDN for the Azure Portal resources when Private Link has been enabled, which is only resolvable inside the Virtual Network used by the Kubernetes Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"private_fully_qualified_domain_name": {
			Description: "The FQDN for the Kubernetes Cluster when Private Link has been enabled, which is only resolvable inside the Virtual Network used by the Kubernetes Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
