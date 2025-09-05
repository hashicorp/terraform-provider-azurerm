// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceArmPolicySetDefinition() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmPolicySetDefinitionRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
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

			"description": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"metadata": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"parameters": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"policy_definitions": { // TODO -- remove in the next major version
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"policy_definition_reference": { // TODO -- rename this back to `policy_definition` after the deprecation
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"policy_definition_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"parameters": { // TODO -- remove this attribute in the next major version
							Type:     pluginsdk.TypeMap,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"parameter_values": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"reference_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"policy_group_names": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"policy_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"policy_definition_group": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"display_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"category": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"description": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"additional_metadata_resource_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmPolicySetDefinitionRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

	id, err := parse.PolicySetDefinitionID(*setDefinition.ID)
	if err != nil {
		return fmt.Errorf("parsing Policy Set Definition %q: %+v", *setDefinition.ID, err)
	}

	d.SetId(id.Id)
	d.Set("name", setDefinition.Name)
	d.Set("display_name", setDefinition.DisplayName)
	d.Set("description", setDefinition.Description)
	d.Set("policy_type", setDefinition.PolicyType)
	d.Set("metadata", flattenJSON(setDefinition.Metadata))

	if paramsStr, err := flattenParameterDefinitionsValueToStringTrack1(setDefinition.Parameters); err != nil {
		return fmt.Errorf("flattening JSON for `parameters`: %+v", err)
	} else {
		d.Set("parameters", paramsStr)
	}

	definitionBytes, err := json.Marshal(setDefinition.PolicyDefinitions)
	if err != nil {
		return fmt.Errorf("flattening JSON for `policy_defintions`: %+v", err)
	}
	d.Set("policy_definitions", string(definitionBytes))

	references, err := flattenAzureRMPolicySetDefinitionPolicyDefinitionsTrack1(setDefinition.PolicyDefinitions)
	if err != nil {
		return fmt.Errorf("flattening `policy_definition_reference`: %+v", err)
	}
	if err := d.Set("policy_definition_reference", references); err != nil {
		return fmt.Errorf("setting `policy_definition_reference`: %+v", err)
	}

	if err := d.Set("policy_definition_group", flattenAzureRMPolicySetDefinitionPolicyGroupsTrack1(setDefinition.PolicyDefinitionGroups)); err != nil {
		return fmt.Errorf("setting `policy_definition_group`: %+v", err)
	}

	return nil
}

func flattenAzureRMPolicySetDefinitionPolicyDefinitionsTrack1(input *[]policy.DefinitionReference) ([]interface{}, error) {
	result := make([]interface{}, 0)
	if input == nil {
		return result, nil
	}

	for _, definition := range *input {
		policyDefinitionID := ""
		if definition.PolicyDefinitionID != nil {
			policyDefinitionID = *definition.PolicyDefinitionID
		}

		parametersMap := make(map[string]interface{})
		for k, v := range definition.Parameters {
			if v == nil {
				continue
			}
			parametersMap[k] = fmt.Sprintf("%v", v.Value) // map in terraform only accepts string as its values, therefore we have to convert the value to string
		}

		parameterValues, err := flattenParameterValuesValueToStringTrack1(definition.Parameters)
		if err != nil {
			return nil, fmt.Errorf("serializing JSON from `parameter_values`: %+v", err)
		}

		policyDefinitionReference := ""
		if definition.PolicyDefinitionReferenceID != nil {
			policyDefinitionReference = *definition.PolicyDefinitionReferenceID
		}

		result = append(result, map[string]interface{}{
			"policy_definition_id": policyDefinitionID,
			"parameter_values":     parameterValues,
			"reference_id":         policyDefinitionReference,
			"policy_group_names":   utils.FlattenStringSlice(definition.GroupNames),
		})
	}
	return result, nil
}

func flattenAzureRMPolicySetDefinitionPolicyGroupsTrack1(input *[]policy.DefinitionGroup) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	for _, group := range *input {
		name := ""
		if group.Name != nil {
			name = *group.Name
		}
		displayName := ""
		if group.DisplayName != nil {
			displayName = *group.DisplayName
		}
		category := ""
		if group.Category != nil {
			category = *group.Category
		}
		description := ""
		if group.Description != nil {
			description = *group.Description
		}
		metadataID := ""
		if group.AdditionalMetadataID != nil {
			metadataID = *group.AdditionalMetadataID
		}

		result = append(result, map[string]interface{}{
			"name":                            name,
			"display_name":                    displayName,
			"category":                        category,
			"description":                     description,
			"additional_metadata_resource_id": metadataID,
		})
	}

	return result
}
