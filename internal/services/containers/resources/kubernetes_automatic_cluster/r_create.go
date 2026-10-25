// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_automatic_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clusterClient := metadata.Client.Containers.KubernetesClustersClient
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

func expandKubernetesAutomaticClusterPrivateCluster(model []PrivateClusterModel) (bool, bool, string) {
	if len(model) == 0 {
		return false, false, ""
	}

	config := model[0]
	return true, config.PrivateClusterPublicFQDNEnabled, config.PrivateDNSZoneID
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
				ProxyRedirectionMechanism: pointer.ToEnum[managedclusters.ProxyRedirectionMechanism](config.ProxyRedirectMechanism),
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
