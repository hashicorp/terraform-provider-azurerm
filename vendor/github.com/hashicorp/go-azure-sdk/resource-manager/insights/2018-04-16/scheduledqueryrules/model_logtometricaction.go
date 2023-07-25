package scheduledqueryrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Action = LogToMetricAction{}

type LogToMetricAction struct {
	Criteria []Criteria `json:"criteria"`

	// Fields inherited from Action
}

var _ json.Marshaler = LogToMetricAction{}

func (s LogToMetricAction) MarshalJSON() ([]byte, error) {
	type wrapper LogToMetricAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling LogToMetricAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling LogToMetricAction: %+v", err)
	}
	decoded["odata.type"] = "Microsoft.WindowsAzure.Management.Monitoring.Alerts.Models.Microsoft.AppInsights.Nexus.DataContracts.Resources.ScheduledQueryRules.LogToMetricAction"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling LogToMetricAction: %+v", err)
	}

	return encoded, nil
}
