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

func resourceSynapseSqlScript() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseSqlScriptCreateUpdate,
		Read:   resourceSynapseSqlScriptRead,
		Update: resourceSynapseSqlScriptCreateUpdate,
		Delete: resourceSynapseSqlScriptDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SqlScriptID(id)
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

			"query": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"sql_connection": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(artifacts.SQLConnectionTypeSQLOnDemand),
								string(artifacts.SQLConnectionTypeSQLPool),
							}, false),
						},
					},
				},
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(artifacts.SQLQuery),
			},
		},
	}
}

func resourceSynapseSqlScriptCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	workspaceId, err := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	client, err := synapseClient.SQLScriptClient(workspaceId.Name, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	id := parse.NewSqlScriptID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.GetSQLScript(ctx, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_synapse_sql_script", id.ID())
		}
	}

	sqlScriptResource := &artifacts.SQLScriptResource{
		Name: utils.String(id.Name),
		Properties: &artifacts.SQLScript{
			Description: utils.String(d.Get("description").(string)),
			Type:        artifacts.SQLScriptType(d.Get("type").(string)),
			Content: &artifacts.SQLScriptContent{
				Query:             utils.String(d.Get("query").(string)),
				CurrentConnection: expandSynapseSqlScriptConnection(d.Get("sql_connection").([]interface{})),
				Metadata: &artifacts.SQLScriptMetadata{
					Language: utils.String(d.Get("language").(string)),
				},
			},
		},
	}
	future, err := client.CreateOrUpdateSQLScript(ctx, id.Name, *sqlScriptResource, "")
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creation for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSynapseSqlScriptRead(d, meta)
}

func resourceSynapseSqlScriptRead(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	id, err := parse.SqlScriptID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.SQLScriptClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	resp, err := client.GetSQLScript(ctx, id.Name, "")
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
		d.Set("type", props.Type)
		if content := props.Content; content != nil {
			d.Set("query", content.Query)
			d.Set("sql_connection", flattenSynapseSqlScript(content.CurrentConnection))
			if content.Metadata != nil {
				d.Set("language", content.Metadata.Language)
			}
		}
	}

	return nil
}

func resourceSynapseSqlScriptDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	id, err := parse.SqlScriptID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.SQLScriptClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	future, err := client.DeleteSQLScript(ctx, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}

func expandSynapseSqlScriptConnection(input []interface{}) *artifacts.SQLConnection {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	value := input[0].(map[string]interface{})
	return &artifacts.SQLConnection{
		Type: artifacts.SQLConnectionType(value["type"].(string)),
		Name: utils.String(value["name"].(string)),
	}
}

func flattenSynapseSqlScript(input *artifacts.SQLConnection) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var name string
	if input.Name != nil {
		name = *input.Name
	}

	return []interface{}{
		map[string]interface{}{
			"name": name,
			"type": input.Type,
		},
	}
}
