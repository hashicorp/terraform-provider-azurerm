package virtualmachines

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnableStatus string

const (
	EnableStatusDisabled EnableStatus = "Disabled"
	EnableStatusEnabled  EnableStatus = "Enabled"
)

func PossibleValuesForEnableStatus() []string {
	return []string{
		string(EnableStatusDisabled),
		string(EnableStatusEnabled),
	}
}

func parseEnableStatus(input string) (*EnableStatus, error) {
	vals := map[string]EnableStatus{
		"disabled": EnableStatusDisabled,
		"enabled":  EnableStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnableStatus(input)
	return &out, nil
}

type HostCachingOptions string

const (
	HostCachingOptionsNone      HostCachingOptions = "None"
	HostCachingOptionsReadOnly  HostCachingOptions = "ReadOnly"
	HostCachingOptionsReadWrite HostCachingOptions = "ReadWrite"
)

func PossibleValuesForHostCachingOptions() []string {
	return []string{
		string(HostCachingOptionsNone),
		string(HostCachingOptionsReadOnly),
		string(HostCachingOptionsReadWrite),
	}
}

func parseHostCachingOptions(input string) (*HostCachingOptions, error) {
	vals := map[string]HostCachingOptions{
		"none":      HostCachingOptionsNone,
		"readonly":  HostCachingOptionsReadOnly,
		"readwrite": HostCachingOptionsReadWrite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HostCachingOptions(input)
	return &out, nil
}

type StorageType string

const (
	StorageTypePremium     StorageType = "Premium"
	StorageTypeStandard    StorageType = "Standard"
	StorageTypeStandardSSD StorageType = "StandardSSD"
)

func PossibleValuesForStorageType() []string {
	return []string{
		string(StorageTypePremium),
		string(StorageTypeStandard),
		string(StorageTypeStandardSSD),
	}
}

func parseStorageType(input string) (*StorageType, error) {
	vals := map[string]StorageType{
		"premium":     StorageTypePremium,
		"standard":    StorageTypeStandard,
		"standardssd": StorageTypeStandardSSD,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageType(input)
	return &out, nil
}

type TransportProtocol string

const (
	TransportProtocolTcp TransportProtocol = "Tcp"
	TransportProtocolUdp TransportProtocol = "Udp"
)

func PossibleValuesForTransportProtocol() []string {
	return []string{
		string(TransportProtocolTcp),
		string(TransportProtocolUdp),
	}
}

func parseTransportProtocol(input string) (*TransportProtocol, error) {
	vals := map[string]TransportProtocol{
		"tcp": TransportProtocolTcp,
		"udp": TransportProtocolUdp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TransportProtocol(input)
	return &out, nil
}

type VirtualMachineCreationSource string

const (
	VirtualMachineCreationSourceFromCustomImage        VirtualMachineCreationSource = "FromCustomImage"
	VirtualMachineCreationSourceFromGalleryImage       VirtualMachineCreationSource = "FromGalleryImage"
	VirtualMachineCreationSourceFromSharedGalleryImage VirtualMachineCreationSource = "FromSharedGalleryImage"
)

func PossibleValuesForVirtualMachineCreationSource() []string {
	return []string{
		string(VirtualMachineCreationSourceFromCustomImage),
		string(VirtualMachineCreationSourceFromGalleryImage),
		string(VirtualMachineCreationSourceFromSharedGalleryImage),
	}
}

func parseVirtualMachineCreationSource(input string) (*VirtualMachineCreationSource, error) {
	vals := map[string]VirtualMachineCreationSource{
		"fromcustomimage":        VirtualMachineCreationSourceFromCustomImage,
		"fromgalleryimage":       VirtualMachineCreationSourceFromGalleryImage,
		"fromsharedgalleryimage": VirtualMachineCreationSourceFromSharedGalleryImage,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualMachineCreationSource(input)
	return &out, nil
}
