package operationprogress

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperationProgressResponseType interface {
	OperationProgressResponseType() BaseOperationProgressResponseTypeImpl
}

var _ OperationProgressResponseType = BaseOperationProgressResponseTypeImpl{}

type BaseOperationProgressResponseTypeImpl struct {
	ObjectType ObjectType `json:"objectType"`
}

func (s BaseOperationProgressResponseTypeImpl) OperationProgressResponseType() BaseOperationProgressResponseTypeImpl {
	return s
}

var _ OperationProgressResponseType = RawOperationProgressResponseTypeImpl{}

// RawOperationProgressResponseTypeImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawOperationProgressResponseTypeImpl struct {
	operationProgressResponseType BaseOperationProgressResponseTypeImpl
	Type                          string
	Values                        map[string]interface{}
}

func (s RawOperationProgressResponseTypeImpl) OperationProgressResponseType() BaseOperationProgressResponseTypeImpl {
	return s.operationProgressResponseType
}

func UnmarshalOperationProgressResponseTypeImplementation(input []byte) (OperationProgressResponseType, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling OperationProgressResponseType into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["objectType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "BackupAndExportResponse") {
		var out BackupAndExportResponseType
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BackupAndExportResponseType: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ImportFromStorageResponse") {
		var out ImportFromStorageResponseType
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImportFromStorageResponseType: %+v", err)
		}
		return out, nil
	}

	var parent BaseOperationProgressResponseTypeImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseOperationProgressResponseTypeImpl: %+v", err)
	}

	return RawOperationProgressResponseTypeImpl{
		operationProgressResponseType: parent,
		Type:                          value,
		Values:                        temp,
	}, nil

}
