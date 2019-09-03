package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/hdinsight/mgmt/2018-06-01-preview/hdinsight"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func hdinsightClusterUpdate(clusterKind string, readFunc schema.ReadFunc) schema.UpdateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		client := meta.(*ArmClient).hdinsight.ClustersClient
		ctx := meta.(*ArmClient).StopContext

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

		if d.HasChange("roles") {
			log.Printf("[DEBUG] Resizing the HDInsight %q Cluster", clusterKind)
			rolesRaw := d.Get("roles").([]interface{})
			roles := rolesRaw[0].(map[string]interface{})
			headNodes := roles["worker_node"].([]interface{})
			headNode := headNodes[0].(map[string]interface{})
			targetInstanceCount := headNode["target_instance_count"].(int)
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

		return readFunc(d, meta)
	}
}

func hdinsightClusterDelete(clusterKind string) schema.DeleteFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		client := meta.(*ArmClient).hdinsight.ClustersClient
		ctx := meta.(*ArmClient).StopContext

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
