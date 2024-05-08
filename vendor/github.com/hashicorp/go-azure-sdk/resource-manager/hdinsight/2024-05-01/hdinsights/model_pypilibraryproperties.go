package hdinsights

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ClusterLibraryProperties = PyPiLibraryProperties{}

type PyPiLibraryProperties struct {
	Name    string  `json:"name"`
	Version *string `json:"version,omitempty"`

	// Fields inherited from ClusterLibraryProperties
	Message   *string `json:"message,omitempty"`
	Remarks   *string `json:"remarks,omitempty"`
	Status    *Status `json:"status,omitempty"`
	Timestamp *string `json:"timestamp,omitempty"`
}

func (o *PyPiLibraryProperties) GetTimestampAsTime() (*time.Time, error) {
	if o.Timestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Timestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *PyPiLibraryProperties) SetTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Timestamp = &formatted
}

var _ json.Marshaler = PyPiLibraryProperties{}

func (s PyPiLibraryProperties) MarshalJSON() ([]byte, error) {
	type wrapper PyPiLibraryProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PyPiLibraryProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PyPiLibraryProperties: %+v", err)
	}
	decoded["type"] = "pypi"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PyPiLibraryProperties: %+v", err)
	}

	return encoded, nil
}
