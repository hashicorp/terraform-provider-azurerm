package policyfragment

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyFragmentContentFormat string

const (
	PolicyFragmentContentFormatRawxml PolicyFragmentContentFormat = "rawxml"
	PolicyFragmentContentFormatXml    PolicyFragmentContentFormat = "xml"
)

func PossibleValuesForPolicyFragmentContentFormat() []string {
	return []string{
		string(PolicyFragmentContentFormatRawxml),
		string(PolicyFragmentContentFormatXml),
	}
}

func (s *PolicyFragmentContentFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePolicyFragmentContentFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePolicyFragmentContentFormat(input string) (*PolicyFragmentContentFormat, error) {
	vals := map[string]PolicyFragmentContentFormat{
		"rawxml": PolicyFragmentContentFormatRawxml,
		"xml":    PolicyFragmentContentFormatXml,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PolicyFragmentContentFormat(input)
	return &out, nil
}
