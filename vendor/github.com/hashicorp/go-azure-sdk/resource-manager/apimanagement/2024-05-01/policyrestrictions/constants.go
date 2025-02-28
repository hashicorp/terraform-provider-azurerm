package policyrestrictions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyRestrictionRequireBase string

const (
	PolicyRestrictionRequireBaseFalse PolicyRestrictionRequireBase = "false"
	PolicyRestrictionRequireBaseTrue  PolicyRestrictionRequireBase = "true"
)

func PossibleValuesForPolicyRestrictionRequireBase() []string {
	return []string{
		string(PolicyRestrictionRequireBaseFalse),
		string(PolicyRestrictionRequireBaseTrue),
	}
}

func (s *PolicyRestrictionRequireBase) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePolicyRestrictionRequireBase(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePolicyRestrictionRequireBase(input string) (*PolicyRestrictionRequireBase, error) {
	vals := map[string]PolicyRestrictionRequireBase{
		"false": PolicyRestrictionRequireBaseFalse,
		"true":  PolicyRestrictionRequireBaseTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PolicyRestrictionRequireBase(input)
	return &out, nil
}
