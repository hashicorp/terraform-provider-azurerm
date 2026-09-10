// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"encoding/json"
	"fmt"
)

func convertDataProtectionModel[T any](input any) (*T, error) {
	encoded, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("marshaling Data Protection model: %+v", err)
	}

	var output T
	if err := json.Unmarshal(encoded, &output); err != nil {
		return nil, fmt.Errorf("unmarshaling Data Protection model: %+v", err)
	}

	return &output, nil
}
