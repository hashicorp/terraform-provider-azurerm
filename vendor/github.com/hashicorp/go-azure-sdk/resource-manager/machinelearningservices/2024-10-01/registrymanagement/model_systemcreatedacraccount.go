package registrymanagement

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SystemCreatedAcrAccount struct {
	AcrAccountName *string        `json:"acrAccountName,omitempty"`
	AcrAccountSku  *string        `json:"acrAccountSku,omitempty"`
	ArmResourceId  *ArmResourceId `json:"armResourceId,omitempty"`
}
