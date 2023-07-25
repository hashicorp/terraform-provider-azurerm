package alertrules

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AlertRule = MLBehaviorAnalyticsAlertRule{}

type MLBehaviorAnalyticsAlertRule struct {
	Properties *MLBehaviorAnalyticsAlertRuleProperties `json:"properties,omitempty"`

	// Fields inherited from AlertRule
	Etag       *string                `json:"etag,omitempty"`
	Id         *string                `json:"id,omitempty"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

var _ json.Marshaler = MLBehaviorAnalyticsAlertRule{}

func (s MLBehaviorAnalyticsAlertRule) MarshalJSON() ([]byte, error) {
	type wrapper MLBehaviorAnalyticsAlertRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MLBehaviorAnalyticsAlertRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MLBehaviorAnalyticsAlertRule: %+v", err)
	}
	decoded["kind"] = "MLBehaviorAnalytics"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MLBehaviorAnalyticsAlertRule: %+v", err)
	}

	return encoded, nil
}
