package alertruletemplates

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AlertRuleTemplate = ScheduledAlertRuleTemplate{}

type ScheduledAlertRuleTemplate struct {
	Properties *ScheduledAlertRuleTemplateProperties `json:"properties,omitempty"`

	// Fields inherited from AlertRuleTemplate
	Id         *string                `json:"id,omitempty"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

var _ json.Marshaler = ScheduledAlertRuleTemplate{}

func (s ScheduledAlertRuleTemplate) MarshalJSON() ([]byte, error) {
	type wrapper ScheduledAlertRuleTemplate
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ScheduledAlertRuleTemplate: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ScheduledAlertRuleTemplate: %+v", err)
	}
	decoded["kind"] = "Scheduled"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ScheduledAlertRuleTemplate: %+v", err)
	}

	return encoded, nil
}
