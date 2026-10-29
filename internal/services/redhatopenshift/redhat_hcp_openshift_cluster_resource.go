// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

//go:generate go run ../../tools/generator-tests resourceidentity -test-expect-non-empty

package redhatopenshift

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2026-10-01/hcpopenshiftclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.ResourceWithCustomizeDiff = RedHatHcpOpenShiftCluster{}
	_ sdk.ResourceWithIdentity      = RedHatHcpOpenShiftCluster{}
	_ sdk.ResourceWithUpdate        = RedHatHcpOpenShiftCluster{}
)

type RedHatHcpOpenShiftCluster struct{}

type RedHatHcpOpenShiftClusterModel struct {
	Name                    string                                     `tfschema:"name"`
	ResourceGroup           string                                     `tfschema:"resource_group_name"`
	Location                string                                     `tfschema:"location"`
	Tags                    map[string]string                          `tfschema:"tags"`
	Identity                []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Version                 []HcpOpenShiftClusterVersionProfile        `tfschema:"version"`
	Dns                     []HcpOpenShiftClusterDnsProfile            `tfschema:"dns"`
	Network                 []HcpOpenShiftClusterNetworkProfile        `tfschema:"network"`
	ConsoleUrl              string                                     `tfschema:"console_url"`
	Api                     []HcpOpenShiftClusterApiProfile            `tfschema:"api"`
	Ingress                 []HcpOpenShiftClusterIngressProfile        `tfschema:"ingress"`
	Platform                []HcpOpenShiftClusterPlatformProfile       `tfschema:"platform"`
	Autoscaling             []HcpOpenShiftClusterAutoscalingProfile    `tfschema:"autoscaling"`
	Etcd                    []HcpOpenShiftClusterEtcdProfile           `tfschema:"etcd"`
	ImageDigestMirrors      []HcpOpenShiftClusterImageDigestMirror     `tfschema:"image_digest_mirrors"`
	NodeDrainTimeoutMinutes int64                                      `tfschema:"node_drain_timeout_minutes"`
	ClusterImageRegistry    []HcpOpenShiftClusterImageRegistryProfile  `tfschema:"cluster_image_registry"`
	CryptoRestrictions      string                                     `tfschema:"crypto_restrictions"`
}

type HcpOpenShiftClusterVersionProfile struct {
	Id           string `tfschema:"id"`
	ChannelGroup string `tfschema:"channel_group"`
}

type HcpOpenShiftClusterDnsProfile struct {
	BaseDomain       string `tfschema:"base_domain"`
	BaseDomainPrefix string `tfschema:"base_domain_prefix"`
}

type HcpOpenShiftClusterNetworkProfile struct {
	NetworkType string `tfschema:"network_type"`
	PodCidr     string `tfschema:"pod_cidr"`
	ServiceCidr string `tfschema:"service_cidr"`
	MachineCidr string `tfschema:"machine_cidr"`
	HostPrefix  int64  `tfschema:"host_prefix"`
}

type HcpOpenShiftClusterApiProfile struct {
	Url             string   `tfschema:"url"`
	Visibility      string   `tfschema:"visibility"`
	AuthorizedCidrs []string `tfschema:"authorized_cidrs"`
}

type HcpOpenShiftClusterIngressProfile struct {
	Type string `tfschema:"type"`
}

type HcpOpenShiftClusterPlatformProfile struct {
	ContainerRegistry       []HcpOpenShiftClusterContainerRegistryProfile `tfschema:"container_registry"`
	ManagedResourceGroup    string                                        `tfschema:"managed_resource_group"`
	SubnetId                string                                        `tfschema:"subnet_id"`
	VnetIntegrationSubnetId string                                        `tfschema:"vnet_integration_subnet_id"`
	OutboundType            string                                        `tfschema:"outbound_type"`
	NetworkSecurityGroupId  string                                        `tfschema:"network_security_group_id"`
	OperatorsAuthentication []HcpOpenShiftClusterOperatorsAuthentication  `tfschema:"operators_authentication"`
	IssuerUrl               string                                        `tfschema:"issuer_url"`
}

type HcpOpenShiftClusterContainerRegistryProfile struct {
	ManagedIdentity string `tfschema:"managed_identity"`
}

type HcpOpenShiftClusterOperatorsAuthentication struct {
	UserAssignedIdentities []HcpOpenShiftClusterUserAssignedIdentities `tfschema:"user_assigned_identities"`
}

type HcpOpenShiftClusterUserAssignedIdentities struct {
	ControlPlaneOperators  map[string]string `tfschema:"control_plane_operators"`
	DataPlaneOperators     map[string]string `tfschema:"data_plane_operators"`
	ServiceManagedIdentity string            `tfschema:"service_managed_identity"`
}

type HcpOpenShiftClusterAutoscalingProfile struct {
	MaxNodesTotal               int64 `tfschema:"max_nodes_total"`
	MaxPodGracePeriodSeconds    int64 `tfschema:"max_pod_grace_period_seconds"`
	MaxNodeProvisionTimeSeconds int64 `tfschema:"max_node_provision_time_seconds"`
	PodPriorityThreshold        int64 `tfschema:"pod_priority_threshold"`
}

