package identity

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SystemAssignedUserAssigned struct {
	Type                    Type     `tfschema:"type"`
	PrincipalId             string   `tfschema:"principal_id"`
	TenantId                string   `tfschema:"tenant_id"`
	UserAssignedIdentityIds []string `tfschema:"identity_ids"`
}

func ExpandSystemAssignedUserAssigned(input []interface{}) (*SystemAssignedUserAssigned, error) {
	if len(input) == 0 || input[0] == nil {
		return &SystemAssignedUserAssigned{
			Type: TypeNone,
		}, nil
	}

	v := input[0].(map[string]interface{})

	config := &SystemAssignedUserAssigned{
		Type: Type(v["type"].(string)),
	}

	identityIdsRaw := v["identity_ids"].(*schema.Set).List()

	if len(identityIdsRaw) != 0 {
		if config.Type != TypeUserAssigned && config.Type != TypeSystemAssignedUserAssigned {
			return nil, fmt.Errorf("`identity_ids` can only be specified when `type` includes `UserAssigned`")
		}

		identityIds := make([]string, 0)
		for _, v := range identityIdsRaw {
			identityIds = append(identityIds, v.(string))
		}

		config.UserAssignedIdentityIds = identityIds
	}

	return config, nil
}

func FlattenSystemAssignedUserAssigned(input *SystemAssignedUserAssigned) []interface{} {
	if input == nil || strings.EqualFold(string(input.Type), string(TypeNone)) {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"type":         input.Type,
			"identity_ids": input.UserAssignedIdentityIds,
			"principal_id": input.PrincipalId,
			"tenant_id":    input.TenantId,
		},
	}
}
