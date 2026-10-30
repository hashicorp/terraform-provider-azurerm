package oraclesubscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OracleSubscriptionProperties struct {
	AddSubscriptionOperationState *AddSubscriptionOperationState       `json:"addSubscriptionOperationState,omitempty"`
	AzureSubscriptionIds          *[]string                            `json:"azureSubscriptionIds,omitempty"`
	CloudAccountId                *string                              `json:"cloudAccountId,omitempty"`
	CloudAccountState             *CloudAccountProvisioningState       `json:"cloudAccountState,omitempty"`
	Intent                        *Intent                              `json:"intent,omitempty"`
	LastOperationStatusDetail     *string                              `json:"lastOperationStatusDetail,omitempty"`
	ProductCode                   *string                              `json:"productCode,omitempty"`
	ProvisioningState             *OracleSubscriptionProvisioningState `json:"provisioningState,omitempty"`
	SaasSubscriptionId            *string                              `json:"saasSubscriptionId,omitempty"`
	TermUnit                      *string                              `json:"termUnit,omitempty"`
}
