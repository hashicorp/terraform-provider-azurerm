package customdomains

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomDomainProperties struct {
	CustomHTTPSParameters           CustomDomainHTTPSParameters      `json:"customHttpsParameters"`
	CustomHTTPSProvisioningState    *CustomHTTPSProvisioningState    `json:"customHttpsProvisioningState,omitempty"`
	CustomHTTPSProvisioningSubstate *CustomHTTPSProvisioningSubstate `json:"customHttpsProvisioningSubstate,omitempty"`
	HostName                        string                           `json:"hostName"`
	ProvisioningState               *CustomHTTPSProvisioningState    `json:"provisioningState,omitempty"`
	ResourceState                   *CustomDomainResourceState       `json:"resourceState,omitempty"`
	ValidationData                  *string                          `json:"validationData,omitempty"`
}

var _ json.Unmarshaler = &CustomDomainProperties{}

func (s *CustomDomainProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		CustomHTTPSProvisioningState    *CustomHTTPSProvisioningState    `json:"customHttpsProvisioningState,omitempty"`
		CustomHTTPSProvisioningSubstate *CustomHTTPSProvisioningSubstate `json:"customHttpsProvisioningSubstate,omitempty"`
		HostName                        string                           `json:"hostName"`
		ProvisioningState               *CustomHTTPSProvisioningState    `json:"provisioningState,omitempty"`
		ResourceState                   *CustomDomainResourceState       `json:"resourceState,omitempty"`
		ValidationData                  *string                          `json:"validationData,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.CustomHTTPSProvisioningState = decoded.CustomHTTPSProvisioningState
	s.CustomHTTPSProvisioningSubstate = decoded.CustomHTTPSProvisioningSubstate
	s.HostName = decoded.HostName
	s.ProvisioningState = decoded.ProvisioningState
	s.ResourceState = decoded.ResourceState
	s.ValidationData = decoded.ValidationData

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling CustomDomainProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["customHttpsParameters"]; ok {
		impl, err := UnmarshalCustomDomainHTTPSParametersImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'CustomHTTPSParameters' for 'CustomDomainProperties': %+v", err)
		}
		s.CustomHTTPSParameters = impl
	}

	return nil
}
