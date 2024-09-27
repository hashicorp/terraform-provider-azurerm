// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
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
	ExternalRefrence           []externalReferenceModel `tfschema:"external_reference"`
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
	return validate.ThreatIntelligenceIndicatorID
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
			Type:         pluginsdk.TypeString,
			Optional:     true,
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

			client := azuresdkhacks.ThreatIntelligenceIndicatorClient{
				BaseClient: metadata.Client.Sentinel.ThreatIntelligenceClient.BaseClient,
			}
			workspaceId, err := workspaces.ParseWorkspaceID(model.WorkspaceId)
			if err != nil {
				return fmt.Errorf("parsing Workspace id %s: %+v", model.WorkspaceId, err)
			}

			patternValue, err := expandIndicatorPattern(model.PatternType, model.Pattern)
			if err != nil {
				return err
			}

			// it could not get the indicator by name before it has been created, because the name is generated by service side.
			// but we can not create duplicated indicator with same values, so list the values and find the existing one.
			existingIndicators, err := queryIndicatorsList(ctx, client, workspaceId)
			if err != nil {
				return fmt.Errorf("listing indicators: %+v", err)
			}
			for _, indicator := range existingIndicators {
				if indicator.PatternType != nil && *indicator.PatternType == model.PatternType {
					if indicator.Pattern != nil && *indicator.Pattern == patternValue {
						if indicator.ID != nil && *indicator.ID != "" {
							return tf.ImportAsExistsError("azurerm_sentinel_threat_intelligence_indicator", *indicator.ID)
						}
						return fmt.Errorf("checking existing indicator: `id` is nil")
					}
				}
			}

			properties := azuresdkhacks.ThreatIntelligenceIndicatorModel{
				Kind: securityinsight.KindBasicThreatIntelligenceInformationKindIndicator,
				ThreatIntelligenceIndicatorProperties: &azuresdkhacks.ThreatIntelligenceIndicatorProperties{
					PatternType: &model.PatternType,
					Revoked:     &model.Revoked,
				},
			}

			props := properties.ThreatIntelligenceIndicatorProperties

			props.Pattern = &patternValue

			if model.Confidence != -1 {
				props.Confidence = pointer.To(int32(model.Confidence))
			}

			if model.CreatedByRef != "" {
				props.CreatedByRef = &model.CreatedByRef
			}

			if model.Description != "" {
				props.Description = &model.Description
			}

			if model.DisplayName != "" {
				props.DisplayName = &model.DisplayName
			}

			if model.Extensions != "" {
				extensionsValue, err := pluginsdk.ExpandJsonFromString(model.Extensions)
				if err != nil {
					return err
				}
				props.Extensions = extensionsValue
			}

			props.ExternalReferences = expandThreatIntelligenceExternalReferenceModel(model.ExternalRefrence)

			props.GranularMarkings = expandThreatIntelligenceGranularMarkingModelModel(model.GranularMarkings)

			props.KillChainPhases = expandThreatIntelligenceKillChainPhaseModel(model.KillChainPhases)

			if model.Language != "" {
				props.Language = &model.Language
			}

			if model.PatternVersion != "" {
				props.PatternVersion = &model.PatternVersion
			}

			if model.Source != "" {
				props.Source = &model.Source
			}

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

			resp, err := client.CreateIndicator(ctx, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, properties)
			if err != nil {
				return fmt.Errorf("creating threaten intelligence indicator in workspace %s: %+v", workspaceId, err)
			}

			info, ok := resp.Value.AsThreatIntelligenceIndicatorModel()
			if !ok {
				return fmt.Errorf("creating threaten intelligence indicator in workspace %s: `model` type mismatch", workspaceId)
			}

			id, err := parse.ThreatIntelligenceIndicatorID(*info.ID)
			if err != nil {
				return fmt.Errorf("parsing threat intelligence indicator id %s: %+v", *info.ID, err)
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
			client := azuresdkhacks.ThreatIntelligenceIndicatorClient{
				BaseClient: metadata.Client.Sentinel.ThreatIntelligenceClient.BaseClient,
			}

			id, err := parse.ThreatIntelligenceIndicatorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model IndicatorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.IndicatorName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties, ok := resp.Value.AsThreatIntelligenceIndicatorModel()
			if !ok {
				return fmt.Errorf("retrieving %s: type mismatch", id)
			}

			if metadata.ResourceData.HasChange("confidence") {
				if model.Confidence == -1 {
					properties.Confidence = nil
				} else {
					properties.Confidence = pointer.To(int32(model.Confidence))
				}
			}

			if metadata.ResourceData.HasChange("created_by") {
				properties.CreatedByRef = &model.CreatedByRef
			}

			if metadata.ResourceData.HasChange("description") {
				properties.Description = &model.Description
			}

			if metadata.ResourceData.HasChange("display_name") {
				properties.DisplayName = &model.DisplayName
			}

			if metadata.ResourceData.HasChange("extension") {
				extensionsValue, err := pluginsdk.ExpandJsonFromString(model.Extensions)
				if err != nil {
					return err
				}
				properties.Extensions = extensionsValue
			}

			if metadata.ResourceData.HasChange("external_reference") {
				properties.ExternalReferences = expandThreatIntelligenceExternalReferenceModel(model.ExternalRefrence)

			}

			if metadata.ResourceData.HasChange("granular_marking") {
				properties.GranularMarkings = expandThreatIntelligenceGranularMarkingModelModel(model.GranularMarkings)
			}

			if metadata.ResourceData.HasChange("indicator_type") {
				properties.IndicatorTypes = &model.IndicatorTypes
			}

			if metadata.ResourceData.HasChange("kill_chain_phase") {
				properties.KillChainPhases = expandThreatIntelligenceKillChainPhaseModel(model.KillChainPhases)
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Labels = &model.Labels
			}

			if metadata.ResourceData.HasChange("language") {
				properties.Language = &model.Language
			}

			if metadata.ResourceData.HasChange("object_marking_refs") {
				properties.ObjectMarkingRefs = &model.ObjectMarking
			}

			if metadata.ResourceData.HasChange("pattern") {
				patternValue, err := expandIndicatorPattern(model.PatternType, model.Pattern)
				if err != nil {
					return err
				}
				properties.Pattern = &patternValue
			}

			if metadata.ResourceData.HasChange("pattern_type") {
				properties.PatternType = &model.PatternType
			}

			if metadata.ResourceData.HasChange("pattern_version") {
				properties.PatternVersion = &model.PatternVersion
			}

			if metadata.ResourceData.HasChange("revoked") {
				properties.Revoked = &model.Revoked
			}

			if metadata.ResourceData.HasChange("source") {
				properties.Source = &model.Source
			}

			if metadata.ResourceData.HasChange("threat_types") {
				properties.ThreatTypes = &model.ThreatTypes
			}

			if metadata.ResourceData.HasChange("validate_from_utc") {
				properties.ValidFrom = &model.ValidFrom
			}

			if metadata.ResourceData.HasChange("validate_until_utc") {
				properties.ValidUntil = &model.ValidUntil
			}

			if _, err := client.Create(ctx, id.ResourceGroup, id.WorkspaceName, id.IndicatorName, *properties); err != nil {
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
			client := azuresdkhacks.ThreatIntelligenceIndicatorClient{
				BaseClient: metadata.Client.Sentinel.ThreatIntelligenceClient.BaseClient,
			}
			id, err := parse.ThreatIntelligenceIndicatorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)
			resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.IndicatorName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model, ok := resp.Value.AsThreatIntelligenceIndicatorModel()
			if !ok {
				return fmt.Errorf("retrieving %s: type mismatch", id)
			}

			state := IndicatorModel{
				Name:        pointer.From(model.Name),
				CreatedOn:   pointer.From(model.Created),
				WorkspaceId: workspaceId.ID(),
				PatternType: pointer.From(model.PatternType),
				Revoked:     pointer.From(model.Revoked),
			}

			patternValue, err := flattenIndicatorPattern(*model.Pattern)
			if err != nil {
				return err
			}
			state.Pattern = patternValue

			if model.Confidence != nil {
				state.Confidence = int64(*model.Confidence)
			} else {
				state.Confidence = -1
			}

			if model.CreatedByRef != nil {
				state.CreatedByRef = *model.CreatedByRef
			}

			if model.Description != nil {
				state.Description = *model.Description
			}

			if model.DisplayName != nil {
				state.DisplayName = *model.DisplayName
			}

			if model.Extensions != nil && len(model.Extensions) > 0 {
				extensionsValue, err := pluginsdk.FlattenJsonToString(model.Extensions)
				if err != nil {
					return err
				}
				state.Extensions = extensionsValue
			}

			state.ExternalRefrence = flattenThreatIntelligenceExternalReferenceModel(model.ExternalReferences)

			state.GranularMarkings = flattenThreatIntelligenceGranularMarkingModelModel(model.GranularMarkings)

			state.KillChainPhases = flattenThreatIntelligenceKillChainPhaseModel(model.KillChainPhases)

			if model.IndicatorTypes != nil && len(*model.IndicatorTypes) > 0 {
				state.IndicatorTypes = *model.IndicatorTypes
			}

			if model.Language != nil {
				state.Language = *model.Language
			}

			if model.PatternVersion != nil {
				state.PatternVersion = *model.PatternVersion
			}

			if model.Source != nil {
				state.Source = *model.Source
			}

			if model.ObjectMarkingRefs != nil && len(*model.ObjectMarkingRefs) > 0 {
				state.ObjectMarking = *model.ObjectMarkingRefs
			}

			if model.Labels != nil && len(*model.Labels) > 0 {
				state.Labels = *model.Labels
			}

			if model.ThreatTypes != nil && len(*model.ThreatTypes) > 0 {
				state.ThreatTypes = *model.ThreatTypes
			}

			if model.ValidFrom != nil && *model.ValidFrom != "" {
				t, err := time.Parse(time.RFC3339, *model.ValidFrom)
				if err != nil {
					return err
				}
				state.ValidFrom = t.Format(time.RFC3339)
			}

			if model.ValidUntil != nil && *model.ValidUntil != "" {
				t, err := time.Parse(time.RFC3339, *model.ValidUntil)
				if err != nil {
					return err
				}
				state.ValidUntil = t.Format(time.RFC3339)
			}

			if model.Defanged != nil {
				state.Defanged = *model.Defanged
			}

			if model.ExternalID != nil {
				state.ExternalId = *model.ExternalID
			}

			if model.ExternalLastUpdatedTimeUtc != nil {
				state.ExternalLastUpdatedTimeUtc = *model.ExternalLastUpdatedTimeUtc
			}

			if model.LastUpdatedTimeUtc != nil {
				state.LastUpdatedTimeUtc = *model.LastUpdatedTimeUtc
			}

			state.ParsedPattern = flattenIndicatorParsedPattern(model.ParsedPattern)

			return metadata.Encode(&state)
		},
	}
}

