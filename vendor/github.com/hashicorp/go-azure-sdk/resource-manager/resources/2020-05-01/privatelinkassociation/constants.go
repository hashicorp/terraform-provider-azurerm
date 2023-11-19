package privatelinkassociation

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PublicNetworkAccessOptions string

const (
	PublicNetworkAccessOptionsDisabled PublicNetworkAccessOptions = "Disabled"
	PublicNetworkAccessOptionsEnabled  PublicNetworkAccessOptions = "Enabled"
)

func PossibleValuesForPublicNetworkAccessOptions() []string {
	return []string{
		string(PublicNetworkAccessOptionsDisabled),
		string(PublicNetworkAccessOptionsEnabled),
	}
}

func (s *PublicNetworkAccessOptions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicNetworkAccessOptions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicNetworkAccessOptions(input string) (*PublicNetworkAccessOptions, error) {
	vals := map[string]PublicNetworkAccessOptions{
		"disabled": PublicNetworkAccessOptionsDisabled,
		"enabled":  PublicNetworkAccessOptionsEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccessOptions(input)
	return &out, nil
}
