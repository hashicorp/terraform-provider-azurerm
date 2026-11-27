// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_spark_pool

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/bigdatapools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	synapseClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseSparkPoolUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := synapseClient.From(meta).SparkPoolClient

	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := bigdatapools.ParseBigDataPoolID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}
	props := existing.Model.Properties

	if d.HasChange("node_size_family") {
		props.NodeSizeFamily = pointer.ToEnum[bigdatapools.NodeSizeFamily](d.Get("node_size_family").(string))
	}

	if d.HasChange("node_size") {
		props.NodeSize = pointer.ToEnum[bigdatapools.NodeSize](d.Get("node_size").(string))
	}

	if d.HasChange("cache_size") {
		props.CacheSize = pointer.To(int64(d.Get("cache_size").(int)))
	}

	if d.HasChange("compute_isolation_enabled") {
		props.IsComputeIsolationEnabled = pointer.To(d.Get("compute_isolation_enabled").(bool))
	}

	if d.HasChanges("dynamic_executor_allocation_enabled", "min_executors", "max_executors") {
		if props.DynamicExecutorAllocation == nil {
			props.DynamicExecutorAllocation = &bigdatapools.DynamicExecutorAllocation{}
		}

		if d.HasChange("dynamic_executor_allocation_enabled") {
			props.DynamicExecutorAllocation.Enabled = pointer.To(d.Get("dynamic_executor_allocation_enabled").(bool))
		}

		if d.HasChange("min_executors") {
			props.DynamicExecutorAllocation.MinExecutors = pointer.To(int64(d.Get("min_executors").(int)))
		}

		if d.HasChange("max_executors") {
			props.DynamicExecutorAllocation.MaxExecutors = pointer.To(int64(d.Get("max_executors").(int)))
		}
	}

	if d.HasChanges("auto_scale", "node_count") {
		if d.HasChange("auto_scale") {
			props.AutoScale = expandSynapseSparkPoolAutoScaleProperties(d.Get("auto_scale").([]any))
		}

		if d.HasChange("node_count") {
			props.NodeCount = nil
			if !*props.AutoScale.Enabled {
				props.NodeCount = pointer.To(int64(d.Get("node_count").(int)))
			}
		}
	}

	if d.HasChange("auto_pause") {
		props.AutoPause = expandSynapseSparkPoolAutoPauseProperties(d.Get("auto_pause").([]any))
	}

	if d.HasChange("session_level_packages_enabled") {
		props.SessionLevelPackagesEnabled = pointer.To(d.Get("session_level_packages_enabled").(bool))
	}

	if d.HasChange("spark_config") {
		props.SparkConfigProperties = expandSynapseSparkPoolSparkConfig(d.Get("spark_config").([]any))
	}

	if d.HasChange("spark_events_folder") {
		props.SparkEventsFolder = pointer.To(d.Get("spark_events_folder").(string))
	}

	if d.HasChange("spark_log_folder") {
		props.DefaultSparkLogFolder = pointer.To(d.Get("spark_log_folder").(string))
	}

	if d.HasChange("library_requirement") {
		props.LibraryRequirements = expandSynapseSparkPoolLibraryRequirements(d.Get("library_requirement").([]any))
	}

	if d.HasChange("spark_version") {
		props.SparkVersion = pointer.To(d.Get("spark_version").(string))
	}

	if d.HasChange("tags") {
		existing.Model.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model, bigdatapools.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("updating %s: %v", id, err)
	}

	return resourceSynapseSparkPoolRead(d, meta)
}