type HcpOpenShiftClusterEtcdProfile struct {
	DataEncryption []HcpOpenShiftClusterEtcdDataEncryptionProfile `tfschema:"data_encryption"`
}

type HcpOpenShiftClusterEtcdDataEncryptionProfile struct {
	KeyManagementMode string                                         `tfschema:"key_management_mode"`
	CustomerManaged   []HcpOpenShiftClusterCustomerManagedEncryption `tfschema:"customer_managed"`
}

type HcpOpenShiftClusterCustomerManagedEncryption struct {
	EncryptionType string                             `tfschema:"encryption_type"`
	Kms            []HcpOpenShiftClusterKmsEncryption `tfschema:"kms"`
}

type HcpOpenShiftClusterKmsEncryption struct {
	VaultName    string                      `tfschema:"vault_name"`
	KeyVaultType string                      `tfschema:"key_vault_type"`
	Visibility   string                      `tfschema:"visibility"`
	ActiveKey    []HcpOpenShiftClusterKmsKey `tfschema:"active_key"`
}

type HcpOpenShiftClusterKmsKey struct {
	Name    string `tfschema:"name"`
	Version string `tfschema:"version"`
}

type HcpOpenShiftClusterImageDigestMirror struct {
	Source  string   `tfschema:"source"`
	Mirrors []string `tfschema:"mirrors"`
}

type HcpOpenShiftClusterImageRegistryProfile struct {
	State string `tfschema:"state"`
}

func (r RedHatHcpOpenShiftCluster) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diff := metadata.ResourceDiff
			identityTypePath := "identity.0.type"
			if diff.NewValueKnown(identityTypePath) && diff.Get(identityTypePath).(string) != string(identity.TypeUserAssigned) {
				return fmt.Errorf("`%s` must be `UserAssigned`", identityTypePath)
			}

			modePath := "etcd.0.data_encryption.0.key_management_mode"
			customerManagedPath := "etcd.0.data_encryption.0.customer_managed"

			if !diff.NewValueKnown(modePath) || !diff.NewValueKnown(customerManagedPath) {
				return nil
			}

			mode := diff.Get(modePath).(string)
			customerManaged := diff.Get(customerManagedPath).([]any)
			isCustomerManaged := mode == string(hcpopenshiftclusters.EtcdDataEncryptionKeyManagementModeTypeCustomerManaged)

			if isCustomerManaged && len(customerManaged) == 0 {
				return fmt.Errorf("`%s` is required when `%s` is `CustomerManaged`", customerManagedPath, modePath)
			}
			if !isCustomerManaged && len(customerManaged) > 0 {
				return fmt.Errorf("`%s` must not be set unless `%s` is `CustomerManaged`", customerManagedPath, modePath)
			}

			encryptionTypePath := customerManagedPath + ".0.encryption_type"
			kmsPath := customerManagedPath + ".0.kms"
			if len(customerManaged) > 0 && diff.NewValueKnown(encryptionTypePath) && diff.NewValueKnown(kmsPath) {
				encryptionType := diff.Get(encryptionTypePath).(string)
				kms := diff.Get(kmsPath).([]any)
				if encryptionType == string(hcpopenshiftclusters.CustomerManagedEncryptionTypeKMS) && len(kms) == 0 {
					return fmt.Errorf("`%s` is required when `%s` is `KMS`", kmsPath, encryptionTypePath)
				}
			}

			return nil
		},
	}
}

