package apimanagementgatewayskus

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayResourceSkuResultOperationPredicate struct {
	ResourceType *string
}

func (p GatewayResourceSkuResultOperationPredicate) Matches(input GatewayResourceSkuResult) bool {

	if p.ResourceType != nil && (input.ResourceType == nil || *p.ResourceType != *input.ResourceType) {
		return false
	}

	return true
}