func (r ThreatIntelligenceIndicator) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.ThreatIntelligenceClient

			id, err := parse.ThreatIntelligenceIndicatorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.IndicatorName); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandThreatIntelligenceExternalReferenceModel(inputList []externalReferenceModel) *[]securityinsight.ThreatIntelligenceExternalReference {
	var outputList []securityinsight.ThreatIntelligenceExternalReference
	for _, v := range inputList {
		input := v
		hashesValue := make(map[string]*string, 0)
		for k, hash := range input.Hashes {
			hashesValue[k] = &hash
		}

		output := securityinsight.ThreatIntelligenceExternalReference{
			Hashes: hashesValue,
		}

		if input.Description != "" {
			output.Description = &input.Description
		}

		if input.ExternalId != "" {
			output.ExternalID = &input.ExternalId
		}

		if input.SourceName != "" {
			output.SourceName = &input.SourceName
		}

		if input.Url != "" {
			output.URL = &input.Url
		}

		outputList = append(outputList, output)
	}

	return &outputList
}

func flattenThreatIntelligenceExternalReferenceModel(input *[]securityinsight.ThreatIntelligenceExternalReference) []externalReferenceModel {
	output := make([]externalReferenceModel, 0)
	if input == nil {
		return output
	}
	for _, v := range *input {
		o := externalReferenceModel{}
		if v.ExternalID != nil {
			o.ExternalId = *v.ExternalID
		}
		if v.URL != nil {
			o.Url = *v.URL
		}
		if v.SourceName != nil {
			o.SourceName = *v.SourceName
		}
		if v.Description != nil {
			o.Description = *v.Description
		}
		if len(v.Hashes) > 0 {
			o.Hashes = make(map[string]string, 0)
			for k, hash := range v.Hashes {
				o.Hashes[k] = *hash
			}
		}
		output = append(output, o)
	}
	return output
}

