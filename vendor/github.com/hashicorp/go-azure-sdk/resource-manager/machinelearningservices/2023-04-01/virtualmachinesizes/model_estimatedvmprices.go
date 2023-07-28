package virtualmachinesizes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EstimatedVMPrices struct {
	BillingCurrency BillingCurrency    `json:"billingCurrency"`
	UnitOfMeasure   UnitOfMeasure      `json:"unitOfMeasure"`
	Values          []EstimatedVMPrice `json:"values"`
}
