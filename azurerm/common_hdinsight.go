package azurerm

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/preview/hdinsight/mgmt/2018-06-01-preview/hdinsight"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"log"
	"time"
)

func hdinsightClusterUpdate(clusterKind string, readFunc schema.ReadFunc) schema.UpdateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		client := meta.(*ArmClient).hdinsight.ClustersClient
		ctx := meta.(*ArmClient).StopContext

		id, err := parseAzureResourceID(d.Id())
		if err != nil {
			return err
		}

		resourceGroup := id.ResourceGroup
		name := id.Path["clusters"]

		if d.HasChange("tags") {
			tags := d.Get("tags").(map[string]interface{})
			params := hdinsight.ClusterPatchParameters{
				Tags: expandTags(tags),
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

		// The API can add an edge node but can't remove them without force newing the resource. We'll check for adding here
		// and can come back to removing if that functionality gets added. https://feedback.azure.com/forums/217335-hdinsight/suggestions/5663773-start-stop-cluster-hdinsight?page=3&per_page=20
		if clusterKind == "Hadoop" {
			if d.HasChange("roles.0.edge_node") {
				o, n := d.GetChange("roles.0.edge_node.#")
				edgeNodeRaw := d.Get("roles.0.edge_node").([]interface{})
				edgeNodeConfig := edgeNodeRaw[0].(map[string]interface{})
				applicationsClient := meta.(*ArmClient).hdinsight.ApplicationsClient

				// Create an edge node
				if o.(int) < n.(int) {
					err := createHDInsightEdgeNode(applicationsClient, ctx, resourceGroup, name, edgeNodeConfig)
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
				}
			}
		}

		return readFunc(d, meta)
	}
}

func hdinsightClusterDelete(clusterKind string) schema.DeleteFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		client := meta.(*ArmClient).hdinsight.ClustersClient
		ctx := meta.(*ArmClient).StopContext

		id, err := parseAzureResourceID(d.Id())
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

func createHDInsightEdgeNode(client hdinsight.ApplicationsClient, ctx context.Context, resourceGroup string, name string, input map[string]interface{}) error {
	installScriptActions := expandHDInsightApplicationEdgeNodeInstallScriptActions(input["install_script_action"].([]interface{}))

	application := hdinsight.Application{
		Properties: &hdinsight.ApplicationProperties{
			ComputeProfile: &hdinsight.ComputeProfile{
				Roles: &[]hdinsight.Role{{
					Name: utils.String("edgenode"),
					HardwareProfile: &hdinsight.HardwareProfile{
						VMSize: utils.String(input["vm_size"].(string)),
					},
					// The TargetInstanceCount must be one for edge nodes.
					TargetInstanceCount: utils.Int32(1),
				}},
			},
			InstallScriptActions: installScriptActions,
			ApplicationType:      utils.String("CustomApplication"),
		},
	}
	future, err := client.Create(ctx, resourceGroup, name, name, application)
	if err != nil {
		return fmt.Errorf("Error creating edge node for HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of edge node for HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func hdInsightWaitForReadyRefreshFunc(ctx context.Context, client hdinsight.ClustersClient, resourceGroupName string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, name)
		if err != nil {
			return nil, "Error", fmt.Errorf("Error issuing read request in relayNamespaceDeleteRefreshFunc to Relay Namespace %q (Resource Group %q): %s", name, resourceGroupName, err)
		}
		if props := res.Properties; props != nil {
			if state := props.ClusterState; state != nil {
				return res, *state, nil
			}
		}

		return res, "Pending", nil
	}
}
