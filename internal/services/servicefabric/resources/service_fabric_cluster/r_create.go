// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package service_fabric_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabric/2021-06-01/cluster"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceFabricClusterCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ServiceFabric.ClustersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := cluster.NewClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_service_fabric_cluster", id.ID())
			}
		}
	}

	addOnFeaturesRaw := d.Get("add_on_features").(*pluginsdk.Set).List()
	addOnFeatures := expandServiceFabricClusterAddOnFeatures(addOnFeaturesRaw)

	azureActiveDirectoryRaw := d.Get("azure_active_directory").([]any)
	azureActiveDirectory := expandServiceFabricClusterAzureActiveDirectory(azureActiveDirectoryRaw)

	diagnosticsRaw := d.Get("diagnostics_config").([]any)
	diagnostics := expandServiceFabricClusterDiagnosticsConfig(diagnosticsRaw)

	upgradePolicyRaw := d.Get("upgrade_policy").([]any)
	upgradePolicy := expandServiceFabricClusterUpgradePolicy(upgradePolicyRaw)

	fabricSettingsRaw := d.Get("fabric_settings").([]any)
	fabricSettings := expandServiceFabricClusterFabricSettings(fabricSettingsRaw)

	nodeTypesRaw := d.Get("node_type").([]any)
	nodeTypes := expandServiceFabricClusterNodeTypes(d, nodeTypesRaw)

	location := d.Get("location").(string)
	reliabilityLevel := cluster.ReliabilityLevel(d.Get("reliability_level").(string))
	managementEndpoint := d.Get("management_endpoint").(string)
	upgradeMode := cluster.UpgradeMode(d.Get("upgrade_mode").(string))
	clusterCodeVersion := d.Get("cluster_code_version").(string)
	vmImage := d.Get("vm_image").(string)
	t := d.Get("tags").(map[string]any)

	clusterModel := cluster.Cluster{
		Location: location,
		Tags:     tags.Expand(t),
		Properties: &cluster.ClusterProperties{
			AddOnFeatures:                      addOnFeatures,
			AzureActiveDirectory:               azureActiveDirectory,
			CertificateCommonNames:             expandServiceFabricClusterCertificateCommonNames(d),
			ReverseProxyCertificateCommonNames: expandServiceFabricClusterReverseProxyCertificateCommonNames(d),
			DiagnosticsStorageAccountConfig:    diagnostics,
			FabricSettings:                     fabricSettings,
			ManagementEndpoint:                 managementEndpoint,
			NodeTypes:                          nodeTypes,
			ReliabilityLevel:                   &reliabilityLevel,
			UpgradeDescription:                 upgradePolicy,
			UpgradeMode:                        &upgradeMode,
			VmImage:                            pointer.To(vmImage),
		},
	}

	if sfZonalUpgradeMode, ok := d.GetOk("service_fabric_zonal_upgrade_mode"); ok {
		clusterModel.Properties.SfZonalUpgradeMode = pointer.ToEnum[cluster.SfZonalUpgradeMode](sfZonalUpgradeMode.(string))
	}

	if vmssZonalUpgradeMode, ok := d.GetOk("vmss_zonal_upgrade_mode"); ok {
		clusterModel.Properties.VMSSZonalUpgradeMode = pointer.ToEnum[cluster.VMSSZonalUpgradeMode](vmssZonalUpgradeMode.(string))
	}

	if certificateRaw, ok := d.GetOk("certificate"); ok {
		clusterModel.Properties.Certificate = expandServiceFabricClusterCertificate(certificateRaw.([]any))
	}

	if reverseProxyCertificateRaw, ok := d.GetOk("reverse_proxy_certificate"); ok {
		clusterModel.Properties.ReverseProxyCertificate = expandServiceFabricClusterReverseProxyCertificate(reverseProxyCertificateRaw.([]any))
	}

	if clientCertificateThumbprintRaw, ok := d.GetOk("client_certificate_thumbprint"); ok {
		clusterModel.Properties.ClientCertificateThumbprints = expandServiceFabricClusterClientCertificateThumbprints(d, clientCertificateThumbprintRaw.([]any))
	}

	if clientCertificateCommonNamesRaw, ok := d.GetOk("client_certificate_common_name"); ok {
		clusterModel.Properties.ClientCertificateCommonNames = expandServiceFabricClusterClientCertificateCommonNames(d, clientCertificateCommonNamesRaw.([]any))
	}

	if clusterCodeVersion != "" {
		clusterModel.Properties.ClusterCodeVersion = pointer.To(clusterCodeVersion)
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, clusterModel, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, clusterModel); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}
	return resourceServiceFabricClusterRead(d, meta)
}
