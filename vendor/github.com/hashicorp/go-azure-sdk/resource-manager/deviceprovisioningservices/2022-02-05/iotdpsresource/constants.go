package iotdpsresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessRightsDescription string

const (
	AccessRightsDescriptionDeviceConnect           AccessRightsDescription = "DeviceConnect"
	AccessRightsDescriptionEnrollmentRead          AccessRightsDescription = "EnrollmentRead"
	AccessRightsDescriptionEnrollmentWrite         AccessRightsDescription = "EnrollmentWrite"
	AccessRightsDescriptionRegistrationStatusRead  AccessRightsDescription = "RegistrationStatusRead"
	AccessRightsDescriptionRegistrationStatusWrite AccessRightsDescription = "RegistrationStatusWrite"
	AccessRightsDescriptionServiceConfig           AccessRightsDescription = "ServiceConfig"
)

func PossibleValuesForAccessRightsDescription() []string {
	return []string{
		string(AccessRightsDescriptionDeviceConnect),
		string(AccessRightsDescriptionEnrollmentRead),
		string(AccessRightsDescriptionEnrollmentWrite),
		string(AccessRightsDescriptionRegistrationStatusRead),
		string(AccessRightsDescriptionRegistrationStatusWrite),
		string(AccessRightsDescriptionServiceConfig),
	}
}

