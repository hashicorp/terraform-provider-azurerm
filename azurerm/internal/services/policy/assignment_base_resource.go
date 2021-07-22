package policy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-09-01/policy"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/identity"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/sdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type policyAssignmentIdentity = identity.SystemAssigned

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
			// lintignore:R001
			metadata.ResourceData.Set(scopeFieldName, id.Scope)

			if err := metadata.ResourceData.Set("identity", br.flattenIdentity(resp.Identity)); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			if props := resp.AssignmentProperties; props != nil {
				metadata.ResourceData.Set("description", props.Description)
				metadata.ResourceData.Set("display_name", props.DisplayName)
				metadata.ResourceData.Set("enforce", props.EnforcementMode == policy.Default)
				metadata.ResourceData.Set("not_scopes", props.NotScopes)
				metadata.ResourceData.Set("policy_definition_id", props.PolicyDefinitionID)

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
					Type: existing.Identity.Type,
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

		"location": azure.SchemaLocationOptional(),

		"identity": policyAssignmentIdentity{}.Schema(),

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

		"parameters": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
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
	expanded, err := policyAssignmentIdentity{}.Expand(input)
	if err != nil {
		return nil, err
	}

	return &policy.Identity{
		Type: policy.ResourceIdentityType(expanded.Type),
	}, nil
}

func (br assignmentBaseResource) flattenIdentity(input *policy.Identity) []interface{} {
	var config *identity.ExpandedConfig
	if input != nil {
		config = &identity.ExpandedConfig{
			Type:        string(input.Type),
			PrincipalId: input.PrincipalID,
			TenantId:    input.TenantID,
		}
	}
	return policyAssignmentIdentity{}.Flatten(config)
}
