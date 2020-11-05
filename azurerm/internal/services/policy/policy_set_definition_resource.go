package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-09-01/policy"
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPolicySetDefinition() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPolicySetDefinitionCreate,
		Update: resourceArmPolicySetDefinitionUpdate,
		Read:   resourceArmPolicySetDefinitionRead,
		Delete: resourceArmPolicySetDefinitionDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.PolicySetDefinitionID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"policy_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(policy.BuiltIn),
					string(policy.Custom),
					string(policy.NotSpecified),
					string(policy.Static),
				}, false),
			},

			"management_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"management_group_name"},
				Deprecated:    "Deprecated in favour of `management_group_name`", // TODO -- remove this in next major version
			},

			"management_group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true, // TODO -- remove this when deprecation resolves
				ConflictsWith: []string{"management_group_id"},
			},

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"metadata": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: policySetDefinitionsMetadataDiffSuppressFunc,
			},

			"parameters": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			"policy_definitions": { // TODO -- remove in the next major version
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: policyDefinitionsDiffSuppressFunc,
				ExactlyOneOf:     []string{"policy_definitions", "policy_definition_reference"},
				Deprecated:       "Deprecated in favor of `policy_definition_reference`",
			},

			"policy_definition_reference": { // TODO -- rename this back to `policy_definition` after the deprecation
				Type:         schema.TypeList,
				Optional:     true,                                                          // TODO -- change this to Required after the deprecation
				Computed:     true,                                                          // TODO -- remove Computed after the deprecation
				ExactlyOneOf: []string{"policy_definitions", "policy_definition_reference"}, // TODO -- remove after the deprecation
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_definition_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.PolicyDefinitionID,
						},

						"parameters": { // TODO -- remove this attribute after the deprecation
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Deprecated: "Deprecated in favour of `parameter_values`",
						},

						"parameter_values": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true, // TODO -- remove Computed after the deprecation
							ValidateFunc:     validation.StringIsJSON,
							DiffSuppressFunc: structure.SuppressJsonDiff,
						},

						"reference_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func policySetDefinitionsMetadataDiffSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
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

// This function only serves the deprecated attribute `policy_definitions` in the old api-version.
// The old api-version only support two attribute - `policy_definition_id` and `parameters` in each element.
// Therefore this function is used for ignoring any other keys and then compare if there is a diff
func policyDefinitionsDiffSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	var oldPolicyDefinitions []DefinitionReferenceInOldApiVersion
	errOld := json.Unmarshal([]byte(old), &oldPolicyDefinitions)
	if errOld != nil {
		return false
	}

	var newPolicyDefinitions []DefinitionReferenceInOldApiVersion
	errNew := json.Unmarshal([]byte(new), &newPolicyDefinitions)
	if errNew != nil {
		return false
	}

	return reflect.DeepEqual(oldPolicyDefinitions, newPolicyDefinitions)
}

type DefinitionReferenceInOldApiVersion struct {
	// PolicyDefinitionID - The ID of the policy definition or policy set definition.
	PolicyDefinitionID *string `json:"policyDefinitionId,omitempty"`
	// Parameters - The parameter values for the referenced policy rule. The keys are the parameter names.
	Parameters map[string]*policy.ParameterValuesValue `json:"parameters"`
}

func resourceArmPolicySetDefinitionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.SetDefinitionsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	managementGroupName := ""
	if v, ok := d.GetOk("management_group_name"); ok {
		managementGroupName = v.(string)
	}
	if v, ok := d.GetOk("management_group_id"); ok {
		managementGroupName = v.(string)
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
		metaData, err := structure.ExpandJsonFromString(metaDataString)
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

	if v, ok := d.GetOk("policy_definitions"); ok {
		var policyDefinitions []policy.DefinitionReference
		err := json.Unmarshal([]byte(v.(string)), &policyDefinitions)
		if err != nil {
			return fmt.Errorf("expanding JSON for `policy_definitions`: %+v", err)
		}
		properties.PolicyDefinitions = &policyDefinitions
	}
	if v, ok := d.GetOk("policy_definition_reference"); ok {
		definitions, err := expandAzureRMPolicySetDefinitionPolicyDefinitions(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `policy_definition_reference`: %+v", err)
		}
		properties.PolicyDefinitions = definitions
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
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"404"},
		Target:                    []string{"200"},
		Refresh:                   policySetDefinitionRefreshFunc(ctx, client, name, managementGroupName),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
	}

	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(schema.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(schema.TimeoutUpdate)
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for Policy Set Definition %q to become available: %+v", name, err)
	}

	var resp policy.SetDefinition
	resp, err = getPolicySetDefinitionByName(ctx, client, name, managementGroupName)
	if err != nil {
		return fmt.Errorf("retrieving Policy Set Definition %q: %+v", name, err)
	}

	d.SetId(*resp.ID)

	return resourceArmPolicySetDefinitionRead(d, meta)
}

func resourceArmPolicySetDefinitionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.SetDefinitionsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicySetDefinitionID(d.Id())
	if err != nil {
		return err
	}

	managementGroupName := ""
	if scopeId, ok := id.PolicyScopeId.(parse.ScopeAtManagementGroup); ok {
		managementGroupName = scopeId.ManagementGroupName
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
			metaData, err := structure.ExpandJsonFromString(metaDataString)
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

	if d.HasChange("policy_definitions") {
		var policyDefinitions []policy.DefinitionReference
		err := json.Unmarshal([]byte(d.Get("policy_definitions").(string)), &policyDefinitions)
		if err != nil {
			return fmt.Errorf("expanding JSON for `policy_definitions`: %+v", err)
		}
		existing.SetDefinitionProperties.PolicyDefinitions = &policyDefinitions
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

	d.SetId(*resp.ID)

	return resourceArmPolicySetDefinitionRead(d, meta)
}

func resourceArmPolicySetDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.SetDefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicySetDefinitionID(d.Id())
	if err != nil {
		return err
	}

	managementGroupName := ""
	if scopeId, ok := id.PolicyScopeId.(parse.ScopeAtManagementGroup); ok {
		managementGroupName = scopeId.ManagementGroupName
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
	d.Set("management_group_name", managementGroupName)

	if props := resp.SetDefinitionProperties; props != nil {
		d.Set("policy_type", string(props.PolicyType))
		d.Set("display_name", props.DisplayName)
		d.Set("description", props.Description)

		if metadata := props.Metadata; metadata != nil {
			metadataVal := metadata.(map[string]interface{})
			metadataStr, err := structure.FlattenJsonToString(metadataVal)
			if err != nil {
				return fmt.Errorf("flattening JSON for `metadata`: %+v", err)
			}

			d.Set("metadata", metadataStr)
		}

		if parameters := props.Parameters; parameters != nil {
			parametersStr, err := flattenParameterDefintionsValueToString(parameters)
			if err != nil {
				return fmt.Errorf("flattening JSON for `parameters`: %+v", err)
			}

			d.Set("parameters", parametersStr)
		}

		if policyDefinitions := props.PolicyDefinitions; policyDefinitions != nil {
			policyDefinitionsRes, err := json.Marshal(policyDefinitions)
			if err != nil {
				return fmt.Errorf("flattening JSON for `policy_defintions`: %+v", err)
			}

			d.Set("policy_definitions", string(policyDefinitionsRes))
		}
		references, err := flattenAzureRMPolicySetDefinitionPolicyDefinitions(props.PolicyDefinitions)
		if err != nil {
			return fmt.Errorf("flattening `policy_definition_reference`: %+v", err)
		}
		if err := d.Set("policy_definition_reference", references); err != nil {
			return fmt.Errorf("setting `policy_definition_reference`: %+v", err)
		}
	}

	return nil
}

func resourceArmPolicySetDefinitionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.SetDefinitionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicySetDefinitionID(d.Id())
	if err != nil {
		return err
	}

	managementGroupName := ""
	switch scopeId := id.PolicyScopeId.(type) { // nolint gocritic
	case parse.ScopeAtManagementGroup:
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

func policySetDefinitionRefreshFunc(ctx context.Context, client *policy.SetDefinitionsClient, name string, managementGroupId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := getPolicySetDefinitionByName(ctx, client, name, managementGroupId)
		if err != nil {
			return nil, strconv.Itoa(res.StatusCode), fmt.Errorf("issuing read request in policySetDefinitionRefreshFunc for Policy Set Definition %q: %+v", name, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}

func expandAzureRMPolicySetDefinitionPolicyDefinitionsUpdate(d *schema.ResourceData) (*[]policy.DefinitionReference, error) {
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
		})
	}

	return &result, nil
}

func expandAzureRMPolicySetDefinitionPolicyDefinitions(input []interface{}) (*[]policy.DefinitionReference, error) {
	result := make([]policy.DefinitionReference, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		parameters := make(map[string]*policy.ParameterValuesValue)
		if p, ok := v["parameter_values"].(string); ok && p != "" {
			if err := json.Unmarshal([]byte(p), &parameters); err != nil {
				return nil, fmt.Errorf("unmarshalling `parameter_values`: %+v", err)
			}
		}
		if p, ok := v["parameters"].(map[string]interface{}); ok {
			if len(parameters) > 0 && len(p) > 0 {
				return nil, fmt.Errorf("cannot set both `parameters` and `parameter_values`")
			}
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
			"parameters":           parametersMap,
			"parameter_values":     parameterValues,
			"reference_id":         policyDefinitionReference,
		})
	}
	return result, nil
}
