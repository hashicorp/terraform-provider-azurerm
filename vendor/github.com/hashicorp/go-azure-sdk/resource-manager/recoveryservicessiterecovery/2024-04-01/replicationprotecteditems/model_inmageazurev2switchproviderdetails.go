package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InMageAzureV2SwitchProviderDetails struct {
	TargetApplianceId *string `json:"targetApplianceId,omitempty"`
	TargetFabricId    *string `json:"targetFabricId,omitempty"`
	TargetResourceId  *string `json:"targetResourceId,omitempty"`
	TargetVaultId     *string `json:"targetVaultId,omitempty"`
}
