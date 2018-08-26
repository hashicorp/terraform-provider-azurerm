package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-01-01-preview/authorization"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRoleDefinition() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRoleDefinitionCreateUpdate,
		Read:   resourceArmRoleDefinitionRead,
		Update: resourceArmRoleDefinitionCreateUpdate,
		Delete: resourceArmRoleDefinitionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
					},
				},
			},

			"assignable_scopes": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceArmRoleDefinitionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).roleDefinitionsClient
	ctx := meta.(*ArmClient).StopContext

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

	properties := authorization.RoleDefinition{
		RoleDefinitionProperties: &authorization.RoleDefinitionProperties{
			RoleName:         utils.String(name),
			Description:      utils.String(description),
			RoleType:         utils.String(roleType),
			Permissions:      &permissions,
			AssignableScopes: &assignableScopes,
		},
	}

	_, err := client.CreateOrUpdate(ctx, scope, roleDefinitionId, properties)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, scope, roleDefinitionId)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Role Definition ID for %q (Scope %q)", name, scope)
	}

	d.SetId(*read.ID)
	return resourceArmRoleDefinitionRead(d, meta)
}

func resourceArmRoleDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).roleDefinitionsClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.GetByID(ctx, d.Id())
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
	client := meta.(*ArmClient).roleDefinitionsClient
	ctx := meta.(*ArmClient).StopContext

	roleDefinitionId := d.Get("role_definition_id").(string)
	scope := d.Get("scope").(string)

	resp, err := client.Delete(ctx, scope, roleDefinitionId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}

		return fmt.Errorf("Error deleting Role Definition %q: %+v", roleDefinitionId, err)
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

		notActionsOutput := make([]string, 0)
		notActions := input["not_actions"].([]interface{})
		for _, a := range notActions {
			notActionsOutput = append(notActionsOutput, a.(string))
		}
		permission.NotActions = &notActionsOutput

		output = append(output, permission)
	}

	return output
}

func expandRoleDefinitionAssignableScopes(d *schema.ResourceData) []string {
	scopes := make([]string, 0)

	assignableScopes := d.Get("assignable_scopes").([]interface{})
	for _, scope := range assignableScopes {
		scopes = append(scopes, scope.(string))
	}

	return scopes
}

func flattenRoleDefinitionPermissions(input *[]authorization.Permission) []interface{} {
	permissions := make([]interface{}, 0)

	for _, permission := range *input {
		output := make(map[string]interface{}, 0)

		actions := make([]string, 0)
		if permission.Actions != nil {
			for _, action := range *permission.Actions {
				actions = append(actions, action)
			}
		}
		output["actions"] = actions

		notActions := make([]string, 0)
		if permission.NotActions != nil {
			for _, action := range *permission.NotActions {
				notActions = append(notActions, action)
			}
		}
		output["not_actions"] = notActions

		permissions = append(permissions, output)
	}

	return permissions
}

func flattenRoleDefinitionAssignableScopes(input *[]string) []interface{} {
	scopes := make([]interface{}, 0)

	for _, scope := range *input {
		scopes = append(scopes, scope)
	}

	return scopes
}
