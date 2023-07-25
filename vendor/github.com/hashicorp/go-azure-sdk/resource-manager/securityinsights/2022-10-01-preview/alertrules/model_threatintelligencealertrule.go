package alertrules

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AlertRule = ThreatIntelligenceAlertRule{}

type ThreatIntelligenceAlertRule struct {
	Properties *ThreatIntelligenceAlertRuleProperties `json:"properties,omitempty"`

	// Fields inherited from AlertRule
	Etag       *string                `json:"etag,omitempty"`
	Id         *string                `json:"id,omitempty"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

var _ json.Marshaler = ThreatIntelligenceAlertRule{}

func (s ThreatIntelligenceAlertRule) MarshalJSON() ([]byte, error) {
	type wrapper ThreatIntelligenceAlertRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ThreatIntelligenceAlertRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ThreatIntelligenceAlertRule: %+v", err)
	}
	decoded["kind"] = "ThreatIntelligence"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ThreatIntelligenceAlertRule: %+v", err)
	}

	return encoded, nil
}
