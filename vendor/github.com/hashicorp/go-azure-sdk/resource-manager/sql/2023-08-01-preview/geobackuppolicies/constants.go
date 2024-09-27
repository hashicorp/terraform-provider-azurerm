package geobackuppolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GeoBackupPolicyState string

const (
	GeoBackupPolicyStateDisabled GeoBackupPolicyState = "Disabled"
	GeoBackupPolicyStateEnabled  GeoBackupPolicyState = "Enabled"
)

func PossibleValuesForGeoBackupPolicyState() []string {
	return []string{
		string(GeoBackupPolicyStateDisabled),
		string(GeoBackupPolicyStateEnabled),
	}
}

func (s *GeoBackupPolicyState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGeoBackupPolicyState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGeoBackupPolicyState(input string) (*GeoBackupPolicyState, error) {
	vals := map[string]GeoBackupPolicyState{
		"disabled": GeoBackupPolicyStateDisabled,
		"enabled":  GeoBackupPolicyStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GeoBackupPolicyState(input)
	return &out, nil
}
