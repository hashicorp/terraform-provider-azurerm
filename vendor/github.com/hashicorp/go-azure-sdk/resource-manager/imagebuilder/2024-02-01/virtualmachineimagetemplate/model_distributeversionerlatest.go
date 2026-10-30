package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DistributeVersioner = DistributeVersionerLatest{}

type DistributeVersionerLatest struct {
	Major *int64 `json:"major,omitempty"`

	// Fields inherited from DistributeVersioner

	Scheme string `json:"scheme"`
}

func (s DistributeVersionerLatest) DistributeVersioner() BaseDistributeVersionerImpl {
	return BaseDistributeVersionerImpl{
		Scheme: s.Scheme,
	}
}

var _ json.Marshaler = DistributeVersionerLatest{}

func (s DistributeVersionerLatest) MarshalJSON() ([]byte, error) {
	type wrapper DistributeVersionerLatest
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DistributeVersionerLatest: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DistributeVersionerLatest: %+v", err)
	}

	decoded["scheme"] = "Latest"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DistributeVersionerLatest: %+v", err)
	}

	return encoded, nil
}
