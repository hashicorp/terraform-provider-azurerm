package synapse

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/synapse/mgmt/2021-03-01/synapse"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
					string(synapse.NodeSizeFamilyMemoryOptimized),
					string(synapse.NodeSizeFamilyNone),
				}, false),
			},

			"node_size": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(synapse.NodeSizeSmall),
					string(synapse.NodeSizeMedium),
					string(synapse.NodeSizeLarge),
					string(synapse.NodeSizeNone),
					string(synapse.NodeSizeXLarge),
					string(synapse.NodeSizeXXLarge),
					string(synapse.NodeSizeXXXLarge),
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

			"node_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
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
				Optional: true,
				Default:  "2.4",
				ValidateFunc: validation.StringInSlice([]string{
					"2.4",
					"3.1",
					"3.2",
				}, false),
			},

			"tags": tags.Schema(),
		},
	}
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

	id := parse.NewSparkPoolID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.BigDataPoolName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_synapse_spark_pool", id.ID())
		}
	}

	workspace, err := workspaceClient.Get(ctx, id.ResourceGroup, id.WorkspaceName)
	if err != nil {
		return fmt.Errorf("reading Synapse workspace %q (Workspace %q / Resource Group %q): %+v", workspaceId.Name, workspaceId.Name, workspaceId.ResourceGroup, err)
	}

	autoScale := expandArmSparkPoolAutoScaleProperties(d.Get("auto_scale").([]interface{}))
	bigDataPoolInfo := synapse.BigDataPoolResourceInfo{
		Location: workspace.Location,
		BigDataPoolResourceProperties: &synapse.BigDataPoolResourceProperties{
			AutoPause:                 expandArmSparkPoolAutoPauseProperties(d.Get("auto_pause").([]interface{})),
			AutoScale:                 autoScale,
			CacheSize:                 utils.Int32(int32(d.Get("cache_size").(int))),
			IsComputeIsolationEnabled: utils.Bool(d.Get("compute_isolation_enabled").(bool)),
			DynamicExecutorAllocation: &synapse.DynamicExecutorAllocation{
				Enabled: utils.Bool(d.Get("dynamic_executor_allocation_enabled").(bool)),
			},
			DefaultSparkLogFolder:       utils.String(d.Get("spark_log_folder").(string)),
			NodeSize:                    synapse.NodeSize(d.Get("node_size").(string)),
			NodeSizeFamily:              synapse.NodeSizeFamily(d.Get("node_size_family").(string)),
			SessionLevelPackagesEnabled: utils.Bool(d.Get("session_level_packages_enabled").(bool)),
			SparkConfigProperties:       expandSparkPoolSparkConfig(d.Get("spark_config").([]interface{})),
			SparkEventsFolder:           utils.String(d.Get("spark_events_folder").(string)),
			SparkVersion:                utils.String(d.Get("spark_version").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if !*autoScale.Enabled {
		bigDataPoolInfo.NodeCount = utils.Int32(int32(d.Get("node_count").(int)))
	}

	force := utils.Bool(false)
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.BigDataPoolName, bigDataPoolInfo, force)
	if err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// Library Requirements can't be specified on Create so we'll call update after we've confirmed the Spark Pool has been created.
	return resourceSynapseSparkPoolUpdate(d, meta)
}

func resourceSynapseSparkPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.SparkPoolClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SparkPoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.BigDataPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Synapse Spark Pool %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Synapse Spark Pool %q (Workspace %q / Resource Group %q): %+v", id.BigDataPoolName, id.WorkspaceName, id.ResourceGroup, err)
	}
	d.Set("name", id.BigDataPoolName)
	workspaceId := parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID()
	d.Set("synapse_workspace_id", workspaceId)

	if props := resp.BigDataPoolResourceProperties; props != nil {
		if err := d.Set("auto_pause", flattenArmSparkPoolAutoPauseProperties(props.AutoPause)); err != nil {
			return fmt.Errorf("setting `auto_pause`: %+v", err)
		}
		if err := d.Set("auto_scale", flattenArmSparkPoolAutoScaleProperties(props.AutoScale)); err != nil {
			return fmt.Errorf("setting `auto_scale`: %+v", err)
		}
		if err := d.Set("library_requirement", flattenArmSparkPoolLibraryRequirements(props.LibraryRequirements)); err != nil {
			return fmt.Errorf("setting `library_requirement`: %+v", err)
		}
		d.Set("cache_size", props.CacheSize)
		d.Set("compute_isolation_enabled", props.IsComputeIsolationEnabled)

		dynamicExecutorAllocationEnabled := false
		if props.DynamicExecutorAllocation != nil {
			dynamicExecutorAllocationEnabled = *props.DynamicExecutorAllocation.Enabled
		}
		d.Set("dynamic_executor_allocation_enabled", dynamicExecutorAllocationEnabled)

		d.Set("node_count", props.NodeCount)
		d.Set("node_size", props.NodeSize)
		d.Set("node_size_family", string(props.NodeSizeFamily))
		d.Set("session_level_packages_enabled", props.SessionLevelPackagesEnabled)
		d.Set("spark_config", flattenSparkPoolSparkConfig(props.SparkConfigProperties))
		d.Set("spark_version", props.SparkVersion)
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSynapseSparkPoolUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.SparkPoolClient
	workspaceClient := meta.(*clients.Client).Synapse.WorkspaceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SparkPoolID(d.Id())
	if err != nil {
		return err
	}

	workspace, err := workspaceClient.Get(ctx, id.ResourceGroup, id.WorkspaceName)
	if err != nil {
		return fmt.Errorf("reading Synapse workspace %q (Workspace %q / Resource Group %q): %+v", id.WorkspaceName, id.WorkspaceName, id.ResourceGroup, err)
	}

	autoScale := expandArmSparkPoolAutoScaleProperties(d.Get("auto_scale").([]interface{}))
	bigDataPoolInfo := synapse.BigDataPoolResourceInfo{
		Location: workspace.Location,
		BigDataPoolResourceProperties: &synapse.BigDataPoolResourceProperties{
			AutoPause:                 expandArmSparkPoolAutoPauseProperties(d.Get("auto_pause").([]interface{})),
			AutoScale:                 autoScale,
			CacheSize:                 utils.Int32(int32(d.Get("cache_size").(int))),
			IsComputeIsolationEnabled: utils.Bool(d.Get("compute_isolation_enabled").(bool)),
			DynamicExecutorAllocation: &synapse.DynamicExecutorAllocation{
				Enabled: utils.Bool(d.Get("dynamic_executor_allocation_enabled").(bool)),
			},
			DefaultSparkLogFolder:       utils.String(d.Get("spark_log_folder").(string)),
			LibraryRequirements:         expandArmSparkPoolLibraryRequirements(d.Get("library_requirement").([]interface{})),
			NodeSize:                    synapse.NodeSize(d.Get("node_size").(string)),
			NodeSizeFamily:              synapse.NodeSizeFamily(d.Get("node_size_family").(string)),
			SessionLevelPackagesEnabled: utils.Bool(d.Get("session_level_packages_enabled").(bool)),
			SparkConfigProperties:       expandSparkPoolSparkConfig(d.Get("spark_config").([]interface{})),
			SparkEventsFolder:           utils.String(d.Get("spark_events_folder").(string)),
			SparkVersion:                utils.String(d.Get("spark_version").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if !*autoScale.Enabled {
		bigDataPoolInfo.NodeCount = utils.Int32(int32(d.Get("node_count").(int)))
	}

	force := utils.Bool(false)
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.BigDataPoolName, bigDataPoolInfo, force)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", *id, err)
	}

	return resourceSynapseSparkPoolRead(d, meta)
}

func resourceSynapseSparkPoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.SparkPoolClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SparkPoolID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.BigDataPoolName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandArmSparkPoolAutoPauseProperties(input []interface{}) *synapse.AutoPauseProperties {
	if len(input) == 0 {
		return &synapse.AutoPauseProperties{
			Enabled: utils.Bool(false),
		}
	}
	v := input[0].(map[string]interface{})
	return &synapse.AutoPauseProperties{
		DelayInMinutes: utils.Int32(int32(v["delay_in_minutes"].(int))),
		Enabled:        utils.Bool(true),
	}
}

func expandArmSparkPoolAutoScaleProperties(input []interface{}) *synapse.AutoScaleProperties {
	if len(input) == 0 || input[0] == nil {
		return &synapse.AutoScaleProperties{
			Enabled: utils.Bool(false),
		}
	}
	v := input[0].(map[string]interface{})
	return &synapse.AutoScaleProperties{
		MinNodeCount: utils.Int32(int32(v["min_node_count"].(int))),
		Enabled:      utils.Bool(true),
		MaxNodeCount: utils.Int32(int32(v["max_node_count"].(int))),
	}
}

func expandArmSparkPoolLibraryRequirements(input []interface{}) *synapse.LibraryRequirements {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &synapse.LibraryRequirements{
		Content:  utils.String(v["content"].(string)),
		Filename: utils.String(v["filename"].(string)),
	}
}

func expandSparkPoolSparkConfig(input []interface{}) *synapse.LibraryRequirements {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	value := input[0].(map[string]interface{})
	return &synapse.LibraryRequirements{
		Content:  utils.String(value["content"].(string)),
		Filename: utils.String(value["filename"].(string)),
	}
}

func flattenArmSparkPoolAutoPauseProperties(input *synapse.AutoPauseProperties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var delayInMinutes int32
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

func flattenArmSparkPoolAutoScaleProperties(input *synapse.AutoScaleProperties) []interface{} {
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

	var maxNodeCount int32
	if input.MaxNodeCount != nil {
		maxNodeCount = *input.MaxNodeCount
	}
	var minNodeCount int32
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

func flattenArmSparkPoolLibraryRequirements(input *synapse.LibraryRequirements) []interface{} {
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

func flattenSparkPoolSparkConfig(input *synapse.LibraryRequirements) []interface{} {
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
