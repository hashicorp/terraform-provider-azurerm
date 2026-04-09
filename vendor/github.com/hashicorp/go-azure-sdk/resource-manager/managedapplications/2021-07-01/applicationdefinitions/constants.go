package applicationdefinitions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationArtifactType string

const (
	ApplicationArtifactTypeCustom       ApplicationArtifactType = "Custom"
	ApplicationArtifactTypeNotSpecified ApplicationArtifactType = "NotSpecified"
	ApplicationArtifactTypeTemplate     ApplicationArtifactType = "Template"
)

func PossibleValuesForApplicationArtifactType() []string {
	return []string{
		string(ApplicationArtifactTypeCustom),
		string(ApplicationArtifactTypeNotSpecified),
		string(ApplicationArtifactTypeTemplate),
	}
}

func (s *ApplicationArtifactType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationArtifactType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationArtifactType(input string) (*ApplicationArtifactType, error) {
	vals := map[string]ApplicationArtifactType{
		"custom":       ApplicationArtifactTypeCustom,
		"notspecified": ApplicationArtifactTypeNotSpecified,
		"template":     ApplicationArtifactTypeTemplate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationArtifactType(input)
	return &out, nil
}

type ApplicationDefinitionArtifactName string

const (
	ApplicationDefinitionArtifactNameApplicationResourceTemplate ApplicationDefinitionArtifactName = "ApplicationResourceTemplate"
	ApplicationDefinitionArtifactNameCreateUiDefinition          ApplicationDefinitionArtifactName = "CreateUiDefinition"
	ApplicationDefinitionArtifactNameMainTemplateParameters      ApplicationDefinitionArtifactName = "MainTemplateParameters"
	ApplicationDefinitionArtifactNameNotSpecified                ApplicationDefinitionArtifactName = "NotSpecified"
)

func PossibleValuesForApplicationDefinitionArtifactName() []string {
	return []string{
		string(ApplicationDefinitionArtifactNameApplicationResourceTemplate),
		string(ApplicationDefinitionArtifactNameCreateUiDefinition),
		string(ApplicationDefinitionArtifactNameMainTemplateParameters),
		string(ApplicationDefinitionArtifactNameNotSpecified),
	}
}

func (s *ApplicationDefinitionArtifactName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationDefinitionArtifactName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationDefinitionArtifactName(input string) (*ApplicationDefinitionArtifactName, error) {
	vals := map[string]ApplicationDefinitionArtifactName{
		"applicationresourcetemplate": ApplicationDefinitionArtifactNameApplicationResourceTemplate,
		"createuidefinition":          ApplicationDefinitionArtifactNameCreateUiDefinition,
		"maintemplateparameters":      ApplicationDefinitionArtifactNameMainTemplateParameters,
		"notspecified":                ApplicationDefinitionArtifactNameNotSpecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationDefinitionArtifactName(input)
	return &out, nil
}

type ApplicationLockLevel string

const (
	ApplicationLockLevelCanNotDelete ApplicationLockLevel = "CanNotDelete"
	ApplicationLockLevelNone         ApplicationLockLevel = "None"
	ApplicationLockLevelReadOnly     ApplicationLockLevel = "ReadOnly"
)

func PossibleValuesForApplicationLockLevel() []string {
	return []string{
		string(ApplicationLockLevelCanNotDelete),
		string(ApplicationLockLevelNone),
		string(ApplicationLockLevelReadOnly),
	}
}

func (s *ApplicationLockLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationLockLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationLockLevel(input string) (*ApplicationLockLevel, error) {
	vals := map[string]ApplicationLockLevel{
		"cannotdelete": ApplicationLockLevelCanNotDelete,
		"none":         ApplicationLockLevelNone,
		"readonly":     ApplicationLockLevelReadOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationLockLevel(input)
	return &out, nil
}

type ApplicationManagementMode string

const (
	ApplicationManagementModeManaged      ApplicationManagementMode = "Managed"
	ApplicationManagementModeNotSpecified ApplicationManagementMode = "NotSpecified"
	ApplicationManagementModeUnmanaged    ApplicationManagementMode = "Unmanaged"
)

func PossibleValuesForApplicationManagementMode() []string {
	return []string{
		string(ApplicationManagementModeManaged),
		string(ApplicationManagementModeNotSpecified),
		string(ApplicationManagementModeUnmanaged),
	}
}

func (s *ApplicationManagementMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationManagementMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationManagementMode(input string) (*ApplicationManagementMode, error) {
	vals := map[string]ApplicationManagementMode{
		"managed":      ApplicationManagementModeManaged,
		"notspecified": ApplicationManagementModeNotSpecified,
		"unmanaged":    ApplicationManagementModeUnmanaged,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationManagementMode(input)
	return &out, nil
}

type DeploymentMode string

const (
	DeploymentModeComplete     DeploymentMode = "Complete"
	DeploymentModeIncremental  DeploymentMode = "Incremental"
	DeploymentModeNotSpecified DeploymentMode = "NotSpecified"
)

func PossibleValuesForDeploymentMode() []string {
	return []string{
		string(DeploymentModeComplete),
		string(DeploymentModeIncremental),
		string(DeploymentModeNotSpecified),
	}
}

func (s *DeploymentMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeploymentMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeploymentMode(input string) (*DeploymentMode, error) {
	vals := map[string]DeploymentMode{
		"complete":     DeploymentModeComplete,
		"incremental":  DeploymentModeIncremental,
		"notspecified": DeploymentModeNotSpecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentMode(input)
	return &out, nil
}
