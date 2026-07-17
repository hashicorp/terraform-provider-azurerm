// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-04-01/managedclusters"
	dnsValidate "github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/kubernetes"
	containerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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
	Tags                   map[string]interface{}                     `tfschema:"tags"`
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
type HostedSystemProfile struct {
	NodeSubnetID       string `tfschema:"node_subnet_id"`
	SystemNodeSubnetID string `tfschema:"system_node_subnet_id"`
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

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name azurerm_kubernetes_automatic_cluster -properties "name,resource_group_name"
type KubernetesAutomaticClusterResource struct{}

var (
	_ sdk.ResourceWithUpdate         = KubernetesAutomaticClusterResource{}
	_ sdk.ResourceWithIdentity       = KubernetesAutomaticClusterResource{}
	_ sdk.ResourceWithCustomImporter = KubernetesAutomaticClusterResource{}
)

func (r KubernetesAutomaticClusterResource) ResourceType() string {
	return "azurerm_kubernetes_automatic_cluster"
}

func (r KubernetesAutomaticClusterResource) ModelObject() interface{} {
	return &KubernetesAutomaticClusterModel{}
}

func (r KubernetesAutomaticClusterResource) Identity() resourceids.ResourceId {
	return &commonids.KubernetesClusterId{}
}

func (r KubernetesAutomaticClusterResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := commonids.ParseKubernetesClusterID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		client := metadata.Client.Containers.KubernetesClustersClient_v2026_04_01
		resp, err := client.Get(ctx, *id)
		if err != nil || resp.Model == nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if resp.Model.Sku == nil || resp.Model.Sku.Name == nil {
			return fmt.Errorf("importing %s: SKU information is missing", id)
		}

		if pointer.From(resp.Model.Sku.Name) != managedclusters.ManagedClusterSKUNameAutomatic {
			return fmt.Errorf("importing %s: specified Kubernetes Cluster is not using the SKU `Automatic`, got `%s`", id, pointer.From(resp.Model.Sku.Name))
		}

		return nil
	}
}

func (r KubernetesAutomaticClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateKubernetesClusterID
}

