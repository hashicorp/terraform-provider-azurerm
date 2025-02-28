package authorization

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthorizationType string

const (
	AuthorizationTypeOAuthTwo AuthorizationType = "OAuth2"
)

func PossibleValuesForAuthorizationType() []string {
	return []string{
		string(AuthorizationTypeOAuthTwo),
	}
}

func (s *AuthorizationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAuthorizationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAuthorizationType(input string) (*AuthorizationType, error) {
	vals := map[string]AuthorizationType{
		"oauth2": AuthorizationTypeOAuthTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthorizationType(input)
	return &out, nil
}

type OAuth2GrantType string

const (
	OAuth2GrantTypeAuthorizationCode OAuth2GrantType = "AuthorizationCode"
	OAuth2GrantTypeClientCredentials OAuth2GrantType = "ClientCredentials"
)

func PossibleValuesForOAuth2GrantType() []string {
	return []string{
		string(OAuth2GrantTypeAuthorizationCode),
		string(OAuth2GrantTypeClientCredentials),
	}
}

func (s *OAuth2GrantType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOAuth2GrantType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOAuth2GrantType(input string) (*OAuth2GrantType, error) {
	vals := map[string]OAuth2GrantType{
		"authorizationcode": OAuth2GrantTypeAuthorizationCode,
		"clientcredentials": OAuth2GrantTypeClientCredentials,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OAuth2GrantType(input)
	return &out, nil
}
