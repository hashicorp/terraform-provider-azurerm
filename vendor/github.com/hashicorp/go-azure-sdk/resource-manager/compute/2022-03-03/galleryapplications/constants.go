package galleryapplications

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryApplicationCustomActionParameterType string

const (
	GalleryApplicationCustomActionParameterTypeConfigurationDataBlob GalleryApplicationCustomActionParameterType = "ConfigurationDataBlob"
	GalleryApplicationCustomActionParameterTypeLogOutputBlob         GalleryApplicationCustomActionParameterType = "LogOutputBlob"
	GalleryApplicationCustomActionParameterTypeString                GalleryApplicationCustomActionParameterType = "String"
)

func PossibleValuesForGalleryApplicationCustomActionParameterType() []string {
	return []string{
		string(GalleryApplicationCustomActionParameterTypeConfigurationDataBlob),
		string(GalleryApplicationCustomActionParameterTypeLogOutputBlob),
		string(GalleryApplicationCustomActionParameterTypeString),
	}
}

func parseGalleryApplicationCustomActionParameterType(input string) (*GalleryApplicationCustomActionParameterType, error) {
	vals := map[string]GalleryApplicationCustomActionParameterType{
		"configurationdatablob": GalleryApplicationCustomActionParameterTypeConfigurationDataBlob,
		"logoutputblob":         GalleryApplicationCustomActionParameterTypeLogOutputBlob,
		"string":                GalleryApplicationCustomActionParameterTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GalleryApplicationCustomActionParameterType(input)
	return &out, nil
}

type OperatingSystemTypes string

const (
	OperatingSystemTypesLinux   OperatingSystemTypes = "Linux"
	OperatingSystemTypesWindows OperatingSystemTypes = "Windows"
)

func PossibleValuesForOperatingSystemTypes() []string {
	return []string{
		string(OperatingSystemTypesLinux),
		string(OperatingSystemTypesWindows),
	}
}

func parseOperatingSystemTypes(input string) (*OperatingSystemTypes, error) {
	vals := map[string]OperatingSystemTypes{
		"linux":   OperatingSystemTypesLinux,
		"windows": OperatingSystemTypesWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperatingSystemTypes(input)
	return &out, nil
}
