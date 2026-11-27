// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package service_fabric_managed_cluster

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2024-04-01/managedcluster"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2024-04-01/nodetype"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicefabricmanaged/client"
)

func (k Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ClusterResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			clusterClient := client.From(metadata.Client).ManagedClusterClient
			nodeTypeClient := client.From(metadata.Client).NodeTypeClient

			subscriptionId := account.From(metadata.Client).SubscriptionId

			managedClusterId := managedcluster.NewManagedClusterID(subscriptionId, model.ResourceGroup, model.Name)
			cluster := managedcluster.ManagedCluster{
				Location:   model.Location,
				Name:       pointer.To(model.Name),
				Properties: expandClusterProperties(&model),
				Sku:        managedcluster.Sku{Name: model.Sku},
			}

			tagsMap := make(map[string]string)
			for k, v := range model.Tags {
				tagsMap[k] = v.(string)
			}
			cluster.Tags = &tagsMap

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := clusterClient.Get(ctx, managedClusterId)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("while checking if cluster %q already exists: %+v", managedClusterId.String(), err)
					}
				} else {
					return metadata.ResourceRequiresImport("azurerm_service_fabric_managed_cluster", managedClusterId)
				}
			}

			if err := clusterClient.CreateOrUpdateCallbackThenPoll(ctx, managedClusterId, cluster, metadata.SetIDCallback(&managedClusterId)); err != nil {
				return fmt.Errorf("creating %s: %+v", managedClusterId, err)
			}
			metadata.SetID(managedClusterId)

			toDelete := make([]string, 0)
			if metadata.ResourceData.HasChange("node_type") {
				o, n := metadata.ResourceData.GetChange("node_type")
				ont := o.([]any)
				nnt := n.([]any)

				for _, on := range ont {
					oldNodeType := on.(map[string]any)
					oldId := oldNodeType["name"].(string)
					found := false
					for _, nt := range nnt {
						newNodeType := nt.(map[string]any)
						newId := newNodeType["name"].(string)
						if oldId == newId {
							found = true
							break
						}
					}

					if !found {
						toDelete = append(toDelete, oldId)
					}
				}
			}

			// Delete the old Node Types
			for _, nodeType := range toDelete {
				nodeTypeId := nodetype.NewNodeTypeID(subscriptionId, model.ResourceGroup, model.Name, nodeType)
				if err := nodeTypeClient.DeleteThenPoll(ctx, nodeTypeId); err != nil {
					return fmt.Errorf("deleting %s: %+v", nodeTypeId, err)
				}
			}

			// Then provision the Node Types
			for _, nt := range model.NodeTypes {
				nodeTypeProperties, err := expandNodeTypeProperties(&nt)
				if err != nil {
					return fmt.Errorf("while expanding node type %q: %+v", nt.Name, err)
				}

				nodeTypeId := nodetype.NewNodeTypeID(subscriptionId, model.ResourceGroup, model.Name, nt.Name)
				payload := nodetype.NodeType{
					Properties: nodeTypeProperties,
				}
				if err := nodeTypeClient.CreateOrUpdateThenPoll(ctx, nodeTypeId, payload); err != nil {
					return fmt.Errorf("adding %s: %+v", nodeTypeId, err)
				}
			}

			return nil
		},

		Timeout: 90 * time.Minute,
	}
}

