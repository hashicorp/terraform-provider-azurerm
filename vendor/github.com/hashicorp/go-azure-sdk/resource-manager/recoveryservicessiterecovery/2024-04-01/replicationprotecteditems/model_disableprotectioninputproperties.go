package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DisableProtectionInputProperties struct {
	DisableProtectionReason  *DisableProtectionReason               `json:"disableProtectionReason,omitempty"`
	ReplicationProviderInput DisableProtectionProviderSpecificInput `json:"replicationProviderInput"`
}

var _ json.Unmarshaler = &DisableProtectionInputProperties{}

func (s *DisableProtectionInputProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		DisableProtectionReason *DisableProtectionReason `json:"disableProtectionReason,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.DisableProtectionReason = decoded.DisableProtectionReason

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling DisableProtectionInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["replicationProviderInput"]; ok {
		impl, err := UnmarshalDisableProtectionProviderSpecificInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ReplicationProviderInput' for 'DisableProtectionInputProperties': %+v", err)
		}
		s.ReplicationProviderInput = impl
	}

	return nil
}
