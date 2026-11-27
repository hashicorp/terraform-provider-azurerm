// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_spark_pool

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/bigdatapools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	synapseClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseSparkPoolCreate(d *pluginsdk.ResourceData, meta any) error {
	client := synapseClient.From(meta).SparkPoolClient
	workspaceClient := synapseClient.From(meta).WorkspacesClient
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	// todo 6.0 - move to the case-sensitive parser when validation.AsGeneratedID is removed: this parses a config
	// value which the paired AsGeneratedID validator accepts with legacy casing, and configs cannot be migrated.
	workspaceId, err := workspaces.ParseWorkspaceIDInsensitively(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	id := bigdatapools.NewBigDataPoolID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_synapse_spark_pool", id.ID())
		}
	}

	workspace, err := workspaceClient.Get(ctx, *workspaceId)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", workspaceId, err)
	}

	if workspace.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", workspaceId)
	}

	payload := bigdatapools.BigDataPoolResourceInfo{
		Location: workspace.Model.Location,
		Properties: &bigdatapools.BigDataPoolResourceProperties{
			AutoPause:                 expandSynapseSparkPoolAutoPauseProperties(d.Get("auto_pause").([]any)),
			AutoScale:                 expandSynapseSparkPoolAutoScaleProperties(d.Get("auto_scale").([]any)),
			CacheSize:                 pointer.To(int64(d.Get("cache_size").(int))),
			IsComputeIsolationEnabled: pointer.To(d.Get("compute_isolation_enabled").(bool)),
			DynamicExecutorAllocation: &bigdatapools.DynamicExecutorAllocation{
				Enabled:      pointer.To(d.Get("dynamic_executor_allocation_enabled").(bool)),
				MinExecutors: pointer.To(int64(d.Get("min_executors").(int))),
				MaxExecutors: pointer.To(int64(d.Get("max_executors").(int))),
			},
			DefaultSparkLogFolder:       pointer.To(d.Get("spark_log_folder").(string)),
			NodeSize:                    pointer.ToEnum[bigdatapools.NodeSize](d.Get("node_size").(string)),
			NodeSizeFamily:              pointer.ToEnum[bigdatapools.NodeSizeFamily](d.Get("node_size_family").(string)),
			SessionLevelPackagesEnabled: pointer.To(d.Get("session_level_packages_enabled").(bool)),
			SparkConfigProperties:       expandSynapseSparkPoolSparkConfig(d.Get("spark_config").([]any)),
			SparkEventsFolder:           pointer.To(d.Get("spark_events_folder").(string)),
			SparkVersion:                pointer.To(d.Get("spark_version").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if !*payload.Properties.AutoScale.Enabled {
		payload.Properties.NodeCount = pointer.To(int64(d.Get("node_count").(int)))
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, bigdatapools.DefaultCreateOrUpdateOperationOptions(), sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}
	d.SetId(id.ID())

	// Library Requirements can't be specified on Create so we'll make an additional request after we've confirmed the Spark Pool has been created.
	payload.Properties.LibraryRequirements = expandSynapseSparkPoolLibraryRequirements(d.Get("library_requirement").([]any))
	if err := client.CreateOrUpdateThenPoll(ctx, id, payload, bigdatapools.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	return resourceSynapseSparkPoolRead(d, meta)
}
