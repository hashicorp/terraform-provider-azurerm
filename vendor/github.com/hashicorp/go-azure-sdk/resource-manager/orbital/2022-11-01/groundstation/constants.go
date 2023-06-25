package groundstation

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapabilityParameter string

const (
	CapabilityParameterCommunication    CapabilityParameter = "Communication"
	CapabilityParameterEarthObservation CapabilityParameter = "EarthObservation"
)

func PossibleValuesForCapabilityParameter() []string {
	return []string{
		string(CapabilityParameterCommunication),
		string(CapabilityParameterEarthObservation),
	}
}

func parseCapabilityParameter(input string) (*CapabilityParameter, error) {
	vals := map[string]CapabilityParameter{
		"communication":    CapabilityParameterCommunication,
		"earthobservation": CapabilityParameterEarthObservation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CapabilityParameter(input)
	return &out, nil
}

type ReleaseMode string

const (
	ReleaseModeGA      ReleaseMode = "GA"
	ReleaseModePreview ReleaseMode = "Preview"
)

func PossibleValuesForReleaseMode() []string {
	return []string{
		string(ReleaseModeGA),
		string(ReleaseModePreview),
	}
}

func parseReleaseMode(input string) (*ReleaseMode, error) {
	vals := map[string]ReleaseMode{
		"ga":      ReleaseModeGA,
		"preview": ReleaseModePreview,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReleaseMode(input)
	return &out, nil
}
