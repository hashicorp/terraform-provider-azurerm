package hybridrunbookworkergroup

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GroupTypeEnum string

const (
	GroupTypeEnumSystem GroupTypeEnum = "System"
	GroupTypeEnumUser   GroupTypeEnum = "User"
)

func PossibleValuesForGroupTypeEnum() []string {
	return []string{
		string(GroupTypeEnumSystem),
		string(GroupTypeEnumUser),
	}
}

func parseGroupTypeEnum(input string) (*GroupTypeEnum, error) {
	vals := map[string]GroupTypeEnum{
		"system": GroupTypeEnumSystem,
		"user":   GroupTypeEnumUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GroupTypeEnum(input)
	return &out, nil
}
