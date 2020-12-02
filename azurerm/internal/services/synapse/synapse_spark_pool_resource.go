package synapse

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/2019-06-01-preview/synapse"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/synapse/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/synapse/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSynapseSparkPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSynapseSparkPoolCreateUpdate,
		Read:   resourceArmSynapseSparkPoolRead,
		Update: resourceArmSynapseSparkPoolCreateUpdate,
		Delete: resourceArmSynapseSparkPoolDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SparkPoolID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SparkPoolName,
			},

			"synapse_workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SynapseWorkspaceID,
			},

			"node_size_family": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(synapse.NodeSizeFamilyMemoryOptimized),
				}, false),
			},

			"node_size": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(synapse.NodeSizeSmall),
					string(synapse.NodeSizeMedium),
					string(synapse.NodeSizeLarge),
				}, false),
			},

			"node_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(3, 200),
				ExactlyOneOf: []string{"node_count", "auto_scale"},
			},

			"auto_scale": {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: []string{"node_count", "auto_scale"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_node_count": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(3, 200),
						},

						"max_node_count": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(3, 200),
						},
					},
				},
			},

			"auto_pause": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delay_in_minutes": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(5, 10080),
						},
					},
				},
			},

			"spark_events_folder": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/events",
			},

			"spark_log_folder": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/logs",
			},

			"library_requirement": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},

						"filename": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"spark_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "2.4",
				ValidateFunc: validation.StringInSlice([]string{
					"2.4",
				}, false),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmSynapseSparkPoolCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.SparkPoolClient
	workspaceClient := meta.(*clients.Client).Synapse.WorkspaceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	workspaceId, _ := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, workspaceId.ResourceGroup, workspaceId.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Synapse Spark Pool %q (Workspace %q / Resource Group %q): %+v", name, workspaceId.Name, workspaceId.ResourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_synapse_spark_pool", *existing.ID)
		}
	}

	workspace, err := workspaceClient.Get(ctx, workspaceId.ResourceGroup, workspaceId.Name)
	if err != nil {
		return fmt.Errorf("reading Synapse workspace %q (Workspace %q / Resource Group %q): %+v", workspaceId.Name, workspaceId.Name, workspaceId.ResourceGroup, err)
	}

	autoScale := expandArmSparkPoolAutoScaleProperties(d.Get("auto_scale").([]interface{}))
	bigDataPoolInfo := synapse.BigDataPoolResourceInfo{
		Location: workspace.Location,
		BigDataPoolResourceProperties: &synapse.BigDataPoolResourceProperties{
			AutoPause:             expandArmSparkPoolAutoPauseProperties(d.Get("auto_pause").([]interface{})),
			AutoScale:             autoScale,
			DefaultSparkLogFolder: utils.String(d.Get("spark_log_folder").(string)),
			LibraryRequirements:   expandArmSparkPoolLibraryRequirements(d.Get("library_requirement").([]interface{})),
			NodeSize:              synapse.NodeSize(d.Get("node_size").(string)),
			NodeSizeFamily:        synapse.NodeSizeFamily(d.Get("node_size_family").(string)),
			SparkEventsFolder:     utils.String(d.Get("spark_events_folder").(string)),
			SparkVersion:          utils.String(d.Get("spark_version").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if !*autoScale.Enabled {
		bigDataPoolInfo.NodeCount = utils.Int32(int32(d.Get("node_count").(int)))
	}

	force := utils.Bool(false)
	future, err := client.CreateOrUpdate(ctx, workspaceId.ResourceGroup, workspaceId.Name, name, bigDataPoolInfo, force)
	if err != nil {
		return fmt.Errorf("creating Synapse Spark Pool %q (Workspace %q / Resource Group %q): %+v", name, workspaceId.Name, workspaceId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating future for Synapse Spark Pool %q (Workspace %q / Resource Group %q): %+v", name, workspaceId.Name, workspaceId.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, workspaceId.ResourceGroup, workspaceId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Synapse Spark Pool %q (Workspace %q / Resource Group %q): %+v", name, workspaceId.Name, workspaceId.ResourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Synapse Spark Pool %q (Workspace %q / Resource Group %q) ID", name, workspaceId.Name, workspaceId.ResourceGroup)
	}

	d.SetId(*resp.ID)
	return resourceArmSynapseSparkPoolRead(d, meta)
}

func resourceArmSynapseSparkPoolRead(d *schema.ResourceData, meta interface{}) error {
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
	workspaceId := parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID("")
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
		d.Set("node_count", props.NodeCount)
		d.Set("node_size", props.NodeSize)
		d.Set("node_size_family", string(props.NodeSizeFamily))
		d.Set("spark_version", props.SparkVersion)
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmSynapseSparkPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.SparkPoolClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SparkPoolID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.BigDataPoolName)
	if err != nil {
		return fmt.Errorf("deleting Synapse Spark Pool %q (Workspace %q / Resource Group %q): %+v", id.BigDataPoolName, id.WorkspaceName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deleting Synapse Spark Pool %q (Workspace %q / Resource Group %q): %+v", id.BigDataPoolName, id.WorkspaceName, id.ResourceGroup, err)
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
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &synapse.LibraryRequirements{
		Content:  utils.String(v["content"].(string)),
		Filename: utils.String(v["filename"].(string)),
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