func (r KubernetesAutomaticClusterResource) CustomizeDiff() sdk.ResourceFunc {
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
				hostedSystem := make([]interface{}, 0)
				if v, ok := rd.Get("hosted_system").([]interface{}); ok {
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

			privateCluster := rd.Get("private_cluster").([]interface{})
			if len(privateCluster) > 0 && privateCluster[0] != nil {
				privateClusterConfig := privateCluster[0].(map[string]interface{})
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

func (r KubernetesAutomaticClusterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: containerValidate.KubernetesClusterName,
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
							ValidateFunc: validate.CIDR,
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
			// O+C if no subnet ids are supplied, it will return the new, managed subnet ids
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
								dnsValidate.ValidateDnsZoneID,
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

func (r KubernetesAutomaticClusterResource) Attributes() map[string]*pluginsdk.Schema {
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

func (r KubernetesAutomaticClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clusterClient := metadata.Client.Containers.KubernetesClustersClient_v2026_04_01
			subscriptionID := metadata.Client.Account.SubscriptionId

			var model KubernetesAutomaticClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := commonids.NewKubernetesClusterID(subscriptionID, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := clusterClient.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			clusterIdentity, err := identity.ExpandSystemOrUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding identity: %+v", err)
			}

			parameters := managedclusters.ManagedCluster{
				Location: location.Normalize(model.Location),
				Sku: &managedclusters.ManagedClusterSKU{
					Name: pointer.To(managedclusters.ManagedClusterSKUNameAutomatic),
					Tier: pointer.To(managedclusters.ManagedClusterSKUTierStandard),
				},
				Properties: &managedclusters.ManagedClusterProperties{
					ApiServerAccessProfile: expandKubernetesAutomaticClusterAPIAccessProfile(model),
					HostedSystemProfile:    expandKubernetesAutomaticClusterHostedSystemProfile(model.HostedSystemProfile),
					IngressProfile:         expandKubernetesAutomaticClusterWebAppRoutingIngress(model.WebAppRoutingIngress),
					ServiceMeshProfile:     expandKubernetesAutomaticClusterServiceMeshProfile(model.ServiceMeshProfile, nil),
				},
				Identity: clusterIdentity,
				Tags:     tags.Expand(model.Tags),
			}

			if err := clusterClient.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, managedclusters.DefaultCreateOrUpdateOperationOptions(), metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r KubernetesAutomaticClusterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clusterClient := metadata.Client.Containers.KubernetesClustersClient_v2026_04_01

			id, err := commonids.ParseKubernetesClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := clusterClient.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			return r.flatten(ctx, metadata, id, resp.Model, true)
		},
	}
}

func (r KubernetesAutomaticClusterResource) flatten(ctx context.Context, metadata sdk.ResourceMetaData, id *commonids.KubernetesClusterId, model *managedclusters.ManagedCluster, includeResource bool) error {
	client := metadata.Client.Containers.KubernetesClustersClient_v2026_04_01

	state := KubernetesAutomaticClusterModel{
		Name:              id.ManagedClusterName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)

		state.Tags = tags.Flatten(model.Tags)

		if props := model.Properties; props != nil {
			state.FQDN = pointer.From(props.Fqdn)
			state.PortalFQDN = pointer.From(props.AzurePortalFQDN)
			state.PrivateFQDN = pointer.From(props.PrivateFQDN)

			state.CurrentKubernetesVersion = pointer.From(props.CurrentKubernetesVersion)
			if props.OidcIssuerProfile != nil {
				state.OIDCIssuerURL = pointer.From(props.OidcIssuerProfile.IssuerURL)
			}

			if nodeResourceGroup := pointer.From(props.NodeResourceGroup); nodeResourceGroup != "" {
				state.NodeResourceGroupID = commonids.NewResourceGroupID(id.SubscriptionId, nodeResourceGroup).ID()
			}
			var err error

			state.APIServerAccessProfile, state.PrivateCluster, err = flattenKubernetesAutomaticClusterAPIAccessProfile(props.ApiServerAccessProfile)
			if err != nil {
				return fmt.Errorf("flattening API access profile: %w", err)
			}

			state.HostedSystemProfile = flattenKubernetesAutomaticClusterHostedSystemProfile(props.HostedSystemProfile)

			flattenedIdentity, err := identity.FlattenSystemOrUserAssignedMapToModel(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening identity: %w", err)
			}
			state.Identity = pointer.From(flattenedIdentity)

			state.ServiceMeshProfile = flattenKubernetesAutomaticClusterServiceMeshProfile(props.ServiceMeshProfile)

			webAppRoutingIngress, err := flattenKubernetesAutomaticClusterWebAppRoutingIngress(props.IngressProfile)
			if err != nil {
				return fmt.Errorf("flattening Web App Routing Ingress: %w", err)
			}
			state.WebAppRoutingIngress = webAppRoutingIngress
		}

		if includeResource {
			credentials, err := client.ListClusterUserCredentials(ctx, *id, managedclusters.ListClusterUserCredentialsOperationOptions{})
			if err != nil {
				return fmt.Errorf("retrieving User Credentials for %s: %+v", id, err)
			}
			if credentials.Model == nil {
				return fmt.Errorf("retrieving User Credentials for %s: payload is empty", id)
			}

			kubeConfigRaw, kubeConfig := flattenKubernetesClusterCredentialsTyped(credentials.Model, "clusterUser")
			state.KubeConfigRaw = pointer.From(kubeConfigRaw)
			state.KubeConfig = kubeConfig
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r KubernetesAutomaticClusterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clusterClient := metadata.Client.Containers.KubernetesClustersClient_v2026_04_01

			id, err := commonids.ParseKubernetesClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := clusterClient.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			if existing.Model.Identity != nil && existing.Model.Identity.IdentityIds != nil {
				for k := range existing.Model.Identity.IdentityIds {
					existing.Model.Identity.IdentityIds[k] = identity.UserAssignedIdentityDetails{}
				}
			}

			var model KubernetesAutomaticClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %w", err)
			}

			props := existing.Model.Properties

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = tags.Expand(model.Tags)
			}

			if metadata.ResourceData.HasChanges("api_server_access", "private_cluster") {
				props.ApiServerAccessProfile = expandKubernetesAutomaticClusterAPIAccessProfile(model)
			}

			if metadata.ResourceData.HasChange("identity") {
				existing.Model.Identity, err = identity.ExpandSystemOrUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding identity: %+v", err)
				}
			}

			if metadata.ResourceData.HasChange("service_mesh") {
				props.ServiceMeshProfile = expandKubernetesAutomaticClusterServiceMeshProfile(model.ServiceMeshProfile, props.ServiceMeshProfile)
			}

			if metadata.ResourceData.HasChange("web_app_routing_ingress") {
				props.IngressProfile = expandKubernetesAutomaticClusterWebAppRoutingIngress(model.WebAppRoutingIngress)
			}

			if err := clusterClient.CreateOrUpdateThenPoll(ctx, *id, *existing.Model, managedclusters.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %w", *id, err)
			}

			return nil
		},
	}
}

