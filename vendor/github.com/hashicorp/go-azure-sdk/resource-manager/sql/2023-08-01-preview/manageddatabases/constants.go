package manageddatabases

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CatalogCollationType string

const (
	CatalogCollationTypeDATABASEDEFAULT             CatalogCollationType = "DATABASE_DEFAULT"
	CatalogCollationTypeSQLLatinOneGeneralCPOneCIAS CatalogCollationType = "SQL_Latin1_General_CP1_CI_AS"
)

func PossibleValuesForCatalogCollationType() []string {
	return []string{
		string(CatalogCollationTypeDATABASEDEFAULT),
		string(CatalogCollationTypeSQLLatinOneGeneralCPOneCIAS),
	}
}

func (s *CatalogCollationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCatalogCollationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCatalogCollationType(input string) (*CatalogCollationType, error) {
	vals := map[string]CatalogCollationType{
		"database_default":             CatalogCollationTypeDATABASEDEFAULT,
		"sql_latin1_general_cp1_ci_as": CatalogCollationTypeSQLLatinOneGeneralCPOneCIAS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CatalogCollationType(input)
	return &out, nil
}

type ManagedDatabaseCreateMode string

const (
	ManagedDatabaseCreateModeDefault                        ManagedDatabaseCreateMode = "Default"
	ManagedDatabaseCreateModePointInTimeRestore             ManagedDatabaseCreateMode = "PointInTimeRestore"
	ManagedDatabaseCreateModeRecovery                       ManagedDatabaseCreateMode = "Recovery"
	ManagedDatabaseCreateModeRestoreExternalBackup          ManagedDatabaseCreateMode = "RestoreExternalBackup"
	ManagedDatabaseCreateModeRestoreLongTermRetentionBackup ManagedDatabaseCreateMode = "RestoreLongTermRetentionBackup"
)

func PossibleValuesForManagedDatabaseCreateMode() []string {
	return []string{
		string(ManagedDatabaseCreateModeDefault),
		string(ManagedDatabaseCreateModePointInTimeRestore),
		string(ManagedDatabaseCreateModeRecovery),
		string(ManagedDatabaseCreateModeRestoreExternalBackup),
		string(ManagedDatabaseCreateModeRestoreLongTermRetentionBackup),
	}
}

func (s *ManagedDatabaseCreateMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedDatabaseCreateMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedDatabaseCreateMode(input string) (*ManagedDatabaseCreateMode, error) {
	vals := map[string]ManagedDatabaseCreateMode{
		"default":                        ManagedDatabaseCreateModeDefault,
		"pointintimerestore":             ManagedDatabaseCreateModePointInTimeRestore,
		"recovery":                       ManagedDatabaseCreateModeRecovery,
		"restoreexternalbackup":          ManagedDatabaseCreateModeRestoreExternalBackup,
		"restorelongtermretentionbackup": ManagedDatabaseCreateModeRestoreLongTermRetentionBackup,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedDatabaseCreateMode(input)
	return &out, nil
}

type ManagedDatabaseStatus string

const (
	ManagedDatabaseStatusCreating     ManagedDatabaseStatus = "Creating"
	ManagedDatabaseStatusDbCopying    ManagedDatabaseStatus = "DbCopying"
	ManagedDatabaseStatusDbMoving     ManagedDatabaseStatus = "DbMoving"
	ManagedDatabaseStatusInaccessible ManagedDatabaseStatus = "Inaccessible"
	ManagedDatabaseStatusOffline      ManagedDatabaseStatus = "Offline"
	ManagedDatabaseStatusOnline       ManagedDatabaseStatus = "Online"
	ManagedDatabaseStatusRestoring    ManagedDatabaseStatus = "Restoring"
	ManagedDatabaseStatusShutdown     ManagedDatabaseStatus = "Shutdown"
	ManagedDatabaseStatusStarting     ManagedDatabaseStatus = "Starting"
	ManagedDatabaseStatusStopped      ManagedDatabaseStatus = "Stopped"
	ManagedDatabaseStatusStopping     ManagedDatabaseStatus = "Stopping"
	ManagedDatabaseStatusUpdating     ManagedDatabaseStatus = "Updating"
)

func PossibleValuesForManagedDatabaseStatus() []string {
	return []string{
		string(ManagedDatabaseStatusCreating),
		string(ManagedDatabaseStatusDbCopying),
		string(ManagedDatabaseStatusDbMoving),
		string(ManagedDatabaseStatusInaccessible),
		string(ManagedDatabaseStatusOffline),
		string(ManagedDatabaseStatusOnline),
		string(ManagedDatabaseStatusRestoring),
		string(ManagedDatabaseStatusShutdown),
		string(ManagedDatabaseStatusStarting),
		string(ManagedDatabaseStatusStopped),
		string(ManagedDatabaseStatusStopping),
		string(ManagedDatabaseStatusUpdating),
	}
}

func (s *ManagedDatabaseStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedDatabaseStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedDatabaseStatus(input string) (*ManagedDatabaseStatus, error) {
	vals := map[string]ManagedDatabaseStatus{
		"creating":     ManagedDatabaseStatusCreating,
		"dbcopying":    ManagedDatabaseStatusDbCopying,
		"dbmoving":     ManagedDatabaseStatusDbMoving,
		"inaccessible": ManagedDatabaseStatusInaccessible,
		"offline":      ManagedDatabaseStatusOffline,
		"online":       ManagedDatabaseStatusOnline,
		"restoring":    ManagedDatabaseStatusRestoring,
		"shutdown":     ManagedDatabaseStatusShutdown,
		"starting":     ManagedDatabaseStatusStarting,
		"stopped":      ManagedDatabaseStatusStopped,
		"stopping":     ManagedDatabaseStatusStopping,
		"updating":     ManagedDatabaseStatusUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedDatabaseStatus(input)
	return &out, nil
}

type MoveOperationMode string

const (
	MoveOperationModeCopy MoveOperationMode = "Copy"
	MoveOperationModeMove MoveOperationMode = "Move"
)

func PossibleValuesForMoveOperationMode() []string {
	return []string{
		string(MoveOperationModeCopy),
		string(MoveOperationModeMove),
	}
}

func (s *MoveOperationMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMoveOperationMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMoveOperationMode(input string) (*MoveOperationMode, error) {
	vals := map[string]MoveOperationMode{
		"copy": MoveOperationModeCopy,
		"move": MoveOperationModeMove,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MoveOperationMode(input)
	return &out, nil
}
