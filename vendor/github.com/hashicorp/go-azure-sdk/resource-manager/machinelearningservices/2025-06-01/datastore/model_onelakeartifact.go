package datastore

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OneLakeArtifact interface {
	OneLakeArtifact() BaseOneLakeArtifactImpl
}

var _ OneLakeArtifact = BaseOneLakeArtifactImpl{}

type BaseOneLakeArtifactImpl struct {
	ArtifactName string              `json:"artifactName"`
	ArtifactType OneLakeArtifactType `json:"artifactType"`
}

func (s BaseOneLakeArtifactImpl) OneLakeArtifact() BaseOneLakeArtifactImpl {
	return s
}

var _ OneLakeArtifact = RawOneLakeArtifactImpl{}

// RawOneLakeArtifactImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawOneLakeArtifactImpl struct {
	oneLakeArtifact BaseOneLakeArtifactImpl
	Type            string
	Values          map[string]interface{}
}

func (s RawOneLakeArtifactImpl) OneLakeArtifact() BaseOneLakeArtifactImpl {
	return s.oneLakeArtifact
}

func UnmarshalOneLakeArtifactImplementation(input []byte) (OneLakeArtifact, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling OneLakeArtifact into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["artifactType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "LakeHouse") {
		var out LakeHouseArtifact
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LakeHouseArtifact: %+v", err)
		}
		return out, nil
	}

	var parent BaseOneLakeArtifactImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseOneLakeArtifactImpl: %+v", err)
	}

	return RawOneLakeArtifactImpl{
		oneLakeArtifact: parent,
		Type:            value,
		Values:          temp,
	}, nil

}