func (r KubernetesAutomaticClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clusterClient := metadata.Client.Containers.KubernetesClustersClient_v2026_04_01

			id, err := commonids.ParseKubernetesClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := clusterClient.DeleteThenPoll(ctx, *id, managedclusters.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %w", *id, err)
			}

			return nil
		},
	}
}

func expandKubernetesAutomaticClusterHostedSystemProfile(input []HostedSystemProfile) *managedclusters.ManagedClusterHostedSystemProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0]

	profile := &managedclusters.ManagedClusterHostedSystemProfile{
		Enabled: pointer.To(true),
	}

	if config.NodeSubnetID != "" {
		profile.NodeSubnetID = pointer.To(config.NodeSubnetID)
	}

	if config.SystemNodeSubnetID != "" {
		profile.SystemNodeSubnetID = pointer.To(config.SystemNodeSubnetID)
	}

	return profile
}

func flattenKubernetesAutomaticClusterHostedSystemProfile(profile *managedclusters.ManagedClusterHostedSystemProfile) []HostedSystemProfile {
	if profile == nil {
		return []HostedSystemProfile{}
	}

	return []HostedSystemProfile{{
		NodeSubnetID:       pointer.From(profile.NodeSubnetID),
		SystemNodeSubnetID: pointer.From(profile.SystemNodeSubnetID),
	}}
}

func expandKubernetesAutomaticClusterAPIAccessProfile(model KubernetesAutomaticClusterModel) *managedclusters.ManagedClusterAPIServerAccessProfile {
	enablePrivateCluster, enablePrivateClusterPublicFQDN, privateDNSZoneID := expandKubernetesAutomaticClusterPrivateCluster(model.PrivateCluster)

	apiAccessProfile := &managedclusters.ManagedClusterAPIServerAccessProfile{
		EnablePrivateCluster:           pointer.To(enablePrivateCluster),
		EnablePrivateClusterPublicFQDN: pointer.To(enablePrivateClusterPublicFQDN),
	}

	if privateDNSZoneID != "" {
		apiAccessProfile.PrivateDNSZone = pointer.To(privateDNSZoneID)
	}

	if len(model.APIServerAccessProfile) == 0 {
		return apiAccessProfile
	}

	config := model.APIServerAccessProfile[0]
	apiAccessProfile.AuthorizedIPRanges = pointer.To(config.AuthorizedIPRanges)

	if config.SubnetID != "" {
		apiAccessProfile.SubnetId = pointer.To(config.SubnetID)
	}

	return apiAccessProfile
}

func flattenKubernetesAutomaticClusterAPIAccessProfile(profile *managedclusters.ManagedClusterAPIServerAccessProfile) ([]APIServerAccessProfileModel, []PrivateClusterModel, error) {
	apiServerAccessProfile := make([]APIServerAccessProfileModel, 0, 1)

	if profile == nil {
		return apiServerAccessProfile, []PrivateClusterModel{}, nil
	}

	privateCluster, err := flattenKubernetesAutomaticClusterPrivateCluster(pointer.From(profile.EnablePrivateCluster), pointer.From(profile.EnablePrivateClusterPublicFQDN), pointer.From(profile.PrivateDNSZone))
	if err != nil {
		return nil, nil, err
	}

	// API access profile can be managed by other properties, only return it if one of the properties has been set
	hasAuthorizedIPRanges := profile.AuthorizedIPRanges != nil && len(*profile.AuthorizedIPRanges) > 0
	hasSubnetId := profile.SubnetId != nil && *profile.SubnetId != ""

	if !hasAuthorizedIPRanges && !hasSubnetId {
		return apiServerAccessProfile, privateCluster, nil
	}

	subnetId := ""
	if hasSubnetId {
		parsedSubnetId, err := parse.SubnetID(pointer.From(profile.SubnetId))
		if err != nil {
			return nil, nil, fmt.Errorf("parsing `api_server_access.0.subnet_id`: %+v", err)
		}
		subnetId = parsedSubnetId.ID()
	}

	apiServerAccessProfile = append(apiServerAccessProfile, APIServerAccessProfileModel{
		AuthorizedIPRanges: pointer.From(profile.AuthorizedIPRanges),
		SubnetID:           subnetId,
	})

	return apiServerAccessProfile, privateCluster, nil
}

