package dataversionregistry

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataVersionBase interface {
}

func unmarshalDataVersionBaseImplementation(input []byte) (DataVersionBase, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DataVersionBase into map[string]interface: %+v", err)
	}

	value, ok := temp["dataType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "mltable") {
		var out MLTableData
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MLTableData: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "uri_file") {
		var out UriFileDataVersion
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into UriFileDataVersion: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "uri_folder") {
		var out UriFolderDataVersion
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into UriFolderDataVersion: %+v", err)
		}
		return out, nil
	}

	type RawDataVersionBaseImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawDataVersionBaseImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
