package synapse

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/sdk/2021-06-01-preview/artifacts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSynapseSparkJobDefinition() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseSparkJobDefinitionCreateUpdate,
		Read:   resourceSynapseSparkJobDefinitionRead,
		Update: resourceSynapseSparkJobDefinitionCreateUpdate,
		Delete: resourceSynapseSparkJobDefinitionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SparkJobDefinitionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"synapse_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"language": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"spark_pool_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"spark_version": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"job": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"archives": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"arguments": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"class_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"driver_cores": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
						},

						"driver_memory": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"executor_cores": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
						},

						"executor_memory": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"file": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"files": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"jars": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"num_executors": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceSynapseSparkJobDefinitionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	workspaceId, err := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	client, err := synapseClient.SparkJobDefinitionClient(workspaceId.Name, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	id := parse.NewSparkJobDefinitionID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.GetSparkJobDefinition(ctx, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_synapse_spark_job_definition", id.ID())
		}
	}

	jobProperty := expandSynapseSparkJobDefinitionJobProperties(d.Get("job").([]interface{}))
	if jobProperty != nil {
		jobProperty.Name = utils.String(id.Name)
	}
	sparkJobDefinition := &artifacts.SparkJobDefinitionResource{
		Name: utils.String(id.Name),
		Properties: &artifacts.SparkJobDefinition{
			Description: utils.String(d.Get("description").(string)),
			TargetBigDataPool: &artifacts.BigDataPoolReference{
				ReferenceName: utils.String(d.Get("spark_pool_name").(string)),
			},
			RequiredSparkVersion: utils.String(d.Get("spark_version").(string)),
			Language:             utils.String(d.Get("language").(string)),
			JobProperties:        jobProperty,
		},
	}
	future, err := client.CreateOrUpdateSparkJobDefinition(ctx, id.Name, *sparkJobDefinition, "")
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creation for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSynapseSparkJobDefinitionRead(d, meta)
}

func resourceSynapseSparkJobDefinitionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	id, err := parse.SparkJobDefinitionID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.SparkJobDefinitionClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	resp, err := client.GetSparkJobDefinition(ctx, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("synapse_workspace_id", parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID())
	if props := resp.Properties; props != nil {
		d.Set("description", props.Description)
		d.Set("language", props.Language)
		d.Set("spark_version", props.RequiredSparkVersion)
		if props.TargetBigDataPool != nil {
			d.Set("spark_pool_name", props.TargetBigDataPool.ReferenceName)
		}
		d.Set("job", flattenSynapseSparkJobDefinitionJobProperties(props.JobProperties))
	}

	return nil
}

func resourceSynapseSparkJobDefinitionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	id, err := parse.SparkJobDefinitionID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.SparkJobDefinitionClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	future, err := client.DeleteSparkJobDefinition(ctx, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}

func expandSynapseSparkJobDefinitionJobProperties(input []interface{}) *artifacts.SparkJobProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	value := input[0].(map[string]interface{})
	return &artifacts.SparkJobProperties{
		File:           utils.String(value["file"].(string)),
		ClassName:      utils.String(value["class_name"].(string)),
		Conf:           nil,
		Args:           utils.ExpandStringSlice(value["arguments"].([]interface{})),
		Jars:           utils.ExpandStringSlice(value["jars"].([]interface{})),
		Files:          utils.ExpandStringSlice(value["files"].([]interface{})),
		Archives:       utils.ExpandStringSlice(value["archives"].([]interface{})),
		DriverMemory:   utils.String(value["driver_memory"].(string)),
		DriverCores:    utils.Int32(int32(value["driver_cores"].(int))),
		ExecutorMemory: utils.String(value["executor_memory"].(string)),
		ExecutorCores:  utils.Int32(int32(value["executor_cores"].(int))),
		NumExecutors:   utils.Int32(int32(value["num_executors"].(int))),
	}
}

func flattenSynapseSparkJobDefinitionJobProperties(input *artifacts.SparkJobProperties) interface{} {
	if input == nil {
		return []interface{}{}
	}

	var file string
	if input.File != nil {
		file = *input.File
	}
	var className string
	if input.ClassName != nil {
		className = *input.ClassName
	}
	var driverMemory string
	if input.DriverMemory != nil {
		driverMemory = *input.DriverMemory
	}
	var executorMemory string
	if input.ExecutorMemory != nil {
		executorMemory = *input.ExecutorMemory
	}
	var driverCores int32
	if input.DriverCores != nil {
		driverCores = *input.DriverCores
	}
	var executorCores int32
	if input.ExecutorCores != nil {
		executorCores = *input.ExecutorCores
	}
	var numExecutors int32
	if input.NumExecutors != nil {
		numExecutors = *input.NumExecutors
	}
	return []interface{}{
		map[string]interface{}{
			"file":            file,
			"class_name":      className,
			"arguments":       utils.FlattenStringSlice(input.Args),
			"jars":            utils.FlattenStringSlice(input.Jars),
			"files":           utils.FlattenStringSlice(input.Files),
			"archives":        utils.FlattenStringSlice(input.Archives),
			"driver_memory":   driverMemory,
			"executor_memory": executorMemory,
			"driver_cores":    driverCores,
			"executor_cores":  executorCores,
			"num_executors":   numExecutors,
		},
	}
}
