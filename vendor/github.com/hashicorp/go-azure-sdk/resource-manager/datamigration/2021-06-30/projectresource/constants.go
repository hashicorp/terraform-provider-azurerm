package projectresource

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

type ProjectProvisioningState string

const (
	ProjectProvisioningStateDeleting  ProjectProvisioningState = "Deleting"
	ProjectProvisioningStateSucceeded ProjectProvisioningState = "Succeeded"
)

func PossibleValuesForProjectProvisioningState() []string {
	return []string{
		string(ProjectProvisioningStateDeleting),
		string(ProjectProvisioningStateSucceeded),
	}
}

func (s *ProjectProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProjectProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProjectProvisioningState(input string) (*ProjectProvisioningState, error) {
	vals := map[string]ProjectProvisioningState{
		"deleting":  ProjectProvisioningStateDeleting,
		"succeeded": ProjectProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProjectProvisioningState(input)
	return &out, nil
}

type ProjectSourcePlatform string

const (
	ProjectSourcePlatformMongoDb    ProjectSourcePlatform = "MongoDb"
	ProjectSourcePlatformMySQL      ProjectSourcePlatform = "MySQL"
	ProjectSourcePlatformPostgreSql ProjectSourcePlatform = "PostgreSql"
	ProjectSourcePlatformSQL        ProjectSourcePlatform = "SQL"
	ProjectSourcePlatformUnknown    ProjectSourcePlatform = "Unknown"
)

func PossibleValuesForProjectSourcePlatform() []string {
	return []string{
		string(ProjectSourcePlatformMongoDb),
		string(ProjectSourcePlatformMySQL),
		string(ProjectSourcePlatformPostgreSql),
		string(ProjectSourcePlatformSQL),
		string(ProjectSourcePlatformUnknown),
	}
}

func (s *ProjectSourcePlatform) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProjectSourcePlatform(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProjectSourcePlatform(input string) (*ProjectSourcePlatform, error) {
	vals := map[string]ProjectSourcePlatform{
		"mongodb":    ProjectSourcePlatformMongoDb,
		"mysql":      ProjectSourcePlatformMySQL,
		"postgresql": ProjectSourcePlatformPostgreSql,
		"sql":        ProjectSourcePlatformSQL,
		"unknown":    ProjectSourcePlatformUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProjectSourcePlatform(input)
	return &out, nil
}

type ProjectTargetPlatform string

const (
	ProjectTargetPlatformAzureDbForMySql      ProjectTargetPlatform = "AzureDbForMySql"
	ProjectTargetPlatformAzureDbForPostgreSql ProjectTargetPlatform = "AzureDbForPostgreSql"
	ProjectTargetPlatformMongoDb              ProjectTargetPlatform = "MongoDb"
	ProjectTargetPlatformSQLDB                ProjectTargetPlatform = "SQLDB"
	ProjectTargetPlatformSQLMI                ProjectTargetPlatform = "SQLMI"
	ProjectTargetPlatformUnknown              ProjectTargetPlatform = "Unknown"
)

func PossibleValuesForProjectTargetPlatform() []string {
	return []string{
		string(ProjectTargetPlatformAzureDbForMySql),
		string(ProjectTargetPlatformAzureDbForPostgreSql),
		string(ProjectTargetPlatformMongoDb),
		string(ProjectTargetPlatformSQLDB),
		string(ProjectTargetPlatformSQLMI),
		string(ProjectTargetPlatformUnknown),
	}
}

func (s *ProjectTargetPlatform) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProjectTargetPlatform(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProjectTargetPlatform(input string) (*ProjectTargetPlatform, error) {
	vals := map[string]ProjectTargetPlatform{
		"azuredbformysql":      ProjectTargetPlatformAzureDbForMySql,
		"azuredbforpostgresql": ProjectTargetPlatformAzureDbForPostgreSql,
		"mongodb":              ProjectTargetPlatformMongoDb,
		"sqldb":                ProjectTargetPlatformSQLDB,
		"sqlmi":                ProjectTargetPlatformSQLMI,
		"unknown":              ProjectTargetPlatformUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProjectTargetPlatform(input)
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
