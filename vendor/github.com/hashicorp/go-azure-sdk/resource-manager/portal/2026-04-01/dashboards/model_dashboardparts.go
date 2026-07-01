package dashboards

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DashboardParts struct {
	Metadata DashboardPartMetadata  `json:"metadata"`
	Position DashboardPartsPosition `json:"position"`
}

var _ json.Unmarshaler = &DashboardParts{}

func (s *DashboardParts) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Position DashboardPartsPosition `json:"position"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Position = decoded.Position

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling DashboardParts into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["metadata"]; ok {
		impl, err := UnmarshalDashboardPartMetadataImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Metadata' for 'DashboardParts': %+v", err)
		}
		s.Metadata = impl
	}

	return nil
}
