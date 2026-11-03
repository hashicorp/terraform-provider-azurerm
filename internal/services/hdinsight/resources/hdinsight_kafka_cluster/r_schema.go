// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package hdinsight_kafka_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func hdinsightKafkaClusterSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": helpers.SchemaHDInsightName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"cluster_version": helpers.SchemaHDInsightClusterVersion(),

		"tier": helpers.SchemaHDInsightTier(),

		"tls_min_version": helpers.SchemaHDInsightTls(),

		"metastores": helpers.SchemaHDInsightsExternalMetastores(),

		"network": helpers.SchemaHDInsightsNetwork(),

		"component_version": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"kafka": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},
				},
			},
		},

		"gateway": helpers.SchemaHDInsightsGateway(),

		"security_profile": helpers.SchemaHDInsightsSecurityProfile(),

		"storage_account": helpers.SchemaHDInsightsStorageAccounts(),

		"storage_account_gen2": helpers.SchemaHDInsightsGen2StorageAccounts(),

		"private_link_configuration": helpers.SchemaHDInsightPrivateLinkConfigurations(),

		"compute_isolation": helpers.SchemaHDInsightsComputeIsolation(),

		"encryption_in_transit_enabled": {
			Type:     pluginsdk.TypeBool,
			ForceNew: true,
			Optional: true,
		},

		"disk_encryption": helpers.SchemaHDInsightsDiskEncryptionProperties(),

		"roles": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"head_node": helpers.SchemaHDInsightNodeDefinition("roles.0.head_node", hdInsightKafkaClusterHeadNodeDefinition, true),

					"worker_node": helpers.SchemaHDInsightNodeDefinition("roles.0.worker_node", hdInsightKafkaClusterWorkerNodeDefinition, true),

					"zookeeper_node": helpers.SchemaHDInsightNodeDefinition("roles.0.zookeeper_node", hdInsightKafkaClusterZookeeperNodeDefinition, true),

					"kafka_management_node": helpers.SchemaHDInsightNodeDefinitionKafka("roles.0.kafka_management_node", hdInsightKafkaClusterKafkaManagementNodeDefinition, false),
				},
			},
		},

		"rest_proxy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"security_group_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsUUID,
					},

					"security_group_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
			RequiredWith: []string{"roles.0.kafka_management_node"},
		},

		"tags": commonschema.Tags(),

		"https_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kafka_rest_proxy_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"ssh_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"monitor": helpers.SchemaHDInsightsMonitor(),

		"extension": helpers.SchemaHDInsightsExtension(),
	}
}
