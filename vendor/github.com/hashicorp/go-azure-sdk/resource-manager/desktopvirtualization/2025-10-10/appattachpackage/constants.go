package appattachpackage

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FailHealthCheckOnStagingFailure string

const (
	FailHealthCheckOnStagingFailureDoNotFail       FailHealthCheckOnStagingFailure = "DoNotFail"
	FailHealthCheckOnStagingFailureNeedsAssistance FailHealthCheckOnStagingFailure = "NeedsAssistance"
	FailHealthCheckOnStagingFailureUnhealthy       FailHealthCheckOnStagingFailure = "Unhealthy"
)

func PossibleValuesForFailHealthCheckOnStagingFailure() []string {
	return []string{
		string(FailHealthCheckOnStagingFailureDoNotFail),
		string(FailHealthCheckOnStagingFailureNeedsAssistance),
		string(FailHealthCheckOnStagingFailureUnhealthy),
	}
}

func (s *FailHealthCheckOnStagingFailure) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFailHealthCheckOnStagingFailure(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFailHealthCheckOnStagingFailure(input string) (*FailHealthCheckOnStagingFailure, error) {
	vals := map[string]FailHealthCheckOnStagingFailure{
		"donotfail":       FailHealthCheckOnStagingFailureDoNotFail,
		"needsassistance": FailHealthCheckOnStagingFailureNeedsAssistance,
		"unhealthy":       FailHealthCheckOnStagingFailureUnhealthy,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FailHealthCheckOnStagingFailure(input)
	return &out, nil
}

type PackageTimestamped string

const (
	PackageTimestampedNotTimestamped PackageTimestamped = "NotTimestamped"
	PackageTimestampedTimestamped    PackageTimestamped = "Timestamped"
)

func PossibleValuesForPackageTimestamped() []string {
	return []string{
		string(PackageTimestampedNotTimestamped),
		string(PackageTimestampedTimestamped),
	}
}

func (s *PackageTimestamped) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePackageTimestamped(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePackageTimestamped(input string) (*PackageTimestamped, error) {
	vals := map[string]PackageTimestamped{
		"nottimestamped": PackageTimestampedNotTimestamped,
		"timestamped":    PackageTimestampedTimestamped,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PackageTimestamped(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateFailed),
		string(ProvisioningStateProvisioning),
		string(ProvisioningStateSucceeded),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"canceled":     ProvisioningStateCanceled,
		"failed":       ProvisioningStateFailed,
		"provisioning": ProvisioningStateProvisioning,
		"succeeded":    ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
