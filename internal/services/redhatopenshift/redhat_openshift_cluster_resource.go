// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name redhat_openshift_cluster -service-package-name redhatopenshift -properties "name,resource_group_name" -known-values "subscription_id:data.Subscriptions.Primary"

package redhatopenshift

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2025-07-25/openshiftclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commonValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redhatopenshift/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.ResourceWithCustomizeDiff = RedHatOpenShiftCluster{}
	_ sdk.ResourceWithIdentity      = RedHatOpenShiftCluster{}
	_ sdk.ResourceWithUpdate        = RedHatOpenShiftCluster{}
)

type RedHatOpenShiftCluster struct{}

type RedHatOpenShiftClusterModel struct {
	Tags                            map[string]string                 `tfschema:"tags"`
	Name                            string                            `tfschema:"name"`
	Location                        string                            `tfschema:"location"`
	ResourceGroup                   string                            `tfschema:"resource_group_name"`
	ConsoleUrl                      string                            `tfschema:"console_url"`
	Identity                        []identity.ModelUserAssigned      `tfschema:"identity"`
	ServicePrincipal                []ServicePrincipal                `tfschema:"service_principal"`
	ClusterProfile                  []ClusterProfile                  `tfschema:"cluster_profile"`
	NetworkProfile                  []NetworkProfile                  `tfschema:"network_profile"`
	MainProfile                     []MainProfile                     `tfschema:"main_profile"`
	WorkerProfile                   []WorkerProfile                   `tfschema:"worker_profile"`
	ApiServerProfile                []ApiServerProfile                `tfschema:"api_server_profile"`
	IngressProfile                  []IngressProfile                  `tfschema:"ingress_profile"`
	PlatformWorkloadIdentityProfile []PlatformWorkloadIdentityProfile `tfschema:"platform_workload_identity_profile"`
}

type ServicePrincipal struct {
	ClientId     string `tfschema:"client_id"`
	ClientSecret string `tfschema:"client_secret"`
}

type ClusterProfile struct {
	PullSecret               string `tfschema:"pull_secret"`
	Domain                   string `tfschema:"domain"`
	ManagedResourceGroupName string `tfschema:"managed_resource_group_name"`
	ResourceGroupId          string `tfschema:"resource_group_id"`
	Version                  string `tfschema:"version"`
	FipsEnabled              bool   `tfschema:"fips_enabled"`
}

type NetworkProfile struct {
	OutboundType                             string                `tfschema:"outbound_type"`
	PodCidr                                  string                `tfschema:"pod_cidr"`
	ServiceCidr                              string                `tfschema:"service_cidr"`
	PreconfiguredNetworkSecurityGroupEnabled bool                  `tfschema:"preconfigured_network_security_group_enabled"`
	LoadBalancerProfile                      []LoadBalancerProfile `tfschema:"load_balancer_profile"`
}

type LoadBalancerProfile struct {
	ManagedOutboundIpCount int64    `tfschema:"managed_outbound_ip_count"`
	EffectiveOutboundIps   []string `tfschema:"effective_outbound_ips"`
}

type PlatformWorkloadIdentityProfile struct {
	UpgradeableTo            string                     `tfschema:"upgradeable_to"`
	PlatformWorkloadIdentity []PlatformWorkloadIdentity `tfschema:"platform_workload_identity"`
}

type PlatformWorkloadIdentity struct {
	Name       string `tfschema:"name"`
	IdentityId string `tfschema:"identity_id"`
	ClientId   string `tfschema:"client_id"`
	ObjectId   string `tfschema:"object_id"`
}

type MainProfile struct {
	SubnetId                string `tfschema:"subnet_id"`
	VmSize                  string `tfschema:"vm_size"`
	DiskEncryptionSetId     string `tfschema:"disk_encryption_set_id"`
	EncryptionAtHostEnabled bool   `tfschema:"encryption_at_host_enabled"`
}

