package databases

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CallerRole string

const (
	CallerRoleAdmin CallerRole = "Admin"
	CallerRoleNone  CallerRole = "None"
)

func PossibleValuesForCallerRole() []string {
	return []string{
		string(CallerRoleAdmin),
		string(CallerRoleNone),
	}
}

func (s *CallerRole) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCallerRole(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCallerRole(input string) (*CallerRole, error) {
	vals := map[string]CallerRole{
		"admin": CallerRoleAdmin,
		"none":  CallerRoleNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CallerRole(input)
	return &out, nil
}

type DatabasePrincipalRole string

const (
	DatabasePrincipalRoleAdmin              DatabasePrincipalRole = "Admin"
	DatabasePrincipalRoleIngestor           DatabasePrincipalRole = "Ingestor"
	DatabasePrincipalRoleMonitor            DatabasePrincipalRole = "Monitor"
	DatabasePrincipalRoleUnrestrictedViewer DatabasePrincipalRole = "UnrestrictedViewer"
	DatabasePrincipalRoleUser               DatabasePrincipalRole = "User"
	DatabasePrincipalRoleViewer             DatabasePrincipalRole = "Viewer"
)

func PossibleValuesForDatabasePrincipalRole() []string {
	return []string{
		string(DatabasePrincipalRoleAdmin),
		string(DatabasePrincipalRoleIngestor),
		string(DatabasePrincipalRoleMonitor),
		string(DatabasePrincipalRoleUnrestrictedViewer),
		string(DatabasePrincipalRoleUser),
		string(DatabasePrincipalRoleViewer),
	}
}

func (s *DatabasePrincipalRole) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatabasePrincipalRole(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatabasePrincipalRole(input string) (*DatabasePrincipalRole, error) {
	vals := map[string]DatabasePrincipalRole{
		"admin":              DatabasePrincipalRoleAdmin,
		"ingestor":           DatabasePrincipalRoleIngestor,
		"monitor":            DatabasePrincipalRoleMonitor,
		"unrestrictedviewer": DatabasePrincipalRoleUnrestrictedViewer,
		"user":               DatabasePrincipalRoleUser,
		"viewer":             DatabasePrincipalRoleViewer,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatabasePrincipalRole(input)
	return &out, nil
}

type DatabasePrincipalType string

const (
	DatabasePrincipalTypeApp   DatabasePrincipalType = "App"
	DatabasePrincipalTypeGroup DatabasePrincipalType = "Group"
	DatabasePrincipalTypeUser  DatabasePrincipalType = "User"
)

func PossibleValuesForDatabasePrincipalType() []string {
	return []string{
		string(DatabasePrincipalTypeApp),
		string(DatabasePrincipalTypeGroup),
		string(DatabasePrincipalTypeUser),
	}
}

func (s *DatabasePrincipalType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatabasePrincipalType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatabasePrincipalType(input string) (*DatabasePrincipalType, error) {
	vals := map[string]DatabasePrincipalType{
		"app":   DatabasePrincipalTypeApp,
		"group": DatabasePrincipalTypeGroup,
		"user":  DatabasePrincipalTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatabasePrincipalType(input)
	return &out, nil
}

type DatabaseShareOrigin string

const (
	DatabaseShareOriginDataShare DatabaseShareOrigin = "DataShare"
	DatabaseShareOriginDirect    DatabaseShareOrigin = "Direct"
	DatabaseShareOriginOther     DatabaseShareOrigin = "Other"
)

func PossibleValuesForDatabaseShareOrigin() []string {
	return []string{
		string(DatabaseShareOriginDataShare),
		string(DatabaseShareOriginDirect),
		string(DatabaseShareOriginOther),
	}
}

func (s *DatabaseShareOrigin) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatabaseShareOrigin(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatabaseShareOrigin(input string) (*DatabaseShareOrigin, error) {
	vals := map[string]DatabaseShareOrigin{
		"datashare": DatabaseShareOriginDataShare,
		"direct":    DatabaseShareOriginDirect,
		"other":     DatabaseShareOriginOther,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatabaseShareOrigin(input)
	return &out, nil
}

type Kind string

const (
	KindReadOnlyFollowing Kind = "ReadOnlyFollowing"
	KindReadWrite         Kind = "ReadWrite"
)

func PossibleValuesForKind() []string {
	return []string{
		string(KindReadOnlyFollowing),
		string(KindReadWrite),
	}
}

func (s *Kind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKind(input string) (*Kind, error) {
	vals := map[string]Kind{
		"readonlyfollowing": KindReadOnlyFollowing,
		"readwrite":         KindReadWrite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Kind(input)
	return &out, nil
}

type PrincipalsModificationKind string

const (
	PrincipalsModificationKindNone    PrincipalsModificationKind = "None"
	PrincipalsModificationKindReplace PrincipalsModificationKind = "Replace"
	PrincipalsModificationKindUnion   PrincipalsModificationKind = "Union"
)

func PossibleValuesForPrincipalsModificationKind() []string {
	return []string{
		string(PrincipalsModificationKindNone),
		string(PrincipalsModificationKindReplace),
		string(PrincipalsModificationKindUnion),
	}
}

func (s *PrincipalsModificationKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrincipalsModificationKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrincipalsModificationKind(input string) (*PrincipalsModificationKind, error) {
	vals := map[string]PrincipalsModificationKind{
		"none":    PrincipalsModificationKindNone,
		"replace": PrincipalsModificationKindReplace,
		"union":   PrincipalsModificationKindUnion,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrincipalsModificationKind(input)
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

func (s *Reason) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReason(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type Type string

const (
	TypeMicrosoftPointKustoClustersAttachedDatabaseConfigurations Type = "Microsoft.Kusto/clusters/attachedDatabaseConfigurations"
	TypeMicrosoftPointKustoClustersDatabases                      Type = "Microsoft.Kusto/clusters/databases"
)

func PossibleValuesForType() []string {
	return []string{
		string(TypeMicrosoftPointKustoClustersAttachedDatabaseConfigurations),
		string(TypeMicrosoftPointKustoClustersDatabases),
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
		"microsoft.kusto/clusters/attacheddatabaseconfigurations": TypeMicrosoftPointKustoClustersAttachedDatabaseConfigurations,
		"microsoft.kusto/clusters/databases":                      TypeMicrosoftPointKustoClustersDatabases,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Type(input)
	return &out, nil
}
