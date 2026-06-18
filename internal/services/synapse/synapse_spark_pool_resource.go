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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/bigdatapools"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseSparkPool() *pluginsdk.Resource {
	r := &pluginsdk.Resource{
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
			_, err := parse.SparkPoolID(id)
			return err
		}),

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
				ValidateFunc: validate.WorkspaceID,
			},

			"node_size_family": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(bigdatapools.NodeSizeFamilyHardwareAcceleratedFPGA),
					string(bigdatapools.NodeSizeFamilyHardwareAcceleratedGPU),
					string(bigdatapools.NodeSizeFamilyMemoryOptimized),
					string(bigdatapools.NodeSizeFamilyNone),
				}, false),
			},

			"node_size": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(bigdatapools.NodeSizeSmall),
					string(bigdatapools.NodeSizeMedium),
					string(bigdatapools.NodeSizeLarge),
					string(bigdatapools.NodeSizeNone),
					string(bigdatapools.NodeSizeXLarge),
					string(bigdatapools.NodeSizeXXLarge),
					string(bigdatapools.NodeSizeXXXLarge),
				}, false),
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

	if !features.FivePointOh() {
		r.Schema["spark_version"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.StringInSlice([]string{
					"3.2",
					"3.3",
					"3.4",
					"3.5",
				}, false),
				func(v interface{}, k string) (warnings []string, errors []error) {
					if val, ok := v.(string); ok && (val == "3.2" || val == "3.3") {
						warnings = append(warnings, fmt.Sprintf("Spark version %s is deprecated and will be removed in a future version of the AzureRM provider. Please consider upgrading to version 3.4 or later.", val))
					}
					return
				},
			),
		}
	}

	return r
}

func resourceSynapseSparkPoolCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.SparkPoolClient
	workspaceClient := meta.(*clients.Client).Synapse.WorkspaceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `synapse_workspace_id`: %+v", err)
	}

	id := bigdatapools.NewBigDataPoolID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, d.Get("name").(string))
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

	workspace, err := workspaceClient.Get(ctx, id.ResourceGroupName, id.WorkspaceName)
	if err != nil {
		return fmt.Errorf("reading Synapse workspace %q (Workspace %q / Resource Group %q): %+v", id.WorkspaceName, id.WorkspaceName, id.ResourceGroupName, err)
	}

	autoScale := expandArmSparkPoolAutoScaleProperties(d.Get("auto_scale").([]interface{}))
	bigDataPoolInfo := bigdatapools.BigDataPoolResourceInfo{
		Location: location.Normalize(pointer.From(workspace.Location)),
		Properties: &bigdatapools.BigDataPoolResourceProperties{
			AutoPause:                 expandArmSparkPoolAutoPauseProperties(d.Get("auto_pause").([]interface{})),
			AutoScale:                 autoScale,
			CacheSize:                 pointer.To(int64(d.Get("cache_size").(int))),
			IsComputeIsolationEnabled: pointer.To(d.Get("compute_isolation_enabled").(bool)),
			DynamicExecutorAllocation: &bigdatapools.DynamicExecutorAllocation{
				Enabled:      pointer.To(d.Get("dynamic_executor_allocation_enabled").(bool)),
				MinExecutors: pointer.To(int64(d.Get("min_executors").(int))),
				MaxExecutors: pointer.To(int64(d.Get("max_executors").(int))),
			},
			DefaultSparkLogFolder:       pointer.To(d.Get("spark_log_folder").(string)),
			NodeSize:                    pointer.To(bigdatapools.NodeSize(d.Get("node_size").(string))),
			NodeSizeFamily:              pointer.To(bigdatapools.NodeSizeFamily(d.Get("node_size_family").(string))),
			SessionLevelPackagesEnabled: pointer.To(d.Get("session_level_packages_enabled").(bool)),
			SparkConfigProperties:       expandSparkPoolSparkConfig(d.Get("spark_config").([]interface{})),
			SparkEventsFolder:           pointer.To(d.Get("spark_events_folder").(string)),
			SparkVersion:                pointer.To(d.Get("spark_version").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if !pointer.From(autoScale.Enabled) {
		bigDataPoolInfo.Properties.NodeCount = pointer.To(int64(d.Get("node_count").(int)))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, bigDataPoolInfo, bigdatapools.CreateOrUpdateOperationOptions{Force: pointer.To(false)}); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	d.SetId(id.ID())

	// Library Requirements can't be specified on Create so we'll call update after we've confirmed the Spark Pool has been created.
	return resourceSynapseSparkPoolUpdate(d, meta)
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
			log.Printf("[INFO] Synapse Spark Pool %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	d.Set("name", id.BigDataPoolName)
	workspaceId := parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID()
	d.Set("synapse_workspace_id", workspaceId)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if err := d.Set("auto_pause", flattenArmSparkPoolAutoPauseProperties(props.AutoPause)); err != nil {
				return fmt.Errorf("setting `auto_pause`: %+v", err)
			}
			if err := d.Set("auto_scale", flattenArmSparkPoolAutoScaleProperties(props.AutoScale)); err != nil {
				return fmt.Errorf("setting `auto_scale`: %+v", err)
			}
			if err := d.Set("library_requirement", flattenArmSparkPoolLibraryRequirements(props.LibraryRequirements)); err != nil {
				return fmt.Errorf("setting `library_requirement`: %+v", err)
			}
			d.Set("cache_size", pointer.From(props.CacheSize))
			d.Set("compute_isolation_enabled", pointer.From(props.IsComputeIsolationEnabled))

			dynamicExecutorAllocationEnabled := false
			minExector := 0
			maxExecutor := 0
			if props.DynamicExecutorAllocation != nil {
				dynamicExecutorAllocationEnabled = pointer.From(props.DynamicExecutorAllocation.Enabled)
				minExector = int(pointer.From(props.DynamicExecutorAllocation.MinExecutors))
				maxExecutor = int(pointer.From(props.DynamicExecutorAllocation.MaxExecutors))
			}
			d.Set("dynamic_executor_allocation_enabled", dynamicExecutorAllocationEnabled)
			d.Set("min_executors", minExector)
			d.Set("max_executors", maxExecutor)

			d.Set("node_count", pointer.From(props.NodeCount))
			d.Set("node_size", string(pointer.From(props.NodeSize)))
			d.Set("node_size_family", string(pointer.From(props.NodeSizeFamily)))
			d.Set("session_level_packages_enabled", pointer.From(props.SessionLevelPackagesEnabled))
			d.Set("spark_config", flattenSparkPoolSparkConfig(props.SparkConfigProperties))
			d.Set("spark_version", pointer.From(props.SparkVersion))
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}

func resourceSynapseSparkPoolUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.SparkPoolClient
	workspaceClient := meta.(*clients.Client).Synapse.WorkspaceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := bigdatapools.ParseBigDataPoolID(d.Id())
	if err != nil {
		return err
	}

	workspace, err := workspaceClient.Get(ctx, id.ResourceGroupName, id.WorkspaceName)
	if err != nil {
		return fmt.Errorf("reading Synapse workspace %q (Workspace %q / Resource Group %q): %+v", id.WorkspaceName, id.WorkspaceName, id.ResourceGroupName, err)
	}

	autoScale := expandArmSparkPoolAutoScaleProperties(d.Get("auto_scale").([]interface{}))
	bigDataPoolInfo := bigdatapools.BigDataPoolResourceInfo{
		Location: location.Normalize(pointer.From(workspace.Location)),
		Properties: &bigdatapools.BigDataPoolResourceProperties{
			AutoPause:                 expandArmSparkPoolAutoPauseProperties(d.Get("auto_pause").([]interface{})),
			AutoScale:                 autoScale,
			CacheSize:                 pointer.To(int64(d.Get("cache_size").(int))),
			IsComputeIsolationEnabled: pointer.To(d.Get("compute_isolation_enabled").(bool)),
			DynamicExecutorAllocation: &bigdatapools.DynamicExecutorAllocation{
				Enabled:      pointer.To(d.Get("dynamic_executor_allocation_enabled").(bool)),
				MinExecutors: pointer.To(int64(d.Get("min_executors").(int))),
				MaxExecutors: pointer.To(int64(d.Get("max_executors").(int))),
			},
			DefaultSparkLogFolder:       pointer.To(d.Get("spark_log_folder").(string)),
			LibraryRequirements:         expandArmSparkPoolLibraryRequirements(d.Get("library_requirement").([]interface{})),
			NodeSize:                    pointer.To(bigdatapools.NodeSize(d.Get("node_size").(string))),
			NodeSizeFamily:              pointer.To(bigdatapools.NodeSizeFamily(d.Get("node_size_family").(string))),
			SessionLevelPackagesEnabled: pointer.To(d.Get("session_level_packages_enabled").(bool)),
			SparkConfigProperties:       expandSparkPoolSparkConfig(d.Get("spark_config").([]interface{})),
			SparkEventsFolder:           pointer.To(d.Get("spark_events_folder").(string)),
			SparkVersion:                pointer.To(d.Get("spark_version").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if !pointer.From(autoScale.Enabled) {
		bigDataPoolInfo.Properties.NodeCount = pointer.To(int64(d.Get("node_count").(int)))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, bigDataPoolInfo, bigdatapools.CreateOrUpdateOperationOptions{Force: pointer.To(false)}); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
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
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandArmSparkPoolAutoPauseProperties(input []interface{}) *bigdatapools.AutoPauseProperties {
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

func expandArmSparkPoolAutoScaleProperties(input []interface{}) *bigdatapools.AutoScaleProperties {
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

func expandArmSparkPoolLibraryRequirements(input []interface{}) *bigdatapools.LibraryRequirements {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &bigdatapools.LibraryRequirements{
		Content:  pointer.To(v["content"].(string)),
		Filename: pointer.To(v["filename"].(string)),
	}
}

func expandSparkPoolSparkConfig(input []interface{}) *bigdatapools.SparkConfigProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	value := input[0].(map[string]interface{})
	return &bigdatapools.SparkConfigProperties{
		Content:  pointer.To(value["content"].(string)),
		Filename: pointer.To(value["filename"].(string)),
	}
}

func flattenArmSparkPoolAutoPauseProperties(input *bigdatapools.AutoPauseProperties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var delayInMinutes int64
	if input.DelayInMinutes != nil {
		delayInMinutes = *input.DelayInMinutes
	}
	var enabled bool
	if input.Enabled != nil {
		enabled = *input.Enabled
	}

	if !enabled {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"delay_in_minutes": delayInMinutes,
		},
	}
}

func flattenArmSparkPoolAutoScaleProperties(input *bigdatapools.AutoScaleProperties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var enabled bool
	if input.Enabled != nil {
		enabled = *input.Enabled
	}

	if !enabled {
		return make([]interface{}, 0)
	}

	var maxNodeCount int64
	if input.MaxNodeCount != nil {
		maxNodeCount = *input.MaxNodeCount
	}
	var minNodeCount int64
	if input.MinNodeCount != nil {
		minNodeCount = *input.MinNodeCount
	}
	return []interface{}{
		map[string]interface{}{
			"max_node_count": maxNodeCount,
			"min_node_count": minNodeCount,
		},
	}
}

func flattenArmSparkPoolLibraryRequirements(input *bigdatapools.LibraryRequirements) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var content string
	if input.Content != nil {
		content = *input.Content
	}
	var filename string
	if input.Filename != nil {
		filename = *input.Filename
	}
	return []interface{}{
		map[string]interface{}{
			"content":  content,
			"filename": filename,
		},
	}
}

func flattenSparkPoolSparkConfig(input *bigdatapools.SparkConfigProperties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var content string
	if input.Content != nil {
		content = *input.Content
	}
	var filename string
	if input.Filename != nil {
		filename = *input.Filename
	}
	return []interface{}{
		map[string]interface{}{
			"content":  content,
			"filename": filename,
		},
	}
}
