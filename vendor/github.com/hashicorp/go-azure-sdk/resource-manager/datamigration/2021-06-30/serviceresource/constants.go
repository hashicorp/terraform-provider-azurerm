package serviceresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthenticationType string

const (
	AuthenticationTypeActiveDirectoryIntegrated AuthenticationType = "ActiveDirectoryIntegrated"
	AuthenticationTypeActiveDirectoryPassword   AuthenticationType = "ActiveDirectoryPassword"
	AuthenticationTypeNone                      AuthenticationType = "None"
	AuthenticationTypeSqlAuthentication         AuthenticationType = "SqlAuthentication"
	AuthenticationTypeWindowsAuthentication     AuthenticationType = "WindowsAuthentication"
)

func PossibleValuesForAuthenticationType() []string {
	return []string{
		string(AuthenticationTypeActiveDirectoryIntegrated),
		string(AuthenticationTypeActiveDirectoryPassword),
		string(AuthenticationTypeNone),
		string(AuthenticationTypeSqlAuthentication),
		string(AuthenticationTypeWindowsAuthentication),
	}
}

func (s *AuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAuthenticationType(input string) (*AuthenticationType, error) {
	vals := map[string]AuthenticationType{
		"activedirectoryintegrated": AuthenticationTypeActiveDirectoryIntegrated,
		"activedirectorypassword":   AuthenticationTypeActiveDirectoryPassword,
		"none":                      AuthenticationTypeNone,
		"sqlauthentication":         AuthenticationTypeSqlAuthentication,
		"windowsauthentication":     AuthenticationTypeWindowsAuthentication,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthenticationType(input)
	return &out, nil
}

type BackupFileStatus string

const (
	BackupFileStatusArrived   BackupFileStatus = "Arrived"
	BackupFileStatusCancelled BackupFileStatus = "Cancelled"
	BackupFileStatusQueued    BackupFileStatus = "Queued"
	BackupFileStatusRestored  BackupFileStatus = "Restored"
	BackupFileStatusRestoring BackupFileStatus = "Restoring"
	BackupFileStatusUploaded  BackupFileStatus = "Uploaded"
	BackupFileStatusUploading BackupFileStatus = "Uploading"
)

func PossibleValuesForBackupFileStatus() []string {
	return []string{
		string(BackupFileStatusArrived),
		string(BackupFileStatusCancelled),
		string(BackupFileStatusQueued),
		string(BackupFileStatusRestored),
		string(BackupFileStatusRestoring),
		string(BackupFileStatusUploaded),
		string(BackupFileStatusUploading),
	}
}

func (s *BackupFileStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBackupFileStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBackupFileStatus(input string) (*BackupFileStatus, error) {
	vals := map[string]BackupFileStatus{
		"arrived":   BackupFileStatusArrived,
		"cancelled": BackupFileStatusCancelled,
		"queued":    BackupFileStatusQueued,
		"restored":  BackupFileStatusRestored,
		"restoring": BackupFileStatusRestoring,
		"uploaded":  BackupFileStatusUploaded,
		"uploading": BackupFileStatusUploading,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackupFileStatus(input)
	return &out, nil
}

type BackupMode string

const (
	BackupModeCreateBackup   BackupMode = "CreateBackup"
	BackupModeExistingBackup BackupMode = "ExistingBackup"
)

func PossibleValuesForBackupMode() []string {
	return []string{
		string(BackupModeCreateBackup),
		string(BackupModeExistingBackup),
	}
}

func (s *BackupMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBackupMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBackupMode(input string) (*BackupMode, error) {
	vals := map[string]BackupMode{
		"createbackup":   BackupModeCreateBackup,
		"existingbackup": BackupModeExistingBackup,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackupMode(input)
	return &out, nil
}

type BackupType string

const (
	BackupTypeDatabase             BackupType = "Database"
	BackupTypeDifferentialDatabase BackupType = "DifferentialDatabase"
	BackupTypeDifferentialFile     BackupType = "DifferentialFile"
	BackupTypeDifferentialPartial  BackupType = "DifferentialPartial"
	BackupTypeFile                 BackupType = "File"
	BackupTypePartial              BackupType = "Partial"
	BackupTypeTransactionLog       BackupType = "TransactionLog"
)

func PossibleValuesForBackupType() []string {
	return []string{
		string(BackupTypeDatabase),
		string(BackupTypeDifferentialDatabase),
		string(BackupTypeDifferentialFile),
		string(BackupTypeDifferentialPartial),
		string(BackupTypeFile),
		string(BackupTypePartial),
		string(BackupTypeTransactionLog),
	}
}

func (s *BackupType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBackupType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBackupType(input string) (*BackupType, error) {
	vals := map[string]BackupType{
		"database":             BackupTypeDatabase,
		"differentialdatabase": BackupTypeDifferentialDatabase,
		"differentialfile":     BackupTypeDifferentialFile,
		"differentialpartial":  BackupTypeDifferentialPartial,
		"file":                 BackupTypeFile,
		"partial":              BackupTypePartial,
		"transactionlog":       BackupTypeTransactionLog,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackupType(input)
	return &out, nil
}

type CommandState string

const (
	CommandStateAccepted  CommandState = "Accepted"
	CommandStateFailed    CommandState = "Failed"
	CommandStateRunning   CommandState = "Running"
	CommandStateSucceeded CommandState = "Succeeded"
	CommandStateUnknown   CommandState = "Unknown"
)

func PossibleValuesForCommandState() []string {
	return []string{
		string(CommandStateAccepted),
		string(CommandStateFailed),
		string(CommandStateRunning),
		string(CommandStateSucceeded),
		string(CommandStateUnknown),
	}
}

func (s *CommandState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCommandState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCommandState(input string) (*CommandState, error) {
	vals := map[string]CommandState{
		"accepted":  CommandStateAccepted,
		"failed":    CommandStateFailed,
		"running":   CommandStateRunning,
		"succeeded": CommandStateSucceeded,
		"unknown":   CommandStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CommandState(input)
	return &out, nil
}

type DatabaseCompatLevel string

const (
	DatabaseCompatLevelCompatLevelEightZero    DatabaseCompatLevel = "CompatLevel80"
	DatabaseCompatLevelCompatLevelNineZero     DatabaseCompatLevel = "CompatLevel90"
	DatabaseCompatLevelCompatLevelOneFourZero  DatabaseCompatLevel = "CompatLevel140"
	DatabaseCompatLevelCompatLevelOneHundred   DatabaseCompatLevel = "CompatLevel100"
	DatabaseCompatLevelCompatLevelOneOneZero   DatabaseCompatLevel = "CompatLevel110"
	DatabaseCompatLevelCompatLevelOneThreeZero DatabaseCompatLevel = "CompatLevel130"
	DatabaseCompatLevelCompatLevelOneTwoZero   DatabaseCompatLevel = "CompatLevel120"
)

func PossibleValuesForDatabaseCompatLevel() []string {
	return []string{
		string(DatabaseCompatLevelCompatLevelEightZero),
		string(DatabaseCompatLevelCompatLevelNineZero),
		string(DatabaseCompatLevelCompatLevelOneFourZero),
		string(DatabaseCompatLevelCompatLevelOneHundred),
		string(DatabaseCompatLevelCompatLevelOneOneZero),
		string(DatabaseCompatLevelCompatLevelOneThreeZero),
		string(DatabaseCompatLevelCompatLevelOneTwoZero),
	}
}

func (s *DatabaseCompatLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatabaseCompatLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatabaseCompatLevel(input string) (*DatabaseCompatLevel, error) {
	vals := map[string]DatabaseCompatLevel{
		"compatlevel80":  DatabaseCompatLevelCompatLevelEightZero,
		"compatlevel90":  DatabaseCompatLevelCompatLevelNineZero,
		"compatlevel140": DatabaseCompatLevelCompatLevelOneFourZero,
		"compatlevel100": DatabaseCompatLevelCompatLevelOneHundred,
		"compatlevel110": DatabaseCompatLevelCompatLevelOneOneZero,
		"compatlevel130": DatabaseCompatLevelCompatLevelOneThreeZero,
		"compatlevel120": DatabaseCompatLevelCompatLevelOneTwoZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatabaseCompatLevel(input)
	return &out, nil
}

type DatabaseFileType string

const (
	DatabaseFileTypeFilestream   DatabaseFileType = "Filestream"
	DatabaseFileTypeFulltext     DatabaseFileType = "Fulltext"
	DatabaseFileTypeLog          DatabaseFileType = "Log"
	DatabaseFileTypeNotSupported DatabaseFileType = "NotSupported"
	DatabaseFileTypeRows         DatabaseFileType = "Rows"
)

func PossibleValuesForDatabaseFileType() []string {
	return []string{
		string(DatabaseFileTypeFilestream),
		string(DatabaseFileTypeFulltext),
		string(DatabaseFileTypeLog),
		string(DatabaseFileTypeNotSupported),
		string(DatabaseFileTypeRows),
	}
}

func (s *DatabaseFileType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatabaseFileType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatabaseFileType(input string) (*DatabaseFileType, error) {
	vals := map[string]DatabaseFileType{
		"filestream":   DatabaseFileTypeFilestream,
		"fulltext":     DatabaseFileTypeFulltext,
		"log":          DatabaseFileTypeLog,
		"notsupported": DatabaseFileTypeNotSupported,
		"rows":         DatabaseFileTypeRows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatabaseFileType(input)
	return &out, nil
}

type DatabaseMigrationStage string

const (
	DatabaseMigrationStageBackup     DatabaseMigrationStage = "Backup"
	DatabaseMigrationStageCompleted  DatabaseMigrationStage = "Completed"
	DatabaseMigrationStageFileCopy   DatabaseMigrationStage = "FileCopy"
	DatabaseMigrationStageInitialize DatabaseMigrationStage = "Initialize"
	DatabaseMigrationStageNone       DatabaseMigrationStage = "None"
	DatabaseMigrationStageRestore    DatabaseMigrationStage = "Restore"
)

func PossibleValuesForDatabaseMigrationStage() []string {
	return []string{
		string(DatabaseMigrationStageBackup),
		string(DatabaseMigrationStageCompleted),
		string(DatabaseMigrationStageFileCopy),
		string(DatabaseMigrationStageInitialize),
		string(DatabaseMigrationStageNone),
		string(DatabaseMigrationStageRestore),
	}
}

func (s *DatabaseMigrationStage) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatabaseMigrationStage(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatabaseMigrationStage(input string) (*DatabaseMigrationStage, error) {
	vals := map[string]DatabaseMigrationStage{
		"backup":     DatabaseMigrationStageBackup,
		"completed":  DatabaseMigrationStageCompleted,
		"filecopy":   DatabaseMigrationStageFileCopy,
		"initialize": DatabaseMigrationStageInitialize,
		"none":       DatabaseMigrationStageNone,
		"restore":    DatabaseMigrationStageRestore,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatabaseMigrationStage(input)
	return &out, nil
}

type DatabaseMigrationState string

const (
	DatabaseMigrationStateCANCELLED             DatabaseMigrationState = "CANCELLED"
	DatabaseMigrationStateCOMPLETED             DatabaseMigrationState = "COMPLETED"
	DatabaseMigrationStateCUTOVERSTART          DatabaseMigrationState = "CUTOVER_START"
	DatabaseMigrationStateFAILED                DatabaseMigrationState = "FAILED"
	DatabaseMigrationStateFULLBACKUPUPLOADSTART DatabaseMigrationState = "FULL_BACKUP_UPLOAD_START"
	DatabaseMigrationStateINITIAL               DatabaseMigrationState = "INITIAL"
	DatabaseMigrationStateLOGSHIPPINGSTART      DatabaseMigrationState = "LOG_SHIPPING_START"
	DatabaseMigrationStatePOSTCUTOVERCOMPLETE   DatabaseMigrationState = "POST_CUTOVER_COMPLETE"
	DatabaseMigrationStateUNDEFINED             DatabaseMigrationState = "UNDEFINED"
	DatabaseMigrationStateUPLOADLOGFILESSTART   DatabaseMigrationState = "UPLOAD_LOG_FILES_START"
)

func PossibleValuesForDatabaseMigrationState() []string {
	return []string{
		string(DatabaseMigrationStateCANCELLED),
		string(DatabaseMigrationStateCOMPLETED),
		string(DatabaseMigrationStateCUTOVERSTART),
		string(DatabaseMigrationStateFAILED),
		string(DatabaseMigrationStateFULLBACKUPUPLOADSTART),
		string(DatabaseMigrationStateINITIAL),
		string(DatabaseMigrationStateLOGSHIPPINGSTART),
		string(DatabaseMigrationStatePOSTCUTOVERCOMPLETE),
		string(DatabaseMigrationStateUNDEFINED),
		string(DatabaseMigrationStateUPLOADLOGFILESSTART),
	}
}

func (s *DatabaseMigrationState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatabaseMigrationState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatabaseMigrationState(input string) (*DatabaseMigrationState, error) {
	vals := map[string]DatabaseMigrationState{
		"cancelled":                DatabaseMigrationStateCANCELLED,
		"completed":                DatabaseMigrationStateCOMPLETED,
		"cutover_start":            DatabaseMigrationStateCUTOVERSTART,
		"failed":                   DatabaseMigrationStateFAILED,
		"full_backup_upload_start": DatabaseMigrationStateFULLBACKUPUPLOADSTART,
		"initial":                  DatabaseMigrationStateINITIAL,
		"log_shipping_start":       DatabaseMigrationStateLOGSHIPPINGSTART,
		"post_cutover_complete":    DatabaseMigrationStatePOSTCUTOVERCOMPLETE,
		"undefined":                DatabaseMigrationStateUNDEFINED,
		"upload_log_files_start":   DatabaseMigrationStateUPLOADLOGFILESSTART,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatabaseMigrationState(input)
	return &out, nil
}

type DatabaseState string

const (
	DatabaseStateCopying          DatabaseState = "Copying"
	DatabaseStateEmergency        DatabaseState = "Emergency"
	DatabaseStateOffline          DatabaseState = "Offline"
	DatabaseStateOfflineSecondary DatabaseState = "OfflineSecondary"
	DatabaseStateOnline           DatabaseState = "Online"
	DatabaseStateRecovering       DatabaseState = "Recovering"
	DatabaseStateRecoveryPending  DatabaseState = "RecoveryPending"
	DatabaseStateRestoring        DatabaseState = "Restoring"
	DatabaseStateSuspect          DatabaseState = "Suspect"
)

func PossibleValuesForDatabaseState() []string {
	return []string{
		string(DatabaseStateCopying),
		string(DatabaseStateEmergency),
		string(DatabaseStateOffline),
		string(DatabaseStateOfflineSecondary),
		string(DatabaseStateOnline),
		string(DatabaseStateRecovering),
		string(DatabaseStateRecoveryPending),
		string(DatabaseStateRestoring),
		string(DatabaseStateSuspect),
	}
}

func (s *DatabaseState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatabaseState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatabaseState(input string) (*DatabaseState, error) {
	vals := map[string]DatabaseState{
		"copying":          DatabaseStateCopying,
		"emergency":        DatabaseStateEmergency,
		"offline":          DatabaseStateOffline,
		"offlinesecondary": DatabaseStateOfflineSecondary,
		"online":           DatabaseStateOnline,
		"recovering":       DatabaseStateRecovering,
		"recoverypending":  DatabaseStateRecoveryPending,
		"restoring":        DatabaseStateRestoring,
		"suspect":          DatabaseStateSuspect,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatabaseState(input)
	return &out, nil
}

type LoginMigrationStage string

const (
	LoginMigrationStageAssignRoleMembership       LoginMigrationStage = "AssignRoleMembership"
	LoginMigrationStageAssignRoleOwnership        LoginMigrationStage = "AssignRoleOwnership"
	LoginMigrationStageCompleted                  LoginMigrationStage = "Completed"
	LoginMigrationStageEstablishObjectPermissions LoginMigrationStage = "EstablishObjectPermissions"
	LoginMigrationStageEstablishServerPermissions LoginMigrationStage = "EstablishServerPermissions"
	LoginMigrationStageEstablishUserMapping       LoginMigrationStage = "EstablishUserMapping"
	LoginMigrationStageInitialize                 LoginMigrationStage = "Initialize"
	LoginMigrationStageLoginMigration             LoginMigrationStage = "LoginMigration"
	LoginMigrationStageNone                       LoginMigrationStage = "None"
)

func PossibleValuesForLoginMigrationStage() []string {
	return []string{
		string(LoginMigrationStageAssignRoleMembership),
		string(LoginMigrationStageAssignRoleOwnership),
		string(LoginMigrationStageCompleted),
		string(LoginMigrationStageEstablishObjectPermissions),
		string(LoginMigrationStageEstablishServerPermissions),
		string(LoginMigrationStageEstablishUserMapping),
		string(LoginMigrationStageInitialize),
		string(LoginMigrationStageLoginMigration),
		string(LoginMigrationStageNone),
	}
}

func (s *LoginMigrationStage) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLoginMigrationStage(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLoginMigrationStage(input string) (*LoginMigrationStage, error) {
	vals := map[string]LoginMigrationStage{
		"assignrolemembership":       LoginMigrationStageAssignRoleMembership,
		"assignroleownership":        LoginMigrationStageAssignRoleOwnership,
		"completed":                  LoginMigrationStageCompleted,
		"establishobjectpermissions": LoginMigrationStageEstablishObjectPermissions,
		"establishserverpermissions": LoginMigrationStageEstablishServerPermissions,
		"establishusermapping":       LoginMigrationStageEstablishUserMapping,
		"initialize":                 LoginMigrationStageInitialize,
		"loginmigration":             LoginMigrationStageLoginMigration,
		"none":                       LoginMigrationStageNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LoginMigrationStage(input)
	return &out, nil
}

type LoginType string

const (
	LoginTypeAsymmetricKey LoginType = "AsymmetricKey"
	LoginTypeCertificate   LoginType = "Certificate"
	LoginTypeExternalGroup LoginType = "ExternalGroup"
	LoginTypeExternalUser  LoginType = "ExternalUser"
	LoginTypeSqlLogin      LoginType = "SqlLogin"
	LoginTypeWindowsGroup  LoginType = "WindowsGroup"
	LoginTypeWindowsUser   LoginType = "WindowsUser"
)

func PossibleValuesForLoginType() []string {
	return []string{
		string(LoginTypeAsymmetricKey),
		string(LoginTypeCertificate),
		string(LoginTypeExternalGroup),
		string(LoginTypeExternalUser),
		string(LoginTypeSqlLogin),
		string(LoginTypeWindowsGroup),
		string(LoginTypeWindowsUser),
	}
}

func (s *LoginType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLoginType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLoginType(input string) (*LoginType, error) {
	vals := map[string]LoginType{
		"asymmetrickey": LoginTypeAsymmetricKey,
		"certificate":   LoginTypeCertificate,
		"externalgroup": LoginTypeExternalGroup,
		"externaluser":  LoginTypeExternalUser,
		"sqllogin":      LoginTypeSqlLogin,
		"windowsgroup":  LoginTypeWindowsGroup,
		"windowsuser":   LoginTypeWindowsUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LoginType(input)
	return &out, nil
}

type MigrationState string

const (
	MigrationStateCompleted  MigrationState = "Completed"
	MigrationStateFailed     MigrationState = "Failed"
	MigrationStateInProgress MigrationState = "InProgress"
	MigrationStateNone       MigrationState = "None"
	MigrationStateSkipped    MigrationState = "Skipped"
	MigrationStateStopped    MigrationState = "Stopped"
	MigrationStateWarning    MigrationState = "Warning"
)

func PossibleValuesForMigrationState() []string {
	return []string{
		string(MigrationStateCompleted),
		string(MigrationStateFailed),
		string(MigrationStateInProgress),
		string(MigrationStateNone),
		string(MigrationStateSkipped),
		string(MigrationStateStopped),
		string(MigrationStateWarning),
	}
}

func (s *MigrationState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMigrationState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMigrationState(input string) (*MigrationState, error) {
	vals := map[string]MigrationState{
		"completed":  MigrationStateCompleted,
		"failed":     MigrationStateFailed,
		"inprogress": MigrationStateInProgress,
		"none":       MigrationStateNone,
		"skipped":    MigrationStateSkipped,
		"stopped":    MigrationStateStopped,
		"warning":    MigrationStateWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MigrationState(input)
	return &out, nil
}

type MigrationStatus string

const (
	MigrationStatusCompleted               MigrationStatus = "Completed"
	MigrationStatusCompletedWithWarnings   MigrationStatus = "CompletedWithWarnings"
	MigrationStatusConfigured              MigrationStatus = "Configured"
	MigrationStatusConnecting              MigrationStatus = "Connecting"
	MigrationStatusDefault                 MigrationStatus = "Default"
	MigrationStatusError                   MigrationStatus = "Error"
	MigrationStatusRunning                 MigrationStatus = "Running"
	MigrationStatusSelectLogins            MigrationStatus = "SelectLogins"
	MigrationStatusSourceAndTargetSelected MigrationStatus = "SourceAndTargetSelected"
	MigrationStatusStopped                 MigrationStatus = "Stopped"
)

func PossibleValuesForMigrationStatus() []string {
	return []string{
		string(MigrationStatusCompleted),
		string(MigrationStatusCompletedWithWarnings),
		string(MigrationStatusConfigured),
		string(MigrationStatusConnecting),
		string(MigrationStatusDefault),
		string(MigrationStatusError),
		string(MigrationStatusRunning),
		string(MigrationStatusSelectLogins),
		string(MigrationStatusSourceAndTargetSelected),
		string(MigrationStatusStopped),
	}
}

func (s *MigrationStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMigrationStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMigrationStatus(input string) (*MigrationStatus, error) {
	vals := map[string]MigrationStatus{
		"completed":               MigrationStatusCompleted,
		"completedwithwarnings":   MigrationStatusCompletedWithWarnings,
		"configured":              MigrationStatusConfigured,
		"connecting":              MigrationStatusConnecting,
		"default":                 MigrationStatusDefault,
		"error":                   MigrationStatusError,
		"running":                 MigrationStatusRunning,
		"selectlogins":            MigrationStatusSelectLogins,
		"sourceandtargetselected": MigrationStatusSourceAndTargetSelected,
		"stopped":                 MigrationStatusStopped,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MigrationStatus(input)
	return &out, nil
}

type MongoDbClusterType string

const (
	MongoDbClusterTypeBlobContainer MongoDbClusterType = "BlobContainer"
	MongoDbClusterTypeCosmosDb      MongoDbClusterType = "CosmosDb"
	MongoDbClusterTypeMongoDb       MongoDbClusterType = "MongoDb"
)

func PossibleValuesForMongoDbClusterType() []string {
	return []string{
		string(MongoDbClusterTypeBlobContainer),
		string(MongoDbClusterTypeCosmosDb),
		string(MongoDbClusterTypeMongoDb),
	}
}

func (s *MongoDbClusterType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMongoDbClusterType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMongoDbClusterType(input string) (*MongoDbClusterType, error) {
	vals := map[string]MongoDbClusterType{
		"blobcontainer": MongoDbClusterTypeBlobContainer,
		"cosmosdb":      MongoDbClusterTypeCosmosDb,
		"mongodb":       MongoDbClusterTypeMongoDb,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MongoDbClusterType(input)
	return &out, nil
}

type MongoDbErrorType string

const (
	MongoDbErrorTypeError           MongoDbErrorType = "Error"
	MongoDbErrorTypeValidationError MongoDbErrorType = "ValidationError"
	MongoDbErrorTypeWarning         MongoDbErrorType = "Warning"
)

func PossibleValuesForMongoDbErrorType() []string {
	return []string{
		string(MongoDbErrorTypeError),
		string(MongoDbErrorTypeValidationError),
		string(MongoDbErrorTypeWarning),
	}
}

func (s *MongoDbErrorType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMongoDbErrorType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMongoDbErrorType(input string) (*MongoDbErrorType, error) {
	vals := map[string]MongoDbErrorType{
		"error":           MongoDbErrorTypeError,
		"validationerror": MongoDbErrorTypeValidationError,
		"warning":         MongoDbErrorTypeWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MongoDbErrorType(input)
	return &out, nil
}

type MongoDbMigrationState string

const (
	MongoDbMigrationStateCanceled        MongoDbMigrationState = "Canceled"
	MongoDbMigrationStateComplete        MongoDbMigrationState = "Complete"
	MongoDbMigrationStateCopying         MongoDbMigrationState = "Copying"
	MongoDbMigrationStateFailed          MongoDbMigrationState = "Failed"
	MongoDbMigrationStateFinalizing      MongoDbMigrationState = "Finalizing"
	MongoDbMigrationStateInitialReplay   MongoDbMigrationState = "InitialReplay"
	MongoDbMigrationStateInitializing    MongoDbMigrationState = "Initializing"
	MongoDbMigrationStateNotStarted      MongoDbMigrationState = "NotStarted"
	MongoDbMigrationStateReplaying       MongoDbMigrationState = "Replaying"
	MongoDbMigrationStateRestarting      MongoDbMigrationState = "Restarting"
	MongoDbMigrationStateValidatingInput MongoDbMigrationState = "ValidatingInput"
)

func PossibleValuesForMongoDbMigrationState() []string {
	return []string{
		string(MongoDbMigrationStateCanceled),
		string(MongoDbMigrationStateComplete),
		string(MongoDbMigrationStateCopying),
		string(MongoDbMigrationStateFailed),
		string(MongoDbMigrationStateFinalizing),
		string(MongoDbMigrationStateInitialReplay),
		string(MongoDbMigrationStateInitializing),
		string(MongoDbMigrationStateNotStarted),
		string(MongoDbMigrationStateReplaying),
		string(MongoDbMigrationStateRestarting),
		string(MongoDbMigrationStateValidatingInput),
	}
}

func (s *MongoDbMigrationState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMongoDbMigrationState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMongoDbMigrationState(input string) (*MongoDbMigrationState, error) {
	vals := map[string]MongoDbMigrationState{
		"canceled":        MongoDbMigrationStateCanceled,
		"complete":        MongoDbMigrationStateComplete,
		"copying":         MongoDbMigrationStateCopying,
		"failed":          MongoDbMigrationStateFailed,
		"finalizing":      MongoDbMigrationStateFinalizing,
		"initialreplay":   MongoDbMigrationStateInitialReplay,
		"initializing":    MongoDbMigrationStateInitializing,
		"notstarted":      MongoDbMigrationStateNotStarted,
		"replaying":       MongoDbMigrationStateReplaying,
		"restarting":      MongoDbMigrationStateRestarting,
		"validatinginput": MongoDbMigrationStateValidatingInput,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MongoDbMigrationState(input)
	return &out, nil
}

type MongoDbReplication string

const (
	MongoDbReplicationContinuous MongoDbReplication = "Continuous"
	MongoDbReplicationDisabled   MongoDbReplication = "Disabled"
	MongoDbReplicationOneTime    MongoDbReplication = "OneTime"
)

func PossibleValuesForMongoDbReplication() []string {
	return []string{
		string(MongoDbReplicationContinuous),
		string(MongoDbReplicationDisabled),
		string(MongoDbReplicationOneTime),
	}
}

func (s *MongoDbReplication) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMongoDbReplication(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMongoDbReplication(input string) (*MongoDbReplication, error) {
	vals := map[string]MongoDbReplication{
		"continuous": MongoDbReplicationContinuous,
		"disabled":   MongoDbReplicationDisabled,
		"onetime":    MongoDbReplicationOneTime,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MongoDbReplication(input)
	return &out, nil
}

type MongoDbShardKeyOrder string

const (
	MongoDbShardKeyOrderForward MongoDbShardKeyOrder = "Forward"
	MongoDbShardKeyOrderHashed  MongoDbShardKeyOrder = "Hashed"
	MongoDbShardKeyOrderReverse MongoDbShardKeyOrder = "Reverse"
)

func PossibleValuesForMongoDbShardKeyOrder() []string {
	return []string{
		string(MongoDbShardKeyOrderForward),
		string(MongoDbShardKeyOrderHashed),
		string(MongoDbShardKeyOrderReverse),
	}
}

func (s *MongoDbShardKeyOrder) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMongoDbShardKeyOrder(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMongoDbShardKeyOrder(input string) (*MongoDbShardKeyOrder, error) {
	vals := map[string]MongoDbShardKeyOrder{
		"forward": MongoDbShardKeyOrderForward,
		"hashed":  MongoDbShardKeyOrderHashed,
		"reverse": MongoDbShardKeyOrderReverse,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MongoDbShardKeyOrder(input)
	return &out, nil
}

type MySqlTargetPlatformType string

const (
	MySqlTargetPlatformTypeAzureDbForMySQL MySqlTargetPlatformType = "AzureDbForMySQL"
	MySqlTargetPlatformTypeSqlServer       MySqlTargetPlatformType = "SqlServer"
)

func PossibleValuesForMySqlTargetPlatformType() []string {
	return []string{
		string(MySqlTargetPlatformTypeAzureDbForMySQL),
		string(MySqlTargetPlatformTypeSqlServer),
	}
}

func (s *MySqlTargetPlatformType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMySqlTargetPlatformType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMySqlTargetPlatformType(input string) (*MySqlTargetPlatformType, error) {
	vals := map[string]MySqlTargetPlatformType{
		"azuredbformysql": MySqlTargetPlatformTypeAzureDbForMySQL,
		"sqlserver":       MySqlTargetPlatformTypeSqlServer,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MySqlTargetPlatformType(input)
	return &out, nil
}

type ObjectType string

const (
	ObjectTypeFunction         ObjectType = "Function"
	ObjectTypeStoredProcedures ObjectType = "StoredProcedures"
	ObjectTypeTable            ObjectType = "Table"
	ObjectTypeUser             ObjectType = "User"
	ObjectTypeView             ObjectType = "View"
)

func PossibleValuesForObjectType() []string {
	return []string{
		string(ObjectTypeFunction),
		string(ObjectTypeStoredProcedures),
		string(ObjectTypeTable),
		string(ObjectTypeUser),
		string(ObjectTypeView),
	}
}

func (s *ObjectType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseObjectType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseObjectType(input string) (*ObjectType, error) {
	vals := map[string]ObjectType{
		"function":         ObjectTypeFunction,
		"storedprocedures": ObjectTypeStoredProcedures,
		"table":            ObjectTypeTable,
		"user":             ObjectTypeUser,
		"view":             ObjectTypeView,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ObjectType(input)
	return &out, nil
}

type ReplicateMigrationState string

const (
	ReplicateMigrationStateACTIONREQUIRED ReplicateMigrationState = "ACTION_REQUIRED"
	ReplicateMigrationStateCOMPLETE       ReplicateMigrationState = "COMPLETE"
	ReplicateMigrationStateFAILED         ReplicateMigrationState = "FAILED"
	ReplicateMigrationStatePENDING        ReplicateMigrationState = "PENDING"
	ReplicateMigrationStateUNDEFINED      ReplicateMigrationState = "UNDEFINED"
	ReplicateMigrationStateVALIDATING     ReplicateMigrationState = "VALIDATING"
)

func PossibleValuesForReplicateMigrationState() []string {
	return []string{
		string(ReplicateMigrationStateACTIONREQUIRED),
		string(ReplicateMigrationStateCOMPLETE),
		string(ReplicateMigrationStateFAILED),
		string(ReplicateMigrationStatePENDING),
		string(ReplicateMigrationStateUNDEFINED),
		string(ReplicateMigrationStateVALIDATING),
	}
}

func (s *ReplicateMigrationState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReplicateMigrationState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReplicateMigrationState(input string) (*ReplicateMigrationState, error) {
	vals := map[string]ReplicateMigrationState{
		"action_required": ReplicateMigrationStateACTIONREQUIRED,
		"complete":        ReplicateMigrationStateCOMPLETE,
		"failed":          ReplicateMigrationStateFAILED,
		"pending":         ReplicateMigrationStatePENDING,
		"undefined":       ReplicateMigrationStateUNDEFINED,
		"validating":      ReplicateMigrationStateVALIDATING,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicateMigrationState(input)
	return &out, nil
}

type ResultType string

const (
	ResultTypeCollection ResultType = "Collection"
	ResultTypeDatabase   ResultType = "Database"
	ResultTypeMigration  ResultType = "Migration"
)

func PossibleValuesForResultType() []string {
	return []string{
		string(ResultTypeCollection),
		string(ResultTypeDatabase),
		string(ResultTypeMigration),
	}
}

func (s *ResultType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResultType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResultType(input string) (*ResultType, error) {
	vals := map[string]ResultType{
		"collection": ResultTypeCollection,
		"database":   ResultTypeDatabase,
		"migration":  ResultTypeMigration,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResultType(input)
	return &out, nil
}

type ScenarioSource string

const (
	ScenarioSourceAccess        ScenarioSource = "Access"
	ScenarioSourceDBTwo         ScenarioSource = "DB2"
	ScenarioSourceMongoDB       ScenarioSource = "MongoDB"
	ScenarioSourceMySQL         ScenarioSource = "MySQL"
	ScenarioSourceMySQLRDS      ScenarioSource = "MySQLRDS"
	ScenarioSourceOracle        ScenarioSource = "Oracle"
	ScenarioSourcePostgreSQL    ScenarioSource = "PostgreSQL"
	ScenarioSourcePostgreSQLRDS ScenarioSource = "PostgreSQLRDS"
	ScenarioSourceSQL           ScenarioSource = "SQL"
	ScenarioSourceSQLRDS        ScenarioSource = "SQLRDS"
	ScenarioSourceSybase        ScenarioSource = "Sybase"
)

func PossibleValuesForScenarioSource() []string {
	return []string{
		string(ScenarioSourceAccess),
		string(ScenarioSourceDBTwo),
		string(ScenarioSourceMongoDB),
		string(ScenarioSourceMySQL),
		string(ScenarioSourceMySQLRDS),
		string(ScenarioSourceOracle),
		string(ScenarioSourcePostgreSQL),
		string(ScenarioSourcePostgreSQLRDS),
		string(ScenarioSourceSQL),
		string(ScenarioSourceSQLRDS),
		string(ScenarioSourceSybase),
	}
}

func (s *ScenarioSource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScenarioSource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScenarioSource(input string) (*ScenarioSource, error) {
	vals := map[string]ScenarioSource{
		"access":        ScenarioSourceAccess,
		"db2":           ScenarioSourceDBTwo,
		"mongodb":       ScenarioSourceMongoDB,
		"mysql":         ScenarioSourceMySQL,
		"mysqlrds":      ScenarioSourceMySQLRDS,
		"oracle":        ScenarioSourceOracle,
		"postgresql":    ScenarioSourcePostgreSQL,
		"postgresqlrds": ScenarioSourcePostgreSQLRDS,
		"sql":           ScenarioSourceSQL,
		"sqlrds":        ScenarioSourceSQLRDS,
		"sybase":        ScenarioSourceSybase,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScenarioSource(input)
	return &out, nil
}

type ScenarioTarget string

const (
	ScenarioTargetAzureDBForMySql       ScenarioTarget = "AzureDBForMySql"
	ScenarioTargetAzureDBForPostgresSQL ScenarioTarget = "AzureDBForPostgresSQL"
	ScenarioTargetMongoDB               ScenarioTarget = "MongoDB"
	ScenarioTargetSQLDB                 ScenarioTarget = "SQLDB"
	ScenarioTargetSQLDW                 ScenarioTarget = "SQLDW"
	ScenarioTargetSQLMI                 ScenarioTarget = "SQLMI"
	ScenarioTargetSQLServer             ScenarioTarget = "SQLServer"
)

func PossibleValuesForScenarioTarget() []string {
	return []string{
		string(ScenarioTargetAzureDBForMySql),
		string(ScenarioTargetAzureDBForPostgresSQL),
		string(ScenarioTargetMongoDB),
		string(ScenarioTargetSQLDB),
		string(ScenarioTargetSQLDW),
		string(ScenarioTargetSQLMI),
		string(ScenarioTargetSQLServer),
	}
}

func (s *ScenarioTarget) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScenarioTarget(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScenarioTarget(input string) (*ScenarioTarget, error) {
	vals := map[string]ScenarioTarget{
		"azuredbformysql":       ScenarioTargetAzureDBForMySql,
		"azuredbforpostgressql": ScenarioTargetAzureDBForPostgresSQL,
		"mongodb":               ScenarioTargetMongoDB,
		"sqldb":                 ScenarioTargetSQLDB,
		"sqldw":                 ScenarioTargetSQLDW,
		"sqlmi":                 ScenarioTargetSQLMI,
		"sqlserver":             ScenarioTargetSQLServer,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScenarioTarget(input)
	return &out, nil
}

type ServerLevelPermissionsGroup string

const (
	ServerLevelPermissionsGroupDefault                             ServerLevelPermissionsGroup = "Default"
	ServerLevelPermissionsGroupMigrationFromMySQLToAzureDBForMySQL ServerLevelPermissionsGroup = "MigrationFromMySQLToAzureDBForMySQL"
	ServerLevelPermissionsGroupMigrationFromSqlServerToAzureDB     ServerLevelPermissionsGroup = "MigrationFromSqlServerToAzureDB"
	ServerLevelPermissionsGroupMigrationFromSqlServerToAzureMI     ServerLevelPermissionsGroup = "MigrationFromSqlServerToAzureMI"
)

func PossibleValuesForServerLevelPermissionsGroup() []string {
	return []string{
		string(ServerLevelPermissionsGroupDefault),
		string(ServerLevelPermissionsGroupMigrationFromMySQLToAzureDBForMySQL),
		string(ServerLevelPermissionsGroupMigrationFromSqlServerToAzureDB),
		string(ServerLevelPermissionsGroupMigrationFromSqlServerToAzureMI),
	}
}

func (s *ServerLevelPermissionsGroup) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServerLevelPermissionsGroup(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServerLevelPermissionsGroup(input string) (*ServerLevelPermissionsGroup, error) {
	vals := map[string]ServerLevelPermissionsGroup{
		"default":                             ServerLevelPermissionsGroupDefault,
		"migrationfrommysqltoazuredbformysql": ServerLevelPermissionsGroupMigrationFromMySQLToAzureDBForMySQL,
		"migrationfromsqlservertoazuredb":     ServerLevelPermissionsGroupMigrationFromSqlServerToAzureDB,
		"migrationfromsqlservertoazuremi":     ServerLevelPermissionsGroupMigrationFromSqlServerToAzureMI,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerLevelPermissionsGroup(input)
	return &out, nil
}

type ServiceProvisioningState string

const (
	ServiceProvisioningStateAccepted      ServiceProvisioningState = "Accepted"
	ServiceProvisioningStateDeleting      ServiceProvisioningState = "Deleting"
	ServiceProvisioningStateDeploying     ServiceProvisioningState = "Deploying"
	ServiceProvisioningStateFailed        ServiceProvisioningState = "Failed"
	ServiceProvisioningStateFailedToStart ServiceProvisioningState = "FailedToStart"
	ServiceProvisioningStateFailedToStop  ServiceProvisioningState = "FailedToStop"
	ServiceProvisioningStateStarting      ServiceProvisioningState = "Starting"
	ServiceProvisioningStateStopped       ServiceProvisioningState = "Stopped"
	ServiceProvisioningStateStopping      ServiceProvisioningState = "Stopping"
	ServiceProvisioningStateSucceeded     ServiceProvisioningState = "Succeeded"
)

func PossibleValuesForServiceProvisioningState() []string {
	return []string{
		string(ServiceProvisioningStateAccepted),
		string(ServiceProvisioningStateDeleting),
		string(ServiceProvisioningStateDeploying),
		string(ServiceProvisioningStateFailed),
		string(ServiceProvisioningStateFailedToStart),
		string(ServiceProvisioningStateFailedToStop),
		string(ServiceProvisioningStateStarting),
		string(ServiceProvisioningStateStopped),
		string(ServiceProvisioningStateStopping),
		string(ServiceProvisioningStateSucceeded),
	}
}

func (s *ServiceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServiceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServiceProvisioningState(input string) (*ServiceProvisioningState, error) {
	vals := map[string]ServiceProvisioningState{
		"accepted":      ServiceProvisioningStateAccepted,
		"deleting":      ServiceProvisioningStateDeleting,
		"deploying":     ServiceProvisioningStateDeploying,
		"failed":        ServiceProvisioningStateFailed,
		"failedtostart": ServiceProvisioningStateFailedToStart,
		"failedtostop":  ServiceProvisioningStateFailedToStop,
		"starting":      ServiceProvisioningStateStarting,
		"stopped":       ServiceProvisioningStateStopped,
		"stopping":      ServiceProvisioningStateStopping,
		"succeeded":     ServiceProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceProvisioningState(input)
	return &out, nil
}

type ServiceScalability string

const (
	ServiceScalabilityAutomatic ServiceScalability = "automatic"
	ServiceScalabilityManual    ServiceScalability = "manual"
	ServiceScalabilityNone      ServiceScalability = "none"
)

func PossibleValuesForServiceScalability() []string {
	return []string{
		string(ServiceScalabilityAutomatic),
		string(ServiceScalabilityManual),
		string(ServiceScalabilityNone),
	}
}

func (s *ServiceScalability) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServiceScalability(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServiceScalability(input string) (*ServiceScalability, error) {
	vals := map[string]ServiceScalability{
		"automatic": ServiceScalabilityAutomatic,
		"manual":    ServiceScalabilityManual,
		"none":      ServiceScalabilityNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceScalability(input)
	return &out, nil
}

type Severity string

const (
	SeverityError   Severity = "Error"
	SeverityMessage Severity = "Message"
	SeverityWarning Severity = "Warning"
)

func PossibleValuesForSeverity() []string {
	return []string{
		string(SeverityError),
		string(SeverityMessage),
		string(SeverityWarning),
	}
}

func (s *Severity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSeverity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSeverity(input string) (*Severity, error) {
	vals := map[string]Severity{
		"error":   SeverityError,
		"message": SeverityMessage,
		"warning": SeverityWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Severity(input)
	return &out, nil
}

type SqlSourcePlatform string

const (
	SqlSourcePlatformSqlOnPrem SqlSourcePlatform = "SqlOnPrem"
)

func PossibleValuesForSqlSourcePlatform() []string {
	return []string{
		string(SqlSourcePlatformSqlOnPrem),
	}
}

func (s *SqlSourcePlatform) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSqlSourcePlatform(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSqlSourcePlatform(input string) (*SqlSourcePlatform, error) {
	vals := map[string]SqlSourcePlatform{
		"sqlonprem": SqlSourcePlatformSqlOnPrem,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SqlSourcePlatform(input)
	return &out, nil
}

type SsisMigrationOverwriteOption string

const (
	SsisMigrationOverwriteOptionIgnore    SsisMigrationOverwriteOption = "Ignore"
	SsisMigrationOverwriteOptionOverwrite SsisMigrationOverwriteOption = "Overwrite"
)

func PossibleValuesForSsisMigrationOverwriteOption() []string {
	return []string{
		string(SsisMigrationOverwriteOptionIgnore),
		string(SsisMigrationOverwriteOptionOverwrite),
	}
}

func (s *SsisMigrationOverwriteOption) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSsisMigrationOverwriteOption(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSsisMigrationOverwriteOption(input string) (*SsisMigrationOverwriteOption, error) {
	vals := map[string]SsisMigrationOverwriteOption{
		"ignore":    SsisMigrationOverwriteOptionIgnore,
		"overwrite": SsisMigrationOverwriteOptionOverwrite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SsisMigrationOverwriteOption(input)
	return &out, nil
}

type SsisMigrationStage string

const (
	SsisMigrationStageCompleted  SsisMigrationStage = "Completed"
	SsisMigrationStageInProgress SsisMigrationStage = "InProgress"
	SsisMigrationStageInitialize SsisMigrationStage = "Initialize"
	SsisMigrationStageNone       SsisMigrationStage = "None"
)

func PossibleValuesForSsisMigrationStage() []string {
	return []string{
		string(SsisMigrationStageCompleted),
		string(SsisMigrationStageInProgress),
		string(SsisMigrationStageInitialize),
		string(SsisMigrationStageNone),
	}
}

func (s *SsisMigrationStage) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSsisMigrationStage(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSsisMigrationStage(input string) (*SsisMigrationStage, error) {
	vals := map[string]SsisMigrationStage{
		"completed":  SsisMigrationStageCompleted,
		"inprogress": SsisMigrationStageInProgress,
		"initialize": SsisMigrationStageInitialize,
		"none":       SsisMigrationStageNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SsisMigrationStage(input)
	return &out, nil
}

type SsisStoreType string

const (
	SsisStoreTypeSsisCatalog SsisStoreType = "SsisCatalog"
)

func PossibleValuesForSsisStoreType() []string {
	return []string{
		string(SsisStoreTypeSsisCatalog),
	}
}

func (s *SsisStoreType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSsisStoreType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSsisStoreType(input string) (*SsisStoreType, error) {
	vals := map[string]SsisStoreType{
		"ssiscatalog": SsisStoreTypeSsisCatalog,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SsisStoreType(input)
	return &out, nil
}

type SyncDatabaseMigrationReportingState string

const (
	SyncDatabaseMigrationReportingStateBACKUPCOMPLETED    SyncDatabaseMigrationReportingState = "BACKUP_COMPLETED"
	SyncDatabaseMigrationReportingStateBACKUPINPROGRESS   SyncDatabaseMigrationReportingState = "BACKUP_IN_PROGRESS"
	SyncDatabaseMigrationReportingStateCANCELLED          SyncDatabaseMigrationReportingState = "CANCELLED"
	SyncDatabaseMigrationReportingStateCANCELLING         SyncDatabaseMigrationReportingState = "CANCELLING"
	SyncDatabaseMigrationReportingStateCOMPLETE           SyncDatabaseMigrationReportingState = "COMPLETE"
	SyncDatabaseMigrationReportingStateCOMPLETING         SyncDatabaseMigrationReportingState = "COMPLETING"
	SyncDatabaseMigrationReportingStateCONFIGURING        SyncDatabaseMigrationReportingState = "CONFIGURING"
	SyncDatabaseMigrationReportingStateFAILED             SyncDatabaseMigrationReportingState = "FAILED"
	SyncDatabaseMigrationReportingStateINITIALIAZING      SyncDatabaseMigrationReportingState = "INITIALIAZING"
	SyncDatabaseMigrationReportingStateREADYTOCOMPLETE    SyncDatabaseMigrationReportingState = "READY_TO_COMPLETE"
	SyncDatabaseMigrationReportingStateRESTORECOMPLETED   SyncDatabaseMigrationReportingState = "RESTORE_COMPLETED"
	SyncDatabaseMigrationReportingStateRESTOREINPROGRESS  SyncDatabaseMigrationReportingState = "RESTORE_IN_PROGRESS"
	SyncDatabaseMigrationReportingStateRUNNING            SyncDatabaseMigrationReportingState = "RUNNING"
	SyncDatabaseMigrationReportingStateSTARTING           SyncDatabaseMigrationReportingState = "STARTING"
	SyncDatabaseMigrationReportingStateUNDEFINED          SyncDatabaseMigrationReportingState = "UNDEFINED"
	SyncDatabaseMigrationReportingStateVALIDATING         SyncDatabaseMigrationReportingState = "VALIDATING"
	SyncDatabaseMigrationReportingStateVALIDATIONCOMPLETE SyncDatabaseMigrationReportingState = "VALIDATION_COMPLETE"
	SyncDatabaseMigrationReportingStateVALIDATIONFAILED   SyncDatabaseMigrationReportingState = "VALIDATION_FAILED"
)

func PossibleValuesForSyncDatabaseMigrationReportingState() []string {
	return []string{
		string(SyncDatabaseMigrationReportingStateBACKUPCOMPLETED),
		string(SyncDatabaseMigrationReportingStateBACKUPINPROGRESS),
		string(SyncDatabaseMigrationReportingStateCANCELLED),
		string(SyncDatabaseMigrationReportingStateCANCELLING),
		string(SyncDatabaseMigrationReportingStateCOMPLETE),
		string(SyncDatabaseMigrationReportingStateCOMPLETING),
		string(SyncDatabaseMigrationReportingStateCONFIGURING),
		string(SyncDatabaseMigrationReportingStateFAILED),
		string(SyncDatabaseMigrationReportingStateINITIALIAZING),
		string(SyncDatabaseMigrationReportingStateREADYTOCOMPLETE),
		string(SyncDatabaseMigrationReportingStateRESTORECOMPLETED),
		string(SyncDatabaseMigrationReportingStateRESTOREINPROGRESS),
		string(SyncDatabaseMigrationReportingStateRUNNING),
		string(SyncDatabaseMigrationReportingStateSTARTING),
		string(SyncDatabaseMigrationReportingStateUNDEFINED),
		string(SyncDatabaseMigrationReportingStateVALIDATING),
		string(SyncDatabaseMigrationReportingStateVALIDATIONCOMPLETE),
		string(SyncDatabaseMigrationReportingStateVALIDATIONFAILED),
	}
}

func (s *SyncDatabaseMigrationReportingState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSyncDatabaseMigrationReportingState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSyncDatabaseMigrationReportingState(input string) (*SyncDatabaseMigrationReportingState, error) {
	vals := map[string]SyncDatabaseMigrationReportingState{
		"backup_completed":    SyncDatabaseMigrationReportingStateBACKUPCOMPLETED,
		"backup_in_progress":  SyncDatabaseMigrationReportingStateBACKUPINPROGRESS,
		"cancelled":           SyncDatabaseMigrationReportingStateCANCELLED,
		"cancelling":          SyncDatabaseMigrationReportingStateCANCELLING,
		"complete":            SyncDatabaseMigrationReportingStateCOMPLETE,
		"completing":          SyncDatabaseMigrationReportingStateCOMPLETING,
		"configuring":         SyncDatabaseMigrationReportingStateCONFIGURING,
		"failed":              SyncDatabaseMigrationReportingStateFAILED,
		"initialiazing":       SyncDatabaseMigrationReportingStateINITIALIAZING,
		"ready_to_complete":   SyncDatabaseMigrationReportingStateREADYTOCOMPLETE,
		"restore_completed":   SyncDatabaseMigrationReportingStateRESTORECOMPLETED,
		"restore_in_progress": SyncDatabaseMigrationReportingStateRESTOREINPROGRESS,
		"running":             SyncDatabaseMigrationReportingStateRUNNING,
		"starting":            SyncDatabaseMigrationReportingStateSTARTING,
		"undefined":           SyncDatabaseMigrationReportingStateUNDEFINED,
		"validating":          SyncDatabaseMigrationReportingStateVALIDATING,
		"validation_complete": SyncDatabaseMigrationReportingStateVALIDATIONCOMPLETE,
		"validation_failed":   SyncDatabaseMigrationReportingStateVALIDATIONFAILED,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SyncDatabaseMigrationReportingState(input)
	return &out, nil
}

type SyncTableMigrationState string

const (
	SyncTableMigrationStateBEFORELOAD SyncTableMigrationState = "BEFORE_LOAD"
	SyncTableMigrationStateCANCELED   SyncTableMigrationState = "CANCELED"
	SyncTableMigrationStateCOMPLETED  SyncTableMigrationState = "COMPLETED"
	SyncTableMigrationStateERROR      SyncTableMigrationState = "ERROR"
	SyncTableMigrationStateFAILED     SyncTableMigrationState = "FAILED"
	SyncTableMigrationStateFULLLOAD   SyncTableMigrationState = "FULL_LOAD"
)

func PossibleValuesForSyncTableMigrationState() []string {
	return []string{
		string(SyncTableMigrationStateBEFORELOAD),
		string(SyncTableMigrationStateCANCELED),
		string(SyncTableMigrationStateCOMPLETED),
		string(SyncTableMigrationStateERROR),
		string(SyncTableMigrationStateFAILED),
		string(SyncTableMigrationStateFULLLOAD),
	}
}

func (s *SyncTableMigrationState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSyncTableMigrationState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSyncTableMigrationState(input string) (*SyncTableMigrationState, error) {
	vals := map[string]SyncTableMigrationState{
		"before_load": SyncTableMigrationStateBEFORELOAD,
		"canceled":    SyncTableMigrationStateCANCELED,
		"completed":   SyncTableMigrationStateCOMPLETED,
		"error":       SyncTableMigrationStateERROR,
		"failed":      SyncTableMigrationStateFAILED,
		"full_load":   SyncTableMigrationStateFULLLOAD,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SyncTableMigrationState(input)
	return &out, nil
}

type TaskState string

const (
	TaskStateCanceled              TaskState = "Canceled"
	TaskStateFailed                TaskState = "Failed"
	TaskStateFailedInputValidation TaskState = "FailedInputValidation"
	TaskStateFaulted               TaskState = "Faulted"
	TaskStateQueued                TaskState = "Queued"
	TaskStateRunning               TaskState = "Running"
	TaskStateSucceeded             TaskState = "Succeeded"
	TaskStateUnknown               TaskState = "Unknown"
)

func PossibleValuesForTaskState() []string {
	return []string{
		string(TaskStateCanceled),
		string(TaskStateFailed),
		string(TaskStateFailedInputValidation),
		string(TaskStateFaulted),
		string(TaskStateQueued),
		string(TaskStateRunning),
		string(TaskStateSucceeded),
		string(TaskStateUnknown),
	}
}

func (s *TaskState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTaskState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTaskState(input string) (*TaskState, error) {
	vals := map[string]TaskState{
		"canceled":              TaskStateCanceled,
		"failed":                TaskStateFailed,
		"failedinputvalidation": TaskStateFailedInputValidation,
		"faulted":               TaskStateFaulted,
		"queued":                TaskStateQueued,
		"running":               TaskStateRunning,
		"succeeded":             TaskStateSucceeded,
		"unknown":               TaskStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TaskState(input)
	return &out, nil
}

type UpdateActionType string

const (
	UpdateActionTypeAddedOnTarget   UpdateActionType = "AddedOnTarget"
	UpdateActionTypeChangedOnTarget UpdateActionType = "ChangedOnTarget"
	UpdateActionTypeDeletedOnTarget UpdateActionType = "DeletedOnTarget"
)

func PossibleValuesForUpdateActionType() []string {
	return []string{
		string(UpdateActionTypeAddedOnTarget),
		string(UpdateActionTypeChangedOnTarget),
		string(UpdateActionTypeDeletedOnTarget),
	}
}

func (s *UpdateActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUpdateActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUpdateActionType(input string) (*UpdateActionType, error) {
	vals := map[string]UpdateActionType{
		"addedontarget":   UpdateActionTypeAddedOnTarget,
		"changedontarget": UpdateActionTypeChangedOnTarget,
		"deletedontarget": UpdateActionTypeDeletedOnTarget,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UpdateActionType(input)
	return &out, nil
}

type ValidationStatus string

const (
	ValidationStatusCompleted           ValidationStatus = "Completed"
	ValidationStatusCompletedWithIssues ValidationStatus = "CompletedWithIssues"
	ValidationStatusDefault             ValidationStatus = "Default"
	ValidationStatusFailed              ValidationStatus = "Failed"
	ValidationStatusInProgress          ValidationStatus = "InProgress"
	ValidationStatusInitialized         ValidationStatus = "Initialized"
	ValidationStatusNotStarted          ValidationStatus = "NotStarted"
	ValidationStatusStopped             ValidationStatus = "Stopped"
)

func PossibleValuesForValidationStatus() []string {
	return []string{
		string(ValidationStatusCompleted),
		string(ValidationStatusCompletedWithIssues),
		string(ValidationStatusDefault),
		string(ValidationStatusFailed),
		string(ValidationStatusInProgress),
		string(ValidationStatusInitialized),
		string(ValidationStatusNotStarted),
		string(ValidationStatusStopped),
	}
}

func (s *ValidationStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseValidationStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseValidationStatus(input string) (*ValidationStatus, error) {
	vals := map[string]ValidationStatus{
		"completed":           ValidationStatusCompleted,
		"completedwithissues": ValidationStatusCompletedWithIssues,
		"default":             ValidationStatusDefault,
		"failed":              ValidationStatusFailed,
		"inprogress":          ValidationStatusInProgress,
		"initialized":         ValidationStatusInitialized,
		"notstarted":          ValidationStatusNotStarted,
		"stopped":             ValidationStatusStopped,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ValidationStatus(input)
	return &out, nil
}
