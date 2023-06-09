package networkmanageractiveconfigurations

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActiveSecurityAdminRulesListResult struct {
	SkipToken *string                        `json:"skipToken,omitempty"`
	Value     *[]ActiveBaseSecurityAdminRule `json:"value,omitempty"`
}

var _ json.Unmarshaler = &ActiveSecurityAdminRulesListResult{}

func (s *ActiveSecurityAdminRulesListResult) UnmarshalJSON(bytes []byte) error {
	type alias ActiveSecurityAdminRulesListResult
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ActiveSecurityAdminRulesListResult: %+v", err)
	}

	s.SkipToken = decoded.SkipToken

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ActiveSecurityAdminRulesListResult into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["value"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Value into list []json.RawMessage: %+v", err)
		}

		output := make([]ActiveBaseSecurityAdminRule, 0)
		for i, val := range listTemp {
			impl, err := unmarshalActiveBaseSecurityAdminRuleImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Value' for 'ActiveSecurityAdminRulesListResult': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Value = &output
	}
	return nil
}
