package threatintelligence

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ThreatIntelligenceInformation interface {
	ThreatIntelligenceInformation() BaseThreatIntelligenceInformationImpl
}

var _ ThreatIntelligenceInformation = BaseThreatIntelligenceInformationImpl{}

type BaseThreatIntelligenceInformationImpl struct {
	Etag       *string                            `json:"etag,omitempty"`
	Id         *string                            `json:"id,omitempty"`
	Kind       ThreatIntelligenceResourceKindEnum `json:"kind"`
	Name       *string                            `json:"name,omitempty"`
	SystemData *systemdata.SystemData             `json:"systemData,omitempty"`
	Type       *string                            `json:"type,omitempty"`
}

func (s BaseThreatIntelligenceInformationImpl) ThreatIntelligenceInformation() BaseThreatIntelligenceInformationImpl {
	return s
}

var _ ThreatIntelligenceInformation = RawThreatIntelligenceInformationImpl{}

// RawThreatIntelligenceInformationImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawThreatIntelligenceInformationImpl struct {
	threatIntelligenceInformation BaseThreatIntelligenceInformationImpl
	Type                          string
	Values                        map[string]interface{}
}

func (s RawThreatIntelligenceInformationImpl) ThreatIntelligenceInformation() BaseThreatIntelligenceInformationImpl {
	return s.threatIntelligenceInformation
}

func (s RawThreatIntelligenceInformationImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalThreatIntelligenceInformationImplementation(input []byte) (ThreatIntelligenceInformation, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ThreatIntelligenceInformation into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "indicator") {
		var out ThreatIntelligenceIndicatorModel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ThreatIntelligenceIndicatorModel: %+v", err)
		}
		return out, nil
	}

	var parent BaseThreatIntelligenceInformationImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseThreatIntelligenceInformationImpl: %+v", err)
	}

	return RawThreatIntelligenceInformationImpl{
		threatIntelligenceInformation: parent,
		Type:                          value,
		Values:                        temp,
	}, nil

}
