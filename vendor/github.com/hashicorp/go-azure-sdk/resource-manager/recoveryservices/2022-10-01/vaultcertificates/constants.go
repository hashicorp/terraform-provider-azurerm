package vaultcertificates

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthType string

const (
	AuthTypeAAD                  AuthType = "AAD"
	AuthTypeACS                  AuthType = "ACS"
	AuthTypeAccessControlService AuthType = "AccessControlService"
	AuthTypeAzureActiveDirectory AuthType = "AzureActiveDirectory"
	AuthTypeInvalid              AuthType = "Invalid"
)

func PossibleValuesForAuthType() []string {
	return []string{
		string(AuthTypeAAD),
		string(AuthTypeACS),
		string(AuthTypeAccessControlService),
		string(AuthTypeAzureActiveDirectory),
		string(AuthTypeInvalid),
	}
}

func parseAuthType(input string) (*AuthType, error) {
	vals := map[string]AuthType{
		"aad":                  AuthTypeAAD,
		"acs":                  AuthTypeACS,
		"accesscontrolservice": AuthTypeAccessControlService,
		"azureactivedirectory": AuthTypeAzureActiveDirectory,
		"invalid":              AuthTypeInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthType(input)
	return &out, nil
}
