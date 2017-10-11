package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/authorization"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceArmBuiltInRoleDefinition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmBuiltInRoleDefinitionRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Contributor",
					"Reader",
					"Owner",
					"VirtualMachineContributor",
				}, false),
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
		},
	}
}

func dataSourceArmBuiltInRoleDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).roleDefinitionsClient
	name := d.Get("name").(string)
	roleDefinitionIds := map[string]string{
		"Contributor":               "/providers/Microsoft.Authorization/roledefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c",
		"Owner":                     "/providers/Microsoft.Authorization/roledefinitions/8e3af657-a8ff-443c-a75c-2fe8c4bcb635",
		"Reader":                    "/providers/Microsoft.Authorization/roledefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
		"VirtualMachineContributor": "/providers/Microsoft.Authorization/roledefinitions/d73bb868-a0df-4d4d-bd69-98a00b01fccb",
	}
	roleDefinitionId := roleDefinitionIds[name]

	d.SetId(roleDefinitionId)

	role, err := client.GetByID(roleDefinitionId)
	if err != nil {
		return fmt.Errorf("Error loadng Role Definition: %+v", err)
	}

	if props := role.Properties; props != nil {
		d.Set("name", props.RoleName)
		d.Set("description", props.Description)
		d.Set("type", props.Type)

		permissions := flattenRoleDefinitionPermissions(props.Permissions)
		if err := d.Set("permissions", permissions); err != nil {
			return err
		}
	}

	return nil
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
