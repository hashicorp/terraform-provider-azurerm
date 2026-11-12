package databaseeditions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DbSystemDatabaseEditionType string

const (
	DbSystemDatabaseEditionTypeEnterpriseEdition                DbSystemDatabaseEditionType = "EnterpriseEdition"
	DbSystemDatabaseEditionTypeEnterpriseEditionDeveloper       DbSystemDatabaseEditionType = "EnterpriseEditionDeveloper"
	DbSystemDatabaseEditionTypeEnterpriseEditionExtreme         DbSystemDatabaseEditionType = "EnterpriseEditionExtreme"
	DbSystemDatabaseEditionTypeEnterpriseEditionHighPerformance DbSystemDatabaseEditionType = "EnterpriseEditionHighPerformance"
	DbSystemDatabaseEditionTypeStandardEdition                  DbSystemDatabaseEditionType = "StandardEdition"
)

func PossibleValuesForDbSystemDatabaseEditionType() []string {
	return []string{
		string(DbSystemDatabaseEditionTypeEnterpriseEdition),
		string(DbSystemDatabaseEditionTypeEnterpriseEditionDeveloper),
		string(DbSystemDatabaseEditionTypeEnterpriseEditionExtreme),
		string(DbSystemDatabaseEditionTypeEnterpriseEditionHighPerformance),
		string(DbSystemDatabaseEditionTypeStandardEdition),
	}
}

func (s *DbSystemDatabaseEditionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDbSystemDatabaseEditionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDbSystemDatabaseEditionType(input string) (*DbSystemDatabaseEditionType, error) {
	vals := map[string]DbSystemDatabaseEditionType{
		"enterpriseedition":                DbSystemDatabaseEditionTypeEnterpriseEdition,
		"enterpriseeditiondeveloper":       DbSystemDatabaseEditionTypeEnterpriseEditionDeveloper,
		"enterpriseeditionextreme":         DbSystemDatabaseEditionTypeEnterpriseEditionExtreme,
		"enterpriseeditionhighperformance": DbSystemDatabaseEditionTypeEnterpriseEditionHighPerformance,
		"standardedition":                  DbSystemDatabaseEditionTypeStandardEdition,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DbSystemDatabaseEditionType(input)
	return &out, nil
}
