package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OSProfile struct {
	AdminPassword   *string         `json:"adminPassword,omitempty"`
	AdminUsername   *string         `json:"adminUsername,omitempty"`
	OsConfiguration OSConfiguration `json:"osConfiguration"`
}

var _ json.Unmarshaler = &OSProfile{}

func (s *OSProfile) UnmarshalJSON(bytes []byte) error {
	type alias OSProfile
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into OSProfile: %+v", err)
	}

	s.AdminPassword = decoded.AdminPassword
	s.AdminUsername = decoded.AdminUsername

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling OSProfile into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["osConfiguration"]; ok {
		impl, err := unmarshalOSConfigurationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'OsConfiguration' for 'OSProfile': %+v", err)
		}
		s.OsConfiguration = impl
	}
	return nil
}