func (r RedHatHcpOpenShiftCluster) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z]([-a-zA-Z0-9]{0,52}[a-zA-Z0-9])?$`),
				"must begin with a letter, end with a letter or number, contain only letters, numbers and hyphens, and be between 1 and 54 characters in length",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityRequired(),

		/*
			Once platform managed encryption is implemented, it will be the default, and the etcd block will no longer be required.
			The schema will need to change in the following ways:

			- Make the `etcd` block optional.
			- Make `key_management_mode` optional, and make it default to platform managed.

			Note that the implementation plan described above is based on the assumption that platform managed encryption requires
			no other configuration aside from simply enabling it by setting `key_management_mode` to platform managed.
		*/
		"etcd": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"data_encryption": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key_management_mode": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringInSlice(hcpopenshiftclusters.PossibleValuesForEtcdDataEncryptionKeyManagementModeType(), false),
								},
								"customer_managed": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"encryption_type": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ForceNew:     true,
												Default:      string(hcpopenshiftclusters.CustomerManagedEncryptionTypeKMS),
												ValidateFunc: validation.StringInSlice(hcpopenshiftclusters.PossibleValuesForCustomerManagedEncryptionType(), false),
											},
											"kms": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"key_vault_type": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringInSlice(hcpopenshiftclusters.PossibleValuesForKmsKeyVaultType(), false),
															Default:      string(hcpopenshiftclusters.KmsKeyVaultTypeKeyVault),
														},
														"vault_name": {
															Type:         pluginsdk.TypeString,
															Required:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringLenBetween(1, 255),
														},
														"visibility": {
															Type:         pluginsdk.TypeString,
															Required:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringInSlice(hcpopenshiftclusters.PossibleValuesForKeyVaultVisibility(), false),
														},
														"active_key": {
															Type:     pluginsdk.TypeList,
															Required: true,
															MaxItems: 1,
															Elem: &pluginsdk.Resource{
																Schema: map[string]*pluginsdk.Schema{
																	"name": {
																		Type:         pluginsdk.TypeString,
																		Required:     true,
																		ForceNew:     true,
																		ValidateFunc: validation.StringLenBetween(1, 255),
																	},
																	"version": {
																		Type:         pluginsdk.TypeString,
																		Required:     true,
																		ValidateFunc: validation.StringLenBetween(1, 255),
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},

		"platform": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"container_registry": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"managed_identity": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: commonids.ValidateUserAssignedIdentityID,
								},
							},
						},
					},
					"managed_resource_group": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
					"vnet_integration_subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
					"outbound_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Default:      string(hcpopenshiftclusters.OutboundTypeLoadBalancer),
						ValidateFunc: validation.StringInSlice(hcpopenshiftclusters.PossibleValuesForOutboundType(), false),
					},
					"network_security_group_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: networksecuritygroups.ValidateNetworkSecurityGroupID,
					},
					"operators_authentication": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"user_assigned_identities": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"control_plane_operators": {
												Type:     pluginsdk.TypeMap,
												Required: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: commonids.ValidateUserAssignedIdentityID,
												},
											},
											"data_plane_operators": {
												Type:     pluginsdk.TypeMap,
												Required: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: commonids.ValidateUserAssignedIdentityID,
												},
											},
											"service_managed_identity": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: commonids.ValidateUserAssignedIdentityID,
											},
										},
									},
								},
							},
						},
					},
					"issuer_url": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"version": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^\d+\.\d+$`),
							"must be in the format `X.Y` (for example, `4.20`)",
						),
					},
					"channel_group": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  "stable",
						ValidateFunc: validation.StringInSlice([]string{
							"nightly",
							"candidate",
							"fast",
							"stable",
						}, false),
					},
				},
			},
		},

		"api": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"url": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"visibility": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(hcpopenshiftclusters.PossibleValuesForVisibility(), false),
						Default:      hcpopenshiftclusters.VisibilityPublic,
					},
					"authorized_cidrs": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						MinItems: 1,
						MaxItems: 500,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsCIDRIPv4,
						},
					},
				},
			},
		},

		"autoscaling": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"max_nodes_total": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
						Default:      0,
					},
					"max_pod_grace_period_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
						Default:      600,
					},
					"max_node_provision_time_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
						Default:      900,
					},
					"pod_priority_threshold": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  -10,
					},
				},
			},
		},

		"cluster_image_registry": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"state": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(hcpopenshiftclusters.ClusterImageRegistryStateEnabled),
						ValidateFunc: validation.StringInSlice(hcpopenshiftclusters.PossibleValuesForClusterImageRegistryState(), false),
					},
				},
			},
		},

		"crypto_restrictions": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(hcpopenshiftclusters.CryptoRestrictionsNone),
			ValidateFunc: validation.StringInSlice(hcpopenshiftclusters.PossibleValuesForCryptoRestrictions(), false),
		},

		"dns": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"base_domain": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"base_domain_prefix": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ValidateFunc: validation.All(
							validation.StringLenBetween(1, 15),
							validation.StringMatch(
								regexp.MustCompile(`^[a-z]([-a-z0-9]*[a-z0-9])?$`),
								"must begin with a lowercase letter, end with a lowercase letter or number, and contain only lowercase letters, numbers and hyphens",
							),
						),
					},
				},
			},
		},

		"image_digest_mirrors": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 240,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"source": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 255),
					},
					"mirrors": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						MaxItems: 255,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
					},
				},
			},
		},

		"ingress": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(hcpopenshiftclusters.PossibleValuesForIngressType(), false),
						Default:      hcpopenshiftclusters.IngressTypePublic,
					},
				},
			},
		},

		"network": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"network_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(hcpopenshiftclusters.PossibleValuesForNetworkType(), false),
						Default:      hcpopenshiftclusters.NetworkTypeOVNKubernetes,
					},
					"pod_cidr": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsCIDRIPv4,
						Default:      "10.128.0.0/14",
					},
					"service_cidr": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsCIDRIPv4,
						Default:      "172.30.0.0/16",
					},
					"machine_cidr": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsCIDRIPv4,
						Default:      "10.0.0.0/16",
					},
					"host_prefix": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 32),
						Default:      23,
					},
				},
			},
		},

		"node_drain_timeout_minutes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 10080),
			Default:      0,
		},

		"tags": commonschema.Tags(),
	}
}

func (r RedHatHcpOpenShiftCluster) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"console_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r RedHatHcpOpenShiftCluster) ModelObject() any {
	return &RedHatHcpOpenShiftClusterModel{}
}

func (r RedHatHcpOpenShiftCluster) ResourceType() string {
	return "azurerm_redhat_hcp_openshift_cluster"
}

