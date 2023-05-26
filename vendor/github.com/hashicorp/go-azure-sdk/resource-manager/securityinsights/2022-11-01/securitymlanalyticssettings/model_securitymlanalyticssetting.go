package securitymlanalyticssettings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityMLAnalyticsSetting interface {
}

func unmarshalSecurityMLAnalyticsSettingImplementation(input []byte) (SecurityMLAnalyticsSetting, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SecurityMLAnalyticsSetting into map[string]interface: %+v", err)
	}

	value, ok := temp["kind"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Anomaly") {
		var out AnomalySecurityMLAnalyticsSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AnomalySecurityMLAnalyticsSettings: %+v", err)
		}
		return out, nil
	}

	type RawSecurityMLAnalyticsSettingImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawSecurityMLAnalyticsSettingImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
