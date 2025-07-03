package deployments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentState string

const (
	DeploymentStateActive   DeploymentState = "active"
	DeploymentStateInactive DeploymentState = "inactive"
)

func PossibleValuesForDeploymentState() []string {
	return []string{
		string(DeploymentStateActive),
		string(DeploymentStateInactive),
	}
}

func (s *DeploymentState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeploymentState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeploymentState(input string) (*DeploymentState, error) {
	vals := map[string]DeploymentState{
		"active":   DeploymentStateActive,
		"inactive": DeploymentStateInactive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentState(input)
	return &out, nil
}
