package azurerm

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"

	"github.com/Azure/azure-sdk-for-go/services/preview/hdinsight/mgmt/2018-06-01-preview/hdinsight"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// NOTE: this isn't a recommended way of building resources in Terraform
// this pattern is used to work around a generic but pedantic API endpoint
var hdInsightHadoopClusterHeadNodeDefinition = azure.HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         2,
	MaxInstanceCount:         2,
	CanSpecifyDisks:          false,
	FixedMinInstanceCount:    utils.Int32(int32(1)),
	FixedTargetInstanceCount: utils.Int32(int32(2)),
}

var hdInsightHadoopClusterWorkerNodeDefinition = azure.HDInsightNodeDefinition{
	CanSpecifyInstanceCount: true,
	MinInstanceCount:        1,
	MaxInstanceCount:        25,
	CanSpecifyDisks:         false,
}

var hdInsightHadoopClusterZookeeperNodeDefinition = azure.HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         3,
	MaxInstanceCount:         3,
	CanSpecifyDisks:          false,
	FixedMinInstanceCount:    utils.Int32(int32(1)),
	FixedTargetInstanceCount: utils.Int32(int32(3)),
}

func resourceArmHDInsightHadoopCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmHDInsightHadoopClusterCreate,
		Read:   resourceArmHDInsightHadoopClusterRead,
		Update: hdinsightClusterUpdate("Hadoop", resourceArmHDInsightHadoopClusterRead),
		Delete: hdinsightClusterDelete("Hadoop"),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			// An edge node can be added but can't be update or removed without forcing a new resource to be created
			oldEdgeNodeCount, newEdgeNodeCount := diff.GetChange("roles.0.edge_node.#")
			oldEdgeNodeInt := oldEdgeNodeCount.(int)
			newEdgeNodeInt := newEdgeNodeCount.(int)

			// ForceNew if attempting to remove an edge node
			if newEdgeNodeInt < oldEdgeNodeInt {
				diff.ForceNew("roles.0.edge_node")
			}

			// ForceNew if attempting to update an edge node
			if newEdgeNodeInt == 1 && oldEdgeNodeInt == 1 {
				// DiffSuppressFunc comes after this check so we need to check if the strings aren't the same sans casing here.
				oVMSize, newVMSize := diff.GetChange("roles.0.edge_node.0.vm_size")
				if !strings.EqualFold(oVMSize.(string), newVMSize.(string)) {
					diff.ForceNew("roles.0.edge_node")
				}

				// ForceNew if attempting to update install scripts
				oldInstallScriptCount, newInstallScriptCount := diff.GetChange("roles.0.edge_node.0.install_script_action.#")
				oldInstallScriptInt := oldInstallScriptCount.(int)
				newInstallScriptInt := newInstallScriptCount.(int)
				if newInstallScriptInt == oldInstallScriptInt {
					diff.ForceNew("roles.0.edge_node.0.install_script_action")
				}
			}

			return nil
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
						"hadoop": {
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
						"head_node": azure.SchemaHDInsightNodeDefinition("roles.0.head_node", hdInsightHadoopClusterHeadNodeDefinition),

						"worker_node": azure.SchemaHDInsightNodeDefinition("roles.0.worker_node", hdInsightHadoopClusterWorkerNodeDefinition),

						"zookeeper_node": azure.SchemaHDInsightNodeDefinition("roles.0.zookeeper_node", hdInsightHadoopClusterZookeeperNodeDefinition),

						"edge_node": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vm_size": {
										Type:             schema.TypeString,
										Required:         true,
										DiffSuppressFunc: suppress.CaseDifference,
										ValidateFunc:     azure.ValidateSchemaHDInsightNodeDefinitionVMSize(),
									},

									"install_script_action": {
										Type:     schema.TypeList,
										Required: true,
										MinItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validate.NoEmptyStrings,
												},
												"uri": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validate.NoEmptyStrings,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"tags": tags.Schema(),

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

func resourceArmHDInsightHadoopClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).hdinsight.ClustersClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	clusterVersion := d.Get("cluster_version").(string)
	t := d.Get("tags").(map[string]interface{})
	tier := hdinsight.Tier(d.Get("tier").(string))

	componentVersionsRaw := d.Get("component_version").([]interface{})
	componentVersions := expandHDInsightHadoopComponentVersion(componentVersionsRaw)

	gatewayRaw := d.Get("gateway").([]interface{})
	gateway := azure.ExpandHDInsightsConfigurations(gatewayRaw)

	storageAccountsRaw := d.Get("storage_account").([]interface{})
	storageAccounts, err := azure.ExpandHDInsightsStorageAccounts(storageAccountsRaw)
	if err != nil {
		return fmt.Errorf("Error expanding `storage_account`: %s", err)
	}

	rolesRaw := d.Get("roles").([]interface{})
	hadoopRoles := hdInsightRoleDefinition{
		HeadNodeDef:      hdInsightHadoopClusterHeadNodeDefinition,
		WorkerNodeDef:    hdInsightHadoopClusterWorkerNodeDefinition,
		ZookeeperNodeDef: hdInsightHadoopClusterZookeeperNodeDefinition,
	}
	roles, err := expandHDInsightRoles(rolesRaw, hadoopRoles)
	if err != nil {
		return fmt.Errorf("Error expanding `roles`: %+v", err)
	}

	if features.ShouldResourcesBeImported() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_hdinsight_hadoop_cluster", *existing.ID)
		}
	}

	params := hdinsight.ClusterCreateParametersExtended{
		Location: utils.String(location),
		Properties: &hdinsight.ClusterCreateProperties{
			Tier:           tier,
			OsType:         hdinsight.Linux,
			ClusterVersion: utils.String(clusterVersion),
			ClusterDefinition: &hdinsight.ClusterDefinition{
				Kind:             utils.String("Hadoop"),
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
		Tags: tags.Expand(t),
	}
	future, err := client.Create(ctx, resourceGroup, name, params)
	if err != nil {
		return fmt.Errorf("Error creating HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Error reading ID for HDInsight Hadoop Cluster %q (Resource Group %q)", name, resourceGroup)
	}

	d.SetId(*read.ID)

	// We can only add an edge node after creation
	if v, ok := d.GetOk("roles.0.edge_node"); ok {
		edgeNodeRaw := v.([]interface{})
		applicationsClient := meta.(*ArmClient).hdinsight.ApplicationsClient
		edgeNodeConfig := edgeNodeRaw[0].(map[string]interface{})

		err := createHDInsightEdgeNode(ctx, applicationsClient, resourceGroup, name, edgeNodeConfig)
		if err != nil {
			return err
		}

		// we can't rely on the use of the Future here due to the node being successfully completed but now the cluster is applying those changes.
		log.Printf("[DEBUG] Waiting for Hadoop Cluster to %q (Resource Group %q) to finish applying edge node", name, resourceGroup)
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"AzureVMConfiguration", "Accepted", "HdInsightConfiguration"},
			Target:     []string{"Running"},
			Refresh:    hdInsightWaitForReadyRefreshFunc(ctx, client, resourceGroup, name),
			Timeout:    60 * time.Minute,
			MinTimeout: 15 * time.Second,
		}
		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("Error waiting for HDInsight Cluster %q (Resource Group %q) to be running: %s", name, resourceGroup, err)
		}

		/*
			edgeNodeStateConf := &resource.StateChangeConf{
				Pending:    []string{"Empty"},
				Target:     []string{"Ready"},
				Refresh:    hdInsightEdgeNodeWaitForReadyRefreshFunc(ctx, applicationsClient, resourceGroup, name),
				Timeout:    3 * time.Minute,
				MinTimeout: 15 * time.Second,
			}

			if _, err := edgeNodeStateConf.WaitForState(); err != nil {
				return fmt.Errorf("Error waiting for HDInsight Cluster Edge Node %q (Resource Group %q) to be ready: %s", name, resourceGroup, err)
			} */
	}

	return resourceArmHDInsightHadoopClusterRead(d, meta)
}

