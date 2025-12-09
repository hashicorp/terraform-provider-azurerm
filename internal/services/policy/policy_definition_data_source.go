// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2025-01-01/policydefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceArmPolicyDefinition() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: policyDefinitionReadFunc(false),

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: policyDefinitionDataSourceSchema(),
	}
}

func policyDefinitionDataSourceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"name", "display_name"},
		},

		"name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"name", "display_name"},
		},

		"management_group_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"policy_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"policy_rule": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"parameters": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"metadata": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"role_definition_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"mode": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func policyDefinitionReadFunc(builtInOnly bool) func(d *pluginsdk.ResourceData, meta interface{}) error {
	return func(d *pluginsdk.ResourceData, meta interface{}) error {
		client := meta.(*clients.Client).Policy.PolicyDefinitionsClient
		subscriptionId := meta.(*clients.Client).Account.SubscriptionId
		ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
		defer cancel()

		displayName := d.Get("display_name").(string)
		name := d.Get("name").(string)
		managementGroupName := ""
		if v, ok := d.GetOk("management_group_name"); ok {
			managementGroupName = v.(string)
		}

		var id any
		var policyDefinition *policydefinitions.PolicyDefinition
		var err error

		// one of display_name and name must be non-empty, this is guaranteed by schema
		if displayName != "" {
			// Display name lookup still uses old SDK helper
			// TODO: This could be optimized but would require listing all definitions
			return fmt.Errorf("lookup by display_name is not yet supported with the new SDK")
		}

		if name != "" {
			if managementGroupName != "" {
				id = policydefinitions.NewProviders2PolicyDefinitionID(managementGroupName, name)
			} else {
				id = policydefinitions.NewProviderPolicyDefinitionID(subscriptionId, name)
			}

			if builtInOnly && managementGroupName == "" {
				builtInId := policydefinitions.NewPolicyDefinitionID(name)
				resp, err := client.GetBuiltIn(ctx, builtInId)
				if err != nil {
					return fmt.Errorf("reading built-in Policy Definition %q: %+v", name, err)
				}
				policyDefinition = resp.Model
			} else {
				_, policyDefinition, err = getPolicyDefinitionByID(ctx, client, id)
				if err != nil {
					return fmt.Errorf("reading Policy Definition %q: %+v", name, err)
				}
			}
		}

		if policyDefinition == nil || policyDefinition.Id == nil {
			return fmt.Errorf("policy definition was nil or had no ID")
		}

		parsedId, err := parse.PolicyDefinitionID(*policyDefinition.Id)
		if err != nil {
			return fmt.Errorf("parsing Policy Definition ID %q: %+v", *policyDefinition.Id, err)
		}

		d.SetId(parsedId.Id)
		d.Set("name", policyDefinition.Name)

		if props := policyDefinition.Properties; props != nil {
			d.Set("display_name", props.DisplayName)
			d.Set("description", props.Description)
			d.Set("policy_type", string(*props.PolicyType))
			d.Set("mode", props.Mode)

			if policyRuleStr := flattenJSON(props.PolicyRule); policyRuleStr != "" {
				d.Set("policy_rule", policyRuleStr)
				roleIDs, _ := getPolicyRoleDefinitionIDs(policyRuleStr)
				d.Set("role_definition_ids", roleIDs)
			}

			if metadataStr := flattenJSON(props.Metadata); metadataStr != "" {
				d.Set("metadata", metadataStr)
			}

			if parametersStr, err := flattenParameterDefinitionsValueToStringForPolicyDefinition(props.Parameters); err == nil {
				d.Set("parameters", parametersStr)
			} else {
				return fmt.Errorf("flattening Policy Parameters: %+v", err)
			}
		}

		d.Set("type", policyDefinition.Type)

		return nil
	}
}