func (r RedHatHcpOpenShiftCluster) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return hcpopenshiftclusters.ValidateHcpOpenShiftClusterID
}

func (r RedHatHcpOpenShiftCluster) Identity() resourceids.ResourceId {
	return &hcpopenshiftclusters.HcpOpenShiftClusterId{}
}

func (r RedHatHcpOpenShiftCluster) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RedHatOpenShift.HcpOpenShiftClustersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config RedHatHcpOpenShiftClusterModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := hcpopenshiftclusters.NewHcpOpenShiftClusterID(subscriptionId, config.ResourceGroup, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := hcpopenshiftclusters.HcpOpenShiftCluster{
				Location: location.Normalize(config.Location),
				Properties: &hcpopenshiftclusters.HcpOpenShiftClusterProperties{
					Api:                     expandHcpOpenShiftClusterApi(config.Api),
					Autoscaling:             expandHcpOpenShiftClusterAutoscaling(config.Autoscaling),
					ClusterImageRegistry:    expandHcpOpenShiftClusterImageRegistry(config.ClusterImageRegistry),
					CryptoRestrictions:      pointer.ToEnum[hcpopenshiftclusters.CryptoRestrictions](config.CryptoRestrictions),
					Dns:                     expandHcpOpenShiftClusterDns(config.Dns),
					Etcd:                    expandHcpOpenShiftClusterEtcd(config.Etcd),
					ImageDigestMirrors:      expandHcpOpenShiftClusterImageDigestMirrors(config.ImageDigestMirrors),
					Ingress:                 expandHcpOpenShiftClusterIngress(config.Ingress),
					Network:                 expandHcpOpenShiftClusterNetwork(config.Network),
					NodeDrainTimeoutMinutes: pointer.To(config.NodeDrainTimeoutMinutes),
					Platform:                expandHcpOpenShiftClusterPlatform(config.Platform),
					Version:                 expandHcpOpenShiftClusterVersion(config.Version),
				},
				Tags: pointer.To(config.Tags),
			}

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(config.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			parameters.Identity = expandedIdentity

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func (r RedHatHcpOpenShiftCluster) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RedHatOpenShift.HcpOpenShiftClustersClient

			id, err := hcpopenshiftclusters.ParseHcpOpenShiftClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state RedHatHcpOpenShiftClusterModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			parameters := *existing.Model

			parameters.SystemData = nil
			parameters.Properties.ProvisioningState = nil
			parameters.Properties.Console = nil
			parameters.Properties.Status = nil
			if parameters.Properties.Api != nil {
				parameters.Properties.Api.Url = ""
			}
			if parameters.Properties.Dns != nil {
				parameters.Properties.Dns.BaseDomain = nil
			}
			parameters.Properties.Platform.IssuerURL = ""

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(state.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				parameters.Identity = expandedIdentity
			}
			if metadata.ResourceData.HasChange("api.0.authorized_cidrs") {
				parameters.Properties.Api = expandHcpOpenShiftClusterApi(state.Api)
			}
			if metadata.ResourceData.HasChange("autoscaling") {
				parameters.Properties.Autoscaling = expandHcpOpenShiftClusterAutoscaling(state.Autoscaling)
			}
			if metadata.ResourceData.HasChange("etcd.0.data_encryption.0.customer_managed.0.kms.0.active_key.0.version") {
				parameters.Properties.Etcd = expandHcpOpenShiftClusterEtcd(state.Etcd)
			}
			if metadata.ResourceData.HasChange("image_digest_mirrors") {
				parameters.Properties.ImageDigestMirrors = expandHcpOpenShiftClusterImageDigestMirrors(state.ImageDigestMirrors)
			}
			if metadata.ResourceData.HasChange("node_drain_timeout_minutes") {
				parameters.Properties.NodeDrainTimeoutMinutes = pointer.To(state.NodeDrainTimeoutMinutes)
			}
			if metadata.ResourceData.HasChange("platform.0.operators_authentication") || metadata.ResourceData.HasChange("platform.0.container_registry") {
				parameters.Properties.Platform = expandHcpOpenShiftClusterPlatform(state.Platform)
			}
			if metadata.ResourceData.HasChange("version") {
				parameters.Properties.Version = expandHcpOpenShiftClusterVersion(state.Version)
			}
			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = pointer.To(state.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r RedHatHcpOpenShiftCluster) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RedHatOpenShift.HcpOpenShiftClustersClient

			id, err := hcpopenshiftclusters.ParseHcpOpenShiftClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			return r.flatten(metadata, *id, resp.Model)
		},
	}
}

