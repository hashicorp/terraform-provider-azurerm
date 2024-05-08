package hdinsights

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ClusterLibraryProperties = MavenLibraryProperties{}

type MavenLibraryProperties struct {
	GroupId string  `json:"groupId"`
	Name    string  `json:"name"`
	Version *string `json:"version,omitempty"`

	// Fields inherited from ClusterLibraryProperties
	Message   *string `json:"message,omitempty"`
	Remarks   *string `json:"remarks,omitempty"`
	Status    *Status `json:"status,omitempty"`
	Timestamp *string `json:"timestamp,omitempty"`
}

func (o *MavenLibraryProperties) GetTimestampAsTime() (*time.Time, error) {
	if o.Timestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Timestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *MavenLibraryProperties) SetTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Timestamp = &formatted
}

var _ json.Marshaler = MavenLibraryProperties{}

func (s MavenLibraryProperties) MarshalJSON() ([]byte, error) {
	type wrapper MavenLibraryProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MavenLibraryProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MavenLibraryProperties: %+v", err)
	}
	decoded["type"] = "maven"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MavenLibraryProperties: %+v", err)
	}

	return encoded, nil
}
