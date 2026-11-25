// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_automatic_cluster

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/kubernetes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(ctx context.Context, metadata sdk.ResourceMetaData, id *commonids.KubernetesClusterId, model *managedclusters.ManagedCluster, includeResource bool) error {
	client := metadata.Client.Containers.KubernetesClustersClient

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

func flattenKubernetesAutomaticClusterHostedSystemProfile(profile *managedclusters.ManagedClusterHostedSystemProfile) []HostedSystemProfile {
	if profile == nil {
		return []HostedSystemProfile{}
	}

	return []HostedSystemProfile{{
		NodeSubnetID:       pointer.From(profile.NodeSubnetID),
		SystemNodeSubnetID: pointer.From(profile.SystemNodeSubnetID),
	}}
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
		proxyRedirectMechanism = pointer.FromEnum(profile.Istio.Components.ProxyRedirectionMechanism)
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
			if helpers.Base64IsEncoded(*kubeConfigRaw) {
				rawConfig = helpers.Base64Decode(*kubeConfigRaw)
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
			ClientCertificate:    user.ClientCertificateData,
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
