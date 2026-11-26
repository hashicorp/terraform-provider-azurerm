// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redhat_openshift_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2025-07-25/openshiftclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

func expandOpenshiftServicePrincipalProfile(input []ServicePrincipal) *openshiftclusters.ServicePrincipalProfile {
	if len(input) == 0 {
		return nil
	}

	return &openshiftclusters.ServicePrincipalProfile{
		ClientId:     pointer.To(input[0].ClientId),
		ClientSecret: pointer.To(input[0].ClientSecret),
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

func expandOpenshiftApiServerProfile(input []ApiServerProfile) *openshiftclusters.APIServerProfile {
	if len(input) == 0 {
		return nil
	}

	return &openshiftclusters.APIServerProfile{
		Visibility: pointer.ToEnum[openshiftclusters.Visibility](input[0].Visibility),
	}
}

func expandOpenshiftIngressProfiles(input []IngressProfile) *[]openshiftclusters.IngressProfile {
	if len(input) == 0 {
		return nil
	}

	profiles := make([]openshiftclusters.IngressProfile, 0, 1)

	profile := openshiftclusters.IngressProfile{
		Name:       pointer.To("default"),
		Visibility: pointer.ToEnum[openshiftclusters.Visibility](input[0].Visibility),
	}

	profiles = append(profiles, profile)

	return &profiles
}
