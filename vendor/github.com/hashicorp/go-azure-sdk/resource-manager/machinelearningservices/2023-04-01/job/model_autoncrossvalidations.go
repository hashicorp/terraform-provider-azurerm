package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ NCrossValidations = AutoNCrossValidations{}

type AutoNCrossValidations struct {

	// Fields inherited from NCrossValidations
}

var _ json.Marshaler = AutoNCrossValidations{}

func (s AutoNCrossValidations) MarshalJSON() ([]byte, error) {
	type wrapper AutoNCrossValidations
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AutoNCrossValidations: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AutoNCrossValidations: %+v", err)
	}
	decoded["mode"] = "Auto"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AutoNCrossValidations: %+v", err)
	}

	return encoded, nil
}
