package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LookupActivityTypeProperties struct {
	Dataset              DatasetReference `json:"dataset"`
	FirstRowOnly         *bool            `json:"firstRowOnly,omitempty"`
	Source               CopySource       `json:"source"`
	TreatDecimalAsString *bool            `json:"treatDecimalAsString,omitempty"`
}

var _ json.Unmarshaler = &LookupActivityTypeProperties{}

func (s *LookupActivityTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Dataset              DatasetReference `json:"dataset"`
		FirstRowOnly         *bool            `json:"firstRowOnly,omitempty"`
		TreatDecimalAsString *bool            `json:"treatDecimalAsString,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Dataset = decoded.Dataset
	s.FirstRowOnly = decoded.FirstRowOnly
	s.TreatDecimalAsString = decoded.TreatDecimalAsString

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling LookupActivityTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["source"]; ok {
		impl, err := UnmarshalCopySourceImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Source' for 'LookupActivityTypeProperties': %+v", err)
		}
		s.Source = impl
	}

	return nil
}
