package identity

import "strings"

type SystemAssigned struct {
	Type        Type   `tfschema:"type"`
	PrincipalId string `tfschema:"principal_id"`
	TenantId    string `tfschema:"tenant_id"`
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
	if input == nil || strings.EqualFold(string(input.Type), string(TypeNone)) {
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
