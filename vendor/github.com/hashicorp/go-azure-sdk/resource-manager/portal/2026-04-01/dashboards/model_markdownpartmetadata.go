package dashboards

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DashboardPartMetadata = MarkdownPartMetadata{}

type MarkdownPartMetadata struct {
	Inputs   *[]interface{}                `json:"inputs,omitempty"`
	Settings *MarkdownPartMetadataSettings `json:"settings,omitempty"`

	// Fields inherited from DashboardPartMetadata

	Type DashboardPartMetadataType `json:"type"`
}

func (s MarkdownPartMetadata) DashboardPartMetadata() BaseDashboardPartMetadataImpl {
	return BaseDashboardPartMetadataImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = MarkdownPartMetadata{}

func (s MarkdownPartMetadata) MarshalJSON() ([]byte, error) {
	type wrapper MarkdownPartMetadata
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MarkdownPartMetadata: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MarkdownPartMetadata: %+v", err)
	}

	decoded["type"] = "Extension/HubsExtension/PartType/MarkdownPart"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MarkdownPartMetadata: %+v", err)
	}

	return encoded, nil
}
