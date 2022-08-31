package serveradministrators

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdministratorType string

const (
	AdministratorTypeActiveDirectory AdministratorType = "ActiveDirectory"
)

func PossibleValuesForAdministratorType() []string {
	return []string{
		string(AdministratorTypeActiveDirectory),
	}
}

func parseAdministratorType(input string) (*AdministratorType, error) {
	vals := map[string]AdministratorType{
		"activedirectory": AdministratorTypeActiveDirectory,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AdministratorType(input)
	return &out, nil
}
