// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCosmosDbSQLRoleAssignment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosDbSQLRoleAssignmentCreate,
		Read:   resourceCosmosDbSQLRoleAssignmentRead,
		Update: resourceCosmosDbSQLRoleAssignmentUpdate,
		Delete: resourceCosmosDbSQLRoleAssignmentDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SqlRoleAssignmentID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
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

			"principal_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"scope": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"role_definition_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.SqlRoleDefinitionID,
			},
		},
	}
}

func resourceCosmosDbSQLRoleAssignmentCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Cosmos.SqlResourceClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	if name == "" {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("generating UUID for Cosmos DB SQL Role Assignment: %+v", err)
		}

		name = uuid
	}

	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)

	id := parse.NewSqlRoleAssignmentID(subscriptionId, resourceGroup, accountName, name)

	locks.ByName(id.DatabaseAccountName, CosmosDbAccountResourceName)
	defer locks.UnlockByName(id.DatabaseAccountName, CosmosDbAccountResourceName)

	existing, err := client.GetSQLRoleAssignment(ctx, id.Name, id.ResourceGroup, id.DatabaseAccountName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cosmosdb_sql_role_assignment", id.ID())
	}

	parameters := documentdb.SQLRoleAssignmentCreateUpdateParameters{
		SQLRoleAssignmentResource: &documentdb.SQLRoleAssignmentResource{
			PrincipalID:      utils.String(d.Get("principal_id").(string)),
			RoleDefinitionID: utils.String(d.Get("role_definition_id").(string)),
			Scope:            utils.String(d.Get("scope").(string)),
		},
	}

	future, err := client.CreateUpdateSQLRoleAssignment(ctx, id.Name, id.ResourceGroup, id.DatabaseAccountName, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the completion of the creating/updating of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCosmosDbSQLRoleAssignmentRead(d, meta)
}

func resourceCosmosDbSQLRoleAssignmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlRoleAssignmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetSQLRoleAssignment(ctx, id.Name, id.ResourceGroup, id.DatabaseAccountName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.DatabaseAccountName)

	if props := resp.SQLRoleAssignmentResource; props != nil {
		d.Set("principal_id", props.PrincipalID)
		d.Set("role_definition_id", props.RoleDefinitionID)
		d.Set("scope", props.Scope)
	}

	return nil
}

func resourceCosmosDbSQLRoleAssignmentUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlResourceClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlRoleAssignmentID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.DatabaseAccountName, CosmosDbAccountResourceName)
	defer locks.UnlockByName(id.DatabaseAccountName, CosmosDbAccountResourceName)

	parameters := documentdb.SQLRoleAssignmentCreateUpdateParameters{
		SQLRoleAssignmentResource: &documentdb.SQLRoleAssignmentResource{
			PrincipalID:      utils.String(d.Get("principal_id").(string)),
			RoleDefinitionID: utils.String(d.Get("role_definition_id").(string)),
			Scope:            utils.String(d.Get("scope").(string)),
		},
	}

	future, err := client.CreateUpdateSQLRoleAssignment(ctx, id.Name, id.ResourceGroup, id.DatabaseAccountName, parameters)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the completion of the updating of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCosmosDbSQLRoleAssignmentRead(d, meta)
}

func resourceCosmosDbSQLRoleAssignmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlRoleAssignmentID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.DatabaseAccountName, CosmosDbAccountResourceName)
	defer locks.UnlockByName(id.DatabaseAccountName, CosmosDbAccountResourceName)

	future, err := client.DeleteSQLRoleAssignment(ctx, id.Name, id.ResourceGroup, id.DatabaseAccountName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the completion of the deleting of %s: %+v", id, err)
	}

	return nil
}
