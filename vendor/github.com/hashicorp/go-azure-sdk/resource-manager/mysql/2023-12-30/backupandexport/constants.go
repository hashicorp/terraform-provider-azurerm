package backupandexport

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupFormat string

const (
	BackupFormatCollatedFormat BackupFormat = "CollatedFormat"
	BackupFormatRaw            BackupFormat = "Raw"
)

func PossibleValuesForBackupFormat() []string {
	return []string{
		string(BackupFormatCollatedFormat),
		string(BackupFormatRaw),
	}
}

func (s *BackupFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBackupFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBackupFormat(input string) (*BackupFormat, error) {
	vals := map[string]BackupFormat{
		"collatedformat": BackupFormatCollatedFormat,
		"raw":            BackupFormatRaw,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackupFormat(input)
	return &out, nil
}

type OperationStatus string

const (
	OperationStatusCancelInProgress OperationStatus = "CancelInProgress"
	OperationStatusCanceled         OperationStatus = "Canceled"
	OperationStatusFailed           OperationStatus = "Failed"
	OperationStatusInProgress       OperationStatus = "InProgress"
	OperationStatusPending          OperationStatus = "Pending"
	OperationStatusSucceeded        OperationStatus = "Succeeded"
)

func PossibleValuesForOperationStatus() []string {
	return []string{
		string(OperationStatusCancelInProgress),
		string(OperationStatusCanceled),
		string(OperationStatusFailed),
		string(OperationStatusInProgress),
		string(OperationStatusPending),
		string(OperationStatusSucceeded),
	}
}

func (s *OperationStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperationStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperationStatus(input string) (*OperationStatus, error) {
	vals := map[string]OperationStatus{
		"cancelinprogress": OperationStatusCancelInProgress,
		"canceled":         OperationStatusCanceled,
		"failed":           OperationStatusFailed,
		"inprogress":       OperationStatusInProgress,
		"pending":          OperationStatusPending,
		"succeeded":        OperationStatusSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperationStatus(input)
	return &out, nil
}
