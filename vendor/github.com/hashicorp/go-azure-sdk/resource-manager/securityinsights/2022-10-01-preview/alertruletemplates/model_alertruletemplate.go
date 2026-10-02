package alertruletemplates

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertRuleTemplate interface {
	AlertRuleTemplate() BaseAlertRuleTemplateImpl
}

var _ AlertRuleTemplate = BaseAlertRuleTemplateImpl{}

type BaseAlertRuleTemplateImpl struct {
	Id         *string                `json:"id,omitempty"`
	Kind       AlertRuleKind          `json:"kind"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

func (s BaseAlertRuleTemplateImpl) AlertRuleTemplate() BaseAlertRuleTemplateImpl {
	return s
}

var _ AlertRuleTemplate = RawAlertRuleTemplateImpl{}

// RawAlertRuleTemplateImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawAlertRuleTemplateImpl struct {
	alertRuleTemplate BaseAlertRuleTemplateImpl
	Type              string
	Values            map[string]interface{}
}

func (s RawAlertRuleTemplateImpl) AlertRuleTemplate() BaseAlertRuleTemplateImpl {
	return s.alertRuleTemplate
}

func (s RawAlertRuleTemplateImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalAlertRuleTemplateImplementation(input []byte) (AlertRuleTemplate, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AlertRuleTemplate into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Fusion") {
		var out FusionAlertRuleTemplate
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FusionAlertRuleTemplate: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MLBehaviorAnalytics") {
		var out MLBehaviorAnalyticsAlertRuleTemplate
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MLBehaviorAnalyticsAlertRuleTemplate: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MicrosoftSecurityIncidentCreation") {
		var out MicrosoftSecurityIncidentCreationAlertRuleTemplate
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MicrosoftSecurityIncidentCreationAlertRuleTemplate: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NRT") {
		var out NrtAlertRuleTemplate
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NrtAlertRuleTemplate: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Scheduled") {
		var out ScheduledAlertRuleTemplate
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ScheduledAlertRuleTemplate: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ThreatIntelligence") {
		var out ThreatIntelligenceAlertRuleTemplate
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ThreatIntelligenceAlertRuleTemplate: %+v", err)
		}
		return out, nil
	}

	var parent BaseAlertRuleTemplateImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseAlertRuleTemplateImpl: %+v", err)
	}

	return RawAlertRuleTemplateImpl{
		alertRuleTemplate: parent,
		Type:              value,
		Values:            temp,
	}, nil

}
