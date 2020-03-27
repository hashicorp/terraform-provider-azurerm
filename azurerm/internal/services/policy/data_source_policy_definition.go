package policy

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/policy"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"name", "display_name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"name", "display_name"},
			},
			"management_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				// TODO -- temporary removed this validation, since a management group ID is always not a resource ID. add this back when there is a proper function for validation of mgmt group IDs
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

	displayName := d.Get("display_name").(string)
	name := d.Get("name").(string)
	managementGroupID := d.Get("management_group_id").(string)

	var policyDefinition policy.Definition
	var err error
	if displayName != "" {
		policyDefinition, err = getPolicyDefinitionByDisplayName(ctx, client, displayName, managementGroupID)
		if err != nil {
			return fmt.Errorf("failed to read Policy Definition (Display Name %q): %+v", displayName, err)
		}
	} else if name != "" {
		policyDefinition, err = getPolicyDefinition(ctx, client, name, managementGroupID)
		if err != nil {
			return fmt.Errorf("failed to read Policy Definition %q: %+v", name, err)
		}
	} else {
		return fmt.Errorf("one of `display_name` or `name` must be set")
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

func getPolicyDefinitionByDisplayName(ctx context.Context, client *policy.DefinitionsClient, displayName, managementGroupID string) (policy.Definition, error) {
	var policyDefinitions policy.DefinitionListResultIterator
	var err error

	if managementGroupID != "" {
		policyDefinitions, err = client.ListByManagementGroupComplete(ctx, managementGroupID)
	} else {
		policyDefinitions, err = client.ListComplete(ctx)
	}
	if err != nil {
		return policy.Definition{}, fmt.Errorf("failed to load Policy Definition List: %+v", err)
	}

	for policyDefinitions.NotDone() {
		def := policyDefinitions.Value()
		if def.DisplayName != nil && *def.DisplayName == displayName && def.ID != nil {
			return def, nil
		}

		if err := policyDefinitions.NextWithContext(ctx); err != nil {
			return policy.Definition{}, fmt.Errorf("failed to load Policy Definition List: %s", err)
		}
	}

	return policy.Definition{}, fmt.Errorf("failed to load Policy Definition List: could not find policy '%s'", displayName)
}
