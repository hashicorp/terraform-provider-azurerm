// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_automatic_cluster

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
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
						return fmt.Errorf("when `hosted_system` is not configured, `identity.type` must be `SystemAssigned`")
					}

					if apiServerSubnetIDIsSet {
						return fmt.Errorf("when `hosted_system` is not configured, `api_server_access.subnet_id` can not be set")
					}
				}
				if len(hostedSystem) > 0 && hostedSystem[0] != nil {
					if !strings.EqualFold(identityType, string(identity.TypeUserAssigned)) {
						return fmt.Errorf("`hosted_system` requires `identity.type` to be `UserAssigned`")
					}

					if !apiServerSubnetIDIsSet {
						return fmt.Errorf("`hosted_system` requires `api_server_access.subnet_id` to be set")
					}
				}
			}

			privateCluster := rd.Get("private_cluster").([]any)
			if len(privateCluster) > 0 && privateCluster[0] != nil {
				privateClusterConfig := privateCluster[0].(map[string]any)
				privateDNSZoneID := privateClusterConfig["private_dns_zone_id"].(string)

				if privateDNSZoneID != "" {
					if privateDNSZoneID != "System" && privateDNSZoneID != "None" && !strings.EqualFold(identityType, string(identity.TypeUserAssigned)) {
						return fmt.Errorf("a user assigned identity must be used when using a custom private dns zone")
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
						Type:     pluginsdk.TypeSet,
						Optional: true,
						AtLeastOneOf: []string{
							"api_server_access.0.authorized_ip_ranges",
							"api_server_access.0.subnet_id",
						},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsCIDRIPv4,
						},
						ConflictsWith: []string{"private_cluster"},
					},

					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"api_server_access.0.authorized_ip_ranges",
							"api_server_access.0.subnet_id",
						},
						ValidateFunc: commonids.ValidateSubnetID,
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},

					"system_node_subnet_id": {
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
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"private_dns_zone_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  "System",
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
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"external_ingress_gateway_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"proxy_redirect_mechanism": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      managedclusters.ProxyRedirectionMechanismInitContainers,
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForProxyRedirectionMechanism(), false),
					},

					"certificate_authority": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),

								"root_certificate_object_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"certificate_chain_object_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"certificate_object_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"key_object_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
					"revisions": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						MaxItems: 2,
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
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.Any(
								zones.ValidateDnsZoneID,
								privatezones.ValidatePrivateDnsZoneID,
							),
						},
					},

					"default_nginx_controller": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(managedclusters.NginxIngressControllerTypeAnnotationControlled),
							string(managedclusters.NginxIngressControllerTypeInternal),
							string(managedclusters.NginxIngressControllerTypeExternal),
						}, false),
						AtLeastOneOf: []string{"web_app_routing_ingress.0.default_nginx_controller", "web_app_routing_ingress.0.istio_enabled"},
					},

					"istio_enabled": {
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
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"object_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"user_assigned_identity_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
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
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"fully_qualified_domain_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kube_config": {
			Type:      pluginsdk.TypeList,
			Computed:  true,
			Sensitive: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"host": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"username": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"password": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"client_certificate": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"client_key": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"cluster_ca_certificate": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
				},
			},
		},

		"kube_config_raw": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"node_resource_group_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"oidc_issuer_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"portal_fully_qualified_domain_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"private_fully_qualified_domain_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
