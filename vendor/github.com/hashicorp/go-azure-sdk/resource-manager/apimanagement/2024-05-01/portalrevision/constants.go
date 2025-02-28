package portalrevision

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PortalRevisionStatus string

const (
	PortalRevisionStatusCompleted  PortalRevisionStatus = "completed"
	PortalRevisionStatusFailed     PortalRevisionStatus = "failed"
	PortalRevisionStatusPending    PortalRevisionStatus = "pending"
	PortalRevisionStatusPublishing PortalRevisionStatus = "publishing"
)

func PossibleValuesForPortalRevisionStatus() []string {
	return []string{
		string(PortalRevisionStatusCompleted),
		string(PortalRevisionStatusFailed),
		string(PortalRevisionStatusPending),
		string(PortalRevisionStatusPublishing),
	}
}

func (s *PortalRevisionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePortalRevisionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePortalRevisionStatus(input string) (*PortalRevisionStatus, error) {
	vals := map[string]PortalRevisionStatus{
		"completed":  PortalRevisionStatusCompleted,
		"failed":     PortalRevisionStatusFailed,
		"pending":    PortalRevisionStatusPending,
		"publishing": PortalRevisionStatusPublishing,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PortalRevisionStatus(input)
	return &out, nil
}
