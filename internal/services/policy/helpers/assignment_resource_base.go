// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"
	"log"
	"maps"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-06-01/policyassignments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AssignmentBaseResource struct{}

type AssignmentNonComplianceMessageModel struct {
	Content                     string `tfschema:"content"`
	PolicyDefinitionReferenceId string `tfschema:"policy_definition_reference_id"`
}

type AssignmentOverrideModel struct {
	Value     string                            `tfschema:"value"`
	Selectors []AssignmentOverrideSelectorModel `tfschema:"selectors"`
}

type AssignmentOverrideSelectorModel struct {
	In    []string `tfschema:"in"`
	NotIn []string `tfschema:"not_in"`
	Kind  string   `tfschema:"kind"`
}

type AssignmentResourceSelectorModel struct {
	Name      string                            `tfschema:"name"`
	Selectors []AssignmentOverrideSelectorModel `tfschema:"selectors"`
}

type AssignmentBaseModel struct {
	Name                 string                                     `tfschema:"name"`
	PolicyDefinitionId   string                                     `tfschema:"policy_definition_id"`
	Description          string                                     `tfschema:"description"`
	DisplayName          string                                     `tfschema:"display_name"`
	Location             string                                     `tfschema:"location"`
	Enforce              bool                                       `tfschema:"enforce"`
	Metadata             string                                     `tfschema:"metadata"`
	Parameters           string                                     `tfschema:"parameters"`
	NotScopes            []string                                   `tfschema:"not_scopes"`
	Identity             []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	NonComplianceMessage []AssignmentNonComplianceMessageModel      `tfschema:"non_compliance_message"`
	Overrides            []AssignmentOverrideModel                  `tfschema:"overrides"`
	ResourceSelectors    []AssignmentResourceSelectorModel          `tfschema:"resource_selectors"`
}

