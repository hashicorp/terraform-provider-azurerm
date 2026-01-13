package dnsprivatezones

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DnsPrivateZonesLifecycleState string

const (
	DnsPrivateZonesLifecycleStateActive   DnsPrivateZonesLifecycleState = "Active"
	DnsPrivateZonesLifecycleStateCreating DnsPrivateZonesLifecycleState = "Creating"
	DnsPrivateZonesLifecycleStateDeleted  DnsPrivateZonesLifecycleState = "Deleted"
	DnsPrivateZonesLifecycleStateDeleting DnsPrivateZonesLifecycleState = "Deleting"
	DnsPrivateZonesLifecycleStateUpdating DnsPrivateZonesLifecycleState = "Updating"
)

func PossibleValuesForDnsPrivateZonesLifecycleState() []string {
	return []string{
		string(DnsPrivateZonesLifecycleStateActive),
		string(DnsPrivateZonesLifecycleStateCreating),
		string(DnsPrivateZonesLifecycleStateDeleted),
		string(DnsPrivateZonesLifecycleStateDeleting),
		string(DnsPrivateZonesLifecycleStateUpdating),
	}
}

func (s *DnsPrivateZonesLifecycleState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDnsPrivateZonesLifecycleState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDnsPrivateZonesLifecycleState(input string) (*DnsPrivateZonesLifecycleState, error) {
	vals := map[string]DnsPrivateZonesLifecycleState{
		"active":   DnsPrivateZonesLifecycleStateActive,
		"creating": DnsPrivateZonesLifecycleStateCreating,
		"deleted":  DnsPrivateZonesLifecycleStateDeleted,
		"deleting": DnsPrivateZonesLifecycleStateDeleting,
		"updating": DnsPrivateZonesLifecycleStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DnsPrivateZonesLifecycleState(input)
	return &out, nil
}

type ResourceProvisioningState string

const (
	ResourceProvisioningStateCanceled  ResourceProvisioningState = "Canceled"
	ResourceProvisioningStateFailed    ResourceProvisioningState = "Failed"
	ResourceProvisioningStateSucceeded ResourceProvisioningState = "Succeeded"
)

func PossibleValuesForResourceProvisioningState() []string {
	return []string{
		string(ResourceProvisioningStateCanceled),
		string(ResourceProvisioningStateFailed),
		string(ResourceProvisioningStateSucceeded),
	}
}

func (s *ResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceProvisioningState(input string) (*ResourceProvisioningState, error) {
	vals := map[string]ResourceProvisioningState{
		"canceled":  ResourceProvisioningStateCanceled,
		"failed":    ResourceProvisioningStateFailed,
		"succeeded": ResourceProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceProvisioningState(input)
	return &out, nil
}

type ZoneType string

const (
	ZoneTypePrimary   ZoneType = "Primary"
	ZoneTypeSecondary ZoneType = "Secondary"
)

func PossibleValuesForZoneType() []string {
	return []string{
		string(ZoneTypePrimary),
		string(ZoneTypeSecondary),
	}
}

func (s *ZoneType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseZoneType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseZoneType(input string) (*ZoneType, error) {
	vals := map[string]ZoneType{
		"primary":   ZoneTypePrimary,
		"secondary": ZoneTypeSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ZoneType(input)
	return &out, nil
}
