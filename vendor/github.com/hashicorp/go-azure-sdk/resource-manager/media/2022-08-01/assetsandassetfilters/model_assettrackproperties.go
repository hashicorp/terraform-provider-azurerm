package assetsandassetfilters

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetTrackProperties struct {
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	Track             TrackBase          `json:"track"`
}

var _ json.Unmarshaler = &AssetTrackProperties{}

func (s *AssetTrackProperties) UnmarshalJSON(bytes []byte) error {
	type alias AssetTrackProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into AssetTrackProperties: %+v", err)
	}

	s.ProvisioningState = decoded.ProvisioningState

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AssetTrackProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["track"]; ok {
		impl, err := unmarshalTrackBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Track' for 'AssetTrackProperties': %+v", err)
		}
		s.Track = impl
	}
	return nil
}
