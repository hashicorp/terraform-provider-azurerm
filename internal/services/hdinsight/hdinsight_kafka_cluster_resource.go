// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hdinsight

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
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
	FixedTargetInstanceCount: pointer.To(int64(2)),
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
	FixedTargetInstanceCount: pointer.To(int64(3)),
}

var hdInsightKafkaClusterKafkaManagementNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         2,
	CanSpecifyDisks:          false,
	FixedTargetInstanceCount: pointer.To(int64(2)),
}

func resourceHDInsightKafkaCluster() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceHDInsightKafkaClusterCreate,
		Read:   resourceHDInsightKafkaClusterRead,
		Update: hdinsightClusterUpdate("Kafka", resourceHDInsightKafkaClusterRead),
		Delete: hdinsightClusterDelete("Kafka"),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseHDInsightClusterID(id)
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

			"private_link_configuration": SchemaHDInsightPrivateLinkConfigurations(),

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

						"kafka_management_node": SchemaHDInsightNodeDefinitionKafka("roles.0.kafka_management_node", hdInsightKafkaClusterKafkaManagementNodeDefinition, false),
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

			"monitor": SchemaHDInsightsMonitor(),

			"extension": SchemaHDInsightsExtension(),
		},
	}

	if !features.FourPointOhBeta() {
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
		}

		resource.Schema["roles"].Elem.(*pluginsdk.Resource).Schema["kafka_management_node"].Elem.(*pluginsdk.Resource).Schema["username"].Deprecated = "`username` will become Computed only in version 4.0 of the AzureRM Provider as the service auto-generates a value for this property"
	}

	return resource
}

func resourceHDInsightKafkaClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HDInsight.Clusters
	extensionsClient := meta.(*clients.Client).HDInsight.Extensions

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewHDInsightClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := location.Normalize(d.Get("location").(string))
	clusterVersion := d.Get("cluster_version").(string)
	t := d.Get("tags").(map[string]interface{})
	tier := clusters.Tier(d.Get("tier").(string))
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
	storageAccounts, expandedIdentity, err := ExpandHDInsightsStorageAccounts(storageAccountsRaw, storageAccountsGen2Raw)
	if err != nil {
		return fmt.Errorf("failure expanding `storage_account`: %s", err)
	}

	networkPropertiesRaw := d.Get("network").([]interface{})
	networkProperties := ExpandHDInsightsNetwork(networkPropertiesRaw)

	privateLinkConfigurationsRaw := d.Get("private_link_configuration").([]interface{})
	privateLinkConfigurations := ExpandHDInsightPrivateLinkConfigurations(privateLinkConfigurationsRaw)

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

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for the presence of an existing Kafka %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_hdinsight_kafka_cluster", id.ID())
	}

	kafkaRestProperty := expandKafkaRestProxyProperty(d.Get("rest_proxy").([]interface{}))

	computeIsolationProperties := ExpandHDInsightComputeIsolationProperties(d.Get("compute_isolation").([]interface{}))

	var configurationsRaw interface{} = configurations
	payload := clusters.ClusterCreateParametersExtended{
		Location: utils.String(location),
		Properties: &clusters.ClusterCreateProperties{
			Tier:                      pointer.To(tier),
			OsType:                    pointer.To(clusters.OSTypeLinux),
			ClusterVersion:            utils.String(clusterVersion),
			MinSupportedTlsVersion:    utils.String(tls),
			NetworkProperties:         networkProperties,
			PrivateLinkConfigurations: privateLinkConfigurations,
			ClusterDefinition: &clusters.ClusterDefinition{
				Kind:             pointer.To(clusters.ClusterKindKafka),
				ComponentVersion: pointer.To(componentVersions),
				Configurations:   pointer.To(configurationsRaw),
			},
			StorageProfile: &clusters.StorageProfile{
				Storageaccounts: storageAccounts,
			},
			ComputeProfile: &clusters.ComputeProfile{
				Roles: roles,
			},
			KafkaRestProperties:        kafkaRestProperty,
			ComputeIsolationProperties: computeIsolationProperties,
		},
		Tags:     tags.Expand(t),
		Identity: expandedIdentity,
	}

	if encryptionInTransit, ok := d.GetOk("encryption_in_transit_enabled"); ok {
		payload.Properties.EncryptionInTransitProperties = &clusters.EncryptionInTransitProperties{
			IsEncryptionInTransitEnabled: utils.Bool(encryptionInTransit.(bool)),
		}
	}

	if diskEncryptionPropertiesRaw, ok := d.GetOk("disk_encryption"); ok {
		payload.Properties.DiskEncryptionProperties, err = ExpandHDInsightsDiskEncryptionProperties(diskEncryptionPropertiesRaw.([]interface{}))
		if err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("security_profile"); ok {
		payload.Properties.SecurityProfile = ExpandHDInsightSecurityProfile(v.([]interface{}))

		// @tombuildsstuff: this behaviour is likely wrong and wants reevaluating - users should need to explicitly define this in the config?
		payload.Identity = &identity.SystemAndUserAssignedMap{
			Type:        identity.TypeUserAssigned,
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		if payload.Properties.SecurityProfile != nil && payload.Properties.SecurityProfile.MsiResourceId != nil {
			payload.Identity.IdentityIds[*payload.Properties.SecurityProfile.MsiResourceId] = identity.UserAssignedIdentityDetails{
				// intentionally empty
			}
		}
	}

	if err := client.CreateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating Kafka %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// We can only enable monitoring after creation
	if v, ok := d.GetOk("monitor"); ok {
		monitorRaw := v.([]interface{})
		if err := enableHDInsightMonitoring(ctx, extensionsClient, id, monitorRaw); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("extension"); ok {
		extensionRaw := v.([]interface{})
		if err := enableHDInsightAzureMonitor(ctx, extensionsClient, id, extensionRaw); err != nil {
			return err
		}
	}

	return resourceHDInsightKafkaClusterRead(d, meta)
}

func resourceHDInsightKafkaClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	clustersClient := meta.(*clients.Client).HDInsight.Clusters
	configurationsClient := meta.(*clients.Client).HDInsight.Configurations
	extensionsClient := meta.(*clients.Client).HDInsight.Extensions
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseHDInsightClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := clustersClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Kafka %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Kafka %s: %+v", id, err)
	}

	// Each call to configurationsClient methods is HTTP request. Getting all settings in one operation
	configurationsResp, err := configurationsClient.List(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Configurations for Kafka %s: %+v", id, err)
	}

	configurations := make(map[string]map[string]string)
	if model := configurationsResp.Model; model != nil && model.Configurations != nil {
		configurations = *model.Configurations
	}
	gateway, exists := configurations["gateway"]
	if !exists {
		return fmt.Errorf("retrieving Gateway Configuration Kafka %s: %+v", id, err)
	}

	monitor, err := extensionsClient.GetMonitoringStatus(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Monitoring Status for Kafka %s: %+v", id, err)
	}

	extension, err := extensionsClient.GetAzureMonitorStatus(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Azure Monitor Status for Kafka %s: %+v", id, err)
	}

	d.Set("name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {

		d.Set("location", location.Normalize(model.Location))

		// storage_account isn't returned so I guess we just leave it ¯\_(ツ)_/¯
		if props := model.Properties; props != nil {
			d.Set("tier", string(pointer.From(props.Tier)))
			d.Set("cluster_version", props.ClusterVersion)
			d.Set("tls_min_version", props.MinSupportedTlsVersion)

			if err := d.Set("component_version", flattenHDInsightKafkaComponentVersion(props.ClusterDefinition.ComponentVersion)); err != nil {
				return fmt.Errorf("failure flattening `component_version`: %+v", err)
			}

			if err := d.Set("gateway", FlattenHDInsightsConfigurations(gateway, d)); err != nil {
				return fmt.Errorf("failure flattening `gateway`: %+v", err)
			}

			flattenHDInsightsMetastores(d, configurations)

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

			httpEndpoint := findHDInsightConnectivityEndpoint("HTTPS", props.ConnectivityEndpoints)
			d.Set("https_endpoint", httpEndpoint)
			sshEndpoint := findHDInsightConnectivityEndpoint("SSH", props.ConnectivityEndpoints)
			d.Set("ssh_endpoint", sshEndpoint)
			kafkaRestProxyEndpoint := findHDInsightConnectivityEndpoint("KafkaRestProxyPublicEndpoint", props.ConnectivityEndpoints)
			d.Set("kafka_rest_proxy_endpoint", kafkaRestProxyEndpoint)

			if props.EncryptionInTransitProperties != nil {
				d.Set("encryption_in_transit_enabled", props.EncryptionInTransitProperties.IsEncryptionInTransitEnabled)
			}

			if err := d.Set("network", flattenHDInsightsNetwork(props.NetworkProperties)); err != nil {
				return fmt.Errorf("flatten `network`: %+v", err)
			}
			if err := d.Set("private_link_configuration", flattenHDInsightPrivateLinkConfigurations(props.PrivateLinkConfigurations)); err != nil {
				return fmt.Errorf("flattening `private_link_configuration`: %+v", err)
			}
			if err := d.Set("compute_isolation", flattenHDInsightComputeIsolationProperties(props.ComputeIsolationProperties)); err != nil {
				return fmt.Errorf("failed setting `compute_isolation`: %+v", err)
			}

			diskEncryptionProps, err := flattenHDInsightsDiskEncryptionProperties(props.DiskEncryptionProperties)
			if err != nil {
				return err
			}
			if err := d.Set("disk_encryption", diskEncryptionProps); err != nil {
				return fmt.Errorf("flattening `disk_encryption`: %+v", err)
			}

			d.Set("monitor", flattenHDInsightMonitoring(monitor.Model))

			if err = d.Set("rest_proxy", flattenKafkaRestProxyProperty(props.KafkaRestProperties)); err != nil {
				return fmt.Errorf("setting `rest_proxy`: %+v", err)
			}

			d.Set("extension", flattenHDInsightAzureMonitor(extension.Model))

			if err := d.Set("security_profile", flattenHDInsightSecurityProfile(props.SecurityProfile, d)); err != nil {
				return fmt.Errorf("setting `security_profile`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}

func expandHDInsightKafkaComponentVersion(input []interface{}) map[string]string {
	kafkaVersion := ""
	if len(input) > 0 && input[0] != nil {
		vs := input[0].(map[string]interface{})
		kafkaVersion = vs["kafka"].(string)
	}
	return map[string]string{
		"kafka": kafkaVersion,
	}
}

func flattenHDInsightKafkaComponentVersion(input *map[string]string) []interface{} {
	output := make([]interface{}, 0)
	if input != nil {
		kafkaVersion := ""
		if v, ok := (*input)["kafka"]; ok {
			kafkaVersion = v
		}
		output = append(output, map[string]interface{}{
			"kafka": kafkaVersion,
		})
	}
	return output
}

func expandKafkaRestProxyProperty(input []interface{}) *clusters.KafkaRestProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	groupId := raw["security_group_id"].(string)
	groupName := raw["security_group_name"].(string)

	return &clusters.KafkaRestProperties{
		ClientGroupInfo: &clusters.ClientGroupInfo{
			GroupId:   &groupId,
			GroupName: &groupName,
		},
	}
}

func flattenKafkaRestProxyProperty(input *clusters.KafkaRestProperties) []interface{} {
	if input == nil || input.ClientGroupInfo == nil {
		return []interface{}{}
	}

	groupInfo := input.ClientGroupInfo

	groupId := ""
	if groupInfo.GroupId != nil {
		groupId = *groupInfo.GroupId
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
