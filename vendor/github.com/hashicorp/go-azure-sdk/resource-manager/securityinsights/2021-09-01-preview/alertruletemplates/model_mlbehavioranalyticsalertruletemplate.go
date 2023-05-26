package alertruletemplates

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AlertRuleTemplate = MLBehaviorAnalyticsAlertRuleTemplate{}

type MLBehaviorAnalyticsAlertRuleTemplate struct {
	Properties *MLBehaviorAnalyticsAlertRuleTemplateProperties `json:"properties,omitempty"`

	// Fields inherited from AlertRuleTemplate
	Id         *string                `json:"id,omitempty"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

var _ json.Marshaler = MLBehaviorAnalyticsAlertRuleTemplate{}

func (s MLBehaviorAnalyticsAlertRuleTemplate) MarshalJSON() ([]byte, error) {
	type wrapper MLBehaviorAnalyticsAlertRuleTemplate
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MLBehaviorAnalyticsAlertRuleTemplate: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MLBehaviorAnalyticsAlertRuleTemplate: %+v", err)
	}
	decoded["kind"] = "MLBehaviorAnalytics"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MLBehaviorAnalyticsAlertRuleTemplate: %+v", err)
	}

	return encoded, nil
}
