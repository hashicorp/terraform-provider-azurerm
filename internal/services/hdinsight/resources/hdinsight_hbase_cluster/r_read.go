// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package hdinsight_hbase_cluster

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHDInsightHBaseClusterRead(d *pluginsdk.ResourceData, meta any) error {
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
			log.Printf("[DEBUG] HBase %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving HBase %s: %+v", id, err)
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
		return fmt.Errorf("retrieving Gateway Configuration for HBase %s: %+v", id, err)
	}

	monitor, err := extensionsClient.GetMonitoringStatus(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Monitoring Status for HBase %s: %+v", id, err)
	}

	extension, err := extensionsClient.GetAzureMonitorStatus(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Azure Monitor Status for HBase %s: %+v", id, err)
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

			if err := d.Set("component_version", flattenHDInsightHBaseComponentVersion(props.ClusterDefinition.ComponentVersion)); err != nil {
				return fmt.Errorf("failure flattening `component_version`: %+v", err)
			}

			if err := d.Set("gateway", helpers.FlattenHDInsightsConfigurations(gateway, d)); err != nil {
				return fmt.Errorf("failure flattening `gateway`: %+v", err)
			}

			helpers.FlattenHDInsightsMetastores(d, configurations)

			if err := d.Set("network", helpers.FlattenHDInsightsNetwork(props.NetworkProperties)); err != nil {
				return fmt.Errorf("flattening `network`: %+v", err)
			}

			if err := d.Set("private_link_configuration", helpers.FlattenHDInsightPrivateLinkConfigurations(props.PrivateLinkConfigurations)); err != nil {
				return fmt.Errorf("flattening `private_link_configuration`: %+v", err)
			}

			diskEncryptionProps, err := helpers.FlattenHDInsightsDiskEncryptionProperties(props.DiskEncryptionProperties)
			if err != nil {
				return err
			}
			if err := d.Set("disk_encryption", diskEncryptionProps); err != nil {
				return fmt.Errorf("flattening `disk_encryption`: %+v", err)
			}

			hbaseRoles := helpers.HdInsightRoleDefinition{
				HeadNodeDef:      hdInsightHBaseClusterHeadNodeDefinition,
				WorkerNodeDef:    hdInsightHBaseClusterWorkerNodeDefinition,
				ZookeeperNodeDef: hdInsightHBaseClusterZookeeperNodeDefinition,
			}
			if err := d.Set("roles", helpers.FlattenHDInsightRoles(d, props.ComputeProfile, hbaseRoles)); err != nil {
				return fmt.Errorf("failure flattening `roles`: %+v", err)
			}

			if err := d.Set("compute_isolation", helpers.FlattenHDInsightComputeIsolationProperties(props.ComputeIsolationProperties)); err != nil {
				return fmt.Errorf("failed setting `compute_isolation`: %+v", err)
			}

			d.Set("https_endpoint", helpers.FindHDInsightConnectivityEndpoint("HTTPS", props.ConnectivityEndpoints))
			d.Set("ssh_endpoint", helpers.FindHDInsightConnectivityEndpoint("SSH", props.ConnectivityEndpoints))

			d.Set("monitor", helpers.FlattenHDInsightMonitoring(monitor.Model))

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

func flattenHDInsightHBaseComponentVersion(input *map[string]string) []any {
	output := make([]any, 0)
	if input != nil {
		hbaseVersion := ""
		if v, ok := (*input)["hbase"]; ok {
			hbaseVersion = v
		}
		output = append(output, map[string]any{
			"hbase": hbaseVersion,
		})
	}
	return output
}

// NOTE: this isn't a recommended way of building resources in Terraform
// this pattern is used to work around a generic but pedantic API endpoint
var hdInsightHBaseClusterHeadNodeDefinition = helpers.HDInsightNodeDefinition{
	MinInstanceCount:         2,
	MaxInstanceCount:         pointer.To(2),
	FixedTargetInstanceCount: pointer.To(int64(2)),
}

var hdInsightHBaseClusterWorkerNodeDefinition = helpers.HDInsightNodeDefinition{
	CanSpecifyInstanceCount: true,
	MinInstanceCount:        1,
	CanAutoScaleOnSchedule:  true,
}

var hdInsightHBaseClusterZookeeperNodeDefinition = helpers.HDInsightNodeDefinition{
	MinInstanceCount:         3,
	MaxInstanceCount:         pointer.To(3),
	FixedTargetInstanceCount: pointer.To(int64(3)),
}
