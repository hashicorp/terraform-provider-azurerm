package clients

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClientCertificateValidationScheme string

const (
	ClientCertificateValidationSchemeDnsMatchesAuthenticationName     ClientCertificateValidationScheme = "DnsMatchesAuthenticationName"
	ClientCertificateValidationSchemeEmailMatchesAuthenticationName   ClientCertificateValidationScheme = "EmailMatchesAuthenticationName"
	ClientCertificateValidationSchemeIPMatchesAuthenticationName      ClientCertificateValidationScheme = "IpMatchesAuthenticationName"
	ClientCertificateValidationSchemeSubjectMatchesAuthenticationName ClientCertificateValidationScheme = "SubjectMatchesAuthenticationName"
	ClientCertificateValidationSchemeThumbprintMatch                  ClientCertificateValidationScheme = "ThumbprintMatch"
	ClientCertificateValidationSchemeUriMatchesAuthenticationName     ClientCertificateValidationScheme = "UriMatchesAuthenticationName"
)

func PossibleValuesForClientCertificateValidationScheme() []string {
	return []string{
		string(ClientCertificateValidationSchemeDnsMatchesAuthenticationName),
		string(ClientCertificateValidationSchemeEmailMatchesAuthenticationName),
		string(ClientCertificateValidationSchemeIPMatchesAuthenticationName),
		string(ClientCertificateValidationSchemeSubjectMatchesAuthenticationName),
		string(ClientCertificateValidationSchemeThumbprintMatch),
		string(ClientCertificateValidationSchemeUriMatchesAuthenticationName),
	}
}

func (s *ClientCertificateValidationScheme) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClientCertificateValidationScheme(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClientCertificateValidationScheme(input string) (*ClientCertificateValidationScheme, error) {
	vals := map[string]ClientCertificateValidationScheme{
		"dnsmatchesauthenticationname":     ClientCertificateValidationSchemeDnsMatchesAuthenticationName,
		"emailmatchesauthenticationname":   ClientCertificateValidationSchemeEmailMatchesAuthenticationName,
		"ipmatchesauthenticationname":      ClientCertificateValidationSchemeIPMatchesAuthenticationName,
		"subjectmatchesauthenticationname": ClientCertificateValidationSchemeSubjectMatchesAuthenticationName,
		"thumbprintmatch":                  ClientCertificateValidationSchemeThumbprintMatch,
		"urimatchesauthenticationname":     ClientCertificateValidationSchemeUriMatchesAuthenticationName,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClientCertificateValidationScheme(input)
	return &out, nil
}

type ClientProvisioningState string

const (
	ClientProvisioningStateCanceled  ClientProvisioningState = "Canceled"
	ClientProvisioningStateCreating  ClientProvisioningState = "Creating"
	ClientProvisioningStateDeleted   ClientProvisioningState = "Deleted"
	ClientProvisioningStateDeleting  ClientProvisioningState = "Deleting"
	ClientProvisioningStateFailed    ClientProvisioningState = "Failed"
	ClientProvisioningStateSucceeded ClientProvisioningState = "Succeeded"
	ClientProvisioningStateUpdating  ClientProvisioningState = "Updating"
)

func PossibleValuesForClientProvisioningState() []string {
	return []string{
		string(ClientProvisioningStateCanceled),
		string(ClientProvisioningStateCreating),
		string(ClientProvisioningStateDeleted),
		string(ClientProvisioningStateDeleting),
		string(ClientProvisioningStateFailed),
		string(ClientProvisioningStateSucceeded),
		string(ClientProvisioningStateUpdating),
	}
}

func (s *ClientProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClientProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClientProvisioningState(input string) (*ClientProvisioningState, error) {
	vals := map[string]ClientProvisioningState{
		"canceled":  ClientProvisioningStateCanceled,
		"creating":  ClientProvisioningStateCreating,
		"deleted":   ClientProvisioningStateDeleted,
		"deleting":  ClientProvisioningStateDeleting,
		"failed":    ClientProvisioningStateFailed,
		"succeeded": ClientProvisioningStateSucceeded,
		"updating":  ClientProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClientProvisioningState(input)
	return &out, nil
}

type ClientState string

const (
	ClientStateDisabled ClientState = "Disabled"
	ClientStateEnabled  ClientState = "Enabled"
)

func PossibleValuesForClientState() []string {
	return []string{
		string(ClientStateDisabled),
		string(ClientStateEnabled),
	}
}

func (s *ClientState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClientState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClientState(input string) (*ClientState, error) {
	vals := map[string]ClientState{
		"disabled": ClientStateDisabled,
		"enabled":  ClientStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClientState(input)
	return &out, nil
}