func (RedHatHcpOpenShiftCluster) flatten(metadata sdk.ResourceMetaData, id hcpopenshiftclusters.HcpOpenShiftClusterId, model *hcpopenshiftclusters.HcpOpenShiftCluster) error {
	state := RedHatHcpOpenShiftClusterModel{
		Name:          id.HcpOpenShiftClusterName,
		ResourceGroup: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.Tags = pointer.From(model.Tags)

		flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		state.Identity = flattenedIdentity

		if props := model.Properties; props != nil {
			state.Api = flattenHcpOpenShiftClusterApi(props.Api)
			state.Autoscaling = flattenHcpOpenShiftClusterAutoscaling(props.Autoscaling)
			state.ClusterImageRegistry = flattenHcpOpenShiftClusterImageRegistry(props.ClusterImageRegistry)
			state.CryptoRestrictions = pointer.FromEnum(props.CryptoRestrictions)
			state.Dns = flattenHcpOpenShiftClusterDns(props.Dns)
			state.Etcd = flattenHcpOpenShiftClusterEtcd(props.Etcd)
			state.ImageDigestMirrors = flattenHcpOpenShiftClusterImageDigestMirrors(props.ImageDigestMirrors)
			state.Ingress = flattenHcpOpenShiftClusterIngress(props.Ingress)
			state.Network = flattenHcpOpenShiftClusterNetwork(props.Network)
			state.NodeDrainTimeoutMinutes = pointer.From(props.NodeDrainTimeoutMinutes)
			platform, err := flattenHcpOpenShiftClusterPlatform(&props.Platform)
			if err != nil {
				return fmt.Errorf("flattening `platform`: %+v", err)
			}
			state.Platform = platform
			state.Version = flattenHcpOpenShiftClusterVersion(&props.Version)
			if props.Console != nil {
				state.ConsoleUrl = props.Console.Url
			}
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}

func (r RedHatHcpOpenShiftCluster) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := hcpopenshiftclusters.ParseHcpOpenShiftClusterID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing resource ID: %+v", err)
			}

			client := metadata.Client.RedHatOpenShift.HcpOpenShiftClustersClient
			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			return nil
		},
	}
}

func expandHcpOpenShiftClusterApi(input []HcpOpenShiftClusterApiProfile) *hcpopenshiftclusters.ApiProfile {
	if len(input) == 0 {
		return nil
	}
	result := &hcpopenshiftclusters.ApiProfile{
		Visibility: pointer.ToEnum[hcpopenshiftclusters.Visibility](input[0].Visibility),
	}
	if len(input[0].AuthorizedCidrs) > 0 {
		result.AuthorizedCidrs = pointer.To(input[0].AuthorizedCidrs)
	}
	return result
}

func flattenHcpOpenShiftClusterApi(input *hcpopenshiftclusters.ApiProfile) []HcpOpenShiftClusterApiProfile {
	if input == nil {
		return []HcpOpenShiftClusterApiProfile{}
	}
	return []HcpOpenShiftClusterApiProfile{{
		Url:             input.Url,
		Visibility:      pointer.FromEnum(input.Visibility),
		AuthorizedCidrs: pointer.From(input.AuthorizedCidrs),
	}}
}

func expandHcpOpenShiftClusterAutoscaling(input []HcpOpenShiftClusterAutoscalingProfile) *hcpopenshiftclusters.ClusterAutoscalingProfile {
	if len(input) == 0 {
		return nil
	}
	return &hcpopenshiftclusters.ClusterAutoscalingProfile{
		MaxNodesTotal:               pointer.To(input[0].MaxNodesTotal),
		MaxPodGracePeriodSeconds:    pointer.To(input[0].MaxPodGracePeriodSeconds),
		MaxNodeProvisionTimeSeconds: pointer.To(input[0].MaxNodeProvisionTimeSeconds),
		PodPriorityThreshold:        pointer.To(input[0].PodPriorityThreshold),
	}
}

func flattenHcpOpenShiftClusterAutoscaling(input *hcpopenshiftclusters.ClusterAutoscalingProfile) []HcpOpenShiftClusterAutoscalingProfile {
	if input == nil {
		return []HcpOpenShiftClusterAutoscalingProfile{}
	}
	return []HcpOpenShiftClusterAutoscalingProfile{{
		MaxNodesTotal:               pointer.From(input.MaxNodesTotal),
		MaxPodGracePeriodSeconds:    pointer.From(input.MaxPodGracePeriodSeconds),
		MaxNodeProvisionTimeSeconds: pointer.From(input.MaxNodeProvisionTimeSeconds),
		PodPriorityThreshold:        pointer.From(input.PodPriorityThreshold),
	}}
}

func expandHcpOpenShiftClusterDns(input []HcpOpenShiftClusterDnsProfile) *hcpopenshiftclusters.DnsProfile {
	if len(input) == 0 {
		return nil
	}
	return &hcpopenshiftclusters.DnsProfile{
		BaseDomainPrefix: pointer.To(input[0].BaseDomainPrefix),
	}
}

func flattenHcpOpenShiftClusterDns(input *hcpopenshiftclusters.DnsProfile) []HcpOpenShiftClusterDnsProfile {
	if input == nil {
		return []HcpOpenShiftClusterDnsProfile{}
	}
	return []HcpOpenShiftClusterDnsProfile{{
		BaseDomain:       pointer.From(input.BaseDomain),
		BaseDomainPrefix: pointer.From(input.BaseDomainPrefix),
	}}
}

