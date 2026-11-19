// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package hdinsight_spark_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func hdinsightSparkClusterSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": helpers.SchemaHDInsightName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"cluster_version": helpers.SchemaHDInsightClusterVersion(),

		"tier": helpers.SchemaHDInsightTier(),

		"tls_min_version": helpers.SchemaHDInsightTls(),

		"encryption_in_transit_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"disk_encryption": helpers.SchemaHDInsightsDiskEncryptionProperties(),

		"component_version": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"spark": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},
				},
			},
		},

		"compute_isolation": helpers.SchemaHDInsightsComputeIsolation(),

		"gateway": helpers.SchemaHDInsightsGateway(),

		"metastores": helpers.SchemaHDInsightsExternalMetastores(),

		"network": helpers.SchemaHDInsightsNetwork(),

		"security_profile": helpers.SchemaHDInsightsSecurityProfile(),

		"storage_account": helpers.SchemaHDInsightsStorageAccounts(),

		"storage_account_gen2": helpers.SchemaHDInsightsGen2StorageAccounts(),

		"private_link_configuration": helpers.SchemaHDInsightPrivateLinkConfigurations(),

		"roles": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"head_node": helpers.SchemaHDInsightNodeDefinition("roles.0.head_node", hdInsightSparkClusterHeadNodeDefinition, true),

					"worker_node": helpers.SchemaHDInsightNodeDefinition("roles.0.worker_node", hdInsightSparkClusterWorkerNodeDefinition, true),

					"zookeeper_node": helpers.SchemaHDInsightNodeDefinition("roles.0.zookeeper_node", hdInsightSparkClusterZookeeperNodeDefinition, true),
				},
			},
		},

		"tags": commonschema.Tags(),

		"https_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"ssh_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"monitor": helpers.SchemaHDInsightsMonitor(),

		"extension": helpers.SchemaHDInsightsExtension(),

		"zones": commonschema.ZonesMultipleOptionalForceNew(),
	}
}
