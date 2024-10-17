package backuppolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CopyOption = ImmediateCopyOption{}

type ImmediateCopyOption struct {

	// Fields inherited from CopyOption

	ObjectType string `json:"objectType"`
}

func (s ImmediateCopyOption) CopyOption() BaseCopyOptionImpl {
	return BaseCopyOptionImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = ImmediateCopyOption{}

func (s ImmediateCopyOption) MarshalJSON() ([]byte, error) {
	type wrapper ImmediateCopyOption
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImmediateCopyOption: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImmediateCopyOption: %+v", err)
	}

	decoded["objectType"] = "ImmediateCopyOption"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImmediateCopyOption: %+v", err)
	}

	return encoded, nil
}
