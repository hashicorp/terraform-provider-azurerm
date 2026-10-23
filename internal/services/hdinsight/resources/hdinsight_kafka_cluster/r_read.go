// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package hdinsight_kafka_cluster

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHDInsightKafkaClusterRead(d *pluginsdk.ResourceData, meta any) error {
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
			log.Printf("[DEBUG] Kafka %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Kafka %s: %+v", id, err)
	}

	// Each call to configurationsClient methods is HTTP request. Getting all settings in one operation
	configurationsResp, err := configurationsClient.List(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Configurations for Kafka %s: %+v", id, err)
	}

	configurations := make(map[string]map[string]string)
	if model := configurationsResp.Model; model != nil && model.Configurations != nil {
		configurations = *model.Configurations
	}
	gateway, exists := configurations["gateway"]
	if !exists {
		return fmt.Errorf("retrieving Gateway Configuration Kafka %s: %+v", id, err)
	}

	monitor, err := extensionsClient.GetMonitoringStatus(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Monitoring Status for Kafka %s: %+v", id, err)
	}

	extension, err := extensionsClient.GetAzureMonitorStatus(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Azure Monitor Status for Kafka %s: %+v", id, err)
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

			if err := d.Set("component_version", flattenHDInsightKafkaComponentVersion(props.ClusterDefinition.ComponentVersion)); err != nil {
				return fmt.Errorf("failure flattening `component_version`: %+v", err)
			}

			if err := d.Set("gateway", helpers.FlattenHDInsightsConfigurations(gateway, d)); err != nil {
				return fmt.Errorf("failure flattening `gateway`: %+v", err)
			}

			helpers.FlattenHDInsightsMetastores(d, configurations)

			kafkaRoles := helpers.HdInsightRoleDefinition{
				HeadNodeDef:            hdInsightKafkaClusterHeadNodeDefinition,
				WorkerNodeDef:          hdInsightKafkaClusterWorkerNodeDefinition,
				ZookeeperNodeDef:       hdInsightKafkaClusterZookeeperNodeDefinition,
				KafkaManagementNodeDef: &hdInsightKafkaClusterKafkaManagementNodeDefinition,
			}
			if err := d.Set("roles", helpers.FlattenHDInsightRoles(d, props.ComputeProfile, kafkaRoles)); err != nil {
				return fmt.Errorf("failure flattening `roles`: %+v", err)
			}

			d.Set("https_endpoint", helpers.FindHDInsightConnectivityEndpoint("HTTPS", props.ConnectivityEndpoints))
			d.Set("ssh_endpoint", helpers.FindHDInsightConnectivityEndpoint("SSH", props.ConnectivityEndpoints))
			d.Set("kafka_rest_proxy_endpoint", helpers.FindHDInsightConnectivityEndpoint("KafkaRestProxyPublicEndpoint", props.ConnectivityEndpoints))

			if props.EncryptionInTransitProperties != nil {
				d.Set("encryption_in_transit_enabled", props.EncryptionInTransitProperties.IsEncryptionInTransitEnabled)
			}

			if err := d.Set("network", helpers.FlattenHDInsightsNetwork(props.NetworkProperties)); err != nil {
				return fmt.Errorf("flatten `network`: %+v", err)
			}
			if err := d.Set("private_link_configuration", helpers.FlattenHDInsightPrivateLinkConfigurations(props.PrivateLinkConfigurations)); err != nil {
				return fmt.Errorf("flattening `private_link_configuration`: %+v", err)
			}
			if err := d.Set("compute_isolation", helpers.FlattenHDInsightComputeIsolationProperties(props.ComputeIsolationProperties)); err != nil {
				return fmt.Errorf("failed setting `compute_isolation`: %+v", err)
			}

			diskEncryptionProps, err := helpers.FlattenHDInsightsDiskEncryptionProperties(props.DiskEncryptionProperties)
			if err != nil {
				return err
			}
			if err := d.Set("disk_encryption", diskEncryptionProps); err != nil {
				return fmt.Errorf("flattening `disk_encryption`: %+v", err)
			}

			d.Set("monitor", helpers.FlattenHDInsightMonitoring(monitor.Model))

			if err = d.Set("rest_proxy", flattenKafkaRestProxyProperty(props.KafkaRestProperties)); err != nil {
				return fmt.Errorf("setting `rest_proxy`: %+v", err)
			}

			d.Set("extension", helpers.FlattenHDInsightAzureMonitor(extension.Model))

			if err := d.Set("security_profile", helpers.FlattenHDInsightSecurityProfile(props.SecurityProfile, d)); err != nil {
				return fmt.Errorf("setting `security_profile`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}

func flattenHDInsightKafkaComponentVersion(input *map[string]string) []any {
	output := make([]any, 0)
	if input != nil {
		kafkaVersion := ""
		if v, ok := (*input)["kafka"]; ok {
			kafkaVersion = v
		}
		output = append(output, map[string]any{
			"kafka": kafkaVersion,
		})
	}
	return output
}

func flattenKafkaRestProxyProperty(input *clusters.KafkaRestProperties) []any {
	if input == nil || input.ClientGroupInfo == nil {
		return []any{}
	}

	groupInfo := input.ClientGroupInfo

	return []any{
		map[string]any{
			"security_group_id":   pointer.From(groupInfo.GroupId),
			"security_group_name": pointer.From(groupInfo.GroupName),
		},
	}
}

// NOTE: this isn't a recommended way of building resources in Terraform
// this pattern is used to work around a generic but pedantic API endpoint
var hdInsightKafkaClusterHeadNodeDefinition = helpers.HDInsightNodeDefinition{
	MinInstanceCount:         2,
	MaxInstanceCount:         pointer.To(2),
	FixedTargetInstanceCount: pointer.To(int64(2)),
}

var hdInsightKafkaClusterWorkerNodeDefinition = helpers.HDInsightNodeDefinition{
	CanSpecifyInstanceCount: true,
	MinInstanceCount:        1,
	CanSpecifyDisks:         true,
	MaxNumberOfDisksPerNode: pointer.To(8),
}

var hdInsightKafkaClusterZookeeperNodeDefinition = helpers.HDInsightNodeDefinition{
	MinInstanceCount:         3,
	MaxInstanceCount:         pointer.To(3),
	FixedTargetInstanceCount: pointer.To(int64(3)),
}

var hdInsightKafkaClusterKafkaManagementNodeDefinition = helpers.HDInsightNodeDefinition{
	MinInstanceCount:         2,
	FixedTargetInstanceCount: pointer.To(int64(2)),
}
