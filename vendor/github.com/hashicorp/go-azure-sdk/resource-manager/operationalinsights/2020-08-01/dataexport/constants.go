package dataexport

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Type string

const (
	TypeEventHub       Type = "EventHub"
	TypeStorageAccount Type = "StorageAccount"
)

func PossibleValuesForType() []string {
	return []string{
		string(TypeEventHub),
		string(TypeStorageAccount),
	}
}

func parseType(input string) (*Type, error) {
	vals := map[string]Type{
		"eventhub":       TypeEventHub,
		"storageaccount": TypeStorageAccount,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Type(input)
	return &out, nil
}
