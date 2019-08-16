package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/hdinsight/mgmt/2018-06-01-preview/hdinsight"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// NOTE: this isn't a recommended way of building resources in Terraform
// this pattern is used to work around a generic but pedantic API endpoint
var hdInsightKafkaClusterHeadNodeDefinition = azure.HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         2,
	MaxInstanceCount:         2,
	CanSpecifyDisks:          false,
	FixedTargetInstanceCount: utils.Int32(int32(2)),
}

var hdInsightKafkaClusterWorkerNodeDefinition = azure.HDInsightNodeDefinition{
	CanSpecifyInstanceCount: true,
	MinInstanceCount:        1,
	MaxInstanceCount:        57,
	CanSpecifyDisks:         true,
	MaxNumberOfDisksPerNode: utils.Int(8),
}

var hdInsightKafkaClusterZookeeperNodeDefinition = azure.HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         3,
	MaxInstanceCount:         3,
	CanSpecifyDisks:          false,
	FixedTargetInstanceCount: utils.Int32(int32(3)),
}

func resourceArmHDInsightKafkaCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmHDInsightKafkaClusterCreate,
		Read:   resourceArmHDInsightKafkaClusterRead,
		Update: hdinsightClusterUpdate("Kafka", resourceArmHDInsightKafkaClusterRead),
		Delete: hdinsightClusterDelete("Kafka"),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": azure.SchemaHDInsightName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"cluster_version": azure.SchemaHDInsightClusterVersion(),

			"tier": azure.SchemaHDInsightTier(),

			"component_version": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kafka": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"gateway": azure.SchemaHDInsightsGateway(),

			"storage_account": azure.SchemaHDInsightsStorageAccounts(),

			"roles": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"head_node": azure.SchemaHDInsightNodeDefinition("roles.0.head_node", hdInsightKafkaClusterHeadNodeDefinition),

						"worker_node": azure.SchemaHDInsightNodeDefinition("roles.0.worker_node", hdInsightKafkaClusterWorkerNodeDefinition),

						"zookeeper_node": azure.SchemaHDInsightNodeDefinition("roles.0.zookeeper_node", hdInsightKafkaClusterZookeeperNodeDefinition),
					},
				},
			},

			"tags": tagsSchema(),

			"https_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ssh_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmHDInsightKafkaClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).hdinsight.ClustersClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	clusterVersion := d.Get("cluster_version").(string)
	tags := d.Get("tags").(map[string]interface{})
	tier := hdinsight.Tier(d.Get("tier").(string))

	componentVersionsRaw := d.Get("component_version").([]interface{})
	componentVersions := expandHDInsightKafkaComponentVersion(componentVersionsRaw)

	gatewayRaw := d.Get("gateway").([]interface{})
	gateway := azure.ExpandHDInsightsConfigurations(gatewayRaw)

	storageAccountsRaw := d.Get("storage_account").([]interface{})
	storageAccounts, err := azure.ExpandHDInsightsStorageAccounts(storageAccountsRaw)
	if err != nil {
		return fmt.Errorf("Error expanding `storage_account`: %s", err)
	}

	kafkaRoles := hdInsightRoleDefinition{
		HeadNodeDef:      hdInsightKafkaClusterHeadNodeDefinition,
		WorkerNodeDef:    hdInsightKafkaClusterWorkerNodeDefinition,
		ZookeeperNodeDef: hdInsightKafkaClusterZookeeperNodeDefinition,
	}
	rolesRaw := d.Get("roles").([]interface{})
	roles, err := expandHDInsightRoles(rolesRaw, kafkaRoles)
	if err != nil {
		return fmt.Errorf("Error expanding `roles`: %+v", err)
	}

	if requireResourcesToBeImported {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing HDInsight Kafka Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_hdinsight_kafka_cluster", *existing.ID)
		}
	}

	params := hdinsight.ClusterCreateParametersExtended{
		Location: utils.String(location),
		Properties: &hdinsight.ClusterCreateProperties{
			Tier:           tier,
			OsType:         hdinsight.Linux,
			ClusterVersion: utils.String(clusterVersion),
			ClusterDefinition: &hdinsight.ClusterDefinition{
				Kind:             utils.String("Kafka"),
				ComponentVersion: componentVersions,
				Configurations:   gateway,
			},
			StorageProfile: &hdinsight.StorageProfile{
				Storageaccounts: storageAccounts,
			},
			ComputeProfile: &hdinsight.ComputeProfile{
				Roles: roles,
			},
		},
		Tags: expandTags(tags),
	}
	future, err := client.Create(ctx, resourceGroup, name, params)
	if err != nil {
		return fmt.Errorf("Error creating HDInsight Kafka Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of HDInsight Kafka Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving HDInsight Kafka Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Error reading ID for HDInsight Kafka Cluster %q (Resource Group %q)", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmHDInsightKafkaClusterRead(d, meta)
}

func resourceArmHDInsightKafkaClusterRead(d *schema.ResourceData, meta interface{}) error {
	clustersClient := meta.(*ArmClient).hdinsight.ClustersClient
	configurationsClient := meta.(*ArmClient).hdinsight.ConfigurationsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["clusters"]

	resp, err := clustersClient.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] HDInsight Kafka Cluster %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving HDInsight Kafka Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	configuration, err := configurationsClient.Get(ctx, resourceGroup, name, "gateway")
	if err != nil {
		return fmt.Errorf("Error retrieving Configuration for HDInsight Kafka Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	// storage_account isn't returned so I guess we just leave it ¯\_(ツ)_/¯
	if props := resp.Properties; props != nil {
		d.Set("cluster_version", props.ClusterVersion)
		d.Set("tier", string(props.Tier))

		if def := props.ClusterDefinition; def != nil {
			if err := d.Set("component_version", flattenHDInsightKafkaComponentVersion(def.ComponentVersion)); err != nil {
				return fmt.Errorf("Error flattening `component_version`: %+v", err)
			}

			if err := d.Set("gateway", azure.FlattenHDInsightsConfigurations(configuration.Value)); err != nil {
				return fmt.Errorf("Error flattening `gateway`: %+v", err)
			}
		}

		kafkaRoles := hdInsightRoleDefinition{
			HeadNodeDef:      hdInsightKafkaClusterHeadNodeDefinition,
			WorkerNodeDef:    hdInsightKafkaClusterWorkerNodeDefinition,
			ZookeeperNodeDef: hdInsightKafkaClusterZookeeperNodeDefinition,
		}
		flattenedRoles := flattenHDInsightRoles(d, props.ComputeProfile, kafkaRoles)
		if err := d.Set("roles", flattenedRoles); err != nil {
			return fmt.Errorf("Error flattening `roles`: %+v", err)
		}

		httpEndpoint := azure.FindHDInsightConnectivityEndpoint("HTTPS", props.ConnectivityEndpoints)
		d.Set("https_endpoint", httpEndpoint)
		sshEndpoint := azure.FindHDInsightConnectivityEndpoint("SSH", props.ConnectivityEndpoints)
		d.Set("ssh_endpoint", sshEndpoint)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func expandHDInsightKafkaComponentVersion(input []interface{}) map[string]*string {
	vs := input[0].(map[string]interface{})
	return map[string]*string{
		"kafka": utils.String(vs["kafka"].(string)),
	}
}

func flattenHDInsightKafkaComponentVersion(input map[string]*string) []interface{} {
	kafkaVersion := ""
	if v, ok := input["kafka"]; ok {
		if v != nil {
			kafkaVersion = *v
		}
	}
	return []interface{}{
		map[string]interface{}{
			"kafka": kafkaVersion,
		},
	}
}
