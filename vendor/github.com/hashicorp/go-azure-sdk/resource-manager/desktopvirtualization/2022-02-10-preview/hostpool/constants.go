package hostpool

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DayOfWeek string

const (
	DayOfWeekFriday    DayOfWeek = "Friday"
	DayOfWeekMonday    DayOfWeek = "Monday"
	DayOfWeekSaturday  DayOfWeek = "Saturday"
	DayOfWeekSunday    DayOfWeek = "Sunday"
	DayOfWeekThursday  DayOfWeek = "Thursday"
	DayOfWeekTuesday   DayOfWeek = "Tuesday"
	DayOfWeekWednesday DayOfWeek = "Wednesday"
)

func PossibleValuesForDayOfWeek() []string {
	return []string{
		string(DayOfWeekFriday),
		string(DayOfWeekMonday),
		string(DayOfWeekSaturday),
		string(DayOfWeekSunday),
		string(DayOfWeekThursday),
		string(DayOfWeekTuesday),
		string(DayOfWeekWednesday),
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
		"friday":    DayOfWeekFriday,
		"monday":    DayOfWeekMonday,
		"saturday":  DayOfWeekSaturday,
		"sunday":    DayOfWeekSunday,
		"thursday":  DayOfWeekThursday,
		"tuesday":   DayOfWeekTuesday,
		"wednesday": DayOfWeekWednesday,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DayOfWeek(input)
	return &out, nil
}

type HostPoolType string

const (
	HostPoolTypeBYODesktop HostPoolType = "BYODesktop"
	HostPoolTypePersonal   HostPoolType = "Personal"
	HostPoolTypePooled     HostPoolType = "Pooled"
)

func PossibleValuesForHostPoolType() []string {
	return []string{
		string(HostPoolTypeBYODesktop),
		string(HostPoolTypePersonal),
		string(HostPoolTypePooled),
	}
}

func (s *HostPoolType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHostPoolType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHostPoolType(input string) (*HostPoolType, error) {
	vals := map[string]HostPoolType{
		"byodesktop": HostPoolTypeBYODesktop,
		"personal":   HostPoolTypePersonal,
		"pooled":     HostPoolTypePooled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HostPoolType(input)
	return &out, nil
}

type HostpoolPublicNetworkAccess string

const (
	HostpoolPublicNetworkAccessDisabled                   HostpoolPublicNetworkAccess = "Disabled"
	HostpoolPublicNetworkAccessEnabled                    HostpoolPublicNetworkAccess = "Enabled"
	HostpoolPublicNetworkAccessEnabledForClientsOnly      HostpoolPublicNetworkAccess = "EnabledForClientsOnly"
	HostpoolPublicNetworkAccessEnabledForSessionHostsOnly HostpoolPublicNetworkAccess = "EnabledForSessionHostsOnly"
)

func PossibleValuesForHostpoolPublicNetworkAccess() []string {
	return []string{
		string(HostpoolPublicNetworkAccessDisabled),
		string(HostpoolPublicNetworkAccessEnabled),
		string(HostpoolPublicNetworkAccessEnabledForClientsOnly),
		string(HostpoolPublicNetworkAccessEnabledForSessionHostsOnly),
	}
}

func (s *HostpoolPublicNetworkAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHostpoolPublicNetworkAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHostpoolPublicNetworkAccess(input string) (*HostpoolPublicNetworkAccess, error) {
	vals := map[string]HostpoolPublicNetworkAccess{
		"disabled":                   HostpoolPublicNetworkAccessDisabled,
		"enabled":                    HostpoolPublicNetworkAccessEnabled,
		"enabledforclientsonly":      HostpoolPublicNetworkAccessEnabledForClientsOnly,
		"enabledforsessionhostsonly": HostpoolPublicNetworkAccessEnabledForSessionHostsOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HostpoolPublicNetworkAccess(input)
	return &out, nil
}

type LoadBalancerType string

const (
	LoadBalancerTypeBreadthFirst LoadBalancerType = "BreadthFirst"
	LoadBalancerTypeDepthFirst   LoadBalancerType = "DepthFirst"
	LoadBalancerTypePersistent   LoadBalancerType = "Persistent"
)

func PossibleValuesForLoadBalancerType() []string {
	return []string{
		string(LoadBalancerTypeBreadthFirst),
		string(LoadBalancerTypeDepthFirst),
		string(LoadBalancerTypePersistent),
	}
}

func (s *LoadBalancerType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLoadBalancerType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLoadBalancerType(input string) (*LoadBalancerType, error) {
	vals := map[string]LoadBalancerType{
		"breadthfirst": LoadBalancerTypeBreadthFirst,
		"depthfirst":   LoadBalancerTypeDepthFirst,
		"persistent":   LoadBalancerTypePersistent,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LoadBalancerType(input)
	return &out, nil
}

type Operation string

const (
	OperationComplete Operation = "Complete"
	OperationHide     Operation = "Hide"
	OperationRevoke   Operation = "Revoke"
	OperationStart    Operation = "Start"
	OperationUnhide   Operation = "Unhide"
)

func PossibleValuesForOperation() []string {
	return []string{
		string(OperationComplete),
		string(OperationHide),
		string(OperationRevoke),
		string(OperationStart),
		string(OperationUnhide),
	}
}

func (s *Operation) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperation(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperation(input string) (*Operation, error) {
	vals := map[string]Operation{
		"complete": OperationComplete,
		"hide":     OperationHide,
		"revoke":   OperationRevoke,
		"start":    OperationStart,
		"unhide":   OperationUnhide,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Operation(input)
	return &out, nil
}

type PersonalDesktopAssignmentType string

const (
	PersonalDesktopAssignmentTypeAutomatic PersonalDesktopAssignmentType = "Automatic"
	PersonalDesktopAssignmentTypeDirect    PersonalDesktopAssignmentType = "Direct"
)

func PossibleValuesForPersonalDesktopAssignmentType() []string {
	return []string{
		string(PersonalDesktopAssignmentTypeAutomatic),
		string(PersonalDesktopAssignmentTypeDirect),
	}
}

func (s *PersonalDesktopAssignmentType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePersonalDesktopAssignmentType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePersonalDesktopAssignmentType(input string) (*PersonalDesktopAssignmentType, error) {
	vals := map[string]PersonalDesktopAssignmentType{
		"automatic": PersonalDesktopAssignmentTypeAutomatic,
		"direct":    PersonalDesktopAssignmentTypeDirect,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PersonalDesktopAssignmentType(input)
	return &out, nil
}

type PreferredAppGroupType string

const (
	PreferredAppGroupTypeDesktop          PreferredAppGroupType = "Desktop"
	PreferredAppGroupTypeNone             PreferredAppGroupType = "None"
	PreferredAppGroupTypeRailApplications PreferredAppGroupType = "RailApplications"
)

func PossibleValuesForPreferredAppGroupType() []string {
	return []string{
		string(PreferredAppGroupTypeDesktop),
		string(PreferredAppGroupTypeNone),
		string(PreferredAppGroupTypeRailApplications),
	}
}

func (s *PreferredAppGroupType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePreferredAppGroupType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePreferredAppGroupType(input string) (*PreferredAppGroupType, error) {
	vals := map[string]PreferredAppGroupType{
		"desktop":          PreferredAppGroupTypeDesktop,
		"none":             PreferredAppGroupTypeNone,
		"railapplications": PreferredAppGroupTypeRailApplications,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PreferredAppGroupType(input)
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

type RegistrationTokenOperation string

const (
	RegistrationTokenOperationDelete RegistrationTokenOperation = "Delete"
	RegistrationTokenOperationNone   RegistrationTokenOperation = "None"
	RegistrationTokenOperationUpdate RegistrationTokenOperation = "Update"
)

func PossibleValuesForRegistrationTokenOperation() []string {
	return []string{
		string(RegistrationTokenOperationDelete),
		string(RegistrationTokenOperationNone),
		string(RegistrationTokenOperationUpdate),
	}
}

func (s *RegistrationTokenOperation) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRegistrationTokenOperation(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRegistrationTokenOperation(input string) (*RegistrationTokenOperation, error) {
	vals := map[string]RegistrationTokenOperation{
		"delete": RegistrationTokenOperationDelete,
		"none":   RegistrationTokenOperationNone,
		"update": RegistrationTokenOperationUpdate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RegistrationTokenOperation(input)
	return &out, nil
}

type SSOSecretType string

const (
	SSOSecretTypeCertificate           SSOSecretType = "Certificate"
	SSOSecretTypeCertificateInKeyVault SSOSecretType = "CertificateInKeyVault"
	SSOSecretTypeSharedKey             SSOSecretType = "SharedKey"
	SSOSecretTypeSharedKeyInKeyVault   SSOSecretType = "SharedKeyInKeyVault"
)

func PossibleValuesForSSOSecretType() []string {
	return []string{
		string(SSOSecretTypeCertificate),
		string(SSOSecretTypeCertificateInKeyVault),
		string(SSOSecretTypeSharedKey),
		string(SSOSecretTypeSharedKeyInKeyVault),
	}
}

func (s *SSOSecretType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSSOSecretType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSSOSecretType(input string) (*SSOSecretType, error) {
	vals := map[string]SSOSecretType{
		"certificate":           SSOSecretTypeCertificate,
		"certificateinkeyvault": SSOSecretTypeCertificateInKeyVault,
		"sharedkey":             SSOSecretTypeSharedKey,
		"sharedkeyinkeyvault":   SSOSecretTypeSharedKeyInKeyVault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SSOSecretType(input)
	return &out, nil
}

type SessionHostComponentUpdateType string

const (
	SessionHostComponentUpdateTypeDefault   SessionHostComponentUpdateType = "Default"
	SessionHostComponentUpdateTypeScheduled SessionHostComponentUpdateType = "Scheduled"
)

func PossibleValuesForSessionHostComponentUpdateType() []string {
	return []string{
		string(SessionHostComponentUpdateTypeDefault),
		string(SessionHostComponentUpdateTypeScheduled),
	}
}

func (s *SessionHostComponentUpdateType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSessionHostComponentUpdateType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSessionHostComponentUpdateType(input string) (*SessionHostComponentUpdateType, error) {
	vals := map[string]SessionHostComponentUpdateType{
		"default":   SessionHostComponentUpdateTypeDefault,
		"scheduled": SessionHostComponentUpdateTypeScheduled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SessionHostComponentUpdateType(input)
	return &out, nil
}

type SkuTier string

const (
	SkuTierBasic    SkuTier = "Basic"
	SkuTierFree     SkuTier = "Free"
	SkuTierPremium  SkuTier = "Premium"
	SkuTierStandard SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		string(SkuTierBasic),
		string(SkuTierFree),
		string(SkuTierPremium),
		string(SkuTierStandard),
	}
}

func (s *SkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"basic":    SkuTierBasic,
		"free":     SkuTierFree,
		"premium":  SkuTierPremium,
		"standard": SkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuTier(input)
	return &out, nil
}
