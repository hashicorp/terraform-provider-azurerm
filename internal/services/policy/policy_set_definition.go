package policy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2025-01-01/policysetdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PolicyDefinitionReferenceModel struct {
	PolicyDefinitionID string   `tfschema:"policy_definition_id"`
	ParameterValues    string   `tfschema:"parameter_values"`
	ReferenceID        string   `tfschema:"reference_id"`
	PolicyGroupNames   []string `tfschema:"policy_group_names"`
}

type PolicyDefinitionGroupModel struct {
	Name                         string `tfschema:"name"`
	DisplayName                  string `tfschema:"display_name"`
	Category                     string `tfschema:"category"`
	Description                  string `tfschema:"description"`
	AdditionalMetadataResourceID string `tfschema:"additional_metadata_resource_id"`
}

func policyDefinitionReferenceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
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
	}
}

func policyDefinitionGroupSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
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
		Set: policySetDefinitionPolicyDefinitionGroupHash,
	}
}

func policySetDefinitionPolicyDefinitionGroupHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(m["name"].(string))
	}

	return pluginsdk.HashString(buf.String())
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

func expandPolicyDefinitionReference(input []PolicyDefinitionReferenceModel) ([]policysetdefinitions.PolicyDefinitionReference, error) {
	result := make([]policysetdefinitions.PolicyDefinitionReference, 0)

	for _, v := range input {
		expandedParameters, err := expandPolicyDefinitionReferenceParameterValues(v.ParameterValues)
		if err != nil {
			return nil, fmt.Errorf("expanding `parameter_values`: %+v", err)
		}

		result = append(result, policysetdefinitions.PolicyDefinitionReference{
			GroupNames:                  pointer.To(v.PolicyGroupNames),
			Parameters:                  expandedParameters,
			PolicyDefinitionId:          v.PolicyDefinitionID,
			PolicyDefinitionReferenceId: pointer.To(v.ReferenceID),
		})
	}

	return result, nil
}

func flattenPolicyDefinitionReference(input []policysetdefinitions.PolicyDefinitionReference) ([]PolicyDefinitionReferenceModel, error) {
	result := make([]PolicyDefinitionReferenceModel, 0)

	for _, definition := range input {
		parameterValues, err := flattenPolicyDefinitionReferenceParameterValues(definition.Parameters)
		if err != nil {
			return nil, fmt.Errorf("flattening `parameter_values`: %+v", err)
		}

		result = append(result, PolicyDefinitionReferenceModel{
			PolicyDefinitionID: definition.PolicyDefinitionId,
			ParameterValues:    parameterValues,
			ReferenceID:        pointer.From(definition.PolicyDefinitionReferenceId),
			PolicyGroupNames:   pointer.From(definition.GroupNames),
		})
	}

	return result, nil
}

func expandPolicyDefinitionReferenceParameterValues(input string) (*map[string]policysetdefinitions.ParameterValuesValue, error) {
	if input == "" {
		return nil, nil
	}

	result := make(map[string]policysetdefinitions.ParameterValuesValue)
	if err := json.Unmarshal([]byte(input), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func flattenPolicyDefinitionReferenceParameterValues(input *map[string]policysetdefinitions.ParameterValuesValue) (string, error) {
	if input == nil {
		return "", nil
	}

	result, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	return string(result), err
}

func expandPolicyDefinitionGroup(input []PolicyDefinitionGroupModel) *[]policysetdefinitions.PolicyDefinitionGroup {
	result := make([]policysetdefinitions.PolicyDefinitionGroup, 0)

	for _, v := range input {
		group := policysetdefinitions.PolicyDefinitionGroup{
			Category:    pointer.To(v.Category),
			Description: pointer.To(v.Description),
			DisplayName: pointer.To(v.DisplayName),
			Name:        v.Name,
		}

		// The API returns an error if we send an empty string for `AdditionalMetadataResourceID`
		if v.AdditionalMetadataResourceID != "" {
			group.AdditionalMetadataId = pointer.To(v.AdditionalMetadataResourceID)
		}

		result = append(result, group)
	}

	return &result
}

func flattenPolicyDefinitionGroup(input *[]policysetdefinitions.PolicyDefinitionGroup) []PolicyDefinitionGroupModel {
	result := make([]PolicyDefinitionGroupModel, 0)

	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, PolicyDefinitionGroupModel{
			AdditionalMetadataResourceID: pointer.From(v.AdditionalMetadataId),
			Category:                     pointer.From(v.Category),
			Description:                  pointer.From(v.Description),
			DisplayName:                  pointer.From(v.DisplayName),
			Name:                         v.Name,
		})
	}

	return result
}

func flattenParameterDefinitionsValue(input *map[string]policysetdefinitions.ParameterDefinitionsValue) (string, error) {
	if input == nil {
		return "", nil
	}

	result, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	compactJson := bytes.Buffer{}
	if err := json.Compact(&compactJson, result); err != nil {
		return "", err
	}

	return compactJson.String(), nil
}

func expandParameterDefinitionsValue(input string) (*map[string]policysetdefinitions.ParameterDefinitionsValue, error) {
	var result map[string]policysetdefinitions.ParameterDefinitionsValue

	err := json.Unmarshal([]byte(input), &result)

	return &result, err
}

func getPolicySetDefinitionByID(ctx context.Context, client *policysetdefinitions.PolicySetDefinitionsClient, id any) (*http.Response, *policysetdefinitions.PolicySetDefinition, error) {
	// TODO: Remove post 5.0
	switch id := id.(type) {
	case policysetdefinitions.ProviderPolicySetDefinitionId:
		return getPolicySetDefinition(ctx, client, id)
	case policysetdefinitions.Providers2PolicySetDefinitionId:
		resp, err := client.GetAtManagementGroup(ctx, id, policysetdefinitions.DefaultGetAtManagementGroupOperationOptions())
		return resp.HttpResponse, resp.Model, err
	default:
		return nil, nil, fmt.Errorf("`id` was not one of the expected types: %T", id)
	}
}

func policySetDefinitionRefreshFunc(ctx context.Context, client *policysetdefinitions.PolicySetDefinitionsClient, id any) pluginsdk.StateRefreshFunc {
	// TODO: Remove post 5.0
	return func() (interface{}, string, error) {
		resp, _, err := getPolicySetDefinitionByID(ctx, client, id)
		if err != nil && !response.WasNotFound(resp) {
			return nil, strconv.Itoa(resp.StatusCode), fmt.Errorf("retrieving %s: %+v", id, err)
		}

		return resp, strconv.Itoa(resp.StatusCode), nil
	}
}
