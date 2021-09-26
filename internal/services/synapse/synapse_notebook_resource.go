package synapse

import (
	"encoding/json"
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

func resourceSynapseNotebook() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseNotebookCreateUpdate,
		Read:   resourceSynapseNotebookRead,
		Update: resourceSynapseNotebookCreateUpdate,
		Delete: resourceSynapseNotebookDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.NotebookID(id)
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

			"cells": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				StateFunc:        normalizeCellsJSON,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			},

			"codemirror_mode": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"display_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"language": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"major_version": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
			},

			"minor_version": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
			},

			"spark_pool_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"session_config": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
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

func resourceSynapseNotebookCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	workspaceId, err := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	client, err := synapseClient.NotebookClient(workspaceId.Name, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	id := parse.NewNotebookID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.GetNotebook(ctx, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_synapse_notebook", id.ID())
		}
	}

	notebook := &artifacts.NotebookResource{
		Name: utils.String(id.Name),
		Properties: &artifacts.Notebook{
			Description: utils.String(d.Get("description").(string)),
			BigDataPool: &artifacts.BigDataPoolReference{
				ReferenceName: utils.String(d.Get("spark_pool_name").(string)),
			},
			SessionProperties: expandSynapseNotebookSessionProperties(d.Get("session_config").([]interface{})),
			Metadata: &artifacts.NotebookMetadata{
				Kernelspec: &artifacts.NotebookKernelSpec{
					Name:        utils.String(id.Name),
					DisplayName: utils.String(d.Get("display_name").(string)),
				},
				LanguageInfo: &artifacts.NotebookLanguageInfo{
					Name:           utils.String(d.Get("language").(string)),
					CodemirrorMode: utils.String(d.Get("codemirror_mode").(string)),
				},
			},
			Nbformat:      utils.Int32(int32(d.Get("major_version").(int))),
			NbformatMinor: utils.Int32(int32(d.Get("minor_version").(int))),
		},
	}
	if cells, err := expandSynapseNotebookCells(d.Get("cells").(string)); err == nil {
		notebook.Properties.Cells = cells
	} else {
		return fmt.Errorf("error unmarshal `cells`: %+v", err)
	}
	future, err := client.CreateOrUpdateNotebook(ctx, id.Name, *notebook, "")
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creation for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSynapseNotebookRead(d, meta)
}

func resourceSynapseNotebookRead(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	id, err := parse.NotebookID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.NotebookClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	resp, err := client.GetNotebook(ctx, id.Name, "")
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
		d.Set("major_version", props.Nbformat)
		d.Set("minor_version", props.NbformatMinor)
		if props.BigDataPool != nil {
			d.Set("spark_pool_name", props.BigDataPool.ReferenceName)
		}
		if metadata := props.Metadata; metadata != nil {
			if metadata.Kernelspec != nil {
				d.Set("display_name", metadata.Kernelspec.DisplayName)
			}
			if metadata.LanguageInfo != nil {
				d.Set("language", metadata.LanguageInfo.Name)
				d.Set("codemirror_mode", metadata.LanguageInfo.CodemirrorMode)
			}
		}
		d.Set("session_config", flattenSynapseNotebookSessionProperties(props.SessionProperties))
		if cells, err := flattenSynapseNotebookCells(props.Cells); err == nil {
			d.Set("cells", cells)
		} else {
			return fmt.Errorf("error marshal `cells`: %+v", err)
		}
	}

	return nil
}

func resourceSynapseNotebookDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	id, err := parse.NotebookID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.NotebookClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	future, err := client.DeleteNotebook(ctx, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}

func expandSynapseNotebookCells(input string) (*[]artifacts.NotebookCell, error) {
	if len(input) == 0 {
		return nil, nil
	}
	var result []artifacts.NotebookCell
	err := json.Unmarshal([]byte(input), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func expandSynapseNotebookSessionProperties(input []interface{}) *artifacts.NotebookSessionProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	value := input[0].(map[string]interface{})
	return &artifacts.NotebookSessionProperties{
		DriverMemory:   utils.String(value["driver_memory"].(string)),
		DriverCores:    utils.Int32(int32(value["driver_cores"].(int))),
		ExecutorMemory: utils.String(value["executor_memory"].(string)),
		ExecutorCores:  utils.Int32(int32(value["executor_cores"].(int))),
		NumExecutors:   utils.Int32(int32(value["num_executors"].(int))),
	}
}

func flattenSynapseNotebookCells(cells *[]artifacts.NotebookCell) (interface{}, error) {
	if cells == nil {
		return nil, nil
	}
	data, err := json.Marshal(cells)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func flattenSynapseNotebookSessionProperties(input *artifacts.NotebookSessionProperties) interface{} {
	if input == nil {
		return []interface{}{}
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
			"driver_memory":   driverMemory,
			"executor_memory": executorMemory,
			"driver_cores":    driverCores,
			"executor_cores":  executorCores,
			"num_executors":   numExecutors,
		},
	}
}

func normalizeCellsJSON(jsonString interface{}) string {
	if jsonString == nil || jsonString == "" {
		return ""
	}
	var result []artifacts.NotebookCell
	err := json.Unmarshal([]byte(jsonString.(string)), &result)
	if err != nil {
		return fmt.Sprintf("unable to parse JSON: %+v", err)
	}
	b, _ := json.Marshal(result)
	return string(b)
}
