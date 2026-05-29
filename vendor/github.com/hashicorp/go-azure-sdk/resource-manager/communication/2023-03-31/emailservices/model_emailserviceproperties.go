package emailservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EmailServiceProperties struct {
	DataLocation      string                          `json:"dataLocation"`
	ProvisioningState *EmailServicesProvisioningState `json:"provisioningState,omitempty"`
}