func expandHcpOpenShiftClusterNetwork(input []HcpOpenShiftClusterNetworkProfile) *hcpopenshiftclusters.NetworkProfile {
	if len(input) == 0 {
		return nil
	}
	return &hcpopenshiftclusters.NetworkProfile{
		NetworkType: pointer.ToEnum[hcpopenshiftclusters.NetworkType](input[0].NetworkType),
		PodCidr:     pointer.To(input[0].PodCidr),
		ServiceCidr: pointer.To(input[0].ServiceCidr),
		MachineCidr: pointer.To(input[0].MachineCidr),
		HostPrefix:  pointer.To(input[0].HostPrefix),
	}
}

func flattenHcpOpenShiftClusterNetwork(input *hcpopenshiftclusters.NetworkProfile) []HcpOpenShiftClusterNetworkProfile {
	if input == nil {
		return []HcpOpenShiftClusterNetworkProfile{}
	}
	return []HcpOpenShiftClusterNetworkProfile{{
		NetworkType: pointer.FromEnum(input.NetworkType),
		PodCidr:     pointer.From(input.PodCidr),
		ServiceCidr: pointer.From(input.ServiceCidr),
		MachineCidr: pointer.From(input.MachineCidr),
		HostPrefix:  pointer.From(input.HostPrefix),
	}}
}

func expandHcpOpenShiftClusterIngress(input []HcpOpenShiftClusterIngressProfile) *hcpopenshiftclusters.IngressProfile {
	if len(input) == 0 {
		return nil
	}
	return &hcpopenshiftclusters.IngressProfile{
		Type: pointer.ToEnum[hcpopenshiftclusters.IngressType](input[0].Type),
	}
}

func flattenHcpOpenShiftClusterIngress(input *hcpopenshiftclusters.IngressProfile) []HcpOpenShiftClusterIngressProfile {
	if input == nil {
		return []HcpOpenShiftClusterIngressProfile{}
	}
	return []HcpOpenShiftClusterIngressProfile{{Type: pointer.FromEnum(input.Type)}}
}

func expandHcpOpenShiftClusterPlatform(input []HcpOpenShiftClusterPlatformProfile) hcpopenshiftclusters.PlatformProfile {
	if len(input) == 0 {
		return hcpopenshiftclusters.PlatformProfile{}
	}
	result := hcpopenshiftclusters.PlatformProfile{
		ContainerRegistry:       expandHcpOpenShiftClusterContainerRegistry(input[0].ContainerRegistry),
		SubnetId:                input[0].SubnetId,
		VnetIntegrationSubnetId: input[0].VnetIntegrationSubnetId,
		OutboundType:            pointer.ToEnum[hcpopenshiftclusters.OutboundType](input[0].OutboundType),
		NetworkSecurityGroupId:  input[0].NetworkSecurityGroupId,
		OperatorsAuthentication: expandHcpOpenShiftClusterOperatorsAuthentication(input[0].OperatorsAuthentication),
	}
	if input[0].ManagedResourceGroup != "" {
		result.ManagedResourceGroup = pointer.To(input[0].ManagedResourceGroup)
	}
	return result
}

func flattenHcpOpenShiftClusterPlatform(input *hcpopenshiftclusters.PlatformProfile) ([]HcpOpenShiftClusterPlatformProfile, error) {
	if input == nil {
		return []HcpOpenShiftClusterPlatformProfile{}, nil
	}

	subnetId, err := commonids.ParseSubnetIDInsensitively(input.SubnetId)
	if err != nil {
		return nil, fmt.Errorf("parsing `subnet_id`: %+v", err)
	}

	vnetIntegrationSubnetId, err := commonids.ParseSubnetIDInsensitively(input.VnetIntegrationSubnetId)
	if err != nil {
		return nil, fmt.Errorf("parsing `vnet_integration_subnet_id`: %+v", err)
	}

	networkSecurityGroupId, err := networksecuritygroups.ParseNetworkSecurityGroupIDInsensitively(input.NetworkSecurityGroupId)
	if err != nil {
		return nil, fmt.Errorf("parsing `network_security_group_id`: %+v", err)
	}

	operatorsAuthentication, err := flattenHcpOpenShiftClusterOperatorsAuthentication(&input.OperatorsAuthentication)
	if err != nil {
		return nil, fmt.Errorf("flattening `operators_authentication`: %+v", err)
	}

	containerRegistry, err := flattenHcpOpenShiftClusterContainerRegistry(input.ContainerRegistry)
	if err != nil {
		return nil, fmt.Errorf("flattening `container_registry`: %+v", err)
	}

	return []HcpOpenShiftClusterPlatformProfile{{
		ContainerRegistry:       containerRegistry,
		ManagedResourceGroup:    pointer.From(input.ManagedResourceGroup),
		SubnetId:                subnetId.ID(),
		VnetIntegrationSubnetId: vnetIntegrationSubnetId.ID(),
		OutboundType:            pointer.FromEnum(input.OutboundType),
		NetworkSecurityGroupId:  networkSecurityGroupId.ID(),
		OperatorsAuthentication: operatorsAuthentication,
		IssuerUrl:               input.IssuerURL,
	}}, nil
}

