package mongorbacs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoRoleDefinitionType string

const (
	MongoRoleDefinitionTypeBuiltInRole MongoRoleDefinitionType = "BuiltInRole"
	MongoRoleDefinitionTypeCustomRole  MongoRoleDefinitionType = "CustomRole"
)

func PossibleValuesForMongoRoleDefinitionType() []string {
	return []string{
		string(MongoRoleDefinitionTypeBuiltInRole),
		string(MongoRoleDefinitionTypeCustomRole),
	}
}

func (s *MongoRoleDefinitionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMongoRoleDefinitionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMongoRoleDefinitionType(input string) (*MongoRoleDefinitionType, error) {
	vals := map[string]MongoRoleDefinitionType{
		"builtinrole": MongoRoleDefinitionTypeBuiltInRole,
		"customrole":  MongoRoleDefinitionTypeCustomRole,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MongoRoleDefinitionType(input)
	return &out, nil
}
