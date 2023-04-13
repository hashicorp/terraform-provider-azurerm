package timeseriesdatabaseconnections

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TimeSeriesDatabaseConnectionProperties interface {
}

func unmarshalTimeSeriesDatabaseConnectionPropertiesImplementation(input []byte) (TimeSeriesDatabaseConnectionProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling TimeSeriesDatabaseConnectionProperties into map[string]interface: %+v", err)
	}

	value, ok := temp["connectionType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AzureDataExplorer") {
		var out AzureDataExplorerConnectionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDataExplorerConnectionProperties: %+v", err)
		}
		return out, nil
	}

	type RawTimeSeriesDatabaseConnectionPropertiesImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawTimeSeriesDatabaseConnectionPropertiesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
