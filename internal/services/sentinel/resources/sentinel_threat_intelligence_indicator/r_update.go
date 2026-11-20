// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_threat_intelligence_indicator

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/threatintelligence"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Update() sdk.ResourceFunc {
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
			previousEtag := pointer.From(v.Etag)
			props := v.Properties
			currentEtag := v.Etag

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

			props.ExternalLastUpdatedTimeUtc = nil
			props.LastUpdatedTimeUtc = nil
			if metadata.ResourceData.HasChangesExcept("tags") {
				resp, err := client.IndicatorCreate(ctx, *id, v)
				if err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}

				if updated, ok := resp.Model.(threatintelligence.ThreatIntelligenceIndicatorModel); ok {
					currentEtag = updated.Etag
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				payload := threatintelligence.ThreatIntelligenceIndicatorModel{
					Etag: currentEtag,
					Kind: threatintelligence.ThreatIntelligenceResourceKindEnumIndicator,
					Properties: &threatintelligence.ThreatIntelligenceIndicatorProperties{
						ThreatIntelligenceTags: &model.Labels,
					},
				}
				if _, err := client.IndicatorReplaceTags(ctx, *id, payload); err != nil {
					return fmt.Errorf("replacing tags for %s: %+v", *id, err)
				}
			}

			// Normal indicator updates change the ETag but do not change LastUpdatedTimeUtc. Wait until GET consistently
			// returns a version newer than the one read before the update.
			pollerType := custompollers.NewThreatIntelligenceIndicatorUpdatePoller(client, *id, previousEtag)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", *id, err)
			}

			return nil
		},
	}
}
