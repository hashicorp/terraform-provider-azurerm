// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/threatintelligence"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type IndicatorPatternType string

const (
	PatternTypeDomainName IndicatorPatternType = "domain-name"
	PatternTypeFile       IndicatorPatternType = "file"
	PatternTypeIpV4Addr   IndicatorPatternType = "ipv4-addr"
	PatternTypeIpV6Addr   IndicatorPatternType = "ipv6-addr"
	PatternTypeUrl        IndicatorPatternType = "url"
)

const killChainName = "Lockheed Martin - Intrusion Kill Chain"

type IndicatorModel struct {
	Name                       string                   `tfschema:"guid"`
	WorkspaceId                string                   `tfschema:"workspace_id"`
	Confidence                 int64                    `tfschema:"confidence"`
	CreatedByRef               string                   `tfschema:"created_by"`
	Description                string                   `tfschema:"description"`
	DisplayName                string                   `tfschema:"display_name"`
	Extensions                 string                   `tfschema:"extension"`
	ExternalReference          []externalReferenceModel `tfschema:"external_reference"`
	GranularMarkings           []granularMarkingModel   `tfschema:"granular_marking"`
	IndicatorTypes             []string                 `tfschema:"indicator_type"`
	KillChainPhases            []killChainPhaseModel    `tfschema:"kill_chain_phase"`
	Labels                     []string                 `tfschema:"tags"`
	Language                   string                   `tfschema:"language"`
	ObjectMarking              []string                 `tfschema:"object_marking_refs"`
	ParsedPattern              []parsedPatternModel     `tfschema:"parsed_pattern"`
	Pattern                    string                   `tfschema:"pattern"`
	PatternType                string                   `tfschema:"pattern_type"`
	PatternVersion             string                   `tfschema:"pattern_version"`
	Revoked                    bool                     `tfschema:"revoked"`
	Source                     string                   `tfschema:"source"`
	ThreatTypes                []string                 `tfschema:"threat_types"`
	ValidFrom                  string                   `tfschema:"validate_from_utc"`
	ValidUntil                 string                   `tfschema:"validate_until_utc"`
	CreatedOn                  string                   `tfschema:"created_on"`
	ExternalId                 string                   `tfschema:"external_id"`
	ExternalLastUpdatedTimeUtc string                   `tfschema:"external_last_updated_time_utc"`
	LastUpdatedTimeUtc         string                   `tfschema:"last_updated_time_utc"`
	Defanged                   bool                     `tfschema:"defanged"`
}

type externalReferenceModel struct {
	SourceName  string            `tfschema:"source_name"`
	Url         string            `tfschema:"url"`
	Hashes      map[string]string `tfschema:"hashes"`
	Description string            `tfschema:"description"`
	ExternalId  string            `tfschema:"id"`
}

type granularMarkingModel struct {
	MarkingRef string   `tfschema:"marking_ref"`
	Selectors  []string `tfschema:"selectors"`
	Language   string   `tfschema:"language"`
}
type killChainPhaseModel struct {
	PhaseName string `tfschema:"name"`
}
type parsedPatternModel struct {
	PatternTypeValues []patternTypeValuesModel `tfschema:"pattern_type_values"`
	PatternTypeKey    string                   `tfschema:"pattern_type_key"`
}

type patternTypeValuesModel struct {
	Value     string `tfschema:"value"`
	ValueType string `tfschema:"value_type"`
}

type ThreatIntelligenceIndicator struct{}

var _ sdk.ResourceWithUpdate = ThreatIntelligenceIndicator{}

func (r ThreatIntelligenceIndicator) ResourceType() string {
	return "azurerm_sentinel_threat_intelligence_indicator"
}

func (r ThreatIntelligenceIndicator) ModelObject() interface{} {
	return &IndicatorModel{}
}

func (r ThreatIntelligenceIndicator) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return threatintelligence.ValidateIndicatorID
}

func (r ThreatIntelligenceIndicator) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"confidence": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 100),
			Default:      -1, // set the default value to -1 to split `nil` and `0`.
		},

		"created_by": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"extension": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C API sets this if omitted without issues for overwriting/reverting to default so this can remain
			Computed:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"external_reference": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"description": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"hashes": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"source_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"granular_marking": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"language": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"marking_ref": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"selectors": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"kill_chain_phase": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"tags": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"language": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"object_marking_refs": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"pattern": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"pattern_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(PatternTypeDomainName),
				string(PatternTypeFile),
				string(PatternTypeIpV4Addr),
				string(PatternTypeIpV6Addr),
				string(PatternTypeUrl),
			}, false),
		},

		"pattern_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"revoked": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"source": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"threat_types": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"validate_from_utc": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"validate_until_utc": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because the API returns a date around 1 year in the future if omitted from config
			Computed:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},
	}
}

