package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Nodes = AllNodes{}

type AllNodes struct {

	// Fields inherited from Nodes
}

var _ json.Marshaler = AllNodes{}

func (s AllNodes) MarshalJSON() ([]byte, error) {
	type wrapper AllNodes
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AllNodes: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AllNodes: %+v", err)
	}
	decoded["nodesValueType"] = "All"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AllNodes: %+v", err)
	}

	return encoded, nil
}
