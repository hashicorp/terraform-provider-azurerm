package service

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServicePropertiesFormat struct {
	PccRules          []PccRuleConfiguration `json:"pccRules"`
	ProvisioningState *ProvisioningState     `json:"provisioningState,omitempty"`
	ServicePrecedence int64                  `json:"servicePrecedence"`
	ServiceQosPolicy  *QosPolicy             `json:"serviceQosPolicy,omitempty"`
}
