package machines

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentConfigurationMode string

const (
	AgentConfigurationModeFull    AgentConfigurationMode = "full"
	AgentConfigurationModeMonitor AgentConfigurationMode = "monitor"
)

func PossibleValuesForAgentConfigurationMode() []string {
	return []string{
		string(AgentConfigurationModeFull),
		string(AgentConfigurationModeMonitor),
	}
}

func parseAgentConfigurationMode(input string) (*AgentConfigurationMode, error) {
	vals := map[string]AgentConfigurationMode{
		"full":    AgentConfigurationModeFull,
		"monitor": AgentConfigurationModeMonitor,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AgentConfigurationMode(input)
	return &out, nil
}

type AssessmentModeTypes string

const (
	AssessmentModeTypesAutomaticByPlatform AssessmentModeTypes = "AutomaticByPlatform"
	AssessmentModeTypesImageDefault        AssessmentModeTypes = "ImageDefault"
)

func PossibleValuesForAssessmentModeTypes() []string {
	return []string{
		string(AssessmentModeTypesAutomaticByPlatform),
		string(AssessmentModeTypesImageDefault),
	}
}

func parseAssessmentModeTypes(input string) (*AssessmentModeTypes, error) {
	vals := map[string]AssessmentModeTypes{
		"automaticbyplatform": AssessmentModeTypesAutomaticByPlatform,
		"imagedefault":        AssessmentModeTypesImageDefault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AssessmentModeTypes(input)
	return &out, nil
}

type InstanceViewTypes string

const (
	InstanceViewTypesInstanceView InstanceViewTypes = "instanceView"
)

func PossibleValuesForInstanceViewTypes() []string {
	return []string{
		string(InstanceViewTypesInstanceView),
	}
}

func parseInstanceViewTypes(input string) (*InstanceViewTypes, error) {
	vals := map[string]InstanceViewTypes{
		"instanceview": InstanceViewTypesInstanceView,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InstanceViewTypes(input)
	return &out, nil
}

type PatchModeTypes string

const (
	PatchModeTypesAutomaticByOS       PatchModeTypes = "AutomaticByOS"
	PatchModeTypesAutomaticByPlatform PatchModeTypes = "AutomaticByPlatform"
	PatchModeTypesImageDefault        PatchModeTypes = "ImageDefault"
	PatchModeTypesManual              PatchModeTypes = "Manual"
)

func PossibleValuesForPatchModeTypes() []string {
	return []string{
		string(PatchModeTypesAutomaticByOS),
		string(PatchModeTypesAutomaticByPlatform),
		string(PatchModeTypesImageDefault),
		string(PatchModeTypesManual),
	}
}

func parsePatchModeTypes(input string) (*PatchModeTypes, error) {
	vals := map[string]PatchModeTypes{
		"automaticbyos":       PatchModeTypesAutomaticByOS,
		"automaticbyplatform": PatchModeTypesAutomaticByPlatform,
		"imagedefault":        PatchModeTypesImageDefault,
		"manual":              PatchModeTypesManual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PatchModeTypes(input)
	return &out, nil
}

type StatusLevelTypes string

const (
	StatusLevelTypesError   StatusLevelTypes = "Error"
	StatusLevelTypesInfo    StatusLevelTypes = "Info"
	StatusLevelTypesWarning StatusLevelTypes = "Warning"
)

func PossibleValuesForStatusLevelTypes() []string {
	return []string{
		string(StatusLevelTypesError),
		string(StatusLevelTypesInfo),
		string(StatusLevelTypesWarning),
	}
}

func parseStatusLevelTypes(input string) (*StatusLevelTypes, error) {
	vals := map[string]StatusLevelTypes{
		"error":   StatusLevelTypesError,
		"info":    StatusLevelTypesInfo,
		"warning": StatusLevelTypesWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StatusLevelTypes(input)
	return &out, nil
}

type StatusTypes string

const (
	StatusTypesConnected    StatusTypes = "Connected"
	StatusTypesDisconnected StatusTypes = "Disconnected"
	StatusTypesError        StatusTypes = "Error"
)

func PossibleValuesForStatusTypes() []string {
	return []string{
		string(StatusTypesConnected),
		string(StatusTypesDisconnected),
		string(StatusTypesError),
	}
}

func parseStatusTypes(input string) (*StatusTypes, error) {
	vals := map[string]StatusTypes{
		"connected":    StatusTypesConnected,
		"disconnected": StatusTypesDisconnected,
		"error":        StatusTypesError,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StatusTypes(input)
	return &out, nil
}
