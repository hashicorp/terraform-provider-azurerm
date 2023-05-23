package pricings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Extension struct {
	AdditionalExtensionProperties *interface{}     `json:"additionalExtensionProperties,omitempty"`
	IsEnabled                     IsEnabled        `json:"isEnabled"`
	Name                          string           `json:"name"`
	OperationStatus               *OperationStatus `json:"operationStatus,omitempty"`
}
