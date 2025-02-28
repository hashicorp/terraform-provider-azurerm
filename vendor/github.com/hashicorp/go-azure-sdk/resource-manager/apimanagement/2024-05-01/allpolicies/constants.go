package allpolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyComplianceState string

const (
	PolicyComplianceStateCompliant    PolicyComplianceState = "Compliant"
	PolicyComplianceStateNonCompliant PolicyComplianceState = "NonCompliant"
	PolicyComplianceStatePending      PolicyComplianceState = "Pending"
)

func PossibleValuesForPolicyComplianceState() []string {
	return []string{
		string(PolicyComplianceStateCompliant),
		string(PolicyComplianceStateNonCompliant),
		string(PolicyComplianceStatePending),
	}
}

func (s *PolicyComplianceState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePolicyComplianceState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePolicyComplianceState(input string) (*PolicyComplianceState, error) {
	vals := map[string]PolicyComplianceState{
		"compliant":    PolicyComplianceStateCompliant,
		"noncompliant": PolicyComplianceStateNonCompliant,
		"pending":      PolicyComplianceStatePending,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PolicyComplianceState(input)
	return &out, nil
}
