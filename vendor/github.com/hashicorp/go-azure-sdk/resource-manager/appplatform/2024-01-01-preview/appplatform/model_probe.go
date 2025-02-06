package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Probe struct {
	DisableProbe        bool        `json:"disableProbe"`
	FailureThreshold    *int64      `json:"failureThreshold,omitempty"`
	InitialDelaySeconds *int64      `json:"initialDelaySeconds,omitempty"`
	PeriodSeconds       *int64      `json:"periodSeconds,omitempty"`
	ProbeAction         ProbeAction `json:"probeAction"`
	SuccessThreshold    *int64      `json:"successThreshold,omitempty"`
	TimeoutSeconds      *int64      `json:"timeoutSeconds,omitempty"`
}

var _ json.Unmarshaler = &Probe{}

func (s *Probe) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		DisableProbe        bool   `json:"disableProbe"`
		FailureThreshold    *int64 `json:"failureThreshold,omitempty"`
		InitialDelaySeconds *int64 `json:"initialDelaySeconds,omitempty"`
		PeriodSeconds       *int64 `json:"periodSeconds,omitempty"`
		SuccessThreshold    *int64 `json:"successThreshold,omitempty"`
		TimeoutSeconds      *int64 `json:"timeoutSeconds,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.DisableProbe = decoded.DisableProbe
	s.FailureThreshold = decoded.FailureThreshold
	s.InitialDelaySeconds = decoded.InitialDelaySeconds
	s.PeriodSeconds = decoded.PeriodSeconds
	s.SuccessThreshold = decoded.SuccessThreshold
	s.TimeoutSeconds = decoded.TimeoutSeconds

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling Probe into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["probeAction"]; ok {
		impl, err := UnmarshalProbeActionImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProbeAction' for 'Probe': %+v", err)
		}
		s.ProbeAction = impl
	}

	return nil
}
