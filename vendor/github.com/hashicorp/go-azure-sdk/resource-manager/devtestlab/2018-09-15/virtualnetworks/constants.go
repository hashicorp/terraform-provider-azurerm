package virtualnetworks

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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

type UsagePermissionType string

const (
	UsagePermissionTypeAllow   UsagePermissionType = "Allow"
	UsagePermissionTypeDefault UsagePermissionType = "Default"
	UsagePermissionTypeDeny    UsagePermissionType = "Deny"
)

func PossibleValuesForUsagePermissionType() []string {
	return []string{
		string(UsagePermissionTypeAllow),
		string(UsagePermissionTypeDefault),
		string(UsagePermissionTypeDeny),
	}
}

func parseUsagePermissionType(input string) (*UsagePermissionType, error) {
	vals := map[string]UsagePermissionType{
		"allow":   UsagePermissionTypeAllow,
		"default": UsagePermissionTypeDefault,
		"deny":    UsagePermissionTypeDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UsagePermissionType(input)
	return &out, nil
}
