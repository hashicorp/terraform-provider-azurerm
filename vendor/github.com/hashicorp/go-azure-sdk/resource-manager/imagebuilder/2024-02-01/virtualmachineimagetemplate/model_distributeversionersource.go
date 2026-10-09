package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DistributeVersioner = DistributeVersionerSource{}

type DistributeVersionerSource struct {

	// Fields inherited from DistributeVersioner

	Scheme string `json:"scheme"`
}

func (s DistributeVersionerSource) DistributeVersioner() BaseDistributeVersionerImpl {
	return BaseDistributeVersionerImpl{
		Scheme: s.Scheme,
	}
}

var _ json.Marshaler = DistributeVersionerSource{}

func (s DistributeVersionerSource) MarshalJSON() ([]byte, error) {
	type wrapper DistributeVersionerSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DistributeVersionerSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DistributeVersionerSource: %+v", err)
	}

	decoded["scheme"] = "Source"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DistributeVersionerSource: %+v", err)
	}

	return encoded, nil
}
