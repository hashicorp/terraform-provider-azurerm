package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TransformOutput struct {
	OnError          *OnErrorType `json:"onError,omitempty"`
	Preset           Preset       `json:"preset"`
	RelativePriority *Priority    `json:"relativePriority,omitempty"`
}

var _ json.Unmarshaler = &TransformOutput{}

func (s *TransformOutput) UnmarshalJSON(bytes []byte) error {
	type alias TransformOutput
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into TransformOutput: %+v", err)
	}

	s.OnError = decoded.OnError
	s.RelativePriority = decoded.RelativePriority

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling TransformOutput into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["preset"]; ok {
		impl, err := unmarshalPresetImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Preset' for 'TransformOutput': %+v", err)
		}
		s.Preset = impl
	}
	return nil
}
