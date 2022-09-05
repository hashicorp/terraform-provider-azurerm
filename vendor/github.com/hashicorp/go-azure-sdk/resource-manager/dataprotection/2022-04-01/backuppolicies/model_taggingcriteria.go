package backuppolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TaggingCriteria struct {
	Criteria        *[]BackupCriteria `json:"criteria,omitempty"`
	IsDefault       bool              `json:"isDefault"`
	TagInfo         RetentionTag      `json:"tagInfo"`
	TaggingPriority int64             `json:"taggingPriority"`
}

var _ json.Unmarshaler = &TaggingCriteria{}

func (s *TaggingCriteria) UnmarshalJSON(bytes []byte) error {
	type alias TaggingCriteria
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into TaggingCriteria: %+v", err)
	}

	s.IsDefault = decoded.IsDefault
	s.TagInfo = decoded.TagInfo
	s.TaggingPriority = decoded.TaggingPriority

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling TaggingCriteria into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["criteria"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Criteria into list []json.RawMessage: %+v", err)
		}

		output := make([]BackupCriteria, 0)
		for i, val := range listTemp {
			impl, err := unmarshalBackupCriteriaImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Criteria' for 'TaggingCriteria': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Criteria = &output
	}
	return nil
}
