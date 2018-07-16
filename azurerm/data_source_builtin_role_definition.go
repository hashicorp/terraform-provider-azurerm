package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmBuiltInRoleDefinition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmBuiltInRoleDefinitionRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permissions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"not_actions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"assignable_scopes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceArmBuiltInRoleDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).roleDefinitionsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	if name == "VirtualMachineContributor" {
		name = "Virtual Machine Contributor"
	}
	filter := fmt.Sprintf("roleName eq '%s'", name)
	roleDefinitions, err := client.List(ctx, "", filter)
	if err != nil {
		return fmt.Errorf("Error loading Role Definition List: %+v", err)
	}
	if len(roleDefinitions.Values()) != 1 {
		return fmt.Errorf("Error loading Role Definition List: could not find role '%s'", name)
	}

	roleDefinitionId := *roleDefinitions.Values()[0].ID

	d.SetId(roleDefinitionId)

	role, err := client.GetByID(ctx, roleDefinitionId)
	if err != nil {
		return fmt.Errorf("Error loading Role Definition: %+v", err)
	}

	if props := role.RoleDefinitionProperties; props != nil {
		d.Set("name", props.RoleName)
		d.Set("description", props.Description)
		d.Set("type", props.RoleType)

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
