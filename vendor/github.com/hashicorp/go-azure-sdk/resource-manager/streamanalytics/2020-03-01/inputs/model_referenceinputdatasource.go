package inputs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReferenceInputDataSource interface {
}

func unmarshalReferenceInputDataSourceImplementation(input []byte) (ReferenceInputDataSource, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ReferenceInputDataSource into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
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

	type RawReferenceInputDataSourceImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawReferenceInputDataSourceImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