func expandThreatIntelligenceGranularMarkingModelModel(inputList []granularMarkingModel) *[]azuresdkhacks.ThreatIntelligenceGranularMarkingModel {
	var outputList []azuresdkhacks.ThreatIntelligenceGranularMarkingModel
	for _, v := range inputList {
		input := v
		output := azuresdkhacks.ThreatIntelligenceGranularMarkingModel{
			MarkingRef: &input.MarkingRef,
			Selectors:  &input.Selectors,
		}

		if input.Language != "" {
			output.Language = &input.Language
		}

		outputList = append(outputList, output)
	}

	return &outputList
}

func flattenThreatIntelligenceGranularMarkingModelModel(input *[]azuresdkhacks.ThreatIntelligenceGranularMarkingModel) []granularMarkingModel {
	output := make([]granularMarkingModel, 0)
	if input == nil {
		return output
	}
	for _, v := range *input {
		o := granularMarkingModel{}
		if v.MarkingRef != nil {
			o.MarkingRef = *v.MarkingRef
		}
		if v.Selectors != nil {
			o.Selectors = *v.Selectors
		}
		if v.Language != nil {
			o.Language = *v.Language
		}
		output = append(output, o)
	}
	return output
}

func expandThreatIntelligenceKillChainPhaseModel(inputList []killChainPhaseModel) *[]securityinsight.ThreatIntelligenceKillChainPhase {
	var outputList []securityinsight.ThreatIntelligenceKillChainPhase
	for _, v := range inputList {
		input := v
		output := securityinsight.ThreatIntelligenceKillChainPhase{
			KillChainName: utils.String(killChainName),
		}

		if input.PhaseName != "" {
			output.PhaseName = &input.PhaseName
		}

		outputList = append(outputList, output)
	}

	return &outputList
}

