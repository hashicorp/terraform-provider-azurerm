package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Filters struct {
	Crop        *Rectangle   `json:"crop,omitempty"`
	Deinterlace *Deinterlace `json:"deinterlace,omitempty"`
	FadeIn      *Fade        `json:"fadeIn,omitempty"`
	FadeOut     *Fade        `json:"fadeOut,omitempty"`
	Overlays    *[]Overlay   `json:"overlays,omitempty"`
	Rotation    *Rotation    `json:"rotation,omitempty"`
}

var _ json.Unmarshaler = &Filters{}

func (s *Filters) UnmarshalJSON(bytes []byte) error {
	type alias Filters
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into Filters: %+v", err)
	}

	s.Crop = decoded.Crop
	s.Deinterlace = decoded.Deinterlace
	s.FadeIn = decoded.FadeIn
	s.FadeOut = decoded.FadeOut
	s.Rotation = decoded.Rotation

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling Filters into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["overlays"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Overlays into list []json.RawMessage: %+v", err)
		}

		output := make([]Overlay, 0)
		for i, val := range listTemp {
			impl, err := unmarshalOverlayImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Overlays' for 'Filters': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Overlays = &output
	}
	return nil
}