func expandHcpOpenShiftClusterContainerRegistry(input []HcpOpenShiftClusterContainerRegistryProfile) *hcpopenshiftclusters.ContainerRegistryProfile {
	if len(input) == 0 {
		return nil
	}
	return &hcpopenshiftclusters.ContainerRegistryProfile{
		ManagedIdentity: pointer.To(input[0].ManagedIdentity),
	}
}

func flattenHcpOpenShiftClusterContainerRegistry(input *hcpopenshiftclusters.ContainerRegistryProfile) ([]HcpOpenShiftClusterContainerRegistryProfile, error) {
	if input == nil || input.ManagedIdentity == nil {
		return nil, nil
	}

	id, err := commonids.ParseUserAssignedIdentityIDInsensitively(*input.ManagedIdentity)
	if err != nil {
		return nil, fmt.Errorf("parsing `managed_identity`: %+v", err)
	}
	return []HcpOpenShiftClusterContainerRegistryProfile{{ManagedIdentity: id.ID()}}, nil
}

func expandHcpOpenShiftClusterOperatorsAuthentication(input []HcpOpenShiftClusterOperatorsAuthentication) hcpopenshiftclusters.OperatorsAuthenticationProfile {
	result := hcpopenshiftclusters.OperatorsAuthenticationProfile{}
	if len(input) == 0 || len(input[0].UserAssignedIdentities) == 0 {
		return result
	}
	identities := input[0].UserAssignedIdentities[0]
	result.UserAssignedIdentities = hcpopenshiftclusters.UserAssignedIdentitiesProfile{
		ControlPlaneOperators:  identities.ControlPlaneOperators,
		DataPlaneOperators:     identities.DataPlaneOperators,
		ServiceManagedIdentity: identities.ServiceManagedIdentity,
	}
	return result
}

func flattenHcpOpenShiftClusterOperatorsAuthentication(input *hcpopenshiftclusters.OperatorsAuthenticationProfile) ([]HcpOpenShiftClusterOperatorsAuthentication, error) {
	if input == nil {
		return []HcpOpenShiftClusterOperatorsAuthentication{}, nil
	}

	controlPlaneOperators, err := flattenHcpOpenShiftClusterUserAssignedIdentityMap(input.UserAssignedIdentities.ControlPlaneOperators)
	if err != nil {
		return nil, fmt.Errorf("parsing `control_plane_operators`: %+v", err)
	}

	dataPlaneOperators, err := flattenHcpOpenShiftClusterUserAssignedIdentityMap(input.UserAssignedIdentities.DataPlaneOperators)
	if err != nil {
		return nil, fmt.Errorf("parsing `data_plane_operators`: %+v", err)
	}

	serviceManagedIdentity, err := commonids.ParseUserAssignedIdentityIDInsensitively(input.UserAssignedIdentities.ServiceManagedIdentity)
	if err != nil {
		return nil, fmt.Errorf("parsing `service_managed_identity`: %+v", err)
	}

	return []HcpOpenShiftClusterOperatorsAuthentication{{
		UserAssignedIdentities: []HcpOpenShiftClusterUserAssignedIdentities{{
			ControlPlaneOperators:  controlPlaneOperators,
			DataPlaneOperators:     dataPlaneOperators,
			ServiceManagedIdentity: serviceManagedIdentity.ID(),
		}},
	}}, nil
}

func flattenHcpOpenShiftClusterUserAssignedIdentityMap(input map[string]string) (map[string]string, error) {
	if input == nil {
		return nil, nil
	}

	result := make(map[string]string, 0)
	for name, raw := range input {
		id, err := commonids.ParseUserAssignedIdentityIDInsensitively(raw)
		if err != nil {
			return nil, fmt.Errorf("parsing identity %q: %+v", name, err)
		}
		result[name] = id.ID()
	}
	return result, nil
}

func expandHcpOpenShiftClusterImageRegistry(input []HcpOpenShiftClusterImageRegistryProfile) *hcpopenshiftclusters.ClusterImageRegistryProfile {
	if len(input) == 0 {
		return nil
	}
	return &hcpopenshiftclusters.ClusterImageRegistryProfile{
		State: pointer.ToEnum[hcpopenshiftclusters.ClusterImageRegistryState](input[0].State),
	}
}

func flattenHcpOpenShiftClusterImageRegistry(input *hcpopenshiftclusters.ClusterImageRegistryProfile) []HcpOpenShiftClusterImageRegistryProfile {
	if input == nil {
		return []HcpOpenShiftClusterImageRegistryProfile{}
	}
	return []HcpOpenShiftClusterImageRegistryProfile{{State: pointer.FromEnum(input.State)}}
}

func expandHcpOpenShiftClusterVersion(input []HcpOpenShiftClusterVersionProfile) hcpopenshiftclusters.VersionProfile {
	if len(input) == 0 {
		return hcpopenshiftclusters.VersionProfile{}
	}
	return hcpopenshiftclusters.VersionProfile{
		Id:           input[0].Id,
		ChannelGroup: pointer.To(input[0].ChannelGroup),
	}
}

