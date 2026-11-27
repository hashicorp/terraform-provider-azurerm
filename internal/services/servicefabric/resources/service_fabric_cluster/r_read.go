// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package service_fabric_cluster

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabric/2021-06-01/cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	servicefabricClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicefabric/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceFabricClusterRead(d *pluginsdk.ResourceData, meta any) error {
	client := servicefabricClient.From(meta).ClustersClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := cluster.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found - removing from state!", id.ID())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id.ID(), err)
	}

	d.Set("name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("cluster_code_version", props.ClusterCodeVersion)
			d.Set("cluster_endpoint", props.ClusterEndpoint)
			d.Set("management_endpoint", props.ManagementEndpoint)
			d.Set("vm_image", props.VmImage)

			reliabilityLevel := ""
			if props.ReliabilityLevel != nil {
				reliabilityLevel = string(*props.ReliabilityLevel)
			}
			d.Set("reliability_level", reliabilityLevel)

			upgradeMode := ""
			if props.UpgradeMode != nil {
				upgradeMode = string(*props.UpgradeMode)
			}
			d.Set("upgrade_mode", upgradeMode)

			sfZonalMode := ""
			if props.SfZonalUpgradeMode != nil {
				sfZonalMode = string(*props.SfZonalUpgradeMode)
			}
			d.Set("service_fabric_zonal_upgrade_mode", sfZonalMode)

			vmssZonalMode := ""
			if props.VMSSZonalUpgradeMode != nil {
				vmssZonalMode = string(*props.VMSSZonalUpgradeMode)
			}
			d.Set("vmss_zonal_upgrade_mode", vmssZonalMode)

			addOnFeatures := flattenServiceFabricClusterAddOnFeatures(props.AddOnFeatures)
			if err := d.Set("add_on_features", pluginsdk.NewSet(pluginsdk.HashString, addOnFeatures)); err != nil {
				return fmt.Errorf("setting `add_on_features`: %+v", err)
			}

			if err := d.Set("azure_active_directory", flattenServiceFabricClusterAzureActiveDirectory(props.AzureActiveDirectory)); err != nil {
				return fmt.Errorf("setting `azure_active_directory`: %+v", err)
			}

			if err := d.Set("certificate", flattenServiceFabricClusterCertificate(props.Certificate)); err != nil {
				return fmt.Errorf("setting `certificate`: %+v", err)
			}

			if err := d.Set("certificate_common_names", flattenServiceFabricClusterCertificateCommonNames(props.CertificateCommonNames)); err != nil {
				return fmt.Errorf("setting `certificate_common_names`: %+v", err)
			}

			if err := d.Set("reverse_proxy_certificate", flattenServiceFabricClusterReverseProxyCertificate(props.ReverseProxyCertificate)); err != nil {
				return fmt.Errorf("setting `reverse_proxy_certificate`: %+v", err)
			}

			if err := d.Set("reverse_proxy_certificate_common_names", flattenServiceFabricClusterCertificateCommonNames(props.ReverseProxyCertificateCommonNames)); err != nil {
				return fmt.Errorf("setting `reverse_proxy_certificate_common_names`: %+v", err)
			}

			if err := d.Set("client_certificate_thumbprint", flattenServiceFabricClusterClientCertificateThumbprints(props.ClientCertificateThumbprints)); err != nil {
				return fmt.Errorf("setting `client_certificate_thumbprint`: %+v", err)
			}

			if err := d.Set("client_certificate_common_name", flattenServiceFabricClusterClientCertificateCommonNames(props.ClientCertificateCommonNames)); err != nil {
				return fmt.Errorf("setting `client_certificate_common_name`: %+v", err)
			}

			if err := d.Set("diagnostics_config", flattenServiceFabricClusterDiagnosticsConfig(props.DiagnosticsStorageAccountConfig)); err != nil {
				return fmt.Errorf("setting `diagnostics_config`: %+v", err)
			}

			if err := d.Set("upgrade_policy", flattenServiceFabricClusterUpgradePolicy(props.UpgradeDescription)); err != nil {
				return fmt.Errorf("setting `upgrade_policy`: %+v", err)
			}

			if err := d.Set("fabric_settings", flattenServiceFabricClusterFabricSettings(props.FabricSettings)); err != nil {
				return fmt.Errorf("setting `fabric_settings`: %+v", err)
			}

			if err := d.Set("node_type", flattenServiceFabricClusterNodeTypes(props.NodeTypes)); err != nil {
				return fmt.Errorf("setting `node_type`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
