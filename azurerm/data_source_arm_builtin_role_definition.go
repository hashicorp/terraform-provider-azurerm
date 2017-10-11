package azurerm

import (
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
		},
	}
}

func dataSourceArmBuiltInRoleDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	roleDefinitionIds := map[string]string{
		"Contributor":               "/providers/Microsoft.Authorization/roledefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c",
		"Owner":                     "/providers/Microsoft.Authorization/roledefinitions/8e3af657-a8ff-443c-a75c-2fe8c4bcb635",
		"Reader":                    "/providers/Microsoft.Authorization/roledefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
		"VirtualMachineContributor": "/providers/Microsoft.Authorization/roledefinitions/d73bb868-a0df-4d4d-bd69-98a00b01fccb",
	}
	roleDefinitionId := roleDefinitionIds[name]

	// TODO: when the API's fixed - pull out additional information from the API
	// https://github.com/Azure/azure-rest-api-specs/issues/1785

	d.SetId(roleDefinitionId)

	return nil
}
