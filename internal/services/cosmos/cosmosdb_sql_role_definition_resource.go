package cosmos

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-06-15/documentdb"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCosmosDbSQLRoleDefinition() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosDbSQLRoleDefinitionCreateUpdate,
		Read:   resourceCosmosDbSQLRoleDefinitionRead,
		Update: resourceCosmosDbSQLRoleDefinitionCreateUpdate,
		Delete: resourceCosmosDbSQLRoleDefinitionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SqlRoleDefinitionID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosAccountName,
			},

			"assignable_scopes": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"permissions": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"data_actions": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"not_data_actions": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},

			"role_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(documentdb.RoleDefinitionTypeCustomRole),
				ValidateFunc: validation.StringInSlice([]string{
					string(documentdb.RoleDefinitionTypeBuiltInRole),
					string(documentdb.RoleDefinitionTypeCustomRole),
				}, false),
			},
		},
	}
}

func resourceCosmosDbSQLRoleDefinitionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Cosmos.SqlResourceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)

	id := parse.NewSqlRoleDefinitionID(subscriptionId, resourceGroup, accountName, name)

	if d.IsNewResource() {
		existing, err := client.GetSQLRoleDefinition(ctx, name, resourceGroup, accountName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_cosmosdb_sql_role_definition", id.ID())
		}
	}

	parameters := documentdb.SQLRoleDefinitionCreateUpdateParameters{
		SQLRoleDefinitionResource: &documentdb.SQLRoleDefinitionResource{
			RoleName:         utils.String(d.Get("role_name").(string)),
			AssignableScopes: utils.ExpandStringSlice(d.Get("assignable_scopes").(*pluginsdk.Set).List()),
			Permissions:      expandSqlRoleDefinitionPermissions(d.Get("permissions").(*pluginsdk.Set).List()),
		},
	}

	if v, ok := d.GetOk("type"); ok {
		parameters.SQLRoleDefinitionResource.Type = documentdb.RoleDefinitionType(v.(string))
	}

	future, err := client.CreateUpdateSQLRoleDefinition(ctx, name, resourceGroup, accountName, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbSQLRoleDefinitionRead(d, meta)
}

func resourceCosmosDbSQLRoleDefinitionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlRoleDefinitionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetSQLRoleDefinition(ctx, id.Name, id.ResourceGroup, id.DatabaseAccountName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.DatabaseAccountName)
	d.Set("type", resp.Type)

	if props := resp.SQLRoleDefinitionResource; props != nil {
		d.Set("assignable_scopes", utils.FlattenStringSlice(props.AssignableScopes))
		d.Set("role_name", props.RoleName)

		if err := d.Set("permissions", flattenSqlRoleDefinitionPermissions(props.Permissions)); err != nil {
			return fmt.Errorf("setting `permissions`: %+v", err)
		}
	}

	return nil
}

func resourceCosmosDbSQLRoleDefinitionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlRoleDefinitionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.DeleteSQLRoleDefinition(ctx, id.Name, id.ResourceGroup, id.DatabaseAccountName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandSqlRoleDefinitionPermissions(input []interface{}) *[]documentdb.Permission {
	results := make([]documentdb.Permission, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := documentdb.Permission{
			DataActions: utils.ExpandStringSlice(v["data_actions"].(*pluginsdk.Set).List()),
		}

		if notDataActions, ok := v["not_data_actions"]; ok {
			result.NotDataActions = utils.ExpandStringSlice(notDataActions.(*pluginsdk.Set).List())
		}

		results = append(results, result)
	}

	return &results
}

func flattenSqlRoleDefinitionPermissions(input *[]documentdb.Permission) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"data_actions":     utils.FlattenStringSlice(item.DataActions),
			"not_data_actions": utils.FlattenStringSlice(item.NotDataActions),
		})
	}

	return results
}