func expandKubernetesAutomaticClusterPrivateCluster(model []PrivateClusterModel) (bool, bool, string) {
	if len(model) == 0 {
		return false, false, ""
	}

	config := model[0]
	return true, config.PrivateClusterPublicFQDNEnabled, config.PrivateDNSZoneID
}

func flattenKubernetesAutomaticClusterPrivateCluster(enablePrivateCluster bool, enablePrivateClusterPublicFQDN bool, privateDNSZoneID string) ([]PrivateClusterModel, error) {
	if !enablePrivateCluster {
		return []PrivateClusterModel{}, nil
	}

	if privateDNSZoneID != "None" && privateDNSZoneID != "System" && privateDNSZoneID != "" {
		parsedPrivateDNSZoneID, err := privatezones.ParsePrivateDnsZoneIDInsensitively(privateDNSZoneID)
		if err != nil {
			return nil, fmt.Errorf("parsing `private_cluster.0.private_dns_zone_id`: %+v", err)
		}
		privateDNSZoneID = parsedPrivateDNSZoneID.ID()
	}

	return []PrivateClusterModel{
		{
			PrivateClusterPublicFQDNEnabled: enablePrivateClusterPublicFQDN,
			PrivateDNSZoneID:                privateDNSZoneID,
		},
	}, nil
}

func expandKubernetesAutomaticClusterWebAppRoutingIngress(input []WebAppRoutingIngressModel) *managedclusters.ManagedClusterIngressProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0]

	istioEnabled := managedclusters.GatewayAPIIstioEnabledDisabled
	if config.IstioEnabled {
		istioEnabled = managedclusters.GatewayAPIIstioEnabledEnabled
	}

	nginxControllerType := managedclusters.NginxIngressControllerTypeNone
	if config.DefaultNginxController != "" {
		nginxControllerType = managedclusters.NginxIngressControllerType(config.DefaultNginxController)
	}

	ingress := managedclusters.ManagedClusterIngressProfile{
		WebAppRouting: &managedclusters.ManagedClusterIngressProfileWebAppRouting{
			Enabled: pointer.To(true),
			Nginx: &managedclusters.ManagedClusterIngressProfileNginx{
				DefaultIngressControllerType: pointer.To(nginxControllerType),
			},
			GatewayAPIImplementations: &managedclusters.ManagedClusterWebAppRoutingGatewayAPIImplementations{
				AppRoutingIstio: &managedclusters.ManagedClusterAppRoutingIstio{
					Mode: &istioEnabled,
				},
			},
		},
	}

	if len(config.DNSZoneIDs) > 0 {
		ingress.WebAppRouting.DnsZoneResourceIds = pointer.To(config.DNSZoneIDs)
	}

	return &ingress
}

