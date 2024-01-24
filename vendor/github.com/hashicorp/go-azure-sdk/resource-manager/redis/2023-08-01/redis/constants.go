package redis

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessPolicyAssignmentProvisioningState string

const (
	AccessPolicyAssignmentProvisioningStateCanceled  AccessPolicyAssignmentProvisioningState = "Canceled"
	AccessPolicyAssignmentProvisioningStateDeleted   AccessPolicyAssignmentProvisioningState = "Deleted"
	AccessPolicyAssignmentProvisioningStateDeleting  AccessPolicyAssignmentProvisioningState = "Deleting"
	AccessPolicyAssignmentProvisioningStateFailed    AccessPolicyAssignmentProvisioningState = "Failed"
	AccessPolicyAssignmentProvisioningStateSucceeded AccessPolicyAssignmentProvisioningState = "Succeeded"
	AccessPolicyAssignmentProvisioningStateUpdating  AccessPolicyAssignmentProvisioningState = "Updating"
)

func PossibleValuesForAccessPolicyAssignmentProvisioningState() []string {
	return []string{
		string(AccessPolicyAssignmentProvisioningStateCanceled),
		string(AccessPolicyAssignmentProvisioningStateDeleted),
		string(AccessPolicyAssignmentProvisioningStateDeleting),
		string(AccessPolicyAssignmentProvisioningStateFailed),
		string(AccessPolicyAssignmentProvisioningStateSucceeded),
		string(AccessPolicyAssignmentProvisioningStateUpdating),
	}
}

func (s *AccessPolicyAssignmentProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccessPolicyAssignmentProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccessPolicyAssignmentProvisioningState(input string) (*AccessPolicyAssignmentProvisioningState, error) {
	vals := map[string]AccessPolicyAssignmentProvisioningState{
		"canceled":  AccessPolicyAssignmentProvisioningStateCanceled,
		"deleted":   AccessPolicyAssignmentProvisioningStateDeleted,
		"deleting":  AccessPolicyAssignmentProvisioningStateDeleting,
		"failed":    AccessPolicyAssignmentProvisioningStateFailed,
		"succeeded": AccessPolicyAssignmentProvisioningStateSucceeded,
		"updating":  AccessPolicyAssignmentProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessPolicyAssignmentProvisioningState(input)
	return &out, nil
}

type AccessPolicyProvisioningState string

const (
	AccessPolicyProvisioningStateCanceled  AccessPolicyProvisioningState = "Canceled"
	AccessPolicyProvisioningStateDeleted   AccessPolicyProvisioningState = "Deleted"
	AccessPolicyProvisioningStateDeleting  AccessPolicyProvisioningState = "Deleting"
	AccessPolicyProvisioningStateFailed    AccessPolicyProvisioningState = "Failed"
	AccessPolicyProvisioningStateSucceeded AccessPolicyProvisioningState = "Succeeded"
	AccessPolicyProvisioningStateUpdating  AccessPolicyProvisioningState = "Updating"
)

func PossibleValuesForAccessPolicyProvisioningState() []string {
	return []string{
		string(AccessPolicyProvisioningStateCanceled),
		string(AccessPolicyProvisioningStateDeleted),
		string(AccessPolicyProvisioningStateDeleting),
		string(AccessPolicyProvisioningStateFailed),
		string(AccessPolicyProvisioningStateSucceeded),
		string(AccessPolicyProvisioningStateUpdating),
	}
}

func (s *AccessPolicyProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccessPolicyProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccessPolicyProvisioningState(input string) (*AccessPolicyProvisioningState, error) {
	vals := map[string]AccessPolicyProvisioningState{
		"canceled":  AccessPolicyProvisioningStateCanceled,
		"deleted":   AccessPolicyProvisioningStateDeleted,
		"deleting":  AccessPolicyProvisioningStateDeleting,
		"failed":    AccessPolicyProvisioningStateFailed,
		"succeeded": AccessPolicyProvisioningStateSucceeded,
		"updating":  AccessPolicyProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessPolicyProvisioningState(input)
	return &out, nil
}

type AccessPolicyType string

const (
	AccessPolicyTypeBuiltIn AccessPolicyType = "BuiltIn"
	AccessPolicyTypeCustom  AccessPolicyType = "Custom"
)

func PossibleValuesForAccessPolicyType() []string {
	return []string{
		string(AccessPolicyTypeBuiltIn),
		string(AccessPolicyTypeCustom),
	}
}

func (s *AccessPolicyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccessPolicyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccessPolicyType(input string) (*AccessPolicyType, error) {
	vals := map[string]AccessPolicyType{
		"builtin": AccessPolicyTypeBuiltIn,
		"custom":  AccessPolicyTypeCustom,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessPolicyType(input)
	return &out, nil
}

type DayOfWeek string

const (
	DayOfWeekEveryday  DayOfWeek = "Everyday"
	DayOfWeekFriday    DayOfWeek = "Friday"
	DayOfWeekMonday    DayOfWeek = "Monday"
	DayOfWeekSaturday  DayOfWeek = "Saturday"
	DayOfWeekSunday    DayOfWeek = "Sunday"
	DayOfWeekThursday  DayOfWeek = "Thursday"
	DayOfWeekTuesday   DayOfWeek = "Tuesday"
	DayOfWeekWednesday DayOfWeek = "Wednesday"
	DayOfWeekWeekend   DayOfWeek = "Weekend"
)

func PossibleValuesForDayOfWeek() []string {
	return []string{
		string(DayOfWeekEveryday),
		string(DayOfWeekFriday),
		string(DayOfWeekMonday),
		string(DayOfWeekSaturday),
		string(DayOfWeekSunday),
		string(DayOfWeekThursday),
		string(DayOfWeekTuesday),
		string(DayOfWeekWednesday),
		string(DayOfWeekWeekend),
	}
}

func (s *DayOfWeek) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDayOfWeek(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDayOfWeek(input string) (*DayOfWeek, error) {
	vals := map[string]DayOfWeek{
		"everyday":  DayOfWeekEveryday,
		"friday":    DayOfWeekFriday,
		"monday":    DayOfWeekMonday,
		"saturday":  DayOfWeekSaturday,
		"sunday":    DayOfWeekSunday,
		"thursday":  DayOfWeekThursday,
		"tuesday":   DayOfWeekTuesday,
		"wednesday": DayOfWeekWednesday,
		"weekend":   DayOfWeekWeekend,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DayOfWeek(input)
	return &out, nil
}

type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCreating  PrivateEndpointConnectionProvisioningState = "Creating"
	PrivateEndpointConnectionProvisioningStateDeleting  PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateFailed    PrivateEndpointConnectionProvisioningState = "Failed"
	PrivateEndpointConnectionProvisioningStateSucceeded PrivateEndpointConnectionProvisioningState = "Succeeded"
)

func PossibleValuesForPrivateEndpointConnectionProvisioningState() []string {
	return []string{
		string(PrivateEndpointConnectionProvisioningStateCreating),
		string(PrivateEndpointConnectionProvisioningStateDeleting),
		string(PrivateEndpointConnectionProvisioningStateFailed),
		string(PrivateEndpointConnectionProvisioningStateSucceeded),
	}
}

func (s *PrivateEndpointConnectionProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateEndpointConnectionProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateEndpointConnectionProvisioningState(input string) (*PrivateEndpointConnectionProvisioningState, error) {
	vals := map[string]PrivateEndpointConnectionProvisioningState{
		"creating":  PrivateEndpointConnectionProvisioningStateCreating,
		"deleting":  PrivateEndpointConnectionProvisioningStateDeleting,
		"failed":    PrivateEndpointConnectionProvisioningStateFailed,
		"succeeded": PrivateEndpointConnectionProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointConnectionProvisioningState(input)
	return &out, nil
}

type PrivateEndpointServiceConnectionStatus string

const (
	PrivateEndpointServiceConnectionStatusApproved PrivateEndpointServiceConnectionStatus = "Approved"
	PrivateEndpointServiceConnectionStatusPending  PrivateEndpointServiceConnectionStatus = "Pending"
	PrivateEndpointServiceConnectionStatusRejected PrivateEndpointServiceConnectionStatus = "Rejected"
)

func PossibleValuesForPrivateEndpointServiceConnectionStatus() []string {
	return []string{
		string(PrivateEndpointServiceConnectionStatusApproved),
		string(PrivateEndpointServiceConnectionStatusPending),
		string(PrivateEndpointServiceConnectionStatusRejected),
	}
}

func (s *PrivateEndpointServiceConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateEndpointServiceConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateEndpointServiceConnectionStatus(input string) (*PrivateEndpointServiceConnectionStatus, error) {
	vals := map[string]PrivateEndpointServiceConnectionStatus{
		"approved": PrivateEndpointServiceConnectionStatusApproved,
		"pending":  PrivateEndpointServiceConnectionStatusPending,
		"rejected": PrivateEndpointServiceConnectionStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointServiceConnectionStatus(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateConfiguringAAD         ProvisioningState = "ConfiguringAAD"
	ProvisioningStateCreating               ProvisioningState = "Creating"
	ProvisioningStateDeleting               ProvisioningState = "Deleting"
	ProvisioningStateDisabled               ProvisioningState = "Disabled"
	ProvisioningStateFailed                 ProvisioningState = "Failed"
	ProvisioningStateLinking                ProvisioningState = "Linking"
	ProvisioningStateProvisioning           ProvisioningState = "Provisioning"
	ProvisioningStateRecoveringScaleFailure ProvisioningState = "RecoveringScaleFailure"
	ProvisioningStateScaling                ProvisioningState = "Scaling"
	ProvisioningStateSucceeded              ProvisioningState = "Succeeded"
	ProvisioningStateUnlinking              ProvisioningState = "Unlinking"
	ProvisioningStateUnprovisioning         ProvisioningState = "Unprovisioning"
	ProvisioningStateUpdating               ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateConfiguringAAD),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateDisabled),
		string(ProvisioningStateFailed),
		string(ProvisioningStateLinking),
		string(ProvisioningStateProvisioning),
		string(ProvisioningStateRecoveringScaleFailure),
		string(ProvisioningStateScaling),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUnlinking),
		string(ProvisioningStateUnprovisioning),
		string(ProvisioningStateUpdating),
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
		"configuringaad":         ProvisioningStateConfiguringAAD,
		"creating":               ProvisioningStateCreating,
		"deleting":               ProvisioningStateDeleting,
		"disabled":               ProvisioningStateDisabled,
		"failed":                 ProvisioningStateFailed,
		"linking":                ProvisioningStateLinking,
		"provisioning":           ProvisioningStateProvisioning,
		"recoveringscalefailure": ProvisioningStateRecoveringScaleFailure,
		"scaling":                ProvisioningStateScaling,
		"succeeded":              ProvisioningStateSucceeded,
		"unlinking":              ProvisioningStateUnlinking,
		"unprovisioning":         ProvisioningStateUnprovisioning,
		"updating":               ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled  PublicNetworkAccess = "Enabled"
)

func PossibleValuesForPublicNetworkAccess() []string {
	return []string{
		string(PublicNetworkAccessDisabled),
		string(PublicNetworkAccessEnabled),
	}
}

func (s *PublicNetworkAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicNetworkAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicNetworkAccess(input string) (*PublicNetworkAccess, error) {
	vals := map[string]PublicNetworkAccess{
		"disabled": PublicNetworkAccessDisabled,
		"enabled":  PublicNetworkAccessEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccess(input)
	return &out, nil
}

type RebootType string

const (
	RebootTypeAllNodes      RebootType = "AllNodes"
	RebootTypePrimaryNode   RebootType = "PrimaryNode"
	RebootTypeSecondaryNode RebootType = "SecondaryNode"
)

func PossibleValuesForRebootType() []string {
	return []string{
		string(RebootTypeAllNodes),
		string(RebootTypePrimaryNode),
		string(RebootTypeSecondaryNode),
	}
}

func (s *RebootType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRebootType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRebootType(input string) (*RebootType, error) {
	vals := map[string]RebootType{
		"allnodes":      RebootTypeAllNodes,
		"primarynode":   RebootTypePrimaryNode,
		"secondarynode": RebootTypeSecondaryNode,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RebootType(input)
	return &out, nil
}

type RedisKeyType string

const (
	RedisKeyTypePrimary   RedisKeyType = "Primary"
	RedisKeyTypeSecondary RedisKeyType = "Secondary"
)

func PossibleValuesForRedisKeyType() []string {
	return []string{
		string(RedisKeyTypePrimary),
		string(RedisKeyTypeSecondary),
	}
}

func (s *RedisKeyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRedisKeyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRedisKeyType(input string) (*RedisKeyType, error) {
	vals := map[string]RedisKeyType{
		"primary":   RedisKeyTypePrimary,
		"secondary": RedisKeyTypeSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RedisKeyType(input)
	return &out, nil
}

type ReplicationRole string

const (
	ReplicationRolePrimary   ReplicationRole = "Primary"
	ReplicationRoleSecondary ReplicationRole = "Secondary"
)

func PossibleValuesForReplicationRole() []string {
	return []string{
		string(ReplicationRolePrimary),
		string(ReplicationRoleSecondary),
	}
}

func (s *ReplicationRole) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReplicationRole(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReplicationRole(input string) (*ReplicationRole, error) {
	vals := map[string]ReplicationRole{
		"primary":   ReplicationRolePrimary,
		"secondary": ReplicationRoleSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicationRole(input)
	return &out, nil
}

type SkuFamily string

const (
	SkuFamilyC SkuFamily = "C"
	SkuFamilyP SkuFamily = "P"
)

func PossibleValuesForSkuFamily() []string {
	return []string{
		string(SkuFamilyC),
		string(SkuFamilyP),
	}
}

func (s *SkuFamily) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuFamily(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuFamily(input string) (*SkuFamily, error) {
	vals := map[string]SkuFamily{
		"c": SkuFamilyC,
		"p": SkuFamilyP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuFamily(input)
	return &out, nil
}

type SkuName string

const (
	SkuNameBasic    SkuName = "Basic"
	SkuNamePremium  SkuName = "Premium"
	SkuNameStandard SkuName = "Standard"
)

func PossibleValuesForSkuName() []string {
	return []string{
		string(SkuNameBasic),
		string(SkuNamePremium),
		string(SkuNameStandard),
	}
}

func (s *SkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuName(input string) (*SkuName, error) {
	vals := map[string]SkuName{
		"basic":    SkuNameBasic,
		"premium":  SkuNamePremium,
		"standard": SkuNameStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuName(input)
	return &out, nil
}

type TlsVersion string

const (
	TlsVersionOnePointOne  TlsVersion = "1.1"
	TlsVersionOnePointTwo  TlsVersion = "1.2"
	TlsVersionOnePointZero TlsVersion = "1.0"
)

func PossibleValuesForTlsVersion() []string {
	return []string{
		string(TlsVersionOnePointOne),
		string(TlsVersionOnePointTwo),
		string(TlsVersionOnePointZero),
	}
}

func (s *TlsVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTlsVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTlsVersion(input string) (*TlsVersion, error) {
	vals := map[string]TlsVersion{
		"1.1": TlsVersionOnePointOne,
		"1.2": TlsVersionOnePointTwo,
		"1.0": TlsVersionOnePointZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TlsVersion(input)
	return &out, nil
}

type UpdateChannel string

const (
	UpdateChannelPreview UpdateChannel = "Preview"
	UpdateChannelStable  UpdateChannel = "Stable"
)

func PossibleValuesForUpdateChannel() []string {
	return []string{
		string(UpdateChannelPreview),
		string(UpdateChannelStable),
	}
}

func (s *UpdateChannel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUpdateChannel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUpdateChannel(input string) (*UpdateChannel, error) {
	vals := map[string]UpdateChannel{
		"preview": UpdateChannelPreview,
		"stable":  UpdateChannelStable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UpdateChannel(input)
	return &out, nil
}
