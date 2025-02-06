package capabilities

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapabilityTypeProperties struct {
	AzureRbacActions     *[]string                                  `json:"azureRbacActions,omitempty"`
	AzureRbacDataActions *[]string                                  `json:"azureRbacDataActions,omitempty"`
	Description          *string                                    `json:"description,omitempty"`
	DisplayName          *string                                    `json:"displayName,omitempty"`
	Kind                 *string                                    `json:"kind,omitempty"`
	ParametersSchema     *string                                    `json:"parametersSchema,omitempty"`
	Publisher            *string                                    `json:"publisher,omitempty"`
	RuntimeProperties    *CapabilityTypePropertiesRuntimeProperties `json:"runtimeProperties,omitempty"`
	TargetType           *string                                    `json:"targetType,omitempty"`
	Urn                  *string                                    `json:"urn,omitempty"`
}
