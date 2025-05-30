package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ThreeTierCustomResourceNames = ThreeTierFullResourceNames{}

type ThreeTierFullResourceNames struct {
	ApplicationServer *ApplicationServerFullResourceNames `json:"applicationServer,omitempty"`
	CentralServer     *CentralServerFullResourceNames     `json:"centralServer,omitempty"`
	DatabaseServer    *DatabaseServerFullResourceNames    `json:"databaseServer,omitempty"`
	SharedStorage     *SharedStorageResourceNames         `json:"sharedStorage,omitempty"`

	// Fields inherited from ThreeTierCustomResourceNames

	NamingPatternType NamingPatternType `json:"namingPatternType"`
}

func (s ThreeTierFullResourceNames) ThreeTierCustomResourceNames() BaseThreeTierCustomResourceNamesImpl {
	return BaseThreeTierCustomResourceNamesImpl{
		NamingPatternType: s.NamingPatternType,
	}
}

var _ json.Marshaler = ThreeTierFullResourceNames{}

func (s ThreeTierFullResourceNames) MarshalJSON() ([]byte, error) {
	type wrapper ThreeTierFullResourceNames
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ThreeTierFullResourceNames: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ThreeTierFullResourceNames: %+v", err)
	}

	decoded["namingPatternType"] = "FullResourceName"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ThreeTierFullResourceNames: %+v", err)
	}

	return encoded, nil
}
