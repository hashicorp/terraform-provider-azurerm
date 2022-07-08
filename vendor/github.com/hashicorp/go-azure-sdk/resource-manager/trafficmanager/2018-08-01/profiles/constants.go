package profiles

import "strings"

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

type EndpointMonitorStatus string

const (
	EndpointMonitorStatusCheckingEndpoint EndpointMonitorStatus = "CheckingEndpoint"
	EndpointMonitorStatusDegraded         EndpointMonitorStatus = "Degraded"
	EndpointMonitorStatusDisabled         EndpointMonitorStatus = "Disabled"
	EndpointMonitorStatusInactive         EndpointMonitorStatus = "Inactive"
	EndpointMonitorStatusOnline           EndpointMonitorStatus = "Online"
	EndpointMonitorStatusStopped          EndpointMonitorStatus = "Stopped"
)

func PossibleValuesForEndpointMonitorStatus() []string {
	return []string{
		string(EndpointMonitorStatusCheckingEndpoint),
		string(EndpointMonitorStatusDegraded),
		string(EndpointMonitorStatusDisabled),
		string(EndpointMonitorStatusInactive),
		string(EndpointMonitorStatusOnline),
		string(EndpointMonitorStatusStopped),
	}
}

func parseEndpointMonitorStatus(input string) (*EndpointMonitorStatus, error) {
	vals := map[string]EndpointMonitorStatus{
		"checkingendpoint": EndpointMonitorStatusCheckingEndpoint,
		"degraded":         EndpointMonitorStatusDegraded,
		"disabled":         EndpointMonitorStatusDisabled,
		"inactive":         EndpointMonitorStatusInactive,
		"online":           EndpointMonitorStatusOnline,
		"stopped":          EndpointMonitorStatusStopped,
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
