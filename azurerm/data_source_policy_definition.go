package azurerm

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/policy"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceArmPolicyDefinition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPolicyDefinitionRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"management_group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_rule": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parameters": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmPolicyDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.DefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("display_name").(string)
	managementGroupID := d.Get("management_group_id").(string)

	var policyDefinitions policy.DefinitionListResultIterator
	var err error

	if managementGroupID != "" {
		policyDefinitions, err = client.ListByManagementGroupComplete(ctx, managementGroupID)
	} else {
		policyDefinitions, err = client.ListComplete(ctx)
	}

	if err != nil {
		return fmt.Errorf("Error loading Policy Definition List: %+v", err)
	}

	var policyDefinition policy.Definition

	for policyDefinitions.NotDone() {
		def := policyDefinitions.Value()
		if def.DisplayName != nil && *def.DisplayName == name {
			policyDefinition = def
			break
		}

		err = policyDefinitions.NextWithContext(ctx)
		if err != nil {
			return fmt.Errorf("Error loading Policy Definition List: %s", err)
		}
	}

	if policyDefinition.ID == nil {
		return fmt.Errorf("Error loading Policy Definition List: could not find policy '%s'", name)
	}

	d.SetId(*policyDefinition.ID)
	d.Set("name", policyDefinition.Name)
	d.Set("display_name", policyDefinition.DisplayName)
	d.Set("description", policyDefinition.Description)
	d.Set("type", policyDefinition.Type)
	d.Set("policy_type", policyDefinition.PolicyType)

	if policyRuleStr := flattenJSON(policyDefinition.PolicyRule); policyRuleStr != "" {
		d.Set("policy_rule", policyRuleStr)
	}

	if metadataStr := flattenJSON(policyDefinition.Metadata); metadataStr != "" {
		d.Set("metadata", metadataStr)
	}

	if parametersStr := flattenJSON(policyDefinition.Parameters); parametersStr != "" {
		d.Set("parameters", parametersStr)
	}

	return nil
}
