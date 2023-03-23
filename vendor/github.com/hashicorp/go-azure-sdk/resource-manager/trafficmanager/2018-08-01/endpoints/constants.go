package endpoints

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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

type EndpointType string

const (
	EndpointTypeAzureEndpoints    EndpointType = "AzureEndpoints"
	EndpointTypeExternalEndpoints EndpointType = "ExternalEndpoints"
	EndpointTypeNestedEndpoints   EndpointType = "NestedEndpoints"
)

func PossibleValuesForEndpointType() []string {
	return []string{
		string(EndpointTypeAzureEndpoints),
		string(EndpointTypeExternalEndpoints),
		string(EndpointTypeNestedEndpoints),
	}
}

func parseEndpointType(input string) (*EndpointType, error) {
	vals := map[string]EndpointType{
		"azureendpoints":    EndpointTypeAzureEndpoints,
		"externalendpoints": EndpointTypeExternalEndpoints,
		"nestedendpoints":   EndpointTypeNestedEndpoints,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EndpointType(input)
	return &out, nil
}
