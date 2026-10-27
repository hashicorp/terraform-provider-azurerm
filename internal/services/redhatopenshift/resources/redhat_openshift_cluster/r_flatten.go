// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redhat_openshift_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2025-07-25/openshiftclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) flatten(metadata sdk.ResourceMetaData, id openshiftclusters.OpenShiftClusterId, model *openshiftclusters.OpenShiftCluster, config RedHatOpenShiftClusterModel) error {
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

			platformWorkloadIdentityProfile, err := flattenOpenShiftPlatformWorkloadIdentityProfile(props.PlatformWorkloadIdentityProfile)
			if err != nil {
				return fmt.Errorf("flattening `platform_workload_identity_profile`: %+v", err)
			}
			state.PlatformWorkloadIdentityProfile = platformWorkloadIdentityProfile

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
			OutboundType:                             pointer.FromEnum(profile.OutboundType),
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

func flattenOpenShiftPlatformWorkloadIdentityProfile(profile *openshiftclusters.PlatformWorkloadIdentityProfile) ([]PlatformWorkloadIdentityProfile, error) {
	if profile == nil {
		return []PlatformWorkloadIdentityProfile{}, nil
	}

	identities := make([]PlatformWorkloadIdentity, 0)
	if profile.PlatformWorkloadIdentities != nil {
		for name, item := range *profile.PlatformWorkloadIdentities {
			resourceId := pointer.From(item.ResourceId)
			parsed, err := commonids.ParseUserAssignedIdentityIDInsensitively(resourceId)
			if err != nil {
				return nil, err
			}
			identities = append(identities, PlatformWorkloadIdentity{
				Name:       name,
				IdentityId: parsed.ID(),
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
	}, nil
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

func flattenOpenShiftAPIServerProfile(profile *openshiftclusters.APIServerProfile) []ApiServerProfile {
	if profile == nil {
		return []ApiServerProfile{}
	}

	return []ApiServerProfile{
		{
			Visibility: pointer.FromEnum(profile.Visibility),
			Url:        pointer.From(profile.Url),
			IpAddress:  pointer.From(profile.IP),
		},
	}
}

func flattenOpenShiftIngressProfiles(profiles *[]openshiftclusters.IngressProfile) []IngressProfile {
	if profiles == nil {
		return []IngressProfile{}
	}

	results := make([]IngressProfile, 0)

	for _, profile := range *profiles {
		results = append(results, IngressProfile{
			Visibility: pointer.FromEnum(profile.Visibility),
			IpAddress:  pointer.From(profile.IP),
			Name:       pointer.From(profile.Name),
		})
	}

	return results
}
