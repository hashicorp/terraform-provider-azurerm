package inputs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReferenceInputDataSource interface {
	ReferenceInputDataSource() BaseReferenceInputDataSourceImpl
}

var _ ReferenceInputDataSource = BaseReferenceInputDataSourceImpl{}

type BaseReferenceInputDataSourceImpl struct {
	Type string `json:"type"`
}

func (s BaseReferenceInputDataSourceImpl) ReferenceInputDataSource() BaseReferenceInputDataSourceImpl {
	return s
}

var _ ReferenceInputDataSource = RawReferenceInputDataSourceImpl{}

// RawReferenceInputDataSourceImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawReferenceInputDataSourceImpl struct {
	referenceInputDataSource BaseReferenceInputDataSourceImpl
	Type                     string
	Values                   map[string]interface{}
}

func (s RawReferenceInputDataSourceImpl) ReferenceInputDataSource() BaseReferenceInputDataSourceImpl {
	return s.referenceInputDataSource
}

func UnmarshalReferenceInputDataSourceImplementation(input []byte) (ReferenceInputDataSource, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ReferenceInputDataSource into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Microsoft.Sql/Server/Database") {
		var out AzureSqlReferenceInputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSqlReferenceInputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.Storage/Blob") {
		var out BlobReferenceInputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BlobReferenceInputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "File") {
		var out FileReferenceInputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FileReferenceInputDataSource: %+v", err)
		}
		return out, nil
	}

	var parent BaseReferenceInputDataSourceImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseReferenceInputDataSourceImpl: %+v", err)
	}

	return RawReferenceInputDataSourceImpl{
		referenceInputDataSource: parent,
		Type:                     value,
		Values:                   temp,
	}, nil

}
