package connectedclusters

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthenticationMethod string

const (
	AuthenticationMethodAAD   AuthenticationMethod = "AAD"
	AuthenticationMethodToken AuthenticationMethod = "Token"
)

func PossibleValuesForAuthenticationMethod() []string {
	return []string{
		string(AuthenticationMethodAAD),
		string(AuthenticationMethodToken),
	}
}

func (s *AuthenticationMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAuthenticationMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAuthenticationMethod(input string) (*AuthenticationMethod, error) {
	vals := map[string]AuthenticationMethod{
		"aad":   AuthenticationMethodAAD,
		"token": AuthenticationMethodToken,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthenticationMethod(input)
	return &out, nil
}

type AutoUpgradeOptions string

const (
	AutoUpgradeOptionsDisabled AutoUpgradeOptions = "Disabled"
	AutoUpgradeOptionsEnabled  AutoUpgradeOptions = "Enabled"
)

func PossibleValuesForAutoUpgradeOptions() []string {
	return []string{
		string(AutoUpgradeOptionsDisabled),
		string(AutoUpgradeOptionsEnabled),
	}
}

func (s *AutoUpgradeOptions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutoUpgradeOptions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutoUpgradeOptions(input string) (*AutoUpgradeOptions, error) {
	vals := map[string]AutoUpgradeOptions{
		"disabled": AutoUpgradeOptionsDisabled,
		"enabled":  AutoUpgradeOptionsEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoUpgradeOptions(input)
	return &out, nil
}

type AzureHybridBenefit string

const (
	AzureHybridBenefitFalse         AzureHybridBenefit = "False"
	AzureHybridBenefitNotApplicable AzureHybridBenefit = "NotApplicable"
	AzureHybridBenefitTrue          AzureHybridBenefit = "True"
)

func PossibleValuesForAzureHybridBenefit() []string {
	return []string{
		string(AzureHybridBenefitFalse),
		string(AzureHybridBenefitNotApplicable),
		string(AzureHybridBenefitTrue),
	}
}

func (s *AzureHybridBenefit) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureHybridBenefit(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureHybridBenefit(input string) (*AzureHybridBenefit, error) {
	vals := map[string]AzureHybridBenefit{
		"false":         AzureHybridBenefitFalse,
		"notapplicable": AzureHybridBenefitNotApplicable,
		"true":          AzureHybridBenefitTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureHybridBenefit(input)
	return &out, nil
}

type ConnectedClusterKind string

const (
	ConnectedClusterKindProvisionedCluster ConnectedClusterKind = "ProvisionedCluster"
)

func PossibleValuesForConnectedClusterKind() []string {
	return []string{
		string(ConnectedClusterKindProvisionedCluster),
	}
}

func (s *ConnectedClusterKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectedClusterKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectedClusterKind(input string) (*ConnectedClusterKind, error) {
	vals := map[string]ConnectedClusterKind{
		"provisionedcluster": ConnectedClusterKindProvisionedCluster,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectedClusterKind(input)
	return &out, nil
}

type ConnectivityStatus string

const (
	ConnectivityStatusConnected  ConnectivityStatus = "Connected"
	ConnectivityStatusConnecting ConnectivityStatus = "Connecting"
	ConnectivityStatusExpired    ConnectivityStatus = "Expired"
	ConnectivityStatusOffline    ConnectivityStatus = "Offline"
)

func PossibleValuesForConnectivityStatus() []string {
	return []string{
		string(ConnectivityStatusConnected),
		string(ConnectivityStatusConnecting),
		string(ConnectivityStatusExpired),
		string(ConnectivityStatusOffline),
	}
}

func (s *ConnectivityStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectivityStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectivityStatus(input string) (*ConnectivityStatus, error) {
	vals := map[string]ConnectivityStatus{
		"connected":  ConnectivityStatusConnected,
		"connecting": ConnectivityStatusConnecting,
		"expired":    ConnectivityStatusExpired,
		"offline":    ConnectivityStatusOffline,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectivityStatus(input)
	return &out, nil
}

type PrivateLinkState string

const (
	PrivateLinkStateDisabled PrivateLinkState = "Disabled"
	PrivateLinkStateEnabled  PrivateLinkState = "Enabled"
)

func PossibleValuesForPrivateLinkState() []string {
	return []string{
		string(PrivateLinkStateDisabled),
		string(PrivateLinkStateEnabled),
	}
}

func (s *PrivateLinkState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateLinkState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateLinkState(input string) (*PrivateLinkState, error) {
	vals := map[string]PrivateLinkState{
		"disabled": PrivateLinkStateDisabled,
		"enabled":  PrivateLinkStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateLinkState(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateProvisioning),
		string(ProvisioningStateSucceeded),
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
		"accepted":     ProvisioningStateAccepted,
		"canceled":     ProvisioningStateCanceled,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"provisioning": ProvisioningStateProvisioning,
		"succeeded":    ProvisioningStateSucceeded,
		"updating":     ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