type WorkerProfile struct {
	VmSize                  string `tfschema:"vm_size"`
	SubnetId                string `tfschema:"subnet_id"`
	DiskEncryptionSetId     string `tfschema:"disk_encryption_set_id"`
	DiskSizeGb              int64  `tfschema:"disk_size_gb"`
	NodeCount               int64  `tfschema:"node_count"`
	EncryptionAtHostEnabled bool   `tfschema:"encryption_at_host_enabled"`
}

type IngressProfile struct {
	Visibility string `tfschema:"visibility"`
	IpAddress  string `tfschema:"ip_address"`
	Name       string `tfschema:"name"`
}

type ApiServerProfile struct {
	Visibility string `tfschema:"visibility"`
	IpAddress  string `tfschema:"ip_address"`
	Url        string `tfschema:"url"`
}

func (r RedHatOpenShiftCluster) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			outboundType := metadata.ResourceDiff.Get("network_profile.0.outbound_type").(string)
			loadBalancerProfile := metadata.ResourceDiff.Get("network_profile.0.load_balancer_profile")
			if outboundType != string(openshiftclusters.OutboundTypeLoadbalancer) && len(loadBalancerProfile.([]interface{})) > 0 {
				return errors.New("`network_profile.0.load_balancer_profile` requires `network_profile.0.outbound_type` to be `Loadbalancer`")
			}

			if metadata.ResourceDiff.Id() == "" {
				return nil
			}

			// The service principal can be swapped, but it cannot be added or removed in place.
			oldServicePrincipal, newServicePrincipal := metadata.ResourceDiff.GetChange("service_principal")
			if len(oldServicePrincipal.([]interface{})) != len(newServicePrincipal.([]interface{})) {
				return metadata.ResourceDiff.ForceNew("service_principal")
			}

			// The identity can be swapped, but it cannot be added or removed in place.
			oldIdentity, newIdentity := metadata.ResourceDiff.GetChange("identity")
			if len(oldIdentity.([]interface{})) != len(newIdentity.([]interface{})) {
				return metadata.ResourceDiff.ForceNew("identity")
			}

			return nil
		},
	}
}

func (r RedHatOpenShiftCluster) Arguments() map[string]*pluginsdk.Schema {
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
					"version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// NOTE: O+C - Azure selects and returns the current default supported OpenShift version when this is omitted.
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: validate.ClusterVersion,
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
						ValidateFunc: commonValidate.CIDR,
					},
					"service_cidr": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonValidate.CIDR,
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
									Optional:     true,
									Default:      1,
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

		"identity": func() *pluginsdk.Schema {
			s := commonschema.UserAssignedIdentityOptional()
			s.ExactlyOneOf = []string{"service_principal", "identity"}
			s.RequiredWith = []string{"platform_workload_identity_profile"}
			if elem, ok := s.Elem.(*schema.Resource); ok {
				if ids, ok := elem.Schema["identity_ids"]; ok {
					ids.MinItems = 1
					ids.MaxItems = 1
				}
			}
			return s
		}(),

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

func (r RedHatOpenShiftCluster) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"console_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r RedHatOpenShiftCluster) ModelObject() interface{} {
	return &RedHatOpenShiftClusterModel{}
}

func (r RedHatOpenShiftCluster) ResourceType() string {
	return "azurerm_redhat_openshift_cluster"
}

func (r RedHatOpenShiftCluster) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return openshiftclusters.ValidateOpenShiftClusterID
}

func (r RedHatOpenShiftCluster) Identity() resourceids.ResourceId {
	return &openshiftclusters.OpenShiftClusterId{}
}

