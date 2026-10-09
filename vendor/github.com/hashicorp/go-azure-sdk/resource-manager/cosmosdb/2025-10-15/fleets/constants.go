package fleets

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FleetspaceApiKind string

const (
	FleetspaceApiKindNoSQL FleetspaceApiKind = "NoSQL"
)

func PossibleValuesForFleetspaceApiKind() []string {
	return []string{
		string(FleetspaceApiKindNoSQL),
	}
}

func (s *FleetspaceApiKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFleetspaceApiKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFleetspaceApiKind(input string) (*FleetspaceApiKind, error) {
	vals := map[string]FleetspaceApiKind{
		"nosql": FleetspaceApiKindNoSQL,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FleetspaceApiKind(input)
	return &out, nil
}

type ServiceTier string

const (
	ServiceTierBusinessCritical ServiceTier = "BusinessCritical"
	ServiceTierGeneralPurpose   ServiceTier = "GeneralPurpose"
)

func PossibleValuesForServiceTier() []string {
	return []string{
		string(ServiceTierBusinessCritical),
		string(ServiceTierGeneralPurpose),
	}
}

func (s *ServiceTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServiceTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServiceTier(input string) (*ServiceTier, error) {
	vals := map[string]ServiceTier{
		"businesscritical": ServiceTierBusinessCritical,
		"generalpurpose":   ServiceTierGeneralPurpose,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceTier(input)
	return &out, nil
}

type Status string

const (
	StatusCanceled  Status = "Canceled"
	StatusCreating  Status = "Creating"
	StatusFailed    Status = "Failed"
	StatusSucceeded Status = "Succeeded"
	StatusUpdating  Status = "Updating"
)

func PossibleValuesForStatus() []string {
	return []string{
		string(StatusCanceled),
		string(StatusCreating),
		string(StatusFailed),
		string(StatusSucceeded),
		string(StatusUpdating),
	}
}

func (s *Status) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStatus(input string) (*Status, error) {
	vals := map[string]Status{
		"canceled":  StatusCanceled,
		"creating":  StatusCreating,
		"failed":    StatusFailed,
		"succeeded": StatusSucceeded,
		"updating":  StatusUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Status(input)
	return &out, nil
}
