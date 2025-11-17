package integrationruntimes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreateLinkedIntegrationRuntimeRequest struct {
	DataFactoryLocation *string `json:"dataFactoryLocation,omitempty"`
	DataFactoryName     *string `json:"dataFactoryName,omitempty"`
	Name                *string `json:"name,omitempty"`
	SubscriptionId      *string `json:"subscriptionId,omitempty"`
}
