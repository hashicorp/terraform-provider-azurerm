// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy" // nolint: staticcheck
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	mgmtGrpParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmPolicySetDefinition() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmPolicySetDefinitionCreate,
		Update: resourceArmPolicySetDefinitionUpdate,
		Read:   resourceArmPolicySetDefinitionRead,
		Delete: resourceArmPolicySetDefinitionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PolicySetDefinitionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourcePolicySetDefinitionSchema(),
	}
}

func resourcePolicySetDefinitionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"policy_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(policy.TypeBuiltIn),
				string(policy.TypeCustom),
				string(policy.TypeNotSpecified),
				string(policy.TypeStatic),
			}, false),
		},

		"management_group_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"metadata": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: policySetDefinitionsMetadataDiffSuppressFunc,
		},

		"parameters": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		// lintignore: S013
		"policy_definition_reference": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"policy_definition_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.PolicyDefinitionID,
					},

					"parameter_values": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						ValidateFunc:     validation.StringIsJSON,
						DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
					},

					"reference_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"policy_group_names": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},

		"policy_definition_group": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"display_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"category": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"description": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"additional_metadata_resource_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
			Set: resourceARMPolicySetDefinitionPolicyDefinitionGroupHash,
		},
	}
}

func policySetDefinitionsMetadataDiffSuppressFunc(_, old, new string, _ *pluginsdk.ResourceData) bool {
	var oldPolicySetDefinitionsMetadata map[string]interface{}
	errOld := json.Unmarshal([]byte(old), &oldPolicySetDefinitionsMetadata)
	if errOld != nil {
		return false
	}

	var newPolicySetDefinitionsMetadata map[string]interface{}
	errNew := json.Unmarshal([]byte(new), &newPolicySetDefinitionsMetadata)
	if errNew != nil {
		return false
	}

	// Ignore the following keys if they're found in the metadata JSON
	ignoreKeys := [4]string{"createdBy", "createdOn", "updatedBy", "updatedOn"}
	for _, key := range ignoreKeys {
		delete(oldPolicySetDefinitionsMetadata, key)
		delete(newPolicySetDefinitionsMetadata, key)
	}

	return reflect.DeepEqual(oldPolicySetDefinitionsMetadata, newPolicySetDefinitionsMetadata)
}

type DefinitionReferenceInOldApiVersion struct {
	// PolicyDefinitionID - The ID of the policy definition or policy set definition.
	PolicyDefinitionID *string `json:"policyDefinitionId,omitempty"`
	// Parameters - The parameter values for the referenced policy rule. The keys are the parameter names.
	Parameters map[string]*policy.ParameterValuesValue `json:"parameters"`
}

func resourceArmPolicySetDefinitionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.SetDefinitionsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	managementGroupName := ""
	if v, ok := d.GetOk("management_group_id"); ok {
		managementGroupID, err := mgmtGrpParse.ManagementGroupID(v.(string))
		if err != nil {
			return err
		}
		managementGroupName = managementGroupID.Name
	}

	existing, err := getPolicySetDefinitionByName(ctx, client, name, managementGroupName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Policy Set Definition %q: %+v", name, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_policy_set_definition", *existing.ID)
	}

	properties := policy.SetDefinitionProperties{
		PolicyType:  policy.Type(d.Get("policy_type").(string)),
		DisplayName: utils.String(d.Get("display_name").(string)),
		Description: utils.String(d.Get("description").(string)),
	}

	if metaDataString := d.Get("metadata").(string); metaDataString != "" {
		metaData, err := pluginsdk.ExpandJsonFromString(metaDataString)
		if err != nil {
			return fmt.Errorf("expanding JSON for `metadata`: %+v", err)
		}
		properties.Metadata = &metaData
	}

	if parametersString := d.Get("parameters").(string); parametersString != "" {
		parameters, err := expandParameterDefinitionsValueFromString(parametersString)
		if err != nil {
			return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
		}
		properties.Parameters = parameters
	}

	if v, ok := d.GetOk("policy_definition_reference"); ok {
		definitions, err := expandAzureRMPolicySetDefinitionPolicyDefinitions(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `policy_definition_reference`: %+v", err)
		}
		properties.PolicyDefinitions = definitions
	}

	if v, ok := d.GetOk("policy_definition_group"); ok {
		properties.PolicyDefinitionGroups = expandAzureRMPolicySetDefinitionPolicyGroups(v.(*pluginsdk.Set).List())
	}

	definition := policy.SetDefinition{
		SetDefinitionProperties: &properties,
	}

	if managementGroupName == "" {
		_, err = client.CreateOrUpdate(ctx, name, definition)
	} else {
		_, err = client.CreateOrUpdateAtManagementGroup(ctx, name, definition, managementGroupName)
	}

	if err != nil {
		return fmt.Errorf("creating Policy Set Definition %q: %+v", name, err)
	}

	// Policy Definitions are eventually consistent; wait for them to stabilize
	log.Printf("[DEBUG] Waiting for Policy Set Definition %q to become available", name)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"404"},
		Target:                    []string{"200"},
		Refresh:                   policySetDefinitionRefreshFunc(ctx, client, name, managementGroupName),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
	}

	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Policy Set Definition %q to become available: %+v", name, err)
	}

	var resp policy.SetDefinition
	resp, err = getPolicySetDefinitionByName(ctx, client, name, managementGroupName)
	if err != nil {
		return fmt.Errorf("retrieving Policy Set Definition %q: %+v", name, err)
	}

	id, err := parse.PolicySetDefinitionID(*resp.ID)
	if err != nil {
		return fmt.Errorf("parsing Policy Set Definition %q: %+v", *resp.ID, err)
	}

	d.SetId(id.Id)

	return resourceArmPolicySetDefinitionRead(d, meta)
}

func resourceArmPolicySetDefinitionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.SetDefinitionsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicySetDefinitionID(d.Id())
	if err != nil {
		return err
	}

	managementGroupName := ""
	var managementGroupId mgmtGrpParse.ManagementGroupId
	if scopeId, ok := id.PolicyScopeId.(parse.ScopeAtManagementGroup); ok {
		managementGroupId = mgmtGrpParse.NewManagementGroupId(scopeId.ManagementGroupName)
		managementGroupName = managementGroupId.Name
	}

	// retrieve
	existing, err := getPolicySetDefinitionByName(ctx, client, id.Name, managementGroupName)
	if err != nil {
		return fmt.Errorf("retrieving Policy Set Definition %q (Scope %q): %+v", id.Name, id.ScopeId(), err)
	}
	if existing.SetDefinitionProperties == nil {
		return fmt.Errorf("retrieving Policy Set Definition %q (Scope %q): `properties` was nil", id.Name, id.ScopeId())
	}

	if d.HasChange("policy_type") {
		existing.SetDefinitionProperties.PolicyType = policy.Type(d.Get("policy_type").(string))
	}

	if d.HasChange("display_name") {
		existing.SetDefinitionProperties.DisplayName = utils.String(d.Get("display_name").(string))
	}

	if d.HasChange("description") {
		existing.SetDefinitionProperties.Description = utils.String(d.Get("description").(string))
	}

	if d.HasChange("metadata") {
		metaDataString := d.Get("metadata").(string)
		if metaDataString != "" {
			metaData, err := pluginsdk.ExpandJsonFromString(metaDataString)
			if err != nil {
				return fmt.Errorf("expanding JSON for `metadata`: %+v", err)
			}
			existing.SetDefinitionProperties.Metadata = metaData
		} else {
			existing.SetDefinitionProperties.Metadata = nil
		}
	}

	if d.HasChange("parameters") {
		parametersString := d.Get("parameters").(string)
		if parametersString != "" {
			parameters, err := expandParameterDefinitionsValueFromString(parametersString)
			if err != nil {
				return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
			}
			existing.SetDefinitionProperties.Parameters = parameters
		} else {
			existing.SetDefinitionProperties.Parameters = nil
		}
	}

	if d.HasChange("policy_definition_group") {
		existing.SetDefinitionProperties.PolicyDefinitionGroups = expandAzureRMPolicySetDefinitionPolicyGroups(d.Get("policy_definition_group").(*pluginsdk.Set).List())
	}

	if d.HasChange("policy_definition_reference") {
		definitions, err := expandAzureRMPolicySetDefinitionPolicyDefinitionsUpdate(d)
		if err != nil {
			return fmt.Errorf("expanding `policy_definition_reference`: %+v", err)
		}
		existing.SetDefinitionProperties.PolicyDefinitions = definitions
	}

	if managementGroupName == "" {
		_, err = client.CreateOrUpdate(ctx, id.Name, existing)
	} else {
		_, err = client.CreateOrUpdateAtManagementGroup(ctx, id.Name, existing, managementGroupName)
	}

	if err != nil {
		return fmt.Errorf("updating Policy Set Definition %q: %+v", id.Name, err)
	}

	var resp policy.SetDefinition
	resp, err = getPolicySetDefinitionByName(ctx, client, id.Name, managementGroupName)
	if err != nil {
		return fmt.Errorf("retrieving Policy Set Definition %q: %+v", id.Name, err)
	}

	id, err = parse.PolicySetDefinitionID(*resp.ID)
	if err != nil {
		return fmt.Errorf("parsing Policy Set Definition %q: %+v", *resp.ID, err)
	}

	d.SetId(id.Id)

	return resourceArmPolicySetDefinitionRead(d, meta)
}

func resourceArmPolicySetDefinitionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.SetDefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicySetDefinitionID(d.Id())
	if err != nil {
		return err
	}

	managementGroupName := ""
	var managementGroupId mgmtGrpParse.ManagementGroupId
	switch scopeId := id.PolicyScopeId.(type) { // nolint gocritic
	case parse.ScopeAtManagementGroup:
		managementGroupId = mgmtGrpParse.NewManagementGroupId(scopeId.ManagementGroupName)
		managementGroupName = managementGroupId.Name
	}

	resp, err := getPolicySetDefinitionByName(ctx, client, id.Name, managementGroupName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Policy Set Definition %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Policy Set Definition %+v", err)
	}

	d.Set("name", resp.Name)
	d.Set("management_group_id", managementGroupName)
	if managementGroupName != "" {
		d.Set("management_group_id", managementGroupId.ID())
	}

	if props := resp.SetDefinitionProperties; props != nil {
		d.Set("policy_type", string(props.PolicyType))
		d.Set("display_name", props.DisplayName)
		d.Set("description", props.Description)

		if metadata := props.Metadata; metadata != nil {
			metadataVal := metadata.(map[string]interface{})
			metadataStr, err := pluginsdk.FlattenJsonToString(metadataVal)
			if err != nil {
				return fmt.Errorf("flattening JSON for `metadata`: %+v", err)
			}

			d.Set("metadata", metadataStr)
		}

		if parameters := props.Parameters; parameters != nil {
			parametersStr, err := flattenParameterDefinitionsValueToString(parameters)
			if err != nil {
				return fmt.Errorf("flattening JSON for `parameters`: %+v", err)
			}

			d.Set("parameters", parametersStr)
		}

		references, err := flattenAzureRMPolicySetDefinitionPolicyDefinitions(props.PolicyDefinitions)
		if err != nil {
			return fmt.Errorf("flattening `policy_definition_reference`: %+v", err)
		}
		if err := d.Set("policy_definition_reference", references); err != nil {
			return fmt.Errorf("setting `policy_definition_reference`: %+v", err)
		}

		if err := d.Set("policy_definition_group", flattenAzureRMPolicySetDefinitionPolicyGroups(props.PolicyDefinitionGroups)); err != nil {
			return fmt.Errorf("setting `policy_definition_group`: %+v", err)
		}
	}

	return nil
}

func resourceArmPolicySetDefinitionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.SetDefinitionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicySetDefinitionID(d.Id())
	if err != nil {
		return err
	}

	managementGroupName := ""
	if scopeId, ok := id.PolicyScopeId.(parse.ScopeAtManagementGroup); ok {
		managementGroupName = scopeId.ManagementGroupName
	}

	var resp autorest.Response
	if managementGroupName == "" {
		resp, err = client.Delete(ctx, id.Name)
	} else {
		resp, err = client.DeleteAtManagementGroup(ctx, id.Name, managementGroupName)
	}

	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("deleting Policy Set Definition %q: %+v", id.Name, err)
	}

	return nil
}

