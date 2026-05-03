// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/rbacs"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCosmosDbSQLRoleDefinition() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosDbSQLRoleDefinitionCreate,
		Read:   resourceCosmosDbSQLRoleDefinitionRead,
		Update: resourceCosmosDbSQLRoleDefinitionUpdate,
		Delete: resourceCosmosDbSQLRoleDefinitionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := rbacs.ParseSqlRoleDefinitionID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"role_definition_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosAccountName,
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(rbacs.RoleDefinitionTypeCustomRole),
				ValidateFunc: validation.StringInSlice([]string{
					string(rbacs.RoleDefinitionTypeBuiltInRole),
					string(rbacs.RoleDefinitionTypeCustomRole),
				}, false),
			},

			"assignable_scopes": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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
					},
				},
			},
		},
	}
}

func resourceCosmosDbSQLRoleDefinitionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.RbacsClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	roleDefinitionId := d.Get("role_definition_id").(string)
	if roleDefinitionId == "" {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("generating UUID for Cosmos DB SQL Role Definition: %+v", err)
		}

		roleDefinitionId = uuid
	}

	id := rbacs.NewSqlRoleDefinitionID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), roleDefinitionId)

	locks.ByName(id.DatabaseAccountName, CosmosDbAccountResourceName)
	defer locks.UnlockByName(id.DatabaseAccountName, CosmosDbAccountResourceName)

	existing, err := client.SqlResourcesGetSqlRoleDefinition(ctx, id)
	if !response.WasNotFound(existing.HttpResponse) {
		if err != nil {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
		return tf.ImportAsExistsError("azurerm_cosmosdb_sql_role_definition", id.ID())
	}

	parameters := rbacs.SqlRoleDefinitionCreateUpdateParameters{
		Properties: &rbacs.SqlRoleDefinitionResource{
			RoleName:         pointer.To(d.Get("name").(string)),
			AssignableScopes: utils.ExpandStringSlice(d.Get("assignable_scopes").(*pluginsdk.Set).List()),
			Permissions:      expandSqlRoleDefinitionPermissions(d.Get("permissions").(*pluginsdk.Set).List()),
			Type:             pointer.ToEnum[rbacs.RoleDefinitionType](d.Get("type").(string)),
		},
	}

	if err := client.SqlResourcesCreateUpdateSqlRoleDefinitionThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbSQLRoleDefinitionRead(d, meta)
}

func resourceCosmosDbSQLRoleDefinitionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.RbacsClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := rbacs.ParseSqlRoleDefinitionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.SqlResourcesGetSqlRoleDefinition(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("role_definition_id", id.RoleDefinitionId)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("account_name", id.DatabaseAccountName)

	if resp.Model != nil {
		if props := resp.Model.Properties; props != nil {
			d.Set("assignable_scopes", utils.FlattenStringSlice(props.AssignableScopes))
			d.Set("name", props.RoleName)
			d.Set("type", pointer.FromEnum(props.Type))

			if err := d.Set("permissions", flattenSqlRoleDefinitionPermissions(props.Permissions)); err != nil {
				return fmt.Errorf("setting `permissions`: %+v", err)
			}
		}
	}

	return nil
}

func resourceCosmosDbSQLRoleDefinitionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.RbacsClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := rbacs.ParseSqlRoleDefinitionID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.DatabaseAccountName, CosmosDbAccountResourceName)
	defer locks.UnlockByName(id.DatabaseAccountName, CosmosDbAccountResourceName)

	existing, err := client.SqlResourcesGetSqlRoleDefinition(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %w", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: properties was nil", id)
	}

	parameters := rbacs.SqlRoleDefinitionCreateUpdateParameters{
		Properties: existing.Model.Properties,
	}

	if d.HasChange("assignable_scopes") {
		parameters.Properties.AssignableScopes = utils.ExpandStringSlice(d.Get("assignable_scopes").(*pluginsdk.Set).List())
	}

	if d.HasChange("name") {
		parameters.Properties.RoleName = pointer.To(d.Get("name").(string))
	}

	if d.HasChange("permissions") {
		parameters.Properties.Permissions = expandSqlRoleDefinitionPermissions(d.Get("permissions").(*pluginsdk.Set).List())
	}

	if err := client.SqlResourcesCreateUpdateSqlRoleDefinitionThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceCosmosDbSQLRoleDefinitionRead(d, meta)
}

func resourceCosmosDbSQLRoleDefinitionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.RbacsClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := rbacs.ParseSqlRoleDefinitionID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.DatabaseAccountName, CosmosDbAccountResourceName)
	defer locks.UnlockByName(id.DatabaseAccountName, CosmosDbAccountResourceName)

	if err := client.SqlResourcesDeleteSqlRoleDefinitionThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandSqlRoleDefinitionPermissions(input []interface{}) *[]rbacs.Permission {
	results := make([]rbacs.Permission, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, rbacs.Permission{
			DataActions: utils.ExpandStringSlice(v["data_actions"].(*pluginsdk.Set).List()),
		})
	}

	return &results
}

func flattenSqlRoleDefinitionPermissions(input *[]rbacs.Permission) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"data_actions": utils.FlattenStringSlice(item.DataActions),
		})
	}

	return results
}