func (br AssignmentBaseResource) CreateFunc(resourceName, scopeFieldName string) sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.AssignmentsClient

			var config AssignmentBaseModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := policyassignments.NewScopedPolicyAssignmentID(metadata.ResourceData.Get(scopeFieldName).(string), config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError(resourceName, id.ID())
				}
			}

			assignment := policyassignments.PolicyAssignment{
				Properties: &policyassignments.PolicyAssignmentProperties{
					PolicyDefinitionId: pointer.To(config.PolicyDefinitionId),
					DisplayName:        pointer.To(config.DisplayName),
					Scope:              pointer.To(id.Scope),
					EnforcementMode:    ConvertEnforcementMode(config.Enforce),
				},
			}

			if config.Description != "" {
				assignment.Properties.Description = pointer.To(config.Description)
			}

			if config.Location != "" {
				assignment.Location = pointer.To(location.Normalize(config.Location))
			}

			if len(config.Identity) > 0 {
				if assignment.Location == nil {
					return fmt.Errorf("`location` must be set when `identity` is assigned")
				}
				identityIns, err := identity.ExpandSystemOrUserAssignedMapFromModel(config.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				assignment.Identity = identityIns
			}

			if config.Parameters != "" {
				expandedParams, err := ExpandParameterValuesValueFromString(config.Parameters)
				if err != nil {
					return fmt.Errorf("expanding JSON for `parameters` %q: %+v", config.Parameters, err)
				}

				if expandedParams != nil {
					assignment.Properties.Parameters = &expandedParams
				}
			}

			if config.Metadata != "" {
				metaData, err := pluginsdk.ExpandJsonFromString(config.Metadata)
				if err != nil {
					return fmt.Errorf("unable to parse metadata: %s", err)
				}
				if metaData != nil {
					var d any = metaData
					assignment.Properties.Metadata = &d
				}
			}

			if len(config.NotScopes) > 0 {
				assignment.Properties.NotScopes = &config.NotScopes
			}

			if len(config.NonComplianceMessage) > 0 {
				assignment.Properties.NonComplianceMessages = br.ExpandNonComplianceMessages(config.NonComplianceMessage)
			}

			if len(config.Overrides) > 0 {
				assignment.Properties.Overrides = br.ExpandOverrides(config.Overrides)
			}

			if len(config.ResourceSelectors) > 0 {
				assignment.Properties.ResourceSelectors = br.ExpandResourceSelectors(config.ResourceSelectors)
			}

			if _, err := client.Create(ctx, id, assignment); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			// Policy Assignments are eventually consistent; wait for them to stabilize
			log.Printf("[DEBUG] Waiting for %s to become available..", id)
			if err := WaitForPolicyAssignmentToStabilize(ctx, client, id, true); err != nil {
				return fmt.Errorf("waiting for %s to become available: %s", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (br AssignmentBaseResource) DeleteFunc() sdk.ResourceFunc {
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
			if err := WaitForPolicyAssignmentToStabilize(ctx, client, *id, false); err != nil {
				return fmt.Errorf("waiting for the deletion of %s: %s", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (br AssignmentBaseResource) ReadFunc(scopeFieldName string) sdk.ResourceFunc {
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

			respModel := resp.Model
			state := AssignmentBaseModel{
				Name:     id.PolicyAssignmentName,
				Location: location.NormalizeNilable(respModel.Location),
			}

			// lintignore:R001
			metadata.ResourceData.Set(scopeFieldName, id.Scope)

			identityIns, err := identity.FlattenSystemOrUserAssignedMapToModel(respModel.Identity)
			if err != nil {
				return fmt.Errorf("FlattenSystemOrUserAssignedMapToModel: %+v", err)
			}
			if identityIns != nil {
				state.Identity = *identityIns
			}

			if props := respModel.Properties; props != nil {
				state.Description = pointer.From(props.Description)
				state.DisplayName = pointer.From(props.DisplayName)
				state.Enforce = pointer.From(props.EnforcementMode) == policyassignments.EnforcementModeDefault
				state.NotScopes = pointer.From(props.NotScopes)
				state.PolicyDefinitionId = pointer.From(props.PolicyDefinitionId)

				state.NonComplianceMessage = br.FlattenNonComplianceMessages(props.NonComplianceMessages)

				state.Metadata = FlattenJSON(pointer.From(props.Metadata))

				flattenedParameters, err := FlattenParameterValuesValueToStringV2(props.Parameters)
				if err != nil {
					return fmt.Errorf("serializing JSON from `parameters`: %+v", err)
				}
				state.Parameters = flattenedParameters

				state.Overrides = br.FlattenOverrides(props.Overrides)
				state.ResourceSelectors = br.FlattenResourceSelectors(props.ResourceSelectors)
			}

			return metadata.Encode(&state)
		},
		Timeout: 5 * time.Minute,
	}
}

func (br AssignmentBaseResource) UpdateFunc() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.AssignmentsClient

			id, err := policyassignments.ParseScopedPolicyAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config AssignmentBaseModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
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
				update.Properties.Description = pointer.To(config.Description)
			}
			if metadata.ResourceData.HasChange("display_name") {
				update.Properties.DisplayName = pointer.To(config.DisplayName)
			}
			if metadata.ResourceData.HasChange("enforce") {
				update.Properties.EnforcementMode = ConvertEnforcementMode(config.Enforce)
			}
			if metadata.ResourceData.HasChange("location") {
				update.Location = pointer.To(config.Location)
			}
			if metadata.ResourceData.HasChange("policy_definition_id") {
				update.Properties.PolicyDefinitionId = pointer.To(config.PolicyDefinitionId)
			}

			if metadata.ResourceData.HasChange("identity") {
				if update.Location == nil {
					return fmt.Errorf("`location` must be set when `identity` is assigned")
				}
				identityIns, err := identity.ExpandSystemOrUserAssignedMapFromModel(config.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				update.Identity = identityIns
			}

			if metadata.ResourceData.HasChange("metadata") {
				m := map[string]any{}
				if config.Metadata != "" {
					m, err = pluginsdk.ExpandJsonFromString(config.Metadata)
					if err != nil {
						return fmt.Errorf("parsing metadata: %+v", err)
					}
				}
				var i any = m
				update.Properties.Metadata = &i
			}

			if metadata.ResourceData.HasChange("not_scopes") {
				notScopes := config.NotScopes
				if notScopes == nil {
					notScopes = make([]string, 0)
				}
				update.Properties.NotScopes = &notScopes
			}

			if metadata.ResourceData.HasChange("non_compliance_message") {
				update.Properties.NonComplianceMessages = br.ExpandNonComplianceMessages(config.NonComplianceMessage)
			}

			if metadata.ResourceData.HasChange("parameters") {
				m := map[string]policyassignments.ParameterValuesValue{}

				if config.Parameters != "" {
					m, err = ExpandParameterValuesValueFromString(config.Parameters)
					if err != nil {
						return fmt.Errorf("expanding JSON for `parameters` %q: %+v", config.Parameters, err)
					}
				}
				update.Properties.Parameters = &m
			}

			if metadata.ResourceData.HasChange("overrides") {
				update.Properties.Overrides = br.ExpandOverrides(config.Overrides)
			}

			if metadata.ResourceData.HasChange("resource_selectors") {
				update.Properties.ResourceSelectors = br.ExpandResourceSelectors(config.ResourceSelectors)
			}

			// NOTE: there isn't an Update endpoint
			if _, err := client.Create(ctx, *id, update); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			// Policy Assignments are eventually consistent; wait for them to stabilize
			log.Printf("[DEBUG] Waiting for %s to become available..", id)
			if err := WaitForPolicyAssignmentToStabilize(ctx, client, *id, true); err != nil {
				return fmt.Errorf("waiting for %s to become available: %s", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (br AssignmentBaseResource) Arguments(fields map[string]*pluginsdk.Schema) map[string]*pluginsdk.Schema {
	// reuse the selector schema for Override and ResourceSelector block

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

		"metadata": MetadataSchema(),

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

		"overrides": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"selectors": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"in": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"kind": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  policyassignments.SelectorKindPolicyDefinitionReferenceId,
									ValidateFunc: validation.StringInSlice([]string{
										string(policyassignments.SelectorKindPolicyDefinitionReferenceId),
										string(policyassignments.SelectorKindResourceLocation),
									}, false),
								},

								"not_in": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					// more detail see https://learn.microsoft.com/azure/governance/policy/concepts/effects
					"value": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"resource_selectors": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"selectors": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"in": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"kind": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										// only 3 types supported for resourceSelector Kind
										// https://learn.microsoft.com/azure/governance/policy/concepts/assignment-structure#resource-selectors-preview
										string(policyassignments.SelectorKindResourceLocation),
										string(policyassignments.SelectorKindResourceType),
										string(policyassignments.SelectorKindResourceWithoutLocation),
									}, false),
								},

								"not_in": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	maps.Copy(output, fields)

	return output
}

func (br AssignmentBaseResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (br AssignmentBaseResource) FlattenNonComplianceMessages(input *[]policyassignments.NonComplianceMessage) []AssignmentNonComplianceMessageModel {
	if input == nil {
		return make([]AssignmentNonComplianceMessageModel, 0)
	}

	results := make([]AssignmentNonComplianceMessageModel, 0, len(*input))
	for _, v := range *input {
		results = append(results, AssignmentNonComplianceMessageModel{
			Content:                     v.Message,
			PolicyDefinitionReferenceId: pointer.From(v.PolicyDefinitionReferenceId),
		})
	}

	return results
}

func (br AssignmentBaseResource) ExpandNonComplianceMessages(input []AssignmentNonComplianceMessageModel) *[]policyassignments.NonComplianceMessage {
	if len(input) == 0 {
		return nil
	}

	output := make([]policyassignments.NonComplianceMessage, 0)
	for _, v := range input {
		ncm := policyassignments.NonComplianceMessage{
			Message: v.Content,
		}
		if id := v.PolicyDefinitionReferenceId; id != "" {
			ncm.PolicyDefinitionReferenceId = pointer.To(id)
		}
		output = append(output, ncm)
	}

	return &output
}

func (br AssignmentBaseResource) ExpandOverrides(overrides []AssignmentOverrideModel) *[]policyassignments.Override {
	if len(overrides) == 0 {
		return nil
	}

	var res []policyassignments.Override
	for _, v := range overrides {
		var item policyassignments.Override
		item.Value = pointer.To(v.Value)
		item.Kind = pointer.To(policyassignments.OverrideKindPolicyEffect)
		item.Selectors = br.ExpandSelectors(v.Selectors)
		res = append(res, item)
	}

	return &res
}

func (br AssignmentBaseResource) ExpandSelectors(i []AssignmentOverrideSelectorModel) *[]policyassignments.Selector {
	if len(i) == 0 {
		return nil
	}

	var res []policyassignments.Selector
	for _, v := range i {
		var item policyassignments.Selector
		item.Kind = pointer.ToEnum[policyassignments.SelectorKind](v.Kind)
		if len(v.In) > 0 {
			item.In = pointer.To(v.In)
		}
		if len(v.NotIn) > 0 {
			item.NotIn = pointer.To(v.NotIn)
		}
		res = append(res, item)
	}

	return &res
}

func (br AssignmentBaseResource) ExpandResourceSelectors(rs []AssignmentResourceSelectorModel) *[]policyassignments.ResourceSelector {
	if len(rs) == 0 {
		return nil
	}

	var res []policyassignments.ResourceSelector
	for _, v := range rs {
		var item policyassignments.ResourceSelector
		item.Name = pointer.To(v.Name)
		item.Selectors = br.ExpandSelectors(v.Selectors)
		res = append(res, item)
	}

	return &res
}

func (br AssignmentBaseResource) FlattenOverrides(overrides *[]policyassignments.Override) []AssignmentOverrideModel {
	if overrides == nil || len(*overrides) == 0 {
		return make([]AssignmentOverrideModel, 0)
	}

	res := make([]AssignmentOverrideModel, 0, len(*overrides))
	for _, o := range *overrides {
		item := AssignmentOverrideModel{
			Value:     pointer.From(o.Value),
			Selectors: br.FlattenSelectors(o.Selectors),
		}
		res = append(res, item)
	}

	return res
}

func (br AssignmentBaseResource) FlattenSelectors(selectors *[]policyassignments.Selector) []AssignmentOverrideSelectorModel {
	if selectors == nil || len(*selectors) == 0 {
		return make([]AssignmentOverrideSelectorModel, 0)
	}

	res := make([]AssignmentOverrideSelectorModel, 0, len(*selectors))
	for _, s := range *selectors {
		item := AssignmentOverrideSelectorModel{
			In:    pointer.From(s.In),
			NotIn: pointer.From(s.NotIn),
			Kind:  pointer.FromEnum(s.Kind),
		}
		res = append(res, item)
	}

	return res
}

func (br AssignmentBaseResource) FlattenResourceSelectors(selectors *[]policyassignments.ResourceSelector) []AssignmentResourceSelectorModel {
	if selectors == nil || *selectors == nil {
		return make([]AssignmentResourceSelectorModel, 0)
	}

	res := make([]AssignmentResourceSelectorModel, 0, len(*selectors))
	for _, v := range *selectors {
		item := AssignmentResourceSelectorModel{
			Name:      pointer.From(v.Name),
			Selectors: br.FlattenSelectors(v.Selectors),
		}
		res = append(res, item)
	}

	return res
}
