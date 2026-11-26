// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package hdinsight_hadoop_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func hdinsightHadoopClusterSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": helpers.SchemaHDInsightName("HDInsight Hadoop Cluster"),

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"cluster_version": helpers.SchemaHDInsightClusterVersion(),

		"tier": helpers.SchemaHDInsightTier("HDInsight Hadoop Cluster"),

		"tls_min_version": helpers.SchemaHDInsightTls(),

		"component_version": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"hadoop": {
						Description: "The version of Hadoop which should be used for this HDInsight Hadoop Cluster.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ForceNew:    true,
					},
				},
			},
		},

		"disk_encryption": helpers.SchemaHDInsightsDiskEncryptionProperties(),

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
					"head_node": helpers.SchemaHDInsightNodeDefinition("roles.0.head_node", hdInsightHadoopClusterHeadNodeDefinition, true),

					"worker_node": helpers.SchemaHDInsightNodeDefinition("roles.0.worker_node", hdInsightHadoopClusterWorkerNodeDefinition, true),

					"zookeeper_node": helpers.SchemaHDInsightNodeDefinition("roles.0.zookeeper_node", hdInsightHadoopClusterZookeeperNodeDefinition, true),

					"edge_node": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"target_instance_count": {
									Description:  "The number of instances which should be run for the Worker Nodes.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 25),
								},

								"vm_size": {
									Description:  "The Size of the Virtual Machine which should be used as the Edge Nodes.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(validate.NodeDefinitionVMSize, false),
								},

								"install_script_action": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"name": {
												Description:  "The name of the install script action.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
											"uri": {
												Description:  "The URI pointing to the script to run during the installation of the edge node.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
											"parameters": {
												Description:  "The parameters for the script.",
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"https_endpoints": SchemaHDInsightsHttpsEndpoints(),

								"uninstall_script_actions": helpers.SchemaHDInsightsScriptActions(),
							},
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"https_endpoint": {
			Description: "The HTTPS Connectivity Endpoint for this HDInsight Hadoop Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"ssh_endpoint": {
			Description: "The SSH Connectivity Endpoint for this HDInsight Hadoop Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"monitor": helpers.SchemaHDInsightsMonitor(),

		"extension": helpers.SchemaHDInsightsExtension(),
	}
}
