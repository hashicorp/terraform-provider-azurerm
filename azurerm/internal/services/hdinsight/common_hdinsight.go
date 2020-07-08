package hdinsight

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/hdinsight/mgmt/2018-06-01-preview/hdinsight"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func hdinsightClusterUpdate(clusterKind string, readFunc schema.ReadFunc) schema.UpdateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		client := meta.(*clients.Client).HDInsight.ClustersClient
		ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
		defer cancel()

		id, err := azure.ParseAzureResourceID(d.Id())
		if err != nil {
			return err
		}

		resourceGroup := id.ResourceGroup
		name := id.Path["clusters"]

		if d.HasChange("tags") {
			t := d.Get("tags").(map[string]interface{})
			params := hdinsight.ClusterPatchParameters{
				Tags: tags.Expand(t),
			}
			if _, err := client.Update(ctx, resourceGroup, name, params); err != nil {
				return fmt.Errorf("Error updating Tags for HDInsight %q Cluster %q (Resource Group %q): %+v", clusterKind, name, resourceGroup, err)
			}
		}

		if d.HasChange("roles.0.worker_node") {
			log.Printf("[DEBUG] Resizing the HDInsight %q Cluster", clusterKind)
			rolesRaw := d.Get("roles").([]interface{})
			roles := rolesRaw[0].(map[string]interface{})
			workerNodes := roles["worker_node"].([]interface{})
			workerNode := workerNodes[0].(map[string]interface{})
			targetInstanceCount := workerNode["target_instance_count"].(int)
			params := hdinsight.ClusterResizeParameters{
				TargetInstanceCount: utils.Int32(int32(targetInstanceCount)),
			}

			future, err := client.Resize(ctx, resourceGroup, name, params)
			if err != nil {
				return fmt.Errorf("Error resizing the HDInsight %q Cluster %q (Resource Group %q): %+v", clusterKind, name, resourceGroup, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("Error waiting for the HDInsight %q Cluster %q (Resource Group %q) to finish resizing: %+v", clusterKind, name, resourceGroup, err)
			}
		}

		// The API can add an edge node but can't remove them without force newing the resource. We'll check for adding here
		// and can come back to removing if that functionality gets added. https://feedback.azure.com/forums/217335-hdinsight/suggestions/5663773-start-stop-cluster-hdinsight?page=3&per_page=20
		if clusterKind == "Hadoop" {
			if d.HasChange("roles.0.edge_node") {
				log.Printf("[DEBUG] Detected change in edge nodes")
				edgeNodeRaw := d.Get("roles.0.edge_node").([]interface{})
				edgeNodeConfig := edgeNodeRaw[0].(map[string]interface{})
				applicationsClient := meta.(*clients.Client).HDInsight.ApplicationsClient

				oldEdgeNodeCount, newEdgeNodeCount := d.GetChange("roles.0.edge_node.0.target_instance_count")
				oldEdgeNodeInt := oldEdgeNodeCount.(int)
				newEdgeNodeInt := newEdgeNodeCount.(int)

				// Note: API currently doesn't support updating number of edge nodes
				// if anything in the edge nodes changes, delete edge nodes then recreate them
				if oldEdgeNodeInt != 0 {
					err := deleteHDInsightEdgeNodes(ctx, applicationsClient, resourceGroup, name)
					if err != nil {
						return err
					}
				}

				if newEdgeNodeInt != 0 {
					err = createHDInsightEdgeNodes(ctx, applicationsClient, resourceGroup, name, edgeNodeConfig)
					if err != nil {
						return err
					}
				}

				// we can't rely on the use of the Future here due to the node being successfully completed but now the cluster is applying those changes.
				log.Printf("[DEBUG] Waiting for Hadoop Cluster to %q (Resource Group %q) to finish applying edge node", name, resourceGroup)
				stateConf := &resource.StateChangeConf{
					Pending:    []string{"AzureVMConfiguration", "Accepted", "HdInsightConfiguration"},
					Target:     []string{"Running"},
					Refresh:    hdInsightWaitForReadyRefreshFunc(ctx, client, resourceGroup, name),
					MinTimeout: 15 * time.Second,
					Timeout:    d.Timeout(schema.TimeoutUpdate),
				}

				if _, err := stateConf.WaitForState(); err != nil {
					return fmt.Errorf("Error waiting for HDInsight Cluster %q (Resource Group %q) to be running: %s", name, resourceGroup, err)
				}
			}
		}

		if d.HasChange("gateway") {
			log.Printf("[DEBUG] Updating the HDInsight %q Cluster gateway", clusterKind)
			vs := d.Get("gateway").([]interface{})[0].(map[string]interface{})

			enabled := true
			username := vs["username"].(string)
			password := vs["password"].(string)

			future, err := client.UpdateGatewaySettings(ctx, resourceGroup, name, hdinsight.UpdateGatewaySettingsParameters{
				IsCredentialEnabled: &enabled,
				UserName:            utils.String(username),
				Password:            utils.String(password),
			})

			if err != nil {
				return err
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("Error waiting for HDInsight Cluster %q (Resource Group %q) Gateway to be updated: %s", name, resourceGroup, err)
			}
		}

		return readFunc(d, meta)
	}
}

