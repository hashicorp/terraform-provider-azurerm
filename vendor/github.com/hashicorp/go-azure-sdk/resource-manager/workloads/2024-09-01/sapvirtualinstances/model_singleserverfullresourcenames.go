package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SingleServerCustomResourceNames = SingleServerFullResourceNames{}

type SingleServerFullResourceNames struct {
	VirtualMachine *VirtualMachineResourceNames `json:"virtualMachine,omitempty"`

	// Fields inherited from SingleServerCustomResourceNames

	NamingPatternType NamingPatternType `json:"namingPatternType"`
}

func (s SingleServerFullResourceNames) SingleServerCustomResourceNames() BaseSingleServerCustomResourceNamesImpl {
	return BaseSingleServerCustomResourceNamesImpl{
		NamingPatternType: s.NamingPatternType,
	}
}

var _ json.Marshaler = SingleServerFullResourceNames{}

func (s SingleServerFullResourceNames) MarshalJSON() ([]byte, error) {
	type wrapper SingleServerFullResourceNames
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SingleServerFullResourceNames: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SingleServerFullResourceNames: %+v", err)
	}

	decoded["namingPatternType"] = "FullResourceName"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SingleServerFullResourceNames: %+v", err)
	}

	return encoded, nil
}
