// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redhat_openshift_cluster

import (
	"context"
	"errors"
	"fmt"
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
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Required:     true,
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(openshiftclusters.PossibleValuesForVisibility(), false),
					},
					"ip_address": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"url": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},

		"cluster_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"domain": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"version": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.ClusterVersion,
					},
					"fips_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},
					"managed_resource_group_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validate.ClusterResourceGroupName,
						DiffSuppressFunc: func(_, old, new string, d *pluginsdk.ResourceData) bool {
							defaultResourceGroupName := fmt.Sprintf("aro-%s", d.Get("cluster_profile.0.domain").(string))
							if old == defaultResourceGroupName && new == "" {
								return true
							}
							return false
						},
					},
					"pull_secret": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"resource_group_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(openshiftclusters.PossibleValuesForVisibility(), false),
					},
					"ip_address": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
					"vm_size": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						ForceNew:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						ValidateFunc:     validation.StringIsNotEmpty,
					},
					"disk_encryption_set_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateDiskEncryptionSetID,
					},
					"encryption_at_host_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsCIDRIPv4,
					},
					"service_cidr": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsCIDRIPv4,
					},
					"load_balancer_profile": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						// NOTE: O+C - Azure returns the default load balancer profile, including effective outbound IPs, when this block is omitted.
						Computed: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"managed_outbound_ip_count": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 20),
								},
								"effective_outbound_ips": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
					"outbound_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  string(openshiftclusters.OutboundTypeLoadbalancer),
						ValidateFunc: validation.StringInSlice(
							openshiftclusters.PossibleValuesForOutboundType(),
							false,
						),
					},
					"preconfigured_network_security_group_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
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
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IntAtLeast(128),
					},
					"node_count": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IntBetween(3, 60),
					},
					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
					"vm_size": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						ForceNew:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						ValidateFunc:     validation.StringIsNotEmpty,
					},
					"disk_encryption_set_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateDiskEncryptionSetID,
					},
					"encryption_at_host_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},
				},
			},
		},

		"identity": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"service_principal", "identity"},
			RequiredWith: []string{"platform_workload_identity_profile"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(identity.TypeUserAssigned),
						}, false),
					},
					// lintignore:S018
					"identity_ids": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						MinItems: 1,
						MaxItems: 1,
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
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: commonids.ValidateUserAssignedIdentityID,
								},
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"client_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"object_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
					"upgradeable_to": {
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"client_secret": {
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
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
