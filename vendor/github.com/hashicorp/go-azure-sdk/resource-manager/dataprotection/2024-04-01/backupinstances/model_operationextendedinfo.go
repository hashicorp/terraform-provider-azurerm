package backupinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperationExtendedInfo interface {
	OperationExtendedInfo() BaseOperationExtendedInfoImpl
}

var _ OperationExtendedInfo = BaseOperationExtendedInfoImpl{}

type BaseOperationExtendedInfoImpl struct {
	ObjectType string `json:"objectType"`
}

func (s BaseOperationExtendedInfoImpl) OperationExtendedInfo() BaseOperationExtendedInfoImpl {
	return s
}

var _ OperationExtendedInfo = RawOperationExtendedInfoImpl{}

// RawOperationExtendedInfoImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawOperationExtendedInfoImpl struct {
	operationExtendedInfo BaseOperationExtendedInfoImpl
	Type                  string
	Values                map[string]interface{}
}

func (s RawOperationExtendedInfoImpl) OperationExtendedInfo() BaseOperationExtendedInfoImpl {
	return s.operationExtendedInfo
}

func UnmarshalOperationExtendedInfoImplementation(input []byte) (OperationExtendedInfo, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling OperationExtendedInfo into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["objectType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "OperationJobExtendedInfo") {
		var out OperationJobExtendedInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OperationJobExtendedInfo: %+v", err)
		}
		return out, nil
	}

	var parent BaseOperationExtendedInfoImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseOperationExtendedInfoImpl: %+v", err)
	}

	return RawOperationExtendedInfoImpl{
		operationExtendedInfo: parent,
		Type:                  value,
		Values:                temp,
	}, nil

}
