package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppFiguresLinkedServiceTypeProperties struct {
	ClientKey SecretBase  `json:"clientKey"`
	Password  SecretBase  `json:"password"`
	UserName  interface{} `json:"userName"`
}

var _ json.Unmarshaler = &AppFiguresLinkedServiceTypeProperties{}

func (s *AppFiguresLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		UserName interface{} `json:"userName"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.UserName = decoded.UserName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AppFiguresLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["clientKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ClientKey' for 'AppFiguresLinkedServiceTypeProperties': %+v", err)
		}
		s.ClientKey = impl
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'AppFiguresLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
