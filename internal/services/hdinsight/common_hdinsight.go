package hdinsight

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/applications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func hdinsightClusterUpdate(clusterKind string, readFunc pluginsdk.ReadFunc) pluginsdk.UpdateFunc {
	return func(d *pluginsdk.ResourceData, meta interface{}) error {
		client := meta.(*clients.Client).HDInsight.ClustersClient
		extensionsClient := meta.(*clients.Client).HDInsight.ExtensionsClient
		ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
		defer cancel()

		id, err := clusters.ParseClusterID(d.Id())
		if err != nil {
			return err
		}

		if d.HasChange("tags") {

			params := clusters.ClusterPatchParameters{
				tags.Expand(d.Get("tags").(map[string]interface{})),
			}
			if _, err := client.Update(ctx, *id, params); err != nil {
				return fmt.Errorf("updating Tags for %s: %+v", *id, err)
			}
		}

		if d.HasChange("roles.0.worker_node") {
			log.Printf("[DEBUG] Resizing the HDInsight %q Cluster", clusterKind)
			rolesRaw := d.Get("roles").([]interface{})
			roles := rolesRaw[0].(map[string]interface{})
			workerNodes := roles["worker_node"].([]interface{})
			workerNode := workerNodes[0].(map[string]interface{})
			if d.HasChange("roles.0.worker_node.0.target_instance_count") {
				targetInstanceCount := workerNode["target_instance_count"].(int)
				params := clusters.ClusterResizeParameters{
					TargetInstanceCount: pointer.To(int64(targetInstanceCount)),
				}

				if err = client.ResizeThenPoll(ctx, *id, params); err != nil {
					return fmt.Errorf("resizing %s: %+v", *id, err)
				}
			}

			if d.HasChange("roles.0.worker_node.0.autoscale") {
				autoscale := ExpandHDInsightNodeAutoScaleDefinition(workerNode["autoscale"].([]interface{}))
				params := clusters.AutoscaleConfigurationUpdateParameter{
					Autoscale: autoscale,
				}

				if err := client.UpdateAutoScaleConfigurationThenPoll(ctx, *id, params); err != nil {
					return fmt.Errorf("changing autoscale of the %s: %+v", *id, err)
				}
			}
		}

		// The API can add an edge node but can't remove them without force newing the pluginsdk. We'll check for adding here
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
					err := deleteHDInsightEdgeNodes(ctx, applicationsClient, *id)
					if err != nil {
						return err
					}
				}

				if newEdgeNodeInt != 0 {
					err = createHDInsightEdgeNodes(ctx, applicationsClient, *id, edgeNodeConfig)
					if err != nil {
						return err
					}
				}

				// we can't rely on the use of the Future here due to the node being successfully completed but now the cluster is applying those changes.
				log.Printf("[DEBUG] Waiting for %s to finish applying edge node", id)
				stateConf := &pluginsdk.StateChangeConf{
					Pending:    []string{"AzureVMConfiguration", "Accepted", "HdInsightConfiguration"},
					Target:     []string{"Running"},
					Refresh:    hdInsightWaitForReadyRefreshFunc(ctx, client, *id),
					MinTimeout: 15 * time.Second,
					Timeout:    d.Timeout(pluginsdk.TimeoutUpdate),
				}

				if _, err := stateConf.WaitForStateContext(ctx); err != nil {
					return fmt.Errorf("waiting for %s to be running: %s", id, err)
				}
			}
		}

		if d.HasChange("monitor") {
			log.Printf("[DEBUG] Change Azure Monitor for the HDInsight %q Cluster", clusterKind)
			if v, ok := d.GetOk("monitor"); ok {
				monitorRaw := v.([]interface{})
				if err := enableHDInsightMonitoring(ctx, extensionsClient, *id, monitorRaw); err != nil {
					return err
				}
			} else if err := disableHDInsightMonitoring(ctx, extensionsClient, *id); err != nil {
				return err
			}
		}
		if d.HasChange("extension") {
			log.Printf("[DEBUG] Change Azure Monitor for the HDInsight %q Cluster", clusterKind)
			if v, ok := d.GetOk("extension"); ok {
				extensionRaw := v.([]interface{})
				if err := enableHDInsightAzureMonitor(ctx, extensionsClient, *id, extensionRaw); err != nil {
					return err
				}
			} else if err := disableHDInsightAzureMonitor(ctx, extensionsClient, *id); err != nil {
				return err
			}
		}
		if d.HasChange("gateway") {
			log.Printf("[DEBUG] Updating the HDInsight %q Cluster gateway", clusterKind)
			vs := d.Get("gateway").([]interface{})[0].(map[string]interface{})

			if err := client.UpdateGatewaySettingsThenPoll(ctx, *id, clusters.UpdateGatewaySettingsParameters{
				RestAuthCredentialIsEnabled: pointer.To(true),
				RestAuthCredentialUsername:  pointer.To(vs["username"].(string)),
				RestAuthCredentialPassword:  pointer.To(vs["password"].(string)),
			}); err != nil {
				return fmt.Errorf("updating %s: %s", id, err)
			}
		}

		return readFunc(d, meta)
	}
}

