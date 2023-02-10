package managedenvironmentsstorages

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessMode string

const (
	AccessModeReadOnly  AccessMode = "ReadOnly"
	AccessModeReadWrite AccessMode = "ReadWrite"
)

func PossibleValuesForAccessMode() []string {
	return []string{
		string(AccessModeReadOnly),
		string(AccessModeReadWrite),
	}
}

func parseAccessMode(input string) (*AccessMode, error) {
	vals := map[string]AccessMode{
		"readonly":  AccessModeReadOnly,
		"readwrite": AccessModeReadWrite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessMode(input)
	return &out, nil
}
