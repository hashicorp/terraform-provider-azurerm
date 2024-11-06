package replicationfabrics

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FabricSpecificDetails = AzureFabricSpecificDetails{}

type AzureFabricSpecificDetails struct {
	ContainerIds      *[]string                           `json:"containerIds,omitempty"`
	ExtendedLocations *[]A2AExtendedLocationDetails       `json:"extendedLocations,omitempty"`
	Location          *string                             `json:"location,omitempty"`
	LocationDetails   *[]A2AFabricSpecificLocationDetails `json:"locationDetails,omitempty"`
	Zones             *[]A2AZoneDetails                   `json:"zones,omitempty"`

	// Fields inherited from FabricSpecificDetails

	InstanceType string `json:"instanceType"`
}

func (s AzureFabricSpecificDetails) FabricSpecificDetails() BaseFabricSpecificDetailsImpl {
	return BaseFabricSpecificDetailsImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = AzureFabricSpecificDetails{}

func (s AzureFabricSpecificDetails) MarshalJSON() ([]byte, error) {
	type wrapper AzureFabricSpecificDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureFabricSpecificDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureFabricSpecificDetails: %+v", err)
	}

	decoded["instanceType"] = "Azure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureFabricSpecificDetails: %+v", err)
	}

	return encoded, nil
}
