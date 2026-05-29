package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationAcceleratorProperties struct {
	Components        *[]ApplicationAcceleratorComponent       `json:"components,omitempty"`
	ProvisioningState *ApplicationAcceleratorProvisioningState `json:"provisioningState,omitempty"`
}
