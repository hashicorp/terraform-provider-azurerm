package automationrules

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationRuleProperties struct {
	Actions             []AutomationRuleAction        `json:"actions"`
	CreatedBy           *ClientInfo                   `json:"createdBy,omitempty"`
	CreatedTimeUtc      *string                       `json:"createdTimeUtc,omitempty"`
	DisplayName         string                        `json:"displayName"`
	LastModifiedBy      *ClientInfo                   `json:"lastModifiedBy,omitempty"`
	LastModifiedTimeUtc *string                       `json:"lastModifiedTimeUtc,omitempty"`
	Order               int64                         `json:"order"`
	TriggeringLogic     AutomationRuleTriggeringLogic `json:"triggeringLogic"`
}

func (o *AutomationRuleProperties) GetCreatedTimeUtcAsTime() (*time.Time, error) {
	if o.CreatedTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *AutomationRuleProperties) SetCreatedTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedTimeUtc = &formatted
}

func (o *AutomationRuleProperties) GetLastModifiedTimeUtcAsTime() (*time.Time, error) {
	if o.LastModifiedTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *AutomationRuleProperties) SetLastModifiedTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTimeUtc = &formatted
}

var _ json.Unmarshaler = &AutomationRuleProperties{}

func (s *AutomationRuleProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		CreatedBy           *ClientInfo                   `json:"createdBy,omitempty"`
		CreatedTimeUtc      *string                       `json:"createdTimeUtc,omitempty"`
		DisplayName         string                        `json:"displayName"`
		LastModifiedBy      *ClientInfo                   `json:"lastModifiedBy,omitempty"`
		LastModifiedTimeUtc *string                       `json:"lastModifiedTimeUtc,omitempty"`
		Order               int64                         `json:"order"`
		TriggeringLogic     AutomationRuleTriggeringLogic `json:"triggeringLogic"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.CreatedBy = decoded.CreatedBy
	s.CreatedTimeUtc = decoded.CreatedTimeUtc
	s.DisplayName = decoded.DisplayName
	s.LastModifiedBy = decoded.LastModifiedBy
	s.LastModifiedTimeUtc = decoded.LastModifiedTimeUtc
	s.Order = decoded.Order
	s.TriggeringLogic = decoded.TriggeringLogic

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AutomationRuleProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["actions"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Actions into list []json.RawMessage: %+v", err)
		}

		output := make([]AutomationRuleAction, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalAutomationRuleActionImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Actions' for 'AutomationRuleProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Actions = output
	}

	return nil
}
