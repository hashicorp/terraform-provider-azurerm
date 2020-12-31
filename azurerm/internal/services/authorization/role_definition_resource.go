package authorization

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-09-01-preview/authorization"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/authorization/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRoleDefinition() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRoleDefinitionCreateUpdate,
		Read:   resourceArmRoleDefinitionRead,
		Update: resourceArmRoleDefinitionCreateUpdate,
		Delete: resourceArmRoleDefinitionDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.RoleDefinitionId(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,

		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceArmRoleDefinitionV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceArmRoleDefinitionStateUpgradeV0,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			"role_definition_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"scope": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"permissions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"not_actions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"data_actions": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: schema.HashString,
						},
						"not_data_actions": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: schema.HashString,
						},
					},
				},
			},

			"assignable_scopes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"role_definition_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmRoleDefinitionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Authorization.RoleDefinitionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	roleDefinitionId := d.Get("role_definition_id").(string)
	if roleDefinitionId == "" {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("Error generating UUID for Role Assignment: %+v", err)
		}

		roleDefinitionId = uuid
	}

	name := d.Get("name").(string)
	scope := d.Get("scope").(string)
	description := d.Get("description").(string)
	roleType := "CustomRole"
	permissions := expandRoleDefinitionPermissions(d)
	assignableScopes := expandRoleDefinitionAssignableScopes(d)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, scope, roleDefinitionId)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Role Definition ID for %q (Scope %q)", name, scope)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			importID := fmt.Sprintf("%s|%s", *existing.ID, scope)
			return tf.ImportAsExistsError("azurerm_role_definition", importID)
		}
	}

	properties := authorization.RoleDefinition{
		RoleDefinitionProperties: &authorization.RoleDefinitionProperties{
			RoleName:         utils.String(name),
			Description:      utils.String(description),
			RoleType:         utils.String(roleType),
			Permissions:      &permissions,
			AssignableScopes: &assignableScopes,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, scope, roleDefinitionId, properties); err != nil {
		return err
	}

	read, err := client.Get(ctx, scope, roleDefinitionId)
	if err != nil {
		return err
	}
	if read.ID == nil || *read.ID == "" {
		return fmt.Errorf("Cannot read Role Definition ID for %q (Scope %q)", name, scope)
	}

	d.SetId(fmt.Sprintf("%s|%s", *read.ID, scope))
	return resourceArmRoleDefinitionRead(d, meta)
}

func resourceArmRoleDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Authorization.RoleDefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	roleDefinitionId, err := parse.RoleDefinitionId(d.Id())
	if err != nil {
		return err
	}

	d.Set("scope", roleDefinitionId.Scope)
	d.Set("role_definition_id", roleDefinitionId.RoleID)
	d.Set("role_definition_resource_id", roleDefinitionId.ResourceID)

	resp, err := client.Get(ctx, roleDefinitionId.Scope, roleDefinitionId.RoleID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Role Definition %q was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error loading Role Definition %q: %+v", d.Id(), err)
	}

	if props := resp.RoleDefinitionProperties; props != nil {
		d.Set("name", props.RoleName)
		d.Set("description", props.Description)

		permissions := flattenRoleDefinitionPermissions(props.Permissions)
		if err := d.Set("permissions", permissions); err != nil {
			return err
		}

		assignableScopes := flattenRoleDefinitionAssignableScopes(props.AssignableScopes)
		if err := d.Set("assignable_scopes", assignableScopes); err != nil {
			return err
		}
	}

	return nil
}

func resourceArmRoleDefinitionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Authorization.RoleDefinitionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, _ := parse.RoleDefinitionId(d.Id())

	resp, err := client.Delete(ctx, id.Scope, id.RoleID)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("deleting Role Definition %q at Scope %q: %+v", id.RoleID, id.Scope, err)
		}
	}

	return nil
}

func expandRoleDefinitionPermissions(d *schema.ResourceData) []authorization.Permission {
	output := make([]authorization.Permission, 0)

	permissions := d.Get("permissions").([]interface{})
	for _, v := range permissions {
		input := v.(map[string]interface{})
		permission := authorization.Permission{}

		actionsOutput := make([]string, 0)
		actions := input["actions"].([]interface{})
		for _, a := range actions {
			actionsOutput = append(actionsOutput, a.(string))
		}
		permission.Actions = &actionsOutput

		dataActionsOutput := make([]string, 0)
		dataActions := input["data_actions"].(*schema.Set)
		for _, a := range dataActions.List() {
			dataActionsOutput = append(dataActionsOutput, a.(string))
		}
		permission.DataActions = &dataActionsOutput

		notActionsOutput := make([]string, 0)
		notActions := input["not_actions"].([]interface{})
		for _, a := range notActions {
			notActionsOutput = append(notActionsOutput, a.(string))
		}
		permission.NotActions = &notActionsOutput

		notDataActionsOutput := make([]string, 0)
		notDataActions := input["not_data_actions"].(*schema.Set)
		for _, a := range notDataActions.List() {
			notDataActionsOutput = append(notDataActionsOutput, a.(string))
		}
		permission.NotDataActions = &notDataActionsOutput

		output = append(output, permission)
	}

	return output
}

func expandRoleDefinitionAssignableScopes(d *schema.ResourceData) []string {
	scopes := make([]string, 0)

	// The first scope in the list must be the target scope as it it not returned in any API call
	assignedScope := d.Get("scope").(string)
	scopes = append(scopes, assignedScope)
	assignableScopes := d.Get("assignable_scopes").([]interface{})
	for _, scope := range assignableScopes {
		// Ensure the assigned scope is not duplicated in the list if also specified in `assignable_scopes`
		if scope != assignedScope {
			scopes = append(scopes, scope.(string))
		}
	}

	return scopes
}

func flattenRoleDefinitionPermissions(input *[]authorization.Permission) []interface{} {
	permissions := make([]interface{}, 0)
	if input == nil {
		return permissions
	}

	for _, permission := range *input {
		output := make(map[string]interface{})

		actions := make([]string, 0)
		if s := permission.Actions; s != nil {
			actions = *s
		}
		output["actions"] = actions

		dataActions := make([]interface{}, 0)
		if permission.DataActions != nil {
			for _, dataAction := range *permission.DataActions {
				dataActions = append(dataActions, dataAction)
			}
		}
		output["data_actions"] = schema.NewSet(schema.HashString, dataActions)

		notActions := make([]string, 0)
		if s := permission.NotActions; s != nil {
			notActions = *s
		}
		output["not_actions"] = notActions

		notDataActions := make([]interface{}, 0)
		if permission.NotDataActions != nil {
			for _, dataAction := range *permission.NotDataActions {
				notDataActions = append(notDataActions, dataAction)
			}
		}
		output["not_data_actions"] = schema.NewSet(schema.HashString, notDataActions)

		permissions = append(permissions, output)
	}

	return permissions
}

func flattenRoleDefinitionAssignableScopes(input *[]string) []interface{} {
	scopes := make([]interface{}, 0)
	if input == nil {
		return scopes
	}

	for _, scope := range *input {
		scopes = append(scopes, scope)
	}

	return scopes
}
