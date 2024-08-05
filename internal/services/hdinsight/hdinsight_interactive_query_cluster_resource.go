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
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// NOTE: this isn't a recommended way of building resources in Terraform
// this pattern is used to work around a generic but pedantic API endpoint
var hdInsightInteractiveQueryClusterHeadNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         2,
	MaxInstanceCount:         pointer.To(2),
	CanSpecifyDisks:          false,
	FixedTargetInstanceCount: pointer.To(int64(2)),
}

var hdInsightInteractiveQueryClusterWorkerNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:                  true,
	MinInstanceCount:                         1,
	CanSpecifyDisks:                          false,
	CanAutoScaleByCapacityDeprecated4PointOh: true,
	CanAutoScaleOnSchedule:                   true,
}

var hdInsightInteractiveQueryClusterZookeeperNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         3,
	MaxInstanceCount:         pointer.To(3),
	CanSpecifyDisks:          false,
	FixedTargetInstanceCount: pointer.To(int64(3)),
}

func resourceHDInsightInteractiveQueryCluster() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceHDInsightInteractiveQueryClusterCreate,
		Read:   resourceHDInsightInteractiveQueryClusterRead,
		Update: hdinsightClusterUpdate("Interactive Query", resourceHDInsightInteractiveQueryClusterRead),
		Delete: hdinsightClusterDelete("Interactive Query"),

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

			"encryption_in_transit_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"disk_encryption": SchemaHDInsightsDiskEncryptionProperties(),

			"component_version": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"interactive_hive": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"compute_isolation": SchemaHDInsightsComputeIsolation(),

			"gateway": SchemaHDInsightsGateway(),

			"metastores": SchemaHDInsightsExternalMetastores(),

			"network": SchemaHDInsightsNetwork(),

			"security_profile": SchemaHDInsightsSecurityProfile(),

			"storage_account": SchemaHDInsightsStorageAccounts(),

			"storage_account_gen2": SchemaHDInsightsGen2StorageAccounts(),

			"private_link_configuration": SchemaHDInsightPrivateLinkConfigurations(),

			"tags": commonschema.Tags(),

			"https_endpoint": {
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
					"head_node": SchemaHDInsightNodeDefinition("roles.0.head_node", hdInsightInteractiveQueryClusterHeadNodeDefinition, true),

					"worker_node": SchemaHDInsightNodeDefinition("roles.0.worker_node", hdInsightInteractiveQueryClusterWorkerNodeDefinition, true),

					"zookeeper_node": SchemaHDInsightNodeDefinition("roles.0.zookeeper_node", hdInsightInteractiveQueryClusterZookeeperNodeDefinition, true),
				},
			},
		}
	} else {
		hdInsightInteractiveQueryClusterWorkerNodeDefinition.CanAutoScaleByCapacityDeprecated4PointOh = false
		resource.Schema["roles"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"head_node": SchemaHDInsightNodeDefinition("roles.0.head_node", hdInsightInteractiveQueryClusterHeadNodeDefinition, true),

					"worker_node": SchemaHDInsightNodeDefinition("roles.0.worker_node", hdInsightInteractiveQueryClusterWorkerNodeDefinition, true),

					"zookeeper_node": SchemaHDInsightNodeDefinition("roles.0.zookeeper_node", hdInsightInteractiveQueryClusterZookeeperNodeDefinition, true),
				},
			},
		}
	}

	return resource
}

func resourceHDInsightInteractiveQueryClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HDInsight.Clusters
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	extensionsClient := meta.(*clients.Client).HDInsight.Extensions
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	id := commonids.NewHDInsightClusterID(subscriptionId, resourceGroup, name)
	location := location.Normalize(d.Get("location").(string))
	clusterVersion := d.Get("cluster_version").(string)
	t := d.Get("tags").(map[string]interface{})
	tier := clusters.Tier(d.Get("tier").(string))
	tls := d.Get("tls_min_version").(string)

	componentVersionsRaw := d.Get("component_version").([]interface{})
	componentVersions := expandHDInsightInteractiveQueryComponentVersion(componentVersionsRaw)

	gatewayRaw := d.Get("gateway").([]interface{})
	configurations := ExpandHDInsightsConfigurations(gatewayRaw)

	metastoresRaw := d.Get("metastores").([]interface{})
	metastores := expandHDInsightsMetastore(metastoresRaw)
	for k, v := range metastores {
		configurations[k] = v
	}

	networkPropertiesRaw := d.Get("network").([]interface{})
	networkProperties := ExpandHDInsightsNetwork(networkPropertiesRaw)

	privateLinkConfigurationsRaw := d.Get("private_link_configuration").([]interface{})
	privateLinkConfigurations := ExpandHDInsightPrivateLinkConfigurations(privateLinkConfigurationsRaw)

	storageAccountsRaw := d.Get("storage_account").([]interface{})
	storageAccountsGen2Raw := d.Get("storage_account_gen2").([]interface{})
	storageAccounts, expandedIdentity, err := ExpandHDInsightsStorageAccounts(storageAccountsRaw, storageAccountsGen2Raw)
	if err != nil {
		return fmt.Errorf("expanding `storage_account`: %s", err)
	}

	interactiveQueryRoles := hdInsightRoleDefinition{
		HeadNodeDef:      hdInsightInteractiveQueryClusterHeadNodeDefinition,
		WorkerNodeDef:    hdInsightInteractiveQueryClusterWorkerNodeDefinition,
		ZookeeperNodeDef: hdInsightInteractiveQueryClusterZookeeperNodeDefinition,
	}
	rolesRaw := d.Get("roles").([]interface{})
	roles, err := expandHDInsightRoles(rolesRaw, interactiveQueryRoles)
	if err != nil {
		return fmt.Errorf("expanding `roles`: %+v", err)
	}

	computeIsolationProperties := ExpandHDInsightComputeIsolationProperties(d.Get("compute_isolation").([]interface{}))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing HDInsight InteractiveQuery Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_hdinsight_interactive_query_cluster", id.ID())
	}

	encryptionInTransit := d.Get("encryption_in_transit_enabled").(bool)

	var configurationsRaw interface{} = configurations
	params := clusters.ClusterCreateParametersExtended{
		Location: utils.String(location),
		Properties: &clusters.ClusterCreateProperties{
			Tier:                      pointer.To(tier),
			OsType:                    pointer.To(clusters.OSTypeLinux),
			ClusterVersion:            utils.String(clusterVersion),
			MinSupportedTlsVersion:    utils.String(tls),
			NetworkProperties:         networkProperties,
			PrivateLinkConfigurations: privateLinkConfigurations,
			EncryptionInTransitProperties: &clusters.EncryptionInTransitProperties{
				IsEncryptionInTransitEnabled: &encryptionInTransit,
			},
			ClusterDefinition: &clusters.ClusterDefinition{
				Kind:             pointer.To(clusters.ClusterKindInteractiveHive),
				ComponentVersion: componentVersions,
				Configurations:   &configurationsRaw,
			},
			StorageProfile: &clusters.StorageProfile{
				Storageaccounts: storageAccounts,
			},
			ComputeProfile: &clusters.ComputeProfile{
				Roles: roles,
			},
			ComputeIsolationProperties: computeIsolationProperties,
		},
		Tags:     tags.Expand(t),
		Identity: expandedIdentity,
	}

	if diskEncryptionPropertiesRaw, ok := d.GetOk("disk_encryption"); ok {
		params.Properties.DiskEncryptionProperties, err = ExpandHDInsightsDiskEncryptionProperties(diskEncryptionPropertiesRaw.([]interface{}))
		if err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("security_profile"); ok {
		params.Properties.SecurityProfile = ExpandHDInsightSecurityProfile(v.([]interface{}))

		// @tombuildsstuff: this behaviour is likely wrong and wants reevaluating - users should need to explicitly define this in the config?
		params.Identity = &identity.SystemAndUserAssignedMap{
			Type:        identity.TypeUserAssigned,
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		if params.Properties.SecurityProfile != nil && params.Properties.SecurityProfile.MsiResourceId != nil {
			params.Identity.IdentityIds[*params.Properties.SecurityProfile.MsiResourceId] = identity.UserAssignedIdentityDetails{
				// intentionally empty
			}
		}
	}

	if err := client.CreateThenPoll(ctx, id, params); err != nil {
		return fmt.Errorf("creating Interactive Query %s: %+v", id, err)
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

	return resourceHDInsightInteractiveQueryClusterRead(d, meta)
}

func resourceHDInsightInteractiveQueryClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] Interactive Query %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Interactive Query %s: %+v", id, err)
	}

	// Each call to configurationsClient methods is HTTP request. Getting all settings in one operation
	configurationsResp, err := configurationsClient.List(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Configuration for Interactive Query %s: %+v", id, err)
	}

	configurations := make(map[string]map[string]string)
	if model := configurationsResp.Model; model != nil && model.Configurations != nil {
		configurations = *model.Configurations
	}

	gateway, exists := configurations["gateway"]
	if !exists {
		return fmt.Errorf("retrieving Gateway Configuration for Interactive Query %s: %+v", id, err)
	}

	monitor, err := extensionsClient.GetMonitoringStatus(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Monitoring Status for Interactive Query %s: %+v", id, err)
	}

	extension, err := extensionsClient.GetAzureMonitorStatus(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Azure Monitor Status for Interactive Query %s: %+v", id, err)
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

			if err := d.Set("component_version", flattenHDInsightInteractiveQueryComponentVersion(props.ClusterDefinition.ComponentVersion)); err != nil {
				return fmt.Errorf("flattening `component_version`: %+v", err)
			}

			if err := d.Set("gateway", FlattenHDInsightsConfigurations(gateway, d)); err != nil {
				return fmt.Errorf("flattening `gateway`: %+v", err)
			}

			flattenHDInsightsMetastores(d, configurations)

			if props.EncryptionInTransitProperties != nil {
				d.Set("encryption_in_transit_enabled", props.EncryptionInTransitProperties.IsEncryptionInTransitEnabled)
			}

			diskEncryptionProps, err := flattenHDInsightsDiskEncryptionProperties(props.DiskEncryptionProperties)
			if err != nil {
				return err
			}
			if err := d.Set("disk_encryption", diskEncryptionProps); err != nil {
				return fmt.Errorf("flattening `disk_encryption`: %+v", err)
			}

			if err := d.Set("network", flattenHDInsightsNetwork(props.NetworkProperties)); err != nil {
				return fmt.Errorf("flattening `network`: %+v", err)
			}

			if err := d.Set("private_link_configuration", flattenHDInsightPrivateLinkConfigurations(props.PrivateLinkConfigurations)); err != nil {
				return fmt.Errorf("flattening `private_link_configuration`: %+v", err)
			}

			interactiveQueryRoles := hdInsightRoleDefinition{
				HeadNodeDef:      hdInsightInteractiveQueryClusterHeadNodeDefinition,
				WorkerNodeDef:    hdInsightInteractiveQueryClusterWorkerNodeDefinition,
				ZookeeperNodeDef: hdInsightInteractiveQueryClusterZookeeperNodeDefinition,
			}
			flattenedRoles := flattenHDInsightRoles(d, props.ComputeProfile, interactiveQueryRoles)
			if err := d.Set("roles", flattenedRoles); err != nil {
				return fmt.Errorf("flattening `roles`: %+v", err)
			}

			if err := d.Set("compute_isolation", flattenHDInsightComputeIsolationProperties(props.ComputeIsolationProperties)); err != nil {
				return fmt.Errorf("failed setting `compute_isolation`: %+v", err)
			}

			httpEndpoint := findHDInsightConnectivityEndpoint("HTTPS", props.ConnectivityEndpoints)
			d.Set("https_endpoint", httpEndpoint)
			sshEndpoint := findHDInsightConnectivityEndpoint("SSH", props.ConnectivityEndpoints)
			d.Set("ssh_endpoint", sshEndpoint)

			d.Set("monitor", flattenHDInsightMonitoring(monitor.Model))

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

func expandHDInsightInteractiveQueryComponentVersion(input []interface{}) *map[string]string {
	vs := input[0].(map[string]interface{})
	return &map[string]string{
		"InteractiveHive": vs["interactive_hive"].(string),
	}
}

func flattenHDInsightInteractiveQueryComponentVersion(input *map[string]string) []interface{} {
	output := make([]interface{}, 0)
	if input != nil {
		interactiveHiveVersion := ""
		if v, ok := (*input)["InteractiveHive"]; ok {
			interactiveHiveVersion = v
		}
		output = append(output, map[string]interface{}{
			"interactive_hive": interactiveHiveVersion,
		})
	}
	return output
}
