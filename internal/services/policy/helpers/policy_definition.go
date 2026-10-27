// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy" //nolint:staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//nolint:staticcheck

//nolint:staticcheck

func FlattenJSON(stringMap any) string {
	if stringMap != nil {
		if v, ok := stringMap.(*any); ok {
			stringMap = *v
		}
		value := stringMap.(map[string]any)
		jsonString, err := pluginsdk.FlattenJsonToString(value)
		if err == nil {
			return jsonString
		}
	}

	return ""
}

func PolicyDefinitionDataSourceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"name", "display_name"},
		},

		"name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
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

func PolicyDefinitionReadFunc(builtInOnly bool) func(d *pluginsdk.ResourceData, meta any) error {
	return func(d *pluginsdk.ResourceData, meta any) error {
		client := meta.(*clients.Client).Policy.DefinitionsClient
		ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
		defer cancel()

		displayName := d.Get("display_name").(string)
		name := d.Get("name").(string)
		managementGroupName := ""
		if v, ok := d.GetOk("management_group_name"); ok {
			managementGroupName = v.(string)
		}

		var policyDefinition policy.Definition
		var err error
		// one of display_name and name must be non-empty, this is guaranteed by schema
		if displayName != "" {
			policyDefinition, err = GetPolicyDefinitionByDisplayName(ctx, client, displayName, managementGroupName, builtInOnly)
			if err != nil {
				return fmt.Errorf("reading Policy Definition (Display Name %q): %+v", displayName, err)
			}
		}
		if name != "" {
			if builtInOnly && managementGroupName == "" {
				policyDefinition, err = client.GetBuiltIn(ctx, name)
			} else {
				policyDefinition, err = GetPolicyDefinitionByName(ctx, client, name, managementGroupName)
			}
			if err != nil {
				return fmt.Errorf("reading Policy Definition %q: %+v", name, err)
			}
		}

		id, err := parse.PolicyDefinitionID(*policyDefinition.ID)
		if err != nil {
			return fmt.Errorf("parsing Policy Definition %q: %+v", *policyDefinition.ID, err)
		}

		d.SetId(id.Id)
		d.Set("name", policyDefinition.Name)
		d.Set("display_name", policyDefinition.DisplayName)
		d.Set("description", policyDefinition.Description)
		d.Set("type", policyDefinition.Type)
		d.Set("policy_type", policyDefinition.PolicyType)
		d.Set("mode", policyDefinition.Mode)

		if policyRuleStr := FlattenJSON(policyDefinition.PolicyRule.(map[string]any)); policyRuleStr != "" {
			d.Set("policy_rule", policyRuleStr)
			roleIDs, _ := GetPolicyRoleDefinitionIDs(policyRuleStr)
			d.Set("role_definition_ids", roleIDs)
		} else {
			return fmt.Errorf("flattening Policy Definition Rule %q: %+v", name, err)
		}

		if metadataStr := FlattenJSON(policyDefinition.Metadata); metadataStr != "" {
			d.Set("metadata", metadataStr)
		}

		if parametersStr, err := FlattenParameterDefinitionsValueToStringTrack1(policyDefinition.Parameters); err == nil {
			d.Set("parameters", parametersStr)
		} else {
			return fmt.Errorf("failed to flatten Policy Parameters %q: %+v", name, err)
		}

		return nil
	}
}
