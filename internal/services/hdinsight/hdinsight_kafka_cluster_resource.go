// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hdinsight

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/hdinsight/mgmt/2018-06-01/hdinsight" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// NOTE: this isn't a recommended way of building resources in Terraform
// this pattern is used to work around a generic but pedantic API endpoint
var hdInsightKafkaClusterHeadNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         2,
	MaxInstanceCount:         pointer.To(2),
	CanSpecifyDisks:          false,
	FixedTargetInstanceCount: pointer.To(int32(2)),
}

var hdInsightKafkaClusterWorkerNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount: true,
	MinInstanceCount:        1,
	CanSpecifyDisks:         true,
	MaxNumberOfDisksPerNode: pointer.To(8),
}

var hdInsightKafkaClusterZookeeperNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         3,
	MaxInstanceCount:         pointer.To(3),
	CanSpecifyDisks:          false,
	FixedTargetInstanceCount: pointer.To(int32(3)),
}

var hdInsightKafkaClusterKafkaManagementNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         2,
	CanSpecifyDisks:          false,
	FixedTargetInstanceCount: pointer.To(int32(2)),
}

func resourceHDInsightKafkaCluster() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceHDInsightKafkaClusterCreate,
		Read:   resourceHDInsightKafkaClusterRead,
		Update: hdinsightClusterUpdate("Kafka", resourceHDInsightKafkaClusterRead),
		Delete: hdinsightClusterDelete("Kafka"),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ClusterID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": SchemaHDInsightName(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"cluster_version": SchemaHDInsightClusterVersion(),

			"tier": SchemaHDInsightTier(),

			"tls_min_version": SchemaHDInsightTls(),

			"metastores": SchemaHDInsightsExternalMetastores(),

			"network": SchemaHDInsightsNetwork(),

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

			"security_profile": SchemaHDInsightsSecurityProfile(),

			"storage_account": SchemaHDInsightsStorageAccounts(),

			"storage_account_gen2": SchemaHDInsightsGen2StorageAccounts(),

			"compute_isolation": SchemaHDInsightsComputeIsolation(),

			"encryption_in_transit_enabled": {
				Type:     pluginsdk.TypeBool,
				ForceNew: true,
				Optional: true,
			},

			"disk_encryption": SchemaHDInsightsDiskEncryptionProperties(),

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

						"security_group_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
				RequiredWith: func() []string {
					if !features.FourPointOh() {
						return []string{"roles.0.kafka_management_node"}
					}
					return []string{}
				}(),
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

			"extension": SchemaHDInsightsExtension(),
		},
	}

	if !features.FourPointOh() {
		resource.Schema["roles"] = &pluginsdk.Schema{
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
			Deprecated: "`kafka_management_node` will be removed in version 4.0 of the AzureRM Provider since it no longer support configurations from the user",
		}
	} else {
		resource.Schema["roles"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"head_node": SchemaHDInsightNodeDefinition("roles.0.head_node", hdInsightKafkaClusterHeadNodeDefinition, true),

					"worker_node": SchemaHDInsightNodeDefinition("roles.0.worker_node", hdInsightKafkaClusterWorkerNodeDefinition, true),

					"zookeeper_node": SchemaHDInsightNodeDefinition("roles.0.zookeeper_node", hdInsightKafkaClusterZookeeperNodeDefinition, true),
				},
			},
		}
	}

	return resource
}

func resourceHDInsightKafkaClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HDInsight.ClustersClient
	extensionsClient := meta.(*clients.Client).HDInsight.ExtensionsClient

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
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

	networkPropertiesRaw := d.Get("network").([]interface{})
	networkProperties := ExpandHDInsightsNetwork(networkPropertiesRaw)

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

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_hdinsight_kafka_cluster", id.ID())
	}

	kafkaRestProperty := expandKafkaRestProxyProperty(d.Get("rest_proxy").([]interface{}))

	computeIsolationProperties := ExpandHDInsightComputeIsolationProperties(d.Get("compute_isolation").([]interface{}))

	params := hdinsight.ClusterCreateParametersExtended{
		Location: utils.String(location),
		Properties: &hdinsight.ClusterCreateProperties{
			Tier:                   tier,
			OsType:                 hdinsight.OSTypeLinux,
			ClusterVersion:         utils.String(clusterVersion),
			MinSupportedTLSVersion: utils.String(tls),
			NetworkProperties:      networkProperties,
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
			KafkaRestProperties:        kafkaRestProperty,
			ComputeIsolationProperties: computeIsolationProperties,
		},
		Tags:     tags.Expand(t),
		Identity: identity,
	}

	if encryptionInTransit, ok := d.GetOk("encryption_in_transit_enabled"); ok {
		params.Properties.EncryptionInTransitProperties = &hdinsight.EncryptionInTransitProperties{
			IsEncryptionInTransitEnabled: utils.Bool(encryptionInTransit.(bool)),
		}
	}

	if diskEncryptionPropertiesRaw, ok := d.GetOk("disk_encryption"); ok {
		params.Properties.DiskEncryptionProperties, err = ExpandHDInsightsDiskEncryptionProperties(diskEncryptionPropertiesRaw.([]interface{}))
		if err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("security_profile"); ok {
		params.Properties.SecurityProfile = ExpandHDInsightSecurityProfile(v.([]interface{}))

		params.Identity = &hdinsight.ClusterIdentity{
			Type:                   hdinsight.ResourceIdentityTypeUserAssigned,
			UserAssignedIdentities: make(map[string]*hdinsight.ClusterIdentityUserAssignedIdentitiesValue),
		}

		if params.Properties.SecurityProfile != nil && params.Properties.SecurityProfile.MsiResourceID != nil {
			params.Identity.UserAssignedIdentities[*params.Properties.SecurityProfile.MsiResourceID] = &hdinsight.ClusterIdentityUserAssignedIdentitiesValue{}
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

	if v, ok := d.GetOk("extension"); ok {
		extensionRaw := v.([]interface{})
		if err := enableHDInsightAzureMonitor(ctx, extensionsClient, resourceGroup, name, extensionRaw); err != nil {
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
		tier := ""
		// the Azure API is inconsistent here, so rewrite this into the casing we expect
		for _, v := range hdinsight.PossibleTierValues() {
			if strings.EqualFold(string(v), string(props.Tier)) {
				tier = string(v)
			}
		}
		d.Set("tier", tier)
		d.Set("cluster_version", props.ClusterVersion)
		d.Set("tls_min_version", props.MinSupportedTLSVersion)

		if def := props.ClusterDefinition; def != nil {
			if err := d.Set("component_version", flattenHDInsightKafkaComponentVersion(def.ComponentVersion)); err != nil {
				return fmt.Errorf("failure flattening `component_version`: %+v", err)
			}

			if err := d.Set("gateway", FlattenHDInsightsConfigurations(gateway, d)); err != nil {
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

		if props.NetworkProperties != nil {
			if err := d.Set("network", FlattenHDInsightsNetwork(props.NetworkProperties)); err != nil {
				return fmt.Errorf("flatten `network`: %+v", err)
			}
		}
		if props.ComputeIsolationProperties.EnableComputeIsolation != nil {
			if err := d.Set("compute_isolation", FlattenHDInsightComputeIsolationProperties(*props.ComputeIsolationProperties)); err != nil {
				return fmt.Errorf("failed setting `compute_isolation`: %+v", err)
			}
		}

		if props.DiskEncryptionProperties != nil {
			diskEncryptionProps, err := FlattenHDInsightsDiskEncryptionProperties(*props.DiskEncryptionProperties)
			if err != nil {
				return err
			}
			if err := d.Set("disk_encryption", diskEncryptionProps); err != nil {
				return fmt.Errorf("flattening `disk_encryption`: %+v", err)
			}
		}

		monitor, err := extensionsClient.GetMonitoringStatus(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("failed reading monitor configuration for HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		d.Set("monitor", flattenHDInsightMonitoring(monitor))

		if err = d.Set("rest_proxy", flattenKafkaRestProxyProperty(props.KafkaRestProperties)); err != nil {
			return fmt.Errorf(`failed setting "rest_proxy" for HDInsight Kafka Cluster %q (Resource Group %q): %+v`, name, resourceGroup, err)
		}

		extension, err := extensionsClient.GetAzureMonitorStatus(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("reading extension configuration for HDInsight Hadoop Cluster %q (Resource Group %q) %+v", name, resourceGroup, err)
		}

		d.Set("extension", flattenHDInsightAzureMonitor(extension))

		if err := d.Set("security_profile", flattenHDInsightSecurityProfile(props.SecurityProfile, d)); err != nil {
			return fmt.Errorf("setting `security_profile`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func expandHDInsightKafkaComponentVersion(input []interface{}) map[string]*string {
	if len(input) == 0 || input[0] == nil {
		return map[string]*string{"kafka": utils.String("")}
	}
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

func expandKafkaRestProxyProperty(input []interface{}) *hdinsight.KafkaRestProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	groupId := raw["security_group_id"].(string)
	groupName := raw["security_group_name"].(string)

	return &hdinsight.KafkaRestProperties{
		ClientGroupInfo: &hdinsight.ClientGroupInfo{
			GroupID:   &groupId,
			GroupName: &groupName,
		},
	}
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

	groupName := ""
	if groupInfo.GroupName != nil {
		groupName = *groupInfo.GroupName
	}

	return []interface{}{
		map[string]interface{}{
			"security_group_id":   groupId,
			"security_group_name": groupName,
		},
	}
}
