package endpoints

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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

func (s *EndpointType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEndpointType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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
