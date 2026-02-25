package keyvault

type NestedItemType string

const (
	NestedItemTypeAny         NestedItemType = "any"
	NestedItemTypeCertificate NestedItemType = "certificates"
	NestedItemTypeKey         NestedItemType = "keys"
	NestedItemTypeSecret      NestedItemType = "secrets"
)

// PossibleNestedItemTypeValues returns a string slice of possible "NestedItemType" values.
func PossibleNestedItemTypeValues() []string {
	return []string{
		string(NestedItemTypeCertificate),
		string(NestedItemTypeKey),
		string(NestedItemTypeSecret),
	}
}

type VersionType int

const (
	VersionTypeAny = iota
	VersionTypeVersioned
	VersionTypeVersionless
)
