package graphservicesprods

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountResourceProperties struct {
	AppId             string             `json:"appId"`
	BillingPlanId     *string            `json:"billingPlanId,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
