package hdinsight

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// NOTE: this isn't a recommended way of building resources in Terraform
// this pattern is used to work around a generic but pedantic API endpoint
var hdInsightInteractiveQueryClusterHeadNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         2,
	MaxInstanceCount:         utils.Int(2),
	CanSpecifyDisks:          false,
	FixedTargetInstanceCount: pointer.To(int64(2)),
}

var hdInsightInteractiveQueryClusterWorkerNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount: true,
	MinInstanceCount:        1,
	CanSpecifyDisks:         false,
	CanAutoScaleByCapacity:  true,
	CanAutoScaleOnSchedule:  true,
}

var hdInsightInteractiveQueryClusterZookeeperNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         3,
	MaxInstanceCount:         utils.Int(3),
	CanSpecifyDisks:          false,
	FixedTargetInstanceCount: pointer.To(int64(3)),
}

func resourceHDInsightInteractiveQueryCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHDInsightInteractiveQueryClusterCreate,
		Read:   resourceHDInsightInteractiveQueryClusterRead,
		Update: hdinsightClusterUpdate("Interactive Query", resourceHDInsightInteractiveQueryClusterRead),
		Delete: hdinsightClusterDelete("Interactive Query"),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := clusters.ParseClusterID(id)
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

			"roles": {
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

			"monitor": SchemaHDInsightsMonitor(),

			"extension": SchemaHDInsightsExtension(),
		},
	}
}

func resourceHDInsightInteractiveQueryClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HDInsight.ClustersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	extensionsClient := meta.(*clients.Client).HDInsight.ExtensionsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := clusters.NewClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_hdinsight_interactive_query_cluster", id.ID())
	}

	configurations := expandHDInsightsConfigurations(d.Get("gateway").([]interface{}))

	metastoresRaw := d.Get("metastores").([]interface{})
	metastores := expandHDInsightsMetastore(metastoresRaw)
	for k, v := range metastores {
		configurations[k] = v
	}

	storageAccountsRaw := d.Get("storage_account").([]interface{})
	storageAccountsGen2Raw := d.Get("storage_account_gen2").([]interface{})
	storageAccounts, i, err := expandHDInsightsStorageAccounts(storageAccountsRaw, storageAccountsGen2Raw)
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

	encryptionInTransit := d.Get("encryption_in_transit_enabled").(bool)

	params := clusters.ClusterCreateParametersExtended{
		Location: pointer.To(azure.NormalizeLocation(d.Get("location").(string))),
		Properties: &clusters.ClusterCreateProperties{
			Tier:                   pointer.To(clusters.Tier(d.Get("tier").(string))),
			OsType:                 pointer.To(clusters.OSTypeLinux),
			ClusterVersion:         pointer.To(d.Get("cluster_version").(string)),
			MinSupportedTlsVersion: pointer.To(d.Get("tls_min_version").(string)),
			NetworkProperties:      expandHDInsightsNetwork(d.Get("network").([]interface{})),
			EncryptionInTransitProperties: &clusters.EncryptionInTransitProperties{
				IsEncryptionInTransitEnabled: &encryptionInTransit,
			},
			ClusterDefinition: &clusters.ClusterDefinition{
				Kind:             pointer.To("INTERACTIVEHIVE"),
				ComponentVersion: expandHDInsightInteractiveQueryComponentVersion(d.Get("component_version").([]interface{})),
				Configurations:   configurations,
			},
			StorageProfile: &clusters.StorageProfile{
				Storageaccounts: storageAccounts,
			},
			ComputeProfile: &clusters.ComputeProfile{
				Roles: roles,
			},
			ComputeIsolationProperties: expandHDInsightComputeIsolationProperties(d.Get("compute_isolation").([]interface{})),
		},
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Identity: i,
	}

	if diskEncryptionPropertiesRaw, ok := d.GetOk("disk_encryption"); ok {
		params.Properties.DiskEncryptionProperties, err = expandHDInsightsDiskEncryptionProperties(diskEncryptionPropertiesRaw.([]interface{}))
		if err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("security_profile"); ok {
		params.Properties.SecurityProfile = expandHDInsightSecurityProfile(v.([]interface{}))

		params.Identity = &identity.SystemAndUserAssignedMap{
			Type:        identity.TypeUserAssigned,
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}

		if params.Properties.SecurityProfile != nil && params.Properties.SecurityProfile.MsiResourceId != nil {
			params.Identity.IdentityIds[*params.Properties.SecurityProfile.MsiResourceId] = identity.UserAssignedIdentityDetails{}
		}
	}

	if err := client.CreateThenPoll(ctx, id, params); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
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
	clustersClient := meta.(*clients.Client).HDInsight.ClustersClient
	configurationsClient := meta.(*clients.Client).HDInsight.ConfigurationsClient
	extensionsClient := meta.(*clients.Client).HDInsight.ExtensionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := clustersClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found  - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	// Each call to configurationsClient methods is HTTP request. Getting all settings in one operation
	configId := configurations.NewClusterID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName)
	configurations, err := configurationsClient.List(ctx, configId)
	if err != nil {
		return fmt.Errorf("retrieving Configuration for %s: %+v", *id, err)
	}

	gateway, exists := configurations.Configurations["gateway"]
	if !exists {
		return fmt.Errorf("retrieving gateway for %s: %+v", *id, err)
	}

	d.Set("name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {

		d.Set("location", location.Normalize(model.Location))

		// storage_account isn't returned, so I guess we just leave it ¯\_(ツ)_/¯
		if props := model.Properties; props != nil {
			tier := ""
			// the Azure API is inconsistent here, so rewrite this into the casing we expect
			for _, v := range clusters.PossibleValuesForTier() {
				if strings.EqualFold(v, string(pointer.From(props.Tier))) {
					tier = v
				}
			}
			d.Set("tier", tier)
			d.Set("cluster_version", props.ClusterVersion)
			d.Set("tls_min_version", props.MinSupportedTlsVersion)

			def := props.ClusterDefinition
			if err := d.Set("component_version", flattenHDInsightInteractiveQueryComponentVersion(def.ComponentVersion)); err != nil {
				return fmt.Errorf("flattening `component_version`: %+v", err)
			}

			if err := d.Set("gateway", FlattenHDInsightsConfigurations(gateway, d)); err != nil {
				return fmt.Errorf("flattening `gateway`: %+v", err)
			}

			flattenHDInsightsMetastores(d, configurations.Configurations)

			if props.EncryptionInTransitProperties != nil {
				d.Set("encryption_in_transit_enabled", props.EncryptionInTransitProperties.IsEncryptionInTransitEnabled)
			}

			if props.DiskEncryptionProperties != nil {
				diskEncryptionProps, err := flattenHDInsightsDiskEncryptionProperties(*props.DiskEncryptionProperties)
				if err != nil {
					return err
				}
				if err := d.Set("disk_encryption", diskEncryptionProps); err != nil {
					return fmt.Errorf("flattening `disk_encryption`: %+v", err)
				}
			}

			if props.NetworkProperties != nil {
				if err := d.Set("network", FlattenHDInsightsNetwork(props.NetworkProperties)); err != nil {
					return fmt.Errorf("flattening `network`: %+v", err)
				}
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

			if props.ComputeIsolationProperties != nil {
				if err := d.Set("compute_isolation", FlattenHDInsightComputeIsolationProperties(*props.ComputeIsolationProperties)); err != nil {
					return fmt.Errorf("failed setting `compute_isolation`: %+v", err)
				}
			}

			httpEndpoint := FindHDInsightConnectivityEndpoint("HTTPS", props.ConnectivityEndpoints)
			d.Set("https_endpoint", httpEndpoint)
			sshEndpoint := FindHDInsightConnectivityEndpoint("SSH", props.ConnectivityEndpoints)
			d.Set("ssh_endpoint", sshEndpoint)

			monId := extensions.NewClusterID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName)
			monitor, err := extensionsClient.GetMonitoringStatus(ctx, monId)
			if err != nil {
				return fmt.Errorf("reading monitor configuration for %s: %+v", *id, err)
			}

			d.Set("monitor", flattenHDInsightMonitoring(monitor.Model))

			extension, err := extensionsClient.GetAzureMonitorStatus(ctx, monId)
			if err != nil {
				return fmt.Errorf("reading extension configuration for %s: %+v", *id, err)
			}

			d.Set("extension", flattenHDInsightAzureMonitor(extension.Model))

			if err := d.Set("security_profile", flattenHDInsightSecurityProfile(props.SecurityProfile, d)); err != nil {
				return fmt.Errorf("setting `security_profile`: %+v", err)
			}
		}

		if err = tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting tags: %+v", err)
		}
	}

	return nil
}

func expandHDInsightInteractiveQueryComponentVersion(input []interface{}) *map[string]string {
	vs := input[0].(map[string]interface{})
	return pointer.To(map[string]string{
		"InteractiveHive": vs["interactive_hive"].(string),
	})
}

func flattenHDInsightInteractiveQueryComponentVersion(input *map[string]string) []interface{} {
	interactiveHiveVersion := ""

	if input != nil {
		if v, ok := (*input)["InteractiveHive"]; ok {
			interactiveHiveVersion = v
		}
	}

	return []interface{}{
		map[string]interface{}{
			"interactive_hive": interactiveHiveVersion,
		},
	}
}
