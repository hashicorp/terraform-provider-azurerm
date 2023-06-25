package vaultcertificates

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VaultCertificateResponse struct {
	Id         *string                    `json:"id,omitempty"`
	Name       *string                    `json:"name,omitempty"`
	Properties ResourceCertificateDetails `json:"properties"`
	Type       *string                    `json:"type,omitempty"`
}

var _ json.Unmarshaler = &VaultCertificateResponse{}

func (s *VaultCertificateResponse) UnmarshalJSON(bytes []byte) error {
	type alias VaultCertificateResponse
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into VaultCertificateResponse: %+v", err)
	}

	s.Id = decoded.Id
	s.Name = decoded.Name
	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling VaultCertificateResponse into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := unmarshalResourceCertificateDetailsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'VaultCertificateResponse': %+v", err)
		}
		s.Properties = impl
	}
	return nil
}
