package hdinsight

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/azure-sdk-for-go/services/hdinsight/mgmt/2018-06-01/hdinsight"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hdinsight/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// NOTE: this isn't a recommended way of building resources in Terraform
// this pattern is used to work around a generic but pedantic API endpoint
var hdInsightKafkaClusterHeadNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         2,
	MaxInstanceCount:         utils.Int(2),
	CanSpecifyDisks:          false,
	FixedTargetInstanceCount: utils.Int32(int32(2)),
}

var hdInsightKafkaClusterWorkerNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount: true,
	MinInstanceCount:        1,
	CanSpecifyDisks:         true,
	MaxNumberOfDisksPerNode: utils.Int(8),
}

var hdInsightKafkaClusterZookeeperNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         3,
	MaxInstanceCount:         utils.Int(3),
	CanSpecifyDisks:          false,
	FixedTargetInstanceCount: utils.Int32(int32(3)),
}

var hdInsightKafkaClusterKafkaManagementNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         2,
	CanSpecifyDisks:          false,
	FixedTargetInstanceCount: utils.Int32(int32(2)),
}

func resourceHDInsightKafkaCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHDInsightKafkaClusterCreate,
		Read:   resourceHDInsightKafkaClusterRead,
		Update: hdinsightClusterUpdate("Kafka", resourceHDInsightKafkaClusterRead),
		Delete: hdinsightClusterDelete("Kafka"),
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": SchemaHDInsightName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"cluster_version": SchemaHDInsightClusterVersion(),

			"tier": SchemaHDInsightTier(),

			"tls_min_version": SchemaHDInsightTls(),

			"metastores": SchemaHDInsightsExternalMetastores(),

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

			"gateway": SchemaHDInsightsGateway(),

			"storage_account": SchemaHDInsightsStorageAccounts(),

			"storage_account_gen2": SchemaHDInsightsGen2StorageAccounts(),

			"encryption_in_transit_enabled": {
				Type:     pluginsdk.TypeBool,
				ForceNew: true,
				Optional: true,
			},

			"roles": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"head_node": SchemaHDInsightNodeDefinition("roles.0.head_node", hdInsightKafkaClusterHeadNodeDefinition, true),

						"worker_node": SchemaHDInsightNodeDefinition("roles.0.worker_node", hdInsightKafkaClusterWorkerNodeDefinition, true),

						"zookeeper_node": SchemaHDInsightNodeDefinition("roles.0.zookeeper_node", hdInsightKafkaClusterZookeeperNodeDefinition, true),

						"kafka_management_node": SchemaHDInsightNodeDefinition("roles.0.kafka_management_node", hdInsightKafkaClusterKafkaManagementNodeDefinition, false),
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
					},
				},
				RequiredWith: []string{"roles.0.kafka_management_node"},
			},

			"tags": tags.Schema(),

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

			"monitor": SchemaHDInsightsMonitor(),
		},
	}
}

func resourceHDInsightKafkaClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	groupClient := meta.(*clients.Client).Authorization.GroupsClient
	client := meta.(*clients.Client).HDInsight.ClustersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	extensionsClient := meta.(*clients.Client).HDInsight.ExtensionsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	id := parse.NewClusterID(subscriptionId, resourceGroup, name)
	location := azure.NormalizeLocation(d.Get("location").(string))
	clusterVersion := d.Get("cluster_version").(string)
	t := d.Get("tags").(map[string]interface{})
	tier := hdinsight.Tier(d.Get("tier").(string))
	tls := d.Get("tls_min_version").(string)

	componentVersionsRaw := d.Get("component_version").([]interface{})
	componentVersions := expandHDInsightKafkaComponentVersion(componentVersionsRaw)

	gatewayRaw := d.Get("gateway").([]interface{})
	configurations := ExpandHDInsightsConfigurations(gatewayRaw)

	metastoresRaw := d.Get("metastores").([]interface{})
	metastores := expandHDInsightsMetastore(metastoresRaw)
	for k, v := range metastores {
		configurations[k] = v
	}

	storageAccountsRaw := d.Get("storage_account").([]interface{})
	storageAccountsGen2Raw := d.Get("storage_account_gen2").([]interface{})
	storageAccounts, identity, err := ExpandHDInsightsStorageAccounts(storageAccountsRaw, storageAccountsGen2Raw)
	if err != nil {
		return fmt.Errorf("failure expanding `storage_account`: %s", err)
	}

	kafkaRoles := hdInsightRoleDefinition{
		HeadNodeDef:            hdInsightKafkaClusterHeadNodeDefinition,
		WorkerNodeDef:          hdInsightKafkaClusterWorkerNodeDefinition,
		ZookeeperNodeDef:       hdInsightKafkaClusterZookeeperNodeDefinition,
		KafkaManagementNodeDef: &hdInsightKafkaClusterKafkaManagementNodeDefinition,
	}
	rolesRaw := d.Get("roles").([]interface{})
	roles, err := expandHDInsightRoles(rolesRaw, kafkaRoles)
	if err != nil {
		return fmt.Errorf("failure expanding `roles`: %+v", err)
	}

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("failure checking for presence of existing HDInsight Kafka Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_hdinsight_kafka_cluster", *existing.ID)
	}

	kafkaRestProperty, err := expandKafkaRestProxyProperty(ctx, groupClient, d.Get("rest_proxy").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding kafka rest proxy property")
	}

	params := hdinsight.ClusterCreateParametersExtended{
		Location: utils.String(location),
		Properties: &hdinsight.ClusterCreateProperties{
			Tier:                   tier,
			OsType:                 hdinsight.OSTypeLinux,
			ClusterVersion:         utils.String(clusterVersion),
			MinSupportedTLSVersion: utils.String(tls),
			ClusterDefinition: &hdinsight.ClusterDefinition{
				Kind:             utils.String("Kafka"),
				ComponentVersion: componentVersions,
				Configurations:   configurations,
			},
			StorageProfile: &hdinsight.StorageProfile{
				Storageaccounts: storageAccounts,
			},
			ComputeProfile: &hdinsight.ComputeProfile{
				Roles: roles,
			},
			KafkaRestProperties: kafkaRestProperty,
		},
		Tags:     tags.Expand(t),
		Identity: identity,
	}

	if encryptionInTransit, ok := d.GetOk("encryption_in_transit_enabled"); ok {
		params.Properties.EncryptionInTransitProperties = &hdinsight.EncryptionInTransitProperties{
			IsEncryptionInTransitEnabled: utils.Bool(encryptionInTransit.(bool)),
		}
	}

	future, err := client.Create(ctx, resourceGroup, name, params)
	if err != nil {
		return fmt.Errorf("failure creating HDInsight Kafka Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("failed waiting for creation of HDInsight Kafka Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("failure retrieving HDInsight Kafka Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("failure reading ID for HDInsight Kafka Cluster %q (Resource Group %q)", name, resourceGroup)
	}

	d.SetId(id.ID())

	// We can only enable monitoring after creation
	if v, ok := d.GetOk("monitor"); ok {
		monitorRaw := v.([]interface{})
		if err := enableHDInsightMonitoring(ctx, extensionsClient, resourceGroup, name, monitorRaw); err != nil {
			return err
		}
	}

	return resourceHDInsightKafkaClusterRead(d, meta)
}

func resourceHDInsightKafkaClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	clustersClient := meta.(*clients.Client).HDInsight.ClustersClient
	configurationsClient := meta.(*clients.Client).HDInsight.ConfigurationsClient
	extensionsClient := meta.(*clients.Client).HDInsight.ExtensionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Name

	resp, err := clustersClient.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] HDInsight Kafka Cluster %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("failure retrieving HDInsight Kafka Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	// Each call to configurationsClient methods is HTTP request. Getting all settings in one operation
	configurations, err := configurationsClient.List(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("failure retrieving Configuration for HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	gateway, exists := configurations.Configurations["gateway"]
	if !exists {
		return fmt.Errorf("failure retrieving gateway for HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
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
		d.Set("tls_min_version", props.MinSupportedTLSVersion)

		if def := props.ClusterDefinition; def != nil {
			if err := d.Set("component_version", flattenHDInsightKafkaComponentVersion(def.ComponentVersion)); err != nil {
				return fmt.Errorf("failure flattening `component_version`: %+v", err)
			}

			if err := d.Set("gateway", FlattenHDInsightsConfigurations(gateway)); err != nil {
				return fmt.Errorf("failure flattening `gateway`: %+v", err)
			}

			flattenHDInsightsMetastores(d, configurations.Configurations)
		}

		kafkaRoles := hdInsightRoleDefinition{
			HeadNodeDef:            hdInsightKafkaClusterHeadNodeDefinition,
			WorkerNodeDef:          hdInsightKafkaClusterWorkerNodeDefinition,
			ZookeeperNodeDef:       hdInsightKafkaClusterZookeeperNodeDefinition,
			KafkaManagementNodeDef: &hdInsightKafkaClusterKafkaManagementNodeDefinition,
		}
		flattenedRoles := flattenHDInsightRoles(d, props.ComputeProfile, kafkaRoles)
		if err := d.Set("roles", flattenedRoles); err != nil {
			return fmt.Errorf("failure flattening `roles`: %+v", err)
		}

		httpEndpoint := FindHDInsightConnectivityEndpoint("HTTPS", props.ConnectivityEndpoints)
		d.Set("https_endpoint", httpEndpoint)
		sshEndpoint := FindHDInsightConnectivityEndpoint("SSH", props.ConnectivityEndpoints)
		d.Set("ssh_endpoint", sshEndpoint)
		kafkaRestProxyEndpoint := FindHDInsightConnectivityEndpoint("KafkaRestProxyPublicEndpoint", props.ConnectivityEndpoints)
		d.Set("kafka_rest_proxy_endpoint", kafkaRestProxyEndpoint)

		if props.EncryptionInTransitProperties != nil {
			d.Set("encryption_in_transit_enabled", props.EncryptionInTransitProperties.IsEncryptionInTransitEnabled)
		}

		monitor, err := extensionsClient.GetMonitoringStatus(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("failed reading monitor configuration for HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		d.Set("monitor", flattenHDInsightMonitoring(monitor))
		if err := d.Set("rest_proxy", flattenKafkaRestProxyProperty(props.KafkaRestProperties)); err != nil {
			return fmt.Errorf(`failed setting "rest_proxy" for HDInsight Kafka Cluster %q (Resource Group %q): %+v`, name, resourceGroup, err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
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

func expandKafkaRestProxyProperty(ctx context.Context, client *graphrbac.GroupsClient, input []interface{}) (*hdinsight.KafkaRestProperties, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	raw := input[0].(map[string]interface{})
	groupId := raw["security_group_id"].(string)

	// Current API requires users further specify the "security_group_name" in the client group info of the kafka rest property,
	// which is unnecessary as user already specify the "security_group_id".
	// https://github.com/Azure/azure-rest-api-specs/issues/10667
	res, err := client.Get(ctx, groupId)
	if err != nil {
		return nil, fmt.Errorf("retrieving AAD gruop %s: %v", groupId, err)
	}

	return &hdinsight.KafkaRestProperties{
		ClientGroupInfo: &hdinsight.ClientGroupInfo{
			GroupID:   &groupId,
			GroupName: res.DisplayName,
		},
	}, nil
}

func flattenKafkaRestProxyProperty(input *hdinsight.KafkaRestProperties) []interface{} {
	if input == nil || input.ClientGroupInfo == nil {
		return []interface{}{}
	}

	groupInfo := input.ClientGroupInfo
	groupId := ""
	if groupInfo.GroupID != nil {
		groupId = *groupInfo.GroupID
	}

	return []interface{}{
		map[string]interface{}{
			"security_group_id": groupId,
		},
	}
}
