package serviceresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateSsisTaskOutput interface {
	MigrateSsisTaskOutput() BaseMigrateSsisTaskOutputImpl
}

var _ MigrateSsisTaskOutput = BaseMigrateSsisTaskOutputImpl{}

type BaseMigrateSsisTaskOutputImpl struct {
	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s BaseMigrateSsisTaskOutputImpl) MigrateSsisTaskOutput() BaseMigrateSsisTaskOutputImpl {
	return s
}

var _ MigrateSsisTaskOutput = RawMigrateSsisTaskOutputImpl{}

// RawMigrateSsisTaskOutputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawMigrateSsisTaskOutputImpl struct {
	migrateSsisTaskOutput BaseMigrateSsisTaskOutputImpl
	Type                  string
	Values                map[string]interface{}
}

func (s RawMigrateSsisTaskOutputImpl) MigrateSsisTaskOutput() BaseMigrateSsisTaskOutputImpl {
	return s.migrateSsisTaskOutput
}

func UnmarshalMigrateSsisTaskOutputImplementation(input []byte) (MigrateSsisTaskOutput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSsisTaskOutput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["resultType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "MigrationLevelOutput") {
		var out MigrateSsisTaskOutputMigrationLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSsisTaskOutputMigrationLevel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SsisProjectLevelOutput") {
		var out MigrateSsisTaskOutputProjectLevel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MigrateSsisTaskOutputProjectLevel: %+v", err)
		}
		return out, nil
	}

	var parent BaseMigrateSsisTaskOutputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseMigrateSsisTaskOutputImpl: %+v", err)
	}

	return RawMigrateSsisTaskOutputImpl{
		migrateSsisTaskOutput: parent,
		Type:                  value,
		Values:                temp,
	}, nil

}