func hdinsightClusterDelete(clusterKind string) pluginsdk.DeleteFunc {
	return func(d *pluginsdk.ResourceData, meta interface{}) error {
		client := meta.(*clients.Client).HDInsight.ClustersClient
		ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
		defer cancel()

		id, err := clusters.ParseClusterID(d.Id())
		if err != nil {
			return err
		}

		if err := client.DeleteThenPoll(ctx, *id); err != nil {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}

		return nil
	}
}

type hdInsightRoleDefinition struct {
	HeadNodeDef            HDInsightNodeDefinition
	WorkerNodeDef          HDInsightNodeDefinition
	ZookeeperNodeDef       HDInsightNodeDefinition
	KafkaManagementNodeDef *HDInsightNodeDefinition
	EdgeNodeDef            *HDInsightNodeDefinition
}

func expandHDInsightRoles(input []interface{}, definition hdInsightRoleDefinition) (*[]clusters.Role, error) {
	v := input[0].(map[string]interface{})

	headNodeRaw := v["head_node"].([]interface{})
	headNode, err := ExpandHDInsightNodeDefinition("headnode", headNodeRaw, definition.HeadNodeDef)
	if err != nil {
		return nil, fmt.Errorf("expanding `head_node`: %+v", err)
	}

	workerNodeRaw := v["worker_node"].([]interface{})
	workerNode, err := ExpandHDInsightNodeDefinition("workernode", workerNodeRaw, definition.WorkerNodeDef)
	if err != nil {
		return nil, fmt.Errorf("expanding `worker_node`: %+v", err)
	}

	zookeeperNodeRaw := v["zookeeper_node"].([]interface{})
	zookeeperNode, err := ExpandHDInsightNodeDefinition("zookeepernode", zookeeperNodeRaw, definition.ZookeeperNodeDef)
	if err != nil {
		return nil, fmt.Errorf("expanding `zookeeper_node`: %+v", err)
	}

	roles := []clusters.Role{
		*headNode,
		*workerNode,
		*zookeeperNode,
	}

	if definition.EdgeNodeDef != nil {
		edgeNodeRaw := v["edge_node"].([]interface{})
		edgeNode, err := ExpandHDInsightNodeDefinition("edgenode", edgeNodeRaw, *definition.EdgeNodeDef)
		if err != nil {
			return nil, fmt.Errorf("expanding `edge_node`: %+v", err)
		}
		roles = append(roles, *edgeNode)
	}

	if definition.KafkaManagementNodeDef != nil {
		kafkaManagementNodeRaw := v["kafka_management_node"].([]interface{})
		// "kafka_management_node" is optional, we expand it only when user has specified it.
		if len(kafkaManagementNodeRaw) != 0 {
			kafkaManagementNode, err := ExpandHDInsightNodeDefinition("kafkamanagementnode", kafkaManagementNodeRaw, *definition.KafkaManagementNodeDef)
			if err != nil {
				return nil, fmt.Errorf("expanding `kafka_management_node`: %+v", err)
			}
			roles = append(roles, *kafkaManagementNode)
		}
	}

	return &roles, nil
}

