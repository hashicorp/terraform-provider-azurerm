// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_spark_pool

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/bigdatapools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	synapseClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseSparkPoolRead(d *pluginsdk.ResourceData, meta any) error {
	client := synapseClient.From(meta).SparkPoolClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := bigdatapools.ParseBigDataPoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.BigDataPoolName)
	d.Set("synapse_workspace_id", workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID())

	if resp.Model != nil {
		if props := resp.Model.Properties; props != nil {
			if err := d.Set("auto_pause", flattenSynapseSparkPoolAutoPauseProperties(props.AutoPause)); err != nil {
				return fmt.Errorf("setting `auto_pause`: %+v", err)
			}
			if err := d.Set("auto_scale", flattenSynapseSparkPoolAutoScaleProperties(props.AutoScale)); err != nil {
				return fmt.Errorf("setting `auto_scale`: %+v", err)
			}
			if err := d.Set("library_requirement", flattenSynapseSparkPoolLibraryRequirements(props.LibraryRequirements)); err != nil {
				return fmt.Errorf("setting `library_requirement`: %+v", err)
			}
			d.Set("cache_size", props.CacheSize)
			d.Set("compute_isolation_enabled", props.IsComputeIsolationEnabled)

			var dynamicExecutorAllocationEnabled bool
			var minExecutor, maxExecutor int64
			if props.DynamicExecutorAllocation != nil {
				dynamicExecutorAllocationEnabled = pointer.From(props.DynamicExecutorAllocation.Enabled)
				minExecutor = pointer.From(props.DynamicExecutorAllocation.MinExecutors)
				maxExecutor = pointer.From(props.DynamicExecutorAllocation.MaxExecutors)
			}
			d.Set("dynamic_executor_allocation_enabled", dynamicExecutorAllocationEnabled)
			d.Set("min_executors", minExecutor)
			d.Set("max_executors", maxExecutor)

			d.Set("node_count", props.NodeCount)
			d.Set("node_size", props.NodeSize)
			d.Set("node_size_family", pointer.FromEnum(props.NodeSizeFamily))
			d.Set("session_level_packages_enabled", props.SessionLevelPackagesEnabled)
			d.Set("spark_config", flattenSynapseSparkPoolSparkConfig(props.SparkConfigProperties))
			d.Set("spark_version", props.SparkVersion)
		}

		if err := tags.FlattenAndSet(d, resp.Model.Tags); err != nil {
			return fmt.Errorf("flattening `tags`: %+v", err)
		}
	}

	return nil
}
