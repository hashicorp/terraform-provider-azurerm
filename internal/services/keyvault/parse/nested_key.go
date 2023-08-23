package parse

import "fmt"

func NewNestedKeyID(keyVaultBaseUrl string, name, version string) (*NestedItemId, error) {
	return NewNestedItemID(keyVaultBaseUrl, NestedItemTypeKey, name, version)
}

// ParseNestedKeyID parses a Key Vault Nested Key ID containing a version into a NestedItemId object
func ParseNestedKeyID(input string) (*NestedItemId, error) {
	item, err := ParseOptionallyVersionedNestedKeyID(input)
	if err != nil {
		return nil, err
	}

	if item.Version == "" {
		return nil, fmt.Errorf("expected a key vault versioned ID but no version information was found in: %q", input)
	}

	return item, nil
}

// ParseOptionallyVersionedNestedKeyID parses a Key Vault Nested Key ID
// optionally containing a version into a NestedItemId object
func ParseOptionallyVersionedNestedKeyID(input string) (*NestedItemId, error) {
	item, err := parseNestedItemId(input)
	if err != nil {
		return nil, err
	}

	if item.NestedItemType != NestedItemTypeKey {
		return nil, fmt.Errorf("expected a Key Vault Key but got %q", string(item.NestedItemType))
	}

	return item, nil
}
