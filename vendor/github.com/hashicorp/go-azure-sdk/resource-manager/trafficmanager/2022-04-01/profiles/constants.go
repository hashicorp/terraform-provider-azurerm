package profiles

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AllowedEndpointRecordType string

const (
	AllowedEndpointRecordTypeAny            AllowedEndpointRecordType = "Any"
	AllowedEndpointRecordTypeDomainName     AllowedEndpointRecordType = "DomainName"
	AllowedEndpointRecordTypeIPvFourAddress AllowedEndpointRecordType = "IPv4Address"
	AllowedEndpointRecordTypeIPvSixAddress  AllowedEndpointRecordType = "IPv6Address"
)

func PossibleValuesForAllowedEndpointRecordType() []string {
	return []string{
		string(AllowedEndpointRecordTypeAny),
		string(AllowedEndpointRecordTypeDomainName),
		string(AllowedEndpointRecordTypeIPvFourAddress),
		string(AllowedEndpointRecordTypeIPvSixAddress),
	}
}

func (s *AllowedEndpointRecordType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAllowedEndpointRecordType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAllowedEndpointRecordType(input string) (*AllowedEndpointRecordType, error) {
	vals := map[string]AllowedEndpointRecordType{
		"any":         AllowedEndpointRecordTypeAny,
		"domainname":  AllowedEndpointRecordTypeDomainName,
		"ipv4address": AllowedEndpointRecordTypeIPvFourAddress,
		"ipv6address": AllowedEndpointRecordTypeIPvSixAddress,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AllowedEndpointRecordType(input)
	return &out, nil
}

type AlwaysServe string

const (
	AlwaysServeDisabled AlwaysServe = "Disabled"
	AlwaysServeEnabled  AlwaysServe = "Enabled"
)

func PossibleValuesForAlwaysServe() []string {
	return []string{
		string(AlwaysServeDisabled),
		string(AlwaysServeEnabled),
	}
}

func (s *AlwaysServe) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAlwaysServe(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAlwaysServe(input string) (*AlwaysServe, error) {
	vals := map[string]AlwaysServe{
		"disabled": AlwaysServeDisabled,
		"enabled":  AlwaysServeEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AlwaysServe(input)
	return &out, nil
}

type EndpointMonitorStatus string

const (
	EndpointMonitorStatusCheckingEndpoint EndpointMonitorStatus = "CheckingEndpoint"
	EndpointMonitorStatusDegraded         EndpointMonitorStatus = "Degraded"
	EndpointMonitorStatusDisabled         EndpointMonitorStatus = "Disabled"
	EndpointMonitorStatusInactive         EndpointMonitorStatus = "Inactive"
	EndpointMonitorStatusOnline           EndpointMonitorStatus = "Online"
	EndpointMonitorStatusStopped          EndpointMonitorStatus = "Stopped"
	EndpointMonitorStatusUnmonitored      EndpointMonitorStatus = "Unmonitored"
)

func PossibleValuesForEndpointMonitorStatus() []string {
	return []string{
		string(EndpointMonitorStatusCheckingEndpoint),
		string(EndpointMonitorStatusDegraded),
		string(EndpointMonitorStatusDisabled),
		string(EndpointMonitorStatusInactive),
		string(EndpointMonitorStatusOnline),
		string(EndpointMonitorStatusStopped),
		string(EndpointMonitorStatusUnmonitored),
	}
}

func (s *EndpointMonitorStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEndpointMonitorStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEndpointMonitorStatus(input string) (*EndpointMonitorStatus, error) {
	vals := map[string]EndpointMonitorStatus{
		"checkingendpoint": EndpointMonitorStatusCheckingEndpoint,
		"degraded":         EndpointMonitorStatusDegraded,
		"disabled":         EndpointMonitorStatusDisabled,
		"inactive":         EndpointMonitorStatusInactive,
		"online":           EndpointMonitorStatusOnline,
		"stopped":          EndpointMonitorStatusStopped,
		"unmonitored":      EndpointMonitorStatusUnmonitored,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EndpointMonitorStatus(input)
	return &out, nil
}

type EndpointStatus string

const (
	EndpointStatusDisabled EndpointStatus = "Disabled"
	EndpointStatusEnabled  EndpointStatus = "Enabled"
)

func PossibleValuesForEndpointStatus() []string {
	return []string{
		string(EndpointStatusDisabled),
		string(EndpointStatusEnabled),
	}
}

func (s *EndpointStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEndpointStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEndpointStatus(input string) (*EndpointStatus, error) {
	vals := map[string]EndpointStatus{
		"disabled": EndpointStatusDisabled,
		"enabled":  EndpointStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EndpointStatus(input)
	return &out, nil
}

type MonitorProtocol string

const (
	MonitorProtocolHTTP  MonitorProtocol = "HTTP"
	MonitorProtocolHTTPS MonitorProtocol = "HTTPS"
	MonitorProtocolTCP   MonitorProtocol = "TCP"
)

func PossibleValuesForMonitorProtocol() []string {
	return []string{
		string(MonitorProtocolHTTP),
		string(MonitorProtocolHTTPS),
		string(MonitorProtocolTCP),
	}
}

func (s *MonitorProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMonitorProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMonitorProtocol(input string) (*MonitorProtocol, error) {
	vals := map[string]MonitorProtocol{
		"http":  MonitorProtocolHTTP,
		"https": MonitorProtocolHTTPS,
		"tcp":   MonitorProtocolTCP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MonitorProtocol(input)
	return &out, nil
}

type ProfileMonitorStatus string

const (
	ProfileMonitorStatusCheckingEndpoints ProfileMonitorStatus = "CheckingEndpoints"
	ProfileMonitorStatusDegraded          ProfileMonitorStatus = "Degraded"
	ProfileMonitorStatusDisabled          ProfileMonitorStatus = "Disabled"
	ProfileMonitorStatusInactive          ProfileMonitorStatus = "Inactive"
	ProfileMonitorStatusOnline            ProfileMonitorStatus = "Online"
)

func PossibleValuesForProfileMonitorStatus() []string {
	return []string{
		string(ProfileMonitorStatusCheckingEndpoints),
		string(ProfileMonitorStatusDegraded),
		string(ProfileMonitorStatusDisabled),
		string(ProfileMonitorStatusInactive),
		string(ProfileMonitorStatusOnline),
	}
}

func (s *ProfileMonitorStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProfileMonitorStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProfileMonitorStatus(input string) (*ProfileMonitorStatus, error) {
	vals := map[string]ProfileMonitorStatus{
		"checkingendpoints": ProfileMonitorStatusCheckingEndpoints,
		"degraded":          ProfileMonitorStatusDegraded,
		"disabled":          ProfileMonitorStatusDisabled,
		"inactive":          ProfileMonitorStatusInactive,
		"online":            ProfileMonitorStatusOnline,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProfileMonitorStatus(input)
	return &out, nil
}

type ProfileStatus string

const (
	ProfileStatusDisabled ProfileStatus = "Disabled"
	ProfileStatusEnabled  ProfileStatus = "Enabled"
)

func PossibleValuesForProfileStatus() []string {
	return []string{
		string(ProfileStatusDisabled),
		string(ProfileStatusEnabled),
	}
}

func (s *ProfileStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProfileStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProfileStatus(input string) (*ProfileStatus, error) {
	vals := map[string]ProfileStatus{
		"disabled": ProfileStatusDisabled,
		"enabled":  ProfileStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProfileStatus(input)
	return &out, nil
}

type TrafficRoutingMethod string

const (
	TrafficRoutingMethodGeographic  TrafficRoutingMethod = "Geographic"
	TrafficRoutingMethodMultiValue  TrafficRoutingMethod = "MultiValue"
	TrafficRoutingMethodPerformance TrafficRoutingMethod = "Performance"
	TrafficRoutingMethodPriority    TrafficRoutingMethod = "Priority"
	TrafficRoutingMethodSubnet      TrafficRoutingMethod = "Subnet"
	TrafficRoutingMethodWeighted    TrafficRoutingMethod = "Weighted"
)

func PossibleValuesForTrafficRoutingMethod() []string {
	return []string{
		string(TrafficRoutingMethodGeographic),
		string(TrafficRoutingMethodMultiValue),
		string(TrafficRoutingMethodPerformance),
		string(TrafficRoutingMethodPriority),
		string(TrafficRoutingMethodSubnet),
		string(TrafficRoutingMethodWeighted),
	}
}

func (s *TrafficRoutingMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTrafficRoutingMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTrafficRoutingMethod(input string) (*TrafficRoutingMethod, error) {
	vals := map[string]TrafficRoutingMethod{
		"geographic":  TrafficRoutingMethodGeographic,
		"multivalue":  TrafficRoutingMethodMultiValue,
		"performance": TrafficRoutingMethodPerformance,
		"priority":    TrafficRoutingMethodPriority,
		"subnet":      TrafficRoutingMethodSubnet,
		"weighted":    TrafficRoutingMethodWeighted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TrafficRoutingMethod(input)
	return &out, nil
}

type TrafficViewEnrollmentStatus string

const (
	TrafficViewEnrollmentStatusDisabled TrafficViewEnrollmentStatus = "Disabled"
	TrafficViewEnrollmentStatusEnabled  TrafficViewEnrollmentStatus = "Enabled"
)

func PossibleValuesForTrafficViewEnrollmentStatus() []string {
	return []string{
		string(TrafficViewEnrollmentStatusDisabled),
		string(TrafficViewEnrollmentStatusEnabled),
	}
}

func (s *TrafficViewEnrollmentStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTrafficViewEnrollmentStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTrafficViewEnrollmentStatus(input string) (*TrafficViewEnrollmentStatus, error) {
	vals := map[string]TrafficViewEnrollmentStatus{
		"disabled": TrafficViewEnrollmentStatusDisabled,
		"enabled":  TrafficViewEnrollmentStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TrafficViewEnrollmentStatus(input)
	return &out, nil
}
