package policy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
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
			id := parse.NewPolicyAssignmentId(metadata.ResourceData.Get(scopeFieldName).(string), metadata.ResourceData.Get("name").(string))
			existing, err := client.Get(ctx, id.Scope, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return tf.ImportAsExistsError(resourceName, id.ID())
			}

			assignment := policy.Assignment{
				AssignmentProperties: &policy.AssignmentProperties{
					PolicyDefinitionID: utils.String(metadata.ResourceData.Get("policy_definition_id").(string)),
					DisplayName:        utils.String(metadata.ResourceData.Get("display_name").(string)),
					Scope:              utils.String(id.Scope),
					EnforcementMode:    convertEnforcementMode(metadata.ResourceData.Get("enforce").(bool)),
				},
			}

			if v := metadata.ResourceData.Get("description").(string); v != "" {
				assignment.AssignmentProperties.Description = utils.String(v)
			}

			if v := metadata.ResourceData.Get("location").(string); v != "" {
				assignment.Location = utils.String(azure.NormalizeLocation(v))
			}

			if v, ok := metadata.ResourceData.GetOk("identity"); ok {
				if assignment.Location == nil {
					return fmt.Errorf("`location` must be set when `identity` is assigned")
				}
				identity, err := br.expandIdentity(v.([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				assignment.Identity = identity
			}

			if v := metadata.ResourceData.Get("parameters").(string); v != "" {
				expandedParams, err := expandParameterValuesValueFromString(v)
				if err != nil {
					return fmt.Errorf("expanding JSON for `parameters` %q: %+v", v, err)
				}

				assignment.AssignmentProperties.Parameters = expandedParams
			}

			if metaDataString := metadata.ResourceData.Get("metadata").(string); metaDataString != "" {
				metaData, err := pluginsdk.ExpandJsonFromString(metaDataString)
				if err != nil {
					return fmt.Errorf("unable to parse metadata: %s", err)
				}
				assignment.AssignmentProperties.Metadata = &metaData
			}

			if v, ok := metadata.ResourceData.GetOk("not_scopes"); ok {
				assignment.AssignmentProperties.NotScopes = expandAzureRmPolicyNotScopes(v.([]interface{}))
			}

			if msgs := metadata.ResourceData.Get("non_compliance_message").([]interface{}); len(msgs) > 0 {
				assignment.NonComplianceMessages = br.expandNonComplianceMessages(msgs)
			}

			if _, err := client.Create(ctx, id.Scope, id.Name, assignment); err != nil {
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

			id, err := parse.PolicyAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, id.Scope, id.Name); err != nil {
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

			id, err := parse.PolicyAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.Scope, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			metadata.ResourceData.Set("name", id.Name)
			metadata.ResourceData.Set("location", location.NormalizeNilable(resp.Location))
			//lintignore:R001
			metadata.ResourceData.Set(scopeFieldName, id.Scope)

			identity, _ := br.flattenIdentity(resp.Identity)
			if err := metadata.ResourceData.Set("identity", identity); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			if props := resp.AssignmentProperties; props != nil {
				metadata.ResourceData.Set("description", props.Description)
				metadata.ResourceData.Set("display_name", props.DisplayName)
				metadata.ResourceData.Set("enforce", props.EnforcementMode == policy.EnforcementModeDefault)
				metadata.ResourceData.Set("not_scopes", props.NotScopes)
				metadata.ResourceData.Set("policy_definition_id", props.PolicyDefinitionID)

				metadata.ResourceData.Set("non_compliance_message", br.flattenNonComplianceMessages(props.NonComplianceMessages))

				flattenedMetaData := flattenJSON(props.Metadata)
				metadata.ResourceData.Set("metadata", flattenedMetaData)

				flattenedParameters, err := flattenParameterValuesValueToString(props.Parameters)
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

			id, err := parse.PolicyAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.Scope, id.Name)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.AssignmentProperties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			update := policy.Assignment{
				Location:             existing.Location,
				AssignmentProperties: existing.AssignmentProperties,
			}
			if existing.Identity != nil {
				update.Identity = &policy.Identity{
					Type:                   existing.Identity.Type,
					UserAssignedIdentities: existing.Identity.UserAssignedIdentities,
				}
			}

			if metadata.ResourceData.HasChange("description") {
				update.AssignmentProperties.Description = utils.String(metadata.ResourceData.Get("description").(string))
			}
			if metadata.ResourceData.HasChange("display_name") {
				update.AssignmentProperties.DisplayName = utils.String(metadata.ResourceData.Get("display_name").(string))
			}
			if metadata.ResourceData.HasChange("enforce") {
				update.AssignmentProperties.EnforcementMode = convertEnforcementMode(metadata.ResourceData.Get("enforce").(bool))
			}
			if metadata.ResourceData.HasChange("location") {
				update.Location = utils.String(metadata.ResourceData.Get("location").(string))
			}
			if metadata.ResourceData.HasChange("policy_definition_id") {
				update.AssignmentProperties.PolicyDefinitionID = utils.String(metadata.ResourceData.Get("policy_definition_id").(string))
			}

			if metadata.ResourceData.HasChange("identity") {
				if update.Location == nil {
					return fmt.Errorf("`location` must be set when `identity` is assigned")
				}
				identityRaw := metadata.ResourceData.Get("identity").([]interface{})
				identity, err := br.expandIdentity(identityRaw)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				update.Identity = identity
			}

			if metadata.ResourceData.HasChange("metadata") {
				v := metadata.ResourceData.Get("metadata").(string)
				update.AssignmentProperties.Metadata = map[string]interface{}{}
				if v != "" {
					metaData, err := pluginsdk.ExpandJsonFromString(v)
					if err != nil {
						return fmt.Errorf("parsing metadata: %+v", err)
					}
					update.AssignmentProperties.Metadata = &metaData
				}
			}

			if metadata.ResourceData.HasChange("not_scopes") {
				update.AssignmentProperties.NotScopes = expandAzureRmPolicyNotScopes(metadata.ResourceData.Get("not_scopes").([]interface{}))
			}

			if metadata.ResourceData.HasChange("non_compliance_message") {
				update.AssignmentProperties.NonComplianceMessages = br.expandNonComplianceMessages(metadata.ResourceData.Get("non_compliance_message").([]interface{}))
			}

			if metadata.ResourceData.HasChange("parameters") {
				update.AssignmentProperties.Parameters = map[string]*policy.ParameterValuesValue{}

				if v := metadata.ResourceData.Get("parameters").(string); v != "" {
					expandedParams, err := expandParameterValuesValueFromString(v)
					if err != nil {
						return fmt.Errorf("expanding JSON for `parameters` %q: %+v", v, err)
					}
					update.AssignmentProperties.Parameters = expandedParams
				}
			}

			// NOTE: there isn't an Update endpoint
			if _, err := client.Create(ctx, id.Scope, id.Name, update); err != nil {
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

func (br assignmentBaseResource) expandIdentity(input []interface{}) (*policy.Identity, error) {
	expanded, err := identity.ExpandSystemOrUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := policy.Identity{
		Type: policy.ResourceIdentityType(string(expanded.Type)),
	}
	if expanded.Type == identity.TypeUserAssigned {
		out.UserAssignedIdentities = make(map[string]*policy.IdentityUserAssignedIdentitiesValue)
		for k := range expanded.IdentityIds {
			out.UserAssignedIdentities[k] = &policy.IdentityUserAssignedIdentitiesValue{
				// intentionally empty
			}
		}
	}
	return &out, nil
}

func (br assignmentBaseResource) flattenIdentity(input *policy.Identity) (*[]interface{}, error) {
	var config *identity.SystemOrUserAssignedMap
	if input != nil {
		config = &identity.SystemOrUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}

		if input.PrincipalID != nil {
			config.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			config.TenantId = *input.TenantID
		}
		for k, v := range input.UserAssignedIdentities {
			config.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientID,
				PrincipalId: v.PrincipalID,
			}
		}
	}
	return identity.FlattenSystemOrUserAssignedMap(config)
}

func (br assignmentBaseResource) flattenNonComplianceMessages(input *[]policy.NonComplianceMessage) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		content := ""
		if v.Message != nil {
			content = *v.Message
		}
		policyDefinitionReferenceId := ""
		if v.PolicyDefinitionReferenceID != nil {
			policyDefinitionReferenceId = *v.PolicyDefinitionReferenceID
		}
		results = append(results, map[string]interface{}{
			"content":                        content,
			"policy_definition_reference_id": policyDefinitionReferenceId,
		})
	}

	return results
}

func (br assignmentBaseResource) expandNonComplianceMessages(input []interface{}) *[]policy.NonComplianceMessage {
	if len(input) == 0 {
		return nil
	}

	output := make([]policy.NonComplianceMessage, 0)
	for _, v := range input {
		if m, ok := v.(map[string]interface{}); ok {
			ncm := policy.NonComplianceMessage{
				Message: utils.String(m["content"].(string)),
			}
			if id := m["policy_definition_reference_id"].(string); id != "" {
				ncm.PolicyDefinitionReferenceID = utils.String(id)
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
