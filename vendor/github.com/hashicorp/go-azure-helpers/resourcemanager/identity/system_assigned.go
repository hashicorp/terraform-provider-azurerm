// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity

import (
	"encoding/json"
)

var _ json.Marshaler = &SystemAssigned{}

type SystemAssigned struct {
	Type        Type   `json:"type" tfschema:"type"`
	PrincipalId string `json:"principalId" tfschema:"principal_id"`
	TenantId    string `json:"tenantId" tfschema:"tenant_id"`
}

func (s *SystemAssigned) MarshalJSON() ([]byte, error) {
	// we use a custom marshal function here since we can only send the Type field
	out := map[string]interface{}{
		"type": string(TypeNone),
	}
	if s != nil && s.Type == TypeSystemAssigned {
		out["type"] = string(TypeSystemAssigned)
	}
	return json.Marshal(out)
}

func ExpandSystemAssigned(input []interface{}) (*SystemAssigned, error) {
	if len(input) == 0 || input[0] == nil {
		return &SystemAssigned{
			Type: TypeNone,
		}, nil
	}

	return &SystemAssigned{
		Type: TypeSystemAssigned,
	}, nil
}

func FlattenSystemAssigned(input *SystemAssigned) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	input.Type = normalizeType(input.Type)

	if input.Type == TypeNone {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"type":         input.Type,
			"principal_id": input.PrincipalId,
			"tenant_id":    input.TenantId,
		},
	}
}

func ExpandSystemAssignedFromModel(input []ModelSystemAssigned) (*SystemAssigned, error) {
	if len(input) == 0 {
		return &SystemAssigned{
			Type: TypeNone,
		}, nil
	}

	return &SystemAssigned{
		Type: TypeSystemAssigned,
	}, nil
}

func FlattenSystemAssignedToModel(input *SystemAssigned) []ModelSystemAssigned {
	if input == nil {
		return []ModelSystemAssigned{}
	}

	input.Type = normalizeType(input.Type)

	if input.Type == TypeNone {
		return []ModelSystemAssigned{}
	}

	return []ModelSystemAssigned{
		{
			Type:        input.Type,
			PrincipalId: input.PrincipalId,
			TenantId:    input.TenantId,
		},
	}
}
