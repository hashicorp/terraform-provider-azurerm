package afddomains

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AFDDomainMtlsParameters = CompleteMtlsPassthroughToOriginAdvancedSettings{}

type CompleteMtlsPassthroughToOriginAdvancedSettings struct {

	// Fields inherited from AFDDomainMtlsParameters

	Scenario MtlsScenarioType `json:"scenario"`
}

func (s CompleteMtlsPassthroughToOriginAdvancedSettings) AFDDomainMtlsParameters() BaseAFDDomainMtlsParametersImpl {
	return BaseAFDDomainMtlsParametersImpl{
		Scenario: s.Scenario,
	}
}

var _ json.Marshaler = CompleteMtlsPassthroughToOriginAdvancedSettings{}

func (s CompleteMtlsPassthroughToOriginAdvancedSettings) MarshalJSON() ([]byte, error) {
	type wrapper CompleteMtlsPassthroughToOriginAdvancedSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CompleteMtlsPassthroughToOriginAdvancedSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CompleteMtlsPassthroughToOriginAdvancedSettings: %+v", err)
	}

	decoded["scenario"] = "CompleteMtlsPassthroughToOrigin"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CompleteMtlsPassthroughToOriginAdvancedSettings: %+v", err)
	}

	return encoded, nil
}
