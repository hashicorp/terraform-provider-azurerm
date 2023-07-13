// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azuresdkhacks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

// TODO 4.0: check if this could be removed.
// workaround for https://github.com/Azure/azure-rest-api-specs/issues/22893

type ThreatIntelligenceInformationModel struct {
	autorest.Response `json:"-"`
	Value             BasicThreatIntelligenceInformation `json:"value,omitempty"`
}

func (tiim *ThreatIntelligenceInformationModel) UnmarshalJSON(body []byte) error {
	tii, err := unmarshalBasicThreatIntelligenceInformation(body)
	if err != nil {
		return err
	}
	tiim.Value = tii

	return nil
}

func unmarshalBasicThreatIntelligenceInformation(body []byte) (BasicThreatIntelligenceInformation, error) {
	var m map[string]interface{}
	err := json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}

	switch m["kind"] {
	case string(securityinsight.KindBasicThreatIntelligenceInformationKindIndicator):
		var tiim ThreatIntelligenceIndicatorModel
		err := json.Unmarshal(body, &tiim)
		return tiim, err
	default:
		// it's not used in this hack.
		return nil, fmt.Errorf("get unknown kind: %s", m["kind"])
	}
}

type BasicThreatIntelligenceInformation interface {
	AsThreatIntelligenceIndicatorModel() (*ThreatIntelligenceIndicatorModel, bool)
	AsThreatIntelligenceInformation() (*securityinsight.ThreatIntelligenceInformation, bool)
}

type ThreatIntelligenceIndicatorModel struct {
	*ThreatIntelligenceIndicatorProperties `json:"properties,omitempty"`
	Etag                                   *string                                                `json:"etag,omitempty"`
	ID                                     *string                                                `json:"id,omitempty"`
	Name                                   *string                                                `json:"name,omitempty"`
	Type                                   *string                                                `json:"type,omitempty"`
	SystemData                             *securityinsight.SystemData                            `json:"systemData,omitempty"`
	Kind                                   securityinsight.KindBasicThreatIntelligenceInformation `json:"kind,omitempty"`
}

func (tiim ThreatIntelligenceIndicatorModel) MarshalJSON() ([]byte, error) {
	tiim.Kind = securityinsight.KindBasicThreatIntelligenceInformationKindIndicator
	objectMap := make(map[string]interface{})
	if tiim.ThreatIntelligenceIndicatorProperties != nil {
		objectMap["properties"] = tiim.ThreatIntelligenceIndicatorProperties
	}
	if tiim.Kind != "" {
		objectMap["kind"] = tiim.Kind
	}
	if tiim.Etag != nil {
		objectMap["etag"] = tiim.Etag
	}
	return json.Marshal(objectMap)
}

func (tiim ThreatIntelligenceIndicatorModel) AsThreatIntelligenceIndicatorModel() (*ThreatIntelligenceIndicatorModel, bool) {
	return &tiim, true
}

func (tiim ThreatIntelligenceIndicatorModel) AsThreatIntelligenceInformation() (*securityinsight.ThreatIntelligenceInformation, bool) {
	return nil, false
}

func (tiim *ThreatIntelligenceIndicatorModel) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "properties":
			if v != nil {
				var threatIntelligenceIndicatorProperties ThreatIntelligenceIndicatorProperties
				err = json.Unmarshal(*v, &threatIntelligenceIndicatorProperties)
				if err != nil {
					return err
				}
				tiim.ThreatIntelligenceIndicatorProperties = &threatIntelligenceIndicatorProperties
			}
		case "kind":
			if v != nil {
				var kind securityinsight.KindBasicThreatIntelligenceInformation
				err = json.Unmarshal(*v, &kind)
				if err != nil {
					return err
				}
				tiim.Kind = kind
			}
		case "etag":
			if v != nil {
				var etag string
				err = json.Unmarshal(*v, &etag)
				if err != nil {
					return err
				}
				tiim.Etag = &etag
			}
		case "id":
			if v != nil {
				var ID string
				err = json.Unmarshal(*v, &ID)
				if err != nil {
					return err
				}
				tiim.ID = &ID
			}
		case "name":
			if v != nil {
				var name string
				err = json.Unmarshal(*v, &name)
				if err != nil {
					return err
				}
				tiim.Name = &name
			}
		case "type":
			if v != nil {
				var typeVar string
				err = json.Unmarshal(*v, &typeVar)
				if err != nil {
					return err
				}
				tiim.Type = &typeVar
			}
		case "systemData":
			if v != nil {
				var systemData securityinsight.SystemData
				err = json.Unmarshal(*v, &systemData)
				if err != nil {
					return err
				}
				tiim.SystemData = &systemData
			}
		}
	}

	return nil
}

