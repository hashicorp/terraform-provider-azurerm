package appplatform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProbeAction interface {
	ProbeAction() BaseProbeActionImpl
}

var _ ProbeAction = BaseProbeActionImpl{}

type BaseProbeActionImpl struct {
	Type ProbeActionType `json:"type"`
}

func (s BaseProbeActionImpl) ProbeAction() BaseProbeActionImpl {
	return s
}

var _ ProbeAction = RawProbeActionImpl{}

// RawProbeActionImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawProbeActionImpl struct {
	probeAction BaseProbeActionImpl
	Type        string
	Values      map[string]interface{}
}

func (s RawProbeActionImpl) ProbeAction() BaseProbeActionImpl {
	return s.probeAction
}

func (s RawProbeActionImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalProbeActionImplementation(input []byte) (ProbeAction, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ProbeAction into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseProbeActionImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseProbeActionImpl: %+v", err)
	}

	return RawProbeActionImpl{
		probeAction: parent,
		Type:        value,
		Values:      temp,
	}, nil

}
