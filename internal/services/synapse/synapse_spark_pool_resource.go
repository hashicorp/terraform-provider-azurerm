// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/bigdatapools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseSparkPool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseSparkPoolCreate,
		Read:   resourceSynapseSparkPoolRead,
		Update: resourceSynapseSparkPoolUpdate,
		Delete: resourceSynapseSparkPoolDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := bigdatapools.ParseBigDataPoolID(id)
			return err
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SynapseSparkPoolV0ToV1{},
		}),

		SchemaVersion: 1,
		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SparkPoolName,
			},

			"synapse_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.AsGeneratedID(workspaces.ParseWorkspaceIDInsensitively),
			},

			"node_size_family": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(bigdatapools.PossibleValuesForNodeSizeFamily(), false),
			},

			"node_size": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(bigdatapools.PossibleValuesForNodeSize(), false),
			},

			"cache_size": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
			},

			"compute_isolation_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"dynamic_executor_allocation_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"min_executors": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 200),
			},

			"max_executors": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 200),
			},

			"node_count": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				// NOTE: O+C There is a bug in the API where this gets set when auto_scale is enabled resulting in a diff
				Computed:     true,
				ValidateFunc: validation.IntBetween(3, 200),
				ExactlyOneOf: []string{"node_count", "auto_scale"},
			},

			"auto_scale": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: []string{"node_count", "auto_scale"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"min_node_count": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(3, 200),
						},

						"max_node_count": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(3, 200),
						},
					},
				},
			},

			"auto_pause": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"delay_in_minutes": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(5, 10080),
						},
					},
				},
			},

			"session_level_packages_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"spark_config": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"content": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"filename": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"spark_events_folder": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "/events",
			},

			"spark_log_folder": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "/logs",
			},

			"library_requirement": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"content": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"filename": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
			},

			"spark_version": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"3.4",
					"3.5",
				}, false),
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceSynapseSparkPoolCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.SparkPoolClient
	workspaceClient := meta.(*clients.Client).Synapse.WorkspacesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// todo 6.0 - move to the case-sensitive parser when validation.AsGeneratedID is removed: this parses a config
	// value which the paired AsGeneratedID validator accepts with legacy casing, and configs cannot be migrated.
	workspaceId, err := workspaces.ParseWorkspaceIDInsensitively(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	id := bigdatapools.NewBigDataPoolID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
			AutoPause:                 expandSynapseSparkPoolAutoPauseProperties(d.Get("auto_pause").([]interface{})),
			AutoScale:                 expandSynapseSparkPoolAutoScaleProperties(d.Get("auto_scale").([]interface{})),
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
			SparkConfigProperties:       expandSynapseSparkPoolSparkConfig(d.Get("spark_config").([]interface{})),
			SparkEventsFolder:           pointer.To(d.Get("spark_events_folder").(string)),
			SparkVersion:                pointer.To(d.Get("spark_version").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if !*payload.Properties.AutoScale.Enabled {
		payload.Properties.NodeCount = pointer.To(int64(d.Get("node_count").(int)))
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, bigdatapools.DefaultCreateOrUpdateOperationOptions(), sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}
	d.SetId(id.ID())

	// Library Requirements can't be specified on Create so we'll make an additional request after we've confirmed the Spark Pool has been created.
	payload.Properties.LibraryRequirements = expandSynapseSparkPoolLibraryRequirements(d.Get("library_requirement").([]interface{}))
	if err := client.CreateOrUpdateThenPoll(ctx, id, payload, bigdatapools.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	return resourceSynapseSparkPoolRead(d, meta)
}

func resourceSynapseSparkPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.SparkPoolClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
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

func resourceSynapseSparkPoolUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.SparkPoolClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := bigdatapools.ParseBigDataPoolID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model != nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties != nil {
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
			props.AutoScale = expandSynapseSparkPoolAutoScaleProperties(d.Get("auto_scale").([]interface{}))
		}

		if d.HasChange("node_count") {
			props.NodeCount = nil
			if !*props.AutoScale.Enabled {
				props.NodeCount = pointer.To(int64(d.Get("node_count").(int)))
			}
		}
	}

	if d.HasChange("auto_pause") {
		props.AutoPause = expandSynapseSparkPoolAutoPauseProperties(d.Get("auto_pause").([]interface{}))
	}

	if d.HasChange("session_level_packages_enabled") {
		props.SessionLevelPackagesEnabled = pointer.To(d.Get("session_level_packages_enabled").(bool))
	}

	if d.HasChange("spark_config") {
		props.SparkConfigProperties = expandSynapseSparkPoolSparkConfig(d.Get("spark_config").([]interface{}))
	}

	if d.HasChange("spark_events_folder") {
		props.SparkEventsFolder = pointer.To(d.Get("spark_events_folder").(string))
	}

	if d.HasChange("spark_log_folder") {
		props.DefaultSparkLogFolder = pointer.To(d.Get("spark_log_folder").(string))
	}

	if d.HasChange("library_requirements") {
		props.LibraryRequirements = expandSynapseSparkPoolLibraryRequirements(d.Get("library_requirement").([]interface{}))
	}

	if d.HasChange("spark_version") {
		props.SparkVersion = pointer.To(d.Get("spark_version").(string))
	}

	if d.HasChange("tags") {
		existing.Model.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model, bigdatapools.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("updating %s: %v", id, err)
	}

	return resourceSynapseSparkPoolRead(d, meta)
}

func resourceSynapseSparkPoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.SparkPoolClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := bigdatapools.ParseBigDataPoolID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandSynapseSparkPoolAutoPauseProperties(input []interface{}) *bigdatapools.AutoPauseProperties {
	if len(input) == 0 {
		return &bigdatapools.AutoPauseProperties{
			Enabled: pointer.To(false),
		}
	}
	v := input[0].(map[string]interface{})
	return &bigdatapools.AutoPauseProperties{
		DelayInMinutes: pointer.To(int64(v["delay_in_minutes"].(int))),
		Enabled:        pointer.To(true),
	}
}

func expandSynapseSparkPoolAutoScaleProperties(input []interface{}) *bigdatapools.AutoScaleProperties {
	if len(input) == 0 || input[0] == nil {
		return &bigdatapools.AutoScaleProperties{
			Enabled: pointer.To(false),
		}
	}
	v := input[0].(map[string]interface{})
	return &bigdatapools.AutoScaleProperties{
		MinNodeCount: pointer.To(int64(v["min_node_count"].(int))),
		Enabled:      pointer.To(true),
		MaxNodeCount: pointer.To(int64(v["max_node_count"].(int))),
	}
}

func expandSynapseSparkPoolLibraryRequirements(input []interface{}) *bigdatapools.LibraryRequirements {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &bigdatapools.LibraryRequirements{
		Content:  pointer.To(v["content"].(string)),
		Filename: pointer.To(v["filename"].(string)),
	}
}

func expandSynapseSparkPoolSparkConfig(input []interface{}) *bigdatapools.SparkConfigProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	value := input[0].(map[string]interface{})
	return &bigdatapools.SparkConfigProperties{
		Content:  pointer.To(value["content"].(string)),
		Filename: pointer.To(value["filename"].(string)),
	}
}

func flattenSynapseSparkPoolAutoPauseProperties(input *bigdatapools.AutoPauseProperties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	if !pointer.From(input.Enabled) {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"delay_in_minutes": pointer.From(input.DelayInMinutes),
		},
	}
}

func flattenSynapseSparkPoolAutoScaleProperties(input *bigdatapools.AutoScaleProperties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	if !pointer.From(input.Enabled) {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"max_node_count": pointer.From(input.MaxNodeCount),
			"min_node_count": pointer.From(input.MinNodeCount),
		},
	}
}

func flattenSynapseSparkPoolLibraryRequirements(input *bigdatapools.LibraryRequirements) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"content":  pointer.From(input.Content),
			"filename": pointer.From(input.Filename),
		},
	}
}

func flattenSynapseSparkPoolSparkConfig(input *bigdatapools.SparkConfigProperties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"content":  pointer.From(input.Content),
			"filename": pointer.From(input.Filename),
		},
	}
}
