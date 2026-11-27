// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_threat_intelligence_indicator

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/threatintelligence"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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

type IndicatorPatternType string

const (
	PatternTypeDomainName IndicatorPatternType = "domain-name"
	PatternTypeFile       IndicatorPatternType = "file"
	PatternTypeIpV4Addr   IndicatorPatternType = "ipv4-addr"
	PatternTypeIpV6Addr   IndicatorPatternType = "ipv6-addr"
	PatternTypeUrl        IndicatorPatternType = "url"
)

const killChainName = "Lockheed Martin - Intrusion Kill Chain"
