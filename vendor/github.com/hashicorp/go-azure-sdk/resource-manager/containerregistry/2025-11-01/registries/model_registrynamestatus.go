package registries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistryNameStatus struct {
	AvailableLoginServerName *string `json:"availableLoginServerName,omitempty"`
	Message                  *string `json:"message,omitempty"`
	NameAvailable            *bool   `json:"nameAvailable,omitempty"`
	Reason                   *string `json:"reason,omitempty"`
}