func (r RedHatOpenShiftCluster) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RedHatOpenShift.OpenShiftClustersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config RedHatOpenShiftClusterModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id := openshiftclusters.NewOpenShiftClusterID(subscriptionId, config.ResourceGroup, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %s", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := openshiftclusters.OpenShiftCluster{
				Name:     pointer.To(id.OpenShiftClusterName),
				Location: location.Normalize(config.Location),
				Properties: &openshiftclusters.OpenShiftClusterProperties{
					ClusterProfile:                  expandOpenshiftClusterProfile(config.ClusterProfile, id.SubscriptionId),
					ServicePrincipalProfile:         expandOpenshiftServicePrincipalProfile(config.ServicePrincipal),
					NetworkProfile:                  expandOpenshiftNetworkProfile(config.NetworkProfile),
					MasterProfile:                   expandOpenshiftMainProfile(config.MainProfile),
					WorkerProfiles:                  expandOpenshiftWorkerProfiles(config.WorkerProfile),
					ApiserverProfile:                expandOpenshiftApiServerProfile(config.ApiServerProfile),
					IngressProfiles:                 expandOpenshiftIngressProfiles(config.IngressProfile),
					PlatformWorkloadIdentityProfile: expandOpenshiftPlatformWorkloadIdentityProfile(config.PlatformWorkloadIdentityProfile),
				},
				Tags: pointer.To(config.Tags),
			}

			if len(config.Identity) > 0 {
				expandedIdentity, err := identity.ExpandUserAssignedMapFromModel(config.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				parameters.Identity = expandedIdentity
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, metadata.SetIDAndIdentityCallback(&id)); err != nil {
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

func (r RedHatOpenShiftCluster) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RedHatOpenShift.OpenShiftClustersClient

			id, err := openshiftclusters.ParseOpenShiftClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state RedHatOpenShiftClusterModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			expandedIdentity, err := identity.ExpandUserAssignedMapFromModel(state.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			parameter := openshiftclusters.OpenShiftCluster{
				Name:     pointer.To(id.OpenShiftClusterName),
				Location: location.Normalize(state.Location),
				Identity: expandedIdentity,
				Properties: &openshiftclusters.OpenShiftClusterProperties{
					ClusterProfile:                  expandOpenshiftClusterProfile(state.ClusterProfile, id.SubscriptionId),
					ServicePrincipalProfile:         expandOpenshiftServicePrincipalProfile(state.ServicePrincipal),
					NetworkProfile:                  expandOpenshiftNetworkProfile(state.NetworkProfile),
					MasterProfile:                   expandOpenshiftMainProfile(state.MainProfile),
					WorkerProfiles:                  expandOpenshiftWorkerProfiles(state.WorkerProfile),
					ApiserverProfile:                expandOpenshiftApiServerProfile(state.ApiServerProfile),
					IngressProfiles:                 expandOpenshiftIngressProfiles(state.IngressProfile),
					PlatformWorkloadIdentityProfile: expandOpenshiftPlatformWorkloadIdentityProfile(state.PlatformWorkloadIdentityProfile),
				},
				Tags: pointer.To(state.Tags),
			}

			// Platform workload identity updates require PUT so a workload identity can be updated in place.
			if err := client.CreateOrUpdateThenPoll(ctx, *id, parameter); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r RedHatOpenShiftCluster) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RedHatOpenShift.OpenShiftClustersClient

			id, err := openshiftclusters.ParseOpenShiftClusterID(metadata.ResourceData.Id())
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

			var config RedHatOpenShiftClusterModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			return r.flatten(metadata, *id, resp.Model, config)
		},
	}
}

