// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var _ json.Marshaler = &LegacySystemAndUserAssignedMap{}

type LegacySystemAndUserAssignedMap struct {
	Type        Type                                   `json:"type" tfschema:"type"`
	PrincipalId string                                 `json:"principalId" tfschema:"principal_id"`
	TenantId    string                                 `json:"tenantId" tfschema:"tenant_id"`
	IdentityIds map[string]UserAssignedIdentityDetails `json:"userAssignedIdentities"`
}

func (s *LegacySystemAndUserAssignedMap) MarshalJSON() ([]byte, error) {
	// we use a custom marshal function here since we can only send the Type / UserAssignedIdentities field
	identityType := TypeNone
	userAssignedIdentityIds := map[string]UserAssignedIdentityDetails{}

	if s != nil {
		if s.Type == typeLegacySystemAssignedUserAssigned {
			return nil, fmt.Errorf("internal error: the legacy `SystemAssigned,UserAssigned` identity type should be being converted to the schema type - this is a bug")
		}

		if s.Type == TypeSystemAssigned {
			identityType = TypeSystemAssigned
		}
		if s.Type == TypeSystemAssignedUserAssigned {
			// convert the Schema value (w/spaces) to the legacy API value (w/o spaces)
			identityType = typeLegacySystemAssignedUserAssigned
		}
		if s.Type == TypeUserAssigned {
			identityType = TypeUserAssigned
		}

		if identityType != TypeNone {
			userAssignedIdentityIds = s.IdentityIds
		}
	}

	out := map[string]interface{}{
		"type":                   string(identityType),
		"userAssignedIdentities": nil,
	}
	if len(userAssignedIdentityIds) > 0 {
		out["userAssignedIdentities"] = userAssignedIdentityIds
	}
	return json.Marshal(out)
}

func (s *LegacySystemAndUserAssignedMap) UnmarshalJSON(input []byte) error {
	if input == nil {
		return nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return fmt.Errorf("unmarshaling LegacySystemAndUserAssignedMap into map[string]interface: %+v", err)
	}
	typeVal := TypeNone
	if v, ok := temp["type"].(string); ok && v != "" {
		if strings.EqualFold(v, string(TypeSystemAssigned)) {
			typeVal = TypeSystemAssigned
		}
		if strings.EqualFold(v, string(TypeUserAssigned)) {
			typeVal = TypeUserAssigned
		}
		if strings.EqualFold(v, string(typeLegacySystemAssignedUserAssigned)) {
			typeVal = TypeSystemAssignedUserAssigned
		}
		if strings.EqualFold(v, string(TypeSystemAssignedUserAssigned)) {
			typeVal = TypeSystemAssignedUserAssigned
		}
	}

	type alias LegacySystemAndUserAssignedMap
	var decoded alias
	if err := json.Unmarshal(input, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Type = typeVal
	s.IdentityIds = decoded.IdentityIds
	s.PrincipalId = decoded.PrincipalId
	s.TenantId = decoded.TenantId

	return nil
}

// ExpandLegacySystemAndUserAssignedMap expands the schema input into a LegacySystemAndUserAssignedMap struct
func ExpandLegacySystemAndUserAssignedMap(input []interface{}) (*LegacySystemAndUserAssignedMap, error) {
	identityType := TypeNone
	identityIds := make(map[string]UserAssignedIdentityDetails, 0)

	if len(input) > 0 {
		raw := input[0].(map[string]interface{})
		typeRaw := raw["type"].(string)
		if typeRaw == string(TypeSystemAssigned) {
			identityType = TypeSystemAssigned
		}
		if typeRaw == string(TypeSystemAssignedUserAssigned) {
			identityType = TypeSystemAssignedUserAssigned
		}
		if typeRaw == string(TypeUserAssigned) {
			identityType = TypeUserAssigned
		}

		identityIdsRaw := raw["identity_ids"].(*schema.Set).List()
		for _, v := range identityIdsRaw {
			identityIds[v.(string)] = UserAssignedIdentityDetails{
				// intentionally empty since the expand shouldn't send these values
			}
		}
	}

	if len(identityIds) > 0 && (identityType != TypeSystemAssignedUserAssigned && identityType != TypeUserAssigned) {
		return nil, fmt.Errorf("`identity_ids` can only be specified when `type` is set to %q or %q", string(TypeSystemAssignedUserAssigned), string(TypeUserAssigned))
	}

	return &LegacySystemAndUserAssignedMap{
		Type:        identityType,
		IdentityIds: identityIds,
	}, nil
}

// FlattenLegacySystemAndUserAssignedMap turns a LegacySystemAndUserAssignedMap into a []interface{}
func FlattenLegacySystemAndUserAssignedMap(input *LegacySystemAndUserAssignedMap) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
	}

	input.Type = normalizeType(input.Type)
	validTypes := map[Type]struct{}{
		TypeSystemAssigned:             {},
		TypeSystemAssignedUserAssigned: {},
		TypeUserAssigned:               {},
	}
	if _, ok := validTypes[input.Type]; !ok {
		return &[]interface{}{}, nil
	}

	identityIds := make([]string, 0)
	for raw := range input.IdentityIds {
		id, err := commonids.ParseUserAssignedIdentityIDInsensitively(raw)
		if err != nil {
			return nil, fmt.Errorf("parsing %q as a User Assigned Identity ID: %+v", raw, err)
		}
		identityIds = append(identityIds, id.ID())
	}

	return &[]interface{}{
		map[string]interface{}{
			"type":         string(input.Type),
			"identity_ids": identityIds,
			"principal_id": input.PrincipalId,
			"tenant_id":    input.TenantId,
		},
	}, nil
}

