package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/policy"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmBuiltInPolicyDefinition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmBuiltInPolicyDefinitionRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"management_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func dataSourceArmBuiltInPolicyDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).policyDefinitionsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	managementGroupID := d.Get("management_group_id").(string)

	var policyDefinitions policy.DefinitionListResultPage
	var err error

	if managementGroupID != "" {
		policyDefinitions, err = client.ListByManagementGroup(ctx, managementGroupID)
	} else {
		policyDefinitions, err = client.ListBuiltIn(ctx)
	}

	if err != nil {
		return fmt.Errorf("Error loading Policy Definition List: %+v", err)
	}

	var policyDefinition policy.Definition

	for policyDefinition.ID == nil && policyDefinitions.NotDone() {
		vals := policyDefinitions.Values()
		for index := 0; index < len(vals); index++ {
			def := vals[index]

			if *def.DisplayName == name {
				policyDefinition = def
				break
			}
		}
	}

	if policyDefinition.ID == nil {
		return fmt.Errorf("Error loading Policy Definition List: could not find policy '%s'", name)
	}

	d.SetId(*policyDefinition.ID)

	return nil
}
