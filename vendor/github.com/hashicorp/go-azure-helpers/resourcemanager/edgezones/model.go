// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package edgezones

import (
	"encoding/json"
	"fmt"
	"strings"
)

var _ json.Marshaler = &Model{}
var _ json.Unmarshaler = &Model{}

type Model struct {
	Name string
}

func (m *Model) MarshalJSON() ([]byte, error) {
	out := map[string]interface{}{}

	if m.Name != "" {
		out["name"] = m.Name
		out["type"] = "EdgeZone"
	}

	return json.Marshal(out)
}

func (m *Model) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Name *string `json:"name"`
		Type *string `json:"type"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("decoding: %+v", err)
	}

	if decoded.Name == nil || decoded.Type == nil || !strings.EqualFold(*decoded.Type, "EdgeZone") {
		return nil
	}

	m.Name = *decoded.Name
	return nil
}
