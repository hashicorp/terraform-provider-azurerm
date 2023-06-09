package adminrules

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ BaseAdminRule = DefaultAdminRule{}

type DefaultAdminRule struct {
	Properties *DefaultAdminPropertiesFormat `json:"properties,omitempty"`

	// Fields inherited from BaseAdminRule
	Etag       *string                `json:"etag,omitempty"`
	Id         *string                `json:"id,omitempty"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

var _ json.Marshaler = DefaultAdminRule{}

func (s DefaultAdminRule) MarshalJSON() ([]byte, error) {
	type wrapper DefaultAdminRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DefaultAdminRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DefaultAdminRule: %+v", err)
	}
	decoded["kind"] = "Default"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DefaultAdminRule: %+v", err)
	}

	return encoded, nil
}
