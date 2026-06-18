package integrationruntime

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CmdkeySetupTypeProperties struct {
	Password   SecretBase  `json:"password"`
	TargetName interface{} `json:"targetName"`
	UserName   interface{} `json:"userName"`
}

var _ json.Unmarshaler = &CmdkeySetupTypeProperties{}

func (s *CmdkeySetupTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		TargetName interface{} `json:"targetName"`
		UserName   interface{} `json:"userName"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.TargetName = decoded.TargetName
	s.UserName = decoded.UserName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling CmdkeySetupTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'CmdkeySetupTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
