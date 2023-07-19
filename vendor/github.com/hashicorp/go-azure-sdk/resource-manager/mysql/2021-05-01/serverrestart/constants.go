package serverrestart

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnableStatusEnum string

const (
	EnableStatusEnumDisabled EnableStatusEnum = "Disabled"
	EnableStatusEnumEnabled  EnableStatusEnum = "Enabled"
)

func PossibleValuesForEnableStatusEnum() []string {
	return []string{
		string(EnableStatusEnumDisabled),
		string(EnableStatusEnumEnabled),
	}
}

func parseEnableStatusEnum(input string) (*EnableStatusEnum, error) {
	vals := map[string]EnableStatusEnum{
		"disabled": EnableStatusEnumDisabled,
		"enabled":  EnableStatusEnumEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnableStatusEnum(input)
	return &out, nil
}
