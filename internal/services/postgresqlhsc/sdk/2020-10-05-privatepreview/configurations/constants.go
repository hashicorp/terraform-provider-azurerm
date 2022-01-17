package configurations

import "strings"

type ConfigurationDataType string

const (
	ConfigurationDataTypeBoolean     ConfigurationDataType = "Boolean"
	ConfigurationDataTypeEnumeration ConfigurationDataType = "Enumeration"
	ConfigurationDataTypeInteger     ConfigurationDataType = "Integer"
	ConfigurationDataTypeNumeric     ConfigurationDataType = "Numeric"
)

func PossibleValuesForConfigurationDataType() []string {
	return []string{
		string(ConfigurationDataTypeBoolean),
		string(ConfigurationDataTypeEnumeration),
		string(ConfigurationDataTypeInteger),
		string(ConfigurationDataTypeNumeric),
	}
}

func parseConfigurationDataType(input string) (*ConfigurationDataType, error) {
	vals := map[string]ConfigurationDataType{
		"boolean":     ConfigurationDataTypeBoolean,
		"enumeration": ConfigurationDataTypeEnumeration,
		"integer":     ConfigurationDataTypeInteger,
		"numeric":     ConfigurationDataTypeNumeric,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConfigurationDataType(input)
	return &out, nil
}

type CreatedByType string

const (
	CreatedByTypeApplication     CreatedByType = "Application"
	CreatedByTypeKey             CreatedByType = "Key"
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	CreatedByTypeUser            CreatedByType = "User"
)

func PossibleValuesForCreatedByType() []string {
	return []string{
		string(CreatedByTypeApplication),
		string(CreatedByTypeKey),
		string(CreatedByTypeManagedIdentity),
		string(CreatedByTypeUser),
	}
}

func parseCreatedByType(input string) (*CreatedByType, error) {
	vals := map[string]CreatedByType{
		"application":     CreatedByTypeApplication,
		"key":             CreatedByTypeKey,
		"managedidentity": CreatedByTypeManagedIdentity,
		"user":            CreatedByTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CreatedByType(input)
	return &out, nil
}

type ServerRole string

const (
	ServerRoleCoordinator ServerRole = "Coordinator"
	ServerRoleWorker      ServerRole = "Worker"
)

func PossibleValuesForServerRole() []string {
	return []string{
		string(ServerRoleCoordinator),
		string(ServerRoleWorker),
	}
}

func parseServerRole(input string) (*ServerRole, error) {
	vals := map[string]ServerRole{
		"coordinator": ServerRoleCoordinator,
		"worker":      ServerRoleWorker,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerRole(input)
	return &out, nil
}