type ThreatIntelligenceIndicatorProperties struct {
	ThreatIntelligenceTags     *[]string                                              `json:"threatIntelligenceTags,omitempty"`
	LastUpdatedTimeUtc         *string                                                `json:"lastUpdatedTimeUtc,omitempty"`
	Source                     *string                                                `json:"source,omitempty"`
	DisplayName                *string                                                `json:"displayName,omitempty"`
	Description                *string                                                `json:"description,omitempty"`
	IndicatorTypes             *[]string                                              `json:"indicatorTypes,omitempty"`
	Pattern                    *string                                                `json:"pattern,omitempty"`
	PatternType                *string                                                `json:"patternType,omitempty"`
	PatternVersion             *string                                                `json:"patternVersion,omitempty"`
	KillChainPhases            *[]securityinsight.ThreatIntelligenceKillChainPhase    `json:"killChainPhases,omitempty"`
	ParsedPattern              *[]securityinsight.ThreatIntelligenceParsedPattern     `json:"parsedPattern,omitempty"`
	ExternalID                 *string                                                `json:"externalId,omitempty"`
	CreatedByRef               *string                                                `json:"createdByRef,omitempty"`
	Defanged                   *bool                                                  `json:"defanged,omitempty"`
	ExternalLastUpdatedTimeUtc *string                                                `json:"externalLastUpdatedTimeUtc,omitempty"`
	ExternalReferences         *[]securityinsight.ThreatIntelligenceExternalReference `json:"externalReferences,omitempty"`
	GranularMarkings           *[]ThreatIntelligenceGranularMarkingModel              `json:"granularMarkings,omitempty"`
	Labels                     *[]string                                              `json:"labels,omitempty"`
	Revoked                    *bool                                                  `json:"revoked,omitempty"`
	Confidence                 *int32                                                 `json:"confidence,omitempty"`
	ObjectMarkingRefs          *[]string                                              `json:"objectMarkingRefs,omitempty"`
	Language                   *string                                                `json:"language,omitempty"`
	ThreatTypes                *[]string                                              `json:"threatTypes,omitempty"`
	ValidFrom                  *string                                                `json:"validFrom,omitempty"`
	ValidUntil                 *string                                                `json:"validUntil,omitempty"`
	Created                    *string                                                `json:"created,omitempty"`
	Modified                   *string                                                `json:"modified,omitempty"`
	Extensions                 map[string]interface{}                                 `json:"extensions"`
	AdditionalData             map[string]interface{}                                 `json:"additionalData"`
	FriendlyName               *string                                                `json:"friendlyName,omitempty"`
}

func (tiip ThreatIntelligenceIndicatorProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if tiip.ThreatIntelligenceTags != nil {
		objectMap["threatIntelligenceTags"] = tiip.ThreatIntelligenceTags
	}
	if tiip.LastUpdatedTimeUtc != nil {
		objectMap["lastUpdatedTimeUtc"] = tiip.LastUpdatedTimeUtc
	}
	if tiip.Source != nil {
		objectMap["source"] = tiip.Source
	}
	if tiip.DisplayName != nil {
		objectMap["displayName"] = tiip.DisplayName
	}
	if tiip.Description != nil {
		objectMap["description"] = tiip.Description
	}
	if tiip.IndicatorTypes != nil {
		objectMap["indicatorTypes"] = tiip.IndicatorTypes
	}
	if tiip.Pattern != nil {
		objectMap["pattern"] = tiip.Pattern
	}
	if tiip.PatternType != nil {
		objectMap["patternType"] = tiip.PatternType
	}
	if tiip.PatternVersion != nil {
		objectMap["patternVersion"] = tiip.PatternVersion
	}
	if tiip.KillChainPhases != nil {
		objectMap["killChainPhases"] = tiip.KillChainPhases
	}
	if tiip.ParsedPattern != nil {
		objectMap["parsedPattern"] = tiip.ParsedPattern
	}
	if tiip.ExternalID != nil {
		objectMap["externalId"] = tiip.ExternalID
	}
	if tiip.CreatedByRef != nil {
		objectMap["createdByRef"] = tiip.CreatedByRef
	}
	if tiip.Defanged != nil {
		objectMap["defanged"] = tiip.Defanged
	}
	if tiip.ExternalLastUpdatedTimeUtc != nil {
		objectMap["externalLastUpdatedTimeUtc"] = tiip.ExternalLastUpdatedTimeUtc
	}
	if tiip.ExternalReferences != nil {
		objectMap["externalReferences"] = tiip.ExternalReferences
	}
	if tiip.GranularMarkings != nil {
		objectMap["granularMarkings"] = tiip.GranularMarkings
	}
	if tiip.Labels != nil {
		objectMap["labels"] = tiip.Labels
	}
	if tiip.Revoked != nil {
		objectMap["revoked"] = tiip.Revoked
	}
	if tiip.Confidence != nil {
		objectMap["confidence"] = tiip.Confidence
	}
	if tiip.ObjectMarkingRefs != nil {
		objectMap["objectMarkingRefs"] = tiip.ObjectMarkingRefs
	}
	if tiip.Language != nil {
		objectMap["language"] = tiip.Language
	}
	if tiip.ThreatTypes != nil {
		objectMap["threatTypes"] = tiip.ThreatTypes
	}
	if tiip.ValidFrom != nil {
		objectMap["validFrom"] = tiip.ValidFrom
	}
	if tiip.ValidUntil != nil {
		objectMap["validUntil"] = tiip.ValidUntil
	}
	if tiip.Created != nil {
		objectMap["created"] = tiip.Created
	}
	if tiip.Modified != nil {
		objectMap["modified"] = tiip.Modified
	}
	if tiip.Extensions != nil {
		objectMap["extensions"] = tiip.Extensions
	}
	return json.Marshal(objectMap)
}

