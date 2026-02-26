package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewaySslPredefinedPolicy struct {
	Id         *string                                                `json:"id,omitempty"`
	Name       *string                                                `json:"name,omitempty"`
	Properties *ApplicationGatewaySslPredefinedPolicyPropertiesFormat `json:"properties,omitempty"`
}
