package edgedevices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ValidateRequest struct {
	AdditionalInfo *string  `json:"additionalInfo,omitempty"`
	EdgeDeviceIds  []string `json:"edgeDeviceIds"`
}
