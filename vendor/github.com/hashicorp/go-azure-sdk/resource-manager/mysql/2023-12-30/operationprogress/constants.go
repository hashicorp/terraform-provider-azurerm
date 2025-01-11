package operationprogress

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ObjectType string

const (
	ObjectTypeBackupAndExportResponse   ObjectType = "BackupAndExportResponse"
	ObjectTypeImportFromStorageResponse ObjectType = "ImportFromStorageResponse"
)

func PossibleValuesForObjectType() []string {
	return []string{
		string(ObjectTypeBackupAndExportResponse),
		string(ObjectTypeImportFromStorageResponse),
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
		"backupandexportresponse":   ObjectTypeBackupAndExportResponse,
		"importfromstorageresponse": ObjectTypeImportFromStorageResponse,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ObjectType(input)
	return &out, nil
}
