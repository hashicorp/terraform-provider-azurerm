package tokens

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type TokenCertificateName string

const (
	TokenCertificateNameCertificateOne TokenCertificateName = "certificate1"
	TokenCertificateNameCertificateTwo TokenCertificateName = "certificate2"
)

func PossibleValuesForTokenCertificateName() []string {
	return []string{
		string(TokenCertificateNameCertificateOne),
		string(TokenCertificateNameCertificateTwo),
	}
}

func parseTokenCertificateName(input string) (*TokenCertificateName, error) {
	vals := map[string]TokenCertificateName{
		"certificate1": TokenCertificateNameCertificateOne,
		"certificate2": TokenCertificateNameCertificateTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TokenCertificateName(input)
	return &out, nil
}

type TokenPasswordName string

const (
	TokenPasswordNamePasswordOne TokenPasswordName = "password1"
	TokenPasswordNamePasswordTwo TokenPasswordName = "password2"
)

func PossibleValuesForTokenPasswordName() []string {
	return []string{
		string(TokenPasswordNamePasswordOne),
		string(TokenPasswordNamePasswordTwo),
	}
}

func parseTokenPasswordName(input string) (*TokenPasswordName, error) {
	vals := map[string]TokenPasswordName{
		"password1": TokenPasswordNamePasswordOne,
		"password2": TokenPasswordNamePasswordTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TokenPasswordName(input)
	return &out, nil
}

type TokenStatus string

const (
	TokenStatusDisabled TokenStatus = "disabled"
	TokenStatusEnabled  TokenStatus = "enabled"
)

func PossibleValuesForTokenStatus() []string {
	return []string{
		string(TokenStatusDisabled),
		string(TokenStatusEnabled),
	}
}

func parseTokenStatus(input string) (*TokenStatus, error) {
	vals := map[string]TokenStatus{
		"disabled": TokenStatusDisabled,
		"enabled":  TokenStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TokenStatus(input)
	return &out, nil
}
