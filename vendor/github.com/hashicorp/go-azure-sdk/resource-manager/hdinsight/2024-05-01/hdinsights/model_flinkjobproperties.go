package hdinsights

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ClusterJobProperties = FlinkJobProperties{}

type FlinkJobProperties struct {
	Action             *Action            `json:"action,omitempty"`
	ActionResult       *string            `json:"actionResult,omitempty"`
	Args               *string            `json:"args,omitempty"`
	EntryClass         *string            `json:"entryClass,omitempty"`
	FlinkConfiguration *map[string]string `json:"flinkConfiguration,omitempty"`
	JarName            *string            `json:"jarName,omitempty"`
	JobId              *string            `json:"jobId,omitempty"`
	JobJarDirectory    *string            `json:"jobJarDirectory,omitempty"`
	JobName            *string            `json:"jobName,omitempty"`
	JobOutput          *string            `json:"jobOutput,omitempty"`
	LastSavePoint      *string            `json:"lastSavePoint,omitempty"`
	RunId              *string            `json:"runId,omitempty"`
	SavePointName      *string            `json:"savePointName,omitempty"`
	Status             *string            `json:"status,omitempty"`

	// Fields inherited from ClusterJobProperties
}

var _ json.Marshaler = FlinkJobProperties{}

func (s FlinkJobProperties) MarshalJSON() ([]byte, error) {
	type wrapper FlinkJobProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FlinkJobProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FlinkJobProperties: %+v", err)
	}
	decoded["jobType"] = "FlinkJob"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FlinkJobProperties: %+v", err)
	}

	return encoded, nil
}
