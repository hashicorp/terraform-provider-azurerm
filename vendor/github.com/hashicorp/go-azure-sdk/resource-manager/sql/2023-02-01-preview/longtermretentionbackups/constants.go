package longtermretentionbackups

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupStorageRedundancy string

const (
	BackupStorageRedundancyGeo     BackupStorageRedundancy = "Geo"
	BackupStorageRedundancyGeoZone BackupStorageRedundancy = "GeoZone"
	BackupStorageRedundancyLocal   BackupStorageRedundancy = "Local"
	BackupStorageRedundancyZone    BackupStorageRedundancy = "Zone"
)

func PossibleValuesForBackupStorageRedundancy() []string {
	return []string{
		string(BackupStorageRedundancyGeo),
		string(BackupStorageRedundancyGeoZone),
		string(BackupStorageRedundancyLocal),
		string(BackupStorageRedundancyZone),
	}
}

func (s *BackupStorageRedundancy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBackupStorageRedundancy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBackupStorageRedundancy(input string) (*BackupStorageRedundancy, error) {
	vals := map[string]BackupStorageRedundancy{
		"geo":     BackupStorageRedundancyGeo,
		"geozone": BackupStorageRedundancyGeoZone,
		"local":   BackupStorageRedundancyLocal,
		"zone":    BackupStorageRedundancyZone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackupStorageRedundancy(input)
	return &out, nil
}

type DatabaseState string

const (
	DatabaseStateAll     DatabaseState = "All"
	DatabaseStateDeleted DatabaseState = "Deleted"
	DatabaseStateLive    DatabaseState = "Live"
)

func PossibleValuesForDatabaseState() []string {
	return []string{
		string(DatabaseStateAll),
		string(DatabaseStateDeleted),
		string(DatabaseStateLive),
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
		"all":     DatabaseStateAll,
		"deleted": DatabaseStateDeleted,
		"live":    DatabaseStateLive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatabaseState(input)
	return &out, nil
}
