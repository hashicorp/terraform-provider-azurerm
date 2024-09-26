package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayProperties struct {
	AddonConfigs            *map[string]interface{}                `json:"addonConfigs,omitempty"`
	ApiMetadataProperties   *GatewayApiMetadataProperties          `json:"apiMetadataProperties,omitempty"`
	ApmTypes                *[]ApmType                             `json:"apmTypes,omitempty"`
	Apms                    *[]ApmReference                        `json:"apms,omitempty"`
	ClientAuth              *GatewayPropertiesClientAuth           `json:"clientAuth,omitempty"`
	CorsProperties          *GatewayCorsProperties                 `json:"corsProperties,omitempty"`
	EnvironmentVariables    *GatewayPropertiesEnvironmentVariables `json:"environmentVariables,omitempty"`
	HTTPSOnly               *bool                                  `json:"httpsOnly,omitempty"`
	Instances               *[]GatewayInstance                     `json:"instances,omitempty"`
	OperatorProperties      *GatewayOperatorProperties             `json:"operatorProperties,omitempty"`
	ProvisioningState       *GatewayProvisioningState              `json:"provisioningState,omitempty"`
	Public                  *bool                                  `json:"public,omitempty"`
	ResourceRequests        *GatewayResourceRequests               `json:"resourceRequests,omitempty"`
	ResponseCacheProperties GatewayResponseCacheProperties         `json:"responseCacheProperties"`
	SsoProperties           *SsoProperties                         `json:"ssoProperties,omitempty"`
	Url                     *string                                `json:"url,omitempty"`
}

var _ json.Unmarshaler = &GatewayProperties{}

func (s *GatewayProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AddonConfigs          *map[string]interface{}                `json:"addonConfigs,omitempty"`
		ApiMetadataProperties *GatewayApiMetadataProperties          `json:"apiMetadataProperties,omitempty"`
		ApmTypes              *[]ApmType                             `json:"apmTypes,omitempty"`
		Apms                  *[]ApmReference                        `json:"apms,omitempty"`
		ClientAuth            *GatewayPropertiesClientAuth           `json:"clientAuth,omitempty"`
		CorsProperties        *GatewayCorsProperties                 `json:"corsProperties,omitempty"`
		EnvironmentVariables  *GatewayPropertiesEnvironmentVariables `json:"environmentVariables,omitempty"`
		HTTPSOnly             *bool                                  `json:"httpsOnly,omitempty"`
		Instances             *[]GatewayInstance                     `json:"instances,omitempty"`
		OperatorProperties    *GatewayOperatorProperties             `json:"operatorProperties,omitempty"`
		ProvisioningState     *GatewayProvisioningState              `json:"provisioningState,omitempty"`
		Public                *bool                                  `json:"public,omitempty"`
		ResourceRequests      *GatewayResourceRequests               `json:"resourceRequests,omitempty"`
		SsoProperties         *SsoProperties                         `json:"ssoProperties,omitempty"`
		Url                   *string                                `json:"url,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AddonConfigs = decoded.AddonConfigs
	s.ApiMetadataProperties = decoded.ApiMetadataProperties
	s.ApmTypes = decoded.ApmTypes
	s.Apms = decoded.Apms
	s.ClientAuth = decoded.ClientAuth
	s.CorsProperties = decoded.CorsProperties
	s.EnvironmentVariables = decoded.EnvironmentVariables
	s.HTTPSOnly = decoded.HTTPSOnly
	s.Instances = decoded.Instances
	s.OperatorProperties = decoded.OperatorProperties
	s.ProvisioningState = decoded.ProvisioningState
	s.Public = decoded.Public
	s.ResourceRequests = decoded.ResourceRequests
	s.SsoProperties = decoded.SsoProperties
	s.Url = decoded.Url

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling GatewayProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["responseCacheProperties"]; ok {
		impl, err := UnmarshalGatewayResponseCachePropertiesImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ResponseCacheProperties' for 'GatewayProperties': %+v", err)
		}
		s.ResponseCacheProperties = impl
	}

	return nil
}
