package job

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Nodes interface {
}

func unmarshalNodesImplementation(input []byte) (Nodes, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Nodes into map[string]interface: %+v", err)
	}

	value, ok := temp["nodesValueType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "All") {
		var out AllNodes
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AllNodes: %+v", err)
		}
		return out, nil
	}

	type RawNodesImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawNodesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
