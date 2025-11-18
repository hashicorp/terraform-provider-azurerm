package sourcecontrol

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SourceType string

const (
	SourceTypeGitHub  SourceType = "GitHub"
	SourceTypeVsoGit  SourceType = "VsoGit"
	SourceTypeVsoTfvc SourceType = "VsoTfvc"
)

func PossibleValuesForSourceType() []string {
	return []string{
		string(SourceTypeGitHub),
		string(SourceTypeVsoGit),
		string(SourceTypeVsoTfvc),
	}
}

func (s *SourceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSourceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSourceType(input string) (*SourceType, error) {
	vals := map[string]SourceType{
		"github":  SourceTypeGitHub,
		"vsogit":  SourceTypeVsoGit,
		"vsotfvc": SourceTypeVsoTfvc,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SourceType(input)
	return &out, nil
}

type TokenType string

const (
	TokenTypeOauth               TokenType = "Oauth"
	TokenTypePersonalAccessToken TokenType = "PersonalAccessToken"
)

func PossibleValuesForTokenType() []string {
	return []string{
		string(TokenTypeOauth),
		string(TokenTypePersonalAccessToken),
	}
}

func (s *TokenType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTokenType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTokenType(input string) (*TokenType, error) {
	vals := map[string]TokenType{
		"oauth":               TokenTypeOauth,
		"personalaccesstoken": TokenTypePersonalAccessToken,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TokenType(input)
	return &out, nil
}
