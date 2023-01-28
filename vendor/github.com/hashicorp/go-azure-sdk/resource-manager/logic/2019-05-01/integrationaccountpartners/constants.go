package integrationaccountpartners

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyType string

const (
	KeyTypeNotSpecified KeyType = "NotSpecified"
	KeyTypePrimary      KeyType = "Primary"
	KeyTypeSecondary    KeyType = "Secondary"
)

func PossibleValuesForKeyType() []string {
	return []string{
		string(KeyTypeNotSpecified),
		string(KeyTypePrimary),
		string(KeyTypeSecondary),
	}
}

func parseKeyType(input string) (*KeyType, error) {
	vals := map[string]KeyType{
		"notspecified": KeyTypeNotSpecified,
		"primary":      KeyTypePrimary,
		"secondary":    KeyTypeSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyType(input)
	return &out, nil
}

type PartnerType string

const (
	PartnerTypeBTwoB        PartnerType = "B2B"
	PartnerTypeNotSpecified PartnerType = "NotSpecified"
)

func PossibleValuesForPartnerType() []string {
	return []string{
		string(PartnerTypeBTwoB),
		string(PartnerTypeNotSpecified),
	}
}

func parsePartnerType(input string) (*PartnerType, error) {
	vals := map[string]PartnerType{
		"b2b":          PartnerTypeBTwoB,
		"notspecified": PartnerTypeNotSpecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PartnerType(input)
	return &out, nil
}
