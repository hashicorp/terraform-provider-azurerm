package appplatform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProbeAction interface {
}

// RawProbeActionImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawProbeActionImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalProbeActionImplementation(input []byte) (ProbeAction, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ProbeAction into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "ExecAction") {
		var out ExecAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ExecAction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HTTPGetAction") {
		var out HTTPGetAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HTTPGetAction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TCPSocketAction") {
		var out TCPSocketAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TCPSocketAction: %+v", err)
		}
		return out, nil
	}

	out := RawProbeActionImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
