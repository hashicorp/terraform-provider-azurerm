package remediations

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceDiscoveryMode string

const (
	ResourceDiscoveryModeExistingNonCompliant ResourceDiscoveryMode = "ExistingNonCompliant"
	ResourceDiscoveryModeReEvaluateCompliance ResourceDiscoveryMode = "ReEvaluateCompliance"
)

func PossibleValuesForResourceDiscoveryMode() []string {
	return []string{
		string(ResourceDiscoveryModeExistingNonCompliant),
		string(ResourceDiscoveryModeReEvaluateCompliance),
	}
}

func (s *ResourceDiscoveryMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceDiscoveryMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceDiscoveryMode(input string) (*ResourceDiscoveryMode, error) {
	vals := map[string]ResourceDiscoveryMode{
		"existingnoncompliant": ResourceDiscoveryModeExistingNonCompliant,
		"reevaluatecompliance": ResourceDiscoveryModeReEvaluateCompliance,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceDiscoveryMode(input)
	return &out, nil
}