func policySetDefinitionRefreshFunc(ctx context.Context, client *policy.SetDefinitionsClient, name, managementGroupId string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := getPolicySetDefinitionByName(ctx, client, name, managementGroupId)
		if err != nil {
			return nil, strconv.Itoa(res.StatusCode), fmt.Errorf("issuing read request in policySetDefinitionRefreshFunc for Policy Set Definition %q: %+v", name, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}

func expandAzureRMPolicySetDefinitionPolicyDefinitionsUpdate(d *pluginsdk.ResourceData) (*[]policy.DefinitionReference, error) {
	result := make([]policy.DefinitionReference, 0)
	input := d.Get("policy_definition_reference").([]interface{})

	for i := range input {
		if d.HasChange(fmt.Sprintf("policy_definition_reference.%d.parameter_values", i)) && d.HasChange(fmt.Sprintf("policy_definition_reference.%d.parameters", i)) {
			return nil, fmt.Errorf("cannot set both `parameters` and `parameter_values`")
		}
		parameters := make(map[string]*policy.ParameterValuesValue)
		if d.HasChange(fmt.Sprintf("policy_definition_reference.%d.parameters", i)) {
			// there is change in `parameters` - the user is will to use this attribute as parameter values
			log.Printf("[DEBUG] updating %s", fmt.Sprintf("policy_definition_reference.%d.parameters", i))
			p := d.Get(fmt.Sprintf("policy_definition_reference.%d.parameters", i)).(map[string]interface{})
			for k, v := range p {
				parameters[k] = &policy.ParameterValuesValue{
					Value: v,
				}
			}
		} else {
			// in this case, it is either parameter_values updated or no update on both, we took the value in `parameter_values` as the final value
			log.Printf("[DEBUG] updating %s", fmt.Sprintf("policy_definition_reference.%d.parameter_values", i))
			if p, ok := d.Get(fmt.Sprintf("policy_definition_reference.%d.parameter_values", i)).(string); ok && p != "" {
				if err := json.Unmarshal([]byte(p), &parameters); err != nil {
					return nil, fmt.Errorf("unmarshalling `parameter_values`: %+v", err)
				}
			}
		}

		result = append(result, policy.DefinitionReference{
			PolicyDefinitionID:          utils.String(d.Get(fmt.Sprintf("policy_definition_reference.%d.policy_definition_id", i)).(string)),
			Parameters:                  parameters,
			PolicyDefinitionReferenceID: utils.String(d.Get(fmt.Sprintf("policy_definition_reference.%d.reference_id", i)).(string)),
			GroupNames:                  utils.ExpandStringSlice(d.Get(fmt.Sprintf("policy_definition_reference.%d.policy_group_names", i)).(*schema.Set).List()),
		})
	}

	return &result, nil
}

func expandAzureRMPolicySetDefinitionPolicyDefinitions(input []interface{}) (*[]policy.DefinitionReference, error) {
	result := make([]policy.DefinitionReference, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		var parameters map[string]*policy.ParameterValuesValue
		if p, ok := v["parameter_values"].(string); ok && p != "" {
			parameters = make(map[string]*policy.ParameterValuesValue)
			if err := json.Unmarshal([]byte(p), &parameters); err != nil {
				return nil, fmt.Errorf("unmarshalling `parameter_values`: %+v", err)
			}
		}
		if p, ok := v["parameters"].(map[string]interface{}); ok {
			if len(parameters) > 0 && len(p) > 0 {
				return nil, fmt.Errorf("cannot set both `parameters` and `parameter_values`")
			}
			parameters = make(map[string]*policy.ParameterValuesValue)
			for k, value := range p {
				parameters[k] = &policy.ParameterValuesValue{
					Value: value,
				}
			}
		}

		result = append(result, policy.DefinitionReference{
			PolicyDefinitionID:          utils.String(v["policy_definition_id"].(string)),
			Parameters:                  parameters,
			PolicyDefinitionReferenceID: utils.String(v["reference_id"].(string)),
			GroupNames:                  utils.ExpandStringSlice(v["policy_group_names"].(*pluginsdk.Set).List()),
		})
	}

	return &result, nil
}

func flattenAzureRMPolicySetDefinitionPolicyDefinitions(input *[]policy.DefinitionReference) ([]interface{}, error) {
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

		parameterValues, err := flattenParameterValuesValueToString(definition.Parameters)
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

func expandAzureRMPolicySetDefinitionPolicyGroups(input []interface{}) *[]policy.DefinitionGroup {
	result := make([]policy.DefinitionGroup, 0)

	for _, item := range input {
		v := item.(map[string]interface{})
		group := policy.DefinitionGroup{}
		if name := v["name"].(string); name != "" {
			group.Name = utils.String(name)
		}
		if displayName := v["display_name"].(string); displayName != "" {
			group.DisplayName = utils.String(displayName)
		}
		if category := v["category"].(string); category != "" {
			group.Category = utils.String(category)
		}
		if description := v["description"].(string); description != "" {
			group.Description = utils.String(description)
		}
		if metadataID := v["additional_metadata_resource_id"].(string); metadataID != "" {
			group.AdditionalMetadataID = utils.String(metadataID)
		}
		result = append(result, group)
	}

	return &result
}

func flattenAzureRMPolicySetDefinitionPolicyGroups(input *[]policy.DefinitionGroup) []interface{} {
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

func resourceARMPolicySetDefinitionPolicyDefinitionGroupHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(m["name"].(string))
	}

	return pluginsdk.HashString(buf.String())
}
