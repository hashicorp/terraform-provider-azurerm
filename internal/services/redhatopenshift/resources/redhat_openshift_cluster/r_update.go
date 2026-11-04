// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redhat_openshift_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2025-07-25/openshiftclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
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
			parameters.Properties.ConsoleProfile = nil
			parameters.Properties.WorkerProfilesStatus = nil
			if parameters.Properties.ClusterProfile != nil {
				parameters.Properties.ClusterProfile.OidcIssuer = nil
			}
			if parameters.Properties.ApiserverProfile != nil {
				parameters.Properties.ApiserverProfile.Url = nil
				parameters.Properties.ApiserverProfile.IP = nil
			}
			if parameters.Properties.IngressProfiles != nil {
				for i := range *parameters.Properties.IngressProfiles {
					(*parameters.Properties.IngressProfiles)[i].IP = nil
				}
			}
			if parameters.Properties.NetworkProfile != nil && parameters.Properties.NetworkProfile.LoadBalancerProfile != nil {
				parameters.Properties.NetworkProfile.LoadBalancerProfile.EffectiveOutboundIPs = nil
			}
			if parameters.Properties.PlatformWorkloadIdentityProfile != nil && parameters.Properties.PlatformWorkloadIdentityProfile.PlatformWorkloadIdentities != nil {
				for name, workloadIdentity := range *parameters.Properties.PlatformWorkloadIdentityProfile.PlatformWorkloadIdentities {
					workloadIdentity.ClientId = nil
					workloadIdentity.ObjectId = nil
					(*parameters.Properties.PlatformWorkloadIdentityProfile.PlatformWorkloadIdentities)[name] = workloadIdentity
				}
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandUserAssignedMapFromModel(state.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				parameters.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("service_principal") {
				parameters.Properties.ServicePrincipalProfile = expandOpenshiftServicePrincipalProfile(state.ServicePrincipal)
			}

			if metadata.ResourceData.HasChange("platform_workload_identity_profile") {
				parameters.Properties.PlatformWorkloadIdentityProfile = expandOpenshiftPlatformWorkloadIdentityProfile(state.PlatformWorkloadIdentityProfile)
			}

			if metadata.ResourceData.HasChange("network_profile") {
				parameters.Properties.NetworkProfile = expandOpenshiftNetworkProfile(state.NetworkProfile)
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
