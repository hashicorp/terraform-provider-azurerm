package integrationaccountmaps

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

type MapType string

const (
	MapTypeLiquid        MapType = "Liquid"
	MapTypeNotSpecified  MapType = "NotSpecified"
	MapTypeXslt          MapType = "Xslt"
	MapTypeXsltThreeZero MapType = "Xslt30"
	MapTypeXsltTwoZero   MapType = "Xslt20"
)

func PossibleValuesForMapType() []string {
	return []string{
		string(MapTypeLiquid),
		string(MapTypeNotSpecified),
		string(MapTypeXslt),
		string(MapTypeXsltThreeZero),
		string(MapTypeXsltTwoZero),
	}
}

func parseMapType(input string) (*MapType, error) {
	vals := map[string]MapType{
		"liquid":       MapTypeLiquid,
		"notspecified": MapTypeNotSpecified,
		"xslt":         MapTypeXslt,
		"xslt30":       MapTypeXsltThreeZero,
		"xslt20":       MapTypeXsltTwoZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MapType(input)
	return &out, nil
}
