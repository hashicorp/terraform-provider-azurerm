package sessionhost

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HealthCheckName string

const (
	HealthCheckNameAppAttachHealthCheck     HealthCheckName = "AppAttachHealthCheck"
	HealthCheckNameDomainJoinedCheck        HealthCheckName = "DomainJoinedCheck"
	HealthCheckNameDomainReachable          HealthCheckName = "DomainReachable"
	HealthCheckNameDomainTrustCheck         HealthCheckName = "DomainTrustCheck"
	HealthCheckNameFSLogixHealthCheck       HealthCheckName = "FSLogixHealthCheck"
	HealthCheckNameMetaDataServiceCheck     HealthCheckName = "MetaDataServiceCheck"
	HealthCheckNameMonitoringAgentCheck     HealthCheckName = "MonitoringAgentCheck"
	HealthCheckNameSupportedEncryptionCheck HealthCheckName = "SupportedEncryptionCheck"
	HealthCheckNameSxSStackListenerCheck    HealthCheckName = "SxSStackListenerCheck"
	HealthCheckNameURLsAccessibleCheck      HealthCheckName = "UrlsAccessibleCheck"
	HealthCheckNameWebRTCRedirectorCheck    HealthCheckName = "WebRTCRedirectorCheck"
)

func PossibleValuesForHealthCheckName() []string {
	return []string{
		string(HealthCheckNameAppAttachHealthCheck),
		string(HealthCheckNameDomainJoinedCheck),
		string(HealthCheckNameDomainReachable),
		string(HealthCheckNameDomainTrustCheck),
		string(HealthCheckNameFSLogixHealthCheck),
		string(HealthCheckNameMetaDataServiceCheck),
		string(HealthCheckNameMonitoringAgentCheck),
		string(HealthCheckNameSupportedEncryptionCheck),
		string(HealthCheckNameSxSStackListenerCheck),
		string(HealthCheckNameURLsAccessibleCheck),
		string(HealthCheckNameWebRTCRedirectorCheck),
	}
}

func (s *HealthCheckName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHealthCheckName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHealthCheckName(input string) (*HealthCheckName, error) {
	vals := map[string]HealthCheckName{
		"appattachhealthcheck":     HealthCheckNameAppAttachHealthCheck,
		"domainjoinedcheck":        HealthCheckNameDomainJoinedCheck,
		"domainreachable":          HealthCheckNameDomainReachable,
		"domaintrustcheck":         HealthCheckNameDomainTrustCheck,
		"fslogixhealthcheck":       HealthCheckNameFSLogixHealthCheck,
		"metadataservicecheck":     HealthCheckNameMetaDataServiceCheck,
		"monitoringagentcheck":     HealthCheckNameMonitoringAgentCheck,
		"supportedencryptioncheck": HealthCheckNameSupportedEncryptionCheck,
		"sxsstacklistenercheck":    HealthCheckNameSxSStackListenerCheck,
		"urlsaccessiblecheck":      HealthCheckNameURLsAccessibleCheck,
		"webrtcredirectorcheck":    HealthCheckNameWebRTCRedirectorCheck,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HealthCheckName(input)
	return &out, nil
}

type HealthCheckResult string

const (
	HealthCheckResultHealthCheckFailed    HealthCheckResult = "HealthCheckFailed"
	HealthCheckResultHealthCheckSucceeded HealthCheckResult = "HealthCheckSucceeded"
	HealthCheckResultSessionHostShutdown  HealthCheckResult = "SessionHostShutdown"
	HealthCheckResultUnknown              HealthCheckResult = "Unknown"
)

func PossibleValuesForHealthCheckResult() []string {
	return []string{
		string(HealthCheckResultHealthCheckFailed),
		string(HealthCheckResultHealthCheckSucceeded),
		string(HealthCheckResultSessionHostShutdown),
		string(HealthCheckResultUnknown),
	}
}

func (s *HealthCheckResult) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHealthCheckResult(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHealthCheckResult(input string) (*HealthCheckResult, error) {
	vals := map[string]HealthCheckResult{
		"healthcheckfailed":    HealthCheckResultHealthCheckFailed,
		"healthchecksucceeded": HealthCheckResultHealthCheckSucceeded,
		"sessionhostshutdown":  HealthCheckResultSessionHostShutdown,
		"unknown":              HealthCheckResultUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HealthCheckResult(input)
	return &out, nil
}

type Status string

const (
	StatusAvailable                   Status = "Available"
	StatusDisconnected                Status = "Disconnected"
	StatusDomainTrustRelationshipLost Status = "DomainTrustRelationshipLost"
	StatusFSLogixNotHealthy           Status = "FSLogixNotHealthy"
	StatusNeedsAssistance             Status = "NeedsAssistance"
	StatusNoHeartbeat                 Status = "NoHeartbeat"
	StatusNotJoinedToDomain           Status = "NotJoinedToDomain"
	StatusShutdown                    Status = "Shutdown"
	StatusSxSStackListenerNotReady    Status = "SxSStackListenerNotReady"
	StatusUnavailable                 Status = "Unavailable"
	StatusUpgradeFailed               Status = "UpgradeFailed"
	StatusUpgrading                   Status = "Upgrading"
)

func PossibleValuesForStatus() []string {
	return []string{
		string(StatusAvailable),
		string(StatusDisconnected),
		string(StatusDomainTrustRelationshipLost),
		string(StatusFSLogixNotHealthy),
		string(StatusNeedsAssistance),
		string(StatusNoHeartbeat),
		string(StatusNotJoinedToDomain),
		string(StatusShutdown),
		string(StatusSxSStackListenerNotReady),
		string(StatusUnavailable),
		string(StatusUpgradeFailed),
		string(StatusUpgrading),
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
		"available":                   StatusAvailable,
		"disconnected":                StatusDisconnected,
		"domaintrustrelationshiplost": StatusDomainTrustRelationshipLost,
		"fslogixnothealthy":           StatusFSLogixNotHealthy,
		"needsassistance":             StatusNeedsAssistance,
		"noheartbeat":                 StatusNoHeartbeat,
		"notjoinedtodomain":           StatusNotJoinedToDomain,
		"shutdown":                    StatusShutdown,
		"sxsstacklistenernotready":    StatusSxSStackListenerNotReady,
		"unavailable":                 StatusUnavailable,
		"upgradefailed":               StatusUpgradeFailed,
		"upgrading":                   StatusUpgrading,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Status(input)
	return &out, nil
}

type UpdateState string

const (
	UpdateStateFailed    UpdateState = "Failed"
	UpdateStateInitial   UpdateState = "Initial"
	UpdateStatePending   UpdateState = "Pending"
	UpdateStateStarted   UpdateState = "Started"
	UpdateStateSucceeded UpdateState = "Succeeded"
)

func PossibleValuesForUpdateState() []string {
	return []string{
		string(UpdateStateFailed),
		string(UpdateStateInitial),
		string(UpdateStatePending),
		string(UpdateStateStarted),
		string(UpdateStateSucceeded),
	}
}

func (s *UpdateState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUpdateState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUpdateState(input string) (*UpdateState, error) {
	vals := map[string]UpdateState{
		"failed":    UpdateStateFailed,
		"initial":   UpdateStateInitial,
		"pending":   UpdateStatePending,
		"started":   UpdateStateStarted,
		"succeeded": UpdateStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UpdateState(input)
	return &out, nil
}
