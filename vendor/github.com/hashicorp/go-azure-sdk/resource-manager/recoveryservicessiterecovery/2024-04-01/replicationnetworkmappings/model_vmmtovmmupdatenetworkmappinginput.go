package replicationnetworkmappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FabricSpecificUpdateNetworkMappingInput = VMmToVMmUpdateNetworkMappingInput{}

type VMmToVMmUpdateNetworkMappingInput struct {

	// Fields inherited from FabricSpecificUpdateNetworkMappingInput

	InstanceType string `json:"instanceType"`
}

func (s VMmToVMmUpdateNetworkMappingInput) FabricSpecificUpdateNetworkMappingInput() BaseFabricSpecificUpdateNetworkMappingInputImpl {
	return BaseFabricSpecificUpdateNetworkMappingInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = VMmToVMmUpdateNetworkMappingInput{}

func (s VMmToVMmUpdateNetworkMappingInput) MarshalJSON() ([]byte, error) {
	type wrapper VMmToVMmUpdateNetworkMappingInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VMmToVMmUpdateNetworkMappingInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VMmToVMmUpdateNetworkMappingInput: %+v", err)
	}

	decoded["instanceType"] = "VmmToVmm"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VMmToVMmUpdateNetworkMappingInput: %+v", err)
	}

	return encoded, nil
}