func flattenHDInsightRoles(d *pluginsdk.ResourceData, input *clusters.ComputeProfile, definition hdInsightRoleDefinition) []interface{} {
	if input == nil || input.Roles == nil {
		return []interface{}{}
	}

	var existingKafkaManagementNodes, existingEdgeNodes, existingHeadNodes, existingWorkerNodes, existingZookeeperNodes []interface{}

	existingVs := d.Get("roles").([]interface{})
	if len(existingVs) > 0 {
		existingV := existingVs[0].(map[string]interface{})

		if definition.EdgeNodeDef != nil {
			existingEdgeNodes = existingV["edge_node"].([]interface{})
		}

		if definition.KafkaManagementNodeDef != nil {
			existingKafkaManagementNodes = existingV["kafka_management_node"].([]interface{})
		}

		existingHeadNodes = existingV["head_node"].([]interface{})
		existingWorkerNodes = existingV["worker_node"].([]interface{})
		existingZookeeperNodes = existingV["zookeeper_node"].([]interface{})
	}

	headNode := FindHDInsightRole(input.Roles, "headnode")
	headNodes := FlattenHDInsightNodeDefinition(headNode, existingHeadNodes, definition.HeadNodeDef)

	workerNode := FindHDInsightRole(input.Roles, "workernode")
	workerNodes := FlattenHDInsightNodeDefinition(workerNode, existingWorkerNodes, definition.WorkerNodeDef)

	zookeeperNode := FindHDInsightRole(input.Roles, "zookeepernode")
	zookeeperNodes := FlattenHDInsightNodeDefinition(zookeeperNode, existingZookeeperNodes, definition.ZookeeperNodeDef)

	result := map[string]interface{}{
		"head_node":      headNodes,
		"worker_node":    workerNodes,
		"zookeeper_node": zookeeperNodes,
	}

	if definition.EdgeNodeDef != nil {
		edgeNode := FindHDInsightRole(input.Roles, "edgenode")
		edgeNodes := FlattenHDInsightNodeDefinition(edgeNode, existingEdgeNodes, *definition.EdgeNodeDef)
		result["edge_node"] = edgeNodes
	}

	if definition.KafkaManagementNodeDef != nil {
		kafkaManagementNode := FindHDInsightRole(input.Roles, "kafkamanagementnode")
		kafkaManagementNodes := FlattenHDInsightNodeDefinition(kafkaManagementNode, existingKafkaManagementNodes, *definition.KafkaManagementNodeDef)
		result["kafka_management_node"] = kafkaManagementNodes
	}

	return []interface{}{
		result,
	}
}

func createHDInsightEdgeNodes(ctx context.Context, client *applications.ApplicationsClient, id clusters.ClusterId, input map[string]interface{}) error {
	installScriptActions := expandHDInsightApplicationEdgeNodeInstallScriptActions(input["install_script_action"].([]interface{}))

	appId := applications.NewApplicationID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.ClusterName)

	application := applications.Application{
		Properties: &applications.ApplicationProperties{
			ComputeProfile: &applications.ComputeProfile{
				Roles: &[]applications.Role{{
					Name: pointer.To("edgenode"),
					HardwareProfile: &applications.HardwareProfile{
						VMSize: pointer.To(input["vm_size"].(string)),
					},
					TargetInstanceCount: pointer.To(int64(input["target_instance_count"].(int))),
				}},
			},
			InstallScriptActions: installScriptActions,
			ApplicationType:      pointer.To("CustomApplication"),
		},
	}

	if v, ok := input["https_endpoints"]; ok {
		application.Properties.HTTPSEndpoints = expandHDInsightApplicationEdgeNodeHttpsEndpoints(v.([]interface{}))
	}

	if v, ok := input["uninstall_script_actions"]; ok {
		application.Properties.UninstallScriptActions = expandHDInsightApplicationEdgeNodeUninstallScriptActions(v.([]interface{}))
	}

	if err := client.CreateThenPoll(ctx, appId, application); err != nil {
		return fmt.Errorf("creating edge nodes for %s: %+v", appId, err)
	}

	return nil
}

