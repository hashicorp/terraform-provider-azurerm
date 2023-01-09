package clusterprincipalassignments

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterPrincipalRole string

const (
	ClusterPrincipalRoleAllDatabasesAdmin  ClusterPrincipalRole = "AllDatabasesAdmin"
	ClusterPrincipalRoleAllDatabasesViewer ClusterPrincipalRole = "AllDatabasesViewer"
)

func PossibleValuesForClusterPrincipalRole() []string {
	return []string{
		string(ClusterPrincipalRoleAllDatabasesAdmin),
		string(ClusterPrincipalRoleAllDatabasesViewer),
	}
}

func parseClusterPrincipalRole(input string) (*ClusterPrincipalRole, error) {
	vals := map[string]ClusterPrincipalRole{
		"alldatabasesadmin":  ClusterPrincipalRoleAllDatabasesAdmin,
		"alldatabasesviewer": ClusterPrincipalRoleAllDatabasesViewer,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterPrincipalRole(input)
	return &out, nil
}

type PrincipalAssignmentType string

const (
	PrincipalAssignmentTypeMicrosoftPointKustoClustersPrincipalAssignments PrincipalAssignmentType = "Microsoft.Kusto/clusters/principalAssignments"
)

func PossibleValuesForPrincipalAssignmentType() []string {
	return []string{
		string(PrincipalAssignmentTypeMicrosoftPointKustoClustersPrincipalAssignments),
	}
}

func parsePrincipalAssignmentType(input string) (*PrincipalAssignmentType, error) {
	vals := map[string]PrincipalAssignmentType{
		"microsoft.kusto/clusters/principalassignments": PrincipalAssignmentTypeMicrosoftPointKustoClustersPrincipalAssignments,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrincipalAssignmentType(input)
	return &out, nil
}

type PrincipalType string

const (
	PrincipalTypeApp   PrincipalType = "App"
	PrincipalTypeGroup PrincipalType = "Group"
	PrincipalTypeUser  PrincipalType = "User"
)

func PossibleValuesForPrincipalType() []string {
	return []string{
		string(PrincipalTypeApp),
		string(PrincipalTypeGroup),
		string(PrincipalTypeUser),
	}
}

func parsePrincipalType(input string) (*PrincipalType, error) {
	vals := map[string]PrincipalType{
		"app":   PrincipalTypeApp,
		"group": PrincipalTypeGroup,
		"user":  PrincipalTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrincipalType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateMoving    ProvisioningState = "Moving"
	ProvisioningStateRunning   ProvisioningState = "Running"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateMoving),
		string(ProvisioningStateRunning),
		string(ProvisioningStateSucceeded),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"moving":    ProvisioningStateMoving,
		"running":   ProvisioningStateRunning,
		"succeeded": ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type Reason string

const (
	ReasonAlreadyExists Reason = "AlreadyExists"
	ReasonInvalid       Reason = "Invalid"
)

func PossibleValuesForReason() []string {
	return []string{
		string(ReasonAlreadyExists),
		string(ReasonInvalid),
	}
}

func parseReason(input string) (*Reason, error) {
	vals := map[string]Reason{
		"alreadyexists": ReasonAlreadyExists,
		"invalid":       ReasonInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Reason(input)
	return &out, nil
}