type ThreatIntelligenceGranularMarkingModel struct {
	Language   *string   `json:"language,omitempty"`
	MarkingRef *string   `json:"markingRef,omitempty"`
	Selectors  *[]string `json:"selectors,omitempty"`
}

type ThreatIntelligenceInformationListPage struct {
	fn   func(context.Context, ThreatIntelligenceInformationList) (ThreatIntelligenceInformationList, error)
	tiil ThreatIntelligenceInformationList
}

type ThreatIntelligenceInformationList struct {
	autorest.Response `json:"-"`
	NextLink          *string                               `json:"nextLink,omitempty"`
	Value             *[]BasicThreatIntelligenceInformation `json:"value,omitempty"`
}

func (tiil ThreatIntelligenceInformationList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if tiil.Value != nil {
		objectMap["value"] = tiil.Value
	}
	return json.Marshal(objectMap)
}

func (tiil *ThreatIntelligenceInformationList) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "nextLink":
			if v != nil {
				var nextLink string
				err = json.Unmarshal(*v, &nextLink)
				if err != nil {
					return err
				}
				tiil.NextLink = &nextLink
			}
		case "value":
			if v != nil {
				value, err := unmarshalBasicThreatIntelligenceInformationArray(*v)
				if err != nil {
					return err
				}
				tiil.Value = &value
			}
		}
	}

	return nil
}

func unmarshalBasicThreatIntelligenceInformationArray(body []byte) ([]BasicThreatIntelligenceInformation, error) {
	var rawMessages []*json.RawMessage
	err := json.Unmarshal(body, &rawMessages)
	if err != nil {
		return nil, err
	}

	tiiArray := make([]BasicThreatIntelligenceInformation, len(rawMessages))

	for index, rawMessage := range rawMessages {
		tii, err := unmarshalBasicThreatIntelligenceInformation(*rawMessage)
		if err != nil {
			return nil, err
		}
		tiiArray[index] = tii
	}
	return tiiArray, nil
}

func (page *ThreatIntelligenceInformationListPage) NextWithContext(ctx context.Context) (err error) {
	for {
		next, err := page.fn(ctx, page.tiil)
		if err != nil {
			return err
		}
		page.tiil = next
		if !next.hasNextLink() || !next.IsEmpty() {
			break
		}
	}
	return nil
}

func (tiil ThreatIntelligenceInformationList) IsEmpty() bool {
	return tiil.Value == nil || len(*tiil.Value) == 0
}

func (tiil ThreatIntelligenceInformationList) hasNextLink() bool {
	return tiil.NextLink != nil && len(*tiil.NextLink) != 0
}

func (page *ThreatIntelligenceInformationListPage) Next() error {
	return page.NextWithContext(context.Background())
}

func (page ThreatIntelligenceInformationListPage) NotDone() bool {
	return !page.tiil.IsEmpty()
}

func (page ThreatIntelligenceInformationListPage) Response() ThreatIntelligenceInformationList {
	return page.tiil
}

// Values returns the slice of values for the current page or nil if there are no values.
func (page ThreatIntelligenceInformationListPage) Values() []BasicThreatIntelligenceInformation {
	if page.tiil.IsEmpty() {
		return nil
	}
	return *page.tiil.Value
}