func (s *AccessRightsDescription) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccessRightsDescription(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccessRightsDescription(input string) (*AccessRightsDescription, error) {
	vals := map[string]AccessRightsDescription{
		"deviceconnect":           AccessRightsDescriptionDeviceConnect,
		"enrollmentread":          AccessRightsDescriptionEnrollmentRead,
		"enrollmentwrite":         AccessRightsDescriptionEnrollmentWrite,
		"registrationstatusread":  AccessRightsDescriptionRegistrationStatusRead,
		"registrationstatuswrite": AccessRightsDescriptionRegistrationStatusWrite,
		"serviceconfig":           AccessRightsDescriptionServiceConfig,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessRightsDescription(input)
	return &out, nil
}

type AllocationPolicy string

const (
	AllocationPolicyGeoLatency AllocationPolicy = "GeoLatency"
	AllocationPolicyHashed     AllocationPolicy = "Hashed"
	AllocationPolicyStatic     AllocationPolicy = "Static"
)

func PossibleValuesForAllocationPolicy() []string {
	return []string{
		string(AllocationPolicyGeoLatency),
		string(AllocationPolicyHashed),
		string(AllocationPolicyStatic),
	}
}

func (s *AllocationPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAllocationPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAllocationPolicy(input string) (*AllocationPolicy, error) {
	vals := map[string]AllocationPolicy{
		"geolatency": AllocationPolicyGeoLatency,
		"hashed":     AllocationPolicyHashed,
		"static":     AllocationPolicyStatic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AllocationPolicy(input)
	return &out, nil
}

type IPFilterActionType string

const (
	IPFilterActionTypeAccept IPFilterActionType = "Accept"
	IPFilterActionTypeReject IPFilterActionType = "Reject"
)

func PossibleValuesForIPFilterActionType() []string {
	return []string{
		string(IPFilterActionTypeAccept),
		string(IPFilterActionTypeReject),
	}
}

func (s *IPFilterActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIPFilterActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIPFilterActionType(input string) (*IPFilterActionType, error) {
	vals := map[string]IPFilterActionType{
		"accept": IPFilterActionTypeAccept,
		"reject": IPFilterActionTypeReject,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IPFilterActionType(input)
	return &out, nil
}

type IPFilterTargetType string

const (
	IPFilterTargetTypeAll        IPFilterTargetType = "all"
	IPFilterTargetTypeDeviceApi  IPFilterTargetType = "deviceApi"
	IPFilterTargetTypeServiceApi IPFilterTargetType = "serviceApi"
)

func PossibleValuesForIPFilterTargetType() []string {
	return []string{
		string(IPFilterTargetTypeAll),
		string(IPFilterTargetTypeDeviceApi),
		string(IPFilterTargetTypeServiceApi),
	}
}

func (s *IPFilterTargetType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIPFilterTargetType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIPFilterTargetType(input string) (*IPFilterTargetType, error) {
	vals := map[string]IPFilterTargetType{
		"all":        IPFilterTargetTypeAll,
		"deviceapi":  IPFilterTargetTypeDeviceApi,
		"serviceapi": IPFilterTargetTypeServiceApi,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IPFilterTargetType(input)
	return &out, nil
}

type IotDpsSku string

const (
	IotDpsSkuSOne IotDpsSku = "S1"
)

func PossibleValuesForIotDpsSku() []string {
	return []string{
		string(IotDpsSkuSOne),
	}
}

func (s *IotDpsSku) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIotDpsSku(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIotDpsSku(input string) (*IotDpsSku, error) {
	vals := map[string]IotDpsSku{
		"s1": IotDpsSkuSOne,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IotDpsSku(input)
	return &out, nil
}

type NameUnavailabilityReason string

const (
	NameUnavailabilityReasonAlreadyExists NameUnavailabilityReason = "AlreadyExists"
	NameUnavailabilityReasonInvalid       NameUnavailabilityReason = "Invalid"
)

func PossibleValuesForNameUnavailabilityReason() []string {
	return []string{
		string(NameUnavailabilityReasonAlreadyExists),
		string(NameUnavailabilityReasonInvalid),
	}
}

func (s *NameUnavailabilityReason) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNameUnavailabilityReason(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNameUnavailabilityReason(input string) (*NameUnavailabilityReason, error) {
	vals := map[string]NameUnavailabilityReason{
		"alreadyexists": NameUnavailabilityReasonAlreadyExists,
		"invalid":       NameUnavailabilityReasonInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NameUnavailabilityReason(input)
	return &out, nil
}

type PrivateLinkServiceConnectionStatus string

const (
	PrivateLinkServiceConnectionStatusApproved     PrivateLinkServiceConnectionStatus = "Approved"
	PrivateLinkServiceConnectionStatusDisconnected PrivateLinkServiceConnectionStatus = "Disconnected"
	PrivateLinkServiceConnectionStatusPending      PrivateLinkServiceConnectionStatus = "Pending"
	PrivateLinkServiceConnectionStatusRejected     PrivateLinkServiceConnectionStatus = "Rejected"
)

func PossibleValuesForPrivateLinkServiceConnectionStatus() []string {
	return []string{
		string(PrivateLinkServiceConnectionStatusApproved),
		string(PrivateLinkServiceConnectionStatusDisconnected),
		string(PrivateLinkServiceConnectionStatusPending),
		string(PrivateLinkServiceConnectionStatusRejected),
	}
}

func (s *PrivateLinkServiceConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateLinkServiceConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateLinkServiceConnectionStatus(input string) (*PrivateLinkServiceConnectionStatus, error) {
	vals := map[string]PrivateLinkServiceConnectionStatus{
		"approved":     PrivateLinkServiceConnectionStatusApproved,
		"disconnected": PrivateLinkServiceConnectionStatusDisconnected,
		"pending":      PrivateLinkServiceConnectionStatusPending,
		"rejected":     PrivateLinkServiceConnectionStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateLinkServiceConnectionStatus(input)
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

type State string

const (
	StateActivating       State = "Activating"
	StateActivationFailed State = "ActivationFailed"
	StateActive           State = "Active"
	StateDeleted          State = "Deleted"
	StateDeleting         State = "Deleting"
	StateDeletionFailed   State = "DeletionFailed"
	StateFailingOver      State = "FailingOver"
	StateFailoverFailed   State = "FailoverFailed"
	StateResuming         State = "Resuming"
	StateSuspended        State = "Suspended"
	StateSuspending       State = "Suspending"
	StateTransitioning    State = "Transitioning"
)

func PossibleValuesForState() []string {
	return []string{
		string(StateActivating),
		string(StateActivationFailed),
		string(StateActive),
		string(StateDeleted),
		string(StateDeleting),
		string(StateDeletionFailed),
		string(StateFailingOver),
		string(StateFailoverFailed),
		string(StateResuming),
		string(StateSuspended),
		string(StateSuspending),
		string(StateTransitioning),
	}
}

func (s *State) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseState(input string) (*State, error) {
	vals := map[string]State{
		"activating":       StateActivating,
		"activationfailed": StateActivationFailed,
		"active":           StateActive,
		"deleted":          StateDeleted,
		"deleting":         StateDeleting,
		"deletionfailed":   StateDeletionFailed,
		"failingover":      StateFailingOver,
		"failoverfailed":   StateFailoverFailed,
		"resuming":         StateResuming,
		"suspended":        StateSuspended,
		"suspending":       StateSuspending,
		"transitioning":    StateTransitioning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := State(input)
	return &out, nil
}
