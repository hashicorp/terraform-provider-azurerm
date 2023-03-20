package policy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-06-01/policyassignments"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type assignmentBaseResource struct{}

func (br assignmentBaseResource) createFunc(resourceName, scopeFieldName string) sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.AssignmentsClient
			id := policyassignments.NewScopedPolicyAssignmentID(metadata.ResourceData.Get(scopeFieldName).(string), metadata.ResourceData.Get("name").(string))
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(resourceName, id.ID())
			}

			assignment := policyassignments.PolicyAssignment{
				Properties: &policyassignments.PolicyAssignmentProperties{
					PolicyDefinitionId: utils.String(metadata.ResourceData.Get("policy_definition_id").(string)),
					DisplayName:        utils.String(metadata.ResourceData.Get("display_name").(string)),
					Scope:              utils.String(id.Scope),
					EnforcementMode:    convertEnforcementMode(metadata.ResourceData.Get("enforce").(bool)),
				},
			}

			if v := metadata.ResourceData.Get("description").(string); v != "" {
				assignment.Properties.Description = utils.String(v)
			}

			if v := metadata.ResourceData.Get("location").(string); v != "" {
				assignment.Location = utils.String(azure.NormalizeLocation(v))
			}

			if v, ok := metadata.ResourceData.GetOk("identity"); ok {
				if assignment.Location == nil {
					return fmt.Errorf("`location` must be set when `identity` is assigned")
				}
				identityIns, err := identity.ExpandSystemOrUserAssignedMap(v.([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				assignment.Identity = identityIns
			}

			if v := metadata.ResourceData.Get("parameters").(string); v != "" {
				expandedParams, err := expandParameterValuesValueFromString(v)
				if err != nil {
					return fmt.Errorf("expanding JSON for `parameters` %q: %+v", v, err)
				}

				if expandedParams != nil {
					assignment.Properties.Parameters = &expandedParams
				}
			}

			if metaDataString := metadata.ResourceData.Get("metadata").(string); metaDataString != "" {
				metaData, err := pluginsdk.ExpandJsonFromString(metaDataString)
				if err != nil {
					return fmt.Errorf("unable to parse metadata: %s", err)
				}
				if metaData != nil {
					var d interface{} = metaData
					assignment.Properties.Metadata = &d
				}
			}

			if v, ok := metadata.ResourceData.GetOk("not_scopes"); ok {
				assignment.Properties.NotScopes = expandAzureRmPolicyNotScopes(v.([]interface{}))
			}

			if msgs := metadata.ResourceData.Get("non_compliance_message").([]interface{}); len(msgs) > 0 {
				assignment.Properties.NonComplianceMessages = br.expandNonComplianceMessages(msgs)
			}

			if _, err := client.Create(ctx, id, assignment); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// Policy Assignments are eventually consistent; wait for them to stabilize
			log.Printf("[DEBUG] Waiting for %s to become available..", id)
			if err := waitForPolicyAssignmentToStabilize(ctx, client, id, true); err != nil {
				return fmt.Errorf("waiting for %s to become available: %s", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (br assignmentBaseResource) deleteFunc() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.AssignmentsClient

			id, err := policyassignments.ParseScopedPolicyAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting Policy Assignment %q: %+v", id, err)
			}

			// Policy Assignments are eventually consistent; wait for it to be gone
			log.Printf("[DEBUG] Waiting for %s to disappear..", id)
			if err := waitForPolicyAssignmentToStabilize(ctx, client, *id, false); err != nil {
				return fmt.Errorf("waiting for the deletion of %s: %s", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (br assignmentBaseResource) readFunc(scopeFieldName string) sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.AssignmentsClient

			id, err := policyassignments.ParseScopedPolicyAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("reading %s: %+v", *id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("reading nil model")
			}

			model := resp.Model
			metadata.ResourceData.Set("name", id.PolicyAssignmentName)
			metadata.ResourceData.Set("location", location.NormalizeNilable(model.Location))
			// lintignore:R001
			metadata.ResourceData.Set(scopeFieldName, id.Scope)

			identityIns, err := identity.FlattenSystemOrUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("FlattenSystemOrUserAssignedMap: %+v", err)
			}
			if err = metadata.ResourceData.Set("identity", identityIns); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			if props := model.Properties; props != nil {
				metadata.ResourceData.Set("description", props.Description)
				metadata.ResourceData.Set("display_name", props.DisplayName)
				var enforce bool
				if mode := props.EnforcementMode; mode != nil {
					enforce = (*props.EnforcementMode) == policyassignments.EnforcementModeDefault
				}
				metadata.ResourceData.Set("enforce", enforce)
				metadata.ResourceData.Set("not_scopes", props.NotScopes)
				metadata.ResourceData.Set("policy_definition_id", props.PolicyDefinitionId)

				metadata.ResourceData.Set("non_compliance_message", br.flattenNonComplianceMessages(props.NonComplianceMessages))

				flattenedMetaData := flattenJSON(pointer.From(props.Metadata))
				metadata.ResourceData.Set("metadata", flattenedMetaData)

				flattenedParameters, err := flattenParameterValuesValueToStringV2(props.Parameters)
				if err != nil {
					return fmt.Errorf("serializing JSON from `parameters`: %+v", err)
				}
				metadata.ResourceData.Set("parameters", flattenedParameters)
			}

			return nil
		},
		Timeout: 5 * time.Minute,
	}
}

func (br assignmentBaseResource) updateFunc() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.AssignmentsClient

			id, err := policyassignments.ParseScopedPolicyAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			getResp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			existing := getResp.Model
			if existing == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			update := policyassignments.PolicyAssignment{
				Location:   existing.Location,
				Properties: existing.Properties,
			}
			if existing.Identity != nil {
				update.Identity = existing.Identity
			}

			if metadata.ResourceData.HasChange("description") {
				update.Properties.Description = utils.String(metadata.ResourceData.Get("description").(string))
			}
			if metadata.ResourceData.HasChange("display_name") {
				update.Properties.DisplayName = utils.String(metadata.ResourceData.Get("display_name").(string))
			}
			if metadata.ResourceData.HasChange("enforce") {
				update.Properties.EnforcementMode = convertEnforcementMode(metadata.ResourceData.Get("enforce").(bool))
			}
			if metadata.ResourceData.HasChange("location") {
				update.Location = utils.String(metadata.ResourceData.Get("location").(string))
			}
			if metadata.ResourceData.HasChange("policy_definition_id") {
				update.Properties.PolicyDefinitionId = utils.String(metadata.ResourceData.Get("policy_definition_id").(string))
			}

			if metadata.ResourceData.HasChange("identity") {
				if update.Location == nil {
					return fmt.Errorf("`location` must be set when `identity` is assigned")
				}
				identityRaw := metadata.ResourceData.Get("identity").([]interface{})
				identityIns, err := identity.ExpandSystemOrUserAssignedMap(identityRaw)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				update.Identity = identityIns
			}

			if metadata.ResourceData.HasChange("metadata") {
				v := metadata.ResourceData.Get("metadata").(string)
				m := map[string]interface{}{}
				if v != "" {
					m, err = pluginsdk.ExpandJsonFromString(v)
					if err != nil {
						return fmt.Errorf("parsing metadata: %+v", err)
					}
				}
				var i interface{} = m
				update.Properties.Metadata = &i
			}

			if metadata.ResourceData.HasChange("not_scopes") {
				update.Properties.NotScopes = expandAzureRmPolicyNotScopes(metadata.ResourceData.Get("not_scopes").([]interface{}))
			}

			if metadata.ResourceData.HasChange("non_compliance_message") {
				update.Properties.NonComplianceMessages = br.expandNonComplianceMessages(metadata.ResourceData.Get("non_compliance_message").([]interface{}))
			}

			if metadata.ResourceData.HasChange("parameters") {
				m := map[string]policyassignments.ParameterValuesValue{}

				if v := metadata.ResourceData.Get("parameters").(string); v != "" {
					m, err = expandParameterValuesValueFromString(v)
					if err != nil {
						return fmt.Errorf("expanding JSON for `parameters` %q: %+v", v, err)
					}
				}
				update.Properties.Parameters = &m
			}

			// NOTE: there isn't an Update endpoint
			if _, err := client.Create(ctx, *id, update); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			// Policy Assignments are eventually consistent; wait for them to stabilize
			log.Printf("[DEBUG] Waiting for %s to become available..", id)
			if err := waitForPolicyAssignmentToStabilize(ctx, client, *id, true); err != nil {
				return fmt.Errorf("waiting for %s to become available: %s", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (br assignmentBaseResource) arguments(fields map[string]*pluginsdk.Schema) map[string]*pluginsdk.Schema {
	output := map[string]*pluginsdk.Schema{
		// NOTE: `name` isn't included since it varies depending on the resource, so it's expected to be passed in
		"policy_definition_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				validate.PolicyDefinitionID,
				validate.PolicySetDefinitionID,
			),
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"location": commonschema.LocationOptional(),

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"enforce": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"metadata": metadataSchema(),

		"not_scopes": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"non_compliance_message": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"content": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"policy_definition_reference_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"parameters": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},
	}

	for k, v := range fields {
		output[k] = v
	}

	return output
}

func (br assignmentBaseResource) attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (br assignmentBaseResource) flattenNonComplianceMessages(input *[]policyassignments.NonComplianceMessage) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		results = append(results, map[string]interface{}{
			"content":                        v.Message,
			"policy_definition_reference_id": pointer.From(v.PolicyDefinitionReferenceId),
		})
	}

	return results
}

func (br assignmentBaseResource) expandNonComplianceMessages(input []interface{}) *[]policyassignments.NonComplianceMessage {
	if len(input) == 0 {
		return nil
	}

	output := make([]policyassignments.NonComplianceMessage, 0)
	for _, v := range input {
		if m, ok := v.(map[string]interface{}); ok {
			ncm := policyassignments.NonComplianceMessage{
				Message: m["content"].(string),
			}
			if id := m["policy_definition_reference_id"].(string); id != "" {
				ncm.PolicyDefinitionReferenceId = utils.String(id)
			}
			output = append(output, ncm)
		}
	}

	return &output
}
func expandAzureRmPolicyNotScopes(input []interface{}) *[]string {
	notScopesRes := make([]string, 0)

	for _, notScope := range input {
		s, ok := notScope.(string)
		if ok {
			notScopesRes = append(notScopesRes, s)
		}
	}

	return &notScopesRes
}