func expandClusterProperties(model *ClusterResourceModel) *managedcluster.ManagedClusterProperties {
	out := &managedcluster.ManagedClusterProperties{}

	addons := make([]managedcluster.ManagedClusterAddOnFeature, 0)
	if model.DNSService {
		addons = append(addons, managedcluster.ManagedClusterAddOnFeatureDnsService)
	}
	if model.BackupRestoreService {
		addons = append(addons, managedcluster.ManagedClusterAddOnFeatureBackupRestoreService)
	}
	out.AddonFeatures = &addons

	out.AdminPassword = pointer.To(model.Password)
	out.AdminUserName = model.Username

	out.DnsName = model.Name
	if model.DNSName != "" && model.DNSName != model.Name {
		out.DnsName = model.DNSName
	}

	if v := model.SubnetId; v != "" {
		out.SubnetId = pointer.To(v)
	}

	if auth := model.Authentication; len(auth) > 0 {
		if adAuth := auth[0].ADAuth; len(adAuth) > 0 {
			if adAuth[0].ClientApp != "" && adAuth[0].ClusterApp != "" && adAuth[0].TenantId != "" {
				out.AzureActiveDirectory = &managedcluster.AzureActiveDirectory{
					ClientApplication:  pointer.To(adAuth[0].ClientApp),
					ClusterApplication: pointer.To(adAuth[0].ClusterApp),
					TenantId:           pointer.To(adAuth[0].TenantId),
				}
			}
		}
		if certs := auth[0].CertAuthentication; len(certs) > 0 {
			clients := make([]managedcluster.ClientCertificate, len(certs))
			for idx, cert := range certs {
				clients[idx] = managedcluster.ClientCertificate{
					CommonName: pointer.To(cert.CommonName),
					IsAdmin:    cert.CertificateType == CertTypeAdmin,
					Thumbprint: pointer.To(cert.Thumbprint),
				}
			}
			out.Clients = &clients
		}
	}

	out.ClientConnectionPort = &model.ClientConnectionPort
	out.ClusterUpgradeCadence = &model.UpgradeWave

	if customSettings := model.CustomFabricSettings; len(customSettings) > 0 {
		// First we build a map of all settings per section
		fsMap := make(map[string][]managedcluster.SettingsParameterDescription)
		for _, cs := range customSettings {
			spd := managedcluster.SettingsParameterDescription{
				Name:  cs.Parameter,
				Value: cs.Value,
			}
			fsMap[cs.Section] = append(fsMap[cs.Section], spd)
		}

		// Then we update the properties struct
		fs := make([]managedcluster.SettingsSectionDescription, 0)
		for k, v := range fsMap {
			fs = append(fs, managedcluster.SettingsSectionDescription{
				Name:       k,
				Parameters: v,
			})
		}
		out.FabricSettings = &fs
	}

	out.HTTPGatewayConnectionPort = &model.HTTPGatewayPort

	if rules := model.LBRules; len(rules) > 0 {
		lbRules := make([]managedcluster.LoadBalancingRule, len(rules))
		nsRules := make([]managedcluster.NetworkSecurityRule, len(rules))

		for idx, rule := range rules {
			lbRules[idx] = managedcluster.LoadBalancingRule{
				BackendPort:      rule.BackendPort,
				FrontendPort:     rule.FrontendPort,
				ProbeProtocol:    rule.ProbeProtocol,
				ProbeRequestPath: pointer.To(rule.ProbeRequestPath),
				Protocol:         rule.Protocol,
			}

			fePortStr := strconv.FormatInt(rule.FrontendPort, 10)
			var sgProto managedcluster.NsgProtocol
			switch rule.Protocol {
			case managedcluster.ProtocolTcp:
				sgProto = managedcluster.NsgProtocolTcp
			case managedcluster.ProtocolUdp:
				sgProto = managedcluster.NsgProtocolUdp
			}
			nsRules[idx] = managedcluster.NetworkSecurityRule{
				Access:                     managedcluster.AccessAllow,
				SourceAddressPrefixes:      &[]string{"0.0.0.0/0"},
				SourcePortRanges:           &[]string{"1-65535"},
				DestinationPortRanges:      &[]string{fePortStr},
				DestinationAddressPrefixes: &[]string{"0.0.0.0/0"},
				Direction:                  managedcluster.DirectionInbound,
				Name:                       fmt.Sprintf("rule%d-allow-fe", rule.FrontendPort),
				Priority:                   1000,
				Protocol:                   sgProto,
			}
		}
		out.LoadBalancingRules = &lbRules
		out.NetworkSecurityRules = &nsRules
	}
	return out
}

func expandNodeTypeProperties(nt *NodeType) (*nodetype.NodeTypeProperties, error) {
	vmSecrets := make([]nodetype.VaultSecretGroup, len(nt.VmSecrets))
	for idx, secret := range nt.VmSecrets {
		vcs := make([]nodetype.VaultCertificate, len(secret.Certificates))
		for cidx, cert := range secret.Certificates {
			vcs[cidx] = nodetype.VaultCertificate{
				CertificateStore: cert.Store,
				CertificateURL:   cert.Url,
			}
		}
		vmSecrets[idx] = nodetype.VaultSecretGroup{
			SourceVault:       nodetype.SubResource{Id: &secret.SourceVault},
			VaultCertificates: vcs,
		}
	}

	appFrom, appTo, err := parsePortRange(nt.ApplicationPorts)
	if err != nil {
		return nil, fmt.Errorf("while parsing application port range (%q): %+v", nt.ApplicationPorts, err)
	}

	ephemeralFrom, ephemeralTo, err := parsePortRange(nt.EphemeralPorts)
	if err != nil {
		return nil, fmt.Errorf("while parsing ephemeral port range (%q): %+v", nt.EphemeralPorts, err)
	}
	nodeTypeProperties := &nodetype.NodeTypeProperties{
		ApplicationPorts: &nodetype.EndpointRangeDescription{
			EndPort:   appTo,
			StartPort: appFrom,
		},
		Capacities:     &nt.Capacities,
		DataDiskSizeGB: &nt.DataDiskSize,
		DataDiskType:   &nt.DataDiskType,
		EphemeralPorts: &nodetype.EndpointRangeDescription{
			EndPort:   ephemeralTo,
			StartPort: ephemeralFrom,
		},
		IsPrimary:               nt.Primary,
		IsStateless:             &nt.Stateless,
		MultiplePlacementGroups: &nt.MultiplePlacementGroupsEnabled,
		PlacementProperties:     &nt.PlacementProperties,
		VMImageOffer:            &nt.VmImageOffer,
		VMImagePublisher:        &nt.VmImagePublisher,
		VMImageSku:              &nt.VmImageSku,
		VMImageVersion:          &nt.VmImageVersion,
		VMInstanceCount:         nt.VmInstanceCount,
		VMSecrets:               &vmSecrets,
		VMSize:                  &nt.VmSize,
	}

	return nodeTypeProperties, nil
}
