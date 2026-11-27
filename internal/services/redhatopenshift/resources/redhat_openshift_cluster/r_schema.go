// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redhat_openshift_cluster

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2025-07-25/openshiftclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redhatopenshift/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			outboundType := metadata.ResourceDiff.Get("network_profile.0.outbound_type").(string)
			loadBalancerProfile := metadata.ResourceDiff.Get("network_profile.0.load_balancer_profile")
			if outboundType != string(openshiftclusters.OutboundTypeLoadbalancer) && len(loadBalancerProfile.([]any)) > 0 {
				return errors.New("`network_profile.0.load_balancer_profile` requires `network_profile.0.outbound_type` to be `Loadbalancer`")
			}

			if metadata.ResourceDiff.Id() == "" {
				return nil
			}

			oldUpgradeableTo, newUpgradeableTo := metadata.ResourceDiff.GetChange("platform_workload_identity_profile.0.upgradeable_to")
			if oldUpgradeableTo.(string) != "" && newUpgradeableTo.(string) == "" {
				return errors.New("`platform_workload_identity_profile.0.upgradeable_to` cannot be removed once set")
			}

			// The service principal can be swapped, but it cannot be added or removed in place.
			oldServicePrincipal, newServicePrincipal := metadata.ResourceDiff.GetChange("service_principal")
			if len(oldServicePrincipal.([]any)) != len(newServicePrincipal.([]any)) {
				return metadata.ResourceDiff.ForceNew("service_principal")
			}

			// The identity can be swapped, but it cannot be added or removed in place.
			oldIdentity, newIdentity := metadata.ResourceDiff.GetChange("identity")
			if len(oldIdentity.([]any)) != len(newIdentity.([]any)) {
				return metadata.ResourceDiff.ForceNew("identity")
			}

			return nil
		},
	}
}

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Azure Red Hat OpenShift Cluster to create.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"api_server_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"visibility": {
						Description:  "Cluster API server visibility.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(openshiftclusters.PossibleValuesForVisibility(), false),
					},
					"ip_address": {
						Description: "The IP Address the API Server Profile is associated with.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"url": {
						Description: "The URL the API Server Profile is associated with.",
						Type:        schema.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"cluster_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"domain": {
						Description:  "The custom domain for the cluster. For more info, see [Prepare a custom domain for your cluster](https://learn.microsoft.com/azure/openshift/tutorial-create-cluster#prepare-a-custom-domain-for-your-cluster-optional).",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"version": {
						Description:  "The version of the OpenShift cluster. Available versions can be found with the Azure CLI command `az aro get-versions --location <region>`.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.ClusterVersion,
					},
					"fips_enabled": {
						Description: "Whether Federal Information Processing Standard (FIPS) validated cryptographic modules are used.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						ForceNew:    true,
						Default:     false,
					},
					"managed_resource_group_name": {
						Description:  "The name of a Resource Group which will be created to host VMs of Azure Red Hat OpenShift Cluster. The value cannot contain uppercase characters.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validate.ClusterResourceGroupName,
						DiffSuppressFunc: func(_, old, new string, d *pluginsdk.ResourceData) bool {
							defaultResourceGroupName := "aro-" + d.Get("cluster_profile.0.domain").(string)
							if old == defaultResourceGroupName && new == "" {
								return true
							}
							return false
						},
					},
					"pull_secret": {
						Description:  "The Red Hat pull secret for the cluster. For more info, see [Get a Red Hat pull secret](https://learn.microsoft.com/azure/openshift/tutorial-create-cluster#get-a-red-hat-pull-secret-optional).",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"resource_group_id": {
						Description: "The Resource Group that the cluster profile is attached to.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"ingress_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"visibility": {
						Description:  "Cluster Ingress visibility.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(openshiftclusters.PossibleValuesForVisibility(), false),
					},
					"ip_address": {
						Description: "The IP Address the Ingress Profile is associated with.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"name": {
						Description: "The name of the Ingress Profile.",
						Type:        schema.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"main_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Description:  "The ID of the subnet where main nodes will be hosted.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
					"vm_size": {
						Description:      "The size of the Virtual Machines for the main nodes.",
						Type:             pluginsdk.TypeString,
						Required:         true,
						ForceNew:         true,
						ValidateFunc:     validation.StringIsNotEmpty,
						DiffSuppressFunc: suppress.CaseDifference,
					},
					"disk_encryption_set_id": {
						Description:  "The resource ID of an associated Disk Encryption Set.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateDiskEncryptionSetID,
					},
					"encryption_at_host_enabled": {
						Description: "Whether main Virtual Machines are encrypted at host.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						ForceNew:    true,
						Default:     false,
					},
				},
			},
		},

		"network_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"pod_cidr": {
						Description:  "The CIDR to use for pod IP addresses.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsCIDRIPv4,
					},
					"service_cidr": {
						Description:  "The network range used by the OpenShift service.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsCIDRIPv4,
					},
					"load_balancer_profile": {
						Description: "Only applicable when `outbound_type` is set to `Loadbalancer`.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						// NOTE: O+C - Azure returns the default load balancer profile, including effective outbound IPs, when this block is omitted.
						Computed: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"managed_outbound_ip_count": {
									Description:  "The desired number of managed outbound public IP addresses created and assigned to the cluster's outbound load balancer.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 20),
								},
								"effective_outbound_ips": {
									Description: "The list of effective outbound IP resource IDs of the cluster's public load balancer.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
					"outbound_type": {
						Description: "The outbound (egress) routing method.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
						Default:     string(openshiftclusters.OutboundTypeLoadbalancer),
						ValidateFunc: validation.StringInSlice(
							openshiftclusters.PossibleValuesForOutboundType(),
							false,
						),
					},
					"preconfigured_network_security_group_enabled": {
						Description: "Whether a preconfigured Network Security Group is being used on the subnets.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						ForceNew:    true,
						Default:     false,
					},
				},
			},
		},

		"worker_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disk_size_gb": {
						Description:  "The internal OS disk size of the worker Virtual Machines in GB.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IntAtLeast(128),
					},
					"node_count": {
						Description:  "The initial number of worker nodes which should exist in the cluster.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IntBetween(3, 60),
					},
					"subnet_id": {
						Description:  "The ID of the subnet where worker nodes will be hosted.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
					"vm_size": {
						Description:      "The size of the Virtual Machines for the worker nodes.",
						Type:             pluginsdk.TypeString,
						Required:         true,
						ForceNew:         true,
						ValidateFunc:     validation.StringIsNotEmpty,
						DiffSuppressFunc: suppress.CaseDifference,
					},
					"disk_encryption_set_id": {
						Description:  "The resource ID of an associated Disk Encryption Set.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateDiskEncryptionSetID,
					},
					"encryption_at_host_enabled": {
						Description: "Whether worker Virtual Machines are encrypted at host.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						ForceNew:    true,
						Default:     false,
					},
				},
			},
		},

		"identity": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			RequiredWith: []string{"platform_workload_identity_profile"},
			ExactlyOneOf: []string{"service_principal", "identity"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Description: "The type of Managed Service Identity assigned to the cluster.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							string(identity.TypeUserAssigned),
						}, false),
					},
					// lintignore:S018
					"identity_ids": {
						Description: "A set of User Assigned Managed Identity IDs to assign to the cluster. Exactly one cluster identity is required for now.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						MinItems:    1,
						MaxItems:    1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: commonids.ValidateUserAssignedIdentityID,
						},
					},
				},
			},
		},

		"platform_workload_identity_profile": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			RequiredWith: []string{"identity"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"platform_workload_identity": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"identity_id": {
									Description:  "The resource ID of the User Assigned Managed Identity to assign to the operator.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: commonids.ValidateUserAssignedIdentityID,
								},
								"name": {
									Description:  "The name of the platform workload identity operator.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"client_id": {
									Description: "The client ID of the User Assigned Managed Identity assigned to the operator.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"object_id": {
									Description: "The object (principal) ID of the User Assigned Managed Identity assigned to the operator.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},
					"upgradeable_to": {
						Description:  "The target OpenShift version (`x.y.z`) the platform workload identities should satisfy. Once set, this value cannot be cleared.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.ClusterVersion,
					},
				},
			},
		},

		"service_principal": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"service_principal", "identity"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Description:  "The Client ID for the Service Principal.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"client_secret": {
						Description:  "The Client Secret for the Service Principal.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"console_url": {
			Description: "The console URL of the Azure Red Hat OpenShift Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