func flattenThreatIntelligenceKillChainPhaseModel(input *[]securityinsight.ThreatIntelligenceKillChainPhase) []killChainPhaseModel {
	output := make([]killChainPhaseModel, 0)
	if input == nil {
		return output
	}
	for _, v := range *input {
		o := killChainPhaseModel{}
		if v.PhaseName != nil {
			o.PhaseName = *v.PhaseName
		}
		output = append(output, o)
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
	reg := regexp.MustCompile(`\s=\s\'(.+)\'`)
	result := reg.FindStringSubmatch(input)
	if len(result) == 2 {
		return result[1], nil
	}
	return "", fmt.Errorf("unable to parse pattern %s", input)
}

func flattenIndicatorParsedPattern(input *[]securityinsight.ThreatIntelligenceParsedPattern) []parsedPatternModel {
	output := make([]parsedPatternModel, 0)
	if input == nil {
		return output
	}
	for _, v := range *input {
		o := parsedPatternModel{}
		if v.PatternTypeKey != nil {
			o.PatternTypeKey = *v.PatternTypeKey
		}
		if v.PatternTypeValues != nil {
			for _, patternTypeValue := range *v.PatternTypeValues {
				valueOutput := patternTypeValuesModel{}
				if patternTypeValue.Value != nil {
					valueOutput.Value = *patternTypeValue.Value
				}
				if patternTypeValue.ValueType != nil {
					valueOutput.ValueType = *patternTypeValue.ValueType
				}
				o.PatternTypeValues = append(o.PatternTypeValues, valueOutput)
			}
		}
		output = append(output, o)
	}
	return output
}

func queryIndicatorsList(ctx context.Context, client azuresdkhacks.ThreatIntelligenceIndicatorClient, workspaceId *workspaces.WorkspaceId) ([]*azuresdkhacks.ThreatIntelligenceIndicatorModel, error) {
	output := make([]*azuresdkhacks.ThreatIntelligenceIndicatorModel, 0)
	resp, err := client.QueryIndicators(ctx, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, securityinsight.ThreatIntelligenceFilteringCriteria{})
	if err != nil {
		return output, err
	}
	for resp.NotDone() {
		for _, indicator := range resp.Values() {
			indicator, ok := indicator.AsThreatIntelligenceIndicatorModel()
			if !ok {
				continue
			}
			output = append(output, indicator)
		}
		if err := resp.NextWithContext(ctx); err != nil {
			return output, err
		}
	}
	return output, nil
}
