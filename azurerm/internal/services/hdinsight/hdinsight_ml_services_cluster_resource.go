package hdinsight

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/hdinsight/mgmt/2018-06-01/hdinsight"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hdinsight/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// NOTE: this isn't a recommended way of building resources in Terraform
// this pattern is used to work around a generic but pedantic API endpoint
var hdInsightMLServicesClusterHeadNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         2,
	MaxInstanceCount:         utils.Int(2),
	CanSpecifyDisks:          false,
	FixedMinInstanceCount:    utils.Int32(int32(1)),
	FixedTargetInstanceCount: utils.Int32(int32(2)),
}

var hdInsightMLServicesClusterWorkerNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount: true,
	MinInstanceCount:        1,
	CanSpecifyDisks:         false,
}

var hdInsightMLServicesClusterZookeeperNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         3,
	MaxInstanceCount:         utils.Int(3),
	CanSpecifyDisks:          false,
	FixedMinInstanceCount:    utils.Int32(int32(1)),
	FixedTargetInstanceCount: utils.Int32(int32(3)),
}

var hdInsightMLServicesClusterEdgeNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         1,
	MaxInstanceCount:         utils.Int(1),
	CanSpecifyDisks:          false,
	FixedTargetInstanceCount: utils.Int32(int32(1)),
}

func resourceHDInsightMLServicesCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		DeprecationMessage: `HDInsight 3.6 will be retired on 2020-12-31 - MLServices is not supported in HDInsight 4.0 and so this resource will be removed in the next major version of the AzureRM Terraform Provider.
		
More information on the HDInsight 3.6 deprecation can be found at:

https://docs.microsoft.com/en-us/azure/hdinsight/hdinsight-component-versioning#available-versions`,
		Create: resourceHDInsightMLServicesClusterCreate,
		Read:   resourceHDInsightMLServicesClusterRead,
		Update: hdinsightClusterUpdate("MLServices", resourceHDInsightMLServicesClusterRead),
		Delete: hdinsightClusterDelete("MLServices"),
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

			"gateway": SchemaHDInsightsGateway(),

			"rstudio": {
				Type:     pluginsdk.TypeBool,
				Required: true,
				ForceNew: true,
			},

			"storage_account": SchemaHDInsightsStorageAccounts(),

			"roles": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"head_node": SchemaHDInsightNodeDefinition("roles.0.head_node", hdInsightMLServicesClusterHeadNodeDefinition, true),

						"worker_node": SchemaHDInsightNodeDefinition("roles.0.worker_node", hdInsightMLServicesClusterWorkerNodeDefinition, true),

						"zookeeper_node": SchemaHDInsightNodeDefinition("roles.0.zookeeper_node", hdInsightMLServicesClusterZookeeperNodeDefinition, true),

						"edge_node": SchemaHDInsightNodeDefinition("roles.0.edge_node", hdInsightMLServicesClusterEdgeNodeDefinition, true),
					},
				},
			},

			"tags": tags.Schema(),

			"edge_ssh_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"https_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"ssh_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func expandHDInsightsMLServicesConfigurations(gateway []interface{}, rStudio bool) map[string]interface{} {
	config := ExpandHDInsightsConfigurations(gateway)

	config["rserver"] = map[string]interface{}{
		"rstudio": rStudio,
	}

	return config
}

func resourceHDInsightMLServicesClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HDInsight.ClustersClient
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

	gatewayRaw := d.Get("gateway").([]interface{})
	rStudio := d.Get("rstudio").(bool)
	gateway := expandHDInsightsMLServicesConfigurations(gatewayRaw, rStudio)

	storageAccountsRaw := d.Get("storage_account").([]interface{})
	storageAccounts, identity, err := ExpandHDInsightsStorageAccounts(storageAccountsRaw, nil)
	if err != nil {
		return fmt.Errorf("Error expanding `storage_account`: %s", err)
	}

	mlServicesRoles := hdInsightRoleDefinition{
		HeadNodeDef:      hdInsightMLServicesClusterHeadNodeDefinition,
		WorkerNodeDef:    hdInsightMLServicesClusterWorkerNodeDefinition,
		ZookeeperNodeDef: hdInsightMLServicesClusterZookeeperNodeDefinition,
		EdgeNodeDef:      &hdInsightMLServicesClusterEdgeNodeDefinition,
	}
	rolesRaw := d.Get("roles").([]interface{})
	roles, err := expandHDInsightRoles(rolesRaw, mlServicesRoles)
	if err != nil {
		return fmt.Errorf("Error expanding `roles`: %+v", err)
	}

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing HDInsight MLServices Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_hdinsight_ml_server_cluster", *existing.ID)
	}

	params := hdinsight.ClusterCreateParametersExtended{
		Location: utils.String(location),
		Properties: &hdinsight.ClusterCreateProperties{
			Tier:                   tier,
			OsType:                 hdinsight.OSTypeLinux,
			ClusterVersion:         utils.String(clusterVersion),
			MinSupportedTLSVersion: utils.String(tls),
			ClusterDefinition: &hdinsight.ClusterDefinition{
				Kind:           utils.String("MLServices"),
				Configurations: gateway,
			},
			StorageProfile: &hdinsight.StorageProfile{
				Storageaccounts: storageAccounts,
			},
			ComputeProfile: &hdinsight.ComputeProfile{
				Roles: roles,
			},
		},
		Tags:     tags.Expand(t),
		Identity: identity,
	}
	future, err := client.Create(ctx, resourceGroup, name, params)
	if err != nil {
		return fmt.Errorf("Error creating HDInsight MLServices Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of HDInsight MLServices Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving HDInsight MLServices Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Error reading ID for HDInsight MLServices Cluster %q (Resource Group %q)", name, resourceGroup)
	}

	d.SetId(id.ID())

	return resourceHDInsightMLServicesClusterRead(d, meta)
}

func resourceHDInsightMLServicesClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	clustersClient := meta.(*clients.Client).HDInsight.ClustersClient
	configurationsClient := meta.(*clients.Client).HDInsight.ConfigurationsClient
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
			log.Printf("[DEBUG] HDInsight MLServices Cluster %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving HDInsight MLServices Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	configuration, err := configurationsClient.Get(ctx, resourceGroup, name, "gateway")
	if err != nil {
		return fmt.Errorf("Error retrieving Gateway Configuration for HDInsight MLServices Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	rStudioConfig, err := configurationsClient.Get(ctx, resourceGroup, name, "rserver")
	if err != nil {
		return fmt.Errorf("Error retrieving RStudio Configuration for HDInsight MLServices Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
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
			if err := d.Set("gateway", FlattenHDInsightsConfigurations(configuration.Value)); err != nil {
				return fmt.Errorf("Error flattening `gateway`: %+v", err)
			}

			var rStudio bool
			if rStudioStr := rStudioConfig.Value["rstudio"]; rStudioStr != nil {
				rStudioBool, err := strconv.ParseBool(*rStudioStr)
				if err != nil {
					return err
				}

				rStudio = rStudioBool
			}

			d.Set("rstudio", rStudio)
		}

		mlServicesRoles := hdInsightRoleDefinition{
			HeadNodeDef:      hdInsightMLServicesClusterHeadNodeDefinition,
			WorkerNodeDef:    hdInsightMLServicesClusterWorkerNodeDefinition,
			ZookeeperNodeDef: hdInsightMLServicesClusterZookeeperNodeDefinition,
			EdgeNodeDef:      &hdInsightMLServicesClusterEdgeNodeDefinition,
		}
		flattenedRoles := flattenHDInsightRoles(d, props.ComputeProfile, mlServicesRoles)
		if err := d.Set("roles", flattenedRoles); err != nil {
			return fmt.Errorf("Error flattening `roles`: %+v", err)
		}

		edgeSSHEndpoint := FindHDInsightConnectivityEndpoint("EDGESSH", props.ConnectivityEndpoints)
		d.Set("edge_ssh_endpoint", edgeSSHEndpoint)
		httpEndpoint := FindHDInsightConnectivityEndpoint("HTTPS", props.ConnectivityEndpoints)
		d.Set("https_endpoint", httpEndpoint)
		sshEndpoint := FindHDInsightConnectivityEndpoint("SSH", props.ConnectivityEndpoints)
		d.Set("ssh_endpoint", sshEndpoint)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
