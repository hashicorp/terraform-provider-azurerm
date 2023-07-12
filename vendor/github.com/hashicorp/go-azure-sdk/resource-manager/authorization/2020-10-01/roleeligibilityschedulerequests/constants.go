package roleeligibilityschedulerequests

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrincipalType string

const (
	PrincipalTypeDevice           PrincipalType = "Device"
	PrincipalTypeForeignGroup     PrincipalType = "ForeignGroup"
	PrincipalTypeGroup            PrincipalType = "Group"
	PrincipalTypeServicePrincipal PrincipalType = "ServicePrincipal"
	PrincipalTypeUser             PrincipalType = "User"
)

func PossibleValuesForPrincipalType() []string {
	return []string{
		string(PrincipalTypeDevice),
		string(PrincipalTypeForeignGroup),
		string(PrincipalTypeGroup),
		string(PrincipalTypeServicePrincipal),
		string(PrincipalTypeUser),
	}
}

func (s *PrincipalType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrincipalType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrincipalType(input string) (*PrincipalType, error) {
	vals := map[string]PrincipalType{
		"device":           PrincipalTypeDevice,
		"foreigngroup":     PrincipalTypeForeignGroup,
		"group":            PrincipalTypeGroup,
		"serviceprincipal": PrincipalTypeServicePrincipal,
		"user":             PrincipalTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrincipalType(input)
	return &out, nil
}

type RequestType string

const (
	RequestTypeAdminAssign    RequestType = "AdminAssign"
	RequestTypeAdminExtend    RequestType = "AdminExtend"
	RequestTypeAdminRemove    RequestType = "AdminRemove"
	RequestTypeAdminRenew     RequestType = "AdminRenew"
	RequestTypeAdminUpdate    RequestType = "AdminUpdate"
	RequestTypeSelfActivate   RequestType = "SelfActivate"
	RequestTypeSelfDeactivate RequestType = "SelfDeactivate"
	RequestTypeSelfExtend     RequestType = "SelfExtend"
	RequestTypeSelfRenew      RequestType = "SelfRenew"
)

func PossibleValuesForRequestType() []string {
	return []string{
		string(RequestTypeAdminAssign),
		string(RequestTypeAdminExtend),
		string(RequestTypeAdminRemove),
		string(RequestTypeAdminRenew),
		string(RequestTypeAdminUpdate),
		string(RequestTypeSelfActivate),
		string(RequestTypeSelfDeactivate),
		string(RequestTypeSelfExtend),
		string(RequestTypeSelfRenew),
	}
}

func (s *RequestType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRequestType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRequestType(input string) (*RequestType, error) {
	vals := map[string]RequestType{
		"adminassign":    RequestTypeAdminAssign,
		"adminextend":    RequestTypeAdminExtend,
		"adminremove":    RequestTypeAdminRemove,
		"adminrenew":     RequestTypeAdminRenew,
		"adminupdate":    RequestTypeAdminUpdate,
		"selfactivate":   RequestTypeSelfActivate,
		"selfdeactivate": RequestTypeSelfDeactivate,
		"selfextend":     RequestTypeSelfExtend,
		"selfrenew":      RequestTypeSelfRenew,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RequestType(input)
	return &out, nil
}

type Status string

const (
	StatusAccepted                    Status = "Accepted"
	StatusAdminApproved               Status = "AdminApproved"
	StatusAdminDenied                 Status = "AdminDenied"
	StatusCanceled                    Status = "Canceled"
	StatusDenied                      Status = "Denied"
	StatusFailed                      Status = "Failed"
	StatusFailedAsResourceIsLocked    Status = "FailedAsResourceIsLocked"
	StatusGranted                     Status = "Granted"
	StatusInvalid                     Status = "Invalid"
	StatusPendingAdminDecision        Status = "PendingAdminDecision"
	StatusPendingApproval             Status = "PendingApproval"
	StatusPendingApprovalProvisioning Status = "PendingApprovalProvisioning"
	StatusPendingEvaluation           Status = "PendingEvaluation"
	StatusPendingExternalProvisioning Status = "PendingExternalProvisioning"
	StatusPendingProvisioning         Status = "PendingProvisioning"
	StatusPendingRevocation           Status = "PendingRevocation"
	StatusPendingScheduleCreation     Status = "PendingScheduleCreation"
	StatusProvisioned                 Status = "Provisioned"
	StatusProvisioningStarted         Status = "ProvisioningStarted"
	StatusRevoked                     Status = "Revoked"
	StatusScheduleCreated             Status = "ScheduleCreated"
	StatusTimedOut                    Status = "TimedOut"
)

func PossibleValuesForStatus() []string {
	return []string{
		string(StatusAccepted),
		string(StatusAdminApproved),
		string(StatusAdminDenied),
		string(StatusCanceled),
		string(StatusDenied),
		string(StatusFailed),
		string(StatusFailedAsResourceIsLocked),
		string(StatusGranted),
		string(StatusInvalid),
		string(StatusPendingAdminDecision),
		string(StatusPendingApproval),
		string(StatusPendingApprovalProvisioning),
		string(StatusPendingEvaluation),
		string(StatusPendingExternalProvisioning),
		string(StatusPendingProvisioning),
		string(StatusPendingRevocation),
		string(StatusPendingScheduleCreation),
		string(StatusProvisioned),
		string(StatusProvisioningStarted),
		string(StatusRevoked),
		string(StatusScheduleCreated),
		string(StatusTimedOut),
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
		"accepted":                    StatusAccepted,
		"adminapproved":               StatusAdminApproved,
		"admindenied":                 StatusAdminDenied,
		"canceled":                    StatusCanceled,
		"denied":                      StatusDenied,
		"failed":                      StatusFailed,
		"failedasresourceislocked":    StatusFailedAsResourceIsLocked,
		"granted":                     StatusGranted,
		"invalid":                     StatusInvalid,
		"pendingadmindecision":        StatusPendingAdminDecision,
		"pendingapproval":             StatusPendingApproval,
		"pendingapprovalprovisioning": StatusPendingApprovalProvisioning,
		"pendingevaluation":           StatusPendingEvaluation,
		"pendingexternalprovisioning": StatusPendingExternalProvisioning,
		"pendingprovisioning":         StatusPendingProvisioning,
		"pendingrevocation":           StatusPendingRevocation,
		"pendingschedulecreation":     StatusPendingScheduleCreation,
		"provisioned":                 StatusProvisioned,
		"provisioningstarted":         StatusProvisioningStarted,
		"revoked":                     StatusRevoked,
		"schedulecreated":             StatusScheduleCreated,
		"timedout":                    StatusTimedOut,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Status(input)
	return &out, nil
}

type Type string

const (
	TypeAfterDateTime Type = "AfterDateTime"
	TypeAfterDuration Type = "AfterDuration"
	TypeNoExpiration  Type = "NoExpiration"
)

func PossibleValuesForType() []string {
	return []string{
		string(TypeAfterDateTime),
		string(TypeAfterDuration),
		string(TypeNoExpiration),
	}
}

func (s *Type) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseType(input string) (*Type, error) {
	vals := map[string]Type{
		"afterdatetime": TypeAfterDateTime,
		"afterduration": TypeAfterDuration,
		"noexpiration":  TypeNoExpiration,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Type(input)
	return &out, nil
}
