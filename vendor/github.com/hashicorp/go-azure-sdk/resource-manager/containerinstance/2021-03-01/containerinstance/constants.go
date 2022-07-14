package containerinstance

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerGroupIpAddressType string

const (
	ContainerGroupIpAddressTypePrivate ContainerGroupIpAddressType = "Private"
	ContainerGroupIpAddressTypePublic  ContainerGroupIpAddressType = "Public"
)

func PossibleValuesForContainerGroupIpAddressType() []string {
	return []string{
		string(ContainerGroupIpAddressTypePrivate),
		string(ContainerGroupIpAddressTypePublic),
	}
}

func parseContainerGroupIpAddressType(input string) (*ContainerGroupIpAddressType, error) {
	vals := map[string]ContainerGroupIpAddressType{
		"private": ContainerGroupIpAddressTypePrivate,
		"public":  ContainerGroupIpAddressTypePublic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerGroupIpAddressType(input)
	return &out, nil
}

type ContainerGroupNetworkProtocol string

const (
	ContainerGroupNetworkProtocolTCP ContainerGroupNetworkProtocol = "TCP"
	ContainerGroupNetworkProtocolUDP ContainerGroupNetworkProtocol = "UDP"
)

func PossibleValuesForContainerGroupNetworkProtocol() []string {
	return []string{
		string(ContainerGroupNetworkProtocolTCP),
		string(ContainerGroupNetworkProtocolUDP),
	}
}

func parseContainerGroupNetworkProtocol(input string) (*ContainerGroupNetworkProtocol, error) {
	vals := map[string]ContainerGroupNetworkProtocol{
		"tcp": ContainerGroupNetworkProtocolTCP,
		"udp": ContainerGroupNetworkProtocolUDP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerGroupNetworkProtocol(input)
	return &out, nil
}

type ContainerGroupRestartPolicy string

const (
	ContainerGroupRestartPolicyAlways    ContainerGroupRestartPolicy = "Always"
	ContainerGroupRestartPolicyNever     ContainerGroupRestartPolicy = "Never"
	ContainerGroupRestartPolicyOnFailure ContainerGroupRestartPolicy = "OnFailure"
)

func PossibleValuesForContainerGroupRestartPolicy() []string {
	return []string{
		string(ContainerGroupRestartPolicyAlways),
		string(ContainerGroupRestartPolicyNever),
		string(ContainerGroupRestartPolicyOnFailure),
	}
}

func parseContainerGroupRestartPolicy(input string) (*ContainerGroupRestartPolicy, error) {
	vals := map[string]ContainerGroupRestartPolicy{
		"always":    ContainerGroupRestartPolicyAlways,
		"never":     ContainerGroupRestartPolicyNever,
		"onfailure": ContainerGroupRestartPolicyOnFailure,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerGroupRestartPolicy(input)
	return &out, nil
}

type ContainerGroupSku string

const (
	ContainerGroupSkuDedicated ContainerGroupSku = "Dedicated"
	ContainerGroupSkuStandard  ContainerGroupSku = "Standard"
)

func PossibleValuesForContainerGroupSku() []string {
	return []string{
		string(ContainerGroupSkuDedicated),
		string(ContainerGroupSkuStandard),
	}
}

func parseContainerGroupSku(input string) (*ContainerGroupSku, error) {
	vals := map[string]ContainerGroupSku{
		"dedicated": ContainerGroupSkuDedicated,
		"standard":  ContainerGroupSkuStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerGroupSku(input)
	return &out, nil
}

type ContainerNetworkProtocol string

const (
	ContainerNetworkProtocolTCP ContainerNetworkProtocol = "TCP"
	ContainerNetworkProtocolUDP ContainerNetworkProtocol = "UDP"
)

func PossibleValuesForContainerNetworkProtocol() []string {
	return []string{
		string(ContainerNetworkProtocolTCP),
		string(ContainerNetworkProtocolUDP),
	}
}

func parseContainerNetworkProtocol(input string) (*ContainerNetworkProtocol, error) {
	vals := map[string]ContainerNetworkProtocol{
		"tcp": ContainerNetworkProtocolTCP,
		"udp": ContainerNetworkProtocolUDP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerNetworkProtocol(input)
	return &out, nil
}

type GpuSku string

const (
	GpuSkuKEightZero  GpuSku = "K80"
	GpuSkuPOneHundred GpuSku = "P100"
	GpuSkuVOneHundred GpuSku = "V100"
)

func PossibleValuesForGpuSku() []string {
	return []string{
		string(GpuSkuKEightZero),
		string(GpuSkuPOneHundred),
		string(GpuSkuVOneHundred),
	}
}

func parseGpuSku(input string) (*GpuSku, error) {
	vals := map[string]GpuSku{
		"k80":  GpuSkuKEightZero,
		"p100": GpuSkuPOneHundred,
		"v100": GpuSkuVOneHundred,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GpuSku(input)
	return &out, nil
}

type LogAnalyticsLogType string

const (
	LogAnalyticsLogTypeContainerInsights     LogAnalyticsLogType = "ContainerInsights"
	LogAnalyticsLogTypeContainerInstanceLogs LogAnalyticsLogType = "ContainerInstanceLogs"
)

func PossibleValuesForLogAnalyticsLogType() []string {
	return []string{
		string(LogAnalyticsLogTypeContainerInsights),
		string(LogAnalyticsLogTypeContainerInstanceLogs),
	}
}

func parseLogAnalyticsLogType(input string) (*LogAnalyticsLogType, error) {
	vals := map[string]LogAnalyticsLogType{
		"containerinsights":     LogAnalyticsLogTypeContainerInsights,
		"containerinstancelogs": LogAnalyticsLogTypeContainerInstanceLogs,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LogAnalyticsLogType(input)
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

type Scheme string

const (
	SchemeHttp  Scheme = "http"
	SchemeHttps Scheme = "https"
)

func PossibleValuesForScheme() []string {
	return []string{
		string(SchemeHttp),
		string(SchemeHttps),
	}
}

func parseScheme(input string) (*Scheme, error) {
	vals := map[string]Scheme{
		"http":  SchemeHttp,
		"https": SchemeHttps,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Scheme(input)
	return &out, nil
}
