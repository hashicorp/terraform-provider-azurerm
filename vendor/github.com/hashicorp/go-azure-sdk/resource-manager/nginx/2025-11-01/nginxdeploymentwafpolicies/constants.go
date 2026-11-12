package nginxdeploymentwafpolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxDeploymentWafPolicyApplyingStatusCode string

const (
	NginxDeploymentWafPolicyApplyingStatusCodeApplying   NginxDeploymentWafPolicyApplyingStatusCode = "Applying"
	NginxDeploymentWafPolicyApplyingStatusCodeFailed     NginxDeploymentWafPolicyApplyingStatusCode = "Failed"
	NginxDeploymentWafPolicyApplyingStatusCodeNotApplied NginxDeploymentWafPolicyApplyingStatusCode = "NotApplied"
	NginxDeploymentWafPolicyApplyingStatusCodeRemoving   NginxDeploymentWafPolicyApplyingStatusCode = "Removing"
	NginxDeploymentWafPolicyApplyingStatusCodeSucceeded  NginxDeploymentWafPolicyApplyingStatusCode = "Succeeded"
)

func PossibleValuesForNginxDeploymentWafPolicyApplyingStatusCode() []string {
	return []string{
		string(NginxDeploymentWafPolicyApplyingStatusCodeApplying),
		string(NginxDeploymentWafPolicyApplyingStatusCodeFailed),
		string(NginxDeploymentWafPolicyApplyingStatusCodeNotApplied),
		string(NginxDeploymentWafPolicyApplyingStatusCodeRemoving),
		string(NginxDeploymentWafPolicyApplyingStatusCodeSucceeded),
	}
}

func (s *NginxDeploymentWafPolicyApplyingStatusCode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNginxDeploymentWafPolicyApplyingStatusCode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNginxDeploymentWafPolicyApplyingStatusCode(input string) (*NginxDeploymentWafPolicyApplyingStatusCode, error) {
	vals := map[string]NginxDeploymentWafPolicyApplyingStatusCode{
		"applying":   NginxDeploymentWafPolicyApplyingStatusCodeApplying,
		"failed":     NginxDeploymentWafPolicyApplyingStatusCodeFailed,
		"notapplied": NginxDeploymentWafPolicyApplyingStatusCodeNotApplied,
		"removing":   NginxDeploymentWafPolicyApplyingStatusCodeRemoving,
		"succeeded":  NginxDeploymentWafPolicyApplyingStatusCodeSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NginxDeploymentWafPolicyApplyingStatusCode(input)
	return &out, nil
}

type NginxDeploymentWafPolicyCompilingStatusCode string

const (
	NginxDeploymentWafPolicyCompilingStatusCodeFailed     NginxDeploymentWafPolicyCompilingStatusCode = "Failed"
	NginxDeploymentWafPolicyCompilingStatusCodeInProgress NginxDeploymentWafPolicyCompilingStatusCode = "InProgress"
	NginxDeploymentWafPolicyCompilingStatusCodeNotStarted NginxDeploymentWafPolicyCompilingStatusCode = "NotStarted"
	NginxDeploymentWafPolicyCompilingStatusCodeSucceeded  NginxDeploymentWafPolicyCompilingStatusCode = "Succeeded"
)

func PossibleValuesForNginxDeploymentWafPolicyCompilingStatusCode() []string {
	return []string{
		string(NginxDeploymentWafPolicyCompilingStatusCodeFailed),
		string(NginxDeploymentWafPolicyCompilingStatusCodeInProgress),
		string(NginxDeploymentWafPolicyCompilingStatusCodeNotStarted),
		string(NginxDeploymentWafPolicyCompilingStatusCodeSucceeded),
	}
}

func (s *NginxDeploymentWafPolicyCompilingStatusCode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNginxDeploymentWafPolicyCompilingStatusCode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNginxDeploymentWafPolicyCompilingStatusCode(input string) (*NginxDeploymentWafPolicyCompilingStatusCode, error) {
	vals := map[string]NginxDeploymentWafPolicyCompilingStatusCode{
		"failed":     NginxDeploymentWafPolicyCompilingStatusCodeFailed,
		"inprogress": NginxDeploymentWafPolicyCompilingStatusCodeInProgress,
		"notstarted": NginxDeploymentWafPolicyCompilingStatusCodeNotStarted,
		"succeeded":  NginxDeploymentWafPolicyCompilingStatusCodeSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NginxDeploymentWafPolicyCompilingStatusCode(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateCreating     ProvisioningState = "Creating"
	ProvisioningStateDeleted      ProvisioningState = "Deleted"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateNotSpecified ProvisioningState = "NotSpecified"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleted),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateNotSpecified),
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
		"creating":     ProvisioningStateCreating,
		"deleted":      ProvisioningStateDeleted,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"notspecified": ProvisioningStateNotSpecified,
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
