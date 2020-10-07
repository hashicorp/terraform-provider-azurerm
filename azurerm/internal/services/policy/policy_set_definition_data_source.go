package policy

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-09-01/policy"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceArmPolicySetDefinition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPolicySetDefinitionRead,

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

			"management_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"metadata": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"parameters": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"policy_definitions": { // TODO -- remove in the next major version
				Type:     schema.TypeString,
				Computed: true,
			},

			"policy_definition_reference": { // TODO -- rename this back to `policy_definition` after the deprecation
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_definition_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"parameters": { // TODO -- remove this attribute in the next major version
							Type:     schema.TypeMap,
							Computed: true,
						},

						"parameter_values": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"reference_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"policy_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmPolicySetDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.SetDefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	displayName := d.Get("display_name").(string)
	managementGroupID := d.Get("management_group_name").(string)

	var setDefinition policy.SetDefinition
	var err error

	// we marked `display_name` and `name` as `ExactlyOneOf`, therefore there will only be one of display_name and name that have non-empty value here
	if displayName != "" {
		setDefinition, err = getPolicySetDefinitionByDisplayName(ctx, client, displayName, managementGroupID)
		if err != nil {
			return fmt.Errorf("reading Policy Set Definition (Display Name %q): %+v", displayName, err)
		}
	}
	if name != "" {
		setDefinition, err = getPolicySetDefinitionByName(ctx, client, name, managementGroupID)
		if err != nil {
			return fmt.Errorf("reading Policy Set Definition %q: %+v", name, err)
		}
	}

	if setDefinition.ID == nil || *setDefinition.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Policy Set Definition %q", name)
	}
	d.SetId(*setDefinition.ID)
	d.Set("name", setDefinition.Name)
	d.Set("display_name", setDefinition.DisplayName)
	d.Set("description", setDefinition.Description)
	d.Set("policy_type", setDefinition.PolicyType)
	d.Set("metadata", flattenJSON(setDefinition.Metadata))

	if paramsStr, err := flattenParameterDefintionsValueToString(setDefinition.Parameters); err != nil {
		return fmt.Errorf("flattening JSON for `parameters`: %+v", err)
	} else {
		d.Set("parameters", paramsStr)
	}

	definitionBytes, err := json.Marshal(setDefinition.PolicyDefinitions)
	if err != nil {
		return fmt.Errorf("flattening JSON for `policy_defintions`: %+v", err)
	}
	d.Set("policy_definitions", string(definitionBytes))

	references, err := flattenAzureRMPolicySetDefinitionPolicyDefinitions(setDefinition.PolicyDefinitions)
	if err != nil {
		return fmt.Errorf("flattening `policy_definition_reference`: %+v", err)
	}
	if err := d.Set("policy_definition_reference", references); err != nil {
		return fmt.Errorf("setting `policy_definition_reference`: %+v", err)
	}

	return nil
}
