package healthbots

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HealthBotProperties struct {
	BotManagementPortalLink *string             `json:"botManagementPortalLink,omitempty"`
	KeyVaultProperties      *KeyVaultProperties `json:"keyVaultProperties,omitempty"`
	ProvisioningState       *string             `json:"provisioningState,omitempty"`
}
