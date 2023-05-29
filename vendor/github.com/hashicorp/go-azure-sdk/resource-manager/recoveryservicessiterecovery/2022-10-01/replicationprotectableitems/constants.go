package replicationprotectableitems

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HealthErrorCustomerResolvability string

const (
	HealthErrorCustomerResolvabilityAllowed    HealthErrorCustomerResolvability = "Allowed"
	HealthErrorCustomerResolvabilityNotAllowed HealthErrorCustomerResolvability = "NotAllowed"
)

func PossibleValuesForHealthErrorCustomerResolvability() []string {
	return []string{
		string(HealthErrorCustomerResolvabilityAllowed),
		string(HealthErrorCustomerResolvabilityNotAllowed),
	}
}

func parseHealthErrorCustomerResolvability(input string) (*HealthErrorCustomerResolvability, error) {
	vals := map[string]HealthErrorCustomerResolvability{
		"allowed":    HealthErrorCustomerResolvabilityAllowed,
		"notallowed": HealthErrorCustomerResolvabilityNotAllowed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HealthErrorCustomerResolvability(input)
	return &out, nil
}

type PresenceStatus string

const (
	PresenceStatusNotPresent PresenceStatus = "NotPresent"
	PresenceStatusPresent    PresenceStatus = "Present"
	PresenceStatusUnknown    PresenceStatus = "Unknown"
)

func PossibleValuesForPresenceStatus() []string {
	return []string{
		string(PresenceStatusNotPresent),
		string(PresenceStatusPresent),
		string(PresenceStatusUnknown),
	}
}

func parsePresenceStatus(input string) (*PresenceStatus, error) {
	vals := map[string]PresenceStatus{
		"notpresent": PresenceStatusNotPresent,
		"present":    PresenceStatusPresent,
		"unknown":    PresenceStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PresenceStatus(input)
	return &out, nil
}