func deleteHDInsightEdgeNodes(ctx context.Context, client *applications.ApplicationsClient, id clusters.ClusterId) error {
	appId := applications.NewApplicationID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.ClusterName)
	if err := client.DeleteThenPoll(ctx, appId); err != nil {
		return fmt.Errorf("deleting edge nodes for %s: %+v", appId, err)
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
		for k, val := range expandHDInsightsHiveMetastore(hiveRaw.([]interface{})) {
			config[k] = val
		}
	}

	if oozieRaw, ok := v["oozie"]; ok {
		for k, val := range expandHDInsightsOozieMetastore(oozieRaw.([]interface{})) {
			config[k] = val
		}
	}

	if ambariRaw, ok := v["ambari"]; ok {
		for k, val := range expandHDInsightsAmbariMetastore(ambariRaw.([]interface{})) {
			config[k] = val
		}
	}

	return config
}

func flattenHDInsightsMetastores(d *pluginsdk.ResourceData, configurations map[string]map[string]*string) {
	result := map[string]interface{}{}

	hiveEnv, envExists := configurations["hive-env"]
	hiveSite, siteExists := configurations["hive-site"]
	if envExists && siteExists {
		result["hive"] = FlattenHDInsightsHiveMetastore(hiveEnv, hiveSite)
	}

	oozieEnv, envExists := configurations["oozie-env"]
	oozieSite, siteExists := configurations["oozie-site"]
	if envExists && siteExists {
		result["oozie"] = FlattenHDInsightsOozieMetastore(oozieEnv, oozieSite)
	}

	ambari, ambariExists := configurations["ambari-conf"]
	if ambariExists {
		result["ambari"] = FlattenHDInsightsAmbariMetastore(ambari)
	}

	if len(result) > 0 {
		d.Set("metastores", []interface{}{
			result,
		})
	}
}

func flattenHDInsightMonitoring(monitor *extensions.ClusterMonitoringResponse) []interface{} {
	if monitor != nil && monitor.ClusterMonitoringEnabled != nil {
		if *monitor.ClusterMonitoringEnabled {
			return []interface{}{
				map[string]string{
					"log_analytics_workspace_id": *monitor.WorkspaceId,
					"primary_key":                "*****",
				},
			}
		}
	}
	return nil
}

func flattenHDInsightAzureMonitor(extension *extensions.AzureMonitorResponse) []interface{} {
	if extension != nil && extension.ClusterMonitoringEnabled != nil {
		if *extension.ClusterMonitoringEnabled {
			return []interface{}{
				map[string]string{
					"log_analytics_workspace_id": *extension.WorkspaceId,
					"primary_key":                "*****",
				},
			}
		}
	}
	return nil
}

func enableHDInsightMonitoring(ctx context.Context, client *extensions.ExtensionsClient, id clusters.ClusterId, input []interface{}) error {
	extID := extensions.NewClusterID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName)

	if err := client.EnableMonitoringThenPoll(ctx, extID, expandHDInsightsMonitor(input)); err != nil {
		return fmt.Errorf("enabling monitor for  %s: %+v", extID, err)
	}

	return nil
}

func disableHDInsightMonitoring(ctx context.Context, client *extensions.ExtensionsClient, id clusters.ClusterId) error {
	extID := extensions.NewClusterID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName)

	if err := client.DisableMonitoringThenPoll(ctx, extID); err != nil {
		return fmt.Errorf("waiting for disabling monitor for %s: %+v", extID, err)
	}

	return nil
}

func enableHDInsightAzureMonitor(ctx context.Context, client *extensions.ExtensionsClient, id clusters.ClusterId, input []interface{}) error {
	v := input[0].(map[string]interface{})

	extID := extensions.NewClusterID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName)

	extension := extensions.AzureMonitorRequest{
		WorkspaceId: pointer.To(v["log_analytics_workspace_id"].(string)),
		PrimaryKey:  pointer.To(v["primary_key"].(string)),
	}

	if err := client.EnableAzureMonitorThenPoll(ctx, extID, extension); err != nil {
		return fmt.Errorf("creating extension for %s: %+v", extID, err)
	}

	return nil
}

func disableHDInsightAzureMonitor(ctx context.Context, client *extensions.ExtensionsClient, id clusters.ClusterId) error {
	extID := extensions.NewClusterID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName)

	if err := client.DisableAzureMonitorThenPoll(ctx, extID); err != nil {
		return fmt.Errorf("disabling extension for %s): %+v", extID, err)
	}

	return nil
}