// ExpandLegacySystemAndUserAssignedMapFromModel expands the schema input into a LegacySystemAndUserAssignedMap struct
func ExpandLegacySystemAndUserAssignedMapFromModel(input []ModelSystemAssignedUserAssigned) (*LegacySystemAndUserAssignedMap, error) {
	if len(input) == 0 {
		return &LegacySystemAndUserAssignedMap{
			Type:        TypeNone,
			IdentityIds: nil,
		}, nil
	}

	identityType := input[0].Type
	identityIds := make(map[string]UserAssignedIdentityDetails, 0)
	for _, v := range input[0].IdentityIds {
		identityIds[v] = UserAssignedIdentityDetails{
			// intentionally empty since the expand shouldn't send these values
		}
	}

	if len(identityIds) > 0 && (identityType != TypeSystemAssignedUserAssigned && identityType != TypeUserAssigned) {
		return nil, fmt.Errorf("`identity_ids` can only be specified when `type` is set to %q or %q", string(TypeSystemAssignedUserAssigned), string(TypeUserAssigned))
	}

	return &LegacySystemAndUserAssignedMap{
		Type:        identityType,
		IdentityIds: identityIds,
	}, nil
}

// FlattenLegacySystemAndUserAssignedMapToModel turns a LegacySystemAndUserAssignedMap into a []interface{}
func FlattenLegacySystemAndUserAssignedMapToModel(input *LegacySystemAndUserAssignedMap) ([]ModelSystemAssignedUserAssigned, error) {
	if input == nil {
		return []ModelSystemAssignedUserAssigned{}, nil
	}

	input.Type = normalizeType(input.Type)
	if input.Type != TypeSystemAssigned && input.Type != TypeSystemAssignedUserAssigned && input.Type != TypeUserAssigned {
		return []ModelSystemAssignedUserAssigned{}, nil
	}

	identityIds := make([]string, 0)
	for raw, _ := range input.IdentityIds {
		id, err := commonids.ParseUserAssignedIdentityIDInsensitively(raw)
		if err != nil {
			return nil, fmt.Errorf("parsing %q as a User Assigned Identity ID: %+v", raw, err)
		}
		identityIds = append(identityIds, id.ID())
	}

	return []ModelSystemAssignedUserAssigned{
		{
			Type:        input.Type,
			IdentityIds: identityIds,
			PrincipalId: input.PrincipalId,
			TenantId:    input.TenantId,
		},
	}, nil
}
