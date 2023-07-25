// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hdinsight

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/hdinsight/mgmt/2018-06-01/hdinsight" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// NOTE: this isn't a recommended way of building resources in Terraform
// this pattern is used to work around a generic but pedantic API endpoint
var hdInsightSparkClusterHeadNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         2,
	MaxInstanceCount:         utils.Int(2),
	CanSpecifyDisks:          false,
	FixedTargetInstanceCount: utils.Int32(int32(2)),
}

var hdInsightSparkClusterWorkerNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount: true,
	MinInstanceCount:        1,
	CanSpecifyDisks:         false,
	CanAutoScaleByCapacity:  true,
	CanAutoScaleOnSchedule:  true,
}

var hdInsightSparkClusterZookeeperNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         3,
	MaxInstanceCount:         utils.Int(3),
	FixedTargetInstanceCount: utils.Int32(int32(3)),
	CanSpecifyDisks:          false,
}

func resourceHDInsightSparkCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHDInsightSparkClusterCreate,
		Read:   resourceHDInsightSparkClusterRead,
		Update: hdinsightClusterUpdate("Spark", resourceHDInsightSparkClusterRead),
		Delete: hdinsightClusterDelete("Spark"),

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

			"encryption_in_transit_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"disk_encryption": SchemaHDInsightsDiskEncryptionProperties(),

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

			"compute_isolation": SchemaHDInsightsComputeIsolation(),

			"gateway": SchemaHDInsightsGateway(),

			"metastores": SchemaHDInsightsExternalMetastores(),

			"network": SchemaHDInsightsNetwork(),

			"security_profile": SchemaHDInsightsSecurityProfile(),

			"storage_account": SchemaHDInsightsStorageAccounts(),

			"storage_account_gen2": SchemaHDInsightsGen2StorageAccounts(),

			"roles": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"head_node": SchemaHDInsightNodeDefinition("roles.0.head_node", hdInsightSparkClusterHeadNodeDefinition, true),

						"worker_node": SchemaHDInsightNodeDefinition("roles.0.worker_node", hdInsightSparkClusterWorkerNodeDefinition, true),

						"zookeeper_node": SchemaHDInsightNodeDefinition("roles.0.zookeeper_node", hdInsightSparkClusterZookeeperNodeDefinition, true),
					},
				},
			},

			"tags": tags.Schema(),

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
}

func resourceHDInsightSparkClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
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
	componentVersions := expandHDInsightSparkComponentVersion(componentVersionsRaw)

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
		return fmt.Errorf("expanding `storage_account`: %s", err)
	}

	networkPropertiesRaw := d.Get("network").([]interface{})
	networkProperties := ExpandHDInsightsNetwork(networkPropertiesRaw)

	sparkRoles := hdInsightRoleDefinition{
		HeadNodeDef:      hdInsightSparkClusterHeadNodeDefinition,
		WorkerNodeDef:    hdInsightSparkClusterWorkerNodeDefinition,
		ZookeeperNodeDef: hdInsightSparkClusterZookeeperNodeDefinition,
	}
	rolesRaw := d.Get("roles").([]interface{})
	roles, err := expandHDInsightRoles(rolesRaw, sparkRoles)
	if err != nil {
		return fmt.Errorf("expanding `roles`: %+v", err)
	}

	computeIsolationProperties := ExpandHDInsightComputeIsolationProperties(d.Get("compute_isolation").([]interface{}))

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing HDInsight Spark Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_hdinsight_spark_cluster", id.ID())
	}

	encryptionInTransit := d.Get("encryption_in_transit_enabled").(bool)

	params := hdinsight.ClusterCreateParametersExtended{
		Location: utils.String(location),
		Properties: &hdinsight.ClusterCreateProperties{
			Tier:           tier,
			OsType:         hdinsight.OSTypeLinux,
			ClusterVersion: utils.String(clusterVersion),
			EncryptionInTransitProperties: &hdinsight.EncryptionInTransitProperties{
				IsEncryptionInTransitEnabled: &encryptionInTransit,
			},
			MinSupportedTLSVersion: utils.String(tls),
			NetworkProperties:      networkProperties,
			ClusterDefinition: &hdinsight.ClusterDefinition{
				Kind:             utils.String("Spark"),
				ComponentVersion: componentVersions,
				Configurations:   configurations,
			},
			StorageProfile: &hdinsight.StorageProfile{
				Storageaccounts: storageAccounts,
			},
			ComputeProfile: &hdinsight.ComputeProfile{
				Roles: roles,
			},
			ComputeIsolationProperties: computeIsolationProperties,
		},
		Tags:     tags.Expand(t),
		Identity: identity,
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
		return fmt.Errorf("creating HDInsight Spark Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of HDInsight Spark Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving HDInsight Spark Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("reading ID for HDInsight Spark Cluster %q (Resource Group %q)", name, resourceGroup)
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

	return resourceHDInsightSparkClusterRead(d, meta)
}

func resourceHDInsightSparkClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] HDInsight Spark Cluster %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving HDInsight Spark Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	// Each call to configurationsClient methods is HTTP request. Getting all settings in one operation
	configurations, err := configurationsClient.List(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Configuration for HDInsight Spark Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	gateway, exists := configurations.Configurations["gateway"]
	if !exists {
		return fmt.Errorf("retrieving gateway for HDInsight Spark Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
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
			if err := d.Set("component_version", flattenHDInsightSparkComponentVersion(def.ComponentVersion)); err != nil {
				return fmt.Errorf("flattening `component_version`: %+v", err)
			}

			if err := d.Set("gateway", FlattenHDInsightsConfigurations(gateway, d)); err != nil {
				return fmt.Errorf("flattening `gateway`: %+v", err)
			}

			flattenHDInsightsMetastores(d, configurations.Configurations)
		}

		sparkRoles := hdInsightRoleDefinition{
			HeadNodeDef:      hdInsightSparkClusterHeadNodeDefinition,
			WorkerNodeDef:    hdInsightSparkClusterWorkerNodeDefinition,
			ZookeeperNodeDef: hdInsightSparkClusterZookeeperNodeDefinition,
		}

		if props.EncryptionInTransitProperties != nil {
			d.Set("encryption_in_transit_enabled", props.EncryptionInTransitProperties.IsEncryptionInTransitEnabled)
		}

		if props.DiskEncryptionProperties != nil {
			diskEncryptionProps, err := FlattenHDInsightsDiskEncryptionProperties(*props.DiskEncryptionProperties)
			if err != nil {
				return err
			}
			if err := d.Set("disk_encryption", diskEncryptionProps); err != nil {
				return fmt.Errorf("flattening setting `disk_encryption`: %+v", err)
			}
		}

		if props.NetworkProperties != nil {
			if err := d.Set("network", FlattenHDInsightsNetwork(props.NetworkProperties)); err != nil {
				return fmt.Errorf("flattening `network`: %+v", err)
			}
		}

		flattenedRoles := flattenHDInsightRoles(d, props.ComputeProfile, sparkRoles)
		if err := d.Set("roles", flattenedRoles); err != nil {
			return fmt.Errorf("flattening `roles`: %+v", err)
		}

		if props.ComputeIsolationProperties != nil {
			if err := d.Set("compute_isolation", FlattenHDInsightComputeIsolationProperties(*props.ComputeIsolationProperties)); err != nil {
				return fmt.Errorf("failed setting `compute_isolation`: %+v", err)
			}
		}

		httpEndpoint := FindHDInsightConnectivityEndpoint("HTTPS", props.ConnectivityEndpoints)
		d.Set("https_endpoint", httpEndpoint)
		sshEndpoint := FindHDInsightConnectivityEndpoint("SSH", props.ConnectivityEndpoints)
		d.Set("ssh_endpoint", sshEndpoint)

		monitor, err := extensionsClient.GetMonitoringStatus(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("reading monitor configuration for HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		d.Set("monitor", flattenHDInsightMonitoring(monitor))

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

func expandHDInsightSparkComponentVersion(input []interface{}) map[string]*string {
	vs := input[0].(map[string]interface{})
	return map[string]*string{
		"Spark": utils.String(vs["spark"].(string)),
	}
}

func flattenHDInsightSparkComponentVersion(input map[string]*string) []interface{} {
	sparkVersion := ""
	if v, ok := input["Spark"]; ok {
		if v != nil {
			sparkVersion = *v
		}
	}
	return []interface{}{
		map[string]interface{}{
			"spark": sparkVersion,
		},
	}
}
