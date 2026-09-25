package routes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AFDEndpointProtocols string

const (
	AFDEndpointProtocolsHTTP  AFDEndpointProtocols = "Http"
	AFDEndpointProtocolsHTTPS AFDEndpointProtocols = "Https"
)

func PossibleValuesForAFDEndpointProtocols() []string {
	return []string{
		string(AFDEndpointProtocolsHTTP),
		string(AFDEndpointProtocolsHTTPS),
	}
}

func (s *AFDEndpointProtocols) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAFDEndpointProtocols(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAFDEndpointProtocols(input string) (*AFDEndpointProtocols, error) {
	vals := map[string]AFDEndpointProtocols{
		"http":  AFDEndpointProtocolsHTTP,
		"https": AFDEndpointProtocolsHTTPS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AFDEndpointProtocols(input)
	return &out, nil
}

type AfdProvisioningState string

const (
	AfdProvisioningStateCreating  AfdProvisioningState = "Creating"
	AfdProvisioningStateDeleting  AfdProvisioningState = "Deleting"
	AfdProvisioningStateFailed    AfdProvisioningState = "Failed"
	AfdProvisioningStateSucceeded AfdProvisioningState = "Succeeded"
	AfdProvisioningStateUpdating  AfdProvisioningState = "Updating"
)

func PossibleValuesForAfdProvisioningState() []string {
	return []string{
		string(AfdProvisioningStateCreating),
		string(AfdProvisioningStateDeleting),
		string(AfdProvisioningStateFailed),
		string(AfdProvisioningStateSucceeded),
		string(AfdProvisioningStateUpdating),
	}
}

func (s *AfdProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAfdProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAfdProvisioningState(input string) (*AfdProvisioningState, error) {
	vals := map[string]AfdProvisioningState{
		"creating":  AfdProvisioningStateCreating,
		"deleting":  AfdProvisioningStateDeleting,
		"failed":    AfdProvisioningStateFailed,
		"succeeded": AfdProvisioningStateSucceeded,
		"updating":  AfdProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AfdProvisioningState(input)
	return &out, nil
}

type AfdQueryStringCachingBehavior string

const (
	AfdQueryStringCachingBehaviorIgnoreQueryString            AfdQueryStringCachingBehavior = "IgnoreQueryString"
	AfdQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings  AfdQueryStringCachingBehavior = "IgnoreSpecifiedQueryStrings"
	AfdQueryStringCachingBehaviorIncludeSpecifiedQueryStrings AfdQueryStringCachingBehavior = "IncludeSpecifiedQueryStrings"
	AfdQueryStringCachingBehaviorUseQueryString               AfdQueryStringCachingBehavior = "UseQueryString"
)

func PossibleValuesForAfdQueryStringCachingBehavior() []string {
	return []string{
		string(AfdQueryStringCachingBehaviorIgnoreQueryString),
		string(AfdQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings),
		string(AfdQueryStringCachingBehaviorIncludeSpecifiedQueryStrings),
		string(AfdQueryStringCachingBehaviorUseQueryString),
	}
}

func (s *AfdQueryStringCachingBehavior) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAfdQueryStringCachingBehavior(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAfdQueryStringCachingBehavior(input string) (*AfdQueryStringCachingBehavior, error) {
	vals := map[string]AfdQueryStringCachingBehavior{
		"ignorequerystring":            AfdQueryStringCachingBehaviorIgnoreQueryString,
		"ignorespecifiedquerystrings":  AfdQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings,
		"includespecifiedquerystrings": AfdQueryStringCachingBehaviorIncludeSpecifiedQueryStrings,
		"usequerystring":               AfdQueryStringCachingBehaviorUseQueryString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AfdQueryStringCachingBehavior(input)
	return &out, nil
}

type DeploymentStatus string

const (
	DeploymentStatusFailed     DeploymentStatus = "Failed"
	DeploymentStatusInProgress DeploymentStatus = "InProgress"
	DeploymentStatusNotStarted DeploymentStatus = "NotStarted"
	DeploymentStatusSucceeded  DeploymentStatus = "Succeeded"
)

func PossibleValuesForDeploymentStatus() []string {
	return []string{
		string(DeploymentStatusFailed),
		string(DeploymentStatusInProgress),
		string(DeploymentStatusNotStarted),
		string(DeploymentStatusSucceeded),
	}
}

func (s *DeploymentStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeploymentStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeploymentStatus(input string) (*DeploymentStatus, error) {
	vals := map[string]DeploymentStatus{
		"failed":     DeploymentStatusFailed,
		"inprogress": DeploymentStatusInProgress,
		"notstarted": DeploymentStatusNotStarted,
		"succeeded":  DeploymentStatusSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentStatus(input)
	return &out, nil
}

type EnabledState string

const (
	EnabledStateDisabled EnabledState = "Disabled"
	EnabledStateEnabled  EnabledState = "Enabled"
)

func PossibleValuesForEnabledState() []string {
	return []string{
		string(EnabledStateDisabled),
		string(EnabledStateEnabled),
	}
}

func (s *EnabledState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEnabledState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEnabledState(input string) (*EnabledState, error) {
	vals := map[string]EnabledState{
		"disabled": EnabledStateDisabled,
		"enabled":  EnabledStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnabledState(input)
	return &out, nil
}

type ForwardingProtocol string

const (
	ForwardingProtocolHTTPOnly     ForwardingProtocol = "HttpOnly"
	ForwardingProtocolHTTPSOnly    ForwardingProtocol = "HttpsOnly"
	ForwardingProtocolMatchRequest ForwardingProtocol = "MatchRequest"
)

func PossibleValuesForForwardingProtocol() []string {
	return []string{
		string(ForwardingProtocolHTTPOnly),
		string(ForwardingProtocolHTTPSOnly),
		string(ForwardingProtocolMatchRequest),
	}
}

func (s *ForwardingProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseForwardingProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseForwardingProtocol(input string) (*ForwardingProtocol, error) {
	vals := map[string]ForwardingProtocol{
		"httponly":     ForwardingProtocolHTTPOnly,
		"httpsonly":    ForwardingProtocolHTTPSOnly,
		"matchrequest": ForwardingProtocolMatchRequest,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ForwardingProtocol(input)
	return &out, nil
}

type HTTPSRedirect string

const (
	HTTPSRedirectDisabled HTTPSRedirect = "Disabled"
	HTTPSRedirectEnabled  HTTPSRedirect = "Enabled"
)

func PossibleValuesForHTTPSRedirect() []string {
	return []string{
		string(HTTPSRedirectDisabled),
		string(HTTPSRedirectEnabled),
	}
}

func (s *HTTPSRedirect) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHTTPSRedirect(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHTTPSRedirect(input string) (*HTTPSRedirect, error) {
	vals := map[string]HTTPSRedirect{
		"disabled": HTTPSRedirectDisabled,
		"enabled":  HTTPSRedirectEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HTTPSRedirect(input)
	return &out, nil
}

type LinkToDefaultDomain string

const (
	LinkToDefaultDomainDisabled LinkToDefaultDomain = "Disabled"
	LinkToDefaultDomainEnabled  LinkToDefaultDomain = "Enabled"
)

func PossibleValuesForLinkToDefaultDomain() []string {
	return []string{
		string(LinkToDefaultDomainDisabled),
		string(LinkToDefaultDomainEnabled),
	}
}

func (s *LinkToDefaultDomain) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLinkToDefaultDomain(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLinkToDefaultDomain(input string) (*LinkToDefaultDomain, error) {
	vals := map[string]LinkToDefaultDomain{
		"disabled": LinkToDefaultDomainDisabled,
		"enabled":  LinkToDefaultDomainEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LinkToDefaultDomain(input)
	return &out, nil
}
