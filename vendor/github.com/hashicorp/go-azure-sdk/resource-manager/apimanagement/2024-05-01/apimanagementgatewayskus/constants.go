package apimanagementgatewayskus

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

type GatewaySkuCapacityScaleType string

const (
	GatewaySkuCapacityScaleTypeAutomatic GatewaySkuCapacityScaleType = "Automatic"
	GatewaySkuCapacityScaleTypeManual    GatewaySkuCapacityScaleType = "Manual"
	GatewaySkuCapacityScaleTypeNone      GatewaySkuCapacityScaleType = "None"
)

func PossibleValuesForGatewaySkuCapacityScaleType() []string {
	return []string{
		string(GatewaySkuCapacityScaleTypeAutomatic),
		string(GatewaySkuCapacityScaleTypeManual),
		string(GatewaySkuCapacityScaleTypeNone),
	}
}

func (s *GatewaySkuCapacityScaleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGatewaySkuCapacityScaleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGatewaySkuCapacityScaleType(input string) (*GatewaySkuCapacityScaleType, error) {
	vals := map[string]GatewaySkuCapacityScaleType{
		"automatic": GatewaySkuCapacityScaleTypeAutomatic,
		"manual":    GatewaySkuCapacityScaleTypeManual,
		"none":      GatewaySkuCapacityScaleTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GatewaySkuCapacityScaleType(input)
	return &out, nil
}
