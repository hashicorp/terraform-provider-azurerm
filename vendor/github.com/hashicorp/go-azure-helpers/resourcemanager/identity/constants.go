package identity

import "strings"

type Type string

const (
	TypeNone                       Type = "None"
	TypeSystemAssigned             Type = "SystemAssigned"
	TypeUserAssigned               Type = "UserAssigned"
	TypeSystemAssignedUserAssigned Type = "SystemAssigned, UserAssigned"
)

func normalizeType(input Type) Type {
	vals := []Type{
		TypeNone,
		TypeSystemAssigned,
		TypeUserAssigned,
		TypeSystemAssignedUserAssigned,
	}
	for _, v := range vals {
		if strings.EqualFold(string(input), string(v)) {
			return v
		}
	}
	return input
}
