package replicationnetworkmappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FabricSpecificCreateNetworkMappingInput = VMmToVMmCreateNetworkMappingInput{}

type VMmToVMmCreateNetworkMappingInput struct {

	// Fields inherited from FabricSpecificCreateNetworkMappingInput

	InstanceType string `json:"instanceType"`
}

func (s VMmToVMmCreateNetworkMappingInput) FabricSpecificCreateNetworkMappingInput() BaseFabricSpecificCreateNetworkMappingInputImpl {
	return BaseFabricSpecificCreateNetworkMappingInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = VMmToVMmCreateNetworkMappingInput{}

func (s VMmToVMmCreateNetworkMappingInput) MarshalJSON() ([]byte, error) {
	type wrapper VMmToVMmCreateNetworkMappingInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VMmToVMmCreateNetworkMappingInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VMmToVMmCreateNetworkMappingInput: %+v", err)
	}

	decoded["instanceType"] = "VmmToVmm"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VMmToVMmCreateNetworkMappingInput: %+v", err)
	}

	return encoded, nil
}