func resourceArmHDInsightHadoopClusterRead(d *schema.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] HDInsight Hadoop Cluster %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	configuration, err := configurationsClient.Get(ctx, resourceGroup, name, "gateway")
	if err != nil {
		return fmt.Errorf("Error retrieving Configuration for HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
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
			if err := d.Set("component_version", flattenHDInsightHadoopComponentVersion(def.ComponentVersion)); err != nil {
				return fmt.Errorf("Error flattening `component_version`: %+v", err)
			}

			if err := d.Set("gateway", azure.FlattenHDInsightsConfigurations(configuration.Value)); err != nil {
				return fmt.Errorf("Error flattening `gateway`: %+v", err)
			}
		}

		hadoopRoles := hdInsightRoleDefinition{
			HeadNodeDef:      hdInsightHadoopClusterHeadNodeDefinition,
			WorkerNodeDef:    hdInsightHadoopClusterWorkerNodeDefinition,
			ZookeeperNodeDef: hdInsightHadoopClusterZookeeperNodeDefinition,
		}
		flattenedRoles := flattenHDInsightRoles(d, props.ComputeProfile, hadoopRoles)

		applicationsClient := meta.(*ArmClient).hdinsight.ApplicationsClient

		edgeNode, err := applicationsClient.Get(ctx, resourceGroup, name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(edgeNode.Response) {
				return fmt.Errorf("Error reading edge node for HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if edgeNodeProps := edgeNode.Properties; edgeNodeProps != nil {
			flattenedRoles = flattenHDInsightEdgeNode(flattenedRoles, edgeNodeProps)
		}

		if err := d.Set("roles", flattenedRoles); err != nil {
			return fmt.Errorf("Error flattening `roles`: %+v", err)
		}

		httpEndpoint := azure.FindHDInsightConnectivityEndpoint("HTTPS", props.ConnectivityEndpoints)
		d.Set("https_endpoint", httpEndpoint)
		sshEndpoint := azure.FindHDInsightConnectivityEndpoint("SSH", props.ConnectivityEndpoints)
		d.Set("ssh_endpoint", sshEndpoint)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenHDInsightEdgeNode(roles []interface{}, props *hdinsight.ApplicationProperties) []interface{} {
	if len(roles) == 0 || props == nil {
		return roles
	}

	role := roles[0].(map[string]interface{})

	edgeNode := make(map[string]interface{})
	if computeProfile := props.ComputeProfile; computeProfile != nil {
		if roles := computeProfile.Roles; roles != nil {
			for _, role := range *roles {
				if hardwareProfile := role.HardwareProfile; hardwareProfile != nil {
					edgeNode["vm_size"] = hardwareProfile.VMSize
				}
			}
		}
	}

	actions := make(map[string]interface{})
	if installScriptActions := props.InstallScriptActions; installScriptActions != nil {
		for _, action := range *installScriptActions {
			actions["name"] = action.Name
			actions["uri"] = action.URI
		}
	}

	edgeNode["install_script_action"] = []interface{}{actions}

	role["edge_node"] = []interface{}{edgeNode}

	return []interface{}{role}
}

func expandHDInsightHadoopComponentVersion(input []interface{}) map[string]*string {
	vs := input[0].(map[string]interface{})
	return map[string]*string{
		"Hadoop": utils.String(vs["hadoop"].(string)),
	}
}

func flattenHDInsightHadoopComponentVersion(input map[string]*string) []interface{} {
	hadoopVersion := ""
	if v, ok := input["Hadoop"]; ok {
		if v != nil {
			hadoopVersion = *v
		}
	}
	return []interface{}{
		map[string]interface{}{
			"hadoop": hadoopVersion,
		},
	}
}

func expandHDInsightApplicationEdgeNodeInstallScriptActions(input []interface{}) *[]hdinsight.RuntimeScriptAction {
	actions := make([]hdinsight.RuntimeScriptAction, 0)

	for _, v := range input {
		val := v.(map[string]interface{})

		name := val["name"].(string)
		uri := val["uri"].(string)

		action := hdinsight.RuntimeScriptAction{
			Name: utils.String(name),
			URI:  utils.String(uri),
			// The only role available for edge nodes is edgenode
			Roles: &[]string{"edgenode"},
		}

		actions = append(actions, action)
	}

	return &actions
}

func retryHDInsightEdgeNodeGet(resGroup string, name string, meta interface{}) func() *resource.RetryError {
	return func() *resource.RetryError {
		client := meta.(*ArmClient).hdinsight.ApplicationsClient
		ctx := meta.(*ArmClient).StopContext

		if _, err := client.Get(ctx, resGroup, name, name); err != nil {
			return resource.RetryableError(err)
		}

		return nil
	}
}

func hdInsightWaitForReadyRefreshFunc(ctx context.Context, client *hdinsight.ClustersClient, resourceGroupName string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, name)
		if err != nil {
			return nil, "Error", fmt.Errorf("Error issuing read request in hdInsightWaitForReadyRefreshFunc to Hadoop Cluster %q (Resource Group %q): %s", name, resourceGroupName, err)
		}
		if props := res.Properties; props != nil {
			if state := props.ClusterState; state != nil {
				return res, *state, nil
			}
		}

		return res, "Pending", nil
	}
}

func hdInsightEdgeNodeWaitForReadyRefreshFunc(ctx context.Context, client hdinsight.ApplicationsClient, resourceGroupName string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, name, name)
		if err != nil {
			if res.Response.Response != nil {
				return nil, "Error", fmt.Errorf("Error issuing read request in hdInsightEdgeNodeWaitForReadyRefreshFunc to Hadoop Cluster Edge Node %q (Resource Group %q): %s", name, resourceGroupName, err)
			}
			return res, "Empty", nil
		}
		if props := res.Properties; props != nil {
			return nil, "Ready", nil
		}

		return res, "Empty", nil
	}
}
