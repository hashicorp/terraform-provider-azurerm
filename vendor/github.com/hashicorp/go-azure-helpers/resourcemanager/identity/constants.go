package identity

type Type string

const (
	TypeNone                       Type = "None"
	TypeSystemAssigned             Type = "SystemAssigned"
	TypeUserAssigned               Type = "UserAssigned"
	TypeSystemAssignedUserAssigned Type = "SystemAssigned, UserAssigned"
)
