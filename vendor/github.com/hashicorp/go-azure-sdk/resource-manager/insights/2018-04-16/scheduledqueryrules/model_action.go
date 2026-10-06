package scheduledqueryrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Action interface {
	Action() BaseActionImpl
}

var _ Action = BaseActionImpl{}

type BaseActionImpl struct {
	OdataType string `json:"odata.type"`
}

func (s BaseActionImpl) Action() BaseActionImpl {
	return s
}

var _ Action = RawActionImpl{}

// RawActionImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawActionImpl struct {
	action BaseActionImpl
	Type   string
	Values map[string]interface{}
}

func (s RawActionImpl) Action() BaseActionImpl {
	return s.action
}

func (s RawActionImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalActionImplementation(input []byte) (Action, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Action into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["odata.type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Microsoft.WindowsAzure.Management.Monitoring.Alerts.Models.Microsoft.AppInsights.Nexus.DataContracts.Resources.ScheduledQueryRules.AlertingAction") {
		var out AlertingAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AlertingAction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.WindowsAzure.Management.Monitoring.Alerts.Models.Microsoft.AppInsights.Nexus.DataContracts.Resources.ScheduledQueryRules.LogToMetricAction") {
		var out LogToMetricAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LogToMetricAction: %+v", err)
		}
		return out, nil
	}

	var parent BaseActionImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseActionImpl: %+v", err)
	}

	return RawActionImpl{
		action: parent,
		Type:   value,
		Values: temp,
	}, nil

}
