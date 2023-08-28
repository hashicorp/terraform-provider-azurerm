package alertrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertRule interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawAlertRuleImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalAlertRuleImplementation(input []byte) (AlertRule, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AlertRule into map[string]interface: %+v", err)
	}

	value, ok := temp["kind"].(string)
	if !ok {
		return nil, nil
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

	out := RawAlertRuleImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
