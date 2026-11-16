// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package hdinsight_hadoop_cluster

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/applications"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHDInsightHadoopClusterRead(d *pluginsdk.ResourceData, meta any) error {
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
			log.Printf("[DEBUG] Hadoop %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Hadoop %s: %+v", id, err)
	}

	// Each call to configurationsClient methods is HTTP request. Getting all settings in one operation
	configurationsResp, err := configurationsClient.List(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Configurations for HBase %s: %+v", id, err)
	}

	configurations := make(map[string]map[string]string)
	if model := configurationsResp.Model; model != nil && model.Configurations != nil {
		configurations = *model.Configurations
	}

	gateway, exists := configurations["gateway"]
	if !exists {
		return fmt.Errorf("retrieving Gateway Configuration for Hadoop %s: %+v", id, err)
	}

	applicationsClient := meta.(*clients.Client).HDInsight.Applications
	applicationId := applications.NewApplicationID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.ClusterName) // 2 id.ClusterName's are intentional
	edgeNode, err := applicationsClient.Get(ctx, applicationId)
	if err != nil {
		if !response.WasNotFound(edgeNode.HttpResponse) {
			return fmt.Errorf("reading Edge Node for Hadoop %s: %+v", id, err)
		}
	}

	extension, err := extensionsClient.GetAzureMonitorStatus(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Azure Monitor Status for Hadoop %s: %+v", id, err)
	}
	monitor, err := extensionsClient.GetMonitoringStatus(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Monitoring Status for Hadoop %s: %+v", id, err)
	}

	d.Set("name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		// storage_account isn't returned so I guess we just leave it ¯\_(ツ)_/¯
		if props := model.Properties; props != nil {
			d.Set("tier", pointer.FromEnum(props.Tier))
			d.Set("cluster_version", props.ClusterVersion)
			d.Set("tls_min_version", props.MinSupportedTlsVersion)

			if err := d.Set("component_version", flattenHDInsightHadoopComponentVersion(props.ClusterDefinition.ComponentVersion)); err != nil {
				return fmt.Errorf("flattening `component_version`: %+v", err)
			}

			if err := d.Set("gateway", helpers.FlattenHDInsightsConfigurations(gateway, d)); err != nil {
				return fmt.Errorf("flattening `gateway`: %+v", err)
			}

			helpers.FlattenHDInsightsMetastores(d, configurations)

			if err := d.Set("network", helpers.FlattenHDInsightsNetwork(props.NetworkProperties)); err != nil {
				return fmt.Errorf("flattening `network`: %+v", err)
			}

			if err := d.Set("private_link_configuration", helpers.FlattenHDInsightPrivateLinkConfigurations(props.PrivateLinkConfigurations)); err != nil {
				return fmt.Errorf("flattening `private_link_configuration`: %+v", err)
			}

			hadoopRoles := helpers.HdInsightRoleDefinition{
				HeadNodeDef:      hdInsightHadoopClusterHeadNodeDefinition,
				WorkerNodeDef:    hdInsightHadoopClusterWorkerNodeDefinition,
				ZookeeperNodeDef: hdInsightHadoopClusterZookeeperNodeDefinition,
			}
			flattenedRoles := helpers.FlattenHDInsightRoles(d, props.ComputeProfile, hadoopRoles)
			if edgeNode.Model != nil && edgeNode.Model.Properties != nil {
				flattenedRoles = flattenHDInsightEdgeNode(flattenedRoles, edgeNode.Model.Properties, d)
			}

			diskEncryptionProps, err := helpers.FlattenHDInsightsDiskEncryptionProperties(props.DiskEncryptionProperties)
			if err != nil {
				return err
			}
			if err := d.Set("disk_encryption", diskEncryptionProps); err != nil {
				return fmt.Errorf("flattening `disk_encryption`: %+v", err)
			}

			if err := d.Set("compute_isolation", helpers.FlattenHDInsightComputeIsolationProperties(props.ComputeIsolationProperties)); err != nil {
				return fmt.Errorf("failed setting `compute_isolation`: %+v", err)
			}

			if err := d.Set("roles", flattenedRoles); err != nil {
				return fmt.Errorf("flattening `roles`: %+v", err)
			}

			d.Set("https_endpoint", helpers.FindHDInsightConnectivityEndpoint("HTTPS", props.ConnectivityEndpoints))
			d.Set("ssh_endpoint", helpers.FindHDInsightConnectivityEndpoint("SSH", props.ConnectivityEndpoints))

			if err := d.Set("security_profile", helpers.FlattenHDInsightSecurityProfile(props.SecurityProfile, d)); err != nil {
				return fmt.Errorf("setting `security_profile`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	d.Set("extension", helpers.FlattenHDInsightAzureMonitor(extension.Model))

	d.Set("monitor", helpers.FlattenHDInsightMonitoring(monitor.Model))

	return nil
}

func flattenHDInsightEdgeNode(roles []any, props *applications.ApplicationProperties, d *pluginsdk.ResourceData) []any {
	if len(roles) == 0 || props == nil {
		return roles
	}

	role := roles[0].(map[string]any)

	edgeNode := make(map[string]any)
	if computeProfile := props.ComputeProfile; computeProfile != nil {
		if roles := computeProfile.Roles; roles != nil {
			for _, role := range *roles {
				if targetInstanceCount := role.TargetInstanceCount; targetInstanceCount != nil {
					edgeNode["target_instance_count"] = targetInstanceCount
				}
				if hardwareProfile := role.HardwareProfile; hardwareProfile != nil && hardwareProfile.VMSize != nil {
					vmSize := ""
					// the Azure API is inconsistent here, so rewrite this into the casing we expect
					for _, v := range validate.NodeDefinitionVMSize {
						if strings.EqualFold(v, *hardwareProfile.VMSize) {
							vmSize = v
						}
					}
					edgeNode["vm_size"] = vmSize
				}
			}
		}
	}

	installScriptActions := make([]any, 0)
	if props.InstallScriptActions != nil {
		for _, action := range *props.InstallScriptActions {
			installScriptActions = append(installScriptActions, map[string]any{
				"name":       action.Name,
				"uri":        action.Uri,
				"parameters": d.Get("roles.0.edge_node.0.install_script_action.0.parameters").(string),
			})
		}
	}

	uninstallScriptActions := make([]any, 0)
	if props.UninstallScriptActions != nil {
		for _, uninstallAction := range *props.UninstallScriptActions {
			uninstallScriptActions = append(uninstallScriptActions, map[string]any{
				"name":       uninstallAction.Name,
				"uri":        uninstallAction.Uri,
				"parameters": d.Get("roles.0.edge_node.0.uninstall_script_actions.0.parameters").(string),
			})
		}
	}

	httpsEndpoints := make([]any, 0)
	if HTTPSEndpoints := props.HTTPSEndpoints; HTTPSEndpoints != nil && len(*HTTPSEndpoints) != 0 {
		for _, HTTPSEndpoint := range *HTTPSEndpoints {
			httpsEndpoints = append(httpsEndpoints, map[string]any{
				"access_modes":         pointer.From(HTTPSEndpoint.AccessModes),
				"destination_port":     pointer.From(HTTPSEndpoint.DestinationPort),
				"disable_gateway_auth": pointer.From(HTTPSEndpoint.DisableGatewayAuth),
				"private_ip_address":   pointer.From(HTTPSEndpoint.PrivateIPAddress),
				"sub_domain_suffix":    pointer.From(HTTPSEndpoint.SubDomainSuffix),
			})
		}
	}

	role["edge_node"] = []any{
		map[string]any{
			"install_script_action":    installScriptActions,
			"https_endpoints":          httpsEndpoints,
			"uninstall_script_actions": uninstallScriptActions,
		},
	}

	return []any{role}
}

func flattenHDInsightHadoopComponentVersion(input *map[string]string) []any {
	output := make([]any, 0)

	if input != nil {
		hadoopVersion := ""
		if v, ok := (*input)["Hadoop"]; ok {
			hadoopVersion = v
		}
		output = append(output, map[string]any{
			"hadoop": hadoopVersion,
		})
	}

	return output
}

// NOTE: this isn't a recommended way of building resources in Terraform
// this pattern is used to work around a generic but pedantic API endpoint
var hdInsightHadoopClusterHeadNodeDefinition = helpers.HDInsightNodeDefinition{
	MinInstanceCount:         2,
	MaxInstanceCount:         pointer.To(2),
	FixedMinInstanceCount:    pointer.To(int64(1)),
	FixedTargetInstanceCount: pointer.To(int64(2)),
}

var hdInsightHadoopClusterWorkerNodeDefinition = helpers.HDInsightNodeDefinition{
	CanSpecifyInstanceCount: true,
	MinInstanceCount:        1,
	CanAutoScaleByCapacity:  true,
	CanAutoScaleOnSchedule:  true,
}

var hdInsightHadoopClusterZookeeperNodeDefinition = helpers.HDInsightNodeDefinition{
	MinInstanceCount:         3,
	MaxInstanceCount:         pointer.To(3),
	FixedMinInstanceCount:    pointer.To(int64(1)),
	FixedTargetInstanceCount: pointer.To(int64(3)),
}
