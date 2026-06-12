package apigateway

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiGatewaySkuType string

const (
	ApiGatewaySkuTypeStandard                 ApiGatewaySkuType = "Standard"
	ApiGatewaySkuTypeWorkspaceGatewayPremium  ApiGatewaySkuType = "WorkspaceGatewayPremium"
	ApiGatewaySkuTypeWorkspaceGatewayStandard ApiGatewaySkuType = "WorkspaceGatewayStandard"
)

func PossibleValuesForApiGatewaySkuType() []string {
	return []string{
		string(ApiGatewaySkuTypeStandard),
		string(ApiGatewaySkuTypeWorkspaceGatewayPremium),
		string(ApiGatewaySkuTypeWorkspaceGatewayStandard),
	}
}

func (s *ApiGatewaySkuType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApiGatewaySkuType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApiGatewaySkuType(input string) (*ApiGatewaySkuType, error) {
	vals := map[string]ApiGatewaySkuType{
		"standard":                 ApiGatewaySkuTypeStandard,
		"workspacegatewaypremium":  ApiGatewaySkuTypeWorkspaceGatewayPremium,
		"workspacegatewaystandard": ApiGatewaySkuTypeWorkspaceGatewayStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApiGatewaySkuType(input)
	return &out, nil
}

type VirtualNetworkType string

const (
	VirtualNetworkTypeExternal VirtualNetworkType = "External"
	VirtualNetworkTypeInternal VirtualNetworkType = "Internal"
	VirtualNetworkTypeNone     VirtualNetworkType = "None"
)

func PossibleValuesForVirtualNetworkType() []string {
	return []string{
		string(VirtualNetworkTypeExternal),
		string(VirtualNetworkTypeInternal),
		string(VirtualNetworkTypeNone),
	}
}

func (s *VirtualNetworkType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkType(input string) (*VirtualNetworkType, error) {
	vals := map[string]VirtualNetworkType{
		"external": VirtualNetworkTypeExternal,
		"internal": VirtualNetworkTypeInternal,
		"none":     VirtualNetworkTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkType(input)
	return &out, nil
}
