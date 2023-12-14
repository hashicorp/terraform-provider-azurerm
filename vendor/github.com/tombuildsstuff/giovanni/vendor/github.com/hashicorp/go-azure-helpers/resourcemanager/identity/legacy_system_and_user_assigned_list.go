// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var _ json.Marshaler = &LegacySystemAndUserAssignedList{}

type LegacySystemAndUserAssignedList struct {
	Type        Type     `json:"type" tfschema:"type"`
	PrincipalId string   `json:"principalId" tfschema:"principal_id"`
	TenantId    string   `json:"tenantId" tfschema:"tenant_id"`
	IdentityIds []string `json:"userAssignedIdentities" tfschema:"identity_ids"`
}

func (s *LegacySystemAndUserAssignedList) MarshalJSON() ([]byte, error) {
	// we use a custom marshal function here since we can only send the Type / UserAssignedIdentities field
	identityType := TypeNone
	userAssignedIdentityIds := []string{}

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

// ExpandLegacySystemAndUserAssignedList expands the schema input into a LegacySystemAndUserAssignedList struct
func ExpandLegacySystemAndUserAssignedList(input []interface{}) (*LegacySystemAndUserAssignedList, error) {
	identityType := TypeNone
	identityIds := make([]string, 0)

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
			identityIds = append(identityIds, v.(string))
		}
	}

	if len(identityIds) > 0 && (identityType != TypeSystemAssignedUserAssigned && identityType != TypeUserAssigned) {
		return nil, fmt.Errorf("`identity_ids` can only be specified when `type` is set to %q or %q", string(TypeSystemAssignedUserAssigned), string(TypeUserAssigned))
	}

	return &LegacySystemAndUserAssignedList{
		Type:        identityType,
		IdentityIds: identityIds,
	}, nil
}

// FlattenLegacySystemAndUserAssignedList turns a LegacySystemAndUserAssignedList into a []interface{}
func FlattenLegacySystemAndUserAssignedList(input *LegacySystemAndUserAssignedList) (*[]interface{}, error) {
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
	for _, raw := range input.IdentityIds {
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
