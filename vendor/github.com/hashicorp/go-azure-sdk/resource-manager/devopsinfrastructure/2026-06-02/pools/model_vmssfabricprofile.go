package pools

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FabricProfile = VMSSFabricProfile{}

type VMSSFabricProfile struct {
	Images         []PoolImage     `json:"images"`
	NetworkProfile *NetworkProfile `json:"networkProfile,omitempty"`
	OsProfile      *OsProfile      `json:"osProfile,omitempty"`
	Sku            DevOpsAzureSku  `json:"sku"`
	StorageProfile *StorageProfile `json:"storageProfile,omitempty"`

	// Fields inherited from FabricProfile

	Kind string `json:"kind"`
}

func (s VMSSFabricProfile) FabricProfile() BaseFabricProfileImpl {
	return BaseFabricProfileImpl{
		Kind: s.Kind,
	}
}

var _ json.Marshaler = VMSSFabricProfile{}

func (s VMSSFabricProfile) MarshalJSON() ([]byte, error) {
	type wrapper VMSSFabricProfile
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VMSSFabricProfile: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VMSSFabricProfile: %+v", err)
	}

	decoded["kind"] = "Vmss"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VMSSFabricProfile: %+v", err)
	}

	return encoded, nil
}
