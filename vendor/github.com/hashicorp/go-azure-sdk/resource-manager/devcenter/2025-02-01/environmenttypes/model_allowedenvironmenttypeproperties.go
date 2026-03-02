package environmenttypes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AllowedEnvironmentTypeProperties struct {
	DisplayName       *string            `json:"displayName,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