func flattenKubernetesAutomaticClusterWebAppRoutingIngress(input *managedclusters.ManagedClusterIngressProfile) ([]WebAppRoutingIngressModel, error) {
	if input == nil || input.WebAppRouting == nil || input.WebAppRouting.Enabled == nil || !*input.WebAppRouting.Enabled {
		return []WebAppRoutingIngressModel{}, nil
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

	webAppRoutingIdentity := make([]WebAppRoutingIdentityModel, 0)
	if input.WebAppRouting.Identity != nil {
		parsedResourceId, err := commonids.ParseUserAssignedIdentityIDInsensitively(pointer.From(input.WebAppRouting.Identity.ResourceId))
		if err != nil {
			return nil, fmt.Errorf("parsing `web_app_routing_ingress.0.web_app_routing_identity.0.user_assigned_identity_id`: %+v", err)
		}

		webAppRoutingIdentity = append(webAppRoutingIdentity, WebAppRoutingIdentityModel{
			ClientID:               pointer.From(input.WebAppRouting.Identity.ClientId),
			ObjectID:               pointer.From(input.WebAppRouting.Identity.ObjectId),
			UserAssignedIdentityID: parsedResourceId.ID(),
		})
	}

	return []WebAppRoutingIngressModel{{
		DNSZoneIDs:             dnsZoneIDs,
		IstioEnabled:           istioEnabled,
		DefaultNginxController: defaultNginxController,
		WebAppRoutingIdentity:  webAppRoutingIdentity,
	}}, nil
}

func expandKubernetesAutomaticClusterServiceMeshProfile(input []ServiceMeshProfileModel, existing *managedclusters.ServiceMeshProfile) *managedclusters.ServiceMeshProfile {
	if len(input) == 0 {
		// explicitly disable istio if it was enabled before
		if existing != nil && existing.Mode == managedclusters.ServiceMeshModeIstio {
			return &managedclusters.ServiceMeshProfile{
				Mode: managedclusters.ServiceMeshModeDisabled,
			}
		}
		return nil
	}

	config := input[0]

	profile := managedclusters.ServiceMeshProfile{
		Mode: managedclusters.ServiceMeshModeIstio,
		Istio: &managedclusters.IstioServiceMesh{
			Components: &managedclusters.IstioComponents{
				IngressGateways: &[]managedclusters.IstioIngressGateway{
					{
						Enabled: config.InternalIngressGatewayEnabled,
						Mode:    managedclusters.IstioIngressGatewayModeInternal,
					}, {
						Enabled: config.ExternalIngressGatewayEnabled,
						Mode:    managedclusters.IstioIngressGatewayModeExternal,
					},
				},
				ProxyRedirectionMechanism: pointer.To(managedclusters.ProxyRedirectionMechanism(config.ProxyRedirectMechanism)),
			},
		},
	}

	if len(config.CertificateAuthority) > 0 {
		profile.Istio.CertificateAuthority = expandKubernetesAutomaticClusterServiceMeshProfileCertificateAuthority(config.CertificateAuthority)
	}

	if len(config.Revisions) > 0 {
		profile.Istio.Revisions = pointer.To(config.Revisions)
	}

	return &profile
}

func flattenKubernetesAutomaticClusterServiceMeshProfile(profile *managedclusters.ServiceMeshProfile) []ServiceMeshProfileModel {
	if profile == nil || profile.Mode != managedclusters.ServiceMeshModeIstio || profile.Istio == nil {
		return []ServiceMeshProfileModel{}
	}

	revisions := make([]string, 0)
	if profile.Istio.Revisions != nil {
		revisions = pointer.From(profile.Istio.Revisions)
	}

	internalIngressGatewayEnabled := false
	externalIngressGatewayEnabled := false
	proxyRedirectMechanism := ""

	if profile.Istio.Components != nil {
		if profile.Istio.Components.IngressGateways != nil {
			for _, gateway := range pointer.From(profile.Istio.Components.IngressGateways) {
				if gateway.Mode == managedclusters.IstioIngressGatewayModeInternal {
					internalIngressGatewayEnabled = gateway.Enabled
				}
				if gateway.Mode == managedclusters.IstioIngressGatewayModeExternal {
					externalIngressGatewayEnabled = gateway.Enabled
				}
			}
		}
		proxyRedirectMechanism = string(pointer.From(profile.Istio.Components.ProxyRedirectionMechanism))
	}

	certificateAuthority := flattenKubernetesAutomaticClusterServiceMeshProfileCertificateAuthority(profile.Istio.CertificateAuthority)

	return []ServiceMeshProfileModel{{
		Revisions:                     revisions,
		InternalIngressGatewayEnabled: internalIngressGatewayEnabled,
		ExternalIngressGatewayEnabled: externalIngressGatewayEnabled,
		ProxyRedirectMechanism:        proxyRedirectMechanism,
		CertificateAuthority:          certificateAuthority,
	}}
}

func expandKubernetesAutomaticClusterServiceMeshProfileCertificateAuthority(input []CertificateAuthorityModel) *managedclusters.IstioCertificateAuthority {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	return &managedclusters.IstioCertificateAuthority{
		Plugin: &managedclusters.IstioPluginCertificateAuthority{
			KeyVaultId:          pointer.To(config.KeyVaultID),
			RootCertObjectName:  pointer.To(config.RootCertObjectName),
			CertChainObjectName: pointer.To(config.CertChainObjectName),
			CertObjectName:      pointer.To(config.CertObjectName),
			KeyObjectName:       pointer.To(config.KeyObjectName),
		},
	}
}

func flattenKubernetesAutomaticClusterServiceMeshProfileCertificateAuthority(certificateAuthority *managedclusters.IstioCertificateAuthority) []CertificateAuthorityModel {
	if certificateAuthority == nil || certificateAuthority.Plugin == nil {
		return []CertificateAuthorityModel{}
	}

	plugin := certificateAuthority.Plugin

	return []CertificateAuthorityModel{{
		KeyVaultID:          pointer.From(plugin.KeyVaultId),
		RootCertObjectName:  pointer.From(plugin.RootCertObjectName),
		CertChainObjectName: pointer.From(plugin.CertChainObjectName),
		CertObjectName:      pointer.From(plugin.CertObjectName),
		KeyObjectName:       pointer.From(plugin.KeyObjectName),
	}}
}

func flattenKubernetesClusterCredentialsTyped(model *managedclusters.CredentialResults, configName string) (*string, []KubeConfigModel) {
	if model == nil || model.Kubeconfigs == nil || len(*model.Kubeconfigs) < 1 {
		return nil, []KubeConfigModel{}
	}

	for _, c := range *model.Kubeconfigs {
		if c.Name == nil || *c.Name != configName {
			continue
		}
		if kubeConfigRaw := c.Value; kubeConfigRaw != nil {
			rawConfig := *kubeConfigRaw
			if base64IsEncoded(*kubeConfigRaw) {
				rawConfig = base64Decode(*kubeConfigRaw)
			}

			var flattenedKubeConfig []KubeConfigModel

			if strings.Contains(rawConfig, "apiserver-id:") || strings.Contains(rawConfig, "exec") {
				kubeConfigAAD, err := kubernetes.ParseKubeConfigAAD(rawConfig)
				if err != nil {
					return pointer.To(rawConfig), []KubeConfigModel{}
				}

				flattenedKubeConfig = flattenKubernetesAutomaticClusterKubeConfigAAD(*kubeConfigAAD)
			} else {
				kubeConfig, err := kubernetes.ParseKubeConfig(rawConfig)
				if err != nil {
					return pointer.To(rawConfig), []KubeConfigModel{}
				}

				flattenedKubeConfig = flattenKubernetesAutomaticClusterKubeConfig(*kubeConfig)
			}

			return pointer.To(rawConfig), flattenedKubeConfig
		}
	}

	return nil, []KubeConfigModel{}
}

func flattenKubernetesAutomaticClusterKubeConfig(config kubernetes.KubeConfig) []KubeConfigModel {
	if len(config.Clusters) == 0 || len(config.Users) == 0 {
		return []KubeConfigModel{}
	}
	cluster := config.Clusters[0].Cluster
	user := config.Users[0].User
	name := config.Users[0].Name

	return []KubeConfigModel{
		{
			Host:                 cluster.Server,
			Username:             name,
			Password:             user.Token,
			ClientCertificate:    user.ClientCertificteData,
			ClientKey:            user.ClientKeyData,
			ClusterCACertificate: cluster.ClusterAuthorityData,
		},
	}
}

func flattenKubernetesAutomaticClusterKubeConfigAAD(config kubernetes.KubeConfigAAD) []KubeConfigModel {
	if len(config.Clusters) == 0 || len(config.Users) == 0 {
		return []KubeConfigModel{}
	}
	cluster := config.Clusters[0].Cluster
	name := config.Users[0].Name

	return []KubeConfigModel{
		{
			Host:                 cluster.Server,
			Username:             name,
			ClusterCACertificate: cluster.ClusterAuthorityData,
		},
	}
}