func hdinsightClusterDelete(clusterKind string) schema.DeleteFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		client := meta.(*clients.Client).HDInsight.ClustersClient
		ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
		defer cancel()

		id, err := azure.ParseAzureResourceID(d.Id())
		if err != nil {
			return err
		}

		resourceGroup := id.ResourceGroup
		name := id.Path["clusters"]

		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Error deleting HDInsight %q Cluster %q (Resource Group %q): %+v", clusterKind, name, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for deletion of HDInsight %q Cluster %q (Resource Group %q): %+v", clusterKind, name, resourceGroup, err)
		}

		return nil
	}
}

type hdInsightRoleDefinition struct {
	HeadNodeDef      azure.HDInsightNodeDefinition
	WorkerNodeDef    azure.HDInsightNodeDefinition
	ZookeeperNodeDef azure.HDInsightNodeDefinition
	EdgeNodeDef      *azure.HDInsightNodeDefinition
}

func expandHDInsightRoles(input []interface{}, definition hdInsightRoleDefinition) (*[]hdinsight.Role, error) {
	v := input[0].(map[string]interface{})

	headNodeRaw := v["head_node"].([]interface{})
	headNode, err := azure.ExpandHDInsightNodeDefinition("headnode", headNodeRaw, definition.HeadNodeDef)
	if err != nil {
		return nil, fmt.Errorf("Error expanding `head_node`: %+v", err)
	}

	workerNodeRaw := v["worker_node"].([]interface{})
	workerNode, err := azure.ExpandHDInsightNodeDefinition("workernode", workerNodeRaw, definition.WorkerNodeDef)
	if err != nil {
		return nil, fmt.Errorf("Error expanding `worker_node`: %+v", err)
	}

	zookeeperNodeRaw := v["zookeeper_node"].([]interface{})
	zookeeperNode, err := azure.ExpandHDInsightNodeDefinition("zookeepernode", zookeeperNodeRaw, definition.ZookeeperNodeDef)
	if err != nil {
		return nil, fmt.Errorf("Error expanding `zookeeper_node`: %+v", err)
	}

	roles := []hdinsight.Role{
		*headNode,
		*workerNode,
		*zookeeperNode,
	}

	if definition.EdgeNodeDef != nil {
		edgeNodeRaw := v["edge_node"].([]interface{})
		edgeNode, err := azure.ExpandHDInsightNodeDefinition("edgenode", edgeNodeRaw, *definition.EdgeNodeDef)
		if err != nil {
			return nil, fmt.Errorf("Error expanding `edge_node`: %+v", err)
		}
		roles = append(roles, *edgeNode)
	}

	return &roles, nil
}

