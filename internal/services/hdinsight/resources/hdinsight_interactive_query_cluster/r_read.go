// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package hdinsight_interactive_query_cluster

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

func resourceHDInsightInteractiveQueryClusterRead(d *pluginsdk.ResourceData, meta any) error {
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
			log.Printf("[DEBUG] Interactive Query %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Interactive Query %s: %+v", id, err)
	}

	// Each call to configurationsClient methods is HTTP request. Getting all settings in one operation
	configurationsResp, err := configurationsClient.List(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Configuration for Interactive Query %s: %+v", id, err)
	}

	configurations := make(map[string]map[string]string)
	if model := configurationsResp.Model; model != nil && model.Configurations != nil {
		configurations = *model.Configurations
	}

	gateway, exists := configurations["gateway"]
	if !exists {
		return fmt.Errorf("retrieving Gateway Configuration for Interactive Query %s: %+v", id, err)
	}

	monitor, err := extensionsClient.GetMonitoringStatus(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Monitoring Status for Interactive Query %s: %+v", id, err)
	}

	extension, err := extensionsClient.GetAzureMonitorStatus(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Azure Monitor Status for Interactive Query %s: %+v", id, err)
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

			if err := d.Set("component_version", flattenHDInsightInteractiveQueryComponentVersion(props.ClusterDefinition.ComponentVersion)); err != nil {
				return fmt.Errorf("flattening `component_version`: %+v", err)
			}

			if err := d.Set("gateway", helpers.FlattenHDInsightsConfigurations(gateway, d)); err != nil {
				return fmt.Errorf("flattening `gateway`: %+v", err)
			}

			helpers.FlattenHDInsightsMetastores(d, configurations)

			if props.EncryptionInTransitProperties != nil {
				d.Set("encryption_in_transit_enabled", props.EncryptionInTransitProperties.IsEncryptionInTransitEnabled)
			}

			diskEncryptionProps, err := helpers.FlattenHDInsightsDiskEncryptionProperties(props.DiskEncryptionProperties)
			if err != nil {
				return err
			}
			if err := d.Set("disk_encryption", diskEncryptionProps); err != nil {
				return fmt.Errorf("flattening `disk_encryption`: %+v", err)
			}

			if err := d.Set("network", helpers.FlattenHDInsightsNetwork(props.NetworkProperties)); err != nil {
				return fmt.Errorf("flattening `network`: %+v", err)
			}

			if err := d.Set("private_link_configuration", helpers.FlattenHDInsightPrivateLinkConfigurations(props.PrivateLinkConfigurations)); err != nil {
				return fmt.Errorf("flattening `private_link_configuration`: %+v", err)
			}

			interactiveQueryRoles := helpers.HdInsightRoleDefinition{
				HeadNodeDef:      hdInsightInteractiveQueryClusterHeadNodeDefinition,
				WorkerNodeDef:    hdInsightInteractiveQueryClusterWorkerNodeDefinition,
				ZookeeperNodeDef: hdInsightInteractiveQueryClusterZookeeperNodeDefinition,
			}
			if err := d.Set("roles", helpers.FlattenHDInsightRoles(d, props.ComputeProfile, interactiveQueryRoles)); err != nil {
				return fmt.Errorf("flattening `roles`: %+v", err)
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

func flattenHDInsightInteractiveQueryComponentVersion(input *map[string]string) []any {
	output := make([]any, 0)
	if input != nil {
		interactiveHiveVersion := ""
		if v, ok := (*input)["InteractiveHive"]; ok {
			interactiveHiveVersion = v
		}
		output = append(output, map[string]any{
			"interactive_hive": interactiveHiveVersion,
		})
	}
	return output
}

// NOTE: this isn't a recommended way of building resources in Terraform
// this pattern is used to work around a generic but pedantic API endpoint
var hdInsightInteractiveQueryClusterHeadNodeDefinition = helpers.HDInsightNodeDefinition{
	MinInstanceCount:         2,
	MaxInstanceCount:         pointer.To(2),
	FixedTargetInstanceCount: pointer.To(int64(2)),
}

var hdInsightInteractiveQueryClusterWorkerNodeDefinition = helpers.HDInsightNodeDefinition{
	CanSpecifyInstanceCount: true,
	MinInstanceCount:        1,
	CanAutoScaleOnSchedule:  true,
}

var hdInsightInteractiveQueryClusterZookeeperNodeDefinition = helpers.HDInsightNodeDefinition{
	MinInstanceCount:         3,
	MaxInstanceCount:         pointer.To(3),
	FixedTargetInstanceCount: pointer.To(int64(3)),
}
