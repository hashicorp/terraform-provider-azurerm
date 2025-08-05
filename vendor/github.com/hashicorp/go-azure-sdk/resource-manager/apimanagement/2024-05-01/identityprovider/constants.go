package identityprovider

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IdentityProviderType string

const (
	IdentityProviderTypeAad       IdentityProviderType = "aad"
	IdentityProviderTypeAadBTwoC  IdentityProviderType = "aadB2C"
	IdentityProviderTypeFacebook  IdentityProviderType = "facebook"
	IdentityProviderTypeGoogle    IdentityProviderType = "google"
	IdentityProviderTypeMicrosoft IdentityProviderType = "microsoft"
	IdentityProviderTypeTwitter   IdentityProviderType = "twitter"
)

func PossibleValuesForIdentityProviderType() []string {
	return []string{
		string(IdentityProviderTypeAad),
		string(IdentityProviderTypeAadBTwoC),
		string(IdentityProviderTypeFacebook),
		string(IdentityProviderTypeGoogle),
		string(IdentityProviderTypeMicrosoft),
		string(IdentityProviderTypeTwitter),
	}
}

func (s *IdentityProviderType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIdentityProviderType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIdentityProviderType(input string) (*IdentityProviderType, error) {
	vals := map[string]IdentityProviderType{
		"aad":       IdentityProviderTypeAad,
		"aadb2c":    IdentityProviderTypeAadBTwoC,
		"facebook":  IdentityProviderTypeFacebook,
		"google":    IdentityProviderTypeGoogle,
		"microsoft": IdentityProviderTypeMicrosoft,
		"twitter":   IdentityProviderTypeTwitter,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IdentityProviderType(input)
	return &out, nil
}
