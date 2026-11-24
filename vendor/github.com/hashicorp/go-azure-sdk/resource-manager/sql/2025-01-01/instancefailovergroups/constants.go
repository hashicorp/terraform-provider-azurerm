package instancefailovergroups

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InstanceFailoverGroupReplicationRole string

const (
	InstanceFailoverGroupReplicationRolePrimary   InstanceFailoverGroupReplicationRole = "Primary"
	InstanceFailoverGroupReplicationRoleSecondary InstanceFailoverGroupReplicationRole = "Secondary"
)

func PossibleValuesForInstanceFailoverGroupReplicationRole() []string {
	return []string{
		string(InstanceFailoverGroupReplicationRolePrimary),
		string(InstanceFailoverGroupReplicationRoleSecondary),
	}
}

func (s *InstanceFailoverGroupReplicationRole) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseInstanceFailoverGroupReplicationRole(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseInstanceFailoverGroupReplicationRole(input string) (*InstanceFailoverGroupReplicationRole, error) {
	vals := map[string]InstanceFailoverGroupReplicationRole{
		"primary":   InstanceFailoverGroupReplicationRolePrimary,
		"secondary": InstanceFailoverGroupReplicationRoleSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InstanceFailoverGroupReplicationRole(input)
	return &out, nil
}

type ReadOnlyEndpointFailoverPolicy string

const (
	ReadOnlyEndpointFailoverPolicyDisabled ReadOnlyEndpointFailoverPolicy = "Disabled"
	ReadOnlyEndpointFailoverPolicyEnabled  ReadOnlyEndpointFailoverPolicy = "Enabled"
)

func PossibleValuesForReadOnlyEndpointFailoverPolicy() []string {
	return []string{
		string(ReadOnlyEndpointFailoverPolicyDisabled),
		string(ReadOnlyEndpointFailoverPolicyEnabled),
	}
}

func (s *ReadOnlyEndpointFailoverPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReadOnlyEndpointFailoverPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReadOnlyEndpointFailoverPolicy(input string) (*ReadOnlyEndpointFailoverPolicy, error) {
	vals := map[string]ReadOnlyEndpointFailoverPolicy{
		"disabled": ReadOnlyEndpointFailoverPolicyDisabled,
		"enabled":  ReadOnlyEndpointFailoverPolicyEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReadOnlyEndpointFailoverPolicy(input)
	return &out, nil
}

type ReadWriteEndpointFailoverPolicy string

const (
	ReadWriteEndpointFailoverPolicyAutomatic ReadWriteEndpointFailoverPolicy = "Automatic"
	ReadWriteEndpointFailoverPolicyManual    ReadWriteEndpointFailoverPolicy = "Manual"
)

func PossibleValuesForReadWriteEndpointFailoverPolicy() []string {
	return []string{
		string(ReadWriteEndpointFailoverPolicyAutomatic),
		string(ReadWriteEndpointFailoverPolicyManual),
	}
}

func (s *ReadWriteEndpointFailoverPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReadWriteEndpointFailoverPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReadWriteEndpointFailoverPolicy(input string) (*ReadWriteEndpointFailoverPolicy, error) {
	vals := map[string]ReadWriteEndpointFailoverPolicy{
		"automatic": ReadWriteEndpointFailoverPolicyAutomatic,
		"manual":    ReadWriteEndpointFailoverPolicyManual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReadWriteEndpointFailoverPolicy(input)
	return &out, nil
}

type SecondaryInstanceType string

const (
	SecondaryInstanceTypeGeo     SecondaryInstanceType = "Geo"
	SecondaryInstanceTypeStandby SecondaryInstanceType = "Standby"
)

func PossibleValuesForSecondaryInstanceType() []string {
	return []string{
		string(SecondaryInstanceTypeGeo),
		string(SecondaryInstanceTypeStandby),
	}
}

func (s *SecondaryInstanceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecondaryInstanceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecondaryInstanceType(input string) (*SecondaryInstanceType, error) {
	vals := map[string]SecondaryInstanceType{
		"geo":     SecondaryInstanceTypeGeo,
		"standby": SecondaryInstanceTypeStandby,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecondaryInstanceType(input)
	return &out, nil
}
