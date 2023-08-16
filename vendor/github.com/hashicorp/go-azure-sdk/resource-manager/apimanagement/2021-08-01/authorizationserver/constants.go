package authorizationserver

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthorizationMethod string

const (
	AuthorizationMethodDELETE  AuthorizationMethod = "DELETE"
	AuthorizationMethodGET     AuthorizationMethod = "GET"
	AuthorizationMethodHEAD    AuthorizationMethod = "HEAD"
	AuthorizationMethodOPTIONS AuthorizationMethod = "OPTIONS"
	AuthorizationMethodPATCH   AuthorizationMethod = "PATCH"
	AuthorizationMethodPOST    AuthorizationMethod = "POST"
	AuthorizationMethodPUT     AuthorizationMethod = "PUT"
	AuthorizationMethodTRACE   AuthorizationMethod = "TRACE"
)

func PossibleValuesForAuthorizationMethod() []string {
	return []string{
		string(AuthorizationMethodDELETE),
		string(AuthorizationMethodGET),
		string(AuthorizationMethodHEAD),
		string(AuthorizationMethodOPTIONS),
		string(AuthorizationMethodPATCH),
		string(AuthorizationMethodPOST),
		string(AuthorizationMethodPUT),
		string(AuthorizationMethodTRACE),
	}
}

func (s *AuthorizationMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAuthorizationMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAuthorizationMethod(input string) (*AuthorizationMethod, error) {
	vals := map[string]AuthorizationMethod{
		"delete":  AuthorizationMethodDELETE,
		"get":     AuthorizationMethodGET,
		"head":    AuthorizationMethodHEAD,
		"options": AuthorizationMethodOPTIONS,
		"patch":   AuthorizationMethodPATCH,
		"post":    AuthorizationMethodPOST,
		"put":     AuthorizationMethodPUT,
		"trace":   AuthorizationMethodTRACE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthorizationMethod(input)
	return &out, nil
}

type BearerTokenSendingMethod string

const (
	BearerTokenSendingMethodAuthorizationHeader BearerTokenSendingMethod = "authorizationHeader"
	BearerTokenSendingMethodQuery               BearerTokenSendingMethod = "query"
)

func PossibleValuesForBearerTokenSendingMethod() []string {
	return []string{
		string(BearerTokenSendingMethodAuthorizationHeader),
		string(BearerTokenSendingMethodQuery),
	}
}

func (s *BearerTokenSendingMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBearerTokenSendingMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBearerTokenSendingMethod(input string) (*BearerTokenSendingMethod, error) {
	vals := map[string]BearerTokenSendingMethod{
		"authorizationheader": BearerTokenSendingMethodAuthorizationHeader,
		"query":               BearerTokenSendingMethodQuery,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BearerTokenSendingMethod(input)
	return &out, nil
}

type ClientAuthenticationMethod string

const (
	ClientAuthenticationMethodBasic ClientAuthenticationMethod = "Basic"
	ClientAuthenticationMethodBody  ClientAuthenticationMethod = "Body"
)

func PossibleValuesForClientAuthenticationMethod() []string {
	return []string{
		string(ClientAuthenticationMethodBasic),
		string(ClientAuthenticationMethodBody),
	}
}

func (s *ClientAuthenticationMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClientAuthenticationMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClientAuthenticationMethod(input string) (*ClientAuthenticationMethod, error) {
	vals := map[string]ClientAuthenticationMethod{
		"basic": ClientAuthenticationMethodBasic,
		"body":  ClientAuthenticationMethodBody,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClientAuthenticationMethod(input)
	return &out, nil
}

type GrantType string

const (
	GrantTypeAuthorizationCode     GrantType = "authorizationCode"
	GrantTypeClientCredentials     GrantType = "clientCredentials"
	GrantTypeImplicit              GrantType = "implicit"
	GrantTypeResourceOwnerPassword GrantType = "resourceOwnerPassword"
)

func PossibleValuesForGrantType() []string {
	return []string{
		string(GrantTypeAuthorizationCode),
		string(GrantTypeClientCredentials),
		string(GrantTypeImplicit),
		string(GrantTypeResourceOwnerPassword),
	}
}

func (s *GrantType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGrantType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGrantType(input string) (*GrantType, error) {
	vals := map[string]GrantType{
		"authorizationcode":     GrantTypeAuthorizationCode,
		"clientcredentials":     GrantTypeClientCredentials,
		"implicit":              GrantTypeImplicit,
		"resourceownerpassword": GrantTypeResourceOwnerPassword,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GrantType(input)
	return &out, nil
}
