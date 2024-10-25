package blobcontainers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImmutabilityPolicyState string

const (
	ImmutabilityPolicyStateLocked   ImmutabilityPolicyState = "Locked"
	ImmutabilityPolicyStateUnlocked ImmutabilityPolicyState = "Unlocked"
)

func PossibleValuesForImmutabilityPolicyState() []string {
	return []string{
		string(ImmutabilityPolicyStateLocked),
		string(ImmutabilityPolicyStateUnlocked),
	}
}

func (s *ImmutabilityPolicyState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseImmutabilityPolicyState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseImmutabilityPolicyState(input string) (*ImmutabilityPolicyState, error) {
	vals := map[string]ImmutabilityPolicyState{
		"locked":   ImmutabilityPolicyStateLocked,
		"unlocked": ImmutabilityPolicyStateUnlocked,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ImmutabilityPolicyState(input)
	return &out, nil
}

type ImmutabilityPolicyUpdateType string

const (
	ImmutabilityPolicyUpdateTypeExtend ImmutabilityPolicyUpdateType = "extend"
	ImmutabilityPolicyUpdateTypeLock   ImmutabilityPolicyUpdateType = "lock"
	ImmutabilityPolicyUpdateTypePut    ImmutabilityPolicyUpdateType = "put"
)

func PossibleValuesForImmutabilityPolicyUpdateType() []string {
	return []string{
		string(ImmutabilityPolicyUpdateTypeExtend),
		string(ImmutabilityPolicyUpdateTypeLock),
		string(ImmutabilityPolicyUpdateTypePut),
	}
}

func (s *ImmutabilityPolicyUpdateType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseImmutabilityPolicyUpdateType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseImmutabilityPolicyUpdateType(input string) (*ImmutabilityPolicyUpdateType, error) {
	vals := map[string]ImmutabilityPolicyUpdateType{
		"extend": ImmutabilityPolicyUpdateTypeExtend,
		"lock":   ImmutabilityPolicyUpdateTypeLock,
		"put":    ImmutabilityPolicyUpdateTypePut,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ImmutabilityPolicyUpdateType(input)
	return &out, nil
}

type LeaseContainerRequestAction string

const (
	LeaseContainerRequestActionAcquire LeaseContainerRequestAction = "Acquire"
	LeaseContainerRequestActionBreak   LeaseContainerRequestAction = "Break"
	LeaseContainerRequestActionChange  LeaseContainerRequestAction = "Change"
	LeaseContainerRequestActionRelease LeaseContainerRequestAction = "Release"
	LeaseContainerRequestActionRenew   LeaseContainerRequestAction = "Renew"
)

func PossibleValuesForLeaseContainerRequestAction() []string {
	return []string{
		string(LeaseContainerRequestActionAcquire),
		string(LeaseContainerRequestActionBreak),
		string(LeaseContainerRequestActionChange),
		string(LeaseContainerRequestActionRelease),
		string(LeaseContainerRequestActionRenew),
	}
}

func (s *LeaseContainerRequestAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLeaseContainerRequestAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLeaseContainerRequestAction(input string) (*LeaseContainerRequestAction, error) {
	vals := map[string]LeaseContainerRequestAction{
		"acquire": LeaseContainerRequestActionAcquire,
		"break":   LeaseContainerRequestActionBreak,
		"change":  LeaseContainerRequestActionChange,
		"release": LeaseContainerRequestActionRelease,
		"renew":   LeaseContainerRequestActionRenew,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LeaseContainerRequestAction(input)
	return &out, nil
}

type LeaseDuration string

const (
	LeaseDurationFixed    LeaseDuration = "Fixed"
	LeaseDurationInfinite LeaseDuration = "Infinite"
)

func PossibleValuesForLeaseDuration() []string {
	return []string{
		string(LeaseDurationFixed),
		string(LeaseDurationInfinite),
	}
}

func (s *LeaseDuration) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLeaseDuration(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLeaseDuration(input string) (*LeaseDuration, error) {
	vals := map[string]LeaseDuration{
		"fixed":    LeaseDurationFixed,
		"infinite": LeaseDurationInfinite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LeaseDuration(input)
	return &out, nil
}

type LeaseState string

const (
	LeaseStateAvailable LeaseState = "Available"
	LeaseStateBreaking  LeaseState = "Breaking"
	LeaseStateBroken    LeaseState = "Broken"
	LeaseStateExpired   LeaseState = "Expired"
	LeaseStateLeased    LeaseState = "Leased"
)

func PossibleValuesForLeaseState() []string {
	return []string{
		string(LeaseStateAvailable),
		string(LeaseStateBreaking),
		string(LeaseStateBroken),
		string(LeaseStateExpired),
		string(LeaseStateLeased),
	}
}

func (s *LeaseState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLeaseState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLeaseState(input string) (*LeaseState, error) {
	vals := map[string]LeaseState{
		"available": LeaseStateAvailable,
		"breaking":  LeaseStateBreaking,
		"broken":    LeaseStateBroken,
		"expired":   LeaseStateExpired,
		"leased":    LeaseStateLeased,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LeaseState(input)
	return &out, nil
}

type LeaseStatus string

const (
	LeaseStatusLocked   LeaseStatus = "Locked"
	LeaseStatusUnlocked LeaseStatus = "Unlocked"
)

func PossibleValuesForLeaseStatus() []string {
	return []string{
		string(LeaseStatusLocked),
		string(LeaseStatusUnlocked),
	}
}

func (s *LeaseStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLeaseStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLeaseStatus(input string) (*LeaseStatus, error) {
	vals := map[string]LeaseStatus{
		"locked":   LeaseStatusLocked,
		"unlocked": LeaseStatusUnlocked,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LeaseStatus(input)
	return &out, nil
}

type ListContainersInclude string

const (
	ListContainersIncludeDeleted ListContainersInclude = "deleted"
)

func PossibleValuesForListContainersInclude() []string {
	return []string{
		string(ListContainersIncludeDeleted),
	}
}

func (s *ListContainersInclude) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseListContainersInclude(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseListContainersInclude(input string) (*ListContainersInclude, error) {
	vals := map[string]ListContainersInclude{
		"deleted": ListContainersIncludeDeleted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ListContainersInclude(input)
	return &out, nil
}

type MigrationState string

const (
	MigrationStateCompleted  MigrationState = "Completed"
	MigrationStateInProgress MigrationState = "InProgress"
)

func PossibleValuesForMigrationState() []string {
	return []string{
		string(MigrationStateCompleted),
		string(MigrationStateInProgress),
	}
}

func (s *MigrationState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMigrationState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMigrationState(input string) (*MigrationState, error) {
	vals := map[string]MigrationState{
		"completed":  MigrationStateCompleted,
		"inprogress": MigrationStateInProgress,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MigrationState(input)
	return &out, nil
}

type PublicAccess string

const (
	PublicAccessBlob      PublicAccess = "Blob"
	PublicAccessContainer PublicAccess = "Container"
	PublicAccessNone      PublicAccess = "None"
)

func PossibleValuesForPublicAccess() []string {
	return []string{
		string(PublicAccessBlob),
		string(PublicAccessContainer),
		string(PublicAccessNone),
	}
}

func (s *PublicAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicAccess(input string) (*PublicAccess, error) {
	vals := map[string]PublicAccess{
		"blob":      PublicAccessBlob,
		"container": PublicAccessContainer,
		"none":      PublicAccessNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicAccess(input)
	return &out, nil
}
