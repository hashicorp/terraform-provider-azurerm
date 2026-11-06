package virtualmachineimagetemplate

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProvisioningError struct {
	Message               *string                `json:"message,omitempty"`
	ProvisioningErrorCode *ProvisioningErrorCode `json:"provisioningErrorCode,omitempty"`
}