func flattenHcpOpenShiftClusterVersion(input *hcpopenshiftclusters.VersionProfile) []HcpOpenShiftClusterVersionProfile {
	if input == nil {
		return []HcpOpenShiftClusterVersionProfile{}
	}
	return []HcpOpenShiftClusterVersionProfile{{
		Id:           input.Id,
		ChannelGroup: pointer.From(input.ChannelGroup),
	}}
}

func expandHcpOpenShiftClusterImageDigestMirrors(input []HcpOpenShiftClusterImageDigestMirror) *[]hcpopenshiftclusters.ImageDigestMirror {
	result := make([]hcpopenshiftclusters.ImageDigestMirror, 0, len(input))
	for _, item := range input {
		result = append(result, hcpopenshiftclusters.ImageDigestMirror{
			Source:  item.Source,
			Mirrors: item.Mirrors,
		})
	}
	return &result
}

func flattenHcpOpenShiftClusterImageDigestMirrors(input *[]hcpopenshiftclusters.ImageDigestMirror) []HcpOpenShiftClusterImageDigestMirror {
	if input == nil {
		return []HcpOpenShiftClusterImageDigestMirror{}
	}
	result := make([]HcpOpenShiftClusterImageDigestMirror, 0, len(*input))
	for _, item := range *input {
		result = append(result, HcpOpenShiftClusterImageDigestMirror{
			Source:  item.Source,
			Mirrors: item.Mirrors,
		})
	}
	return result
}

func expandHcpOpenShiftClusterEtcd(input []HcpOpenShiftClusterEtcdProfile) *hcpopenshiftclusters.EtcdProfile {
	if len(input) == 0 || len(input[0].DataEncryption) == 0 {
		return nil
	}
	dataEncryption := input[0].DataEncryption[0]
	result := &hcpopenshiftclusters.EtcdProfile{
		DataEncryption: &hcpopenshiftclusters.EtcdDataEncryptionProfile{
			KeyManagementMode: hcpopenshiftclusters.EtcdDataEncryptionKeyManagementModeType(dataEncryption.KeyManagementMode),
		},
	}
	if len(dataEncryption.CustomerManaged) > 0 {
		customerManaged := dataEncryption.CustomerManaged[0]
		result.DataEncryption.CustomerManaged = &hcpopenshiftclusters.CustomerManagedEncryptionProfile{
			EncryptionType: pointer.ToEnum[hcpopenshiftclusters.CustomerManagedEncryptionType](customerManaged.EncryptionType),
		}
		if len(customerManaged.Kms) > 0 {
			kms := customerManaged.Kms[0]
			result.DataEncryption.CustomerManaged.Kms = &hcpopenshiftclusters.KmsEncryptionProfile{
				VaultName:  kms.VaultName,
				Visibility: hcpopenshiftclusters.KeyVaultVisibility(kms.Visibility),
			}
			if kms.KeyVaultType != "" {
				result.DataEncryption.CustomerManaged.Kms.KeyVaultType = pointer.ToEnum[hcpopenshiftclusters.KmsKeyVaultType](kms.KeyVaultType)
			}
			if len(kms.ActiveKey) > 0 {
				result.DataEncryption.CustomerManaged.Kms.ActiveKey = hcpopenshiftclusters.KmsKey{
					Name:    kms.ActiveKey[0].Name,
					Version: kms.ActiveKey[0].Version,
				}
			}
		}
	}
	return result
}

func flattenHcpOpenShiftClusterEtcd(input *hcpopenshiftclusters.EtcdProfile) []HcpOpenShiftClusterEtcdProfile {
	if input == nil || input.DataEncryption == nil {
		return []HcpOpenShiftClusterEtcdProfile{}
	}
	dataEncryption := HcpOpenShiftClusterEtcdDataEncryptionProfile{
		KeyManagementMode: string(input.DataEncryption.KeyManagementMode),
	}
	if input.DataEncryption.CustomerManaged != nil {
		customerManaged := input.DataEncryption.CustomerManaged
		flattenedCustomerManaged := HcpOpenShiftClusterCustomerManagedEncryption{
			EncryptionType: pointer.FromEnum(customerManaged.EncryptionType),
		}
		if customerManaged.Kms != nil {
			flattenedCustomerManaged.Kms = []HcpOpenShiftClusterKmsEncryption{{
				VaultName:    customerManaged.Kms.VaultName,
				KeyVaultType: pointer.FromEnum(customerManaged.Kms.KeyVaultType),
				Visibility:   string(customerManaged.Kms.Visibility),
				ActiveKey: []HcpOpenShiftClusterKmsKey{{
					Name:    customerManaged.Kms.ActiveKey.Name,
					Version: customerManaged.Kms.ActiveKey.Version,
				}},
			}}
		}
		dataEncryption.CustomerManaged = []HcpOpenShiftClusterCustomerManagedEncryption{flattenedCustomerManaged}
	}
	return []HcpOpenShiftClusterEtcdProfile{{DataEncryption: []HcpOpenShiftClusterEtcdDataEncryptionProfile{dataEncryption}}}
}