func flattenHDInsightRoles(d *schema.ResourceData, input *hdinsight.ComputeProfile, definition hdInsightRoleDefinition) []interface{} {
	if input == nil || input.Roles == nil {
		return []interface{}{}
	}

	var existingEdgeNodes, existingHeadNodes, existingWorkerNodes, existingZookeeperNodes []interface{}

	existingVs := d.Get("roles").([]interface{})
	if len(existingVs) > 0 {
		existingV := existingVs[0].(map[string]interface{})

		if definition.EdgeNodeDef != nil {
			existingEdgeNodes = existingV["edge_node"].([]interface{})
		}

		existingHeadNodes = existingV["head_node"].([]interface{})
		existingWorkerNodes = existingV["worker_node"].([]interface{})
		existingZookeeperNodes = existingV["zookeeper_node"].([]interface{})
	}

	headNode := azure.FindHDInsightRole(input.Roles, "headnode")
	headNodes := azure.FlattenHDInsightNodeDefinition(headNode, existingHeadNodes, definition.HeadNodeDef)

	workerNode := azure.FindHDInsightRole(input.Roles, "workernode")
	workerNodes := azure.FlattenHDInsightNodeDefinition(workerNode, existingWorkerNodes, definition.WorkerNodeDef)

	zookeeperNode := azure.FindHDInsightRole(input.Roles, "zookeepernode")
	zookeeperNodes := azure.FlattenHDInsightNodeDefinition(zookeeperNode, existingZookeeperNodes, definition.ZookeeperNodeDef)

	result := map[string]interface{}{
		"head_node":      headNodes,
		"worker_node":    workerNodes,
		"zookeeper_node": zookeeperNodes,
	}

	if definition.EdgeNodeDef != nil {
		edgeNode := azure.FindHDInsightRole(input.Roles, "edgenode")
		edgeNodes := azure.FlattenHDInsightNodeDefinition(edgeNode, existingEdgeNodes, *definition.EdgeNodeDef)
		result["edge_node"] = edgeNodes
	}

	return []interface{}{
		result,
	}
}

func createHDInsightEdgeNodes(ctx context.Context, client *hdinsight.ApplicationsClient, resourceGroup string, name string, input map[string]interface{}) error {
	installScriptActions := expandHDInsightApplicationEdgeNodeInstallScriptActions(input["install_script_action"].([]interface{}))

	application := hdinsight.Application{
		Properties: &hdinsight.ApplicationProperties{
			ComputeProfile: &hdinsight.ComputeProfile{
				Roles: &[]hdinsight.Role{{
					Name: utils.String("edgenode"),
					HardwareProfile: &hdinsight.HardwareProfile{
						VMSize: utils.String(input["vm_size"].(string)),
					},
					TargetInstanceCount: utils.Int32(int32(input["target_instance_count"].(int))),
				}},
			},
			InstallScriptActions: installScriptActions,
			ApplicationType:      utils.String("CustomApplication"),
		},
	}
	future, err := client.Create(ctx, resourceGroup, name, name, application)
	if err != nil {
		return fmt.Errorf("Error creating edge nodes for HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of edge node for HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func deleteHDInsightEdgeNodes(ctx context.Context, client *hdinsight.ApplicationsClient, resourceGroup string, name string) error {
	future, err := client.Delete(ctx, resourceGroup, name, name)

	if err != nil {
		return fmt.Errorf("Error deleting edge nodes for HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of edge nodes for HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func expandHDInsightsMetastore(input []interface{}) map[string]interface{} {
	if len(input) == 0 || input[0] == nil {
		return map[string]interface{}{}
	}

	v := input[0].(map[string]interface{})
	config := map[string]interface{}{}

	if hiveRaw, ok := v["hive"]; ok {
		for k, val := range azure.ExpandHDInsightsHiveMetastore(hiveRaw.([]interface{})) {
			config[k] = val
		}
	}

	if oozieRaw, ok := v["oozie"]; ok {
		for k, val := range azure.ExpandHDInsightsOozieMetastore(oozieRaw.([]interface{})) {
			config[k] = val
		}
	}

	if ambariRaw, ok := v["ambari"]; ok {
		for k, val := range azure.ExpandHDInsightsAmbariMetastore(ambariRaw.([]interface{})) {
			config[k] = val
		}
	}

	return config
}

func flattenHDInsightsMetastores(d *schema.ResourceData, configurations map[string]map[string]*string) {
	result := map[string]interface{}{}

	hiveEnv, envExists := configurations["hive-env"]
	hiveSite, siteExists := configurations["hive-site"]
	if envExists && siteExists {
		result["hive"] = azure.FlattenHDInsightsHiveMetastore(hiveEnv, hiveSite)
	}

	oozieEnv, envExists := configurations["oozie-env"]
	oozieSite, siteExists := configurations["oozie-site"]
	if envExists && siteExists {
		result["oozie"] = azure.FlattenHDInsightsOozieMetastore(oozieEnv, oozieSite)
	}

	ambari, ambariExists := configurations["ambari-conf"]
	if ambariExists {
		result["ambari"] = azure.FlattenHDInsightsAmbariMetastore(ambari)
	}

	if len(result) > 0 {
		d.Set("metastores", []interface{}{
			result,
		})
	}
}
