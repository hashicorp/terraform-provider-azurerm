package applications

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationArtifactName string

const (
	ApplicationArtifactNameAuthorizations       ApplicationArtifactName = "Authorizations"
	ApplicationArtifactNameCustomRoleDefinition ApplicationArtifactName = "CustomRoleDefinition"
	ApplicationArtifactNameNotSpecified         ApplicationArtifactName = "NotSpecified"
	ApplicationArtifactNameViewDefinition       ApplicationArtifactName = "ViewDefinition"
)

func PossibleValuesForApplicationArtifactName() []string {
	return []string{
		string(ApplicationArtifactNameAuthorizations),
		string(ApplicationArtifactNameCustomRoleDefinition),
		string(ApplicationArtifactNameNotSpecified),
		string(ApplicationArtifactNameViewDefinition),
	}
}

func (s *ApplicationArtifactName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationArtifactName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationArtifactName(input string) (*ApplicationArtifactName, error) {
	vals := map[string]ApplicationArtifactName{
		"authorizations":       ApplicationArtifactNameAuthorizations,
		"customroledefinition": ApplicationArtifactNameCustomRoleDefinition,
		"notspecified":         ApplicationArtifactNameNotSpecified,
		"viewdefinition":       ApplicationArtifactNameViewDefinition,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationArtifactName(input)
	return &out, nil
}

type ApplicationArtifactType string

const (
	ApplicationArtifactTypeCustom       ApplicationArtifactType = "Custom"
	ApplicationArtifactTypeNotSpecified ApplicationArtifactType = "NotSpecified"
	ApplicationArtifactTypeTemplate     ApplicationArtifactType = "Template"
)

func PossibleValuesForApplicationArtifactType() []string {
	return []string{
		string(ApplicationArtifactTypeCustom),
		string(ApplicationArtifactTypeNotSpecified),
		string(ApplicationArtifactTypeTemplate),
	}
}

func (s *ApplicationArtifactType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationArtifactType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationArtifactType(input string) (*ApplicationArtifactType, error) {
	vals := map[string]ApplicationArtifactType{
		"custom":       ApplicationArtifactTypeCustom,
		"notspecified": ApplicationArtifactTypeNotSpecified,
		"template":     ApplicationArtifactTypeTemplate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationArtifactType(input)
	return &out, nil
}

type ApplicationManagementMode string

const (
	ApplicationManagementModeManaged      ApplicationManagementMode = "Managed"
	ApplicationManagementModeNotSpecified ApplicationManagementMode = "NotSpecified"
	ApplicationManagementModeUnmanaged    ApplicationManagementMode = "Unmanaged"
)

func PossibleValuesForApplicationManagementMode() []string {
	return []string{
		string(ApplicationManagementModeManaged),
		string(ApplicationManagementModeNotSpecified),
		string(ApplicationManagementModeUnmanaged),
	}
}

func (s *ApplicationManagementMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationManagementMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationManagementMode(input string) (*ApplicationManagementMode, error) {
	vals := map[string]ApplicationManagementMode{
		"managed":      ApplicationManagementModeManaged,
		"notspecified": ApplicationManagementModeNotSpecified,
		"unmanaged":    ApplicationManagementModeUnmanaged,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationManagementMode(input)
	return &out, nil
}

type JitApprovalMode string

const (
	JitApprovalModeAutoApprove   JitApprovalMode = "AutoApprove"
	JitApprovalModeManualApprove JitApprovalMode = "ManualApprove"
	JitApprovalModeNotSpecified  JitApprovalMode = "NotSpecified"
)

func PossibleValuesForJitApprovalMode() []string {
	return []string{
		string(JitApprovalModeAutoApprove),
		string(JitApprovalModeManualApprove),
		string(JitApprovalModeNotSpecified),
	}
}

func (s *JitApprovalMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJitApprovalMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJitApprovalMode(input string) (*JitApprovalMode, error) {
	vals := map[string]JitApprovalMode{
		"autoapprove":   JitApprovalModeAutoApprove,
		"manualapprove": JitApprovalModeManualApprove,
		"notspecified":  JitApprovalModeNotSpecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JitApprovalMode(input)
	return &out, nil
}

type JitApproverType string

const (
	JitApproverTypeGroup JitApproverType = "group"
	JitApproverTypeUser  JitApproverType = "user"
)

func PossibleValuesForJitApproverType() []string {
	return []string{
		string(JitApproverTypeGroup),
		string(JitApproverTypeUser),
	}
}

func (s *JitApproverType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJitApproverType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJitApproverType(input string) (*JitApproverType, error) {
	vals := map[string]JitApproverType{
		"group": JitApproverTypeGroup,
		"user":  JitApproverTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JitApproverType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateDeleted      ProvisioningState = "Deleted"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateNotSpecified ProvisioningState = "NotSpecified"
	ProvisioningStateRunning      ProvisioningState = "Running"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleted),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateNotSpecified),
		string(ProvisioningStateRunning),
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
		"deleted":      ProvisioningStateDeleted,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"notspecified": ProvisioningStateNotSpecified,
		"running":      ProvisioningStateRunning,
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

type ResourceIdentityType string

const (
	ResourceIdentityTypeNone                       ResourceIdentityType = "None"
	ResourceIdentityTypeSystemAssigned             ResourceIdentityType = "SystemAssigned"
	ResourceIdentityTypeSystemAssignedUserAssigned ResourceIdentityType = "SystemAssigned, UserAssigned"
	ResourceIdentityTypeUserAssigned               ResourceIdentityType = "UserAssigned"
)

func PossibleValuesForResourceIdentityType() []string {
	return []string{
		string(ResourceIdentityTypeNone),
		string(ResourceIdentityTypeSystemAssigned),
		string(ResourceIdentityTypeSystemAssignedUserAssigned),
		string(ResourceIdentityTypeUserAssigned),
	}
}

func (s *ResourceIdentityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceIdentityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceIdentityType(input string) (*ResourceIdentityType, error) {
	vals := map[string]ResourceIdentityType{
		"none":                         ResourceIdentityTypeNone,
		"systemassigned":               ResourceIdentityTypeSystemAssigned,
		"systemassigned, userassigned": ResourceIdentityTypeSystemAssignedUserAssigned,
		"userassigned":                 ResourceIdentityTypeUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceIdentityType(input)
	return &out, nil
}

type Status string

const (
	StatusElevate      Status = "Elevate"
	StatusNotSpecified Status = "NotSpecified"
	StatusRemove       Status = "Remove"
)

func PossibleValuesForStatus() []string {
	return []string{
		string(StatusElevate),
		string(StatusNotSpecified),
		string(StatusRemove),
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
		"elevate":      StatusElevate,
		"notspecified": StatusNotSpecified,
		"remove":       StatusRemove,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Status(input)
	return &out, nil
}

type Substatus string

const (
	SubstatusApproved     Substatus = "Approved"
	SubstatusDenied       Substatus = "Denied"
	SubstatusExpired      Substatus = "Expired"
	SubstatusFailed       Substatus = "Failed"
	SubstatusNotSpecified Substatus = "NotSpecified"
	SubstatusTimeout      Substatus = "Timeout"
)

func PossibleValuesForSubstatus() []string {
	return []string{
		string(SubstatusApproved),
		string(SubstatusDenied),
		string(SubstatusExpired),
		string(SubstatusFailed),
		string(SubstatusNotSpecified),
		string(SubstatusTimeout),
	}
}

func (s *Substatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSubstatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSubstatus(input string) (*Substatus, error) {
	vals := map[string]Substatus{
		"approved":     SubstatusApproved,
		"denied":       SubstatusDenied,
		"expired":      SubstatusExpired,
		"failed":       SubstatusFailed,
		"notspecified": SubstatusNotSpecified,
		"timeout":      SubstatusTimeout,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Substatus(input)
	return &out, nil
}
