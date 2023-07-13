package protectionpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProtectionPolicy = GenericProtectionPolicy{}

type GenericProtectionPolicy struct {
	FabricName          *string                `json:"fabricName,omitempty"`
	SubProtectionPolicy *[]SubProtectionPolicy `json:"subProtectionPolicy,omitempty"`
	TimeZone            *string                `json:"timeZone,omitempty"`

	// Fields inherited from ProtectionPolicy
	ProtectedItemsCount            *int64    `json:"protectedItemsCount,omitempty"`
	ResourceGuardOperationRequests *[]string `json:"resourceGuardOperationRequests,omitempty"`
}

var _ json.Marshaler = GenericProtectionPolicy{}

func (s GenericProtectionPolicy) MarshalJSON() ([]byte, error) {
	type wrapper GenericProtectionPolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling GenericProtectionPolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling GenericProtectionPolicy: %+v", err)
	}
	decoded["backupManagementType"] = "GenericProtectionPolicy"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling GenericProtectionPolicy: %+v", err)
	}

	return encoded, nil
}
