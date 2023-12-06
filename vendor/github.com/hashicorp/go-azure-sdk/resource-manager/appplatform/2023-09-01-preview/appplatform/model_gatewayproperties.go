package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayProperties struct {
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
