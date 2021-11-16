package identity

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type UserAssigned struct {
	Type                    Type     `tfschema:"type"`
	UserAssignedIdentityIds []string `tfschema:"identity_ids"`
}

func ExpandUserAssigned(input []interface{}) (*UserAssigned, error) {
	if len(input) == 0 || input[0] == nil {
		return &UserAssigned{
			Type: TypeNone,
		}, nil
	}

	v := input[0].(map[string]interface{})

	identityIds := make([]string, 0)
	for _, v := range v["identity_ids"].(*schema.Set).List() {
		identityIds = append(identityIds, v.(string))
	}

	return &UserAssigned{
		Type:                    TypeUserAssigned,
		UserAssignedIdentityIds: identityIds,
	}, nil
}

func FlattenUserAssigned(input *UserAssigned) []interface{} {
	if input == nil || strings.EqualFold(string(input.Type), string(TypeNone)) {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"type":         input.Type,
			"identity_ids": input.UserAssignedIdentityIds,
		},
	}
}
