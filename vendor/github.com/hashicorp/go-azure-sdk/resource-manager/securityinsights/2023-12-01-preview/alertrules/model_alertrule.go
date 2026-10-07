package alertrules

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertRule interface {
	AlertRule() BaseAlertRuleImpl
}

var _ AlertRule = BaseAlertRuleImpl{}

type BaseAlertRuleImpl struct {
	Etag       *string                `json:"etag,omitempty"`
	Id         *string                `json:"id,omitempty"`
	Kind       AlertRuleKind          `json:"kind"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

func (s BaseAlertRuleImpl) AlertRule() BaseAlertRuleImpl {
	return s
}

var _ AlertRule = RawAlertRuleImpl{}

// RawAlertRuleImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawAlertRuleImpl struct {
	alertRule BaseAlertRuleImpl
	Type      string
	Values    map[string]interface{}
}

func (s RawAlertRuleImpl) AlertRule() BaseAlertRuleImpl {
	return s.alertRule
}

func (s RawAlertRuleImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalAlertRuleImplementation(input []byte) (AlertRule, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AlertRule into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Fusion") {
		var out FusionAlertRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FusionAlertRule: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MLBehaviorAnalytics") {
		var out MLBehaviorAnalyticsAlertRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MLBehaviorAnalyticsAlertRule: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MicrosoftSecurityIncidentCreation") {
		var out MicrosoftSecurityIncidentCreationAlertRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MicrosoftSecurityIncidentCreationAlertRule: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NRT") {
		var out NrtAlertRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NrtAlertRule: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Scheduled") {
		var out ScheduledAlertRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ScheduledAlertRule: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ThreatIntelligence") {
		var out ThreatIntelligenceAlertRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ThreatIntelligenceAlertRule: %+v", err)
		}
		return out, nil
	}

	var parent BaseAlertRuleImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseAlertRuleImpl: %+v", err)
	}

	return RawAlertRuleImpl{
		alertRule: parent,
		Type:      value,
		Values:    temp,
	}, nil

}
