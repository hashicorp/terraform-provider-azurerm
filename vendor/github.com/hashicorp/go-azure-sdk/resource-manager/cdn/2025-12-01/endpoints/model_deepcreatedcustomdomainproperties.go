package endpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeepCreatedCustomDomainProperties struct {
	HostName       string  `json:"hostName"`
	ValidationData *string `json:"validationData,omitempty"`
}
