package identity

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var _ json.Marshaler = &UserAssignedMap{}

type UserAssignedMap struct {
	Type        Type                                   `json:"type" tfschema:"type"`
	IdentityIds map[string]UserAssignedIdentityDetails `json:"userAssignedIdentities"`
}

func (s *UserAssignedMap) MarshalJSON() ([]byte, error) {
	// we use a custom marshal function here since we can only send the Type / UserAssignedIdentities field
	identityType := TypeNone
	userAssignedIdentityIds := map[string]UserAssignedIdentityDetails{}

	if s != nil {
		if s.Type == TypeUserAssigned {
			identityType = TypeUserAssigned
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

// ExpandUserAssignedMap expands the schema input into a UserAssignedMap struct
func ExpandUserAssignedMap(input []interface{}) (*UserAssignedMap, error) {
	identityType := TypeNone
	identityIds := make(map[string]UserAssignedIdentityDetails, 0)

	if len(input) > 0 {
		raw := input[0].(map[string]interface{})
		typeRaw := raw["type"].(string)
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

	if len(identityIds) > 0 && identityType != TypeUserAssigned {
		return nil, fmt.Errorf("`identity_ids` can only be specified when `type` is set to %q", string(TypeUserAssigned))
	}

	return &UserAssignedMap{
		Type:        identityType,
		IdentityIds: identityIds,
	}, nil
}

// FlattenUserAssignedMap turns a UserAssignedMap into a []interface{}
func FlattenUserAssignedMap(input *UserAssignedMap) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
	}

	input.Type = normalizeType(input.Type)

	if input.Type != TypeUserAssigned {
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
		},
	}, nil
}

// ExpandUserAssignedMapFromModel expands the typed schema input into a UserAssignedMap struct
func ExpandUserAssignedMapFromModel(input []ModelUserAssigned) (*UserAssignedMap, error) {
	if len(input) == 0 {
		return &UserAssignedMap{
			Type:        TypeNone,
			IdentityIds: nil,
		}, nil
	}

	identity := input[0]

	identityIds := make(map[string]UserAssignedIdentityDetails, 0)
	for _, v := range identity.IdentityIds {
		identityIds[v] = UserAssignedIdentityDetails{
			// intentionally empty since the expand shouldn't send these values
		}
	}

	return &UserAssignedMap{
		Type:        identity.Type,
		IdentityIds: identityIds,
	}, nil
}

// FlattenUserAssignedMapToModel turns a UserAssignedMap into a typed schema model
func FlattenUserAssignedMapToModel(input *UserAssignedMap) (*[]ModelUserAssigned, error) {
	if input == nil {
		return &[]ModelUserAssigned{}, nil
	}

	input.Type = normalizeType(input.Type)

	if input.Type != TypeUserAssigned {
		return &[]ModelUserAssigned{}, nil
	}

	identityIds := make([]string, 0)
	for raw := range input.IdentityIds {
		id, err := commonids.ParseUserAssignedIdentityIDInsensitively(raw)
		if err != nil {
			return nil, fmt.Errorf("parsing %q as a User Assigned Identity ID: %+v", raw, err)
		}
		identityIds = append(identityIds, id.ID())
	}

	return &[]ModelUserAssigned{
		{
			Type:        input.Type,
			IdentityIds: identityIds,
		},
	}, nil
}
