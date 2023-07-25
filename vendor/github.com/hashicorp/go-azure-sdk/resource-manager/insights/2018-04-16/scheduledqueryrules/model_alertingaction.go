package scheduledqueryrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Action = AlertingAction{}

type AlertingAction struct {
	AznsAction      *AzNsActionGroup `json:"aznsAction,omitempty"`
	Severity        AlertSeverity    `json:"severity"`
	ThrottlingInMin *int64           `json:"throttlingInMin,omitempty"`
	Trigger         TriggerCondition `json:"trigger"`

	// Fields inherited from Action
}

var _ json.Marshaler = AlertingAction{}

func (s AlertingAction) MarshalJSON() ([]byte, error) {
	type wrapper AlertingAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AlertingAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AlertingAction: %+v", err)
	}
	decoded["odata.type"] = "Microsoft.WindowsAzure.Management.Monitoring.Alerts.Models.Microsoft.AppInsights.Nexus.DataContracts.Resources.ScheduledQueryRules.AlertingAction"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AlertingAction: %+v", err)
	}

	return encoded, nil
}
