package accounts

type CreatedByType string

const (
	CreatedByTypeApplication     CreatedByType = "Application"
	CreatedByTypeKey             CreatedByType = "Key"
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	CreatedByTypeUser            CreatedByType = "User"
)

type KeyType string

const (
	KeyTypePrimary   KeyType = "primary"
	KeyTypeSecondary KeyType = "secondary"
)

type Kind string

const (
	KindGenOne Kind = "Gen1"
	KindGenTwo Kind = "Gen2"
)

type Name string

const (
	NameGTwo  Name = "G2"
	NameSOne  Name = "S1"
	NameSZero Name = "S0"
)