func (RedHatOpenShiftCluster) flatten(metadata sdk.ResourceMetaData, id openshiftclusters.OpenShiftClusterId, model *openshiftclusters.OpenShiftCluster, config RedHatOpenShiftClusterModel) error {
	state := RedHatOpenShiftClusterModel{
		Name:          id.OpenShiftClusterName,
		ResourceGroup: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.Tags = pointer.From(model.Tags)

		flattenedIdentity, err := identity.FlattenUserAssignedMapToModel(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		state.Identity = pointer.From(flattenedIdentity)

		if props := model.Properties; props != nil {
			clusterProfile, err := flattenOpenShiftClusterProfile(props.ClusterProfile, config)
			if err != nil {
				return fmt.Errorf("flatten cluster profile: %+v", err)
			}
			state.ClusterProfile = *clusterProfile

			state.ServicePrincipal = flattenOpenShiftServicePrincipalProfile(props.ServicePrincipalProfile, config)
			state.NetworkProfile = flattenOpenShiftNetworkProfile(props.NetworkProfile)
			state.MainProfile = flattenOpenShiftMainProfile(props.MasterProfile)
			state.ApiServerProfile = flattenOpenShiftAPIServerProfile(props.ApiserverProfile)
			state.IngressProfile = flattenOpenShiftIngressProfiles(props.IngressProfiles)
			state.PlatformWorkloadIdentityProfile = flattenOpenShiftPlatformWorkloadIdentityProfile(props.PlatformWorkloadIdentityProfile)

			workerProfiles, err := flattenOpenShiftWorkerProfiles(props.WorkerProfiles)
			if err != nil {
				return fmt.Errorf("flattening worker profiles: %+v", err)
			}
			state.WorkerProfile = workerProfiles

			if props.ConsoleProfile != nil {
				state.ConsoleUrl = pointer.From(props.ConsoleProfile.Url)
			}
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r RedHatOpenShiftCluster) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := openshiftclusters.ParseOpenShiftClusterID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client := metadata.Client.RedHatOpenShift.OpenShiftClustersClient

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandOpenshiftClusterProfile(input []ClusterProfile, subscriptionId string) *openshiftclusters.ClusterProfile {
	if len(input) == 0 {
		return nil
	}

	fipsValidatedModules := openshiftclusters.FipsValidatedModulesDisabled
	if input[0].FipsEnabled {
		fipsValidatedModules = openshiftclusters.FipsValidatedModulesEnabled
	}

	// the api needs a ResourceGroupId value and the portal doesn't allow you to set it but the portal returns the
	// resource id being `aro-{domain}` so we'll follow that here.
	resourceGroupId := commonids.NewResourceGroupID(subscriptionId, fmt.Sprintf("aro-%s", input[0].Domain)).ID()
	if rg := input[0].ManagedResourceGroupName; rg != "" {
		resourceGroupId = commonids.NewResourceGroupID(subscriptionId, rg).ID()
	}

	return &openshiftclusters.ClusterProfile{
		ResourceGroupId:      pointer.To(resourceGroupId),
		Domain:               pointer.To(input[0].Domain),
		PullSecret:           pointer.To(input[0].PullSecret),
		FipsValidatedModules: pointer.To(fipsValidatedModules),
		Version:              pointer.To(input[0].Version),
	}
}

func flattenOpenShiftClusterProfile(profile *openshiftclusters.ClusterProfile, config RedHatOpenShiftClusterModel) (*[]ClusterProfile, error) {
	if profile == nil {
		return &[]ClusterProfile{}, nil
	}

	// pull secret isn't returned by the API so pass the existing value along
	pullSecret := ""
	if len(config.ClusterProfile) != 0 {
		pullSecret = config.ClusterProfile[0].PullSecret
	}

	fipsEnabled := false
	if profile.FipsValidatedModules != nil {
		fipsEnabled = *profile.FipsValidatedModules == openshiftclusters.FipsValidatedModulesEnabled
	}

	resourceGroupId, err := commonids.ParseResourceGroupIDInsensitively(*profile.ResourceGroupId)
	if err != nil {
		return nil, err
	}

	resourceGroupIdString := ""
	resourceGroupName := ""
	if resourceGroupId != nil {
		resourceGroupIdString = resourceGroupId.ID()
		resourceGroupName = resourceGroupId.ResourceGroupName
	}

	return &[]ClusterProfile{
		{
			PullSecret:               pullSecret,
			Domain:                   pointer.From(profile.Domain),
			FipsEnabled:              fipsEnabled,
			ResourceGroupId:          resourceGroupIdString,
			ManagedResourceGroupName: resourceGroupName,
			Version:                  pointer.From(profile.Version),
		},
	}, nil
}

func expandOpenshiftServicePrincipalProfile(input []ServicePrincipal) *openshiftclusters.ServicePrincipalProfile {
	if len(input) == 0 {
		return nil
	}

	return &openshiftclusters.ServicePrincipalProfile{
		ClientId:     pointer.To(input[0].ClientId),
		ClientSecret: pointer.To(input[0].ClientSecret),
	}
}

func flattenOpenShiftServicePrincipalProfile(profile *openshiftclusters.ServicePrincipalProfile, config RedHatOpenShiftClusterModel) []ServicePrincipal {
	if profile == nil {
		return []ServicePrincipal{}
	}

	// client secret isn't returned by the API so pass the existing value along
	clientSecret := ""
	if len(config.ServicePrincipal) != 0 {
		clientSecret = config.ServicePrincipal[0].ClientSecret
	}

	return []ServicePrincipal{
		{
			ClientId:     pointer.From(profile.ClientId),
			ClientSecret: clientSecret,
		},
	}
}

func expandOpenshiftNetworkProfile(input []NetworkProfile) *openshiftclusters.NetworkProfile {
	if len(input) == 0 {
		return nil
	}

	preconfiguredNSG := openshiftclusters.PreconfiguredNSGDisabled
	if input[0].PreconfiguredNetworkSecurityGroupEnabled {
		preconfiguredNSG = openshiftclusters.PreconfiguredNSGEnabled
	}

	return &openshiftclusters.NetworkProfile{
		OutboundType:        pointer.ToEnum[openshiftclusters.OutboundType](input[0].OutboundType),
		PodCidr:             pointer.To(input[0].PodCidr),
		ServiceCidr:         pointer.To(input[0].ServiceCidr),
		PreconfiguredNSG:    pointer.To(preconfiguredNSG),
		LoadBalancerProfile: expandOpenshiftLoadBalancerProfile(input[0].LoadBalancerProfile),
	}
}

func expandOpenshiftLoadBalancerProfile(input []LoadBalancerProfile) *openshiftclusters.LoadBalancerProfile {
	if len(input) == 0 {
		return nil
	}

	result := &openshiftclusters.LoadBalancerProfile{}
	if input[0].ManagedOutboundIpCount > 0 {
		result.ManagedOutboundIPs = &openshiftclusters.ManagedOutboundIPs{
			Count: pointer.To(input[0].ManagedOutboundIpCount),
		}
	}
	return result
}

func flattenOpenShiftNetworkProfile(profile *openshiftclusters.NetworkProfile) []NetworkProfile {
	if profile == nil {
		return []NetworkProfile{}
	}

	preconfiguredNetworkSecurityGroupEnabled := false
	if profile.PreconfiguredNSG != nil {
		preconfiguredNetworkSecurityGroupEnabled = *profile.PreconfiguredNSG == openshiftclusters.PreconfiguredNSGEnabled
	}

	return []NetworkProfile{
		{
			OutboundType:                             string(pointer.From(profile.OutboundType)),
			PodCidr:                                  pointer.From(profile.PodCidr),
			ServiceCidr:                              pointer.From(profile.ServiceCidr),
			PreconfiguredNetworkSecurityGroupEnabled: preconfiguredNetworkSecurityGroupEnabled,
			LoadBalancerProfile:                      flattenOpenShiftLoadBalancerProfile(profile.LoadBalancerProfile),
		},
	}
}

func flattenOpenShiftLoadBalancerProfile(profile *openshiftclusters.LoadBalancerProfile) []LoadBalancerProfile {
	if profile == nil {
		return []LoadBalancerProfile{}
	}

	var count int64
	if profile.ManagedOutboundIPs != nil {
		count = pointer.From(profile.ManagedOutboundIPs.Count)
	}

	effectiveIps := make([]string, 0)
	if profile.EffectiveOutboundIPs != nil {
		for _, ip := range *profile.EffectiveOutboundIPs {
			effectiveIps = append(effectiveIps, pointer.From(ip.Id))
		}
	}

	return []LoadBalancerProfile{
		{
			ManagedOutboundIpCount: count,
			EffectiveOutboundIps:   effectiveIps,
		},
	}
}

func expandOpenshiftPlatformWorkloadIdentityProfile(input []PlatformWorkloadIdentityProfile) *openshiftclusters.PlatformWorkloadIdentityProfile {
	if len(input) == 0 {
		return nil
	}

	identities := make(map[string]openshiftclusters.PlatformWorkloadIdentity, len(input[0].PlatformWorkloadIdentity))
	for _, item := range input[0].PlatformWorkloadIdentity {
		identities[item.Name] = openshiftclusters.PlatformWorkloadIdentity{
			ResourceId: pointer.To(item.IdentityId),
		}
	}

	result := &openshiftclusters.PlatformWorkloadIdentityProfile{
		PlatformWorkloadIdentities: &identities,
	}

	if input[0].UpgradeableTo != "" {
		result.UpgradeableTo = pointer.To(input[0].UpgradeableTo)
	}

	return result
}

func flattenOpenShiftPlatformWorkloadIdentityProfile(profile *openshiftclusters.PlatformWorkloadIdentityProfile) []PlatformWorkloadIdentityProfile {
	if profile == nil {
		return []PlatformWorkloadIdentityProfile{}
	}

	identities := make([]PlatformWorkloadIdentity, 0)
	if profile.PlatformWorkloadIdentities != nil {
		for name, item := range *profile.PlatformWorkloadIdentities {
			resourceId := pointer.From(item.ResourceId)
			if parsed, err := commonids.ParseUserAssignedIdentityIDInsensitively(resourceId); err == nil {
				resourceId = parsed.ID()
			}
			identities = append(identities, PlatformWorkloadIdentity{
				Name:       name,
				IdentityId: resourceId,
				ClientId:   pointer.From(item.ClientId),
				ObjectId:   pointer.From(item.ObjectId),
			})
		}
	}

	return []PlatformWorkloadIdentityProfile{
		{
			UpgradeableTo:            pointer.From(profile.UpgradeableTo),
			PlatformWorkloadIdentity: identities,
		},
	}
}

func expandOpenshiftMainProfile(input []MainProfile) *openshiftclusters.MasterProfile {
	if len(input) == 0 {
		return nil
	}

	encryptionAtHost := openshiftclusters.EncryptionAtHostDisabled
	if input[0].EncryptionAtHostEnabled {
		encryptionAtHost = openshiftclusters.EncryptionAtHostEnabled
	}

	return &openshiftclusters.MasterProfile{
		VMSize:              pointer.To(input[0].VmSize),
		SubnetId:            pointer.To(input[0].SubnetId),
		EncryptionAtHost:    pointer.To(encryptionAtHost),
		DiskEncryptionSetId: pointer.To(input[0].DiskEncryptionSetId),
	}
}

func flattenOpenShiftMainProfile(profile *openshiftclusters.MasterProfile) []MainProfile {
	if profile == nil {
		return []MainProfile{}
	}

	encryptionAtHostEnabled := false
	if profile.EncryptionAtHost != nil {
		encryptionAtHostEnabled = *profile.EncryptionAtHost == openshiftclusters.EncryptionAtHostEnabled
	}

	return []MainProfile{
		{
			VmSize:                  pointer.From(profile.VMSize),
			SubnetId:                pointer.From(profile.SubnetId),
			EncryptionAtHostEnabled: encryptionAtHostEnabled,
			DiskEncryptionSetId:     pointer.From(profile.DiskEncryptionSetId),
		},
	}
}

func expandOpenshiftWorkerProfiles(input []WorkerProfile) *[]openshiftclusters.WorkerProfile {
	if len(input) == 0 {
		return nil
	}

	profiles := make([]openshiftclusters.WorkerProfile, 0, 1)

	encryptionAtHost := openshiftclusters.EncryptionAtHostDisabled
	if input[0].EncryptionAtHostEnabled {
		encryptionAtHost = openshiftclusters.EncryptionAtHostEnabled
	}

	profile := openshiftclusters.WorkerProfile{
		Name:                pointer.To("worker"),
		VMSize:              pointer.To(input[0].VmSize),
		DiskSizeGB:          pointer.To(input[0].DiskSizeGb),
		SubnetId:            pointer.To(input[0].SubnetId),
		Count:               pointer.To(input[0].NodeCount),
		EncryptionAtHost:    pointer.To(encryptionAtHost),
		DiskEncryptionSetId: pointer.To(input[0].DiskEncryptionSetId),
	}

	profiles = append(profiles, profile)

	return &profiles
}

func flattenOpenShiftWorkerProfiles(profiles *[]openshiftclusters.WorkerProfile) ([]WorkerProfile, error) {
	if profiles == nil || len(*profiles) == 0 {
		return []WorkerProfile{}, nil
	}

	rawProfiles := *profiles
	profile := rawProfiles[0]

	encryptionAtHostEnabled := false
	if profile.EncryptionAtHost != nil {
		encryptionAtHostEnabled = *profile.EncryptionAtHost == openshiftclusters.EncryptionAtHostEnabled
	}

	subnetIdString := ""
	if profile.SubnetId != nil {
		subnetId, err := commonids.ParseSubnetIDInsensitively(*profile.SubnetId)
		if err != nil {
			return []WorkerProfile{}, fmt.Errorf("parsing subnet id: %+v", err)
		}
		subnetIdString = subnetId.ID()
	}

	return []WorkerProfile{
		{
			NodeCount:               pointer.From(profile.Count),
			VmSize:                  pointer.From(profile.VMSize),
			DiskSizeGb:              pointer.From(profile.DiskSizeGB),
			SubnetId:                subnetIdString,
			EncryptionAtHostEnabled: encryptionAtHostEnabled,
			DiskEncryptionSetId:     pointer.From(profile.DiskEncryptionSetId),
		},
	}, nil
}

func expandOpenshiftApiServerProfile(input []ApiServerProfile) *openshiftclusters.APIServerProfile {
	if len(input) == 0 {
		return nil
	}

	visibility := openshiftclusters.Visibility(input[0].Visibility)

	return &openshiftclusters.APIServerProfile{
		Visibility: &visibility,
	}
}

func flattenOpenShiftAPIServerProfile(profile *openshiftclusters.APIServerProfile) []ApiServerProfile {
	if profile == nil {
		return []ApiServerProfile{}
	}

	return []ApiServerProfile{
		{
			Visibility: string(pointer.From(profile.Visibility)),
			Url:        pointer.From(profile.Url),
			IpAddress:  pointer.From(profile.IP),
		},
	}
}

func expandOpenshiftIngressProfiles(input []IngressProfile) *[]openshiftclusters.IngressProfile {
	if len(input) == 0 {
		return nil
	}

	profiles := make([]openshiftclusters.IngressProfile, 0, 1)

	profile := openshiftclusters.IngressProfile{
		Name:       pointer.To("default"),
		Visibility: pointer.To(openshiftclusters.Visibility(input[0].Visibility)),
	}

	profiles = append(profiles, profile)

	return &profiles
}

func flattenOpenShiftIngressProfiles(profiles *[]openshiftclusters.IngressProfile) []IngressProfile {
	if profiles == nil {
		return []IngressProfile{}
	}

	results := make([]IngressProfile, 0)

	for _, profile := range *profiles {
		results = append(results, IngressProfile{
			Visibility: string(pointer.From(profile.Visibility)),
			IpAddress:  pointer.From(profile.IP),
			Name:       pointer.From(profile.Name),
		})
	}

	return results
}