func (r ThreatIntelligenceIndicator) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"guid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"created_on": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"defanged": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"external_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"external_last_updated_time_utc": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"last_updated_time_utc": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"parsed_pattern": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"pattern_type_key": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"pattern_type_values": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"value": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"value_type": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},

		"indicator_type": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r ThreatIntelligenceIndicator) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model IndicatorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Sentinel.ThreatIntelligenceClient

			workspaceId, err := threatintelligence.ParseWorkspaceID(model.WorkspaceId)
			if err != nil {
				return fmt.Errorf("parsing Workspace id %s: %+v", model.WorkspaceId, err)
			}

			patternValue, err := expandIndicatorPattern(model.PatternType, model.Pattern)
			if err != nil {
				return err
			}

			// it could not get the indicator by name before it has been created, because the name is generated by service side.
			// but we can not create duplicated indicator with same values, so list the values and find the existing one.
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				items, err := client.IndicatorsListComplete(ctx, *workspaceId, threatintelligence.DefaultIndicatorsListOperationOptions())
				if err != nil {
					return fmt.Errorf("listing indicators on %s: %+v", workspaceId, err)
				}

				for _, indicator := range items.Items {
					v, ok := indicator.(threatintelligence.ThreatIntelligenceIndicatorModel)
					if !ok {
						continue
					}

					if props := v.Properties; props != nil && pointer.From(props.PatternType) == model.PatternType {
						if pointer.From(props.Pattern) == patternValue {
							if indicatorId := pointer.From(v.Id); indicatorId != "" {
								return tf.ImportAsExistsError("azurerm_sentinel_threat_intelligence_indicator", indicatorId)
							}
							return fmt.Errorf("checking for presence of existing indicator: `id` was nil")
						}
					}
				}
			}

			payload := threatintelligence.ThreatIntelligenceIndicatorModel{
				Kind: threatintelligence.ThreatIntelligenceResourceKindEnumIndicator,
				Properties: &threatintelligence.ThreatIntelligenceIndicatorProperties{
					PatternType:    &model.PatternType,
					Pattern:        &patternValue,
					Revoked:        &model.Revoked,
					CreatedByRef:   pointer.ToOrNil(model.CreatedByRef),
					Description:    pointer.ToOrNil(model.Description),
					DisplayName:    pointer.ToOrNil(model.DisplayName),
					Language:       pointer.ToOrNil(model.Language),
					PatternVersion: pointer.ToOrNil(model.PatternVersion),
					Source:         pointer.ToOrNil(model.Source),
				},
			}
			props := payload.Properties

			if model.Confidence != -1 {
				props.Confidence = pointer.To(model.Confidence)
			}

			if model.Extensions != "" {
				extensionsValue, err := pluginsdk.ExpandJsonFromString(model.Extensions)
				if err != nil {
					return err
				}
				props.Extensions = &extensionsValue
			}

			props.ExternalReferences = expandThreatIntelligenceExternalReferenceModel(model.ExternalReference)

			props.GranularMarkings = expandThreatIntelligenceGranularMarkingModelModel(model.GranularMarkings)

			props.KillChainPhases = expandThreatIntelligenceKillChainPhaseModel(model.KillChainPhases)

			if len(model.ObjectMarking) > 0 {
				props.ObjectMarkingRefs = &model.ObjectMarking
			}

			if len(model.Labels) > 0 {
				props.Labels = &model.Labels
			}

			if len(model.ThreatTypes) > 0 {
				props.ThreatTypes = &model.ThreatTypes
			}

			if model.ValidFrom != "" {
				gmtLoc, _ := time.LoadLocation("GMT")
				t, err := time.Parse(time.RFC3339, model.ValidFrom)
				if err != nil {
					return err
				}
				validFromValue := t.In(gmtLoc).Format(time.RFC1123Z)
				props.ValidFrom = &validFromValue
			}

			if model.ValidUntil != "" {
				gmtLoc, _ := time.LoadLocation("GMT")
				t, err := time.Parse(time.RFC3339, model.ValidUntil)
				if err != nil {
					return err
				}
				validUntilValue := t.In(gmtLoc).Format(time.RFC1123Z)
				props.ValidUntil = &validUntilValue
			}

			resp, err := client.IndicatorCreateIndicator(ctx, *workspaceId, payload)
			if err != nil {
				return fmt.Errorf("creating Threat Intelligence Indicator in workspace %s: %+v", workspaceId, err)
			}

			info, ok := resp.Model.(threatintelligence.ThreatIntelligenceIndicatorModel)
			if !ok {
				return fmt.Errorf("creating Threat Intelligence Indicator in workspace %s: `model` type mismatch", workspaceId)
			}

			id, err := threatintelligence.ParseIndicatorID(pointer.From(info.Id))
			if err != nil {
				return err
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ThreatIntelligenceIndicator) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.ThreatIntelligenceClient

			id, err := threatintelligence.ParseIndicatorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model IndicatorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.IndicatorGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			v, ok := resp.Model.(threatintelligence.ThreatIntelligenceIndicatorModel)
			if !ok {
				return fmt.Errorf("retrieving %s: type mismatch, got %T", id, resp.Model)
			}

			if v.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}
			props := v.Properties

			if metadata.ResourceData.HasChange("confidence") {
				props.Confidence = pointer.To(model.Confidence)
				if model.Confidence == -1 {
					props.Confidence = nil
				}
			}

			if metadata.ResourceData.HasChange("created_by") {
				props.CreatedByRef = &model.CreatedByRef
			}

			if metadata.ResourceData.HasChange("description") {
				props.Description = &model.Description
			}

			if metadata.ResourceData.HasChange("display_name") {
				props.DisplayName = &model.DisplayName
			}

			if metadata.ResourceData.HasChange("extension") {
				extensionsValue, err := pluginsdk.ExpandJsonFromString(model.Extensions)
				if err != nil {
					return err
				}
				props.Extensions = &extensionsValue
			}

			if metadata.ResourceData.HasChange("external_reference") {
				props.ExternalReferences = expandThreatIntelligenceExternalReferenceModel(model.ExternalReference)
			}

			if metadata.ResourceData.HasChange("granular_marking") {
				props.GranularMarkings = expandThreatIntelligenceGranularMarkingModelModel(model.GranularMarkings)
			}

			if metadata.ResourceData.HasChange("indicator_type") {
				props.IndicatorTypes = &model.IndicatorTypes
			}

			if metadata.ResourceData.HasChange("kill_chain_phase") {
				props.KillChainPhases = expandThreatIntelligenceKillChainPhaseModel(model.KillChainPhases)
			}

			if metadata.ResourceData.HasChange("tags") {
				props.Labels = &model.Labels
			}

			if metadata.ResourceData.HasChange("language") {
				props.Language = &model.Language
			}

			if metadata.ResourceData.HasChange("object_marking_refs") {
				props.ObjectMarkingRefs = &model.ObjectMarking
			}

			if metadata.ResourceData.HasChange("pattern") {
				patternValue, err := expandIndicatorPattern(model.PatternType, model.Pattern)
				if err != nil {
					return err
				}
				props.Pattern = &patternValue
			}

			if metadata.ResourceData.HasChange("pattern_type") {
				props.PatternType = &model.PatternType
			}

			if metadata.ResourceData.HasChange("pattern_version") {
				props.PatternVersion = &model.PatternVersion
			}

			if metadata.ResourceData.HasChange("revoked") {
				props.Revoked = &model.Revoked
			}

			if metadata.ResourceData.HasChange("source") {
				props.Source = &model.Source
			}

			if metadata.ResourceData.HasChange("threat_types") {
				props.ThreatTypes = &model.ThreatTypes
			}

			if metadata.ResourceData.HasChange("validate_from_utc") {
				props.ValidFrom = &model.ValidFrom
			}

			if metadata.ResourceData.HasChange("validate_until_utc") {
				props.ValidUntil = &model.ValidUntil
			}

			if _, err := client.IndicatorCreate(ctx, *id, v); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ThreatIntelligenceIndicator) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.ThreatIntelligenceClient

			id, err := threatintelligence.ParseIndicatorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := client.IndicatorGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)
			state := IndicatorModel{
				Name:        id.IndicatorName,
				WorkspaceId: workspaceId.ID(),
			}
			if resp.Model != nil {
				if v, ok := resp.Model.(threatintelligence.ThreatIntelligenceIndicatorModel); ok {
					if props := v.Properties; props != nil {
						state.CreatedOn = pointer.From(props.Created)
						state.PatternType = pointer.From(props.PatternType)
						state.Revoked = pointer.From(props.Revoked)

						patternValue, err := flattenIndicatorPattern(pointer.From(props.Pattern))
						if err != nil {
							return err
						}
						state.Pattern = patternValue

						state.Confidence = -1
						if props.Confidence != nil {
							state.Confidence = *props.Confidence
						}

						state.CreatedByRef = pointer.From(props.CreatedByRef)
						state.Description = pointer.From(props.Description)
						state.DisplayName = pointer.From(props.DisplayName)

						extensionsValue, err := pluginsdk.FlattenJsonToString(pointer.From(props.Extensions))
						if err != nil {
							return err
						}
						state.Extensions = extensionsValue

						state.ExternalReference = flattenThreatIntelligenceExternalReferenceModel(props.ExternalReferences)
						state.GranularMarkings = flattenThreatIntelligenceGranularMarkingModelModel(props.GranularMarkings)
						state.KillChainPhases = flattenThreatIntelligenceKillChainPhaseModel(props.KillChainPhases)
						state.IndicatorTypes = pointer.From(props.IndicatorTypes)
						state.Language = pointer.From(props.Language)
						state.PatternVersion = pointer.From(props.PatternVersion)
						state.Source = pointer.From(props.Source)
						state.ObjectMarking = pointer.From(props.ObjectMarkingRefs)
						state.Labels = pointer.From(props.Labels)
						state.ThreatTypes = pointer.From(props.ThreatTypes)

						if pointer.From(props.ValidFrom) != "" {
							t, err := time.Parse(time.RFC3339, pointer.From(props.ValidFrom))
							if err != nil {
								return err
							}
							state.ValidFrom = t.Format(time.RFC3339)
						}

						if pointer.From(props.ValidUntil) != "" {
							t, err := time.Parse(time.RFC3339, pointer.From(props.ValidUntil))
							if err != nil {
								return err
							}
							state.ValidUntil = t.Format(time.RFC3339)
						}

						state.Defanged = pointer.From(props.Defanged)
						state.ExternalId = pointer.From(props.ExternalId)
						state.ExternalLastUpdatedTimeUtc = pointer.From(props.ExternalLastUpdatedTimeUtc)
						state.LastUpdatedTimeUtc = pointer.From(props.LastUpdatedTimeUtc)
						state.ParsedPattern = flattenIndicatorParsedPattern(props.ParsedPattern)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ThreatIntelligenceIndicator) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.ThreatIntelligenceClient

			id, err := threatintelligence.ParseIndicatorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.IndicatorDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandThreatIntelligenceExternalReferenceModel(inputList []externalReferenceModel) *[]threatintelligence.ThreatIntelligenceExternalReference {
	outputList := make([]threatintelligence.ThreatIntelligenceExternalReference, 0, len(inputList))
	for _, v := range inputList {
		input := v
		hashesValue := make(map[string]string)
		for k, hash := range input.Hashes {
			hashesValue[k] = hash
		}

		outputList = append(outputList, threatintelligence.ThreatIntelligenceExternalReference{
			Hashes:      &hashesValue,
			Description: pointer.ToOrNil(input.Description),
			ExternalId:  pointer.ToOrNil(input.ExternalId),
			SourceName:  pointer.ToOrNil(input.SourceName),
			Url:         pointer.ToOrNil(input.Url),
		})
	}

	return &outputList
}

func flattenThreatIntelligenceExternalReferenceModel(input *[]threatintelligence.ThreatIntelligenceExternalReference) []externalReferenceModel {
	output := make([]externalReferenceModel, 0)
	if input == nil {
		return output
	}
	for _, v := range *input {
		o := externalReferenceModel{
			SourceName:  pointer.From(v.SourceName),
			Url:         pointer.From(v.Url),
			Description: pointer.From(v.Description),
			ExternalId:  pointer.From(v.ExternalId),
		}
		if len(pointer.From(v.Hashes)) > 0 {
			o.Hashes = make(map[string]string)
			for k, hash := range pointer.From(v.Hashes) {
				o.Hashes[k] = hash
			}
		}
		output = append(output, o)
	}
	return output
}

func expandThreatIntelligenceGranularMarkingModelModel(inputList []granularMarkingModel) *[]threatintelligence.ThreatIntelligenceGranularMarkingModel {
	outputList := make([]threatintelligence.ThreatIntelligenceGranularMarkingModel, 0, len(inputList))
	for _, v := range inputList {
		outputList = append(outputList, threatintelligence.ThreatIntelligenceGranularMarkingModel{
			Language:   pointer.ToOrNil(v.Language),
			MarkingRef: pointer.To(v.MarkingRef),
			Selectors:  &v.Selectors,
		})
	}

	return &outputList
}

func flattenThreatIntelligenceGranularMarkingModelModel(input *[]threatintelligence.ThreatIntelligenceGranularMarkingModel) []granularMarkingModel {
	output := make([]granularMarkingModel, 0)
	if input == nil {
		return output
	}
	for _, v := range *input {
		output = append(output, granularMarkingModel{
			MarkingRef: pointer.From(v.MarkingRef),
			Selectors:  pointer.From(v.Selectors),
			Language:   pointer.From(v.Language),
		})
	}
	return output
}

func expandThreatIntelligenceKillChainPhaseModel(inputList []killChainPhaseModel) *[]threatintelligence.ThreatIntelligenceKillChainPhase {
	outputList := make([]threatintelligence.ThreatIntelligenceKillChainPhase, 0, len(inputList))
	for _, v := range inputList {
		outputList = append(outputList, threatintelligence.ThreatIntelligenceKillChainPhase{
			KillChainName: pointer.To(killChainName),
			PhaseName:     pointer.ToOrNil(v.PhaseName),
		})
	}

	return &outputList
}

func flattenThreatIntelligenceKillChainPhaseModel(input *[]threatintelligence.ThreatIntelligenceKillChainPhase) []killChainPhaseModel {
	output := make([]killChainPhaseModel, 0)
	if input == nil {
		return output
	}
	for _, v := range *input {
		output = append(output, killChainPhaseModel{
			PhaseName: pointer.From(v.PhaseName),
		})
	}
	return output
}

func expandIndicatorPattern(patternType string, pattern string) (string, error) {
	// possible values get from Portal
	// [domain-name:value = 'example.com']
	// [file:hashes.'MD5' = '6b0770e8133eee220333733931610598' ]
	// although the Portal support multiple hash, the service only accept one, so we ignore this type.
	// [file:hashes.'MD5' = '6b0770e8133eee220333733931610598' OR file:hashes.'MD5' = '6b0770e8133eee220333733931610591' ]
	// [ipv4-addr:value = '1.1.1.1']"
	// [ipv6-addr:value = '::1']
	// [url:value = 'http://www.example.com']
	if patternType == string(PatternTypeFile) {
		reg := regexp.MustCompile(`(MD5|SHA-1|SHA-256|SHA-512|MD6|RIPEMD-160|SHA-224|SHA3-224|SHA3-256|SHA3-384|SHA3-512|SHA-384|SSDEEPWHIRLPOOL):`)
		hashTypes := reg.FindStringSubmatch(pattern)
		if len(hashTypes) != 2 {
			return "", fmt.Errorf("when `pattern_type` set to `file`, `pattern` must combine a hash type with the hash code with a colon, an example is `MD5:78ecc5c05cd8b79af480df2f8fba0b9d`")
		}
		hashType := hashTypes[1]
		return fmt.Sprintf(`[file:hashes.'%s' = '%s']`, hashType, pattern), nil
	}
	return fmt.Sprintf(`[%s:value = '%s']`, patternType, pattern), nil
}

func flattenIndicatorPattern(input string) (string, error) {
	result := regexp.MustCompile(`\s=\s'(.+)'`).FindStringSubmatch(input)
	if len(result) == 2 {
		return result[1], nil
	}
	return "", fmt.Errorf("unable to parse pattern %s", input)
}

func flattenIndicatorParsedPattern(input *[]threatintelligence.ThreatIntelligenceParsedPattern) []parsedPatternModel {
	output := make([]parsedPatternModel, 0)
	if input == nil {
		return output
	}
	for _, v := range *input {
		o := parsedPatternModel{
			PatternTypeKey: pointer.From(v.PatternTypeKey),
		}
		if v.PatternTypeValues != nil {
			for _, p := range *v.PatternTypeValues {
				o.PatternTypeValues = append(o.PatternTypeValues, patternTypeValuesModel{
					Value:     pointer.From(p.Value),
					ValueType: pointer.From(p.ValueType),
				})
			}
		}
		output = append(output, o)
	}
	return output
}
